package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ricsy/gt/pkg/api"
	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/ricsy/gt/pkg/util"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repositories",
	Long:  `Manage Gitee repositories`,
}

var repoListOpts struct {
	Owner string
	Limit int
}

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories",
	RunE:  repoListCommand,
}

var repoViewCmd = &cobra.Command{
	Use:   "view [owner/repo]",
	Short: "View repository details",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  repoViewCommand,
}

var repoCreateOpts struct {
	Name        string
	Description string
	Private     bool
	Public      bool
}

var repoCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new repository",
	RunE:  repoCreateCommand,
}

var repoCloneCmd = &cobra.Command{
	Use:   "clone <repo> [directory]",
	Short: "Clone a repository",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  repoCloneCommand,
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoListCmd, repoViewCmd, repoCreateCmd, repoCloneCmd)

	repoListCmd.Flags().StringVar(&repoListOpts.Owner, "owner", "", "Owner username")
	repoListCmd.Flags().IntVar(&repoListOpts.Limit, "limit", 30, "Maximum number of repos to list")

	repoCreateCmd.Flags().StringVar(&repoCreateOpts.Name, "name", "", "Repository name")
	repoCreateCmd.Flags().StringVar(&repoCreateOpts.Description, "description", "", "Repository description")
	repoCreateCmd.Flags().BoolVar(&repoCreateOpts.Private, "private", false, "Create private repository")
	repoCreateCmd.Flags().BoolVar(&repoCreateOpts.Public, "public", false, "Create public repository")
	_ = repoCreateCmd.MarkFlagRequired("name")
}

func repoListCommand(cmd *cobra.Command, args []string) error {
	host := config.DefaultHost

	token, err := auth.GetToken(host)
	if err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	client := api.NewClient(host, token)

	var repos []api.Repository
	var err2 error

	if repoListOpts.Owner != "" {
		repos, err2 = client.ListUserRepos(repoListOpts.Owner)
	} else {
		repos, err2 = client.ListRepos()
	}

	if err2 != nil {
		return fmt.Errorf("failed to list repos: %w", err2)
	}

	limit := repoListOpts.Limit
	if limit <= 0 {
		limit = len(repos)
	}
	if limit > len(repos) {
		limit = len(repos)
	}

	for i := 0; i < limit; i++ {
		r := repos[i]
		vis := "public"
		if r.Private {
			vis = "private"
		}
		fmt.Printf("%s [%s]\n", r.FullName, vis)
		if r.Description != "" {
			fmt.Printf("  %s\n", r.Description)
		}
	}

	return nil
}

func repoViewCommand(cmd *cobra.Command, args []string) error {
	host := config.DefaultHost

	token, err := auth.GetToken(host)
	if err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	client := api.NewClient(host, token)

	var repo *api.Repository

	if len(args) == 1 {
		owner, repoName := util.SplitOwnerRepo(args[0])
		repo, err = client.GetRepo(owner, repoName)
		if err != nil {
			return fmt.Errorf("failed to get repo: %w", err)
		}
	} else {
		// List current user's repos and show the first one as default
		repos, err := client.ListRepos()
		if err != nil {
			return fmt.Errorf("failed to list repos: %w", err)
		}
		if len(repos) == 0 {
			fmt.Println("No repositories found")
			return nil
		}
		repo = &repos[0]
	}

	fmt.Printf("Name: %s\n", repo.FullName)
	fmt.Printf("Description: %s\n", repo.Description)
	fmt.Printf("URL: %s\n", repo.HTMLURL)
	fmt.Printf("Clone: %s\n", repo.CloneURL)
	fmt.Printf("Stars: %d | Forks: %d\n", repo.StarCount, repo.ForksCount)

	return nil
}

func repoCreateCommand(cmd *cobra.Command, args []string) error {
	host := config.DefaultHost

	token, err := auth.GetToken(host)
	if err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}

	client := api.NewClient(host, token)

	opts := api.CreateRepoOptions{
		Name:        repoCreateOpts.Name,
		Description: repoCreateOpts.Description,
		Private:     repoCreateOpts.Private && !repoCreateOpts.Public,
		AutoInit:    true,
	}

	repo, err := client.CreateRepo(opts)
	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}

	fmt.Printf("Repository created: %s\n", repo.HTMLURL)

	return nil
}

func repoCloneCommand(cmd *cobra.Command, args []string) error {
	repoArg := args[0]

	var directory string
	if len(args) > 1 {
		directory = args[1]
	}

	owner, repoName := util.SplitOwnerRepo(repoArg)

	cloneURL := fmt.Sprintf("https://gitee.com/%s/%s.git", owner, repoName)

	var gitArgs []string
	gitArgs = append(gitArgs, "clone", cloneURL)
	if directory != "" {
		gitArgs = append(gitArgs, directory)
	}

	gitExec := exec.Command("git", gitArgs...)
	gitExec.Stdout = os.Stdout
	gitExec.Stderr = os.Stderr

	err := gitExec.Run()
	if err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}
