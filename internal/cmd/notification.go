package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	notificationRepo          string
	notificationOwner         string
	notificationUnread        bool
	notificationParticipating bool
	notificationType          string
	notificationSince         string
	notificationBefore        string
	notificationIDs           string
	notificationMessageID     string
	notificationUserID        int64
	notificationContent       string
)

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Manage notifications",
	Long:  `Commands for managing Gitee notifications`,
}

var repoNotificationCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repo notifications",
	Long:  `Commands for managing repository notifications`,
}

var threadCmd = &cobra.Command{
	Use:   "thread",
	Short: "Manage notification threads",
	Long:  `Commands for managing notification threads`,
}

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Manage messages",
	Long:  `Commands for managing Gitee messages`,
}

var notificationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List notifications for a repo",
	RunE:  repoNotificationList,
}

var notificationMarkReadCmd = &cobra.Command{
	Use:   "mark-read",
	Short: "Mark repo notifications as read",
	RunE:  repoNotificationMarkRead,
}

var threadListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notifications",
	RunE:  threadList,
}

var threadMarkAllReadCmd = &cobra.Command{
	Use:   "mark-all-read",
	Short: "Mark all notifications as read",
	RunE:  threadMarkAllRead,
}

var threadViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View a notification",
	Args:  cobra.ExactArgs(1),
	RunE:  threadView,
}

var threadMarkReadCmd = &cobra.Command{
	Use:   "mark-read <id>",
	Short: "Mark a notification as read",
	Args:  cobra.ExactArgs(1),
	RunE:  threadMarkRead,
}

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Get notification count",
	RunE:  notificationCount,
}

var messageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List messages",
	RunE:  messageList,
}

var messageCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Send a message",
	RunE:  messageCreate,
}

var messageMarkAllReadCmd = &cobra.Command{
	Use:   "mark-all-read",
	Short: "Mark all messages as read",
	RunE:  messageMarkAllRead,
}

var messageViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View a message",
	Args:  cobra.ExactArgs(1),
	RunE:  messageView,
}

var messageMarkReadCmd = &cobra.Command{
	Use:   "mark-read <id>",
	Short: "Mark a message as read",
	Args:  cobra.ExactArgs(1),
	RunE:  messageMarkRead,
}

