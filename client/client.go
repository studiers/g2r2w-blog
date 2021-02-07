package main

import (
	"context"
	"fmt"
	pb "github.com/studiers/g2r2w-blog/proto"
	grpc "google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	const port = 23000
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
	if err != nil {
		log.Panic(err)
	}
	client := pb.NewBlogClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}

	printPosts(client, ctx)
	printPost(client, ctx)
}

func printPosts(client pb.BlogClient, ctx context.Context) {
	stream, err := client.GetPosts(ctx, &pb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	for {
		post, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		log.Println(post)
	}
}

func printPost(client pb.BlogClient, ctx context.Context) {
	post, err := client.GetPost(ctx, &pb.PostIndex{Index: 1})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(post)
}
