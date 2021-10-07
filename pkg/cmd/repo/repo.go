package repo

import (
	"github.com/MakeNowJust/heredoc"
	repoCloneCmd "github.com/gepis/sm-gh-api/pkg/cmd/repo/clone"
	repoCreateCmd "github.com/gepis/sm-gh-api/pkg/cmd/repo/create"
	creditsCmd "github.com/gepis/sm-gh-api/pkg/cmd/repo/credits"
	repoForkCmd "github.com/gepis/sm-gh-api/pkg/cmd/repo/fork"
	gardenCmd "github.com/gepis/sm-gh-api/pkg/cmd/repo/garden"
	repoListCmd "github.com/gepis/sm-gh-api/pkg/cmd/repo/list"
	repoSyncCmd "github.com/gepis/sm-gh-api/pkg/cmd/repo/sync"
	repoViewCmd "github.com/gepis/sm-gh-api/pkg/cmd/repo/view"
	repoBrowseCmd "github.com/gepis/sm-gh-api/pkg/cmd/repo/browse"
	"github.com/gepis/sm-gh-api/pkg/cmdutil"
	"github.com/spf13/cobra"
	git_config "github.com/gepis/git"
	"github.com/scmn-dev/secman/tools/shared"

	"github.com/gepis/sm-gh-api/pkg/cmd/factory"
)

func NewCmdRepo(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo <command>",
		Short: "Create, clone, fork, and view repositories.",
		Long:  `Work with GitHub repositories`,
		Example: heredoc.Doc(`
			secman repo create
			secman repo clone scmn-dev/gh-api
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

	repoResolvingCmdFactory := *f
	repoResolvingCmdFactory.BaseRepo = factory.SmartBaseRepoFunc(f)

	cmd.AddCommand(repoViewCmd.NewCmdView(f, nil))
	cmd.AddCommand(repoForkCmd.NewCmdFork(f, nil))
	cmd.AddCommand(repoCloneCmd.NewCmdClone(f, nil))
	cmd.AddCommand(repoCreateCmd.NewCmdCreate(f, nil))
	cmd.AddCommand(repoListCmd.NewCmdList(f, nil))
	cmd.AddCommand(repoSyncCmd.NewCmdSync(f, nil))
	cmd.AddCommand(creditsCmd.NewCmdRepoCredits(f, nil))
	cmd.AddCommand(gardenCmd.NewCmdGarden(f, nil))
	cmd.AddCommand(repoBrowseCmd.NewCmdBrowse(&repoResolvingCmdFactory, nil))

	return cmd
}
