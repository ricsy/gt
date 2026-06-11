package cmd

import (
	"strings"
	"testing"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

func TestBuildAuthenticatedCloneURL(t *testing.T) {
	const host = "gitee.com"
	originalConfigDirFunc := config.ConfigDir
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	err := auth.Login(host, "test-token", "test-user")
	if err != nil {
		t.Fatalf("auth.Login() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = auth.Logout(host)
	})

	got, err := buildAuthenticatedCloneURL(host, "https://gitee.com/owner/repo.git")
	if err != nil {
		t.Fatalf("buildAuthenticatedCloneURL() returned error: %v", err)
	}

	if !strings.Contains(got, "https://test-user:test-token@gitee.com/owner/repo.git") {
		t.Fatalf("buildAuthenticatedCloneURL() = %q, want credentials embedded", got)
	}
}

func TestBuildAuthenticatedCloneURLWithoutStoredUser(t *testing.T) {
	originalConfigDirFunc := config.ConfigDir
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	got, err := buildAuthenticatedCloneURL("gitee.com", "https://gitee.com/owner/repo.git")
	if err == nil {
		t.Fatalf("buildAuthenticatedCloneURL() error = nil, want non-nil, got URL %q", got)
	}
}

func TestBuildCreateRepoOptionsRejectsPublicVisibility(t *testing.T) {
	cmd := &cobra.Command{Use: "create"}
	cmd.Flags().Bool("public", false, "")
	cmd.Flags().Bool("private", false, "")
	if err := cmd.Flags().Set("public", "true"); err != nil {
		t.Fatalf("set public flag: %v", err)
	}

	repoCreateOpts = struct {
		Name        string
		Description string
		Private     bool
		Public      bool
	}{
		Name:        "demo",
		Description: "public repo",
		Public:      true,
	}

	_, err := buildCreateRepoOptions(cmd)
	if err == nil {
		t.Fatal("buildCreateRepoOptions() error = nil, want non-nil for --public")
	}
}

func TestBuildCreateRepoOptionsKeepsPrivateVisibility(t *testing.T) {
	cmd := &cobra.Command{Use: "create"}
	cmd.Flags().Bool("public", false, "")
	cmd.Flags().Bool("private", false, "")
	if err := cmd.Flags().Set("private", "true"); err != nil {
		t.Fatalf("set private flag: %v", err)
	}

	repoCreateOpts = struct {
		Name        string
		Description string
		Private     bool
		Public      bool
	}{
		Name:        "demo",
		Description: "private repo",
		Private:     true,
	}

	opts, err := buildCreateRepoOptions(cmd)
	if err != nil {
		t.Fatalf("buildCreateRepoOptions() returned error: %v", err)
	}

	if opts.Private != true {
		t.Fatalf("opts.Private = %v, want true", opts.Private)
	}
}
