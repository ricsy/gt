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
	listCmd.Flags().IntVar(&perPage, "per-page", 0, "Items per page (max 100)")

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

			opts := buildWebhookCreateOpts(webhookURL, webhookTitle, encryptionType, password,
				pushEventsFlag, tagPushEventsFlag, issuesEventsFlag, noteEventsFlag, mergeRequestsEventsFlag)
			webhook, err := client.CreateWebhook(owner, repoName, opts)
			if err != nil {
				return err
			}

			cmd.Printf("Created webhook: %d - %s\n", webhook.ID, webhook.URL)
			return nil
		},
	}
	createCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")
	createCmd.Flags().StringVar(&webhookURL, "url", "", "Webhook URL (required)")
	registerWebhookFlags(createCmd, &webhookURL, &webhookTitle, &encryptionType, &password,
		&pushEventsFlag, &tagPushEventsFlag, &issuesEventsFlag, &noteEventsFlag, &mergeRequestsEventsFlag)

	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a webhook",
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

			opts := buildWebhookUpdateOpts(webhookURL, webhookTitle, encryptionType, password,
				pushEventsFlag, tagPushEventsFlag, issuesEventsFlag, noteEventsFlag, mergeRequestsEventsFlag)
			webhook, err := client.UpdateWebhook(owner, repoName, id, opts)
			if err != nil {
				return err
			}

			cmd.Printf("Updated webhook: %d - %s\n", webhook.ID, webhook.URL)
			return nil
		},
	}
	updateCmd.Flags().StringVar(&repoFlag, "repo", "", "Repository (owner/repo)")
	registerWebhookFlags(updateCmd, &webhookURL, &webhookTitle, &encryptionType, &password,
		&pushEventsFlag, &tagPushEventsFlag, &issuesEventsFlag, &noteEventsFlag, &mergeRequestsEventsFlag)

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
	cmd.AddCommand(listCmd, viewCmd, createCmd, updateCmd, deleteCmd, testCmd)

	return cmd
}

func buildWebhookCreateOpts(url, title string, encType int, pwd string,
	pushEvents, tagPushEvents, issuesEvents, noteEvents, mergeRequestsEvents bool) api.CreateWebhookOptions {
	return api.CreateWebhookOptions{
		URL:                 url,
		Title:               title,
		EncryptionType:      encType,
		Password:            pwd,
		PushEvents:          api.BoolPtr(pushEvents),
		TagPushEvents:       api.BoolPtr(tagPushEvents),
		IssuesEvents:        api.BoolPtr(issuesEvents),
		NoteEvents:          api.BoolPtr(noteEvents),
		MergeRequestsEvents: api.BoolPtr(mergeRequestsEvents),
	}
}

func buildWebhookUpdateOpts(url, title string, encType int, pwd string,
	pushEvents, tagPushEvents, issuesEvents, noteEvents, mergeRequestsEvents bool) api.UpdateWebhookOptions {
	return api.UpdateWebhookOptions{
		URL:                 url,
		Title:               title,
		EncryptionType:      encType,
		Password:            pwd,
		PushEvents:          api.BoolPtr(pushEvents),
		TagPushEvents:       api.BoolPtr(tagPushEvents),
		IssuesEvents:        api.BoolPtr(issuesEvents),
		NoteEvents:          api.BoolPtr(noteEvents),
		MergeRequestsEvents: api.BoolPtr(mergeRequestsEvents),
	}
}

func registerWebhookFlags(cmd *cobra.Command, webhookURL, webhookTitle *string, encryptionType *int, password *string,
	pushEventsFlag, tagPushEventsFlag, issuesEventsFlag, noteEventsFlag, mergeRequestsEventsFlag *bool) {
	cmd.Flags().StringVar(webhookURL, "url", "", "Webhook URL")
	cmd.Flags().StringVar(webhookTitle, "title", "", "Webhook title (max 191 chars)")
	cmd.Flags().IntVar(encryptionType, "encryption-type", 0, "Encryption type (0=secret, 1=signature)")
	cmd.Flags().StringVar(password, "password", "", "Password")
	cmd.Flags().BoolVar(pushEventsFlag, "push-events", false, "Trigger on push events")
	cmd.Flags().BoolVar(tagPushEventsFlag, "tag-push-events", false, "Trigger on tag push events")
	cmd.Flags().BoolVar(issuesEventsFlag, "issues-events", false, "Trigger on issues events")
	cmd.Flags().BoolVar(noteEventsFlag, "note-events", false, "Trigger on note events")
	cmd.Flags().BoolVar(mergeRequestsEventsFlag, "merge-requests-events", false, "Trigger on merge requests events")
}

func init() {
	rootCmd.AddCommand(newWebhookCmd())
}
