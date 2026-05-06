package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestSearchCmd(t *testing.T) {
	cmd := searchCmd
	if cmd == nil {
		t.Fatal("searchCmd returned nil")
	}
	if cmd.Use != "search" {
		t.Errorf("expected use 'search', got %s", cmd.Use)
	}
	if len(cmd.Commands()) != 3 {
		t.Errorf("expected 3 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestSearchReposCommand(t *testing.T) {
	cmd := searchReposCmd
	if cmd.Use != "repos" {
		t.Errorf("expected use 'repos', got %s", cmd.Use)
	}
}

func TestSearchIssuesCommand(t *testing.T) {
	cmd := searchIssuesCmd
	if cmd.Use != "issues" {
		t.Errorf("expected use 'issues', got %s", cmd.Use)
	}
}

func TestSearchUsersCommand(t *testing.T) {
	cmd := searchUsersCmd
	if cmd.Use != "users" {
		t.Errorf("expected use 'users', got %s", cmd.Use)
	}
}

func TestSearchReposFlags(t *testing.T) {
	cmd := searchReposCmd
	flags := cmd.Flags()

	if flags.Lookup("q") == nil {
		t.Error("expected --q flag")
	}
	if flags.Lookup("owner") == nil {
		t.Error("expected --owner flag")
	}
	if flags.Lookup("language") == nil {
		t.Error("expected --language flag")
	}
	if flags.Lookup("sort") == nil {
		t.Error("expected --sort flag")
	}
	if flags.Lookup("order") == nil {
		t.Error("expected --order flag")
	}
	if flags.Lookup("fork") == nil {
		t.Error("expected --fork flag")
	}
	if flags.Lookup("page") == nil {
		t.Error("expected --page flag")
	}
	if flags.Lookup("per-page") == nil {
		t.Error("expected --per-page flag")
	}
}

func TestSearchIssuesFlags(t *testing.T) {
	cmd := searchIssuesCmd
	flags := cmd.Flags()

	if flags.Lookup("q") == nil {
		t.Error("expected --q flag")
	}
	if flags.Lookup("repo") == nil {
		t.Error("expected --repo flag")
	}
	if flags.Lookup("language") == nil {
		t.Error("expected --language flag")
	}
	if flags.Lookup("label") == nil {
		t.Error("expected --label flag")
	}
	if flags.Lookup("state") == nil {
		t.Error("expected --state flag")
	}
	if flags.Lookup("author") == nil {
		t.Error("expected --author flag")
	}
	if flags.Lookup("assignee") == nil {
		t.Error("expected --assignee flag")
	}
}

func TestSearchUsersFlags(t *testing.T) {
	cmd := searchUsersCmd
	flags := cmd.Flags()

	if flags.Lookup("q") == nil {
		t.Error("expected --q flag")
	}
	if flags.Lookup("sort") == nil {
		t.Error("expected --sort flag")
	}
	if flags.Lookup("order") == nil {
		t.Error("expected --order flag")
	}
	if flags.Lookup("page") == nil {
		t.Error("expected --page flag")
	}
	if flags.Lookup("per-page") == nil {
		t.Error("expected --per-page flag")
	}
}

var _ = &cobra.Command{}
