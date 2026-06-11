package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	gistDescription string
	gistPublic      bool
	gistFiles       string
	gistCommentBody string
	gistSince       string
)

var gistCmd = &cobra.Command{
	Use:   "gist",
	Short: "Manage gists",
	Long:  `Commands for managing Gitee gists`,
}

var gistListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your gists",
	RunE:  gistList,
}

var gistStarredCmd = &cobra.Command{
	Use:   "starred",
	Short: "List your starred gists",
	RunE:  gistStarred,
}

var gistViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View a gist",
	Args:  cobra.ExactArgs(1),
	RunE:  gistView,
}

var gistCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a gist",
	RunE:  gistCreate,
}

var gistUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a gist",
	Args:  cobra.ExactArgs(1),
	RunE:  gistUpdate,
}

var gistDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a gist",
	Args:  cobra.ExactArgs(1),
	RunE:  gistDelete,
}

var gistStarCmd = &cobra.Command{
	Use:   "star <id>",
	Short: "Star a gist",
	Args:  cobra.ExactArgs(1),
	RunE:  gistStar,
}

var gistUnstarCmd = &cobra.Command{
	Use:   "unstar <id>",
	Short: "Unstar a gist",
	Args:  cobra.ExactArgs(1),
	RunE:  gistUnstar,
}

var gistForkCmd = &cobra.Command{
	Use:   "fork <id>",
	Short: "Fork a gist",
	Args:  cobra.ExactArgs(1),
	RunE:  gistFork,
}

var gistForksCmd = &cobra.Command{
	Use:   "forks <id>",
	Short: "List forks of a gist",
	Args:  cobra.ExactArgs(1),
	RunE:  gistForks,
}

var gistCommitsCmd = &cobra.Command{
	Use:   "commits <id>",
	Short: "List commits of a gist",
	Args:  cobra.ExactArgs(1),
	RunE:  gistCommits,
}

var gistCommentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage gist comments",
}

var gistCommentListCmd = &cobra.Command{
	Use:   "list <gist_id>",
	Short: "List comments on a gist",
	Args:  cobra.ExactArgs(1),
	RunE:  gistCommentList,
}

var gistCommentViewCmd = &cobra.Command{
	Use:   "view <gist_id> <comment_id>",
	Short: "View a comment",
	Args:  cobra.ExactArgs(2),
	RunE:  gistCommentView,
}

var gistCommentCreateCmd = &cobra.Command{
	Use:   "create <gist_id>",
	Short: "Create a comment on a gist",
	Args:  cobra.ExactArgs(1),
	RunE:  gistCommentCreate,
}

var gistCommentUpdateCmd = &cobra.Command{
	Use:   "update <gist_id> <comment_id>",
	Short: "Update a comment",
	Args:  cobra.ExactArgs(2),
	RunE:  gistCommentUpdate,
}

var gistCommentDeleteCmd = &cobra.Command{
	Use:   "delete <gist_id> <comment_id>",
	Short: "Delete a comment",
	Args:  cobra.ExactArgs(2),
	RunE:  gistCommentDelete,
}

