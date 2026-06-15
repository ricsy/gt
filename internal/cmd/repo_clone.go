package cmd

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

var repoCloneCmd = &cobra.Command{
	Use:   "clone <repo> [directory]",
	Short: "Clone a repository",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  repoCloneCommand,
}

func repoCloneCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := ResolveRepo(args[0])
	if err != nil {
		return err
	}

	var directory string
	if len(args) > 1 {
		directory = args[1]
	}

	cloneURL := config.RepoGitHTTPSURL(resolveCommandHost(), owner, repoName)
	if authenticatedURL, err := buildAuthenticatedCloneURL(resolveCommandHost(), cloneURL); err == nil && authenticatedURL != "" {
		cloneURL = authenticatedURL
	}

	var gitArgs []string
	gitArgs = append(gitArgs, "clone", cloneURL)
	if directory != "" {
		gitArgs = append(gitArgs, directory)
	}

	gitExec := exec.Command("git", gitArgs...)
	gitExec.Stdout = os.Stdout
	gitExec.Stderr = os.Stderr

	err = gitExec.Run()
	if err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

func buildAuthenticatedCloneURL(host, cloneURL string) (string, error) {
	authInfo, err := auth.GetAuth(host)
	if err != nil || authInfo.Token == "" || authInfo.User == "" {
		return "", err
	}

	parsed, err := url.Parse(cloneURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse clone URL: %w", err)
	}
	parsed.User = url.UserPassword(authInfo.User, authInfo.Token)
	return parsed.String(), nil
}
