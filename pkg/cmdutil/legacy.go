package cmdutil

import (
	"fmt"
	"os"
	"github.com/scmn-dev/cluster"
)

// TODO: consider passing via Factory
// TODO: support per-hostname settings
func DetermineEditor(cf func() (cluster.Cluster, error)) (string, error) {
	editorCommand := os.Getenv("GH_EDITOR")
	if editorCommand == "" {
		cfg, err := cf()
		if err != nil {
			return "", fmt.Errorf("could not read cluster: %w", err)
		}
		editorCommand, _ = cfg.Get("", "editor")
	}

	return editorCommand, nil
}
