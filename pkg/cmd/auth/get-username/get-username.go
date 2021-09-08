package get_username

import (
	"fmt"

	"github.com/scmn-dev/git"
	"github.com/spf13/cobra"
)

func GetUsername() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-username",
		Args:  cobra.ExactArgs(0),
		Short: "Get Your Github Username.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(git_config.GitConfig())
		},
	}

	return cmd
}
