package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigDir(t *testing.T) {
	dir := ConfigDir()
	if dir == "" {
		t.Error("ConfigDir returned empty string")
	}
	if !filepath.IsAbs(dir) {
		t.Errorf("ConfigDir should return absolute path, got: %s", dir)
	}
}

func TestHostsFile(t *testing.T) {
	hostsFile := HostsFile()
	expectedSuffix := filepath.Join(".config", "gt", "hosts.yml")
	if !strings.HasSuffix(hostsFile, expectedSuffix) {
		t.Errorf("HostsFile() = %s, want suffix %s", hostsFile, expectedSuffix)
	}
}

func TestConfigFile(t *testing.T) {
	cfgFile := ConfigFile()
	expectedSuffix := filepath.Join(".config", "gt", "config.yml")
	if !strings.HasSuffix(cfgFile, expectedSuffix) {
		t.Errorf("ConfigFile() = %s, want suffix %s", cfgFile, expectedSuffix)
	}
}

func TestLoadSaveHosts(t *testing.T) {
	tmpDir := t.TempDir()
	configDirFunc = func() string { return tmpDir }
	defer func() { configDirFunc = ConfigDirImpl }()

	hosts := map[string]HostAuth{
		"gitee.com": {
			Token: "test-token",
			User:  "testuser",
		},
		"github.com": {
			Token: "gh-token",
			User:  "ghuser",
		},
	}

	if err := SaveHosts(hosts); err != nil {
		t.Fatalf("SaveHosts failed: %v", err)
	}

	loaded, err := LoadHosts()
	if err != nil {
		t.Fatalf("LoadHosts failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(loaded))
	}

	if loaded["gitee.com"].Token != "test-token" {
		t.Errorf("gitee.com token = %s, want test-token", loaded["gitee.com"].Token)
	}

	if loaded["github.com"].User != "ghuser" {
		t.Errorf("github.com user = %s, want ghuser", loaded["github.com"].User)
	}
}

func TestLoadSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDirFunc = func() string { return tmpDir }
	defer func() { configDirFunc = ConfigDirImpl }()

	cfg := &Config{
		DefaultRepo:  "my-repo",
		DefaultOwner: "myowner",
		DefaultHost:  "gitee.example.com",
	}

	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loaded.DefaultRepo != "my-repo" {
		t.Errorf("DefaultRepo = %s, want my-repo", loaded.DefaultRepo)
	}

	if loaded.DefaultOwner != "myowner" {
		t.Errorf("DefaultOwner = %s, want myowner", loaded.DefaultOwner)
	}

	if loaded.DefaultHost != "gitee.example.com" {
		t.Errorf("DefaultHost = %s, want gitee.example.com", loaded.DefaultHost)
	}
}

func TestLoadHostsNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	configDirFunc = func() string { return tmpDir }
	defer func() { configDirFunc = ConfigDirImpl }()

	hosts, err := LoadHosts()
	if err != nil {
		t.Fatalf("LoadHosts failed for non-existent file: %v", err)
	}

	if _, ok := hosts["gitee.com"]; !ok {
		t.Errorf("expected default gitee.com entry in hosts map")
	}
}

func TestDefaultConfigUsesDefaultHost(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.DefaultHost != DefaultHost {
		t.Fatalf("DefaultHost = %s, want %s", cfg.DefaultHost, DefaultHost)
	}
}

func TestRepoGitHTTPSURL(t *testing.T) {
	got := RepoGitHTTPSURL("gitee.example.com", "owner", "repo")
	want := "https://gitee.example.com/owner/repo.git"
	if got != want {
		t.Fatalf("RepoGitHTTPSURL() = %s, want %s", got, want)
	}
}
