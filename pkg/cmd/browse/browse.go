package browse

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/secman-team/gh-api/api"
	"github.com/secman-team/gh-api/core/ghrepo"
	"github.com/secman-team/gh-api/pkg/cmdutil"
	"github.com/secman-team/gh-api/pkg/iostreams"
	"github.com/spf13/cobra"
)

type browser interface {
	Browse(string) error
}

type BrowseOptions struct {
	BaseRepo   func() (ghrepo.Interface, error)
	Browser    browser
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams

	SelectorArg string

	Branch       string
	ProjectsFlag bool
	RepoFlag     bool
	SettingsFlag bool
	WikiFlag     bool
}

func NewCmdBrowse(f *cmdutil.Factory, runF func(*BrowseOptions) error) *cobra.Command {
	opts := &BrowseOptions{
		Browser:    f.Browser,
		HttpClient: f.HttpClient,
		IO:         f.IOStreams,
	}

	cmd := &cobra.Command{
		Long:  "Open the GitHub repository in the web browser.",
		Short: "Open the repository in the browser",
		Use:   "browse [<number> | <path>]",
		Args:  cobra.MaximumNArgs(1),
		Example: heredoc.Doc(`
			$ secman browse
			#=> Open the home page of the current repository
			$ secman browse 61
			#=> Open issue or pull request 61
			$ secman browse --settings
			#=> Open repository settings
			$ secman browse main.go:312
			#=> Open main.go at line 312
			$ secman browse main.go --branch main
			#=> Open main.go in the main branch
		`),
		Annotations: map[string]string{
			"IsCore": "true",
			"help:arguments": heredoc.Doc(`
				A browser location can be specified using arguments in the following format:
				- by number for issue or pull request, e.g. "123"; or
				- by path for opening folders and files, e.g. "pkg/upgrade/upgrade.go"
			`),
			"help:environment": heredoc.Doc(`
				To configure a web browser other than the default, use the BROWSER environment variable.
			`),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.BaseRepo = f.BaseRepo

			if len(args) > 0 {
				opts.SelectorArg = args[0]
			}

			if err := cmdutil.MutuallyExclusive("cannot use --projects with --settings", opts.ProjectsFlag, opts.SettingsFlag); err != nil {
				return err
			}

			if err := cmdutil.MutuallyExclusive("cannot use --projects with --wiki", opts.ProjectsFlag, opts.WikiFlag); err != nil {
				return err
			}

			if err := cmdutil.MutuallyExclusive("cannot use --projects with --branch", opts.ProjectsFlag, opts.Branch != ""); err != nil {
				return err
			}

			if err := cmdutil.MutuallyExclusive("cannot use --settings with --wiki", opts.SettingsFlag, opts.WikiFlag); err != nil {
				return err
			}

			if err := cmdutil.MutuallyExclusive("cannot use --settings with --branch", opts.SettingsFlag, opts.Branch != ""); err != nil {
				return err
			}

			if err := cmdutil.MutuallyExclusive("cannot use --wiki with --branch", opts.WikiFlag, opts.Branch != ""); err != nil {
				return err
			}

			if runF != nil {
				return runF(opts)
			}

			return runBrowse(opts)
		},
	}

	cmdutil.EnableRepoOverride(cmd, f)
	cmd.Flags().BoolVarP(&opts.ProjectsFlag, "projects", "p", false, "Open repository projects")
	cmd.Flags().BoolVarP(&opts.WikiFlag, "wiki", "w", false, "Open repository wiki")
	cmd.Flags().BoolVarP(&opts.SettingsFlag, "settings", "s", false, "Open repository settings")
	cmd.Flags().StringVarP(&opts.Branch, "branch", "b", "", "Select another branch by passing in the branch name")

	return cmd
}

func runBrowse(opts *BrowseOptions) error {
	baseRepo, err := opts.BaseRepo()
	if err != nil {
		return fmt.Errorf("unable to determine base repository: %w\nUse 'secman browse --help' for more information about browse\n", err)
	}

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("unable to create an http client: %w\nUse 'secman browse --help' for more information about browse\n", err)
	}

	url := ghrepo.GenerateRepoURL(baseRepo, "")

	if opts.ProjectsFlag {
		err := opts.Browser.Browse(url + "/projects")
		return err
	}

	if opts.SettingsFlag {
		err := opts.Browser.Browse(url + "/settings")
		return err
	}

	if opts.WikiFlag {
		err := opts.Browser.Browse(url + "/wiki")
		return err
	}

	if isNumber(opts.SelectorArg) {
		url += "/issues/" + opts.SelectorArg
		err := opts.Browser.Browse(url)
		return err
	}

	if opts.Branch != "" {
		url += "/tree/" + opts.Branch + "/"
	} else {
		apiClient := api.NewClientFromHTTP(httpClient)
		branchName, err := api.RepoDefaultBranch(apiClient, baseRepo)
		if err != nil {
			return err
		}

		url += "/tree/" + branchName + "/"
	}

	if opts.SelectorArg != "" {
		arr, err := parseFileArg(opts.SelectorArg)
		if err != nil {
			return err
		}

		if len(arr) > 1 {
			url += arr[0] + "#L" + arr[1]
		} else {
			url += arr[0]
		}
	}

	err = opts.Browser.Browse(url)
	if opts.IO.IsStdoutTTY() && err == nil {
		fmt.Fprintf(opts.IO.Out, "now opening %s in browser\n", url)
	}

	return err
}

func parseFileArg(fileArg string) ([]string, error) {
	arr := strings.Split(fileArg, ":")
	if len(arr) > 2 {
		return arr, fmt.Errorf("invalid use of colon\nUse 'secman browse --help' for more information about browse\n")
	}

	if len(arr) > 1 && !isNumber(arr[1]) {
		return arr, fmt.Errorf("invalid line number after colon\nUse 'secman browse --help' for more information about browse\n")
	}

	return arr, nil
}

func isNumber(arg string) bool {
	_, err := strconv.Atoi(arg)
	return err == nil
}
