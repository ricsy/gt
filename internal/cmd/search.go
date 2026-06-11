package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	searchQuery    string
	searchOwner    string
	searchRepo     string
	searchLanguage string
	searchSort     string
	searchOrder    string
	searchState    string
	searchLabel    string
	searchAuthor   string
	searchAssignee string
	searchFork     bool
	searchPage     int
	searchPerPage  int
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search repositories, issues, and users",
}

var searchReposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Search repositories",
	RunE:  searchRepos,
}

var searchIssuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Search issues",
	RunE:  searchIssues,
}

var searchUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Search users",
	RunE:  searchUsers,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.AddCommand(searchReposCmd, searchIssuesCmd, searchUsersCmd)

	searchReposCmd.Flags().StringVar(&searchQuery, "q", "", "Search keyword (required)")
	searchReposCmd.Flags().StringVar(&searchOwner, "owner", "", "Filter by owner")
	searchReposCmd.Flags().StringVar(&searchLanguage, "language", "", "Filter by language")
	searchReposCmd.Flags().StringVar(&searchSort, "sort", "", "Sort by: last_push_at, stars_count, forks_count, watches_count")
	searchReposCmd.Flags().StringVar(&searchOrder, "order", "", "Sort order: asc, desc")
	searchReposCmd.Flags().BoolVar(&searchFork, "fork", false, "Include forked repos")
	searchReposCmd.Flags().IntVar(&searchPage, "page", 0, "Page number")
	searchReposCmd.Flags().IntVar(&searchPerPage, "per-page", 0, "Items per page (max 100)")

	searchIssuesCmd.Flags().StringVar(&searchQuery, "q", "", "Search keyword (required)")
	searchIssuesCmd.Flags().StringVar(&searchOwner, "owner", "", "Filter by owner")
	searchIssuesCmd.Flags().StringVar(&searchRepo, "repo", "", "Filter by repo")
	searchIssuesCmd.Flags().StringVar(&searchLanguage, "language", "", "Filter by language")
	searchIssuesCmd.Flags().StringVar(&searchLabel, "label", "", "Filter by label")
	searchIssuesCmd.Flags().StringVar(&searchState, "state", "", "Filter by state: open, progressing, closed, rejected")
	searchIssuesCmd.Flags().StringVar(&searchAuthor, "author", "", "Filter by author")
	searchIssuesCmd.Flags().StringVar(&searchAssignee, "assignee", "", "Filter by assignee")
	searchIssuesCmd.Flags().StringVar(&searchSort, "sort", "", "Sort by: created_at, updated_at, notes_count")
	searchIssuesCmd.Flags().StringVar(&searchOrder, "order", "", "Sort order: asc, desc")
	searchIssuesCmd.Flags().IntVar(&searchPage, "page", 0, "Page number")
	searchIssuesCmd.Flags().IntVar(&searchPerPage, "per-page", 0, "Items per page (max 100)")

	searchUsersCmd.Flags().StringVar(&searchQuery, "q", "", "Search keyword (required)")
	searchUsersCmd.Flags().StringVar(&searchSort, "sort", "", "Sort by: joined_at")
	searchUsersCmd.Flags().StringVar(&searchOrder, "order", "", "Sort order: asc, desc")
	searchUsersCmd.Flags().IntVar(&searchPage, "page", 0, "Page number")
	searchUsersCmd.Flags().IntVar(&searchPerPage, "per-page", 0, "Items per page (max 100)")
}

func searchRepos(cmd *cobra.Command, args []string) error {
	if searchQuery == "" {
		return fmt.Errorf("q is required: use --q")
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	forkPtr := api.BoolPtr(searchFork)
	repos, err := client.SearchRepos(api.SearchReposOptions{
		Q:        searchQuery,
		Owner:    searchOwner,
		Fork:     forkPtr,
		Language: searchLanguage,
		Sort:     searchSort,
		Order:    searchOrder,
		Page:     searchPage,
		PerPage:  searchPerPage,
	})
	if err != nil {
		return err
	}

	if len(repos) == 0 {
		fmt.Println("No repositories found")
		return nil
	}

	for _, r := range repos {
		fmt.Printf("%s - %s\n", r.FullName, r.Description)
	}
	return nil
}

func searchIssues(cmd *cobra.Command, args []string) error {
	if searchQuery == "" {
		return fmt.Errorf("q is required: use --q")
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	issues, err := client.SearchIssues(api.SearchIssuesOptions{
		Q:        searchQuery,
		Repo:     searchRepo,
		Language: searchLanguage,
		Label:    searchLabel,
		State:    searchState,
		Author:   searchAuthor,
		Assignee: searchAssignee,
		Sort:     searchSort,
		Order:    searchOrder,
		Page:     searchPage,
		PerPage:  searchPerPage,
	})
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		fmt.Println("No issues found")
		return nil
	}

	for _, i := range issues {
		fmt.Printf("#%s %s [%s]\n", i.Number, i.Title, i.State)
	}
	return nil
}

func searchUsers(cmd *cobra.Command, args []string) error {
	if searchQuery == "" {
		return fmt.Errorf("q is required: use --q")
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	users, err := client.SearchUsers(api.SearchUsersOptions{
		Q:       searchQuery,
		Sort:    searchSort,
		Order:   searchOrder,
		Page:    searchPage,
		PerPage: searchPerPage,
	})
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return nil
	}

	for _, u := range users {
		fmt.Printf("%s - %s\n", u.Login, u.Name)
	}
	return nil
}
