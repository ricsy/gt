package api

import (
	"net/http"
	"strconv"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// Repository is an alias for response.Repository
type Repository = response.Repository

// CreateRepoOptions is an alias for response.CreateRepoOptions
type CreateRepoOptions = response.CreateRepoOptions

// UpdateRepoOptions is an alias for response.UpdateRepoOptions
type UpdateRepoOptions = response.UpdateRepoOptions

// Branch is an alias for response.Branch
type Branch = response.Branch

// CompleteBranch is an alias for response.CompleteBranch
type CompleteBranch = response.CompleteBranch

// ListBranchesOptions is an alias for response.ListBranchesOptions
type ListBranchesOptions = response.ListBranchesOptions

// CreateBranchOptions is an alias for response.CreateBranchOptions
type CreateBranchOptions = response.CreateBranchOptions

type ListRepoCommitsOptions = response.ListRepoCommitsOptions

type BranchCommit = response.BranchCommit

type BranchCommitActor = response.BranchCommitActor

type BranchCommitDetail = response.BranchCommitDetail

// Collaborator is an alias for response.Collaborator
type Collaborator = response.Collaborator

// CollaboratorPermission is an alias for response.CollaboratorPermission
type CollaboratorPermission = response.CollaboratorPermission

// ForkRepository is an alias for response.ForkRepository
type ForkRepository = response.ForkRepository

// ListForksOptions is an alias for response.ListForksOptions
type ListForksOptions = response.ListForksOptions

type RepoCommitHistorySummary struct {
	Count  int
	Latest *BranchCommit
}

// ListRepos lists user repositories
func (c *Client) ListRepos() ([]Repository, error) {
	var repos []Repository
	err := c.DoFromEndpoint(UserRepos.List, nil, nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// GetRepo gets a single repository
func (c *Client) GetRepo(owner, repo string) (*Repository, error) {
	var result Repository
	err := c.DoFromEndpoint(Repo.Get, []interface{}{owner, repo}, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateRepo creates a new repository
func (c *Client) CreateRepo(opts CreateRepoOptions) (*Repository, error) {
	var result Repository
	err := c.DoFromEndpoint(UserRepos.Create, nil, opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUserRepos lists repositories for a specific user
func (c *Client) ListUserRepos(username string) ([]Repository, error) {
	var repos []Repository
	err := c.DoFromEndpoint(UserReposByName.List, []interface{}{username}, nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// UpdateRepo updates a repository
func (c *Client) UpdateRepo(owner, repo string, opts UpdateRepoOptions) (*Repository, error) {
	var result Repository
	err := c.DoFromEndpoint(Repo.Update, []interface{}{owner, repo}, opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteRepo deletes a repository.
func (c *Client) DeleteRepo(owner, repo string) error {
	return c.DoFromEndpoint(Repo.Delete, []interface{}{owner, repo}, nil, nil)
}

func (c *Client) GetRepoCommitHistorySummary(owner, repo, ref string) (*RepoCommitHistorySummary, error) {
	commits, headers, err := c.listRepoCommits(owner, repo, ListRepoCommitsOptions{
		SHA:     ref,
		Page:    1,
		PerPage: 100,
	})
	if err != nil {
		return nil, err
	}

	summary := &RepoCommitHistorySummary{Count: len(commits)}
	if len(commits) > 0 {
		summary.Latest = &commits[0]
	}

	if totalCount := parseHeaderInt(headers, "X-Total-Count"); totalCount > 0 {
		summary.Count = totalCount
		return summary, nil
	}

	if len(commits) < 100 {
		return summary, nil
	}

	for page := 2; ; page++ {
		nextPageCommits, _, err := c.listRepoCommits(owner, repo, ListRepoCommitsOptions{
			SHA:     ref,
			Page:    page,
			PerPage: 100,
		})
		if err != nil {
			return nil, err
		}

		summary.Count += len(nextPageCommits)
		if len(nextPageCommits) < 100 {
			return summary, nil
		}
	}
}

// ListBranches lists repository branches.
func (c *Client) ListBranches(owner, repo string, opts ListBranchesOptions) ([]Branch, error) {
	var branches []Branch
	query := buildListBranchesQuery(opts)
	err := c.doFromEndpointWithQuery(Branches.List, []interface{}{owner, repo}, query, nil, &branches)
	if err != nil {
		return nil, err
	}
	return branches, nil
}

// CreateBranch creates a repository branch.
func (c *Client) CreateBranch(owner, repo string, opts CreateBranchOptions) (*CompleteBranch, error) {
	var branch CompleteBranch
	err := c.DoFromEndpoint(Branches.Create, []interface{}{owner, repo}, opts, &branch)
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

// GetBranch gets a repository branch.
func (c *Client) GetBranch(owner, repo, branchName string) (*CompleteBranch, error) {
	var branch CompleteBranch
	err := c.DoFromEndpoint(Branches.Get, []interface{}{owner, repo, branchName}, nil, &branch)
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

// ProtectBranch enables branch protection.
func (c *Client) ProtectBranch(owner, repo, branchName string) (*CompleteBranch, error) {
	var branch CompleteBranch
	err := c.DoFromEndpoint(Branches.Protection, []interface{}{owner, repo, branchName}, nil, &branch)
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

// UnprotectBranch disables branch protection.
func (c *Client) UnprotectBranch(owner, repo, branchName string) error {
	return c.DoFromEndpoint(Branches.Delete, []interface{}{owner, repo, branchName}, nil, nil)
}

// ListCollaborators lists repository collaborators.
func (c *Client) ListCollaborators(owner, repo string) ([]Collaborator, error) {
	var collabs []Collaborator
	err := c.DoFromEndpoint(Collaborators.List, []interface{}{owner, repo}, nil, &collabs)
	if err != nil {
		return nil, err
	}
	return collabs, nil
}

// GetCollaborator checks if a user is a collaborator.
func (c *Client) GetCollaborator(owner, repo, username string) (*Collaborator, error) {
	var collab Collaborator
	err := c.DoFromEndpoint(Collaborators.Get, []interface{}{owner, repo, username}, nil, &collab)
	if err != nil {
		return nil, err
	}
	return &collab, nil
}

// GetCollaboratorPermission gets a collaborator's permission level.
func (c *Client) GetCollaboratorPermission(owner, repo, username string) (*CollaboratorPermission, error) {
	var perm CollaboratorPermission
	err := c.DoFromEndpoint(Collaborators.Protection, []interface{}{owner, repo, username}, nil, &perm)
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

// AddCollaborator adds or updates a collaborator on a repository.
func (c *Client) AddCollaborator(owner, repo, username, permission string) error {
	req := map[string]string{"permission": permission}
	return c.DoFromEndpoint(Collaborators.Add, []interface{}{owner, repo, username}, req, nil)
}

// RemoveCollaborator removes a collaborator from a repository.
func (c *Client) RemoveCollaborator(owner, repo, username string) error {
	return c.DoFromEndpoint(Collaborators.Remove, []interface{}{owner, repo, username}, nil, nil)
}

// ListForks lists repository forks.
func (c *Client) ListForks(owner, repo string, opts ListForksOptions) ([]ForkRepository, error) {
	var forks []ForkRepository
	query := buildListForksQuery(opts)
	err := c.doFromEndpointWithQuery(RepoForks.List, []interface{}{owner, repo}, query, nil, &forks)
	if err != nil {
		return nil, err
	}
	return forks, nil
}

// ForkRepository forks a repository.
func (c *Client) ForkRepository(owner, repo string) (*ForkRepository, error) {
	var fork ForkRepository
	err := c.DoFromEndpoint(RepoForks.Fork, []interface{}{owner, repo}, nil, &fork)
	if err != nil {
		return nil, err
	}
	return &fork, nil
}

func buildListForksQuery(opts ListForksOptions) string {
	params := []string{"sort", opts.Sort}
	params = append(params, paginationParams(opts.Page, opts.PerPage)...)
	query := util.BuildQuery(params...)
	if query == "" {
		return ""
	}
	return "?" + query
}

func (c *Client) listRepoCommits(owner, repo string, opts ListRepoCommitsOptions) ([]BranchCommit, http.Header, error) {
	var commits []BranchCommit
	headers, err := c.DoFromEndpointWithHeaders(Repo.Commits, []interface{}{owner, repo}, opts, &commits)
	if err != nil {
		return nil, nil, err
	}
	return commits, headers, nil
}

func parseHeaderInt(headers http.Header, key string) int {
	if headers == nil {
		return 0
	}

	value := headers.Get(key)
	if value == "" {
		return 0
	}

	count, err := strconv.Atoi(value)
	if err != nil || count < 0 {
		return 0
	}

	return count
}

func buildListBranchesQuery(opts ListBranchesOptions) string {
	params := []string{"sort", opts.Sort, "direction", opts.Direction}
	params = append(params, paginationParams(opts.Page, opts.PerPage)...)
	query := util.BuildQuery(params...)
	if query == "" {
		return ""
	}
	return "?" + query
}
