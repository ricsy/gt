package cmd

import (
	"fmt"
	"strconv"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	prRepo  string
	prState string
	prTitle string
	prBody  string
	prHead  string
	prBase  string
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
	Long:  `Commands for working with pull requests on Gitee.`,
}

var prListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests",
	RunE:  prList,
}

var prViewCmd = &cobra.Command{
	Use:   "view <number>",
	Short: "View a pull request",
	Args:  cobra.ExactArgs(1),
	RunE:  prView,
}

var prCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pull request",
	RunE:  prCreate,
}

var prMergeCmd = &cobra.Command{
	Use:   "merge <number>",
	Short: "Merge a pull request",
	Args:  cobra.ExactArgs(1),
	RunE:  prMerge,
}

var prCloseCmd = &cobra.Command{
	Use:   "close <number>",
	Short: "Close a pull request",
	Args:  cobra.ExactArgs(1),
	RunE:  prClose,
}

var prCommentCmd = &cobra.Command{
	Use:   "comment <number>",
	Short: "Comment on a pull request",
	Args:  cobra.ExactArgs(1),
	RunE:  prComment,
}

func init() {
	rootCmd.AddCommand(prCmd)

	prCmd.AddCommand(prListCmd, prViewCmd, prCreateCmd, prMergeCmd, prCloseCmd, prCommentCmd)

	prListCmd.Flags().StringVar(&prRepo, "repo", "", "owner/repo")
	prListCmd.Flags().StringVar(&prState, "state", api.StateOpen, "Filter by state: open, closed, all")

	prViewCmd.Flags().StringVar(&prRepo, "repo", "", "owner/repo")

	prCreateCmd.Flags().StringVar(&prRepo, "repo", "", "owner/repo")
	prCreateCmd.Flags().StringVar(&prTitle, "title", "", "PR title")
	prCreateCmd.Flags().StringVar(&prBody, "body", "", "PR body")
	prCreateCmd.Flags().StringVar(&prHead, "head", "", "Head branch")
	prCreateCmd.Flags().StringVar(&prBase, "base", "main", "Base branch")

	prMergeCmd.Flags().StringVar(&prRepo, "repo", "", "owner/repo")

	prCloseCmd.Flags().StringVar(&prRepo, "repo", "", "owner/repo")

	prCommentCmd.Flags().StringVar(&prRepo, "repo", "", "owner/repo")
	prCommentCmd.Flags().StringVar(&prBody, "body", "", "Comment body")
}

func prList(cmd *cobra.Command, args []string) error {
	owner, repo, err := resolveRepoFlag(prRepo)
	if err != nil {
		return err
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	prs, err := client.ListPRs(owner, repo, prState)
	if err != nil {
		return err
	}

	if len(prs) == 0 {
		fmt.Println("No pull requests found")
		return nil
	}

	for _, pr := range prs {
		fmt.Printf("#%d\t%s\t[%s]\n", pr.Number, pr.Title, pr.State)
	}
	return nil
}

func prView(cmd *cobra.Command, args []string) error {
	owner, repo, err := resolveRepoFlag(prRepo)
	if err != nil {
		return err
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid PR number: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	pr, err := client.GetPR(owner, repo, number)
	if err != nil {
		return err
	}

	fmt.Printf("Title: %s\n", pr.Title)
	fmt.Printf("State: %s\n", pr.State)
	fmt.Printf("Author: %s\n", pr.User.Login)
	fmt.Printf("URL: %s\n", pr.HTMLURL)
	fmt.Printf("\n%s\n", pr.Body)
	return nil
}

func prCreate(cmd *cobra.Command, args []string) error {
	owner, repo, err := resolveRepoFlag(prRepo)
	if err != nil {
		return err
	}

	if prTitle == "" {
		return fmt.Errorf("title is required: use --title")
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	pr, err := client.CreatePR(owner, repo, prTitle, prBody, prHead, prBase)
	if err != nil {
		return err
	}

	fmt.Printf("Created PR #%d: %s\n", pr.Number, pr.HTMLURL)
	return nil
}

func prMerge(cmd *cobra.Command, args []string) error {
	owner, repo, err := resolveRepoFlag(prRepo)
	if err != nil {
		return err
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid PR number: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.MergePR(owner, repo, number)
	if err != nil {
		return err
	}

	fmt.Printf("Merged PR #%d\n", number)
	return nil
}

func prClose(cmd *cobra.Command, args []string) error {
	owner, repo, err := resolveRepoFlag(prRepo)
	if err != nil {
		return err
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid PR number: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	_, err = client.ClosePR(owner, repo, number)
	if err != nil {
		return err
	}

	fmt.Printf("Closed PR #%d\n", number)
	return nil
}

func prComment(cmd *cobra.Command, args []string) error {
	owner, repo, err := resolveRepoFlag(prRepo)
	if err != nil {
		return err
	}

	if prBody == "" {
		return fmt.Errorf("body is required: use --body")
	}

	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid PR number: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.CreatePRComment(owner, repo, number, prBody)
	if err != nil {
		return err
	}

	fmt.Printf("Comment added to PR #%d\n", number)
	return nil
}
