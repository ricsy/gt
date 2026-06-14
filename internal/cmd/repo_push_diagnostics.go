package cmd

import (
	"strings"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

// printRepoCreatePushDiagnostics explains the post-create push state without mutating the local repository.
func printRepoCreatePushDiagnostics(cmd *cobra.Command, repo *api.Repository, cloneURLMode string) {
	host := resolveCommandHost()
	cloneURL, err := resolveRepoCloneURL(host, repo.Owner.Login, repo.Name, repo.CloneURL, repo.SSHURL, cloneURLMode)

	cmd.Println("Push diagnostics:")
	cmd.Println("- Repository creation: success")
	if err != nil {
		cmd.Printf("- Remote URL: unavailable (%v)\n", err)
		return
	}
	cmd.Printf("- Suggested remote URL: %s\n", cloneURL)

	report := runAuthDoctor(host)
	if cloneURLMode == cloneURLModeHTTPS {
		if report.StoredAuthErr != "" {
			cmd.Printf("- HTTPS auth: missing stored auth for %s. Run: gt auth login --username <name> --token <token>\n", host)
		} else if report.GitCredentialErr != "" || !report.AuthUserMatchesGit {
			cmd.Printf("- HTTPS auth: not ready for git push. Run: gt auth setup --overwrite\n")
		} else {
			cmd.Printf("- HTTPS auth: ready\n")
		}
	} else {
		cmd.Printf("- Remote mode: SSH (gt does not validate SSH agent keys here)\n")
	}

	insideWorkTree, workTreeErr := gitIsInsideWorkTree()
	if workTreeErr != nil {
		cmd.Printf("- Local repository check: failed (%v)\n", workTreeErr)
		return
	}
	if !insideWorkTree {
		cmd.Printf("- Local repository check: current directory is not a git worktree\n")
		cmd.Printf("- Next steps:\n")
		cmd.Printf("  git remote add origin %s\n", cloneURL)
		cmd.Printf("  git push -u origin %s\n", preferredRepoPushBranch())
		return
	}

	originURL, originErr := gitRemoteGetURL("origin")
	if originErr != nil {
		cmd.Printf("- Origin remote: missing\n")
		cmd.Printf("- Next steps:\n")
		cmd.Printf("  git remote add origin %s\n", cloneURL)
		cmd.Printf("  git push -u origin %s\n", preferredRepoPushBranch())
		return
	}

	if originURL == cloneURL {
		cmd.Printf("- Origin remote: already points to the new repository\n")
	} else {
		cmd.Printf("- Origin remote: points to %s\n", originURL)
		cmd.Printf("- Update command:\n")
		cmd.Printf("  git remote set-url origin %s\n", cloneURL)
	}

	if cloneURLMode == cloneURLModeSSH {
		cmd.Printf("- SSH reachability: skipped\n")
		return
	}

	if report.StoredAuthErr != "" || report.GitCredentialErr != "" || !report.AuthUserMatchesGit {
		cmd.Printf("- HTTPS reachability: skipped until auth is fixed\n")
		return
	}

	targetURL := cloneURL
	if isSSHRemote(originURL) || originURL != cloneURL {
		targetURL = cloneURL
	}
	if err := gitLsRemote(targetURL); err != nil {
		cmd.Printf("- HTTPS reachability: failed (%v)\n", err)
		cmd.Printf("- Recommendation: run gt auth setup --overwrite\n")
		return
	}

	cmd.Printf("- HTTPS reachability: ok\n")
}

func isSSHRemote(remoteURL string) bool {
	return strings.HasPrefix(remoteURL, "git@") || strings.HasPrefix(remoteURL, "ssh://")
}

func preferredRepoPushBranch() string {
	branch, err := gitCurrentBranch()
	if err != nil || strings.TrimSpace(branch) == "" {
		return repoPrimaryBranch
	}

	return strings.TrimSpace(branch)
}
