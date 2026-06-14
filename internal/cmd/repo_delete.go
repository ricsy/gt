package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var repoDeleteOpts struct {
	Yes bool
}

func repoDeleteCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := ResolveRepo(args[0])
	if err != nil {
		return err
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	repo, err := client.GetRepo(owner, repoName)
	if err != nil {
		return fmt.Errorf("failed to get repo: %w", err)
	}

	if err := confirmRepositoryDeletion(cmd, repo, repoDeleteOpts.Yes); err != nil {
		return err
	}

	if err := client.DeleteRepo(owner, repoName); err != nil {
		return fmt.Errorf("failed to delete repo: %w", err)
	}

	cmd.Printf("Deleted repository: %s\n", repo.FullName)
	return nil
}

func repositoryHasCommitHistory(repo *api.Repository) bool {
	if repo == nil {
		return false
	}
	return strings.TrimSpace(repo.DefaultBranch) != "" || strings.TrimSpace(repo.PushedAt) != ""
}

func confirmRepositoryDeletion(cmd *cobra.Command, repo *api.Repository, allowNonInteractive bool) error {
	if repo == nil {
		return fmt.Errorf("repository is required")
	}

	fullRepoName := repo.FullName
	if fullRepoName == "" {
		fullRepoName = strings.Trim(repo.Owner.Login+"/"+repo.Name, "/")
	}

	if !repositoryHasCommitHistory(repo) {
		if allowNonInteractive {
			return nil
		}
		return promptRepositoryDeletionConfirmation(cmd, fullRepoName, false)
	}

	return promptRepositoryDeletionConfirmation(cmd, fullRepoName, true)
}

func promptRepositoryDeletionConfirmation(cmd *cobra.Command, fullRepoName string, hasCommitHistory bool) error {
	input := cmd.InOrStdin()
	file, ok := input.(*os.File)
	if !ok || !term.IsTerminal(int(file.Fd())) {
		if hasCommitHistory {
			return fmt.Errorf("repository %s has commit history; rerun in an interactive terminal and type the full repository name to confirm deletion", fullRepoName)
		}
		return fmt.Errorf("repository deletion requires confirmation; rerun interactively or pass --yes only for repositories without commit history")
	}

	_, _ = fmt.Fprint(cmd.ErrOrStderr(), buildRepositoryDeletionPrompt(fullRepoName, hasCommitHistory))

	reader := bufio.NewReader(input)
	confirmation, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to read deletion confirmation: %w", err)
	}

	if strings.TrimSpace(confirmation) != fullRepoName {
		return fmt.Errorf("deletion aborted: confirmation did not match %s", fullRepoName)
	}

	return nil
}

func buildRepositoryDeletionPrompt(fullRepoName string, hasCommitHistory bool) string {
	var builder strings.Builder
	if hasCommitHistory {
		_, _ = fmt.Fprintf(&builder, "Repository %s has commit history and will be permanently deleted.\n", fullRepoName)
	}
	builder.WriteString("Do not type yes. Enter the full repository name to confirm deletion.\n")
	_, _ = fmt.Fprintf(&builder, "Confirmation (expected: %s): ", fullRepoName)
	return builder.String()
}
