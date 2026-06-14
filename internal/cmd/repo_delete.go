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

type repoDeletionCommitSummary struct {
	LatestAt      string
	LatestTitle   string
	HasLatestInfo bool
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

	var commitSummary *repoDeletionCommitSummary
	if repositoryHasCommitHistory(repo) {
		historySummary, err := client.GetRepoCommitHistorySummary(owner, repoName, repo.DefaultBranch)
		if err != nil {
			return fmt.Errorf("failed to inspect repo commit history before deletion: %w", err)
		}
		commitSummary = buildRepoDeletionCommitSummary(historySummary)
	}

	if err := confirmRepositoryDeletion(cmd, repo, commitSummary, repoDeleteOpts.Yes); err != nil {
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

func confirmRepositoryDeletion(cmd *cobra.Command, repo *api.Repository, commitSummary *repoDeletionCommitSummary, allowNonInteractive bool) error {
	if repo == nil {
		return fmt.Errorf("repository is required")
	}

	fullRepoName := repo.FullName
	if fullRepoName == "" {
		fullRepoName = strings.Trim(repo.Owner.Login+"/"+repo.Name, "/")
	}

	if !repositoryHasCommitHistory(repo) {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "No commit history detected for %s.\n", fullRepoName)
		if allowNonInteractive {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Proceeding with deletion because -y/--yes is allowed only for repositories without commit history.\n")
			return nil
		}
		return promptRepositoryDeletionConfirmation(cmd, fullRepoName, nil)
	}

	return promptRepositoryDeletionConfirmation(cmd, fullRepoName, commitSummary)
}

func promptRepositoryDeletionConfirmation(cmd *cobra.Command, fullRepoName string, commitSummary *repoDeletionCommitSummary) error {
	input := cmd.InOrStdin()
	file, ok := input.(*os.File)
	if !ok || !term.IsTerminal(int(file.Fd())) {
		if commitSummary != nil {
			return fmt.Errorf("repository %s has commit history (%s); rerun in an interactive terminal and type the full repository name to confirm deletion", fullRepoName, formatRepoDeletionCommitSummary(commitSummary))
		}
		return fmt.Errorf("repository deletion requires confirmation; rerun interactively or pass --yes only for repositories without commit history")
	}

	_, _ = fmt.Fprint(cmd.ErrOrStderr(), buildRepositoryDeletionPrompt(fullRepoName, commitSummary))

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

func buildRepositoryDeletionPrompt(fullRepoName string, commitSummary *repoDeletionCommitSummary) string {
	var builder strings.Builder
	if commitSummary != nil {
		_, _ = fmt.Fprintf(&builder, "Repository %s has commit history and will be permanently deleted.\n", fullRepoName)
		_, _ = fmt.Fprintf(&builder, "Commit history: %s.\n", formatRepoDeletionCommitSummary(commitSummary))
	}
	builder.WriteString("Do not type yes. Enter the full repository name to confirm deletion.\n")
	_, _ = fmt.Fprintf(&builder, "Confirmation (expected: %s): ", fullRepoName)
	return builder.String()
}

func buildRepoDeletionCommitSummary(historySummary *api.RepoCommitHistorySummary) *repoDeletionCommitSummary {
	if historySummary == nil {
		return nil
	}

	summary := &repoDeletionCommitSummary{}
	if historySummary.Latest == nil {
		return summary
	}

	latestAt := strings.TrimSpace(historySummary.Latest.Commit.Committer.Date)
	if latestAt == "" {
		latestAt = strings.TrimSpace(historySummary.Latest.Commit.Author.Date)
	}

	latestTitle := extractCommitTitle(historySummary.Latest.Commit.Message)
	summary.LatestAt = latestAt
	summary.LatestTitle = latestTitle
	summary.HasLatestInfo = latestAt != "" || latestTitle != ""
	return summary
}

func formatRepoDeletionCommitSummary(summary *repoDeletionCommitSummary) string {
	if summary == nil {
		return "commit history detected"
	}

	if summary.HasLatestInfo {
		switch {
		case summary.LatestAt != "" && summary.LatestTitle != "":
			return fmt.Sprintf("latest %s - %s", summary.LatestAt, summary.LatestTitle)
		case summary.LatestAt != "":
			return fmt.Sprintf("latest %s", summary.LatestAt)
		case summary.LatestTitle != "":
			return fmt.Sprintf("latest %s", summary.LatestTitle)
		}
	}
	return "commit history detected"
}

func extractCommitTitle(message string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return ""
	}

	if newlineIndex := strings.IndexByte(message, '\n'); newlineIndex >= 0 {
		message = message[:newlineIndex]
	}

	return strings.TrimSpace(message)
}
