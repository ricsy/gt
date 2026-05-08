package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	enterpriseName        string
	enterprisePage        int
	enterprisePerPage     int
	enterpriseAdmin       bool
	enterpriseRole        string
	enterpriseMemberName  string
	enterpriseEmail       string
	enterpriseQueryType   string
	enterpriseQueryValue  string
	enterpriseActive      bool
	enterpriseSearch      string
	enterpriseRepoType    string
	enterpriseDirect      bool
	enterprisePRState     string
	enterprisePRRepo      string
	enterprisePRLabels    string
	enterprisePRSort      string
	enterprisePRDirection string
)

var enterpriseCmd = &cobra.Command{
	Use:   "enterprise",
	Short: "Manage enterprises",
	Long:  `Commands for working with Gitee enterprises`,
}

var enterpriseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List authenticated user's enterprises",
	RunE:  enterpriseList,
}

var enterpriseViewCmd = &cobra.Command{
	Use:   "view <enterprise>",
	Short: "Show enterprise details",
	Args:  cobra.ExactArgs(1),
	RunE:  enterpriseView,
}

var enterpriseMemberCmd = &cobra.Command{
	Use:   "member",
	Short: "Manage enterprise members",
}

var enterpriseMemberListCmd = &cobra.Command{
	Use:   "list <enterprise>",
	Short: "List enterprise members",
	Args:  cobra.ExactArgs(1),
	RunE:  enterpriseMemberList,
}

var enterpriseMemberSearchCmd = &cobra.Command{
	Use:   "search <enterprise>",
	Short: "Search an enterprise member",
	Args:  cobra.ExactArgs(1),
	RunE:  enterpriseMemberSearch,
}

var enterpriseMemberViewCmd = &cobra.Command{
	Use:   "view <enterprise> <username>",
	Short: "Show enterprise member details",
	Args:  cobra.ExactArgs(2),
	RunE:  enterpriseMemberView,
}

var enterpriseMemberAddCmd = &cobra.Command{
	Use:   "add <enterprise>",
	Short: "Add or invite an enterprise member",
	Args:  cobra.ExactArgs(1),
	RunE:  enterpriseMemberAdd,
}

var enterpriseMemberUpdateCmd = &cobra.Command{
	Use:   "update <enterprise> <username>",
	Short: "Update an enterprise member",
	Args:  cobra.ExactArgs(2),
	RunE:  enterpriseMemberUpdate,
}

var enterpriseMemberRemoveCmd = &cobra.Command{
	Use:   "remove <enterprise> <username>",
	Short: "Remove an enterprise member",
	Args:  cobra.ExactArgs(2),
	RunE:  enterpriseMemberRemove,
}

var enterpriseRepoCmd = &cobra.Command{
	Use:   "repo <enterprise>",
	Short: "List enterprise repositories",
	Args:  cobra.ExactArgs(1),
	RunE:  enterpriseRepoList,
}

var enterprisePRCmd = &cobra.Command{
	Use:   "pr <enterprise>",
	Short: "List enterprise pull requests",
	Args:  cobra.ExactArgs(1),
	RunE:  enterprisePRList,
}

func init() {
	rootCmd.AddCommand(enterpriseCmd)
	enterpriseCmd.AddCommand(enterpriseListCmd, enterpriseViewCmd, enterpriseMemberCmd, enterpriseRepoCmd, enterprisePRCmd)
	enterpriseMemberCmd.AddCommand(enterpriseMemberListCmd, enterpriseMemberSearchCmd, enterpriseMemberViewCmd, enterpriseMemberAddCmd, enterpriseMemberUpdateCmd, enterpriseMemberRemoveCmd)

	addEnterprisePaginationFlags(enterpriseListCmd)
	addEnterprisePaginationFlags(enterpriseMemberListCmd)
	addEnterprisePaginationFlags(enterpriseRepoCmd)
	addEnterprisePaginationFlags(enterprisePRCmd)

	enterpriseListCmd.Flags().BoolVar(&enterpriseAdmin, "admin", false, "Only list enterprises administered by the authenticated user")
	enterpriseMemberListCmd.Flags().StringVar(&enterpriseRole, "role", "", "Filter members by role: all, admin, member")

	enterpriseMemberSearchCmd.Flags().StringVar(&enterpriseQueryType, "query-type", "", "Search type: username or email")
	enterpriseMemberSearchCmd.Flags().StringVar(&enterpriseQueryValue, "query-value", "", "Search value")
	_ = enterpriseMemberSearchCmd.MarkFlagRequired("query-type")
	_ = enterpriseMemberSearchCmd.MarkFlagRequired("query-value")

	addEnterpriseMemberWriteFlags(enterpriseMemberAddCmd)
	enterpriseMemberAddCmd.Flags().StringVar(&enterpriseEmail, "email", "", "Member email")
	enterpriseMemberAddCmd.Flags().StringVar(&enterpriseMemberName, "username", "", "Member username")

	addEnterpriseMemberWriteFlags(enterpriseMemberUpdateCmd)
	enterpriseMemberUpdateCmd.Flags().BoolVar(&enterpriseActive, "active", false, "Allow member to access enterprise resources")

	enterpriseRepoCmd.Flags().StringVar(&enterpriseSearch, "search", "", "Repository search string")
	enterpriseRepoCmd.Flags().StringVar(&enterpriseRepoType, "type", "", "Repository type: all, public, internal, private")
	enterpriseRepoCmd.Flags().BoolVar(&enterpriseDirect, "direct", false, "Only list direct enterprise repositories")

	enterprisePRCmd.Flags().StringVar(&enterprisePRState, "state", "", "Pull request state: open, closed, merged, all")
	enterprisePRCmd.Flags().StringVar(&enterprisePRRepo, "repo", "", "Repository path")
	enterprisePRCmd.Flags().StringVar(&enterprisePRLabels, "labels", "", "Comma-separated labels")
	enterprisePRCmd.Flags().StringVar(&enterprisePRSort, "sort", "", "Sort field")
	enterprisePRCmd.Flags().StringVar(&enterprisePRDirection, "direction", "", "Sort direction: asc or desc")
}

func addEnterprisePaginationFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&enterprisePage, "page", 0, "Page number")
	cmd.Flags().IntVar(&enterprisePerPage, "per-page", 0, "Items per page (max 100)")
}

func addEnterpriseMemberWriteFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&enterpriseRole, "role", "", "Enterprise role: admin, member, outsourced")
	cmd.Flags().StringVar(&enterpriseName, "name", "", "Member display name or remark")
}

func enterpriseList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	opts := api.ListEnterprisesOptions{Page: enterprisePage, PerPage: enterprisePerPage}
	if cmd.Flags().Changed("admin") {
		opts.Admin = api.BoolPtr(enterpriseAdmin)
	}
	enterprises, err := client.ListEnterprises(opts)
	if err != nil {
		return err
	}
	for _, enterprise := range enterprises {
		printEnterprise(enterprise)
	}
	return nil
}

func enterpriseView(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	enterprise, err := client.GetEnterprise(args[0])
	if err != nil {
		return err
	}
	printEnterprise(*enterprise)
	return nil
}

func enterpriseMemberList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	members, err := client.ListEnterpriseMembers(args[0], api.ListEnterpriseMembersOptions{
		Role:    enterpriseRole,
		Page:    enterprisePage,
		PerPage: enterprisePerPage,
	})
	if err != nil {
		return err
	}
	for _, member := range members {
		printEnterpriseMember(member)
	}
	return nil
}

func enterpriseMemberSearch(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	member, err := client.SearchEnterpriseMember(args[0], api.SearchEnterpriseMemberOptions{
		QueryType:  enterpriseQueryType,
		QueryValue: enterpriseQueryValue,
	})
	if err != nil {
		return err
	}
	printEnterpriseMember(*member)
	return nil
}

func enterpriseMemberView(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	member, err := client.GetEnterpriseMember(args[0], args[1])
	if err != nil {
		return err
	}
	printEnterpriseMember(*member)
	return nil
}

func enterpriseMemberAdd(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	if err := client.AddEnterpriseMember(args[0], api.AddEnterpriseMemberOptions{
		Username: enterpriseMemberName,
		Email:    enterpriseEmail,
		Role:     enterpriseRole,
		Name:     enterpriseName,
	}); err != nil {
		return err
	}
	fmt.Println("Enterprise member invitation submitted")
	return nil
}

func enterpriseMemberUpdate(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	opts := api.UpdateEnterpriseMemberOptions{
		Role: enterpriseRole,
		Name: enterpriseName,
	}
	if cmd.Flags().Changed("active") {
		opts.Active = api.BoolPtr(enterpriseActive)
	}
	member, err := client.UpdateEnterpriseMember(args[0], args[1], opts)
	if err != nil {
		return err
	}
	printEnterpriseMember(*member)
	return nil
}

func enterpriseMemberRemove(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	if err := client.RemoveEnterpriseMember(args[0], args[1]); err != nil {
		return err
	}
	fmt.Printf("Removed enterprise member: %s\n", args[1])
	return nil
}

func enterpriseRepoList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	opts := api.ListEnterpriseReposOptions{
		Search:  enterpriseSearch,
		Type:    enterpriseRepoType,
		Page:    enterprisePage,
		PerPage: enterprisePerPage,
	}
	if cmd.Flags().Changed("direct") {
		opts.Direct = api.BoolPtr(enterpriseDirect)
	}
	repos, err := client.ListEnterpriseRepos(args[0], opts)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		fmt.Printf("%s\t%s\n", repo.FullName, repo.Description)
	}
	return nil
}

func enterprisePRList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	prs, err := client.ListEnterprisePullRequests(args[0], api.ListEnterprisePullRequestsOptions{
		State:     enterprisePRState,
		Repo:      enterprisePRRepo,
		Labels:    enterprisePRLabels,
		Sort:      enterprisePRSort,
		Direction: enterprisePRDirection,
		Page:      enterprisePage,
		PerPage:   enterprisePerPage,
	})
	if err != nil {
		return err
	}
	for _, pr := range prs {
		fmt.Printf("#%d\t%s\t%s\n", pr.Number, pr.State, pr.Title)
	}
	return nil
}

func printEnterprise(enterprise api.EnterpriseBasic) {
	fmt.Printf("%s\t%s\n", enterprise.Path, enterprise.Name)
}

func printEnterpriseMember(member api.EnterpriseMember) {
	username := ""
	name := ""
	if member.User != nil {
		username = member.User.Login
		name = member.User.Name
	}
	fmt.Printf("%s\t%s\t%s\tactive=%t\n", username, name, member.Role, member.Active)
}