func init() {
	notificationCmd.AddCommand(repoNotificationCmd, threadCmd, messageCmd, countCmd)

	repoNotificationCmd.AddCommand(notificationListCmd, notificationMarkReadCmd)
	threadCmd.AddCommand(threadListCmd, threadMarkAllReadCmd, threadViewCmd, threadMarkReadCmd)
	messageCmd.AddCommand(messageListCmd, messageCreateCmd, messageMarkAllReadCmd, messageViewCmd, messageMarkReadCmd)

	rootCmd.AddCommand(notificationCmd)

	notificationListCmd.Flags().StringVarP(&notificationRepo, "repo", "r", "", "Repository name (required)")
	notificationListCmd.Flags().StringVarP(&notificationOwner, "owner", "o", "", "Owner name (required)")
	notificationListCmd.Flags().BoolVar(&notificationUnread, "unread", false, "Only show unread notifications")
	notificationListCmd.Flags().BoolVar(&notificationParticipating, "participating", false, "Only show participating notifications")
	notificationListCmd.Flags().StringVar(&notificationType, "type", "", "Filter by type: all, event, referer")
	notificationListCmd.Flags().StringVar(&notificationSince, "since", "", "Notifications after this time (ISO 8601)")
	notificationListCmd.Flags().StringVar(&notificationBefore, "before", "", "Notifications before this time (ISO 8601)")
	_ = notificationListCmd.MarkFlagRequired("repo")
	_ = notificationListCmd.MarkFlagRequired("owner")

	notificationMarkReadCmd.Flags().StringVarP(&notificationRepo, "repo", "r", "", "Repository name (required)")
	notificationMarkReadCmd.Flags().StringVarP(&notificationOwner, "owner", "o", "", "Owner name (required)")
	notificationMarkReadCmd.Flags().StringVar(&notificationIDs, "ids", "", "Specific notification IDs to mark (comma-separated)")
	_ = notificationMarkReadCmd.MarkFlagRequired("repo")
	_ = notificationMarkReadCmd.MarkFlagRequired("owner")

	threadListCmd.Flags().BoolVar(&notificationUnread, "unread", false, "Only show unread notifications")
	threadListCmd.Flags().BoolVar(&notificationParticipating, "participating", false, "Only show participating notifications")
	threadListCmd.Flags().StringVar(&notificationType, "type", "", "Filter by type: all, event, referer")
	threadListCmd.Flags().StringVar(&notificationSince, "since", "", "Notifications after this time (ISO 8601)")
	threadListCmd.Flags().StringVar(&notificationBefore, "before", "", "Notifications before this time (ISO 8601)")

	threadMarkAllReadCmd.Flags().StringVar(&notificationIDs, "ids", "", "Specific notification IDs to mark (comma-separated)")

	threadViewCmd.Flags().StringVar(&notificationMessageID, "id", "", "Notification ID")

	threadMarkReadCmd.Flags().StringVar(&notificationMessageID, "id", "", "Notification ID")

	messageListCmd.Flags().BoolVar(&notificationUnread, "unread", false, "Only show unread messages")
	messageListCmd.Flags().StringVar(&notificationSince, "since", "", "Messages after this time (ISO 8601)")
	messageListCmd.Flags().StringVar(&notificationBefore, "before", "", "Messages before this time (ISO 8601)")

	messageCreateCmd.Flags().Int64Var(&notificationUserID, "user-id", 0, "User ID to send message to (required)")
	messageCreateCmd.Flags().StringVarP(&notificationContent, "content", "c", "", "Message content (required)")
	_ = messageCreateCmd.MarkFlagRequired("user-id")
	_ = messageCreateCmd.MarkFlagRequired("content")

	messageViewCmd.Flags().StringVar(&notificationMessageID, "id", "", "Message ID")

	messageMarkReadCmd.Flags().StringVar(&notificationMessageID, "id", "", "Message ID")
}

func repoNotificationList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	opts := api.ListRepoNotificationsOptions{
		Type:   notificationType,
		Since:  notificationSince,
		Before: notificationBefore,
	}
	if notificationUnread {
		opts.Unread = &notificationUnread
	}
	if notificationParticipating {
		opts.Participating = &notificationParticipating
	}

	result, err := client.ListRepoNotifications(notificationOwner, notificationRepo, opts)
	if err != nil {
		return fmt.Errorf("failed to list repo notifications: %w", err)
	}

	if result.TotalCount == 0 {
		fmt.Println("No notifications found")
		return nil
	}

	fmt.Printf("Total: %d notifications\n", result.TotalCount)
	for _, n := range result.List {
		unreadStr := ""
		if n.Unread {
			unreadStr = " [UNREAD]"
		}
		fmt.Printf("#%d %s%s - %s\n", n.ID, n.Subject.Title, unreadStr, n.Repository.FullName)
	}
	return nil
}

func repoNotificationMarkRead(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	opts := api.MarkNotificationsReadOptions{
		IDs: notificationIDs,
	}

	err = client.MarkRepoNotificationsRead(notificationOwner, notificationRepo, opts)
	if err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}

	fmt.Println("Notifications marked as read")
	return nil
}

func threadList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	opts := api.ListNotificationsOptions{
		Type:   notificationType,
		Since:  notificationSince,
		Before: notificationBefore,
	}
	if notificationUnread {
		opts.Unread = &notificationUnread
	}
	if notificationParticipating {
		opts.Participating = &notificationParticipating
	}

	result, err := client.ListNotifications(opts)
	if err != nil {
		return fmt.Errorf("failed to list notifications: %w", err)
	}

	if result.TotalCount == 0 {
		fmt.Println("No notifications found")
		return nil
	}

	fmt.Printf("Total: %d notifications\n", result.TotalCount)
	for _, n := range result.List {
		unreadStr := ""
		if n.Unread {
			unreadStr = " [UNREAD]"
		}
		fmt.Printf("#%d %s%s\n", n.ID, n.Subject.Title, unreadStr)
	}
	return nil
}

