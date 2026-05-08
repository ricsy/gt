package cmd

import (
	"fmt"
	"strconv"

	"github.com/ricsy/gt/pkg/api"
	"github.com/spf13/cobra"
)

var (
	userName      string
	userBlog      string
	userWeibo     string
	userBio       string
	userPage      int
	userPerPage   int
	keyTitle      string
	keyValue      string
	namespace     string
	namespaceMode string
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Long:  `Commands for working with Gitee users`,
}

var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show authenticated user",
	RunE:  userMe,
}

var userViewCmd = &cobra.Command{
	Use:   "view <username>",
	Short: "Show user profile",
	Args:  cobra.ExactArgs(1),
	RunE:  userView,
}

var userUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update authenticated user profile",
	RunE:  userUpdate,
}

var userFollowersCmd = &cobra.Command{
	Use:   "followers [username]",
	Short: "List user followers",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  userFollowers,
}

var userFollowingCmd = &cobra.Command{
	Use:   "following [username]",
	Short: "List users followed by a user",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  userFollowing,
}

var userFollowCmd = &cobra.Command{
	Use:   "follow <username>",
	Short: "Follow a user",
	Args:  cobra.ExactArgs(1),
	RunE:  userFollow,
}

var userUnfollowCmd = &cobra.Command{
	Use:   "unfollow <username>",
	Short: "Unfollow a user",
	Args:  cobra.ExactArgs(1),
	RunE:  userUnfollow,
}

var userKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage SSH keys",
}

var userKeyListCmd = &cobra.Command{
	Use:   "list [username]",
	Short: "List SSH keys",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  userKeyList,
}

var userKeyViewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "Show an SSH key",
	Args:  cobra.ExactArgs(1),
	RunE:  userKeyView,
}

var userKeyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an SSH key",
	RunE:  userKeyCreate,
}

var userKeyDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an SSH key",
	Args:  cobra.ExactArgs(1),
	RunE:  userKeyDelete,
}

var userNamespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage namespaces",
}

var userNamespaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List namespaces",
	RunE:  userNamespaceList,
}

var userNamespaceViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Show a namespace",
	RunE:  userNamespaceView,
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userMeCmd, userViewCmd, userUpdateCmd, userFollowersCmd, userFollowingCmd, userFollowCmd, userUnfollowCmd, userKeyCmd, userNamespaceCmd)
	userKeyCmd.AddCommand(userKeyListCmd, userKeyViewCmd, userKeyCreateCmd, userKeyDeleteCmd)
	userNamespaceCmd.AddCommand(userNamespaceListCmd, userNamespaceViewCmd)

	addUserPaginationFlags(userFollowersCmd)
	addUserPaginationFlags(userFollowingCmd)
	addUserPaginationFlags(userKeyListCmd)
	addUserPaginationFlags(userNamespaceListCmd)

	userUpdateCmd.Flags().StringVar(&userName, "name", "", "Display name")
	userUpdateCmd.Flags().StringVar(&userBlog, "blog", "", "Blog URL")
	userUpdateCmd.Flags().StringVar(&userWeibo, "weibo", "", "Weibo URL")
	userUpdateCmd.Flags().StringVar(&userBio, "bio", "", "Bio")

	userKeyCreateCmd.Flags().StringVar(&keyTitle, "title", "", "SSH key title")
	userKeyCreateCmd.Flags().StringVar(&keyValue, "key", "", "SSH public key")
	_ = userKeyCreateCmd.MarkFlagRequired("title")
	_ = userKeyCreateCmd.MarkFlagRequired("key")

	userNamespaceViewCmd.Flags().StringVar(&namespace, "path", "", "Namespace path")
	_ = userNamespaceViewCmd.MarkFlagRequired("path")
	userNamespaceListCmd.Flags().StringVar(&namespaceMode, "mode", "", "Namespace mode: project, intrant, all")
}

func addUserPaginationFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&userPage, "page", 0, "Page number")
	cmd.Flags().IntVar(&userPerPage, "per-page", 0, "Items per page (max 100)")
}

