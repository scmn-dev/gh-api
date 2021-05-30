package repo

import (
	// "os"
	// "fmt"
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
	// "github.com/abdfnx/git_config"
	// "github.com/secman-team/gh-api/pkg/cmd/factory"
	"github.com/secman-team/gh-api/pkg/iostreams"
)

// type ColorScheme struct {
// 	IO *iostreams.IOStreams
// }

// func opts(f *cmdutil.Factory) ColorScheme {
// 	opts := ColorScheme{
// 		IO: f.IOStreams,
// 	}

// 	return opts
// }

// var cs = opts(factory.New()).IO.ColorScheme()

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
	}

	cmd.AddCommand(repoViewCmd.NewCmdView(f, nil))
	cmd.AddCommand(repoForkCmd.NewCmdFork(f, nil))
	cmd.AddCommand(repoCloneCmd.NewCmdClone(f, nil))
	cmd.AddCommand(repoCreateCmd.NewCmdCreate(f, nil))
	cmd.AddCommand(repoListCmd.NewCmdList(f, nil))
	cmd.AddCommand(creditsCmd.NewCmdRepoCredits(f, nil))
	cmd.AddCommand(gardenCmd.NewCmdGarden(f, nil))

	// username := git_config.GitConfig()
	// if username == ":username" {
	// 	fmt.Println("You're not authenticated, to authenticate run " + cs.Bold("secman auth login"))

	// 	os.Exit(0)
	// }

	return cmd
}
