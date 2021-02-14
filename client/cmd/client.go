package cmd

import (
	"github.com/spf13/cobra"
	"github.com/studiers/g2r2w-blog/client/cmd/post"
	pb "github.com/studiers/g2r2w-blog/proto"
)

func NewCmdClient(client *pb.BlogClient) *cobra.Command {
	cmd := &cobra.Command{
		Use: "client",
	}

	cmd.AddCommand(post.NewCmdPost(client))

	return cmd
}
