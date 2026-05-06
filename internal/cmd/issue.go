package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	issueRepo  string
	issueOwner string
	issueState string
	issueLimit int
	issueTitle string
	issueBody  string
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Long:  `Commands for managing Gitee issues`,
}

var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	RunE:  issueList,
}

var issueViewCmd = &cobra.Command{
	Use:   "view <number>",
	Short: "View an issue",
	Args:  cobra.ExactArgs(1),
	RunE:  issueView,
}

var issueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an issue",
	RunE:  issueCreate,
}

var issueCloseCmd = &cobra.Command{
	Use:   "close <number>",
	Short: "Close an issue",
	Args:  cobra.ExactArgs(1),
	RunE:  issueClose,
}

var issueReopenCmd = &cobra.Command{
	Use:   "reopen <number>",
	Short: "Reopen a closed issue",
	Args:  cobra.ExactArgs(1),
	RunE:  issueReopen,
}

var issueCommentCmd = &cobra.Command{
	Use:   "comment <number>",
	Short: "Add a comment to an issue",
	Args:  cobra.ExactArgs(1),
	RunE:  issueAddComment,
}

func init() {
	issueCmd.AddCommand(issueListCmd, issueViewCmd, issueCreateCmd, issueCloseCmd, issueReopenCmd, issueCommentCmd)

	issueListCmd.Flags().StringVarP(&issueRepo, "repo", "r", "", "Repository name (required)")
	issueListCmd.Flags().StringVarP(&issueOwner, "owner", "o", "", "Owner name (required)")
	issueListCmd.Flags().StringVar(&issueState, "state", api.StateOpen, "Filter by state (open, closed, all)")
	issueListCmd.Flags().IntVarP(&issueLimit, "limit", "l", 10, "Maximum number of issues to list")
	_ = issueListCmd.MarkFlagRequired("repo")
	_ = issueListCmd.MarkFlagRequired("owner")

	issueViewCmd.Flags().StringVarP(&issueRepo, "repo", "r", "", "Repository name (required)")
	issueViewCmd.Flags().StringVarP(&issueOwner, "owner", "o", "", "Owner name (required)")
	_ = issueViewCmd.MarkFlagRequired("repo")
	_ = issueViewCmd.MarkFlagRequired("owner")

	issueCreateCmd.Flags().StringVarP(&issueRepo, "repo", "r", "", "Repository name (required)")
	issueCreateCmd.Flags().StringVarP(&issueOwner, "owner", "o", "", "Owner name (required)")
	issueCreateCmd.Flags().StringVarP(&issueTitle, "title", "t", "", "Issue title (required)")
	issueCreateCmd.Flags().StringVarP(&issueBody, "body", "b", "", "Issue body")
	_ = issueCreateCmd.MarkFlagRequired("repo")
	_ = issueCreateCmd.MarkFlagRequired("owner")
	_ = issueCreateCmd.MarkFlagRequired("title")

	issueCloseCmd.Flags().StringVarP(&issueRepo, "repo", "r", "", "Repository name (required)")
	issueCloseCmd.Flags().StringVarP(&issueOwner, "owner", "o", "", "Owner name (required)")
	_ = issueCloseCmd.MarkFlagRequired("repo")
	_ = issueCloseCmd.MarkFlagRequired("owner")

	issueReopenCmd.Flags().StringVarP(&issueRepo, "repo", "r", "", "Repository name (required)")
	issueReopenCmd.Flags().StringVarP(&issueOwner, "owner", "o", "", "Owner name (required)")
	_ = issueReopenCmd.MarkFlagRequired("repo")
	_ = issueReopenCmd.MarkFlagRequired("owner")

	issueCommentCmd.Flags().StringVarP(&issueRepo, "repo", "r", "", "Repository name (required)")
	issueCommentCmd.Flags().StringVarP(&issueOwner, "owner", "o", "", "Owner name (required)")
	issueCommentCmd.Flags().StringVarP(&issueBody, "body", "b", "", "Comment body (required)")
	_ = issueCommentCmd.MarkFlagRequired("repo")
	_ = issueCommentCmd.MarkFlagRequired("owner")
	_ = issueCommentCmd.MarkFlagRequired("body")

	rootCmd.AddCommand(issueCmd)
}

func issueList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	issues, err := client.ListIssues(issueOwner, issueRepo, issueState, 1, issueLimit)
	if err != nil {
		return fmt.Errorf("failed to list issues: %w", err)
	}

	if len(issues) == 0 {
		fmt.Println("No issues found")
		return nil
	}

	for _, issue := range issues {
		fmt.Printf("#%d [%s] %s (by %s)\n", issue.Number, issue.State, issue.Title, issue.User.Login)
	}
	return nil
}

func issueView(cmd *cobra.Command, args []string) error {
	number := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	issue, err := client.GetIssue(issueOwner, issueRepo, number)
	if err != nil {
		return fmt.Errorf("failed to get issue: %w", err)
	}

	fmt.Printf("#%d %s\n", issue.Number, issue.Title)
	fmt.Printf("State: %s | Comments: %d\n", issue.State, issue.Comments)
	fmt.Printf("Author: %s | Created: %s\n", issue.User.Login, issue.CreatedAt)
	fmt.Println("\n--- Body ---")
	if issue.Body != "" {
		fmt.Println(issue.Body)
	} else {
		fmt.Println("(No description)")
	}

	if issue.Comments > 0 {
		comments, err := client.ListIssueComments(issueOwner, issueRepo, number)
		if err == nil && len(comments) > 0 {
			fmt.Println("--- Comments ---")
			for _, c := range comments {
				fmt.Printf("[%s] %s: %s", c.CreatedAt, c.User.Login, c.Body)
			}
		}
	}

	return nil
}

func issueCreate(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	issue, err := client.CreateIssue(issueOwner, issueRepo, issueTitle, issueBody)
	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	fmt.Printf("Issue created successfully: #%d %s\n", issue.Number, issue.Title)
	fmt.Printf("URL: %s\n", issue.HTMLURL)
	return nil
}

func issueClose(cmd *cobra.Command, args []string) error {
	number := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	issue, err := client.CloseIssue(issueOwner, issueRepo, number)
	if err != nil {
		return fmt.Errorf("failed to close issue: %w", err)
	}

	fmt.Printf("Issue #%d closed successfully\n", issue.Number)
	return nil
}

func issueReopen(cmd *cobra.Command, args []string) error {
	number := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	issue, err := client.ReopenIssue(issueOwner, issueRepo, number)
	if err != nil {
		return fmt.Errorf("failed to reopen issue: %w", err)
	}

	fmt.Printf("Issue #%d reopened successfully\n", issue.Number)
	return nil
}

func issueAddComment(cmd *cobra.Command, args []string) error {
	number := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	comment, err := client.CreateIssueComment(issueOwner, issueRepo, number, issueBody)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	fmt.Printf("Comment added successfully (ID: %d)\n", comment.ID)
	return nil
}
