package cmdutil

import (
	"net/http"

	"github.com/scmn-dev/gh-api/context"
	"github.com/scmn-dev/gh-api/core/config"
	"github.com/scmn-dev/gh-api/core/ghrepo"
	"github.com/scmn-dev/gh-api/pkg/iostreams"
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

	// ExtensionManager extensions.ExtensionManager

	// Executable is the path to the currently invoked gh binary
	Executable string
}
