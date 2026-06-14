package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

func TestRepositoryHasCommitHistory(t *testing.T) {
	testCases := []struct {
		name string
		repo *api.Repository
		want bool
	}{
		{
			name: "empty repo has no history",
			repo: &api.Repository{},
			want: false,
		},
		{
			name: "default branch implies history",
			repo: &api.Repository{DefaultBranch: "main"},
			want: true,
		},
		{
			name: "pushed at implies history",
			repo: &api.Repository{PushedAt: "2026-06-15T00:00:00Z"},
			want: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := repositoryHasCommitHistory(testCase.repo)
			if got != testCase.want {
				t.Fatalf("repositoryHasCommitHistory() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestConfirmRepositoryDeletionAllowsYesForEmptyRepo(t *testing.T) {
	repo := &api.Repository{FullName: "gitee/demo-repo"}

	cmd := &cobra.Command{}
	cmd.SetIn(bytes.NewBufferString(""))
	cmd.SetErr(new(bytes.Buffer))

	if err := confirmRepositoryDeletion(cmd, repo, true); err != nil {
		t.Fatalf("confirmRepositoryDeletion() returned error: %v", err)
	}
}

func TestConfirmRepositoryDeletionRejectsNonInteractiveCommittedRepo(t *testing.T) {
	repo := &api.Repository{
		FullName:      "gitee/demo-repo",
		DefaultBranch: "main",
	}

	cmd := &cobra.Command{}
	cmd.SetIn(bytes.NewBufferString("gitee/demo-repo\n"))
	cmd.SetErr(new(bytes.Buffer))

	if err := confirmRepositoryDeletion(cmd, repo, true); err == nil {
		t.Fatal("confirmRepositoryDeletion() error = nil, want non-nil for non-interactive repo with commit history")
	}
}

func TestPromptRepositoryDeletionConfirmationRejectsNonTerminalInput(t *testing.T) {
	tempFile, err := os.CreateTemp(t.TempDir(), "repo-delete-confirmation")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	defer func() { _ = tempFile.Close() }()

	if _, err := tempFile.WriteString("gitee/demo-repo\n"); err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Seek() error = %v", err)
	}

	cmd := &cobra.Command{}
	cmd.SetIn(tempFile)
	cmd.SetErr(new(bytes.Buffer))

	if err := promptRepositoryDeletionConfirmation(cmd, "gitee/demo-repo", false); err == nil {
		t.Fatal("promptRepositoryDeletionConfirmation() error = nil, want non-nil for non-terminal input")
	}
}

func TestPromptRepositoryDeletionConfirmationMessageIsExplicit(t *testing.T) {
	output := buildRepositoryDeletionPrompt("gitee/demo-repo", true)
	if !bytes.Contains([]byte(output), []byte("Do not type yes")) {
		t.Fatalf("expected explicit no-yes warning in output, got: %s", output)
	}
	if !bytes.Contains([]byte(output), []byte("Confirmation (expected: gitee/demo-repo)")) {
		t.Fatalf("expected full repository name hint in output, got: %s", output)
	}
}

func TestRepoDeleteCommandSupportsYesShortFlag(t *testing.T) {
	flag := repoDeleteCmd.Flags().Lookup("yes")
	if flag == nil {
		t.Fatal("expected repo delete command to define a yes flag")
	}
	if flag.Shorthand != "y" {
		t.Fatalf("expected yes flag shorthand to be y, got %q", flag.Shorthand)
	}
}
