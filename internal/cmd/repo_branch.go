package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var repoBranchOpts struct {
	Repo      string
	Sort      string
	Direction string
	Page      int
	PerPage   int
	Refs      string
}

var repoBranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Manage repository branches",
}

var repoBranchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repository branches",
	RunE:  repoBranchListCommand,
}

var repoBranchViewCmd = &cobra.Command{
	Use:   "view <branch>",
	Short: "View a repository branch",
	Args:  cobra.ExactArgs(1),
	RunE:  repoBranchViewCommand,
}

var repoBranchCreateCmd = &cobra.Command{
	Use:   "create <branch>",
	Short: "Create a repository branch",
	Args:  cobra.ExactArgs(1),
	RunE:  repoBranchCreateCommand,
}

var repoBranchProtectCmd = &cobra.Command{
	Use:   "protect <branch>",
	Short: "Protect a repository branch",
	Args:  cobra.ExactArgs(1),
	RunE:  repoBranchProtectCommand,
}

var repoBranchUnprotectCmd = &cobra.Command{
	Use:   "unprotect <branch>",
	Short: "Remove branch protection",
	Args:  cobra.ExactArgs(1),
	RunE:  repoBranchUnprotectCommand,
}

func addRepoBranchRepoFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&repoBranchOpts.Repo, "repo", "", "Repository (owner/repo)")
}

func repoBranchListCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoBranchOpts.Repo)
	if err != nil {
		return err
	}
	branches, err := client.ListBranches(owner, repoName, api.ListBranchesOptions{
		Sort:      repoBranchOpts.Sort,
		Direction: repoBranchOpts.Direction,
		Page:      repoBranchOpts.Page,
		PerPage:   repoBranchOpts.PerPage,
	})
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}
	for _, branch := range branches {
		printBranch(branch.Name, branch.Commit.SHA, branch.Protected)
	}
	return nil
}

func repoBranchViewCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoBranchOpts.Repo)
	if err != nil {
		return err
	}
	branch, err := client.GetBranch(owner, repoName, args[0])
	if err != nil {
		return fmt.Errorf("failed to get branch: %w", err)
	}
	printBranch(branch.Name, branch.Commit.SHA, branch.Protected)
	return nil
}

func repoBranchCreateCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoBranchOpts.Repo)
	if err != nil {
		return err
	}
	branch, err := client.CreateBranch(owner, repoName, api.CreateBranchOptions{
		Refs:       repoBranchOpts.Refs,
		BranchName: args[0],
	})
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	printBranch(branch.Name, branch.Commit.SHA, branch.Protected)
	return nil
}

func repoBranchProtectCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoBranchOpts.Repo)
	if err != nil {
		return err
	}
	branch, err := client.ProtectBranch(owner, repoName, args[0])
	if err != nil {
		return fmt.Errorf("failed to protect branch: %w", err)
	}
	printBranch(branch.Name, branch.Commit.SHA, branch.Protected)
	return nil
}

func repoBranchUnprotectCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoBranchOpts.Repo)
	if err != nil {
		return err
	}
	if err := client.UnprotectBranch(owner, repoName, args[0]); err != nil {
		return fmt.Errorf("failed to unprotect branch: %w", err)
	}
	fmt.Printf("Unprotected branch: %s\n", args[0])
	return nil
}

func printBranch(name, commit string, protected bool) {
	protection := "unprotected"
	if protected {
		protection = "protected"
	}
	fmt.Printf("%s\t%s\t%s\n", name, commit, protection)
}