func init() {
	rootCmd.AddCommand(gistCmd)
	gistCmd.AddCommand(gistListCmd, gistStarredCmd, gistViewCmd, gistCreateCmd, gistUpdateCmd, gistDeleteCmd)
	gistCmd.AddCommand(gistStarCmd, gistUnstarCmd, gistForkCmd, gistForksCmd, gistCommitsCmd)
	gistCmd.AddCommand(gistCommentCmd)
	gistCommentCmd.AddCommand(gistCommentListCmd, gistCommentViewCmd, gistCommentCreateCmd, gistCommentUpdateCmd, gistCommentDeleteCmd)

	gistListCmd.Flags().StringVar(&gistSince, "since", "", "Start time for updates (ISO 8601)")
	gistStarredCmd.Flags().StringVar(&gistSince, "since", "", "Start time for updates (ISO 8601)")

	gistCreateCmd.Flags().StringVarP(&gistDescription, "description", "d", "", "Gist description (required)")
	gistCreateCmd.Flags().BoolVarP(&gistPublic, "public", "p", false, "Make gist public")
	gistCreateCmd.Flags().StringVarP(&gistFiles, "files", "f", "", "Files content (JSON format, required)")
	_ = gistCreateCmd.MarkFlagRequired("description")
	_ = gistCreateCmd.MarkFlagRequired("files")

	gistUpdateCmd.Flags().StringVarP(&gistDescription, "description", "d", "", "New description")
	gistUpdateCmd.Flags().StringVarP(&gistFiles, "files", "f", "", "New files content (JSON format)")

	gistCommentCreateCmd.Flags().StringVarP(&gistCommentBody, "body", "b", "", "Comment body (required)")
	_ = gistCommentCreateCmd.MarkFlagRequired("body")

	gistCommentUpdateCmd.Flags().StringVarP(&gistCommentBody, "body", "b", "", "New comment body (required)")
	_ = gistCommentUpdateCmd.MarkFlagRequired("body")
}

func gistList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	opts := &api.ListGistsOptions{}
	if gistSince != "" {
		opts.Since = gistSince
	}

	gists, err := client.ListGists(opts)
	if err != nil {
		return fmt.Errorf("failed to list gists: %w", err)
	}

	if len(gists) == 0 {
		fmt.Println("No gists found")
		return nil
	}

	for _, g := range gists {
		fmt.Printf("%s [%s]\n", g.ID, g.Description)
	}
	return nil
}

func gistStarred(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	opts := &api.ListGistsOptions{}
	if gistSince != "" {
		opts.Since = gistSince
	}

	gists, err := client.ListStarredGists(opts)
	if err != nil {
		return fmt.Errorf("failed to list starred gists: %w", err)
	}

	if len(gists) == 0 {
		fmt.Println("No starred gists found")
		return nil
	}

	for _, g := range gists {
		fmt.Printf("%s [%s]\n", g.ID, g.Description)
	}
	return nil
}

func gistView(cmd *cobra.Command, args []string) error {
	id := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	gist, err := client.GetGist(id)
	if err != nil {
		return fmt.Errorf("failed to get gist: %w", err)
	}

	fmt.Printf("ID: %s\n", gist.ID)
	fmt.Printf("Description: %s\n", gist.Description)
	fmt.Printf("Public: %t\n", gist.Public)
	fmt.Printf("Created: %s\n", gist.CreatedAt)
	fmt.Printf("Updated: %s\n", gist.UpdatedAt)
	fmt.Printf("URL: %s\n", gist.URL)
	return nil
}

func gistCreate(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	files, err := parseGistFiles(gistFiles)
	if err != nil {
		return err
	}

	gist, err := client.CreateGist(api.CreateGistOptions{
		Files:       files,
		Description: gistDescription,
		Public:      gistPublic,
	})
	if err != nil {
		return fmt.Errorf("failed to create gist: %w", err)
	}

	fmt.Printf("Gist created: %s\n", gist.ID)
	return nil
}

func gistUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	opts := api.UpdateGistOptions{}
	if cmd.Flags().Changed("description") {
		opts.Description = gistDescription
	}
	if cmd.Flags().Changed("files") {
		files, err := parseGistFiles(gistFiles)
		if err != nil {
			return err
		}
		opts.Files = files
	}

	gist, err := client.UpdateGist(id, opts)
	if err != nil {
		return fmt.Errorf("failed to update gist: %w", err)
	}

	fmt.Printf("Gist updated: %s\n", gist.ID)
	return nil
}

func gistDelete(cmd *cobra.Command, args []string) error {
	id := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.DeleteGist(id)
	if err != nil {
		return fmt.Errorf("failed to delete gist: %w", err)
	}

	fmt.Printf("Gist %s deleted\n", id)
	return nil
}

func gistStar(cmd *cobra.Command, args []string) error {
	id := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.StarGist(id)
	if err != nil {
		return fmt.Errorf("failed to star gist: %w", err)
	}

	fmt.Printf("Gist %s starred\n", id)
	return nil
}

