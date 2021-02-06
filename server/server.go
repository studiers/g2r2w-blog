package main

import (
	"fmt"
	pb "github.com/studiers/g2r2w-blog/proto"
	"google.golang.org/grpc"
	"log"
	"net"
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

type blogServer struct {
	pb.UnimplementedBlogServer
}

func (server *blogServer) GetPosts(_ *pb.Empty, postsServer pb.Blog_GetPostsServer) error {
	postsServer.Send(&pb.Post{Title: "Hello"})
	return nil
}

func newServer() *blogServer {
	s := &blogServer{}
	return s
}