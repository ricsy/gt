package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestMilestoneCmd(t *testing.T) {
	cmd := milestoneCmd
	if cmd == nil {
		t.Fatal("milestoneCmd returned nil")
	}
	if cmd.Use != "milestone" {
		t.Errorf("expected use 'milestone', got %s", cmd.Use)
	}
	if len(cmd.Commands()) != 5 {
		t.Errorf("expected 5 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestMilestoneListCommand(t *testing.T) {
	cmd := milestoneListCmd
	if cmd.Use != "list" {
		t.Errorf("expected use 'list', got %s", cmd.Use)
	}
}

func TestMilestoneViewCommand(t *testing.T) {
	cmd := milestoneViewCmd
	if cmd.Use != "view" {
		t.Errorf("expected use 'view', got %s", cmd.Use)
	}
}

func TestMilestoneCreateCommand(t *testing.T) {
	cmd := milestoneCreateCmd
	if cmd.Use != "create" {
		t.Errorf("expected use 'create', got %s", cmd.Use)
	}
}

func TestMilestoneUpdateCommand(t *testing.T) {
	cmd := milestoneUpdateCmd
	if cmd.Use != "update" {
		t.Errorf("expected use 'update', got %s", cmd.Use)
	}
}

func TestMilestoneDeleteCommand(t *testing.T) {
	cmd := milestoneDeleteCmd
	if cmd.Use != "delete" {
		t.Errorf("expected use 'delete', got %s", cmd.Use)
	}
}

func TestMilestoneListFlags(t *testing.T) {
	cmd := milestoneListCmd
	flags := cmd.Flags()

	if flags.Lookup("repo") == nil {
		t.Error("expected --repo flag")
	}
	if flags.Lookup("owner") == nil {
		t.Error("expected --owner flag")
	}
	if flags.Lookup("state") == nil {
		t.Error("expected --state flag")
	}
}

func TestMilestoneCreateFlags(t *testing.T) {
	cmd := milestoneCreateCmd
	flags := cmd.Flags()

	if flags.Lookup("title") == nil {
		t.Error("expected --title flag")
	}
	if flags.Lookup("description") == nil {
		t.Error("expected --description flag")
	}
	if flags.Lookup("due_on") == nil {
		t.Error("expected --due_on flag")
	}
}

var _ = &cobra.Command{}
