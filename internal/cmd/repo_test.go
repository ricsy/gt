package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/ricsy/gt/pkg/api"
	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

func resetRepoCreateOpts() {
	repoCreateOpts = repoCreateOptions{
		Private:    true,
		HasIssues:  true,
		HasWiki:    true,
		CanComment: true,
	}
}

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

func newRepoCreateTestCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "create"}
	cmd.Flags().Bool("public", false, "")
	cmd.Flags().Bool("private", true, "")
	cmd.Flags().Bool("has-issues", true, "")
	cmd.Flags().Bool("has-wiki", true, "")
	cmd.Flags().Bool("can-comment", true, "")
	cmd.Flags().Bool("auto-init", false, "")
	cmd.Flags().String("gitignore-template", "", "")
	cmd.Flags().String("license-template", "", "")
	cmd.Flags().String("homepage", "", "")
	cmd.Flags().String("path", "", "")
	return cmd
}

func TestBuildCreateRepoOptionsRejectsPublicVisibility(t *testing.T) {
	cmd := newRepoCreateTestCommand()
	if err := cmd.Flags().Set("public", "true"); err != nil {
		t.Fatalf("set public flag: %v", err)
	}

	resetRepoCreateOpts()
	repoCreateOpts = repoCreateOptions{
		Name:        "demo",
		Description: "public repo",
		Public:      true,
		Private:     true,
		HasIssues:   true,
		HasWiki:     true,
		CanComment:  true,
	}

	_, err := buildCreateRepoOptions(cmd)
	if err == nil {
		t.Fatal("buildCreateRepoOptions() error = nil, want non-nil for --public")
	}
}

func TestBuildCreateRepoOptionsRejectsPrivateFalse(t *testing.T) {
	cmd := newRepoCreateTestCommand()
	if err := cmd.Flags().Set("private", "false"); err != nil {
		t.Fatalf("set private flag: %v", err)
	}

	resetRepoCreateOpts()
	repoCreateOpts = repoCreateOptions{
		Name:       "demo",
		Private:    false,
		HasIssues:  true,
		HasWiki:    true,
		CanComment: true,
	}

	_, err := buildCreateRepoOptions(cmd)
	if err == nil {
		t.Fatal("buildCreateRepoOptions() error = nil, want non-nil for --private=false")
	}
}

func TestBuildCreateRepoOptionsDefaultsMatchPersonalRepoAPI(t *testing.T) {
	cmd := newRepoCreateTestCommand()

	resetRepoCreateOpts()
	repoCreateOpts.Name = "demo"

	opts, err := buildCreateRepoOptions(cmd)
	if err != nil {
		t.Fatalf("buildCreateRepoOptions() returned error: %v", err)
	}

	if !opts.Private {
		t.Fatalf("opts.Private = %v, want true", opts.Private)
	}
	if opts.AutoInit {
		t.Fatalf("opts.AutoInit = %v, want false", opts.AutoInit)
	}
	if opts.HasIssues != nil || opts.HasWiki != nil || opts.CanComment != nil {
		t.Fatal("optional capability fields should stay nil when related flags are not changed")
	}
}

func TestBuildCreateRepoOptionsMapsAllSupportedPersonalFlags(t *testing.T) {
	cmd := newRepoCreateTestCommand()
	for name, value := range map[string]string{
		"has-issues":         "false",
		"has-wiki":           "false",
		"can-comment":        "false",
		"auto-init":          "true",
		"gitignore-template": "Go",
		"license-template":   "MIT",
		"homepage":           "https://example.com",
		"path":               "traceops-cli",
	} {
		if err := cmd.Flags().Set(name, value); err != nil {
			t.Fatalf("set %s: %v", name, err)
		}
	}

	resetRepoCreateOpts()
	repoCreateOpts = repoCreateOptions{
		Name:              "demo",
		Description:       "desc",
		Homepage:          "https://example.com",
		Private:           true,
		HasIssues:         false,
		HasWiki:           false,
		CanComment:        false,
		AutoInit:          true,
		GitignoreTemplate: "Go",
		LicenseTemplate:   "MIT",
		Path:              "traceops-cli",
	}

	opts, err := buildCreateRepoOptions(cmd)
	if err != nil {
		t.Fatalf("buildCreateRepoOptions() returned error: %v", err)
	}

	if opts.Homepage != "https://example.com" {
		t.Fatalf("opts.Homepage = %q, want %q", opts.Homepage, "https://example.com")
	}
	if !opts.AutoInit {
		t.Fatalf("opts.AutoInit = %v, want true", opts.AutoInit)
	}
	if opts.Path != "traceops-cli" {
		t.Fatalf("opts.Path = %q, want %q", opts.Path, "traceops-cli")
	}
	if opts.GitignoreTemplate != "Go" {
		t.Fatalf("opts.GitignoreTemplate = %q, want %q", opts.GitignoreTemplate, "Go")
	}
	if opts.LicenseTemplate != "MIT" {
		t.Fatalf("opts.LicenseTemplate = %q, want %q", opts.LicenseTemplate, "MIT")
	}
	if opts.HasIssues == nil || *opts.HasIssues != false {
		t.Fatalf("opts.HasIssues = %v, want false pointer", opts.HasIssues)
	}
	if opts.HasWiki == nil || *opts.HasWiki != false {
		t.Fatalf("opts.HasWiki = %v, want false pointer", opts.HasWiki)
	}
	if opts.CanComment == nil || *opts.CanComment != false {
		t.Fatalf("opts.CanComment = %v, want false pointer", opts.CanComment)
	}
}

