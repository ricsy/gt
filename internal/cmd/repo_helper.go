package cmd

import (
	"fmt"
	"os"

	"github.com/ricsy/gt/pkg/api"
	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/ricsy/gt/pkg/util"
)

func getClient() (*api.Client, error) {
	token, err := auth.GetToken(config.DefaultHost)
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}
	return newCommandAPIClient(config.DefaultHost, token), nil
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
	owner, repoName = util.SplitOwnerRepo(repo)
	if owner == "" || repoName == "" {
		return "", "", fmt.Errorf("invalid repo format: owner/repo expected")
	}
	return owner, repoName, nil
}
