package key

import (
	cmdAdd "github.com/scmn-dev/gh-api/pkg/cmd/ssh-key/add"
	cmdList "github.com/scmn-dev/gh-api/pkg/cmd/ssh-key/list"
	"github.com/scmn-dev/gh-api/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSSHKey(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh-key <command>",
		Short: "Manage SSH keys",
		Long:  "Manage SSH keys registered with your GitHub account",
	}

	cmd.AddCommand(cmdList.NewCmdList(f, nil))
	cmd.AddCommand(cmdAdd.NewCmdAdd(f, nil))

	return cmd
}
