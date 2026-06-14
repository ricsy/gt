package cmd

import (
	"fmt"
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