func userListOptions() api.ListUsersOptions {
	return api.ListUsersOptions{Page: userPage, PerPage: userPerPage}
}

func userMe(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	user, err := client.GetAuthenticatedUser()
	if err != nil {
		return err
	}
	printUser(user.UserBasic)
	return nil
}

func userView(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	user, err := client.GetUser(args[0])
	if err != nil {
		return err
	}
	printUser(*user)
	return nil
}

func userUpdate(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	user, err := client.UpdateAuthenticatedUser(api.UpdateUserOptions{
		Name:  userName,
		Blog:  userBlog,
		Weibo: userWeibo,
		Bio:   userBio,
	})
	if err != nil {
		return err
	}
	printUser(*user)
	return nil
}

func userFollowers(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	var users []api.User
	if len(args) == 1 {
		users, err = client.ListUserFollowers(args[0], userListOptions())
	} else {
		users, err = client.ListFollowers(userListOptions())
	}
	if err != nil {
		return err
	}
	printUsers(users)
	return nil
}

func userFollowing(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	var users []api.User
	if len(args) == 1 {
		users, err = client.ListUserFollowing(args[0], userListOptions())
	} else {
		users, err = client.ListFollowing(userListOptions())
	}
	if err != nil {
		return err
	}
	printUsers(users)
	return nil
}

func userFollow(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	if err := client.FollowUser(args[0]); err != nil {
		return err
	}
	fmt.Printf("Followed %s\n", args[0])
	return nil
}

func userUnfollow(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	if err := client.UnfollowUser(args[0]); err != nil {
		return err
	}
	fmt.Printf("Unfollowed %s\n", args[0])
	return nil
}

func userKeyList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	var keys []api.SSHKey
	if len(args) == 1 {
		keys, err = client.ListUserSSHKeys(args[0], userListOptions())
	} else {
		keys, err = client.ListSSHKeys(userListOptions())
	}
	if err != nil {
		return err
	}

	for _, key := range keys {
		fmt.Printf("%d\t%s\n", key.ID, key.Title)
	}
	return nil
}

func userKeyView(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid key id: %w", err)
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	key, err := client.GetSSHKey(id)
	if err != nil {
		return err
	}
	fmt.Printf("%d\t%s\n%s\n", key.ID, key.Title, key.Key)
	return nil
}

func userKeyCreate(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	key, err := client.CreateSSHKey(api.CreateSSHKeyOptions{Title: keyTitle, Key: keyValue})
	if err != nil {
		return err
	}
	fmt.Printf("Created SSH key: %d %s\n", key.ID, key.Title)
	return nil
}

func userKeyDelete(cmd *cobra.Command, args []string) error {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid key id: %w", err)
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	if err := client.DeleteSSHKey(id); err != nil {
		return err
	}
	fmt.Printf("Deleted SSH key: %d\n", id)
	return nil
}

func userNamespaceList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	namespaces, err := client.ListNamespaces(api.ListNamespacesOptions{
		Mode:    namespaceMode,
		Page:    userPage,
		PerPage: userPerPage,
	})
	if err != nil {
		return err
	}
	for _, ns := range namespaces {
		fmt.Printf("%s\t%s\n", ns.Path, ns.Name)
	}
	return nil
}

func userNamespaceView(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	ns, err := client.GetNamespace(namespace)
	if err != nil {
		return err
	}
	fmt.Printf("%s\t%s\nURL: %s\n", ns.Path, ns.Name, ns.HTMLURL)
	return nil
}

func printUser(user api.User) {
	fmt.Printf("%s (%s)\n", user.Login, user.Name)
	if user.Email != "" {
		fmt.Printf("Email: %s\n", user.Email)
	}
	if user.HTMLURL != "" {
		fmt.Printf("URL: %s\n", user.HTMLURL)
	}
}

func printUsers(users []api.User) {
	for _, user := range users {
		fmt.Printf("%s\t%s\n", user.Login, user.Name)
	}
}
