package cmd

import (
	"testing"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
)

func TestResolveRepoSupportsBareNameWithOrgScope(t *testing.T) {
	originalConfigDirFunc := config.ConfigDirImpl
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	cfg := config.DefaultConfig()
	cfg.RepoScopeMode = repoScopeModeOrg
	cfg.RepoScopeNamespace = "gitee"
	if err := config.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() returned error: %v", err)
	}

	owner, repoName, err := ResolveRepo("demo-repo")
	if err != nil {
		t.Fatalf("ResolveRepo() returned error: %v", err)
	}
	if owner != "gitee" || repoName != "demo-repo" {
		t.Fatalf("ResolveRepo() = %s/%s, want gitee/demo-repo", owner, repoName)
	}
}

func TestResolveRepoRejectsOwnerOutsideOrgScope(t *testing.T) {
	originalConfigDirFunc := config.ConfigDirImpl
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	cfg := config.DefaultConfig()
	cfg.RepoScopeMode = repoScopeModeOrg
	cfg.RepoScopeNamespace = "gitee"
	if err := config.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() returned error: %v", err)
	}

	if _, _, err := ResolveRepo("someone/demo-repo"); err == nil {
		t.Fatal("ResolveRepo() error = nil, want non-nil when owner is outside org scope")
	}
}

func TestBuildCreateRepoOptionsUsesOrgScopeNamespace(t *testing.T) {
	originalConfigDirFunc := config.ConfigDirImpl
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	cfg := config.DefaultConfig()
	cfg.RepoScopeMode = repoScopeModeOrg
	cfg.RepoScopeNamespace = "gitee"
	if err := config.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() returned error: %v", err)
	}

	cmd := newRepoCreateTestCommand()
	resetRepoCreateOpts()
	repoCreateOpts.Name = "demo-repo"

	opts, err := buildCreateRepoOptions(cmd)
	if err != nil {
		t.Fatalf("buildCreateRepoOptions() returned error: %v", err)
	}
	if opts.Namespace != "gitee" {
		t.Fatalf("opts.Namespace = %q, want %q", opts.Namespace, "gitee")
	}
}

func TestBuildCreateRepoOptionsRejectsNamespaceOutsidePersonalScope(t *testing.T) {
	originalConfigDirFunc := config.ConfigDirImpl
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	if err := auth.Login("gitee.com", "test-token", "test-user"); err != nil {
		t.Fatalf("auth.Login() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = auth.Logout("gitee.com")
	})

	cfg := config.DefaultConfig()
	cfg.RepoScopeMode = repoScopeModePersonal
	if err := config.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() returned error: %v", err)
	}

	cmd := newRepoCreateTestCommand()
	if err := cmd.Flags().Set("namespace", "gitee"); err != nil {
		t.Fatalf("set namespace flag: %v", err)
	}

	resetRepoCreateOpts()
	repoCreateOpts.Name = "demo-repo"
	repoCreateOpts.Namespace = "gitee"

	if _, err := buildCreateRepoOptions(cmd); err == nil {
		t.Fatal("buildCreateRepoOptions() error = nil, want non-nil for namespace outside personal scope")
	}
}

func TestResolveRepoUsesEnvScopeOverride(t *testing.T) {
	t.Setenv("GT_REPO_SCOPE_MODE", repoScopeModeOrg)
	t.Setenv("GT_REPO_SCOPE_NAMESPACE", "gitee")

	owner, repoName, err := ResolveRepo("traceops")
	if err != nil {
		t.Fatalf("ResolveRepo() returned error: %v", err)
	}
	if owner != "gitee" || repoName != "traceops" {
		t.Fatalf("ResolveRepo() = %s/%s, want gitee/traceops", owner, repoName)
	}
}
