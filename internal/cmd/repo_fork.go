package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var repoForkOpts struct {
	Repo    string
	Sort    string
	Page    int
	PerPage int
}

var repoForkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Manage repository forks",
}

var repoForkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repository forks",
	RunE:  repoForkListCommand,
}

var repoForkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Fork a repository",
	RunE:  repoForkCreateCommand,
}

func addRepoForkFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&repoForkOpts.Repo, "repo", "", "Repository (owner/repo)")
}

func repoForkListCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoForkOpts.Repo)
	if err != nil {
		return err
	}
	forks, err := client.ListForks(owner, repoName, api.ListForksOptions{
		Sort:    repoForkOpts.Sort,
		Page:    repoForkOpts.Page,
		PerPage: repoForkOpts.PerPage,
	})
	if err != nil {
		return fmt.Errorf("failed to list forks: %w", err)
	}
	for _, fork := range forks {
		cmd.Printf("%s\n", fork.FullName)
	}
	return nil
}

func repoForkCreateCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoForkOpts.Repo)
	if err != nil {
		return err
	}
	fork, err := client.ForkRepository(owner, repoName)
	if err != nil {
		return fmt.Errorf("failed to fork repository: %w", err)
	}
	cmd.Printf("Forked: %s\n", fork.HTMLURL)
	return nil
}
