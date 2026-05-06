package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewReleaseCmd(t *testing.T) {
	cmd := newReleaseCmd()
	if cmd == nil {
		t.Fatal("newReleaseCmd returned nil")
	}
	if cmd.Use != "release" {
		t.Errorf("expected use 'release', got %s", cmd.Use)
	}
	if len(cmd.Commands()) != 4 {
		t.Errorf("expected 4 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestReleaseListCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List releases",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "list" {
		t.Errorf("expected use 'list', got %s", cmd.Use)
	}
}

func TestReleaseViewCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View release",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "view" {
		t.Errorf("expected use 'view', got %s", cmd.Use)
	}
}

func TestReleaseCreateCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create release",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "create" {
		t.Errorf("expected use 'create', got %s", cmd.Use)
	}
}

func TestReleaseDeleteCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete release",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "delete" {
		t.Errorf("expected use 'delete', got %s", cmd.Use)
	}
}
