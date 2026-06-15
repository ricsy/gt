package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	gitDataRepo      string
	gitDataRecursive bool
)

var gitDataCmd = &cobra.Command{
	Use:   "gitdata",
	Short: "Inspect Git data",
	Long:  `Commands for reading Gitee Git data objects`,
}

var gitDataBlobCmd = &cobra.Command{
	Use:   "blob <sha>",
	Short: "Get a Git blob",
	Args:  cobra.ExactArgs(1),
	RunE:  gitDataBlob,
}

var gitDataTreeCmd = &cobra.Command{
	Use:   "tree <sha>",
	Short: "Get a Git tree",
	Args:  cobra.ExactArgs(1),
	RunE:  gitDataTree,
}

var gitDataMetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Get Gitee metrics",
	RunE:  gitDataMetrics,
}

func init() {
	rootCmd.AddCommand(gitDataCmd)
	gitDataCmd.AddCommand(gitDataBlobCmd, gitDataTreeCmd, gitDataMetricsCmd)

	gitDataBlobCmd.Flags().StringVar(&gitDataRepo, "repo", "", "Repository (owner/repo)")
	gitDataTreeCmd.Flags().StringVar(&gitDataRepo, "repo", "", "Repository (owner/repo)")
	gitDataTreeCmd.Flags().BoolVar(&gitDataRecursive, "recursive", false, "Fetch tree recursively")
	gitDataMetricsCmd.Flags().StringVar(&gitDataRepo, "repo", "", "Repository (owner/repo)")
}

func gitDataBlob(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(gitDataRepo)
	if err != nil {
		return err
	}
	blob, err := client.GetBlob(owner, repoName, args[0])
	if err != nil {
		return err
	}
	fmt.Printf("SHA: %s\nSize: %d\nEncoding: %s\n", blob.SHA, blob.Size, blob.Encoding)
	if blob.Content != "" {
		fmt.Println(blob.Content)
	}
	return nil
}

func gitDataTree(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(gitDataRepo)
	if err != nil {
		return err
	}
	recursive := 0
	if gitDataRecursive {
		recursive = 1
	}
	tree, err := client.GetTree(owner, repoName, args[0], api.GetTreeOptions{Recursive: recursive})
	if err != nil {
		return err
	}
	for _, entry := range tree.Tree {
		fmt.Printf("%s\t%s\t%s\n", entry.Type, entry.SHA, entry.Path)
	}
	return nil
}

func gitDataMetrics(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(gitDataRepo)
	if err != nil {
		return err
	}
	metrics, err := client.GetGiteeMetrics(owner, repoName)
	if err != nil {
		return err
	}
	fmt.Printf("Score: %d\nCreated: %s\n", metrics.TotalScore, metrics.CreatedAt)
	if metrics.Data != "" {
		fmt.Println(metrics.Data)
	}
	return nil
}
