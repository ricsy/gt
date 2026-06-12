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

var setupFlags struct {
	Overwrite bool
}

var doctorFlags struct {
	JSON bool
}

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

		authInfo, _ := auth.GetAuth(host)
		if err := auth.Logout(host); err != nil {
			return fmt.Errorf("failed to logout: %w", err)
		}

		if logoutCleanupGit {
			if err := clearGitCredentialsForHost(host, authInfo.User); err != nil {
				return fmt.Errorf("logged out from %s but failed to clear git credentials: %w", host, err)
			}
			fmt.Printf("Logged out from %s and cleared git credentials\n", host)
			return nil
		}

		fmt.Printf("Logged out from %s\n", host)
		return nil
	},
}

var logoutCleanupGit bool

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

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure git credentials for a host",
	Long:  `Write the current host token into git credential storage for HTTPS git operations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := resolveAuthHost()

		authInfo, err := auth.GetAuth(host)
		if err != nil {
			return fmt.Errorf("failed to get stored auth: %w", err)
		}
		if authInfo.User == "" {
			return fmt.Errorf("stored auth for %s does not include a username; run gt auth login --username <name>", host)
		}
		if authInfo.Token == "" {
			return fmt.Errorf("stored auth for %s does not include a token", host)
		}

		if setupFlags.Overwrite {
			if err := clearGitCredentialsForHost(host, authInfo.User); err != nil {
				return fmt.Errorf("failed to clear existing git credentials: %w", err)
			}
		}

		if err := gitCredentialApprove(host, authInfo.User, authInfo.Token); err != nil {
			return err
		}

		fmt.Printf("Configured git credentials for %s as %s\n", host, authInfo.User)
		return nil
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose git authentication for a host",
	Long:  `Inspect stored auth, git credential state, and local git remote access for a host.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		host := resolveAuthHost()
		report := runAuthDoctor(host)

		if doctorFlags.JSON {
			jsonText, err := report.ToJSON()
			if err != nil {
				return fmt.Errorf("failed to encode doctor report: %w", err)
			}
			cmd.Println(jsonText)
			if report.ReadyForHTTPSGit {
				return nil
			}
			return fmt.Errorf("doctor found issues")
		}

		cmd.Printf("Host: %s\n", host)
		if report.StoredAuthErr != "" {
			cmd.Printf("Stored auth: missing (%s)\n", report.StoredAuthErr)
		} else {
			cmd.Printf("Stored auth: ok as %s\n", report.StoredAuth.User)
		}

		if report.CredentialHelperErr != "" {
			cmd.Printf("Git credential helper: unavailable (%s)\n", report.CredentialHelperErr)
		} else if report.CredentialHelper == "" {
			cmd.Printf("Git credential helper: not configured\n")
		} else {
			cmd.Printf("Git credential helper: %s\n", report.CredentialHelper)
		}

		if report.GitCredentialErr != "" {
			cmd.Printf("Git credential lookup: failed (%s)\n", report.GitCredentialErr)
		} else {
			cmd.Printf("Git credential lookup: ok as %s\n", report.GitCredential.User)
		}

		if len(report.CredentialTargets) > 1 {
			cmd.Printf("Credential targets: multiple entries found for %s\n", host)
			for _, target := range report.CredentialTargets {
				cmd.Printf("  - %s (%s)\n", target.Target, target.User)
			}
		}

		if report.StoredAuthErr == "" && report.GitCredentialErr == "" && !report.AuthUserMatchesGit {
			cmd.Printf("Credential mismatch: stored auth user is %s, git credential user is %s\n", report.StoredAuth.User, report.GitCredential.User)
			cmd.Printf("Recommendation: gt auth setup --overwrite\n")
		}

		if report.CurrentDirectoryRepoErr != "" {
			cmd.Printf("Current directory git status: failed (%s)\n", report.CurrentDirectoryRepoErr)
		} else if !report.CurrentDirectoryIsRepo {
			cmd.Printf("Current directory git status: not a git worktree\n")
		} else {
			cmd.Printf("Current directory git status: git worktree\n")
		}

		if report.OriginURL != "" {
			cmd.Printf("Origin URL: %s\n", report.OriginURL)
			if report.OriginAccessErr != "" {
				cmd.Printf("Origin access: failed (%s)\n", report.OriginAccessErr)
			} else if report.OriginAccessOK {
				cmd.Printf("Origin access: ok\n")
			}
		} else if report.OriginURLErr != "" {
			cmd.Printf("Origin URL: unavailable (%s)\n", report.OriginURLErr)
		}

		if report.StoredAuthErr != "" {
			return fmt.Errorf("doctor found issues: run gt auth login --username <name> --token <token>")
		}
		if report.GitCredentialErr != "" || !report.AuthUserMatchesGit {
			return fmt.Errorf("doctor found issues: run gt auth setup --overwrite")
		}

		cmd.Println("Doctor: OK")
		return nil
	},
}

func clearGitCredentialsForHost(host, username string) error {
	if err := gitCredentialReject(host, ""); err != nil {
		return err
	}
	if username != "" {
		if err := gitCredentialReject(host, username); err != nil {
			return err
		}
	}

	targets, err := gitCredentialTargetsList(host)
	if err != nil {
		return err
	}
	for _, target := range targets {
		if err := gitCredentialTargetDelete(target.Target); err != nil {
			return err
		}
	}

	return nil
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
	authCmd.AddCommand(setupCmd)
	authCmd.AddCommand(doctorCmd)

	loginCmd.Flags().StringVar(&loginFlags.host, "host", "", "Host (default: gitee.com)")
	loginCmd.Flags().StringVarP(&loginFlags.token, "token", "t", "", "Authentication token")
	loginCmd.Flags().StringVarP(&loginFlags.username, "username", "u", "", "Username")

	logoutCmd.Flags().StringVar(&loginFlags.host, "host", "", "Host (default: gitee.com)")
	logoutCmd.Flags().BoolVar(&logoutCleanupGit, "cleanup-git", false, "Also clear git credentials for the current host")

	statusCmd.Flags().StringVar(&loginFlags.host, "host", "", "Host (default: gitee.com)")

	tokenCmd.Flags().StringVar(&loginFlags.host, "host", "", "Host (default: gitee.com)")
	tokenCmd.Flags().BoolVar(&tokenShow, "show", false, "Print the full token value")

	setupCmd.Flags().StringVar(&loginFlags.host, "host", "", "Host (default: gitee.com)")
	setupCmd.Flags().BoolVar(&setupFlags.Overwrite, "overwrite", false, "Clear existing git credentials for the host before writing the new one")

	doctorCmd.Flags().StringVar(&loginFlags.host, "host", "", "Host (default: gitee.com)")
	doctorCmd.Flags().BoolVar(&doctorFlags.JSON, "json", false, "Print the doctor report as JSON")
}
