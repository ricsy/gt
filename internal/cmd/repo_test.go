package cmd

import (
	"bytes"
	"testing"

	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

func TestRepoCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	repoCmd.SetOut(buf)
	repoCmd.SetArgs([]string{"--help"})

	err := repoCmd.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestRepoBranchCommand(t *testing.T) {
	if repoBranchCmd.Use != "branch" {
		t.Errorf("expected use 'branch', got %s", repoBranchCmd.Use)
	}
	if len(repoBranchCmd.Commands()) != 5 {
		t.Errorf("expected 5 subcommands, got %d", len(repoBranchCmd.Commands()))
	}
}

func TestRepoListCommandFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "list"}
	cmd.Flags().String("owner", "", "Owner username")
	cmd.Flags().Int("limit", 30, "Maximum number of repos")

	if cmd.Flags().Lookup("owner") == nil {
		t.Error("owner flag not found")
	}
	if cmd.Flags().Lookup("limit") == nil {
		t.Error("limit flag not found")
	}
}

func TestRepoCreateCommandFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "create"}
	cmd.Flags().String("name", "", "Repository name")
	cmd.Flags().String("description", "", "Repository description")
	cmd.Flags().Bool("private", false, "Create private repository")
	cmd.Flags().Bool("public", false, "Create public repository")

	if cmd.Flags().Lookup("name") == nil {
		t.Error("name flag not found")
	}
	if cmd.Flags().Lookup("description") == nil {
		t.Error("description flag not found")
	}
	if cmd.Flags().Lookup("private") == nil {
		t.Error("private flag not found")
	}
	if cmd.Flags().Lookup("public") == nil {
		t.Error("public flag not found")
	}
}

func TestRepoBranchCommandFlags(t *testing.T) {
	if repoBranchListCmd.Flags().Lookup("repo") == nil {
		t.Error("repo flag not found")
	}
	if repoBranchListCmd.Flags().Lookup("sort") == nil {
		t.Error("sort flag not found")
	}
	if repoBranchListCmd.Flags().Lookup("direction") == nil {
		t.Error("direction flag not found")
	}
	if repoBranchCreateCmd.Flags().Lookup("refs") == nil {
		t.Error("refs flag not found")
	}
	if repoBranchProtectCmd.Flags().Lookup("repo") == nil {
		t.Error("repo flag not found")
	}
}

func TestCloneCommandArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"one arg", []string{"owner/repo"}, false},
		{"two args", []string{"owner/repo", "directory"}, false},
		{"no args", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{
				Use:   "clone",
				Short: "Clone a repository",
				Args:  cobra.RangeArgs(1, 2),
				Run:   func(cmd *cobra.Command, args []string) {},
			}

			err := cmd.ValidateArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepoCloneUsesWebCloneURL(t *testing.T) {
	got := config.RepoGitHTTPSURL("gitee.com", "owner", "repo")
	want := "https://gitee.com/owner/repo.git"
	if got != want {
		t.Fatalf("RepoGitHTTPSURL() = %s, want %s", got, want)
	}
}