func TestPrintRepoCreatePushDiagnosticsSuggestsAuthSetup(t *testing.T) {
	originalConfigDirFunc := config.ConfigDir
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	originalFill := gitCredentialFill
	originalHelper := gitCredentialHelperGet
	originalInside := gitIsInsideWorkTree
	originalRemoteGetURL := gitRemoteGetURL
	originalLsRemote := gitLsRemote
	t.Cleanup(func() {
		gitCredentialFill = originalFill
		gitCredentialHelperGet = originalHelper
		gitIsInsideWorkTree = originalInside
		gitRemoteGetURL = originalRemoteGetURL
		gitLsRemote = originalLsRemote
	})

	if err := auth.Login("gitee.com", "test-token", "test-user"); err != nil {
		t.Fatalf("auth.Login() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = auth.Logout("gitee.com")
	})

	gitCredentialHelperGet = func() (string, error) { return "manager", nil }
	gitCredentialFill = func(host string) (config.HostAuth, error) {
		return config.HostAuth{}, errors.New("missing credential")
	}
	gitIsInsideWorkTree = func() (bool, error) { return false, nil }
	gitRemoteGetURL = func(name string) (string, error) { return "", errors.New("no remote") }
	gitLsRemote = func(target string) error { return nil }

	cmd := &cobra.Command{}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	printRepoCreatePushDiagnostics(cmd, &api.Repository{CloneURL: "https://gitee.com/ricsy/traceops.git"})

	output := buf.String()
	if !strings.Contains(output, "gt auth setup") {
		t.Fatalf("expected auth setup guidance in output, got: %s", output)
	}
}

func TestPrintRepoCreatePushDiagnosticsDetectsRemoteMismatch(t *testing.T) {
	originalConfigDirFunc := config.ConfigDir
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	originalFill := gitCredentialFill
	originalHelper := gitCredentialHelperGet
	originalInside := gitIsInsideWorkTree
	originalRemoteGetURL := gitRemoteGetURL
	originalLsRemote := gitLsRemote
	t.Cleanup(func() {
		gitCredentialFill = originalFill
		gitCredentialHelperGet = originalHelper
		gitIsInsideWorkTree = originalInside
		gitRemoteGetURL = originalRemoteGetURL
		gitLsRemote = originalLsRemote
	})

	if err := auth.Login("gitee.com", "test-token", "test-user"); err != nil {
		t.Fatalf("auth.Login() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = auth.Logout("gitee.com")
	})

	gitCredentialHelperGet = func() (string, error) { return "manager", nil }
	gitCredentialFill = func(host string) (config.HostAuth, error) {
		return config.HostAuth{Token: "test-token", User: "test-user"}, nil
	}
	gitIsInsideWorkTree = func() (bool, error) { return true, nil }
	gitRemoteGetURL = func(name string) (string, error) {
		return "https://gitee.com/ricsy/old.git", nil
	}
	gitLsRemote = func(target string) error { return nil }

	cmd := &cobra.Command{}
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	printRepoCreatePushDiagnostics(cmd, &api.Repository{CloneURL: "https://gitee.com/ricsy/traceops.git"})

	output := buf.String()
	if !strings.Contains(output, "git remote set-url origin https://gitee.com/ricsy/traceops.git") {
		t.Fatalf("expected remote set-url guidance in output, got: %s", output)
	}
}
