package post

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	pb "github.com/studiers/g2r2w-blog/proto"
	"io"
	"log"
)

func (o *ListOptions) Run(client *pb.BlogClient) {
	for {
		if stream, err := (*client).ListPosts(context.Background(), &pb.Empty{}); err == nil {
			if resp, err := stream.Recv(); err == nil {
				fmt.Printf("Post#%d: %v\n", resp.Id, resp.Post)
			} else if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
}

type ListOptions struct {
	client *pb.BlogClient
}

func NewListOptions(client *pb.BlogClient) *ListOptions {
	return &ListOptions{
		client,
	}
}

func NewCmdList(client *pb.BlogClient) *cobra.Command {
	o := NewListOptions(client)

	cmd := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(client)
		},
	}

	return cmd
}