func threadMarkAllRead(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	opts := api.MarkNotificationsReadOptions{
		IDs: notificationIDs,
	}

	err = client.MarkAllNotificationsRead(opts)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	fmt.Println("All notifications marked as read")
	return nil
}

func threadView(cmd *cobra.Command, args []string) error {
	id := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	notification, err := client.GetNotification(id)
	if err != nil {
		return fmt.Errorf("failed to get notification: %w", err)
	}

	fmt.Printf("ID: %d\n", notification.ID)
	fmt.Printf("Type: %s\n", notification.Type)
	fmt.Printf("Unread: %v\n", notification.Unread)
	fmt.Printf("Subject: %s\n", notification.Subject.Title)
	fmt.Printf("Repository: %s\n", notification.Repository.FullName)
	if notification.Subject.URL != "" {
		fmt.Printf("URL: %s\n", notification.Subject.URL)
	}
	if notification.HTMLURL != "" {
		fmt.Printf("HTML URL: %s\n", notification.HTMLURL)
	}
	fmt.Printf("Updated: %s\n", notification.UpdatedAt)
	if notification.Actor.Name != "" {
		fmt.Printf("Actor: %s\n", notification.Actor.Name)
	}
	return nil
}

func threadMarkRead(cmd *cobra.Command, args []string) error {
	id := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.MarkNotificationRead(id)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	fmt.Printf("Notification %s marked as read\n", id)
	return nil
}

func notificationCount(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	var unread *bool
	if notificationUnread {
		unread = &notificationUnread
	}

	count, err := client.GetNotificationCount(unread)
	if err != nil {
		return fmt.Errorf("failed to get notification count: %w", err)
	}

	fmt.Printf("Total: %d\n", count.TotalCount)
	fmt.Printf("Notifications: %d\n", count.NotificationCount)
	fmt.Printf("Messages: %d\n", count.MessageCount)
	return nil
}

func messageList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	opts := api.ListMessagesOptions{
		Since:  notificationSince,
		Before: notificationBefore,
	}
	if notificationUnread {
		opts.Unread = &notificationUnread
	}

	result, err := client.ListMessages(opts)
	if err != nil {
		return fmt.Errorf("failed to list messages: %w", err)
	}

	if result.TotalCount == 0 {
		fmt.Println("No messages found")
		return nil
	}

	fmt.Printf("Total: %d messages\n", result.TotalCount)
	for _, m := range result.List {
		unreadStr := ""
		if m.Unread {
			unreadStr = " [UNREAD]"
		}
		fmt.Printf("#%d %s%s - %s\n", m.ID, m.Sender.Name, unreadStr, m.Content)
	}
	return nil
}

func messageCreate(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	msg, err := client.CreateMessage(api.CreateMessageOptions{
		UserID:  notificationUserID,
		Content: notificationContent,
	})
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	fmt.Printf("Message sent: #%d\n", msg.ID)
	return nil
}

func messageMarkAllRead(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.MarkAllMessagesRead()
	if err != nil {
		return fmt.Errorf("failed to mark all messages as read: %w", err)
	}

	fmt.Println("All messages marked as read")
	return nil
}

func messageView(cmd *cobra.Command, args []string) error {
	id := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	msg, err := client.GetMessage(id)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	fmt.Printf("ID: %d\n", msg.ID)
	fmt.Printf("Unread: %v\n", msg.Unread)
	fmt.Printf("Content: %s\n", msg.Content)
	fmt.Printf("Sender: %s\n", msg.Sender.Name)
	if msg.URL != "" {
		fmt.Printf("URL: %s\n", msg.URL)
	}
	if msg.HTMLURL != "" {
		fmt.Printf("HTML URL: %s\n", msg.HTMLURL)
	}
	fmt.Printf("Updated: %s\n", msg.UpdatedAt)
	return nil
}

func messageMarkRead(cmd *cobra.Command, args []string) error {
	id := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.MarkMessageRead(id)
	if err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	fmt.Printf("Message %s marked as read\n", id)
	return nil
}
