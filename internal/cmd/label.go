package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	labelRepo          string
	labelOwner         string
	labelNames         []string
	labelName          string
	labelColor         string
	labelNewName       string
	labelNewColor      string
	enterpriseLabelEnt string
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage labels",
	Long:  `Commands for managing Gitee labels`,
}

var issueLabelCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issue labels",
	Long:  `Commands for managing issue labels`,
}

var projectLabelCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project labels",
	Long:  `Commands for managing project labels`,
}

var enterpriseLabelCmd = &cobra.Command{
	Use:   "enterprise",
	Short: "Manage enterprise labels",
	Long:  `Commands for managing enterprise labels`,
}

var labelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repository labels",
	RunE:  labelList,
}

var labelViewCmd = &cobra.Command{
	Use:   "view <name>",
	Short: "View a label",
	Args:  cobra.ExactArgs(1),
	RunE:  labelView,
}

var labelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a label",
	RunE:  labelCreate,
}

var labelUpdateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Update a label",
	Args:  cobra.ExactArgs(1),
	RunE:  labelUpdate,
}

var labelDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a label",
	Args:  cobra.ExactArgs(1),
	RunE:  labelDelete,
}

var issueLabelListCmd = &cobra.Command{
	Use:   "list <number>",
	Short: "List labels for an issue",
	Args:  cobra.ExactArgs(1),
	RunE:  issueLabelList,
}

var issueLabelAddCmd = &cobra.Command{
	Use:   "add <number>",
	Short: "Add labels to an issue",
	Args:  cobra.ExactArgs(1),
	RunE:  issueLabelAdd,
}

var issueLabelReplaceCmd = &cobra.Command{
	Use:   "replace <number>",
	Short: "Replace all labels for an issue",
	Args:  cobra.ExactArgs(1),
	RunE:  issueLabelReplace,
}

var issueLabelRemoveAllCmd = &cobra.Command{
	Use:   "remove-all <number>",
	Short: "Remove all labels from an issue",
	Args:  cobra.ExactArgs(1),
	RunE:  issueLabelRemoveAll,
}

var issueLabelRemoveCmd = &cobra.Command{
	Use:   "remove <number> <name>",
	Short: "Remove a label from an issue",
	Args:  cobra.ExactArgs(2),
	RunE:  issueLabelRemove,
}

var projectLabelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List project labels",
	RunE:  projectLabelList,
}

var projectLabelAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add labels to project",
	RunE:  projectLabelAdd,
}

var projectLabelReplaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Replace all project labels",
	RunE:  projectLabelReplace,
}

var projectLabelRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove project labels",
	RunE:  projectLabelRemove,
}

var enterpriseLabelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List enterprise labels",
	RunE:  enterpriseLabelList,
}

var enterpriseLabelViewCmd = &cobra.Command{
	Use:   "view <name>",
	Short: "View an enterprise label",
	Args:  cobra.ExactArgs(1),
	RunE:  enterpriseLabelView,
}

