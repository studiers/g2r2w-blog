package main

import (
	"fmt"
	"github.com/studiers/g2r2w-blog/client/cmd"
	pb "github.com/studiers/g2r2w-blog/proto"
	"google.golang.org/grpc"
	"log"
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

	cmd := cmd.NewCmdClient(&client)
	cmd.Execute()
}
