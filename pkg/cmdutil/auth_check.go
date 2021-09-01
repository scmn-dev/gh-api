package cmdutil

import (
	"github.com/spf13/cobra"
	"github.com/scmn-dev/secman/cluster"
)

// TODO can have this set a PersistentPreRun so we don't have to set for all child commands of auth,
// cluster

func DisableAuthCheck(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}

	cmd.Annotations["skipAuthCheck"] = "true"
}

func CheckAuth(cfg cluster.Cluster) bool {
	if cluster.AuthTokenProvidedFromEnv() {
		return true
	}

	hosts, err := cfg.Hosts()
	if err != nil {
		return false
	}

	for _, hostname := range hosts {
		token, _ := cfg.Get(hostname, "oauth_token")
		if token != "" {
			return true
		}
	}

	return false
}

func IsAuthCheckEnabled(cmd *cobra.Command) bool {
	switch cmd.Name() {
		case "help", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
			return false
	}

	for c := cmd; c.Parent() != nil; c = c.Parent() {
		if c.Annotations != nil && c.Annotations["skipAuthCheck"] == "true" {
			return false
		}
	}

	return true
}
