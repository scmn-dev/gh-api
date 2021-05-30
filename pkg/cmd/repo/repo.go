package repo

import (
	"github.com/MakeNowJust/heredoc"
	repoCloneCmd "github.com/secman-team/gh-api/pkg/cmd/repo/clone"
	repoCreateCmd "github.com/secman-team/gh-api/pkg/cmd/repo/create"
	creditsCmd "github.com/secman-team/gh-api/pkg/cmd/repo/credits"
	repoForkCmd "github.com/secman-team/gh-api/pkg/cmd/repo/fork"
	gardenCmd "github.com/secman-team/gh-api/pkg/cmd/repo/garden"
	repoListCmd "github.com/secman-team/gh-api/pkg/cmd/repo/list"
	repoViewCmd "github.com/secman-team/gh-api/pkg/cmd/repo/view"
	"github.com/secman-team/gh-api/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/abdfnx/git_config"
	"github.com/secman-team/secman/tools/shared"
)

func NewCmdRepo(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo <command>",
		Short: "Create, clone, fork, and view repositories.",
		Long:  `Work with GitHub repositories`,
		Example: heredoc.Doc(`
			secman repo create
			secman repo clone secman-team/gh-api
		`),
		Annotations: map[string]string{
			"help:arguments": heredoc.Doc(`
				A repository can be supplied as an argument in any of the following formats:
				- "OWNER/REPO"
				- by URL, e.g. "https://github.com/OWNER/REPO"
			`),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			username := git_config.GitConfig()

			if username == ":username" {
				shared.AuthMessage()

			} else {
				cmd.Help()
			}

			return nil
		},
	}

	cmd.AddCommand(repoViewCmd.NewCmdView(f, nil))
	cmd.AddCommand(repoForkCmd.NewCmdFork(f, nil))
	cmd.AddCommand(repoCloneCmd.NewCmdClone(f, nil))
	cmd.AddCommand(repoCreateCmd.NewCmdCreate(f, nil))
	cmd.AddCommand(repoListCmd.NewCmdList(f, nil))
	cmd.AddCommand(creditsCmd.NewCmdRepoCredits(f, nil))
	cmd.AddCommand(gardenCmd.NewCmdGarden(f, nil))

	return cmd
}
