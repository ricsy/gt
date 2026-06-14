package cmd

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestRootCommandHasEnvFileFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("env-file")
	if flag == nil {
		t.Fatal("expected root command to define an env-file flag")
	}
}

func TestRootCommandSilencesUsageOnError(t *testing.T) {
	if !rootCmd.SilenceUsage {
		t.Fatal("expected root command to silence usage output on command errors")
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

func TestResolveCommandEnvFilePrecedence(t *testing.T) {
	originalCommandEnvFile := commandEnvFile
	commandEnvFile = ""
	t.Cleanup(func() {
		commandEnvFile = originalCommandEnvFile
		_ = os.Unsetenv("GT_ENV_FILE")
	})

	if err := os.Setenv("GT_ENV_FILE", "env.from.var"); err != nil {
		t.Fatalf("Setenv() error = %v", err)
	}
	if got := resolveCommandEnvFile(); got != "env.from.var" {
		t.Fatalf("resolveCommandEnvFile() with env = %s, want env.from.var", got)
	}

	commandEnvFile = "env.from.flag"
	if got := resolveCommandEnvFile(); got != "env.from.flag" {
		t.Fatalf("resolveCommandEnvFile() with flag = %s, want env.from.flag", got)
	}
}

func TestLoadCommandEnvFileLoadsVariables(t *testing.T) {
	originalCommandEnvFile := commandEnvFile
	commandEnvFile = ""
	t.Cleanup(func() {
		commandEnvFile = originalCommandEnvFile
		_ = os.Unsetenv("GT_ENV_FILE")
		_ = os.Unsetenv("GT_SCOPE_CHECK")
	})

	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env.test")
	if err := os.WriteFile(envFile, []byte("GT_SCOPE_CHECK=loaded\n"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	commandEnvFile = envFile

	if err := loadCommandEnvFile(); err != nil {
		t.Fatalf("loadCommandEnvFile() error = %v", err)
	}
	if got := os.Getenv("GT_SCOPE_CHECK"); got != "loaded" {
		t.Fatalf("GT_SCOPE_CHECK = %q, want %q", got, "loaded")
	}
}
