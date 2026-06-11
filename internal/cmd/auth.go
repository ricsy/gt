package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

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

var tokenShow bool

func resolveAuthHost() string {
	if loginFlags.host != "" {
		return loginFlags.host
	}
	return config.DefaultHost
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a host",
	Long:  `Login to a git host with a token.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := resolveAuthHost()
		token := loginFlags.token
		username := loginFlags.username

		if token == "" {
			var err error
			token, err = readTokenFromInput(cmd, host)
			if err != nil {
				return err
			}
		}

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
		host := resolveAuthHost()

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
		host := resolveAuthHost()

		a, err := auth.GetAuth(host)
		if err != nil {
			fmt.Printf("Not logged in to %s\n", host)
			return nil
		}
		fmt.Printf("Logged in to %s as %s\n", host, a.User)
		return nil
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Show authentication token",
	Long:  `Show the authentication token for a host. Output is masked unless --show is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := resolveAuthHost()

		token, err := auth.GetToken(host)
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}

		if tokenShow {
			fmt.Println(token)
			return nil
		}

		fmt.Println(maskToken(token))
		return nil
	},
}

func readTokenFromInput(cmd *cobra.Command, host string) (string, error) {
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Enter token for %s: ", host)

	input := cmd.InOrStdin()
	if file, ok := input.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		tokenBytes, err := term.ReadPassword(int(file.Fd()))
		_, _ = fmt.Fprintln(cmd.ErrOrStderr())
		if err != nil {
			return "", fmt.Errorf("failed to read token: %w", err)
		}

		token := strings.TrimSpace(string(tokenBytes))
		if token == "" {
			return "", fmt.Errorf("token cannot be empty")
		}
		return token, nil
	}

	reader := bufio.NewReader(input)
	token, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("failed to read token: %w", err)
	}
	_, _ = fmt.Fprintln(cmd.ErrOrStderr())

	token = strings.TrimSpace(token)
	if token == "" {
		return "", fmt.Errorf("token cannot be empty")
	}
	return token, nil
}

func maskToken(token string) string {
	if token == "" {
		return ""
	}
	if len(token) <= 4 {
		return strings.Repeat("*", len(token))
	}
	return token[:2] + strings.Repeat("*", len(token)-4) + token[len(token)-2:]
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
	tokenCmd.Flags().BoolVar(&tokenShow, "show", false, "Print the full token value")
}
