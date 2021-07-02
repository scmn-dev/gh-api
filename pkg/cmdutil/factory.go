package cmdutil

import (
	"net/http"

	"github.com/secman-team/gh-api/context"
	"github.com/secman-team/gh-api/core/config"
	"github.com/secman-team/gh-api/core/ghrepo"
	"github.com/secman-team/gh-api/pkg/iostreams"
	"github.com/secman-team/gh-api/pkg/extensions"
)

type Browser interface {
	Browse(string) error
}

type Factory struct {
	IOStreams *iostreams.IOStreams
	Browser   Browser

	HttpClient func() (*http.Client, error)
	BaseRepo   func() (ghrepo.Interface, error)
	Remotes    func() (context.Remotes, error)
	Config     func() (config.Config, error)
	Branch     func() (string, error)

	ExtensionManager extensions.ExtensionManager

	// Executable is the path to the currently invoked gh binary
	Executable string
}
