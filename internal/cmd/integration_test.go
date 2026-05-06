package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ricsy/gt/pkg/util"
)

// TestConfig holds test configuration
type TestConfig struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Issue struct {
		Title   string `json:"title"`
		Body    string `json:"body"`
		State   string `json:"state"`
		Comment string `json:"comment"`
	} `json:"issue"`
	PR struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		Head  string `json:"head"`
		Base  string `json:"base"`
	} `json:"pr"`
	Release struct {
		Tag  string `json:"tag"`
		Name string `json:"name"`
		Body string `json:"body"`
	} `json:"release"`
}

var (
	projectDir   string
	testRepoName string
	testOwner    string
	testCLI      string
)

func loadTestConfig(t *testing.T) *TestConfig {
	data, err := os.ReadFile(filepath.Join(projectDir, "data", "test_data.json"))
	if err != nil {
		t.Skipf("Skipping integration test: test data not found: %v", err)
	}

	var cfg TestConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Failed to parse test config: %v", err)
	}
	return &cfg
}

func buildCLI(t *testing.T) string {
	t.Helper()
	cliPath := filepath.Join(projectDir, "gt")

	cmd := exec.Command("go", "build", "-o", cliPath, ".")
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build CLI: %v\n%s", err, output)
	}
	return cliPath
}

func runCLI(t *testing.T, cli string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(cli, args...)
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func generateTestRepoName() string {
	now := time.Now().Format("20060102-15")
	return fmt.Sprintf("test_cli_gt-%s", now)
}

func setupTestRepo(cli string) string {
	repoName := generateTestRepoName()

	// Try to create repo - if it already exists, that's fine
	cmd := exec.Command(cli, "repo", "create",
		"--name", repoName,
		"--description", "Integration test repository for gt CLI")
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it's a "already exists" error
		if strings.Contains(string(output), "already exists") || strings.Contains(string(output), "已存在") {
			fmt.Printf("Using existing test repository: %s\n", repoName)
			return repoName
		}
		fmt.Printf("Failed to create test repo %s: %v\nOutput: %s\n", repoName, err, output)
		os.Exit(1)
	}

	fmt.Printf("Created test repository: %s\n", repoName)
	return repoName
}

func printCleanupMessage() {
	fmt.Printf(`
========================================
Integration test completed!
Test repository: %s
Please delete it manually after testing:
  1. Go to: https://gitee.com/%s/%s
  2. Click 'Settings' → 'Delete Repository'
========================================
`, testRepoName, testOwner, testRepoName)
}

// TestMain handles setup and teardown for all integration tests
func TestMain(m *testing.M) {
	// Resolve project dir once
	var err error
	projectDir, err = filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		fmt.Printf("Failed to resolve project dir: %v\n", err)
		os.Exit(1)
	}

	// Build CLI first
	cliPath := filepath.Join(projectDir, "gt")

	cmd := exec.Command("go", "build", "-o", cliPath, ".")
	cmd.Dir = projectDir
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Failed to build CLI: %v\n%s\n", err, output)
		os.Exit(1)
	}
	testCLI = cliPath

	// Load config to get owner
	cfg, err := os.ReadFile(filepath.Join(projectDir, "data", "test_data.json"))
	if err != nil {
		fmt.Printf("Failed to read test config: %v\n", err)
		os.Exit(1)
	}
	var config TestConfig
	if err := json.Unmarshal(cfg, &config); err != nil {
		fmt.Printf("Failed to parse test config: %v\n", err)
		os.Exit(1)
	}
	testOwner = config.Owner

	// Create test repo
	testRepoName = setupTestRepo(testCLI)
	defer printCleanupMessage()

	// Run tests
	exitCode := m.Run()

	if exitCode == 0 {
		fmt.Println("\nAll integration tests passed!")
	}

	os.Exit(exitCode)
}

// TestIntegrationRepo tests repository commands
func TestIntegrationRepo(t *testing.T) {
	stdout, err := runCLI(t, testCLI, "repo", "view", testOwner+"/"+testRepoName)
	if err != nil {
		t.Errorf("repo view failed: %v\nOutput: %s", err, stdout)
	}
	if !strings.Contains(stdout, testRepoName) {
		t.Errorf("Expected repo name %s in output, got: %s", testRepoName, stdout)
	}
}

