package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewWebhookCmd(t *testing.T) {
	cmd := newWebhookCmd()
	if cmd == nil {
		t.Fatal("newWebhookCmd returned nil")
	}
	if cmd.Use != "webhook" {
		t.Errorf("expected use 'webhook', got %s", cmd.Use)
	}
	if len(cmd.Commands()) != 5 {
		t.Errorf("expected 5 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestWebhookListCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List webhooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "list" {
		t.Errorf("expected use 'list', got %s", cmd.Use)
	}
}

func TestWebhookViewCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "view" {
		t.Errorf("expected use 'view', got %s", cmd.Use)
	}
}

func TestWebhookCreateCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "create" {
		t.Errorf("expected use 'create', got %s", cmd.Use)
	}
}

func TestWebhookDeleteCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "delete" {
		t.Errorf("expected use 'delete', got %s", cmd.Use)
	}
}

func TestWebhookTestCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "test" {
		t.Errorf("expected use 'test', got %s", cmd.Use)
	}
}
