package cmd

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("gt version", version)
		},
	}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	expected := "gt version " + version + "\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestRootCommandHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestRootCommandHasTimeoutFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("timeout")
	if flag == nil {
		t.Fatal("expected root command to define a timeout flag")
	}
}

func TestRootCommandHasHostFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("host")
	if flag == nil {
		t.Fatal("expected root command to define a host flag")
	}
}

func TestCommandHTTPClientUsesRequestTimeout(t *testing.T) {
	originalTimeout := requestTimeout
	t.Cleanup(func() {
		requestTimeout = originalTimeout
	})

	requestTimeout = 7 * time.Second

	client := newCommandHTTPClient()
	if client.Timeout != requestTimeout {
		t.Fatalf("expected timeout %s, got %s", requestTimeout, client.Timeout)
	}
}

func TestResolveCommandHostPrecedence(t *testing.T) {
	originalCommandHost := commandHost
	commandHost = ""
	t.Cleanup(func() {
		commandHost = originalCommandHost
		_ = os.Unsetenv("GT_HOST")
		config.SetConfigDirFunc(config.ConfigDirImpl)
	})

	tmpDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return tmpDir })

	cfg := config.DefaultConfig()
	cfg.DefaultHost = "config.example.com"
	if err := config.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	if got := resolveCommandHost(); got != "config.example.com" {
		t.Fatalf("resolveCommandHost() = %s, want config.example.com", got)
	}

	if err := os.Setenv("GT_HOST", "env.example.com"); err != nil {
		t.Fatalf("Setenv() error = %v", err)
	}
	if got := resolveCommandHost(); got != "env.example.com" {
		t.Fatalf("resolveCommandHost() with env = %s, want env.example.com", got)
	}

	commandHost = "flag.example.com"
	if got := resolveCommandHost(); got != "flag.example.com" {
		t.Fatalf("resolveCommandHost() with flag = %s, want flag.example.com", got)
	}
}
