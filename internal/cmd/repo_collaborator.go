package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var repoCollaboratorOpts struct {
	Repo string
}

var repoCollaboratorAddOpts struct {
	Permission string
}

var repoCollaboratorCmd = &cobra.Command{
	Use:   "collaborator",
	Short: "Manage repository collaborators",
}

var repoCollaboratorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repository collaborators",
	RunE:  repoCollaboratorListCommand,
}

var repoCollaboratorViewCmd = &cobra.Command{
	Use:   "view <user>",
	Short: "Check if user is a collaborator",
	Args:  cobra.ExactArgs(1),
	RunE:  repoCollaboratorViewCommand,
}

var repoCollaboratorPermCmd = &cobra.Command{
	Use:   "perm <user>",
	Short: "Get collaborator permission",
	Args:  cobra.ExactArgs(1),
	RunE:  repoCollaboratorPermCommand,
}

var repoCollaboratorAddCmd = &cobra.Command{
	Use:   "add <user>",
	Short: "Add a collaborator to repository",
	Args:  cobra.ExactArgs(1),
	RunE:  repoCollaboratorAddCommand,
}

var repoCollaboratorRemoveCmd = &cobra.Command{
	Use:   "remove <user>",
	Short: "Remove a collaborator from repository",
	Args:  cobra.ExactArgs(1),
	RunE:  repoCollaboratorRemoveCommand,
}

func addRepoFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&repoCollaboratorOpts.Repo, "repo", "", "Repository (owner/repo)")
}

func repoCollaboratorListCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoCollaboratorOpts.Repo)
	if err != nil {
		return err
	}
	collabs, err := client.ListCollaborators(owner, repoName)
	if err != nil {
		return fmt.Errorf("failed to list collaborators: %w", err)
	}
	for _, collaborator := range collabs {
		cmd.Printf("%s (%s)\n", collaborator.Login, collaborator.Name)
	}
	return nil
}

func repoCollaboratorViewCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoCollaboratorOpts.Repo)
	if err != nil {
		return err
	}
	collab, err := client.GetCollaborator(owner, repoName, args[0])
	if err != nil {
		return fmt.Errorf("failed to get collaborator: %w", err)
	}
	if collab.Login == "" && collab.Name == "" {
		cmd.Printf("%s is a collaborator\n", args[0])
		return nil
	}
	cmd.Printf("%s (%s)\n", collab.Login, collab.Name)
	return nil
}

func repoCollaboratorPermCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoCollaboratorOpts.Repo)
	if err != nil {
		return err
	}
	perm, err := client.GetCollaboratorPermission(owner, repoName, args[0])
	if err != nil {
		return fmt.Errorf("failed to get collaborator permission: %w", err)
	}
	cmd.Printf("Permission: %s\n", perm.Permission)
	return nil
}

func repoCollaboratorAddCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoCollaboratorOpts.Repo)
	if err != nil {
		return err
	}
	if err := client.AddCollaborator(owner, repoName, args[0], repoCollaboratorAddOpts.Permission); err != nil {
		return fmt.Errorf("failed to add collaborator: %w", err)
	}
	cmd.Printf("Added collaborator: %s\n", args[0])
	return nil
}

func repoCollaboratorRemoveCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(repoCollaboratorOpts.Repo)
	if err != nil {
		return err
	}
	if err := client.RemoveCollaborator(owner, repoName, args[0]); err != nil {
		return fmt.Errorf("failed to remove collaborator: %w", err)
	}
	cmd.Printf("Removed collaborator: %s\n", args[0])
	return nil
}
