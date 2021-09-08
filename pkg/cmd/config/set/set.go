package set

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/scmn-dev/cluster"
	"github.com/scmn-dev/gh-api/pkg/cmdutil"
	"github.com/scmn-dev/gh-api/pkg/iostreams"
	"github.com/spf13/cobra"
)

type SetOptions struct {
	IO     *iostreams.IOStreams
	Cluster cluster.Cluster

	Key      string
	Value    string
	Hostname string
}

func NewCmdClusterSet(f *cmdutil.Factory, runF func(*SetOptions) error) *cobra.Command {
	opts := &SetOptions{
		IO: f.IOStreams,
	}

	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Update Clusteruration with a value for the given key",
		Example: heredoc.Doc(`
			secman Cluster set editor vim
			secman Cluster set editor "code --wait"
			secman Cluster set git_protocol ssh --host github.com
			secman Cluster set prompt disabled
		`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			Cluster, err := f.Cluster()
			if err != nil {
				return err
			}
			opts.Cluster = Cluster
			opts.Key = args[0]
			opts.Value = args[1]

			if runF != nil {
				return runF(opts)
			}

			return setRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Hostname, "host", "", "", "Set per-host setting")

	return cmd
}

func setRun(opts *SetOptions) error {
	err := cluster.ValidateKey(opts.Key)
	if err != nil {
		warningIcon := opts.IO.ColorScheme().WarningIcon()
		fmt.Fprintf(opts.IO.ErrOut, "%s warning: '%s' is not a known cluster key\n", warningIcon, opts.Key)
	}

	err = cluster.ValidateValue(opts.Key, opts.Value)
	if err != nil {
		var invalidValue *cluster.InvalidValueError
		if errors.As(err, &invalidValue) {
			var values []string
			for _, v := range invalidValue.ValidValues {
				values = append(values, fmt.Sprintf("'%s'", v))
			}
			return fmt.Errorf("failed to set %q to %q: valid values are %v", opts.Key, opts.Value, strings.Join(values, ", "))
		}
	}

	err = opts.Cluster.Set(opts.Hostname, opts.Key, opts.Value)
	if err != nil {
		return fmt.Errorf("failed to set %q to %q: %w", opts.Key, opts.Value, err)
	}

	err = opts.Cluster.Write()
	if err != nil {
		return fmt.Errorf("failed to write Cluster to disk: %w", err)
	}
	return nil
}
