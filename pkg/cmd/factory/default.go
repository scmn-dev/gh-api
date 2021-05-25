package factory

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/secman-team/gh-api/git"
	"github.com/secman-team/gh-api/core/config"
	"github.com/secman-team/gh-api/core/ghrepo"
	"github.com/secman-team/gh-api/pkg/cmdutil"
	"github.com/secman-team/gh-api/pkg/iostreams"
)

func New() *cmdutil.Factory {
	io := iostreams.System()

	var cachedConfig config.Config
	var configError error
	configFunc := func() (config.Config, error) {
		if cachedConfig != nil || configError != nil {
			return cachedConfig, configError
		}
		cachedConfig, configError = config.ParseDefaultConfig()
		if errors.Is(configError, os.ErrNotExist) {
			cachedConfig = config.NewBlankConfig()
			configError = nil
		}
		cachedConfig = config.InheritEnv(cachedConfig)
		return cachedConfig, configError
	}

	rr := &remoteResolver{
		readRemotes: git.Remotes,
		getConfig:   configFunc,
	}
	remotesFunc := rr.Resolver()

	ghExecutable := "gh"
	if exe, err := os.Executable(); err == nil {
		ghExecutable = exe
	}

	return &cmdutil.Factory{
		IOStreams: io,
		Config:    configFunc,
		Remotes:   remotesFunc,
		HttpClient: func() (*http.Client, error) {
			cfg, err := configFunc()
			if err != nil {
				return nil, err
			}

			return NewHTTPClient(io, cfg, "x", true), nil
		},
		BaseRepo: func() (ghrepo.Interface, error) {
			remotes, err := remotesFunc()
			if err != nil {
				return nil, err
			}
			return remotes[0], nil
		},
		Branch: func() (string, error) {
			currentBranch, err := git.CurrentBranch()
			if err != nil {
				return "", fmt.Errorf("could not determine current branch: %w", err)
			}
			return currentBranch, nil
		},
		Executable: ghExecutable,
		Browser:    cmdutil.NewBrowser(os.Getenv("BROWSER"), io.Out, io.ErrOut),
	}
}
