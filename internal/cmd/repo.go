package cmd

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/ricsy/gt/pkg/api"
	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

const (
	repoPrimaryBranch = "main"
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

type repoCreateOptions struct {
	Name              string
	Description       string
	Homepage          string
	Private           bool
	Public            bool
	HasIssues         bool
	HasWiki           bool
	CanComment        bool
	AutoInit          bool
	GitignoreTemplate string
	LicenseTemplate   string
	Path              string
	Namespace         string
	CloneURLMode      string
}

var repoCreateOpts = repoCreateOptions{}

var repoCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new repository",
	RunE:  repoCreateCommand,
}

var repoDeleteCmd = &cobra.Command{
	Use:   "delete <repo>",
	Short: "Delete a repository",
	Args:  cobra.ExactArgs(1),
	RunE:  repoDeleteCommand,
}

var repoBootstrapOpts struct {
	RemoteName string
	Push       bool
}

var repoBootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Create a repository and connect the current git worktree",
	RunE:  repoBootstrapCommand,
}

var repoBranchOpts struct {
	Repo      string
	Sort      string
	Direction string
	Page      int
	PerPage   int
	Refs      string
}

var repoBranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Manage repository branches",
}

var repoBranchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repository branches",
	RunE:  repoBranchListCommand,
}

var repoBranchViewCmd = &cobra.Command{
	Use:   "view <branch>",
	Short: "View a repository branch",
	Args:  cobra.ExactArgs(1),
	RunE:  repoBranchViewCommand,
}

var repoBranchCreateCmd = &cobra.Command{
	Use:   "create <branch>",
	Short: "Create a repository branch",
	Args:  cobra.ExactArgs(1),
	RunE:  repoBranchCreateCommand,
}

var repoBranchProtectCmd = &cobra.Command{
	Use:   "protect <branch>",
	Short: "Protect a repository branch",
	Args:  cobra.ExactArgs(1),
	RunE:  repoBranchProtectCommand,
}

var repoBranchUnprotectCmd = &cobra.Command{
	Use:   "unprotect <branch>",
	Short: "Remove branch protection",
	Args:  cobra.ExactArgs(1),
	RunE:  repoBranchUnprotectCommand,
}

var repoCloneCmd = &cobra.Command{
	Use:   "clone <repo> [directory]",
	Short: "Clone a repository",
	Args:  cobra.RangeArgs(1, 2),
	RunE:  repoCloneCommand,
}

var repoCollaboratorOpts struct {
	Repo string
}

var repoCollaboratorCmd = &cobra.Command{
	Use:   "collaborator",
	Short: "Manage repository collaborators",
}

var repoCollaboratorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repository collaborators",
	RunE:  repoCollaboratorListCommand,
}

var repoCollaboratorViewCmd = &cobra.Command{
	Use:   "view <user>",
	Short: "Check if user is a collaborator",
	Args:  cobra.ExactArgs(1),
	RunE:  repoCollaboratorViewCommand,
}

var repoCollaboratorPermCmd = &cobra.Command{
	Use:   "perm <user>",
	Short: "Get collaborator permission",
	Args:  cobra.ExactArgs(1),
	RunE:  repoCollaboratorPermCommand,
}

var repoCollaboratorAddCmd = &cobra.Command{
	Use:   "add <user>",
	Short: "Add a collaborator to repository",
	Args:  cobra.ExactArgs(1),
	RunE:  repoCollaboratorAddCommand,
}

var repoCollaboratorRemoveCmd = &cobra.Command{
	Use:   "remove <user>",
	Short: "Remove a collaborator from repository",
	Args:  cobra.ExactArgs(1),
	RunE:  repoCollaboratorRemoveCommand,
}

var repoCollaboratorAddOpts struct {
	Permission string
}

var repoForkOpts struct {
	Repo    string
	Sort    string
	Page    int
	PerPage int
}

var repoForkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Manage repository forks",
}

var repoForkListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repository forks",
	RunE:  repoForkListCommand,
}

var repoForkCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Fork a repository",
	RunE:  repoForkCreateCommand,
}

var repoModeCmd = &cobra.Command{
	Use:   "mode",
	Short: "Manage repository scope mode",
}

var repoModeShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show active repository scope mode",
	RunE:  repoModeShowCommand,
}

var repoModePersonalCmd = &cobra.Command{
	Use:   "personal",
	Short: "Lock repository operations to the authenticated user",
	RunE:  repoModePersonalCommand,
}

var repoModeOrgCmd = &cobra.Command{
	Use:   "org <namespace>",
	Short: "Lock repository operations to an organization namespace",
	Args:  cobra.ExactArgs(1),
	RunE:  repoModeOrgCommand,
}

var repoModeClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear repository scope mode",
	RunE:  repoModeClearCommand,
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repoListCmd, repoViewCmd, repoCreateCmd, repoDeleteCmd, repoBootstrapCmd, repoBranchCmd, repoCloneCmd, repoCollaboratorCmd, repoForkCmd, repoModeCmd)
	repoBranchCmd.AddCommand(repoBranchListCmd, repoBranchViewCmd, repoBranchCreateCmd, repoBranchProtectCmd, repoBranchUnprotectCmd)
	repoCollaboratorCmd.AddCommand(repoCollaboratorListCmd, repoCollaboratorViewCmd, repoCollaboratorPermCmd, repoCollaboratorAddCmd, repoCollaboratorRemoveCmd)
	repoForkCmd.AddCommand(repoForkListCmd, repoForkCreateCmd)
	repoModeCmd.AddCommand(repoModeShowCmd, repoModePersonalCmd, repoModeOrgCmd, repoModeClearCmd)

	repoListCmd.Flags().StringVar(&repoListOpts.Owner, "owner", "", "Owner username")
	repoListCmd.Flags().IntVar(&repoListOpts.Limit, "limit", 30, "Maximum number of repos to list")

	repoCreateCmd.Flags().StringVar(&repoCreateOpts.Name, "name", "", "Repository name")
	repoCreateCmd.Flags().StringVar(&repoCreateOpts.Description, "description", "", "Repository description")
	repoCreateCmd.Flags().StringVar(&repoCreateOpts.Homepage, "homepage", "", "Repository homepage URL")
	repoCreateCmd.Flags().BoolVar(&repoCreateOpts.Private, "private", true, "Create private repository")
	repoCreateCmd.Flags().BoolVar(&repoCreateOpts.Public, "public", false, "Request public visibility (unsupported for personal repo creation)")
	repoCreateCmd.Flags().BoolVar(&repoCreateOpts.HasIssues, "has-issues", true, "Enable repository issues")
	repoCreateCmd.Flags().BoolVar(&repoCreateOpts.HasWiki, "has-wiki", true, "Enable repository wiki")
	repoCreateCmd.Flags().BoolVar(&repoCreateOpts.CanComment, "can-comment", true, "Allow repository comments")
	repoCreateCmd.Flags().BoolVar(&repoCreateOpts.AutoInit, "auto-init", false, "Initialize repository with README files")
	repoCreateCmd.Flags().StringVar(&repoCreateOpts.GitignoreTemplate, "gitignore-template", "", "Gitignore template name")
	repoCreateCmd.Flags().StringVar(&repoCreateOpts.LicenseTemplate, "license-template", "", "License template name")
	repoCreateCmd.Flags().StringVar(&repoCreateOpts.Path, "path", "", "Repository path")
	repoCreateCmd.Flags().StringVar(&repoCreateOpts.Namespace, "namespace", "", "Repository namespace")
	repoCreateCmd.Flags().StringVar(&repoCreateOpts.CloneURLMode, "clone-url-mode", "https", "Preferred clone URL mode for follow-up diagnostics: https or ssh")
	_ = repoCreateCmd.MarkFlagRequired("name")
	repoDeleteCmd.Flags().BoolVar(&repoDeleteOpts.Yes, "yes", false, "Delete without prompt only when the repository has no commit history")

	repoBootstrapCmd.Flags().StringVar(&repoCreateOpts.Name, "name", "", "Repository name")
	repoBootstrapCmd.Flags().StringVar(&repoCreateOpts.Description, "description", "", "Repository description")
	repoBootstrapCmd.Flags().StringVar(&repoCreateOpts.Homepage, "homepage", "", "Repository homepage URL")
	repoBootstrapCmd.Flags().BoolVar(&repoCreateOpts.Private, "private", true, "Create private repository")
	repoBootstrapCmd.Flags().BoolVar(&repoCreateOpts.Public, "public", false, "Request public visibility (unsupported for personal repo creation)")
	repoBootstrapCmd.Flags().BoolVar(&repoCreateOpts.HasIssues, "has-issues", true, "Enable repository issues")
	repoBootstrapCmd.Flags().BoolVar(&repoCreateOpts.HasWiki, "has-wiki", true, "Enable repository wiki")
	repoBootstrapCmd.Flags().BoolVar(&repoCreateOpts.CanComment, "can-comment", true, "Allow repository comments")
	repoBootstrapCmd.Flags().BoolVar(&repoCreateOpts.AutoInit, "auto-init", false, "Initialize repository with README files")
	repoBootstrapCmd.Flags().StringVar(&repoCreateOpts.GitignoreTemplate, "gitignore-template", "", "Gitignore template name")
	repoBootstrapCmd.Flags().StringVar(&repoCreateOpts.LicenseTemplate, "license-template", "", "License template name")
	repoBootstrapCmd.Flags().StringVar(&repoCreateOpts.Path, "path", "", "Repository path")
	repoBootstrapCmd.Flags().StringVar(&repoCreateOpts.Namespace, "namespace", "", "Repository namespace")
	repoBootstrapCmd.Flags().StringVar(&repoCreateOpts.CloneURLMode, "clone-url-mode", "https", "Remote URL mode to configure: https or ssh")
	repoBootstrapCmd.Flags().StringVar(&repoBootstrapOpts.RemoteName, "remote-name", "origin", "Remote name to create or update")
	repoBootstrapCmd.Flags().BoolVar(&repoBootstrapOpts.Push, "push", true, "Push the current branch after wiring the remote")
	_ = repoBootstrapCmd.MarkFlagRequired("name")

	addRepoBranchRepoFlag(repoBranchListCmd)
	addRepoBranchRepoFlag(repoBranchViewCmd)
	addRepoBranchRepoFlag(repoBranchCreateCmd)
	addRepoBranchRepoFlag(repoBranchProtectCmd)
	addRepoBranchRepoFlag(repoBranchUnprotectCmd)
	repoBranchListCmd.Flags().StringVar(&repoBranchOpts.Sort, "sort", "", "Sort by: name or updated")
	repoBranchListCmd.Flags().StringVar(&repoBranchOpts.Direction, "direction", "", "Sort direction: asc or desc")
	repoBranchListCmd.Flags().IntVar(&repoBranchOpts.Page, "page", 0, "Page number")
	repoBranchListCmd.Flags().IntVar(&repoBranchOpts.PerPage, "per-page", 0, "Items per page (max 100)")
	repoBranchCreateCmd.Flags().StringVar(&repoBranchOpts.Refs, "refs", repoPrimaryBranch, "Starting ref")

	addRepoFlag(repoCollaboratorListCmd)
	addRepoFlag(repoCollaboratorViewCmd)
	addRepoFlag(repoCollaboratorPermCmd)
	addRepoFlag(repoCollaboratorAddCmd)
	addRepoFlag(repoCollaboratorRemoveCmd)
	repoCollaboratorAddCmd.Flags().StringVar(&repoCollaboratorAddOpts.Permission, "permission", "push", "Collaborator permission (push, pull, admin)")

	addRepoForkFlag(repoForkListCmd)
	addRepoForkFlag(repoForkCreateCmd)
	repoForkListCmd.Flags().StringVar(&repoForkOpts.Sort, "sort", "", "Sort by: newest, oldest, stargazers")
	repoForkListCmd.Flags().IntVar(&repoForkOpts.Page, "page", 0, "Page number")
	repoForkListCmd.Flags().IntVar(&repoForkOpts.PerPage, "per-page", 0, "Items per page (max 100)")
}

func addRepoBranchRepoFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&repoBranchOpts.Repo, "repo", "", "Repository (owner/repo)")
}

func addRepoFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&repoCollaboratorOpts.Repo, "repo", "", "Repository (owner/repo)")
}

func addRepoForkFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&repoForkOpts.Repo, "repo", "", "Repository (owner/repo)")
}

func repoListCommand(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	var repos []api.Repository
	var err2 error

	if repoListOpts.Owner != "" {
		if err := enforceRepoScopeOwner(repoListOpts.Owner); err != nil {
			return err
		}
		scope, scopeErr := loadRepoScope()
		if scopeErr != nil {
			return scopeErr
		}
		if scope.Mode == repoScopeModeOrg && repoListOpts.Owner == scope.Namespace {
			repos, err2 = client.ListOrgRepos(repoListOpts.Owner, api.ListOrgReposOptions{})
		} else {
			repos, err2 = client.ListUserRepos(repoListOpts.Owner)
		}
	} else if scopedOwner, scopeErr := resolveScopedRepoOwner(); scopeErr != nil {
		return scopeErr
	} else if scopedOwner != "" {
		scope, scopeErr := loadRepoScope()
		if scopeErr != nil {
			return scopeErr
		}
		if scope.Mode == repoScopeModeOrg {
			repos, err2 = client.ListOrgRepos(scopedOwner, api.ListOrgReposOptions{})
		} else {
			repos, err2 = client.ListUserRepos(scopedOwner)
		}
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
	client, err := getClient()
	if err != nil {
		return err
	}

	var repo *api.Repository

	if len(args) == 1 {
		owner, repoName, err := ResolveRepo(args[0])
		if err != nil {
			return err
		}
		repo, err = client.GetRepo(owner, repoName)
		if err != nil {
			return fmt.Errorf("failed to get repo: %w", err)
		}
	} else {
		owner, repoName, err := resolveRepoFlag("")
		if err != nil {
			return err
		}
		repo, err = client.GetRepo(owner, repoName)
		if err != nil {
			return fmt.Errorf("failed to get repo: %w", err)
		}
	}

	fmt.Printf("Name: %s\n", repo.FullName)
	fmt.Printf("Description: %s\n", repo.Description)
	fmt.Printf("URL: %s\n", repo.HTMLURL)
	fmt.Printf("Clone: %s\n", repo.CloneURL)
	fmt.Printf("Stars: %d | Forks: %d\n", repo.StarCount, repo.ForksCount)

	return nil
}

func repoCreateCommand(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	opts, err := buildCreateRepoOptions(cmd)
	if err != nil {
		return err
	}

	repo, err := client.CreateRepo(opts)
	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}

	fmt.Printf("Repository created: %s\n", repo.HTMLURL)
	printRepoCreatePushDiagnostics(cmd, repo, repoCreateOpts.CloneURLMode)

	return nil
}

