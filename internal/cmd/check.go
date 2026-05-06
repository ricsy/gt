package cmd

import (
	"fmt"
	"strconv"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	checkRepo          string
	checkOwner         string
	checkName          string
	checkHeadSHA       string
	checkDetailsURL    string
	checkStatus        string
	checkConclusion    string
	checkStartedAt     string
	checkCompletedAt   string
	checkOutputTitle   string
	checkOutputSummary string
	checkOutputText    string
	checkPage          int
	checkPerPage       int
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Manage check runs",
	Long:  `Commands for managing Gitee check runs`,
}

var checkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a check run",
	RunE:  checkCreate,
}

var checkViewCmd = &cobra.Command{
	Use:   "view <check_run_id>",
	Short: "View a check run",
	Args:  cobra.ExactArgs(1),
	RunE:  checkView,
}

var checkUpdateCmd = &cobra.Command{
	Use:   "update <check_run_id>",
	Short: "Update a check run",
	Args:  cobra.ExactArgs(1),
	RunE:  checkUpdate,
}

var checkAnnotationsCmd = &cobra.Command{
	Use:   "annotations <check_run_id>",
	Short: "Get check run annotations",
	Args:  cobra.ExactArgs(1),
	RunE:  checkAnnotations,
}

var checkCommitsCmd = &cobra.Command{
	Use:   "commits <ref>",
	Short: "List check runs for a commit",
	Args:  cobra.ExactArgs(1),
	RunE:  checkCommits,
}

func init() {
	checkCmd.AddCommand(checkCreateCmd, checkViewCmd, checkUpdateCmd, checkAnnotationsCmd, checkCommitsCmd)

	checkCreateCmd.Flags().StringVarP(&checkRepo, "repo", "r", "", "Repository name (required)")
	checkCreateCmd.Flags().StringVarP(&checkOwner, "owner", "o", "", "Owner name (required)")
	checkCreateCmd.Flags().StringVarP(&checkName, "name", "n", "", "Check run name (required)")
	checkCreateCmd.Flags().StringVarP(&checkHeadSHA, "head_sha", "s", "", "Commit SHA (required)")
	checkCreateCmd.Flags().StringVar(&checkDetailsURL, "details_url", "", "Details URL")
	checkCreateCmd.Flags().StringVar(&checkStatus, "status", "queued", "Status (queued, in_progress, completed)")
	checkCreateCmd.Flags().StringVarP(&checkOutputTitle, "output_title", "t", "", "Output title (required)")
	checkCreateCmd.Flags().StringVarP(&checkOutputSummary, "output_summary", "m", "", "Output summary (required)")
	checkCreateCmd.Flags().StringVar(&checkOutputText, "output_text", "", "Output text")
	_ = checkCreateCmd.MarkFlagRequired("repo")
	_ = checkCreateCmd.MarkFlagRequired("owner")
	_ = checkCreateCmd.MarkFlagRequired("name")
	_ = checkCreateCmd.MarkFlagRequired("head_sha")
	_ = checkCreateCmd.MarkFlagRequired("output_title")
	_ = checkCreateCmd.MarkFlagRequired("output_summary")

	checkViewCmd.Flags().StringVarP(&checkRepo, "repo", "r", "", "Repository name (required)")
	checkViewCmd.Flags().StringVarP(&checkOwner, "owner", "o", "", "Owner name (required)")
	_ = checkViewCmd.MarkFlagRequired("repo")
	_ = checkViewCmd.MarkFlagRequired("owner")

	checkUpdateCmd.Flags().StringVarP(&checkRepo, "repo", "r", "", "Repository name (required)")
	checkUpdateCmd.Flags().StringVarP(&checkOwner, "owner", "o", "", "Owner name (required)")
	checkUpdateCmd.Flags().StringVar(&checkDetailsURL, "details_url", "", "Details URL")
	checkUpdateCmd.Flags().StringVar(&checkStatus, "status", "", "Status (queued, in_progress, completed)")
	checkUpdateCmd.Flags().StringVar(&checkStartedAt, "started_at", "", "Started at (RFC3339 format)")
	checkUpdateCmd.Flags().StringVar(&checkConclusion, "conclusion", "", "Conclusion (neutral, success, failure, cancelled, action_required, timed_out, skipped)")
	checkUpdateCmd.Flags().StringVar(&checkCompletedAt, "completed_at", "", "Completed at (RFC3339 format)")
	checkUpdateCmd.Flags().StringVar(&checkOutputTitle, "output_title", "", "Output title")
	checkUpdateCmd.Flags().StringVar(&checkOutputSummary, "output_summary", "", "Output summary")
	checkUpdateCmd.Flags().StringVar(&checkOutputText, "output_text", "", "Output text")
	_ = checkUpdateCmd.MarkFlagRequired("repo")
	_ = checkUpdateCmd.MarkFlagRequired("owner")

	checkAnnotationsCmd.Flags().StringVarP(&checkRepo, "repo", "r", "", "Repository name (required)")
	checkAnnotationsCmd.Flags().StringVarP(&checkOwner, "owner", "o", "", "Owner name (required)")
	checkAnnotationsCmd.Flags().IntVar(&checkPage, "page", 0, "Page number")
	checkAnnotationsCmd.Flags().IntVar(&checkPerPage, "per-page", 0, "Items per page (max 100)")
	_ = checkAnnotationsCmd.MarkFlagRequired("repo")
	_ = checkAnnotationsCmd.MarkFlagRequired("owner")

	checkCommitsCmd.Flags().StringVarP(&checkRepo, "repo", "r", "", "Repository name (required)")
	checkCommitsCmd.Flags().StringVarP(&checkOwner, "owner", "o", "", "Owner name (required)")
	checkCommitsCmd.Flags().IntVar(&checkPage, "page", 0, "Page number")
	checkCommitsCmd.Flags().IntVar(&checkPerPage, "per-page", 0, "Items per page (max 100)")
	_ = checkCommitsCmd.MarkFlagRequired("repo")
	_ = checkCommitsCmd.MarkFlagRequired("owner")

	rootCmd.AddCommand(checkCmd)
}

