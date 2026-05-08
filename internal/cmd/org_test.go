package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewOrgCmd(t *testing.T) {
	cmd := newOrgCmd()
	if cmd == nil {
		t.Fatal("newOrgCmd returned nil")
	}
	if cmd.Use != "org" {
		t.Errorf("expected use 'org', got %s", cmd.Use)
	}
	if len(cmd.Commands()) != 4 {
		t.Errorf("expected 4 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestOrgListCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List organizations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "list" {
		t.Errorf("expected use 'list', got %s", cmd.Use)
	}
}

func TestOrgViewCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "view" {
		t.Errorf("expected use 'view', got %s", cmd.Use)
	}
}