// buildCreateRepoOptions aligns CLI flags with the current personal repository API contract.
func buildCreateRepoOptions(cmd *cobra.Command) (api.CreateRepoOptions, error) {
	opts := api.CreateRepoOptions{
		Name:        repoCreateOpts.Name,
		Description: repoCreateOpts.Description,
		Homepage:    repoCreateOpts.Homepage,
		Private:     true,
		AutoInit:    repoCreateOpts.AutoInit,
		Path:        repoCreateOpts.Path,
	}

	namespace, err := resolveRepoCreateNamespace(repoCreateOpts.Namespace)
	if err != nil {
		return api.CreateRepoOptions{}, err
	}
	opts.Namespace = namespace

	if cmd.Flags().Changed("public") && repoCreateOpts.Public {
		return api.CreateRepoOptions{}, fmt.Errorf("public repositories are not supported by the current user repo API; omit --public")
	}
	if cmd.Flags().Changed("private") {
		if !repoCreateOpts.Private {
			return api.CreateRepoOptions{}, fmt.Errorf("personal repositories currently require --private=true")
		}
		opts.Private = repoCreateOpts.Private
	}
	if cmd.Flags().Changed("has-issues") {
		opts.HasIssues = api.BoolPtr(repoCreateOpts.HasIssues)
	}
	if cmd.Flags().Changed("has-wiki") {
		opts.HasWiki = api.BoolPtr(repoCreateOpts.HasWiki)
	}
	if cmd.Flags().Changed("can-comment") {
		opts.CanComment = api.BoolPtr(repoCreateOpts.CanComment)
	}
	if cmd.Flags().Changed("gitignore-template") {
		opts.GitignoreTemplate = repoCreateOpts.GitignoreTemplate
	}
	if cmd.Flags().Changed("license-template") {
		opts.LicenseTemplate = repoCreateOpts.LicenseTemplate
	}

	return opts, nil
}

func repoBranchListCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoBranchOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	branches, err := client.ListBranches(owner, repoName, api.ListBranchesOptions{
		Sort:      repoBranchOpts.Sort,
		Direction: repoBranchOpts.Direction,
		Page:      repoBranchOpts.Page,
		PerPage:   repoBranchOpts.PerPage,
	})
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}
	for _, branch := range branches {
		printBranch(branch.Name, branch.Commit.SHA, branch.Protected)
	}
	return nil
}

func repoBranchViewCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoBranchOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	branch, err := client.GetBranch(owner, repoName, args[0])
	if err != nil {
		return fmt.Errorf("failed to get branch: %w", err)
	}
	printBranch(branch.Name, branch.Commit.SHA, branch.Protected)
	return nil
}

func repoBranchCreateCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoBranchOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	branch, err := client.CreateBranch(owner, repoName, api.CreateBranchOptions{
		Refs:       repoBranchOpts.Refs,
		BranchName: args[0],
	})
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	printBranch(branch.Name, branch.Commit.SHA, branch.Protected)
	return nil
}

func repoBranchProtectCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoBranchOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	branch, err := client.ProtectBranch(owner, repoName, args[0])
	if err != nil {
		return fmt.Errorf("failed to protect branch: %w", err)
	}
	printBranch(branch.Name, branch.Commit.SHA, branch.Protected)
	return nil
}

func repoBranchUnprotectCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoBranchOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	if err := client.UnprotectBranch(owner, repoName, args[0]); err != nil {
		return fmt.Errorf("failed to unprotect branch: %w", err)
	}
	fmt.Printf("Unprotected branch: %s\n", args[0])
	return nil
}

func repoCollaboratorListCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoCollaboratorOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	collabs, err := client.ListCollaborators(owner, repoName)
	if err != nil {
		return fmt.Errorf("failed to list collaborators: %w", err)
	}
	for _, c := range collabs {
		cmd.Printf("%s (%s)\n", c.Login, c.Name)
	}
	return nil
}

func repoCollaboratorViewCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoCollaboratorOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	collab, err := client.GetCollaborator(owner, repoName, args[0])
	if err != nil {
		return fmt.Errorf("failed to get collaborator: %w", err)
	}
	if collab.Login == "" && collab.Name == "" {
		cmd.Printf("%s is a collaborator\n", args[0])
		return nil
	}
	cmd.Printf("%s (%s)\n", collab.Login, collab.Name)
	return nil
}

func repoCollaboratorPermCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoCollaboratorOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	perm, err := client.GetCollaboratorPermission(owner, repoName, args[0])
	if err != nil {
		return fmt.Errorf("failed to get collaborator permission: %w", err)
	}
	cmd.Printf("Permission: %s\n", perm.Permission)
	return nil
}

func repoCollaboratorAddCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoCollaboratorOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	if err := client.AddCollaborator(owner, repoName, args[0], repoCollaboratorAddOpts.Permission); err != nil {
		return fmt.Errorf("failed to add collaborator: %w", err)
	}
	cmd.Printf("Added collaborator: %s\n", args[0])
	return nil
}

func repoCollaboratorRemoveCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoCollaboratorOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	if err := client.RemoveCollaborator(owner, repoName, args[0]); err != nil {
		return fmt.Errorf("failed to remove collaborator: %w", err)
	}
	cmd.Printf("Removed collaborator: %s\n", args[0])
	return nil
}

func repoForkListCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoForkOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	forks, err := client.ListForks(owner, repoName, api.ListForksOptions{
		Sort:    repoForkOpts.Sort,
		Page:    repoForkOpts.Page,
		PerPage: repoForkOpts.PerPage,
	})
	if err != nil {
		return fmt.Errorf("failed to list forks: %w", err)
	}
	for _, f := range forks {
		cmd.Printf("%s\n", f.FullName)
	}
	return nil
}

func repoForkCreateCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := resolveRepoFlag(repoForkOpts.Repo)
	if err != nil {
		return err
	}
	client, err := getClient()
	if err != nil {
		return err
	}
	fork, err := client.ForkRepository(owner, repoName)
	if err != nil {
		return fmt.Errorf("failed to fork repository: %w", err)
	}
	cmd.Printf("Forked: %s\n", fork.HTMLURL)
	return nil
}

func repoModeShowCommand(cmd *cobra.Command, args []string) error {
	scope, err := loadRepoScope()
	if err != nil {
		return err
	}

	cmd.Printf("Repo scope mode: %s\n", describeRepoScope(scope))
	return nil
}

func repoModePersonalCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	cfg.RepoScopeMode = repoScopeModePersonal
	cfg.RepoScopeNamespace = ""
	if err := config.SaveConfig(cfg); err != nil {
		return err
	}

	cmd.Printf("Repo scope mode: %s\n", repoScopeModePersonal)
	return nil
}

func repoModeOrgCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	cfg.RepoScopeMode = repoScopeModeOrg
	cfg.RepoScopeNamespace = args[0]
	if err := config.SaveConfig(cfg); err != nil {
		return err
	}

	cmd.Printf("Repo scope mode: %s:%s\n", repoScopeModeOrg, args[0])
	return nil
}

func repoModeClearCommand(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	cfg.RepoScopeMode = ""
	cfg.RepoScopeNamespace = ""
	if err := config.SaveConfig(cfg); err != nil {
		return err
	}

	cmd.Printf("Repo scope mode: none\n")
	return nil
}

func printBranch(name, commit string, protected bool) {
	protection := "unprotected"
	if protected {
		protection = "protected"
	}
	fmt.Printf("%s\t%s\t%s\n", name, commit, protection)
}

func repoCloneCommand(cmd *cobra.Command, args []string) error {
	owner, repoName, err := ResolveRepo(args[0])
	if err != nil {
		return err
	}

	var directory string
	if len(args) > 1 {
		directory = args[1]
	}

	cloneURL := config.RepoGitHTTPSURL(resolveCommandHost(), owner, repoName)
	if authenticatedURL, err := buildAuthenticatedCloneURL(resolveCommandHost(), cloneURL); err == nil && authenticatedURL != "" {
		cloneURL = authenticatedURL
	}

	var gitArgs []string
	gitArgs = append(gitArgs, "clone", cloneURL)
	if directory != "" {
		gitArgs = append(gitArgs, directory)
	}

	gitExec := exec.Command("git", gitArgs...)
	gitExec.Stdout = os.Stdout
	gitExec.Stderr = os.Stderr

	err = gitExec.Run()
	if err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}

func buildAuthenticatedCloneURL(host, cloneURL string) (string, error) {
	authInfo, err := auth.GetAuth(host)
	if err != nil || authInfo.Token == "" || authInfo.User == "" {
		return "", err
	}

	parsed, err := url.Parse(cloneURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse clone URL: %w", err)
	}
	parsed.User = url.UserPassword(authInfo.User, authInfo.Token)
	return parsed.String(), nil
}
