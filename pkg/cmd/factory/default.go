package factory

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/scmn-dev/gh-api/api"
	"github.com/scmn-dev/gh-api/context"
	"github.com/scmn-dev/gh-api/git"
	"github.com/scmn-dev/secman/cluster"
	"github.com/scmn-dev/gh-api/core/ghrepo"
	"github.com/scmn-dev/gh-api/pkg/cmdutil"
	"github.com/scmn-dev/gh-api/pkg/iostreams"
)

func New() *cmdutil.Factory {
	f := &cmdutil.Factory{
		Cluster:     ClusterFunc(), // No factory dependencies
		Branch:     branchFunc(),   // No factory dependencies
		Executable: executable(),   // No factory dependencies
	}

	f.IOStreams = ioStreams(f)                   // Depends on Cluster
	f.HttpClient = httpClientFunc(f, "x")        // Depends on Cluster, IOStreams, and appVersion
	f.Remotes = remotesFunc(f)                   // Depends on Cluster
	f.BaseRepo = BaseRepoFunc(f)                 // Depends on Remotes
	f.Browser = browser(f)                       // Depends on Cluster, and IOStreams

	return f
}

func BaseRepoFunc(f *cmdutil.Factory) func() (ghrepo.Interface, error) {
	return func() (ghrepo.Interface, error) {
		remotes, err := f.Remotes()
		if err != nil {
			return nil, err
		}

		return remotes[0], nil
	}
}

func SmartBaseRepoFunc(f *cmdutil.Factory) func() (ghrepo.Interface, error) {
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

func remotesFunc(f *cmdutil.Factory) func() (context.Remotes, error) {
	rr := &remoteResolver{
		readRemotes: git.Remotes,
		getCluster:   f.Cluster,
	}

	return rr.Resolver()
}

func httpClientFunc(f *cmdutil.Factory, appVersion string) func() (*http.Client, error) {
	return func() (*http.Client, error) {
		io := f.IOStreams
		cfg, err := f.Cluster()
		if err != nil {
			return nil, err
		}

		return NewHTTPClient(io, cfg, appVersion, true)
	}
}

func browser(f *cmdutil.Factory) cmdutil.Browser {
	io := f.IOStreams
	return cmdutil.NewBrowser(browserLauncher(f), io.Out, io.ErrOut)
}

// Browser precedence
// 1. GH_BROWSER
// 2. browser from Cluster
// 3. BROWSER
func browserLauncher(f *cmdutil.Factory) string {
	if smBrowser := os.Getenv("GH_BROWSER"); smBrowser != "" {
		return smBrowser
	}

	cfg, err := f.Cluster()
	if err == nil {
		if cfgBrowser, _ := cfg.Get("", "browser"); cfgBrowser != "" {
			return cfgBrowser
		}
	}

	return os.Getenv("BROWSER")
}

func executable() string {
	secman := "secman"
	if exe, err := os.Executable(); err == nil {
		secman = exe
	}

	return secman
}

func ClusterFunc() func() (cluster.Cluster, error) {
	var cachedCluster cluster.Cluster
	var ClusterError error
	return func() (cluster.Cluster, error) {
		if cachedCluster != nil || ClusterError != nil {
			return cachedCluster, ClusterError
		}

		cachedCluster, ClusterError = cluster.ParseDefaultCluster()
		if errors.Is(ClusterError, os.ErrNotExist) {
			cachedCluster = cluster.NewBlankCluster()
			ClusterError = nil
		}

		cachedCluster = cluster.InheritEnv(cachedCluster)
		return cachedCluster, ClusterError
	}
}

func branchFunc() func() (string, error) {
	return func() (string, error) {
		currentBranch, err := git.CurrentBranch()
		if err != nil {
			return "", fmt.Errorf("could not determine current branch: %w", err)
		}

		return currentBranch, nil
	}
}

func ioStreams(f *cmdutil.Factory) *iostreams.IOStreams {
	io := iostreams.System()
	cfg, err := f.Cluster()
	if err != nil {
		return io
	}

	if prompt, _ := cfg.Get("", "prompt"); prompt == "disabled" {
		io.SetNeverPrompt(true)
	}

	// Pager precedence
	// 1. GH_PAGER
	// 2. pager from Cluster
	// 3. PAGER
	if ghPager, ghPagerExists := os.LookupEnv("GH_PAGER"); ghPagerExists {
		io.SetPager(ghPager)
	} else if pager, _ := cfg.Get("", "pager"); pager != "" {
		io.SetPager(pager)
	}

	return io
}
