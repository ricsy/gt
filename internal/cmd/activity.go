package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	activityRepo     string
	activityOrg      string
	activityPublic   bool
	activityReceived bool
	activityNetwork  bool
	activityPage     int
	activityPerPage  int
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Manage activity events and watches",
}

var activityEventsCmd = &cobra.Command{
	Use:   "events [username]",
	Short: "List activity events",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  activityEvents,
}

var activityWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Manage watched repositories",
}

var activityWatchListCmd = &cobra.Command{
	Use:   "list [username]",
	Short: "List watched repositories",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  activityWatchList,
}

var activityWatchRepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Watch a repository",
	RunE:  activityWatchRepo,
}

var activityUnwatchRepoCmd = &cobra.Command{
	Use:   "unwatch",
	Short: "Unwatch a repository",
	RunE:  activityUnwatchRepo,
}

var activitySubscribersCmd = &cobra.Command{
	Use:   "subscribers",
	Short: "List repository watchers",
	RunE:  activitySubscribers,
}

func init() {
	rootCmd.AddCommand(activityCmd)
	activityCmd.AddCommand(activityEventsCmd, activityWatchCmd, activitySubscribersCmd)
	activityWatchCmd.AddCommand(activityWatchListCmd, activityWatchRepoCmd, activityUnwatchRepoCmd)

	addActivityPaginationFlags(activityEventsCmd)
	addActivityPaginationFlags(activityWatchListCmd)
	addActivityPaginationFlags(activitySubscribersCmd)

	activityEventsCmd.Flags().StringVar(&activityRepo, "repo", "", "Repository (owner/repo)")
	activityEventsCmd.Flags().StringVar(&activityOrg, "org", "", "Organization login")
	activityEventsCmd.Flags().BoolVar(&activityPublic, "public", false, "List only public events")
	activityEventsCmd.Flags().BoolVar(&activityReceived, "received", false, "List received events for a user")
	activityEventsCmd.Flags().BoolVar(&activityNetwork, "network", false, "List public network events for a repo")

	activityWatchRepoCmd.Flags().StringVar(&activityRepo, "repo", "", "Repository (owner/repo)")
	activityUnwatchRepoCmd.Flags().StringVar(&activityRepo, "repo", "", "Repository (owner/repo)")
	activitySubscribersCmd.Flags().StringVar(&activityRepo, "repo", "", "Repository (owner/repo)")
}

func addActivityPaginationFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&activityPage, "page", 0, "Page number")
	cmd.Flags().IntVar(&activityPerPage, "per-page", 0, "Items per page (max 100)")
}

func activityListOptions() api.ListActivityOptions {
	return api.ListActivityOptions{Page: activityPage, PerPage: activityPerPage}
}

func activityEvents(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	var events []api.Event
	switch {
	case activityRepo != "":
		owner, repoName, err := ResolveRepo(activityRepo)
		if err != nil {
			return err
		}
		if activityNetwork {
			events, err = client.ListNetworkEvents(owner, repoName, activityListOptions())
		} else {
			events, err = client.ListRepoEvents(owner, repoName, activityListOptions())
		}
		if err != nil {
			return err
		}
	case activityOrg != "":
		events, err = client.ListOrgEvents(activityOrg, activityListOptions())
		if err != nil {
			return err
		}
	case len(args) == 1:
		if activityReceived {
			events, err = client.ListUserReceivedEvents(args[0], activityPublic, activityListOptions())
		} else {
			events, err = client.ListUserEvents(args[0], activityPublic, activityListOptions())
		}
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("specify [username], --repo, or --org")
	}

	for _, event := range events {
		fmt.Printf("%s\t%s\t%s\n", event.ID, event.Type, event.CreatedAt)
	}
	return nil
}

func activityWatchList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	var repos []api.Repository
	if len(args) == 1 {
		repos, err = client.ListUserSubscriptions(args[0], activityListOptions())
	} else {
		repos, err = client.ListSubscriptions(activityListOptions())
	}
	return printRepositoriesResult(repos, err)
}

func activityWatchRepo(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(activityRepo)
	if err != nil {
		return err
	}
	if err := client.WatchRepo(owner, repoName); err != nil {
		return err
	}
	fmt.Printf("Watching %s/%s\n", owner, repoName)
	return nil
}

func activityUnwatchRepo(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(activityRepo)
	if err != nil {
		return err
	}
	if err := client.UnwatchRepo(owner, repoName); err != nil {
		return err
	}
	fmt.Printf("Unwatched %s/%s\n", owner, repoName)
	return nil
}

func activitySubscribers(cmd *cobra.Command, args []string) error {
	owner, repoName, client, err := resolveRepoClient(activityRepo)
	if err != nil {
		return err
	}
	users, err := client.ListSubscribers(owner, repoName, activityListOptions())
	return printUsersResult(users, err)
}