func checkCreate(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	checkRun, err := client.CreateCheckRun(checkOwner, checkRepo, api.CreateCheckRunOptions{
		Name:          checkName,
		HeadSHA:       checkHeadSHA,
		DetailsURL:    checkDetailsURL,
		Status:        checkStatus,
		OutputTitle:   checkOutputTitle,
		OutputSummary: checkOutputSummary,
		OutputText:    checkOutputText,
	})
	if err != nil {
		return fmt.Errorf("failed to create check run: %w", err)
	}

	fmt.Printf("Check run created successfully: #%d %s\n", checkRun.ID, checkRun.Name)
	return nil
}

func checkView(cmd *cobra.Command, args []string) error {
	checkRunID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid check run ID: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	checkRun, err := client.GetCheckRun(checkOwner, checkRepo, checkRunID)
	if err != nil {
		return fmt.Errorf("failed to get check run: %w", err)
	}

	fmt.Printf("#%d %s\n", checkRun.ID, checkRun.Name)
	fmt.Printf("Status: %s | Conclusion: %s\n", checkRun.Status, checkRun.Conclusion)
	fmt.Printf("Head SHA: %s\n", checkRun.HeadSHA)
	if checkRun.Output != nil {
		fmt.Printf("\n--- Output ---\n%s\n", checkRun.Output.Summary)
	}
	return nil
}

func checkUpdate(cmd *cobra.Command, args []string) error {
	checkRunID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid check run ID: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	checkRun, err := client.UpdateCheckRun(checkOwner, checkRepo, checkRunID, api.UpdateCheckRunOptions{
		DetailsURL:    checkDetailsURL,
		Status:        checkStatus,
		StartedAt:     checkStartedAt,
		Conclusion:    checkConclusion,
		CompletedAt:   checkCompletedAt,
		OutputTitle:   checkOutputTitle,
		OutputSummary: checkOutputSummary,
		OutputText:    checkOutputText,
	})
	if err != nil {
		return fmt.Errorf("failed to update check run: %w", err)
	}

	fmt.Printf("Check run #%d updated successfully\n", checkRun.ID)
	return nil
}

func checkAnnotations(cmd *cobra.Command, args []string) error {
	checkRunID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid check run ID: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	annotations, err := client.GetCheckRunAnnotations(checkOwner, checkRepo, checkRunID, api.ListCheckRunsOptions{
		Page:    checkPage,
		PerPage: checkPerPage,
	})
	if err != nil {
		return fmt.Errorf("failed to get check run annotations: %w", err)
	}

	if len(annotations) == 0 {
		fmt.Println("No annotations found")
		return nil
	}

	for _, a := range annotations {
		fmt.Printf("[%s] %s:%d-%d %s\n", a.AnnotationLevel, a.Path, a.StartLine, a.EndLine, a.Message)
	}
	return nil
}

func checkCommits(cmd *cobra.Command, args []string) error {
	ref := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	checkRuns, err := client.ListCommitCheckRuns(checkOwner, checkRepo, ref, api.ListCheckRunsOptions{
		Page:    checkPage,
		PerPage: checkPerPage,
	})
	if err != nil {
		return fmt.Errorf("failed to list commit check runs: %w", err)
	}

	if len(checkRuns) == 0 {
		fmt.Println("No check runs found")
		return nil
	}

	for _, cr := range checkRuns {
		fmt.Printf("#%d [%s] %s (%s)\n", cr.ID, cr.Status, cr.Name, cr.HeadSHA)
	}
	return nil
}
