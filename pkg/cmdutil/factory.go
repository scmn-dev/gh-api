package cmdutil

import (
	"net/http"

	"github.com/scmn-dev/gh-api/context"
	"github.com/scmn-dev/gh-api/core/ghrepo"
	"github.com/scmn-dev/gh-api/pkg/iostreams"

	"github.com/scmn-dev/secman/tools/packages"
	"github.com/scmn-dev/cluster"
)

type Browser interface {
	Browse(string) error
}

type Factory struct {
	IOStreams 		*iostreams.IOStreams
	Browser   		Browser

	HttpClient 		func() (*http.Client, error)
	BaseRepo   		func() (ghrepo.Interface, error)
	Remotes    		func() (context.Remotes, error)
	Cluster     		func() (cluster.Cluster, error)
	Branch     	    func() (string, error)
	PackageManager  packages.PackageManager
	Executable string
}
