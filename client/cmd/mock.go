package cmd

import (
	"context"
	pb "github.com/studiers/g2r2w-blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type MockCommentId int64
type MockComment string
type MockComments map[MockCommentId]MockComment
type mockPost struct {
	*pb.Post
	comments MockComments
}

type MockPostId int64
type MockPosts map[MockPostId]*mockPost

type mockBlogClient struct {
	posts MockPosts
}

func NewMockBlogClient(mockPosts MockPosts) *pb.BlogClient {
	var client pb.BlogClient = &mockBlogClient{
		posts: mockPosts,
	}
	return &client
}

func NewMockPostFromProtobuf(p *pb.Post) *mockPost {
	return &mockPost{
		Post:     p,
		comments: make(MockComments),
	}
}

func NewMockPost(title string, content string) *mockPost {
	return &mockPost{
		Post: &pb.Post{
			Title:   title,
			Content: content,
		},
		comments: make(MockComments),
	}
}

func (m *mockBlogClient) GetPost(_ context.Context, in *pb.GetPostRequest, _ ...grpc.CallOption) (*pb.PostResponse, error) {
	if post, ok := m.posts[MockPostId(in.Id)]; ok {
		return &pb.PostResponse{
			Id:   in.Id,
			Post: post.Post,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "The mockPost id %d doesn't exist.", in.Id)
}

func (m *mockBlogClient) CreatePost(_ context.Context, in *pb.CreatePostRequest, _ ...grpc.CallOption) (*pb.PostResponse, error) {
	postId := MockPostId(len(m.posts))

	m.posts[postId] = NewMockPostFromProtobuf(in.Post)

	return &pb.PostResponse{
		Id:   int64(postId),
		Post: in.Post,
	}, nil
}

func (m *mockBlogClient) ModifyPost(ctx context.Context, in *pb.ModifyPostRequest, opts ...grpc.CallOption) (*pb.PostResponse, error) {
	if p, ok := m.posts[MockPostId(in.Id)]; ok {
		p.Post = in.Post
		return &pb.PostResponse{
			Id:   in.Id,
			Post: in.Post,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "The mockPost id %d doesn't exist.", in.Id)
}

func (m *mockBlogClient) CreateComment(ctx context.Context, in *pb.CreateCommentRequest, opts ...grpc.CallOption) (*pb.CreateCommentResponse, error) {
	if p, ok := m.posts[MockPostId(in.PostId)]; ok {
		commentId := MockCommentId(len(p.comments))
		p.comments[commentId] = MockComment(in.Comment)
		return &pb.CreateCommentResponse{
			Id:      int64(commentId),
			PostId:  in.PostId,
			Comment: in.Comment,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "The mockPost id %d doesn't exist.", in.PostId)
}

type mockBlogListCommentsClient struct {
	grpc.ClientStream

	current  int64
	comments MockComments
}

func (m mockBlogListCommentsClient) Recv() (*pb.CommentResponse, error) {
	m.current += 1

	if m.current >= int64(len(m.comments)) {
		return nil, io.EOF
	}

	return &pb.CommentResponse{
		Id:      m.current,
		Comment: string(m.comments[MockCommentId(m.current)]),
	}, nil
}

func newMockBlogListCommentsClient(comments MockComments) pb.Blog_ListCommentsClient {
	return mockBlogListCommentsClient{
		ClientStream: nil,
		current:      -1,
		comments:     comments,
	}
}

type mockBlogListPostsClient struct {
	grpc.ClientStream

	current int64
	posts   MockPosts
}

func (m mockBlogListPostsClient) Recv() (*pb.PostResponse, error) {
	m.current += 1

	if m.current >= int64(len(m.posts)) {
		return nil, io.EOF
	}

	post := m.posts[MockPostId(m.current)]
	return &pb.PostResponse{
		Id:   m.current,
		Post: post.Post,
	}, nil
}

func NewMockBlogListPostsClient(posts MockPosts) pb.Blog_ListPostsClient {
	return mockBlogListPostsClient{
		ClientStream: nil,
		current:      -1,
		posts:        posts,
	}
}

func (m *mockBlogClient) ListComments(ctx context.Context, in *pb.ListCommentsRequest, opts ...grpc.CallOption) (pb.Blog_ListCommentsClient, error) {
	if p, ok := m.posts[MockPostId(in.PostId)]; ok {
		return newMockBlogListCommentsClient(p.comments), nil
	}

	return nil, status.Errorf(codes.NotFound, "The mockPost id %d doesn't exist.", in.PostId)
}

func (m *mockBlogClient) ListPosts(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (pb.Blog_ListPostsClient, error) {
	return NewMockBlogListPostsClient(m.posts), nil
}