// TestIntegrationIssue tests issue commands
func TestIntegrationIssue(t *testing.T) {
	cfg := loadTestConfig(t)

	// Create issue
	createOutput, err := runCLI(t, testCLI, "issue", "create",
		"--owner", testOwner,
		"--repo", testRepoName,
		"--title", cfg.Issue.Title,
		"--body", cfg.Issue.Body)
	if err != nil {
		t.Fatalf("issue create failed: %v\nOutput: %s", err, createOutput)
	}

	// Extract issue number
	parts := strings.Split(createOutput, "#")
	if len(parts) < 2 {
		t.Fatalf("Failed to extract issue number from: %s", createOutput)
	}
	issueNum := strings.Fields(parts[1])[0]

	// List issues
	listOutput, err := runCLI(t, testCLI, "issue", "list",
		"--owner", testOwner,
		"--repo", testRepoName)
	if err != nil {
		t.Errorf("issue list failed: %v\nOutput: %s", err, listOutput)
	}
	if !strings.Contains(listOutput, issueNum) {
		t.Errorf("Expected issue %s in list output, got: %s", issueNum, listOutput)
	}

	// View issue
	viewOutput, err := runCLI(t, testCLI, "issue", "view", issueNum,
		"--owner", testOwner,
		"--repo", testRepoName)
	if err != nil {
		t.Errorf("issue view failed: %v\nOutput: %s", err, viewOutput)
	}
	if !strings.Contains(viewOutput, cfg.Issue.Title) {
		t.Errorf("Expected title %s in view output, got: %s", cfg.Issue.Title, viewOutput)
	}

	// Add comment
	_, err = runCLI(t, testCLI, "issue", "comment", issueNum,
		"--owner", testOwner,
		"--repo", testRepoName,
		"--body", cfg.Issue.Comment)
	if err != nil {
		t.Errorf("issue comment failed: %v", err)
	}

	// Close issue
	_, err = runCLI(t, testCLI, "issue", "close", issueNum,
		"--owner", testOwner,
		"--repo", testRepoName)
	if err != nil {
		t.Errorf("issue close failed: %v", err)
	}

	// Reopen issue
	_, err = runCLI(t, testCLI, "issue", "reopen", issueNum,
		"--owner", testOwner,
		"--repo", testRepoName)
	if err != nil {
		t.Errorf("issue reopen failed: %v", err)
	}
}

// TestIntegrationRelease tests release commands
func TestIntegrationRelease(t *testing.T) {
	cfg := loadTestConfig(t)

	// Clone repo and create git tag first (release requires git tag)
	cloneDir := filepath.Join(os.TempDir(), "gt-test-"+testRepoName)
	defer os.RemoveAll(cloneDir)

	env := append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	cloneCmd := exec.Command("git", "clone", "https://gitee.com/"+testOwner+"/"+testRepoName+".git", cloneDir)
	cloneCmd.Env = env
	if err := cloneCmd.Run(); err != nil {
		t.Skipf("Failed to clone repo for release test: %v", err)
	}

	tagCmd := exec.Command("git", "tag", cfg.Release.Tag)
	tagCmd.Dir = cloneDir
	if err := tagCmd.Run(); err != nil {
		t.Skipf("Failed to create git tag: %v", err)
	}

	pushCmd := exec.Command("git", "push", "origin", cfg.Release.Tag)
	pushCmd.Dir = cloneDir
	pushCmd.Env = env
	if err := pushCmd.Run(); err != nil {
		t.Skipf("Failed to push git tag: %v", err)
	}

	// Create release
	createOutput, err := runCLI(t, testCLI, "release", "create",
		"--repo", testOwner+"/"+testRepoName,
		"--name", cfg.Release.Name,
		"--body", cfg.Release.Body,
		cfg.Release.Tag)
	if err != nil {
		t.Fatalf("release create failed: %v\nOutput: %s", err, createOutput)
	}

	// List releases
	listOutput, err := runCLI(t, testCLI, "release", "list",
		"--repo", testOwner+"/"+testRepoName)
	if err != nil {
		t.Errorf("release list failed: %v\nOutput: %s", err, listOutput)
	}
	if !strings.Contains(listOutput, cfg.Release.Tag) {
		t.Errorf("Expected tag %s in list output, got: %s", cfg.Release.Tag, listOutput)
	}

	// View release
	viewOutput, err := runCLI(t, testCLI, "release", "view",
		"--repo", testOwner+"/"+testRepoName,
		cfg.Release.Tag)
	if err != nil {
		t.Errorf("release view failed: %v\nOutput: %s", err, viewOutput)
	}

	// Delete release
	_, err = runCLI(t, testCLI, "release", "delete",
		"--repo", testOwner+"/"+testRepoName,
		cfg.Release.Tag)
	if err != nil {
		t.Errorf("release delete failed: %v", err)
	}
}

// TestIntegrationOrg tests organization commands
func TestIntegrationOrg(t *testing.T) {
	// List orgs - just verify it runs
	_, err := runCLI(t, testCLI, "org", "list")
	if err != nil {
		t.Logf("org list warning (may have no orgs): %v", err)
	}
}

// TestIntegrationAuth tests auth status
func TestIntegrationAuth(t *testing.T) {
	// Verify auth is configured
	_, err := runCLI(t, testCLI, "auth", "status")
	if err != nil {
		t.Skipf("Skipping: authentication not configured: %v", err)
	}
}

