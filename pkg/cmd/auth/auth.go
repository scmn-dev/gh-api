package auth

import (
	gitCredentialCmd "github.com/secman-team/gh-api/pkg/cmd/auth/gitcredential"
	authLoginCmd "github.com/secman-team/gh-api/pkg/cmd/auth/login"
	authLogoutCmd "github.com/secman-team/gh-api/pkg/cmd/auth/logout"
	authRefreshCmd "github.com/secman-team/gh-api/pkg/cmd/auth/refresh"
	authStatusCmd "github.com/secman-team/gh-api/pkg/cmd/auth/status"
	"github.com/secman-team/gh-api/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Login, logout, and refresh your authentication.",
		Long:  `Manage secman's authentication state.`,
	}

	cmdutil.DisableAuthCheck(cmd)

	cmd.AddCommand(authLoginCmd.NewCmdLogin(f, nil))
	cmd.AddCommand(authLogoutCmd.NewCmdLogout(f, nil))
	cmd.AddCommand(authStatusCmd.NewCmdStatus(f, nil))
	cmd.AddCommand(authRefreshCmd.NewCmdRefresh(f, nil))
	cmd.AddCommand(gitCredentialCmd.NewCmdCredential(f, nil))

	return cmd
}
