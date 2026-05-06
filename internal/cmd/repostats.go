package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	repostatsRepo  string
	repostatsOwner string
	repostatsStart string
	repostatsEnd   string
)

var repostatsCmd = &cobra.Command{
	Use:   "repostats",
	Short: "Repository statistics",
	Long:  `Commands for repository statistics (traffic, languages, contributors)`,
}

var repostatsTrafficCmd = &cobra.Command{
	Use:   "traffic",
	Short: "Get repository traffic data",
	RunE:  repostatsTraffic,
}

var repostatsLanguagesCmd = &cobra.Command{
	Use:   "languages",
	Short: "Get repository languages",
	RunE:  repostatsLanguages,
}

var repostatsContributorsCmd = &cobra.Command{
	Use:   "contributors",
	Short: "Get repository contributors",
	RunE:  repostatsContributors,
}

func init() {
	repostatsCmd.AddCommand(repostatsTrafficCmd, repostatsLanguagesCmd, repostatsContributorsCmd)

	repostatsTrafficCmd.Flags().StringVarP(&repostatsRepo, "repo", "r", "", "Repository name (required)")
	repostatsTrafficCmd.Flags().StringVarP(&repostatsOwner, "owner", "o", "", "Owner name (required)")
	repostatsTrafficCmd.Flags().StringVar(&repostatsStart, "start-day", "", "Start day (yyyy-MM-dd)")
	repostatsTrafficCmd.Flags().StringVar(&repostatsEnd, "end-day", "", "End day (yyyy-MM-dd)")
	_ = repostatsTrafficCmd.MarkFlagRequired("repo")
	_ = repostatsTrafficCmd.MarkFlagRequired("owner")

	repostatsLanguagesCmd.Flags().StringVarP(&repostatsRepo, "repo", "r", "", "Repository name (required)")
	repostatsLanguagesCmd.Flags().StringVarP(&repostatsOwner, "owner", "o", "", "Owner name (required)")
	_ = repostatsLanguagesCmd.MarkFlagRequired("repo")
	_ = repostatsLanguagesCmd.MarkFlagRequired("owner")

	repostatsContributorsCmd.Flags().StringVarP(&repostatsRepo, "repo", "r", "", "Repository name (required)")
	repostatsContributorsCmd.Flags().StringVarP(&repostatsOwner, "owner", "o", "", "Owner name (required)")
	_ = repostatsContributorsCmd.MarkFlagRequired("repo")
	_ = repostatsContributorsCmd.MarkFlagRequired("owner")

	rootCmd.AddCommand(repostatsCmd)
}

func repostatsTraffic(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	traffic, err := client.GetTrafficData(repostatsOwner, repostatsRepo, api.GetTrafficDataOptions{
		StartDay: repostatsStart,
		EndDay:   repostatsEnd,
	})
	if err != nil {
		return fmt.Errorf("failed to get traffic data: %w", err)
	}

	fmt.Printf("Traffic Data for %s/%s\n", repostatsOwner, repostatsRepo)
	fmt.Printf("Total - Views: %d, Pulls: %d, Pushes: %d, Downloads: %d\n",
		traffic.Summary.IP, traffic.Summary.Pull, traffic.Summary.Push, traffic.Summary.DownloadZip)
	if len(traffic.Counts) > 0 {
		fmt.Println("\nDaily breakdown:")
		for _, d := range traffic.Counts {
			fmt.Printf("  Bucket %d: Views: %d, Pulls: %d, Pushes: %d, Downloads: %d\n",
				d.Bucket, d.IP, d.Pull, d.Push, d.DownloadZip)
		}
	}
	return nil
}

func repostatsLanguages(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	languages, err := client.GetRepoLanguages(repostatsOwner, repostatsRepo)
	if err != nil {
		return fmt.Errorf("failed to get languages: %w", err)
	}

	fmt.Printf("Languages for %s/%s\n", repostatsOwner, repostatsRepo)
	if languages != nil && languages.Languages != nil {
		for lang, bytes := range languages.Languages {
			fmt.Printf("  %s: %d bytes\n", lang, bytes)
		}
	}
	return nil
}

func repostatsContributors(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	contributors, err := client.GetRepoContributors(repostatsOwner, repostatsRepo)
	if err != nil {
		return fmt.Errorf("failed to get contributors: %w", err)
	}

	if len(contributors) == 0 {
		fmt.Println("No contributors found")
		return nil
	}

	fmt.Printf("Contributors for %s/%s\n", repostatsOwner, repostatsRepo)
	for _, c := range contributors {
		fmt.Printf("  %s (%s): %d contributions\n", c.Name, c.Email, c.Contributions)
	}
	return nil
}
