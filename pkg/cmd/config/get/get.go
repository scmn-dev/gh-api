package get

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/scmn-dev/secman/cluster"
	"github.com/scmn-dev/gh-api/pkg/cmdutil"
	"github.com/scmn-dev/gh-api/pkg/iostreams"
	"github.com/spf13/cobra"
)

type GetOptions struct {
	IO     *iostreams.IOStreams
	Cluster cluster.Cluster

	Hostname string
	Key      string
}

func NewCmdClusterGet(f *cmdutil.Factory, runF func(*GetOptions) error) *cobra.Command {
	opts := &GetOptions{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Print the value of a given cluster key",
		Example: heredoc.Doc(`
			secman cluster get git_protocol
			https
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			Cluster, err := f.Cluster()
			if err != nil {
				return err
			}

			opts.Cluster = Cluster
			opts.Key = args[0]

			if runF != nil {
				return runF(opts)
			}

			return getRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "host", "", "", "Get per-host setting")

	return cmd
}

func getRun(opts *GetOptions) error {
	val, err := opts.Cluster.Get(opts.Hostname, opts.Key)
	if err != nil {
		return err
	}

	if val != "" {
		fmt.Fprintf(opts.IO.Out, "%s\n", val)
	}

	return nil
}
