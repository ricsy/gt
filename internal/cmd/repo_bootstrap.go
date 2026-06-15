package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

const (
	cloneURLModeHTTPS = "https"
	cloneURLModeSSH   = "ssh"
)

// validateCloneURLMode validates user-facing clone URL mode flags.
func validateCloneURLMode(mode string) error {
	switch mode {
	case "", cloneURLModeHTTPS, cloneURLModeSSH:
		return nil
	default:
		return fmt.Errorf("invalid clone url mode %q: expected https or ssh", mode)
	}
}

// resolveRepoCloneURL returns the effective repository remote URL for the requested mode.
func resolveRepoCloneURL(host string, repoOwner, repoName, httpsURL, sshURL, mode string) (string, error) {
	if err := validateCloneURLMode(mode); err != nil {
		return "", err
	}

	switch mode {
	case cloneURLModeSSH:
		if sshURL != "" {
			return sshURL, nil
		}
		return config.RepoGitSSHURL(host, repoOwner, repoName), nil
	default:
		if httpsURL != "" {
			return httpsURL, nil
		}
		return config.RepoGitHTTPSURL(host, repoOwner, repoName), nil
	}
}

// repoBootstrapCommand creates a remote repository, wires the current worktree, and optionally pushes it.
func repoBootstrapCommand(cmd *cobra.Command, args []string) error {
	host := resolveCommandHost()
	if err := validateCloneURLMode(repoCreateOpts.CloneURLMode); err != nil {
		return err
	}

	repo, err := createRepoFromCommand(cmd)
	if err != nil {
		return err
	}

	cloneURL, err := resolveRepoCloneURL(host, repo.Owner.Login, repo.Name, repo.CloneURL, repo.SSHURL, repoCreateOpts.CloneURLMode)
	if err != nil {
		return err
	}

	insideWorkTree, err := gitIsInsideWorkTree()
	if err != nil {
		return fmt.Errorf("failed to inspect current git worktree: %w", err)
	}
	if !insideWorkTree {
		return fmt.Errorf("repository created at %s, but current directory is not a git worktree", repo.HTMLURL)
	}

	if originURL, err := gitRemoteGetURL(repoBootstrapOpts.RemoteName); err != nil {
		if err := gitRemoteAdd(repoBootstrapOpts.RemoteName, cloneURL); err != nil {
			return fmt.Errorf("repository created at %s, but failed to add remote %s: %w", repo.HTMLURL, repoBootstrapOpts.RemoteName, err)
		}
	} else if originURL != cloneURL {
		if err := gitRemoteSetURL(repoBootstrapOpts.RemoteName, cloneURL); err != nil {
			return fmt.Errorf("repository created at %s, but failed to update remote %s: %w", repo.HTMLURL, repoBootstrapOpts.RemoteName, err)
		}
	}

	fmt.Printf("Repository created: %s\n", repo.HTMLURL)
	fmt.Printf("Remote %s: %s\n", repoBootstrapOpts.RemoteName, cloneURL)

	if !repoBootstrapOpts.Push {
		cmd.Println("Bootstrap: remote configured, push skipped")
		return nil
	}

	branch, err := gitCurrentBranch()
	if err != nil {
		return fmt.Errorf("repository created at %s, but failed to detect current branch: %w", repo.HTMLURL, err)
	}
	if branch == "" {
		return fmt.Errorf("repository created at %s, but current branch is empty", repo.HTMLURL)
	}

	if repoCreateOpts.CloneURLMode == cloneURLModeHTTPS {
		report := runAuthDoctor(host)
		if report.StoredAuthErr != "" {
			return fmt.Errorf("repository created at %s, but HTTPS auth is not ready: run gt auth login --username <name> --token <token>", repo.HTMLURL)
		}
		if report.GitCredentialErr != "" || !report.AuthUserMatchesGit {
			return fmt.Errorf("repository created at %s, but HTTPS git credentials are not ready: run gt auth setup --overwrite", repo.HTMLURL)
		}
	}

	if err := gitPushUpstream(repoBootstrapOpts.RemoteName, branch); err != nil {
		return fmt.Errorf("repository created at %s, remote configured to %s, but push failed: %w", repo.HTMLURL, cloneURL, err)
	}

	cmd.Printf("Bootstrap: pushed %s to %s/%s\n", branch, repoBootstrapOpts.RemoteName, branch)
	return nil
}
