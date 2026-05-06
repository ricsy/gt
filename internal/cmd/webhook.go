package cmd

import (
	"strconv"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

func newWebhookCmd() *cobra.Command {
	var repoFlag string

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List webhooks for a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, repoName, err := resolveRepoFlag(repoFlag)
			if err != nil {
				return err
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			webhooks, err := client.ListWebhooks(owner, repoName)
			if err != nil {
				return err
			}

			for _, w := range webhooks {
				cmd.Printf("%d - %s (%s)\n", w.ID, w.URL, w.Title)
			}
			return nil
		},
	}
	listCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")

	viewCmd := &cobra.Command{
		Use:   "view <id>",
		Short: "View a webhook by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			owner, repoName, err := resolveRepoFlag(repoFlag)
			if err != nil {
				return err
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			webhook, err := client.GetWebhook(owner, repoName, id)
			if err != nil {
				return err
			}

			cmd.Printf("ID: %d\nURL: %s\nTitle: %s\n", webhook.ID, webhook.URL, webhook.Title)
			return nil
		},
	}
	viewCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, repoName, err := resolveRepoFlag(repoFlag)
			if err != nil {
				return err
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			webhook, err := client.CreateWebhook(owner, repoName, api.CreateWebhookOptions{
				URL: "https://example.com/webhook",
			})
			if err != nil {
				return err
			}

			cmd.Printf("Created webhook: %d - %s\n", webhook.ID, webhook.URL)
			return nil
		},
	}
	createCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")

	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			owner, repoName, err := resolveRepoFlag(repoFlag)
			if err != nil {
				return err
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			if err := client.DeleteWebhook(owner, repoName, id); err != nil {
				return err
			}

			cmd.Printf("Deleted webhook: %d\n", id)
			return nil
		},
	}
	deleteCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")

	testCmd := &cobra.Command{
		Use:   "test <id>",
		Short: "Test a webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}

			owner, repoName, err := resolveRepoFlag(repoFlag)
			if err != nil {
				return err
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			if err := client.TestWebhook(owner, repoName, id); err != nil {
				return err
			}

			cmd.Printf("Test webhook: %d\n", id)
			return nil
		},
	}
	testCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")

	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Manage webhooks",
	}
	cmd.AddCommand(listCmd, viewCmd, createCmd, deleteCmd, testCmd)

	return cmd
}

func init() {
	rootCmd.AddCommand(newWebhookCmd())
}
