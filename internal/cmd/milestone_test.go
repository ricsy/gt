package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	if cmd.Use != "view <number>" {
		t.Errorf("expected use 'view <number>', got %s", cmd.Use)
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
	if cmd.Use != "update <number>" {
		t.Errorf("expected use 'update <number>', got %s", cmd.Use)
	}
}

func TestMilestoneDeleteCommand(t *testing.T) {
	cmd := milestoneDeleteCmd
	if cmd.Use != "delete <number>" {
		t.Errorf("expected use 'delete <number>', got %s", cmd.Use)
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

	assertRequiredFlag(t, cmd, "due_on")
}

var _ = &cobra.Command{}

func assertRequiredFlag(t *testing.T, cmd *cobra.Command, name string) {
	t.Helper()

	flag := cmd.Flags().Lookup(name)
	if flag == nil {
		t.Fatalf("flag %q not found", name)
	}

	required := false
	cmd.Flags().VisitAll(func(current *pflag.Flag) {
		if current.Name == name {
			required = current.Annotations != nil && len(current.Annotations[cobra.BashCompOneRequiredFlag]) > 0
		}
	})

	if !required {
		t.Fatalf("expected flag %q to be marked required", name)
	}
}
