package post

import (
	"context"
	"github.com/spf13/cobra"
	pb "github.com/studiers/g2r2w-blog/proto"
)

func (o *CreateOptions) Run(client *pb.BlogClient, args []string) {
	title, content := args[0], args[1]
	(*client).CreatePost(context.Background(), &pb.CreatePostRequest{
		Post: &pb.Post{
			Title:   title,
			Content: content,
		},
	})
}

type CreateOptions struct {
	client *pb.BlogClient
}

func NewCreateOptions(client *pb.BlogClient) *CreateOptions {
	return &CreateOptions{
		client,
	}
}

func NewCmdCreate(client *pb.BlogClient) *cobra.Command {
	o := NewCreateOptions(client)

	cmd := &cobra.Command{
		Use:  "create TITLE CONTENT",
		Args: cobra.ExactValidArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(client, args)
		},
	}

	return cmd
}
