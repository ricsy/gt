package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/util"
)

// ResolveRepo parses owner/repo from a string and validates the format.
// Returns owner, repo, error.
func ResolveRepo(repo string) (owner, repoName string, err error) {
	if repo == "" {
		return "", "", fmt.Errorf("repo is required")
	}
	owner, repoName = util.SplitOwnerRepo(repo)
	if owner == "" || repoName == "" {
		return "", "", fmt.Errorf("invalid repo format: owner/repo expected")
	}
	return owner, repoName, nil
}
