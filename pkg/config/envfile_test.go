package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveEnvFilePrefersExplicitPath(t *testing.T) {
	t.Setenv("GT_ENV_FILE", "env.from.var")

	got := ResolveEnvFile("env.from.flag")
	if got != "env.from.flag" {
		t.Fatalf("ResolveEnvFile() = %q, want %q", got, "env.from.flag")
	}
}

func TestResolveEnvFileFallsBackToEnv(t *testing.T) {
	t.Setenv("GT_ENV_FILE", "env.from.var")

	got := ResolveEnvFile("")
	if got != "env.from.var" {
		t.Fatalf("ResolveEnvFile() = %q, want %q", got, "env.from.var")
	}
}

func TestLoadEnvFileSetsMissingVariablesOnly(t *testing.T) {
	t.Setenv("GT_EXISTING", "keep-me")
	_ = os.Unsetenv("GT_FROM_FILE")
	_ = os.Unsetenv("GT_QUOTED")

	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env.test")
	content := "GT_EXISTING=override-me\nGT_FROM_FILE=loaded\nexport GT_QUOTED=\"quoted value\"\n"
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := LoadEnvFile(envFile); err != nil {
		t.Fatalf("LoadEnvFile() error = %v", err)
	}

	if got := os.Getenv("GT_EXISTING"); got != "keep-me" {
		t.Fatalf("GT_EXISTING = %q, want %q", got, "keep-me")
	}
	if got := os.Getenv("GT_FROM_FILE"); got != "loaded" {
		t.Fatalf("GT_FROM_FILE = %q, want %q", got, "loaded")
	}
	if got := os.Getenv("GT_QUOTED"); got != "quoted value" {
		t.Fatalf("GT_QUOTED = %q, want %q", got, "quoted value")
	}
}
