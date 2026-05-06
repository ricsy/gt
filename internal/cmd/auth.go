package cmd

import (
	"fmt"
	"os"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

const defaultHost = config.DefaultHost

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage git authentication",
	Long:  `Manage authentication for git hosts.`,
}

var loginFlags struct {
	host     string
	token    string
	username string
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a host",
	Long:  `Login to a git host with a token.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := loginFlags.host
		if host == "" {
			host = defaultHost
		}

		token := loginFlags.token
		username := loginFlags.username

		// Interactive mode: prompt for token if not provided
		if token == "" {
			fmt.Printf("Enter token for %s: ", host)
			_, err := fmt.Scanln(&token)
			if err != nil {
				return fmt.Errorf("failed to read token: %w", err)
			}
		}

		// If username not provided, use a placeholder
		if username == "" {
			username = "user"
		}

		if err := auth.Login(host, token, username); err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}

		fmt.Printf("Logged in to %s\n", host)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from a host",
	Long:  `Remove authentication for a git host.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := loginFlags.host
		if host == "" {
			host = defaultHost
		}

		if !auth.IsLoggedIn(host) {
			return fmt.Errorf("not logged in to %s", host)
		}

		if err := auth.Logout(host); err != nil {
			return fmt.Errorf("failed to logout: %w", err)
		}

		fmt.Printf("Logged out from %s\n", host)
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status",
	Long:  `Show the current authentication status for a host.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := loginFlags.host
		if host == "" {
			host = defaultHost
		}

		if auth.IsLoggedIn(host) {
			user, err := auth.CurrentUser(host)
			if err != nil {
				fmt.Printf("Logged in to %s\n", host)
			} else {
				fmt.Printf("Logged in to %s as %s\n", host, user)
			}
		} else {
			fmt.Printf("Not logged in to %s\n", host)
		}
		return nil
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Show or set authentication token",
	Long:  `Show or set the authentication token for a host.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := loginFlags.host
		if host == "" {
			host = defaultHost
		}

		token, err := auth.GetToken(host)
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}

		fmt.Println(token)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(tokenCmd)

	loginCmd.Flags().StringVar(&loginFlags.host, "host", "", "Host (default: gitee.com)")
	loginCmd.Flags().StringVarP(&loginFlags.token, "token", "t", "", "Authentication token")
	loginCmd.Flags().StringVarP(&loginFlags.username, "username", "u", "", "Username")

	logoutCmd.Flags().StringVar(&loginFlags.host, "host", "", "Host (default: gitee.com)")

	statusCmd.Flags().StringVar(&loginFlags.host, "host", "", "Host (default: gitee.com)")

	tokenCmd.Flags().StringVar(&loginFlags.host, "host", "", "Host (default: gitee.com)")

	_ = os.Stdin // suppress unused variable warning
}
