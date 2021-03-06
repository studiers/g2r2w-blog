package main

import (
	"context"
	"fmt"
	pb "github.com/studiers/g2r2w-blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"
)

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterBlogServer(grpcServer, newServer())
	const port = 23000
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Panic(err)
	}
	grpcServer.Serve(lis)
}

type comments map[int64]string

type post struct {
	pb.Post
	comments comments
}

type posts map[int64]*post

type blogServer struct {
	pb.UnimplementedBlogServer
	posts posts
	mutex sync.RWMutex
}

func (server *blogServer) ListPosts(_ *pb.Empty, postsServer pb.Blog_ListPostsServer) error {
	server.mutex.RLock()
	defer server.mutex.RUnlock()

	for postId, post := range server.posts {
		log.Printf("Post#%d %v sent.\n", postId, post)
		postsServer.Send(&pb.PostResponse{Id: postId, Post: &post.Post})
	}

	return nil
}

func (server *blogServer) GetPost(_ context.Context, req *pb.GetPostRequest) (*pb.PostResponse, error) {
	server.mutex.RLock()
	defer server.mutex.RUnlock()

	if post, ok := server.posts[req.Id]; ok {
		return &pb.PostResponse{
			Id:   req.Id,
			Post: &post.Post,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "The post id, %d, doesn't seems existing.", req.Id)
}

func (server *blogServer) CreatePost(_ context.Context, req *pb.CreatePostRequest) (*pb.PostResponse, error) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	postId := int64(len(server.posts))
	server.posts[postId] = &post{
		*req.Post,
		make(comments),
	}

	log.Printf("Post#%d %v created\n", postId, req.Post)

	return &pb.PostResponse{
		Id:   postId,
		Post: req.Post,
	}, nil
}

func (server *blogServer) ModifyPost(_ context.Context, req *pb.ModifyPostRequest) (*pb.PostResponse, error) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	if p, ok := server.posts[req.Id]; ok {
		p.Post = *req.Post
		return &pb.PostResponse{
			Id:   req.Id,
			Post: req.Post,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "The post id, %d, doesn't seems existing.", req.Id)
}

func (server *blogServer) CreateComment(_ context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	if p, ok := server.posts[req.PostId]; ok {
		commentId := int64(len(p.comments))
		p.comments[commentId] = req.Comment
		return &pb.CreateCommentResponse{
			Id:      commentId,
			PostId:  req.PostId,
			Comment: req.Comment,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "The post id, %d, doesn't seems existing.", req.PostId)
}

func (server *blogServer) ListComments(req *pb.ListCommentsRequest, commentsServer pb.Blog_ListCommentsServer) error {
	server.mutex.RLock()
	defer server.mutex.RUnlock()

	if p, ok := server.posts[req.PostId]; ok {
		for index, comment := range p.comments {
			commentsServer.Send(&pb.CommentResponse{
				Id:      index,
				Comment: comment,
			})
		}
	}

	return status.Errorf(codes.NotFound, "The post id, %d, doesn't seems existing.", req.PostId)
}

func newServer() *blogServer {
	s := &blogServer{
		posts: make(posts),
		mutex: sync.RWMutex{},
	}
	return s
}