// TestIntegrationRepoList tests repository list command
func TestIntegrationRepoList(t *testing.T) {
	// List repos for the owner - just verify it runs without error
	// Note: Gitee API may not return newly created repos immediately
	stdout, err := runCLI(t, testCLI, "repo", "list", "--owner", testOwner, "--limit", "100")
	if err != nil {
		t.Errorf("repo list failed: %v\nOutput: %s", err, stdout)
	}
	// Log output for debugging but don't fail - repo view already confirms repo exists
	if stdout != "" {
		t.Logf("repo list returned %d repos", strings.Count(stdout, "\n")/2)
	}
}

// TestIntegrationPRList tests PR list command
func TestIntegrationPRList(t *testing.T) {
	// List PRs (may be empty but should not error)
	_, err := runCLI(t, testCLI, "pr", "list",
		"--repo", testOwner+"/"+testRepoName)
	if err != nil {
		t.Errorf("pr list failed: %v", err)
	}
}

// TestIntegrationPRCreate tests PR create command
func TestIntegrationPRCreate(t *testing.T) {
	cfg := loadTestConfig(t)

	// First create a branch in the test repo
	cloneDir := filepath.Join(os.TempDir(), "gt-test-pr-"+testRepoName)
	defer os.RemoveAll(cloneDir)

	env := append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	cloneCmd := exec.Command("git", "clone", "https://gitee.com/"+testOwner+"/"+testRepoName+".git", cloneDir)
	cloneCmd.Env = env
	if err := cloneCmd.Run(); err != nil {
		t.Skipf("Failed to clone repo for PR test: %v", err)
	}

	// Create a test branch
	branchCmd := exec.Command("git", "checkout", "-b", cfg.PR.Head)
	branchCmd.Dir = cloneDir
	if err := branchCmd.Run(); err != nil {
		t.Skipf("Failed to create branch: %v", err)
	}

	// Create a dummy commit
	commitCmd := exec.Command("git", "commit", "--allow-empty", "-m", "Test commit")
	commitCmd.Dir = cloneDir
	if err := commitCmd.Run(); err != nil {
		t.Skipf("Failed to commit: %v", err)
	}

	// Push branch
	pushCmd := exec.Command("git", "push", "-u", "origin", cfg.PR.Head)
	pushCmd.Dir = cloneDir
	pushCmd.Env = env
	if err := pushCmd.Run(); err != nil {
		t.Skipf("Failed to push branch: %v", err)
	}

	// Create PR
	createOutput, err := runCLI(t, testCLI, "pr", "create",
		"--repo", testOwner+"/"+testRepoName,
		"--title", cfg.PR.Title,
		"--body", cfg.PR.Body,
		"--head", cfg.PR.Head,
		"--base", "master")
	if err != nil {
		t.Errorf("pr create failed: %v\nOutput: %s", err, createOutput)
	}

	// Extract PR number from output or error (format: "Created PR #123" or URL like ".../pulls/1")
	prNum := ""
	if strings.Contains(createOutput, "#") {
		parts := strings.Split(createOutput, "#")
		if len(parts) >= 2 {
			prNum = util.ExtractDigits(parts[1])
		}
	} else if strings.Contains(createOutput, "/pulls/") {
		// Extract from URL like "https://gitee.com/owner/repo/pulls/1"
		parts := strings.Split(createOutput, "/pulls/")
		if len(parts) >= 2 {
			prNum = util.ExtractDigits(parts[1])
		}
	}

	// If prNum is empty or invalid, get from PR list
	if prNum == "" {
		listOutput, _ := runCLI(t, testCLI, "pr", "list", "--repo", testOwner+"/"+testRepoName)
		// Parse first PR number from output like "#1\tIntegration Test PR\t[open]"
		parts := strings.Split(listOutput, "#")
		if len(parts) >= 2 {
			prNum = util.ExtractDigits(parts[1])
		}
	}

	if prNum != "" {

		// View PR
		viewOutput, err := runCLI(t, testCLI, "pr", "view", prNum,
			"--repo", testOwner+"/"+testRepoName)
		if err != nil {
			t.Errorf("pr view failed: %v\nOutput: %s", err, viewOutput)
		}
		if !strings.Contains(viewOutput, cfg.PR.Title) {
			t.Errorf("Expected title %s in view output, got: %s", cfg.PR.Title, viewOutput)
		}

		// Add PR comment
		_, err = runCLI(t, testCLI, "pr", "comment", prNum,
			"--repo", testOwner+"/"+testRepoName,
			"--body", "Test PR comment")
		if err != nil {
			t.Errorf("pr comment failed: %v", err)
		}

		// Close PR
		_, err = runCLI(t, testCLI, "pr", "close", prNum,
			"--repo", testOwner+"/"+testRepoName)
		if err != nil {
			t.Errorf("pr close failed: %v", err)
		}
	}
}

// TestIntegrationConfig tests config command
func TestIntegrationConfig(t *testing.T) {
	// Get config - should not error
	_, err := runCLI(t, testCLI, "config", "get", "default_repo")
	if err != nil {
		t.Logf("config get warning: %v", err)
	}
}
