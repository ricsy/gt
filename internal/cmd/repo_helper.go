package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ricsy/gt/pkg/api"
	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/ricsy/gt/pkg/util"
)

func getClient() (*api.Client, error) {
	host := resolveCommandHost()
	token, err := auth.GetToken(host)
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}
	return newCommandAPIClient(host, token), nil
}

func resolveCommandHost() string {
	if commandHost != "" {
		return commandHost
	}

	if envHost := os.Getenv("GT_HOST"); envHost != "" {
		return envHost
	}

	cfg, err := config.LoadConfig()
	if err == nil && cfg.DefaultHost != "" {
		return cfg.DefaultHost
	}

	return config.DefaultHost
}

func resolveRepoFlag(repoFlag string) (owner, repoName string, err error) {
	repo := repoFlag
	if repo == "" {
		repo = os.Getenv("GT_REPO")
	}
	if repo == "" {
		repo, err = resolveRepoFromGitRemote(resolveCommandHost())
		if err != nil {
			return "", "", err
		}
	}
	if repo == "" {
		cfg, loadErr := config.LoadConfig()
		if loadErr == nil && cfg.DefaultRepo != "" {
			repo = cfg.DefaultRepo
		}
	}
	return ResolveRepo(repo)
}

// ResolveRepo parses owner/repo from a string and validates the format.
func ResolveRepo(repo string) (owner, repoName string, err error) {
	if repo == "" {
		return "", "", fmt.Errorf("repo is required: use --repo, set GT_REPO, or configure default_repo")
	}

	repo = strings.TrimSpace(repo)
	if !strings.Contains(repo, "/") {
		defaultOwner, err := resolveDefaultRepoOwner()
		if err != nil {
			return "", "", err
		}
		if defaultOwner == "" {
			return "", "", fmt.Errorf("owner is required: use owner/repo, configure default_owner, or enter repo scope mode")
		}
		owner = defaultOwner
		repoName = repo
	} else {
		owner, repoName = util.SplitOwnerRepo(repo)
	}

	if owner == "" || repoName == "" {
		return "", "", fmt.Errorf("invalid repo format: owner/repo expected")
	}
	if err := enforceRepoScopeOwner(owner); err != nil {
		return "", "", err
	}
	return owner, repoName, nil
}

func resolveRepoFromGitRemote(host string) (string, error) {
	insideWorkTree, err := gitIsInsideWorkTree()
	if err != nil || !insideWorkTree {
		return "", nil
	}

	originURL, err := gitRemoteGetURL("origin")
	if err != nil {
		return "", nil
	}

	return parseRepoFromRemoteURL(originURL, host)
}

func parseRepoFromRemoteURL(remoteURL, expectedHost string) (string, error) {
	remoteURL = strings.TrimSpace(remoteURL)
	expectedHost = strings.TrimSpace(expectedHost)
	if remoteURL == "" || expectedHost == "" {
		return "", nil
	}

	if strings.HasPrefix(remoteURL, "git@") {
		parts := strings.SplitN(strings.TrimPrefix(remoteURL, "git@"), ":", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], expectedHost) {
			return "", nil
		}
		return normalizeRemoteRepoPath(parts[1]), nil
	}

	parsed, err := url.Parse(remoteURL)
	if err != nil {
		return "", nil
	}
	if !strings.EqualFold(parsed.Hostname(), expectedHost) {
		return "", nil
	}

	return normalizeRemoteRepoPath(parsed.Path), nil
}

func normalizeRemoteRepoPath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, ".git")
	if strings.Count(path, "/") != 1 {
		return ""
	}
	return path
}
