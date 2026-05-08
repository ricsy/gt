package api

import "github.com/ricsy/gt/pkg/api/response"
import "github.com/ricsy/gt/pkg/util"

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

// Collaborator is an alias for response.Collaborator
type Collaborator = response.Collaborator

// CollaboratorPermission is an alias for response.CollaboratorPermission
type CollaboratorPermission = response.CollaboratorPermission

// ForkRepository is an alias for response.ForkRepository
type ForkRepository = response.ForkRepository

// ListForksOptions is an alias for response.ListForksOptions
type ListForksOptions = response.ListForksOptions

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
	return util.BuildQuery("sort", opts.Sort)
}

func buildListBranchesQuery(opts ListBranchesOptions) string {
	params := []string{"sort", opts.Sort, "direction", opts.Direction}
	params = append(params, paginationParams(opts.Page, opts.PerPage)...)
	return util.BuildQuery(params...)
}
