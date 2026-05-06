package cmd

import (
	"fmt"
	"strconv"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	milestoneRepo      string
	milestoneOwner     string
	milestoneState     string
	milestoneSort      string
	milestoneDirection string
	milestoneTitle     string
	milestoneDesc      string
	milestoneDueOn     string
	milestoneLimit     int
)

var milestoneCmd = &cobra.Command{
	Use:   "milestone",
	Short: "Manage milestones",
	Long:  `Commands for managing Gitee milestones`,
}

var milestoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List milestones",
	RunE:  milestoneList,
}

var milestoneViewCmd = &cobra.Command{
	Use:   "view <number>",
	Short: "View a milestone",
	Args:  cobra.ExactArgs(1),
	RunE:  milestoneView,
}

var milestoneCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a milestone",
	RunE:  milestoneCreate,
}

var milestoneUpdateCmd = &cobra.Command{
	Use:   "update <number>",
	Short: "Update a milestone",
	Args:  cobra.ExactArgs(1),
	RunE:  milestoneUpdate,
}

var milestoneDeleteCmd = &cobra.Command{
	Use:   "delete <number>",
	Short: "Delete a milestone",
	Args:  cobra.ExactArgs(1),
	RunE:  milestoneDelete,
}

func init() {
	milestoneCmd.AddCommand(milestoneListCmd, milestoneViewCmd, milestoneCreateCmd, milestoneUpdateCmd, milestoneDeleteCmd)

	milestoneListCmd.Flags().StringVarP(&milestoneRepo, "repo", "r", "", "Repository name (required)")
	milestoneListCmd.Flags().StringVarP(&milestoneOwner, "owner", "o", "", "Owner name (required)")
	milestoneListCmd.Flags().StringVar(&milestoneState, "state", "open", "Filter by state (open, closed, all)")
	milestoneListCmd.Flags().StringVar(&milestoneSort, "sort", "due_on", "Sort by: due_on")
	milestoneListCmd.Flags().StringVar(&milestoneDirection, "direction", "asc", "Sort direction: asc, desc")
	milestoneListCmd.Flags().IntVarP(&milestoneLimit, "limit", "l", 20, "Maximum number of milestones")
	_ = milestoneListCmd.MarkFlagRequired("repo")
	_ = milestoneListCmd.MarkFlagRequired("owner")

	milestoneViewCmd.Flags().StringVarP(&milestoneRepo, "repo", "r", "", "Repository name (required)")
	milestoneViewCmd.Flags().StringVarP(&milestoneOwner, "owner", "o", "", "Owner name (required)")
	_ = milestoneViewCmd.MarkFlagRequired("repo")
	_ = milestoneViewCmd.MarkFlagRequired("owner")

	milestoneCreateCmd.Flags().StringVarP(&milestoneRepo, "repo", "r", "", "Repository name (required)")
	milestoneCreateCmd.Flags().StringVarP(&milestoneOwner, "owner", "o", "", "Owner name (required)")
	milestoneCreateCmd.Flags().StringVarP(&milestoneTitle, "title", "t", "", "Milestone title (required)")
	milestoneCreateCmd.Flags().StringVarP(&milestoneDesc, "description", "d", "", "Milestone description")
	milestoneCreateCmd.Flags().StringVarP(&milestoneDueOn, "due_on", "", "", "Due date (YYYY-MM-DD)")
	_ = milestoneCreateCmd.MarkFlagRequired("repo")
	_ = milestoneCreateCmd.MarkFlagRequired("owner")
	_ = milestoneCreateCmd.MarkFlagRequired("title")

	milestoneUpdateCmd.Flags().StringVarP(&milestoneRepo, "repo", "r", "", "Repository name (required)")
	milestoneUpdateCmd.Flags().StringVarP(&milestoneOwner, "owner", "o", "", "Owner name (required)")
	milestoneUpdateCmd.Flags().StringVarP(&milestoneTitle, "title", "t", "", "Milestone title")
	milestoneUpdateCmd.Flags().StringVarP(&milestoneDesc, "description", "d", "", "Milestone description")
	milestoneUpdateCmd.Flags().StringVar(&milestoneState, "state", "", "Milestone state (open, closed)")
	milestoneUpdateCmd.Flags().StringVarP(&milestoneDueOn, "due_on", "", "", "Due date (YYYY-MM-DD)")
	_ = milestoneUpdateCmd.MarkFlagRequired("repo")
	_ = milestoneUpdateCmd.MarkFlagRequired("owner")

	milestoneDeleteCmd.Flags().StringVarP(&milestoneRepo, "repo", "r", "", "Repository name (required)")
	milestoneDeleteCmd.Flags().StringVarP(&milestoneOwner, "owner", "o", "", "Owner name (required)")
	_ = milestoneDeleteCmd.MarkFlagRequired("repo")
	_ = milestoneDeleteCmd.MarkFlagRequired("owner")

	rootCmd.AddCommand(milestoneCmd)
}

func milestoneList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	milestones, err := client.ListMilestones(milestoneOwner, milestoneRepo, api.ListMilestonesOptions{
		State:     milestoneState,
		Sort:      milestoneSort,
		Direction: milestoneDirection,
		Page:      1,
		PerPage:   milestoneLimit,
	})
	if err != nil {
		return fmt.Errorf("failed to list milestones: %w", err)
	}

	if len(milestones) == 0 {
		fmt.Println("No milestones found")
		return nil
	}

	for _, m := range milestones {
		fmt.Printf("#%d [%s] %s (due: %s)\n", m.Number, m.State, m.Title, m.DueOn)
	}
	return nil
}

func milestoneView(cmd *cobra.Command, args []string) error {
	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid milestone number: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	milestone, err := client.GetMilestone(milestoneOwner, milestoneRepo, number)
	if err != nil {
		return fmt.Errorf("failed to get milestone: %w", err)
	}

	fmt.Printf("#%d %s\n", milestone.Number, milestone.Title)
	fmt.Printf("State: %s | Due: %s\n", milestone.State, milestone.DueOn)
	fmt.Printf("Open issues: %d | Closed issues: %d\n", milestone.OpenIssues, milestone.ClosedIssues)
	if milestone.Description != "" {
		fmt.Printf("\n--- Description ---\n%s\n", milestone.Description)
	}
	return nil
}

func milestoneCreate(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	milestone, err := client.CreateMilestone(milestoneOwner, milestoneRepo, api.CreateMilestoneOptions{
		Title:       milestoneTitle,
		Description: milestoneDesc,
		DueOn:       milestoneDueOn,
	})
	if err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}

	fmt.Printf("Milestone created successfully: #%d %s\n", milestone.Number, milestone.Title)
	return nil
}

func milestoneUpdate(cmd *cobra.Command, args []string) error {
	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid milestone number: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	milestone, err := client.UpdateMilestone(milestoneOwner, milestoneRepo, number, api.UpdateMilestoneOptions{
		Title:       milestoneTitle,
		Description: milestoneDesc,
		State:       milestoneState,
		DueOn:       milestoneDueOn,
	})
	if err != nil {
		return fmt.Errorf("failed to update milestone: %w", err)
	}

	fmt.Printf("Milestone #%d updated successfully\n", milestone.Number)
	return nil
}

func milestoneDelete(cmd *cobra.Command, args []string) error {
	number, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid milestone number: %w", err)
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.DeleteMilestone(milestoneOwner, milestoneRepo, number)
	if err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}

	fmt.Printf("Milestone #%d deleted successfully\n", number)
	return nil
}
