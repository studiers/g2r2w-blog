package post

import (
	"github.com/spf13/cobra"
	pb "github.com/studiers/g2r2w-blog/proto"
)

func NewCmdPost(client *pb.BlogClient) *cobra.Command {
	cmd := &cobra.Command{
		Use: "post",
	}

	cmd.AddCommand(NewCmdCreate(client))
	cmd.AddCommand(NewCmdList(client))

	return cmd
}
