package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

type repoResolver func(string) (string, string, error)

func resolveRepoClientWith(repoValue string, resolver repoResolver) (string, string, *api.Client, error) {
	owner, repoName, err := resolver(repoValue)
	if err != nil {
		return "", "", nil, err
	}

	client, err := getClient()
	if err != nil {
		return "", "", nil, err
	}

	return owner, repoName, client, nil
}

func resolveRepoClient(repoFlag string) (string, string, *api.Client, error) {
	return resolveRepoClientWith(repoFlag, resolveRepoFlag)
}

func printRepositories(repos []api.Repository) {
	for _, repo := range repos {
		fmt.Printf("%s\t%s\n", repo.FullName, repo.Description)
	}
}

func printRepositoriesResult(repos []api.Repository, err error) error {
	if err != nil {
		return err
	}
	printRepositories(repos)
	return nil
}

func printUsersResult(users []api.User, err error) error {
	if err != nil {
		return err
	}
	printUsers(users)
	return nil
}

func createRepoFromCommand(cmd *cobra.Command) (*api.Repository, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	opts, err := buildCreateRepoOptions(cmd)
	if err != nil {
		return nil, err
	}

	repo, err := client.CreateRepo(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create repo: %w", err)
	}

	return repo, nil
}
