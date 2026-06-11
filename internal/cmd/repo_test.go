package cmd

import (
	"strings"
	"testing"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
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
