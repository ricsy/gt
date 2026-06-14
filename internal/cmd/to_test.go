package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

func TestResolveToTargetURLFallsBackToGitRemote(t *testing.T) {
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

	got, err := resolveToTargetURL("", false)
	if err != nil {
		t.Fatalf("resolveToTargetURL() returned error: %v", err)
	}
	if got != "https://gitee.com/gitee/demo-repo" {
		t.Fatalf("resolveToTargetURL() = %q, want %q", got, "https://gitee.com/gitee/demo-repo")
	}
}

func TestResolveToTargetURLSupportsExplicitRepo(t *testing.T) {
	got, err := resolveToTargetURL("gitee/demo-repo", false)
	if err != nil {
		t.Fatalf("resolveToTargetURL() returned error: %v", err)
	}
	if got != "https://gitee.com/gitee/demo-repo" {
		t.Fatalf("resolveToTargetURL() = %q, want %q", got, "https://gitee.com/gitee/demo-repo")
	}
}

func TestResolveToTargetURLSupportsNamespaceHomepage(t *testing.T) {
	got, err := resolveToTargetURL("gitee", false)
	if err != nil {
		t.Fatalf("resolveToTargetURL() returned error: %v", err)
	}
	if got != "https://gitee.com/gitee" {
		t.Fatalf("resolveToTargetURL() = %q, want %q", got, "https://gitee.com/gitee")
	}
}

func TestResolveToTargetURLTreatsSingleSegmentAsRepoWhenRequested(t *testing.T) {
	originalResolveCurrentUser := resolveCurrentUser
	commandHost = ""
	resolveCurrentUser = func(host string) (string, error) {
		return "demo-user", nil
	}
	t.Cleanup(func() {
		resolveCurrentUser = originalResolveCurrentUser
		commandHost = ""
	})

	got, err := resolveToTargetURL("demo-repo", true)
	if err != nil {
		t.Fatalf("resolveToTargetURL() returned error: %v", err)
	}
	if got != "https://gitee.com/demo-user/demo-repo" {
		t.Fatalf("resolveToTargetURL() = %q, want %q", got, "https://gitee.com/demo-user/demo-repo")
	}
}

func TestResolveToTargetURLRequiresCurrentUserForRepoMode(t *testing.T) {
	originalResolveCurrentUser := resolveCurrentUser
	resolveCurrentUser = func(host string) (string, error) {
		return "", errors.New("not logged in")
	}
	t.Cleanup(func() {
		resolveCurrentUser = originalResolveCurrentUser
	})

	if _, err := resolveToTargetURL("demo-repo", true); err == nil {
		t.Fatal("resolveToTargetURL() error = nil, want non-nil when repo mode cannot resolve the current user")
	}
}

func TestToCommandOpensResolvedURL(t *testing.T) {
	originalOpenBrowserURL := openBrowserURL
	originalResolveCurrentUser := resolveCurrentUser
	originalVerifyToRepoExists := verifyToRepoExists
	var openedURL string
	openBrowserURL = func(targetURL string) error {
		openedURL = targetURL
		return nil
	}
	resolveCurrentUser = func(host string) (string, error) {
		return "demo-user", nil
	}
	verifyToRepoExists = func(owner, repo string) error {
		return nil
	}
	t.Cleanup(func() {
		openBrowserURL = originalOpenBrowserURL
		resolveCurrentUser = originalResolveCurrentUser
		verifyToRepoExists = originalVerifyToRepoExists
		toOpts.Repo = false
	})

	toOpts.Repo = true
	cmd := &cobra.Command{}
	cmd.SetOut(new(bytes.Buffer))

	if err := toCommand(cmd, []string{"demo-repo"}); err != nil {
		t.Fatalf("toCommand() returned error: %v", err)
	}
	if openedURL != "https://gitee.com/demo-user/demo-repo" {
		t.Fatalf("openedURL = %q, want %q", openedURL, "https://gitee.com/demo-user/demo-repo")
	}
}

func TestToCommandCanOpenPublicNamespaceWithoutLogin(t *testing.T) {
	originalOpenBrowserURL := openBrowserURL
	originalVerifyToNamespaceExists := verifyToNamespaceExists
	originalResolveCurrentUser := resolveCurrentUser
	var openedURL string
	openBrowserURL = func(targetURL string) error {
		openedURL = targetURL
		return nil
	}
	verifyToNamespaceExists = func(namespace string) error {
		return nil
	}
	resolveCurrentUser = func(host string) (string, error) {
		return "", errors.New("not logged in")
	}
	t.Cleanup(func() {
		openBrowserURL = originalOpenBrowserURL
		verifyToNamespaceExists = originalVerifyToNamespaceExists
		resolveCurrentUser = originalResolveCurrentUser
	})

	cmd := &cobra.Command{}
	cmd.SetOut(new(bytes.Buffer))

	if err := toCommand(cmd, []string{"gitee"}); err != nil {
		t.Fatalf("toCommand() returned error: %v", err)
	}
	if openedURL != "https://gitee.com/gitee" {
		t.Fatalf("openedURL = %q, want %q", openedURL, "https://gitee.com/gitee")
	}
}

func TestToCommandDoesNotOpenMissingRepository(t *testing.T) {
	originalOpenBrowserURL := openBrowserURL
	originalVerifyToRepoExists := verifyToRepoExists
	calledOpen := false
	openBrowserURL = func(targetURL string) error {
		calledOpen = true
		return nil
	}
	verifyToRepoExists = func(owner, repo string) error {
		return fmt.Errorf("repository not found: %s/%s", owner, repo)
	}
	t.Cleanup(func() {
		openBrowserURL = originalOpenBrowserURL
		verifyToRepoExists = originalVerifyToRepoExists
	})

	cmd := &cobra.Command{}
	cmd.SetOut(new(bytes.Buffer))

	err := toCommand(cmd, []string{"gitee/missing-repo"})
	if err == nil {
		t.Fatal("toCommand() error = nil, want non-nil for missing repository")
	}
	if err.Error() != "repository not found: gitee/missing-repo" {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledOpen {
		t.Fatal("expected browser not to open for missing repository")
	}
}

func TestToCommandDoesNotOpenMissingNamespace(t *testing.T) {
	originalOpenBrowserURL := openBrowserURL
	originalVerifyToNamespaceExists := verifyToNamespaceExists
	calledOpen := false
	openBrowserURL = func(targetURL string) error {
		calledOpen = true
		return nil
	}
	verifyToNamespaceExists = func(namespace string) error {
		return fmt.Errorf("namespace not found: %s", namespace)
	}
	t.Cleanup(func() {
		openBrowserURL = originalOpenBrowserURL
		verifyToNamespaceExists = originalVerifyToNamespaceExists
	})

	cmd := &cobra.Command{}
	cmd.SetOut(new(bytes.Buffer))

	err := toCommand(cmd, []string{"missing-namespace"})
	if err == nil {
		t.Fatal("toCommand() error = nil, want non-nil for missing namespace")
	}
	if err.Error() != "namespace not found: missing-namespace" {
		t.Fatalf("unexpected error: %v", err)
	}
	if calledOpen {
		t.Fatal("expected browser not to open for missing namespace")
	}
}

func TestToCommandSupportsRepoFlagShorthand(t *testing.T) {
	flag := toCmd.Flags().Lookup("repo")
	if flag == nil {
		t.Fatal("expected to command to define a repo flag")
	}
	if flag.Shorthand != "r" {
		t.Fatalf("expected repo flag shorthand to be r, got %q", flag.Shorthand)
	}
}
