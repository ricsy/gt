package cmd

import (
	"strings"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

func printRepoCreatePushDiagnostics(cmd *cobra.Command, repo *api.Repository) {
	cmd.Println("Push diagnostic:")

	host := resolveCommandHost()
	report := runAuthDoctor(host)
	if report.StoredAuthErr != nil {
		cmd.Printf("- Missing stored auth for %s. Run: gt auth login --username <name> --token <token>\n", host)
		return
	}
	if report.GitCredentialErr != nil {
		cmd.Printf("- Git credentials are not ready for HTTPS git operations. Run: gt auth setup\n")
	} else {
		cmd.Printf("- Git credentials are ready for HTTPS git operations.\n")
	}

	insideWorkTree, err := gitIsInsideWorkTree()
	if err != nil {
		cmd.Printf("- Could not inspect the current git worktree: %v\n", err)
		return
	}
	if !insideWorkTree {
		cmd.Printf("- Current directory is not a git worktree.\n")
		cmd.Printf("- Add the remote and push later:\n")
		cmd.Printf("  git remote add origin %s\n", repo.CloneURL)
		cmd.Printf("  git push -u origin master\n")
		return
	}

	originURL, err := gitRemoteGetURL("origin")
	if err != nil {
		cmd.Printf("- No origin remote found.\n")
		cmd.Printf("- Add the remote and push:\n")
		cmd.Printf("  git remote add origin %s\n", repo.CloneURL)
		cmd.Printf("  git push -u origin master\n")
		return
	}

	if originURL != repo.CloneURL {
		cmd.Printf("- origin points to %s\n", originURL)
		cmd.Printf("- Expected remote URL: %s\n", repo.CloneURL)
		cmd.Printf("- Update it before push:\n")
		cmd.Printf("  git remote set-url origin %s\n", repo.CloneURL)
	}

	if report.GitCredentialErr != nil {
		cmd.Printf("- HTTPS remote access was not checked because git credentials are not configured yet.\n")
		return
	}

	if isSSHRemote(originURL) {
		cmd.Printf("- origin uses SSH; verify SSH access separately.\n")
		return
	}

	if err := gitLsRemote(repo.CloneURL); err != nil {
		cmd.Printf("- HTTPS remote access check failed: %v\n", err)
		cmd.Printf("- Run: gt auth setup\n")
		return
	}

	cmd.Printf("- HTTPS remote access check passed for %s\n", repo.CloneURL)
}

func isSSHRemote(remoteURL string) bool {
	return strings.HasPrefix(remoteURL, "git@") || strings.HasPrefix(remoteURL, "ssh://")
}
