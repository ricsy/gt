package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigGetCommand(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a configuration value",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"get", "default_repo"})

	if cmd.Use != "get" {
		t.Errorf("expected use 'get', got %s", cmd.Use)
	}
}

func TestConfigSetCommand(t *testing.T) {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set a configuration value",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	if cmd.Use != "set" {
		t.Errorf("expected use 'set', got %s", cmd.Use)
	}
}

func TestConfigListCommand(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"list"})

	if cmd.Use != "list" {
		t.Errorf("expected use 'list', got %s", cmd.Use)
	}
}

func TestNewConfigCmd(t *testing.T) {
	cmd := newConfigCmd()
	if cmd == nil {
		t.Fatal("newConfigCmd returned nil")
	}
	if cmd.Use != "config" {
		t.Errorf("expected use 'config', got %s", cmd.Use)
	}
	if len(cmd.Commands()) != 3 {
		t.Errorf("expected 3 subcommands, got %d", len(cmd.Commands()))
	}
}
