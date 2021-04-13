package root

import (
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/secman-team/gh-api/api"
	"github.com/secman-team/gh-api/context"
	"github.com/secman-team/gh-api/internal/ghrepo"
	actionsCmd "github.com/secman-team/gh-api/pkg/cmd/actions"
	aliasCmd "github.com/secman-team/gh-api/pkg/cmd/alias"
	apiCmd "github.com/secman-team/gh-api/pkg/cmd/api"
	authCmd "github.com/secman-team/gh-api/pkg/cmd/auth"
	completionCmd "github.com/secman-team/gh-api/pkg/cmd/completion"
	configCmd "github.com/secman-team/gh-api/pkg/cmd/config"
	"github.com/secman-team/gh-api/pkg/cmd/factory"
	issueCmd "github.com/secman-team/gh-api/pkg/cmd/issue"
	prCmd "github.com/secman-team/gh-api/pkg/cmd/pr"
	releaseCmd "github.com/secman-team/gh-api/pkg/cmd/release"
	repoCmd "github.com/secman-team/gh-api/pkg/cmd/repo"
	creditsCmd "github.com/secman-team/gh-api/pkg/cmd/repo/credits"
	runCmd "github.com/secman-team/gh-api/pkg/cmd/run"
	versionCmd "github.com/secman-team/gh-api/pkg/cmd/version"
	workflowCmd "github.com/secman-team/gh-api/pkg/cmd/workflow"
	"github.com/secman-team/gh-api/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdRoot(f *cmdutil.Factory, version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gh <command> <subcommand> [flags]",
		Short: "GitHub CLI",
		Long:  `Work seamlessly with GitHub from the command line.`,

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ gh repo clone secman-team/gh-api
		`),
		Annotations: map[string]string{
			"help:feedback": heredoc.Doc(`
				Open an issue using 'gh issue create -R github.com/secman-team/gh-api'
			`),
		},
	}

	cmd.SetOut(f.IOStreams.Out)
	cmd.SetErr(f.IOStreams.ErrOut)

	cs := f.IOStreams.ColorScheme()

	helpHelper := func(command *cobra.Command, args []string) {
		rootHelpFunc(cs, command, args)
	}

	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.SetHelpFunc(helpHelper)
	cmd.SetUsageFunc(rootUsageFunc)
	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	formattedVersion := versionCmd.Format(version, buildDate)
	cmd.SetVersionTemplate(formattedVersion)
	cmd.Version = formattedVersion
	cmd.Flags().Bool("version", false, "Show gh version")

	// Child commands
	cmd.AddCommand(versionCmd.NewCmdVersion(f, version, buildDate))
	cmd.AddCommand(aliasCmd.NewCmdAlias(f))
	cmd.AddCommand(authCmd.NewCmdAuth(f))
	cmd.AddCommand(configCmd.NewCmdConfig(f))
	cmd.AddCommand(creditsCmd.NewCmdCredits(f, nil))
	cmd.AddCommand(completionCmd.NewCmdCompletion(f.IOStreams))

	cmd.AddCommand(actionsCmd.NewCmdActions(f))
	cmd.AddCommand(runCmd.NewCmdRun(f))
	cmd.AddCommand(workflowCmd.NewCmdWorkflow(f))

	// the `api` command should not inherit any extra HTTP headers
	bareHTTPCmdFactory := *f
	bareHTTPCmdFactory.HttpClient = bareHTTPClient(f, version)

	cmd.AddCommand(apiCmd.NewCmdApi(&bareHTTPCmdFactory, nil))

	// below here at the commands that require the "intelligent" BaseRepo resolver
	repoResolvingCmdFactory := *f
	repoResolvingCmdFactory.BaseRepo = resolvedBaseRepo(f)

	cmd.AddCommand(prCmd.NewCmdPR(&repoResolvingCmdFactory))
	cmd.AddCommand(issueCmd.NewCmdIssue(&repoResolvingCmdFactory))
	cmd.AddCommand(releaseCmd.NewCmdRelease(&repoResolvingCmdFactory))
	cmd.AddCommand(repoCmd.NewCmdRepo(&repoResolvingCmdFactory))

	// Help topics
	cmd.AddCommand(NewHelpTopic("environment"))
	referenceCmd := NewHelpTopic("reference")
	referenceCmd.SetHelpFunc(referenceHelpFn(f.IOStreams))
	cmd.AddCommand(referenceCmd)

	cmdutil.DisableAuthCheck(cmd)

	// this needs to appear last:
	referenceCmd.Long = referenceLong(cmd)
	return cmd
}

func bareHTTPClient(f *cmdutil.Factory, version string) func() (*http.Client, error) {
	return func() (*http.Client, error) {
		cfg, err := f.Config()
		if err != nil {
			return nil, err
		}
		return factory.NewHTTPClient(f.IOStreams, cfg, version, false), nil
	}
}

func resolvedBaseRepo(f *cmdutil.Factory) func() (ghrepo.Interface, error) {
	return func() (ghrepo.Interface, error) {
		httpClient, err := f.HttpClient()
		if err != nil {
			return nil, err
		}

		apiClient := api.NewClientFromHTTP(httpClient)

		remotes, err := f.Remotes()
		if err != nil {
			return nil, err
		}
		repoContext, err := context.ResolveRemotesToRepos(remotes, apiClient, "")
		if err != nil {
			return nil, err
		}
		baseRepo, err := repoContext.BaseRepo(f.IOStreams)
		if err != nil {
			return nil, err
		}

		return baseRepo, nil
	}
}
