package cmd

import (
	"errors"
	"testing"

	"github.com/ricsy/gt/pkg/config"
)

func TestParseRepoFromRemoteURL(t *testing.T) {
	testCases := []struct {
		name         string
		remoteURL    string
		expectedHost string
		want         string
	}{
		{
			name:         "https origin on configured host",
			remoteURL:    "https://gitee.com/gitee/demo-repo.git",
			expectedHost: "gitee.com",
			want:         "gitee/demo-repo",
		},
		{
			name:         "ssh origin on configured host",
			remoteURL:    "git@gitee.com:gitee/demo-repo.git",
			expectedHost: "gitee.com",
			want:         "gitee/demo-repo",
		},
		{
			name:         "remote on different host is ignored",
			remoteURL:    "https://github.com/gitee/demo-repo.git",
			expectedHost: "gitee.com",
			want:         "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := parseRepoFromRemoteURL(testCase.remoteURL, testCase.expectedHost)
			if err != nil {
				t.Fatalf("parseRepoFromRemoteURL() returned error: %v", err)
			}
			if got != testCase.want {
				t.Fatalf("parseRepoFromRemoteURL() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestResolveRepoFlagFallsBackToGitRemote(t *testing.T) {
	originalCommandHost := commandHost
	originalInsideWorkTree := gitIsInsideWorkTree
	originalRemoteGetURL := gitRemoteGetURL
	originalConfigDirFunc := config.ConfigDirImpl
	commandHost = ""
	t.Cleanup(func() {
		commandHost = originalCommandHost
		gitIsInsideWorkTree = originalInsideWorkTree
		gitRemoteGetURL = originalRemoteGetURL
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	tmpDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return tmpDir })

	gitIsInsideWorkTree = func() (bool, error) { return true, nil }
	gitRemoteGetURL = func(name string) (string, error) {
		return "https://gitee.com/gitee/demo-repo.git", nil
	}

	owner, repoName, err := resolveRepoFlag("")
	if err != nil {
		t.Fatalf("resolveRepoFlag() returned error: %v", err)
	}
	if owner != "gitee" || repoName != "demo-repo" {
		t.Fatalf("resolveRepoFlag() = %s/%s, want gitee/demo-repo", owner, repoName)
	}
}

func TestResolveRepoFlagPrefersGitRemoteOverDefaultRepo(t *testing.T) {
	originalCommandHost := commandHost
	originalInsideWorkTree := gitIsInsideWorkTree
	originalRemoteGetURL := gitRemoteGetURL
	originalConfigDirFunc := config.ConfigDirImpl
	commandHost = ""
	t.Cleanup(func() {
		commandHost = originalCommandHost
		gitIsInsideWorkTree = originalInsideWorkTree
		gitRemoteGetURL = originalRemoteGetURL
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	tmpDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return tmpDir })

	cfg := config.DefaultConfig()
	cfg.DefaultRepo = "gitee/default-repo"
	if err := config.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() returned error: %v", err)
	}

	gitIsInsideWorkTree = func() (bool, error) { return true, nil }
	gitRemoteGetURL = func(name string) (string, error) {
		return "https://gitee.com/gitee/demo-repo.git", nil
	}

	owner, repoName, err := resolveRepoFlag("")
	if err != nil {
		t.Fatalf("resolveRepoFlag() returned error: %v", err)
	}
	if owner != "gitee" || repoName != "demo-repo" {
		t.Fatalf("resolveRepoFlag() = %s/%s, want gitee/demo-repo", owner, repoName)
	}
}

func TestResolveRepoFlagIgnoresNonGiteeRemote(t *testing.T) {
	originalCommandHost := commandHost
	originalInsideWorkTree := gitIsInsideWorkTree
	originalRemoteGetURL := gitRemoteGetURL
	originalConfigDirFunc := config.ConfigDirImpl
	commandHost = ""
	t.Cleanup(func() {
		commandHost = originalCommandHost
		gitIsInsideWorkTree = originalInsideWorkTree
		gitRemoteGetURL = originalRemoteGetURL
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	tmpDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return tmpDir })

	gitIsInsideWorkTree = func() (bool, error) { return true, nil }
	gitRemoteGetURL = func(name string) (string, error) {
		return "https://github.com/gitee/demo-repo.git", nil
	}

	if _, _, err := resolveRepoFlag(""); err == nil {
		t.Fatal("resolveRepoFlag() error = nil, want non-nil when no repo/default owner can be inferred")
	}
}

func TestResolveRepoFlagIgnoresRemoteLookupErrors(t *testing.T) {
	originalCommandHost := commandHost
	originalInsideWorkTree := gitIsInsideWorkTree
	originalRemoteGetURL := gitRemoteGetURL
	originalConfigDirFunc := config.ConfigDirImpl
	commandHost = ""
	t.Cleanup(func() {
		commandHost = originalCommandHost
		gitIsInsideWorkTree = originalInsideWorkTree
		gitRemoteGetURL = originalRemoteGetURL
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	tmpDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return tmpDir })

	gitIsInsideWorkTree = func() (bool, error) { return true, nil }
	gitRemoteGetURL = func(name string) (string, error) {
		return "", errors.New("missing origin")
	}

	if _, _, err := resolveRepoFlag(""); err == nil {
		t.Fatal("resolveRepoFlag() error = nil, want non-nil when no repo/default owner can be inferred")
	}
}