func gistUnstar(cmd *cobra.Command, args []string) error {
	id := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.UnstarGist(id)
	if err != nil {
		return fmt.Errorf("failed to unstar gist: %w", err)
	}

	fmt.Printf("Gist %s unstarred\n", id)
	return nil
}

func gistFork(cmd *cobra.Command, args []string) error {
	id := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	gist, err := client.ForkGist(id)
	if err != nil {
		return fmt.Errorf("failed to fork gist: %w", err)
	}

	fmt.Printf("Gist forked: %s\n", gist.ID)
	return nil
}

func gistForks(cmd *cobra.Command, args []string) error {
	id := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	gists, err := client.ListGistForks(id)
	if err != nil {
		return fmt.Errorf("failed to list gist forks: %w", err)
	}

	if len(gists) == 0 {
		fmt.Println("No forks found")
		return nil
	}

	for _, g := range gists {
		fmt.Printf("%s [%s]\n", g.ID, g.Description)
	}
	return nil
}

func gistCommits(cmd *cobra.Command, args []string) error {
	id := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	gist, err := client.ListGistCommits(id)
	if err != nil {
		return fmt.Errorf("failed to list gist commits: %w", err)
	}

	if gist == nil || gist.History == nil {
		fmt.Println("No commits found")
		return nil
	}

	fmt.Printf("%v\n", gist.History)
	return nil
}

func gistCommentList(cmd *cobra.Command, args []string) error {
	gistID := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	comments, err := client.ListGistComments(gistID)
	if err != nil {
		return fmt.Errorf("failed to list gist comments: %w", err)
	}

	if len(comments) == 0 {
		fmt.Println("No comments found")
		return nil
	}

	for _, c := range comments {
		fmt.Printf("%d: %s\n", c.ID, c.Body)
	}
	return nil
}

func gistCommentView(cmd *cobra.Command, args []string) error {
	gistID := args[0]
	commentIDStr := args[1]
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	comment, err := client.GetGistComment(gistID, commentID)
	if err != nil {
		return fmt.Errorf("failed to get gist comment: %w", err)
	}

	fmt.Printf("ID: %d\n", comment.ID)
	fmt.Printf("Body: %s\n", comment.Body)
	fmt.Printf("Created: %s\n", comment.CreatedAt)
	fmt.Printf("Updated: %s\n", comment.UpdatedAt)
	return nil
}

func gistCommentCreate(cmd *cobra.Command, args []string) error {
	gistID := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	comment, err := client.CreateGistComment(gistID, api.CreateGistCommentOptions{
		Body: gistCommentBody,
	})
	if err != nil {
		return fmt.Errorf("failed to create gist comment: %w", err)
	}

	fmt.Printf("Comment created: %d\n", comment.ID)
	return nil
}

func gistCommentUpdate(cmd *cobra.Command, args []string) error {
	gistID := args[0]
	commentIDStr := args[1]
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	comment, err := client.UpdateGistComment(gistID, commentID, api.UpdateGistCommentOptions{
		Body: gistCommentBody,
	})
	if err != nil {
		return fmt.Errorf("failed to update gist comment: %w", err)
	}

	fmt.Printf("Comment updated: %d\n", comment.ID)
	return nil
}

func gistCommentDelete(cmd *cobra.Command, args []string) error {
	gistID := args[0]
	commentIDStr := args[1]
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid comment ID: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.DeleteGistComment(gistID, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete gist comment: %w", err)
	}

	fmt.Printf("Comment %d deleted\n", commentID)
	return nil
}

// parseGistFiles 将 CLI 传入的 JSON 文本解析为 API 需要的 files 对象。
func parseGistFiles(raw string) (map[string]map[string]string, error) {
	var files map[string]map[string]string
	if err := json.Unmarshal([]byte(raw), &files); err != nil {
		return nil, fmt.Errorf("invalid --files JSON: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("invalid --files JSON: at least one file is required")
	}
	return files, nil
}
