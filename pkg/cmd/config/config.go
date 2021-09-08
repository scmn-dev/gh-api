package Cluster

import (
	"fmt"
	"strings"

	"github.com/scmn-dev/cluster"
	cmdGet "github.com/scmn-dev/gh-api/pkg/cmd/config/get"
	cmdSet "github.com/scmn-dev/gh-api/pkg/cmd/config/set"
	"github.com/scmn-dev/gh-api/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdCluster(f *cmdutil.Factory) *cobra.Command {
	longDoc := strings.Builder{}
	longDoc.WriteString("Display or change Clusteruration settings for secman.\n\n")
	longDoc.WriteString("Current respected settings:\n")
	for _, co := range cluster.ClusterOptions() {
		longDoc.WriteString(fmt.Sprintf("- %s: %s", co.Key, co.Description))
		if co.DefaultValue != "" {
			longDoc.WriteString(fmt.Sprintf(" (default: %q)", co.DefaultValue))
		}
		longDoc.WriteRune('\n')
	}

	cmd := &cobra.Command{
		Use:   "gh-Cluster <command>",
		Short: "Manage Clusteruration of github for secman.",
		Long:  longDoc.String(),
	}

	cmdutil.DisableAuthCheck(cmd)

	cmd.AddCommand(cmdGet.NewCmdClusterGet(f, nil))
	cmd.AddCommand(cmdSet.NewCmdClusterSet(f, nil))

	return cmd
}
