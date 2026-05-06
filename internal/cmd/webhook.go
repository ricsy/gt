package cmd

import (
	"fmt"
	"strconv"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

func newWebhookCmd() *cobra.Command {
	var repoFlag string
	var page, perPage int
	var webhookURL, webhookTitle, password string
	var encryptionType int
	var pushEventsFlag, tagPushEventsFlag, issuesEventsFlag, noteEventsFlag, mergeRequestsEventsFlag bool

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

			webhooks, err := client.ListWebhooks(owner, repoName, api.ListWebhooksOptions{
				Page:    page,
				PerPage: perPage,
			})
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
	listCmd.Flags().IntVar(&page, "page", 0, "Page number")
	listCmd.Flags().IntVar(&perPage, "per-page", 0, "Items per page")

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

			if webhookURL == "" {
				return fmt.Errorf("url is required: use --url")
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			pushEvents := api.BoolPtr(false)
			tagPushEvents := api.BoolPtr(false)
			issuesEvents := api.BoolPtr(false)
			noteEvents := api.BoolPtr(false)
			mergeRequestsEvents := api.BoolPtr(false)

			if pushEventsFlag {
				pushEvents = api.BoolPtr(true)
			}
			if tagPushEventsFlag {
				tagPushEvents = api.BoolPtr(true)
			}
			if issuesEventsFlag {
				issuesEvents = api.BoolPtr(true)
			}
			if noteEventsFlag {
				noteEvents = api.BoolPtr(true)
			}
			if mergeRequestsEventsFlag {
				mergeRequestsEvents = api.BoolPtr(true)
			}

			webhook, err := client.CreateWebhook(owner, repoName, api.CreateWebhookOptions{
				URL:                 webhookURL,
				Title:               webhookTitle,
				EncryptionType:      encryptionType,
				Password:            password,
				PushEvents:          pushEvents,
				TagPushEvents:       tagPushEvents,
				IssuesEvents:        issuesEvents,
				NoteEvents:          noteEvents,
				MergeRequestsEvents: mergeRequestsEvents,
			})
			if err != nil {
				return err
			}

			cmd.Printf("Created webhook: %d - %s\n", webhook.ID, webhook.URL)
			return nil
		},
	}
	createCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")
	createCmd.Flags().StringVar(&webhookURL, "url", "", "Webhook URL (required)")
	createCmd.Flags().StringVar(&webhookTitle, "title", "", "Webhook title")
	createCmd.Flags().IntVar(&encryptionType, "encryption-type", 0, "Encryption type (0: 密码, 1: 签名密钥)")
	createCmd.Flags().StringVar(&password, "password", "", "Password")
	createCmd.Flags().BoolVar(&pushEventsFlag, "push-events", false, "Trigger on push events")
	createCmd.Flags().BoolVar(&tagPushEventsFlag, "tag-push-events", false, "Trigger on tag push events")
	createCmd.Flags().BoolVar(&issuesEventsFlag, "issues-events", false, "Trigger on issues events")
	createCmd.Flags().BoolVar(&noteEventsFlag, "note-events", false, "Trigger on note events")
	createCmd.Flags().BoolVar(&mergeRequestsEventsFlag, "merge-requests-events", false, "Trigger on merge requests events")

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
