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
)

func NewCmdRepo(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo <command>",
		Short: "Create, clone, fork, and view repositories.",
		Long:  `Work with GitHub repositories`,
		Example: heredoc.Doc(`
			secman repo create
			secman repo clone secman-team/gh-api
			secman repo view --web
		`),
		Annotations: map[string]string{
			"IsCore": "true",
			"help:arguments": heredoc.Doc(`
				A repository can be supplied as an argument in any of the following formats:
				- "OWNER/REPO"
				- by URL, e.g. "https://github.com/OWNER/REPO"
			`),
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
