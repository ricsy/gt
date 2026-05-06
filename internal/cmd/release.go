package cmd

import (
	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

func newReleaseCmd() *cobra.Command {
	var repoFlag string
	var nameFlag string
	var bodyFlag string
	var targetCommitFlag string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List releases for a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, repoName, err := resolveRepoFlag(repoFlag)
			if err != nil {
				return err
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			releases, err := client.ListReleases(owner, repoName)
			if err != nil {
				return err
			}

			for _, r := range releases {
				cmd.Printf("%s - %s (%s)\n", r.TagName, r.Name, r.PublishedAt)
			}
			return nil
		},
	}
	listCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")

	viewCmd := &cobra.Command{
		Use:   "view <tag>",
		Short: "View a release by tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tag := args[0]
			owner, repoName, err := resolveRepoFlag(repoFlag)
			if err != nil {
				return err
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			release, err := client.GetRelease(owner, repoName, tag)
			if err != nil {
				return err
			}

			cmd.Printf("Tag: %s\nName: %s\nBody:\n%s\nURL: %s\n", release.TagName, release.Name, release.Body, release.HtmlUrl)
			return nil
		},
	}
	viewCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")

	createCmd := &cobra.Command{
		Use:   "create <tag>",
		Short: "Create a release",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tag := args[0]
			owner, repoName, err := resolveRepoFlag(repoFlag)
			if err != nil {
				return err
			}

			name := nameFlag
			if name == "" {
				name = tag
			}

			target := targetCommitFlag
			if target == "" {
				target = "main"
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			release, err := client.CreateRelease(owner, repoName, api.CreateReleaseOptions{
				TagName:         tag,
				Name:            name,
				Body:            bodyFlag,
				TargetCommitish: target,
			})
			if err != nil {
				return err
			}

			cmd.Printf("Created release: %s\nURL: %s\n", release.TagName, release.HtmlUrl)
			return nil
		},
	}
	createCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")
	createCmd.Flags().StringVar(&nameFlag, "name", "", "Release name")
	createCmd.Flags().StringVar(&bodyFlag, "body", "", "Release body")
	createCmd.Flags().StringVar(&targetCommitFlag, "target", "main", "Target branch or commit")

	deleteCmd := &cobra.Command{
		Use:   "delete <tag>",
		Short: "Delete a release",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tag := args[0]
			owner, repoName, err := resolveRepoFlag(repoFlag)
			if err != nil {
				return err
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			if err := client.DeleteRelease(owner, repoName, tag); err != nil {
				return err
			}

			cmd.Printf("Deleted release: %s\n", tag)
			return nil
		},
	}
	deleteCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")

	cmd := &cobra.Command{
		Use:   "release",
		Short: "Manage releases",
	}
	cmd.AddCommand(listCmd, viewCmd, createCmd, deleteCmd)

	return cmd
}

func init() {
	rootCmd.AddCommand(newReleaseCmd())
}