func init() {
	labelCmd.AddCommand(issueLabelCmd, projectLabelCmd, enterpriseLabelCmd)
	labelCmd.AddCommand(labelListCmd, labelViewCmd, labelCreateCmd, labelUpdateCmd, labelDeleteCmd)

	issueLabelCmd.AddCommand(issueLabelListCmd, issueLabelAddCmd, issueLabelReplaceCmd, issueLabelRemoveAllCmd, issueLabelRemoveCmd)
	projectLabelCmd.AddCommand(projectLabelListCmd, projectLabelAddCmd, projectLabelReplaceCmd, projectLabelRemoveCmd)
	enterpriseLabelCmd.AddCommand(enterpriseLabelListCmd, enterpriseLabelViewCmd)

	rootCmd.AddCommand(labelCmd)

	labelListCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	labelListCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	_ = labelListCmd.MarkFlagRequired("repo")
	_ = labelListCmd.MarkFlagRequired("owner")

	labelViewCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	labelViewCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	_ = labelViewCmd.MarkFlagRequired("repo")
	_ = labelViewCmd.MarkFlagRequired("owner")

	labelCreateCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	labelCreateCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	labelCreateCmd.Flags().StringVarP(&labelName, "name", "n", "", "Label name (required)")
	labelCreateCmd.Flags().StringVarP(&labelColor, "color", "c", "", "Label color (6-digit hex, required)")
	_ = labelCreateCmd.MarkFlagRequired("repo")
	_ = labelCreateCmd.MarkFlagRequired("owner")
	_ = labelCreateCmd.MarkFlagRequired("name")
	_ = labelCreateCmd.MarkFlagRequired("color")

	labelUpdateCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	labelUpdateCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	labelUpdateCmd.Flags().StringVar(&labelNewName, "new-name", "", "New label name")
	labelUpdateCmd.Flags().StringVar(&labelNewColor, "new-color", "", "New label color")
	_ = labelUpdateCmd.MarkFlagRequired("repo")
	_ = labelUpdateCmd.MarkFlagRequired("owner")

	labelDeleteCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	labelDeleteCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	_ = labelDeleteCmd.MarkFlagRequired("repo")
	_ = labelDeleteCmd.MarkFlagRequired("owner")

	issueLabelListCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	issueLabelListCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	_ = issueLabelListCmd.MarkFlagRequired("owner")
	_ = issueLabelListCmd.MarkFlagRequired("repo")

	issueLabelAddCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	issueLabelAddCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	issueLabelAddCmd.Flags().StringSliceVarP(&labelNames, "names", "n", []string{}, "Label names (required)")
	_ = issueLabelAddCmd.MarkFlagRequired("owner")
	_ = issueLabelAddCmd.MarkFlagRequired("repo")
	_ = issueLabelAddCmd.MarkFlagRequired("names")

	issueLabelReplaceCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	issueLabelReplaceCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	issueLabelReplaceCmd.Flags().StringSliceVarP(&labelNames, "names", "n", []string{}, "Label names (required)")
	_ = issueLabelReplaceCmd.MarkFlagRequired("owner")
	_ = issueLabelReplaceCmd.MarkFlagRequired("repo")
	_ = issueLabelReplaceCmd.MarkFlagRequired("names")

	issueLabelRemoveAllCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	issueLabelRemoveAllCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	_ = issueLabelRemoveAllCmd.MarkFlagRequired("owner")
	_ = issueLabelRemoveAllCmd.MarkFlagRequired("repo")

	issueLabelRemoveCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	issueLabelRemoveCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	_ = issueLabelRemoveCmd.MarkFlagRequired("owner")
	_ = issueLabelRemoveCmd.MarkFlagRequired("repo")

	projectLabelListCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	projectLabelListCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	_ = projectLabelListCmd.MarkFlagRequired("owner")
	_ = projectLabelListCmd.MarkFlagRequired("repo")

	projectLabelAddCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	projectLabelAddCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	projectLabelAddCmd.Flags().StringSliceVarP(&labelNames, "names", "n", []string{}, "Label names (required)")
	_ = projectLabelAddCmd.MarkFlagRequired("owner")
	_ = projectLabelAddCmd.MarkFlagRequired("repo")
	_ = projectLabelAddCmd.MarkFlagRequired("names")

	projectLabelReplaceCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	projectLabelReplaceCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	projectLabelReplaceCmd.Flags().StringSliceVarP(&labelNames, "names", "n", []string{}, "Label names (required)")
	_ = projectLabelReplaceCmd.MarkFlagRequired("owner")
	_ = projectLabelReplaceCmd.MarkFlagRequired("repo")
	_ = projectLabelReplaceCmd.MarkFlagRequired("names")

	projectLabelRemoveCmd.Flags().StringVarP(&labelOwner, "owner", "o", "", "Owner name (required)")
	projectLabelRemoveCmd.Flags().StringVarP(&labelRepo, "repo", "r", "", "Repository name (required)")
	projectLabelRemoveCmd.Flags().StringSliceVarP(&labelNames, "names", "n", []string{}, "Label names (required)")
	_ = projectLabelRemoveCmd.MarkFlagRequired("owner")
	_ = projectLabelRemoveCmd.MarkFlagRequired("repo")
	_ = projectLabelRemoveCmd.MarkFlagRequired("names")

	enterpriseLabelListCmd.Flags().StringVarP(&enterpriseLabelEnt, "enterprise", "e", "", "Enterprise name (required)")
	_ = enterpriseLabelListCmd.MarkFlagRequired("enterprise")

	enterpriseLabelViewCmd.Flags().StringVarP(&enterpriseLabelEnt, "enterprise", "e", "", "Enterprise name (required)")
	_ = enterpriseLabelViewCmd.MarkFlagRequired("enterprise")
}

func labelList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	labels, err := client.ListLabels(labelOwner, labelRepo)
	if err != nil {
		return fmt.Errorf("failed to list labels: %w", err)
	}

	if len(labels) == 0 {
		fmt.Println("No labels found")
		return nil
	}

	for _, l := range labels {
		fmt.Printf("%s [%s]\n", l.Name, l.Color)
	}
	return nil
}

func labelView(cmd *cobra.Command, args []string) error {
	name := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	label, err := client.GetLabel(labelOwner, labelRepo, name)
	if err != nil {
		return fmt.Errorf("failed to get label: %w", err)
	}

	fmt.Printf("Name: %s\n", label.Name)
	fmt.Printf("Color: %s\n", label.Color)
	fmt.Printf("ID: %d\n", label.ID)
	if label.URL != "" {
		fmt.Printf("URL: %s\n", label.URL)
	}
	return nil
}

func labelCreate(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	label, err := client.CreateLabel(labelOwner, labelRepo, api.CreateLabelOptions{
		Name:  labelName,
		Color: labelColor,
	})
	if err != nil {
		return fmt.Errorf("failed to create label: %w", err)
	}

	fmt.Printf("Label created: %s [%s]\n", label.Name, label.Color)
	return nil
}

func labelUpdate(cmd *cobra.Command, args []string) error {
	originalName := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	label, err := client.UpdateLabel(labelOwner, labelRepo, originalName, api.UpdateLabelOptions{
		Name:  labelNewName,
		Color: labelNewColor,
	})
	if err != nil {
		return fmt.Errorf("failed to update label: %w", err)
	}

	fmt.Printf("Label updated: %s [%s]\n", label.Name, label.Color)
	return nil
}

func labelDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.DeleteLabel(labelOwner, labelRepo, name)
	if err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}

	fmt.Printf("Label %s deleted\n", name)
	return nil
}

func issueLabelList(cmd *cobra.Command, args []string) error {
	number := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	labels, err := client.ListIssueLabels(labelOwner, labelRepo, number)
	if err != nil {
		return fmt.Errorf("failed to list issue labels: %w", err)
	}

	if len(labels) == 0 {
		fmt.Println("No labels found")
		return nil
	}

	for _, l := range labels {
		fmt.Printf("%s [%s]\n", l.Name, l.Color)
	}
	return nil
}

func issueLabelAdd(cmd *cobra.Command, args []string) error {
	number := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	labels, err := client.AddIssueLabels(labelOwner, labelRepo, number, labelNames)
	if err != nil {
		return fmt.Errorf("failed to add issue labels: %w", err)
	}

	fmt.Println("Labels added:")
	for _, l := range labels {
		fmt.Printf("  %s\n", l.Name)
	}
	return nil
}

func issueLabelReplace(cmd *cobra.Command, args []string) error {
	number := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	labels, err := client.ReplaceIssueLabels(labelOwner, labelRepo, number, labelNames)
	if err != nil {
		return fmt.Errorf("failed to replace issue labels: %w", err)
	}

	fmt.Println("Labels replaced:")
	for _, l := range labels {
		fmt.Printf("  %s\n", l.Name)
	}
	return nil
}

func issueLabelRemoveAll(cmd *cobra.Command, args []string) error {
	number := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.DeleteAllIssueLabels(labelOwner, labelRepo, number)
	if err != nil {
		return fmt.Errorf("failed to delete all issue labels: %w", err)
	}

	fmt.Printf("All labels removed from issue #%s\n", number)
	return nil
}

func issueLabelRemove(cmd *cobra.Command, args []string) error {
	number := args[0]
	name := args[1]
	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.DeleteIssueLabel(labelOwner, labelRepo, number, name)
	if err != nil {
		return fmt.Errorf("failed to delete issue label: %w", err)
	}

	fmt.Printf("Label '%s' removed from issue #%s\n", name, number)
	return nil
}

func projectLabelList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	labels, err := client.ListProjectLabels(labelOwner, labelRepo)
	if err != nil {
		return fmt.Errorf("failed to list project labels: %w", err)
	}

	if len(labels) == 0 {
		fmt.Println("No project labels found")
		return nil
	}

	for _, l := range labels {
		fmt.Printf("%s (id: %d)\n", l.Name, l.ID)
	}
	return nil
}

func projectLabelAdd(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	labels, err := client.AddProjectLabels(labelOwner, labelRepo, labelNames)
	if err != nil {
		return fmt.Errorf("failed to add project labels: %w", err)
	}

	fmt.Println("Labels added:")
	for _, l := range labels {
		fmt.Printf("  %s (id: %d)\n", l.Name, l.ID)
	}
	return nil
}

func projectLabelReplace(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	labels, err := client.ReplaceProjectLabels(labelOwner, labelRepo, labelNames)
	if err != nil {
		return fmt.Errorf("failed to replace project labels: %w", err)
	}

	fmt.Println("Labels replaced:")
	for _, l := range labels {
		fmt.Printf("  %s (id: %d)\n", l.Name, l.ID)
	}
	return nil
}

func projectLabelRemove(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	err = client.DeleteProjectLabels(labelOwner, labelRepo, labelNames)
	if err != nil {
		return fmt.Errorf("failed to delete project labels: %w", err)
	}

	fmt.Println("Labels removed")
	return nil
}

func enterpriseLabelList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	labels, err := client.ListEnterpriseLabels(enterpriseLabelEnt)
	if err != nil {
		return fmt.Errorf("failed to list enterprise labels: %w", err)
	}

	if len(labels) == 0 {
		fmt.Println("No enterprise labels found")
		return nil
	}

	for _, l := range labels {
		fmt.Printf("%s [%s]\n", l.Name, l.Color)
	}
	return nil
}

func enterpriseLabelView(cmd *cobra.Command, args []string) error {
	name := args[0]
	client, err := getClient()
	if err != nil {
		return err
	}

	label, err := client.GetEnterpriseLabel(enterpriseLabelEnt, name)
	if err != nil {
		return fmt.Errorf("failed to get enterprise label: %w", err)
	}

	fmt.Printf("Name: %s\n", label.Name)
	fmt.Printf("Color: %s\n", label.Color)
	fmt.Printf("ID: %d\n", label.ID)
	if label.URL != "" {
		fmt.Printf("URL: %s\n", label.URL)
	}
	return nil
}
