package api

import (
	"strconv"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// EnterpriseBasic is an alias for response.EnterpriseBasic.
type EnterpriseBasic = response.EnterpriseBasic

// EnterpriseMember is an alias for response.EnterpriseMember.
type EnterpriseMember = response.EnterpriseMember

// ListEnterprisesOptions is an alias for response.ListEnterprisesOptions.
type ListEnterprisesOptions = response.ListEnterprisesOptions

// ListEnterpriseMembersOptions is an alias for response.ListEnterpriseMembersOptions.
type ListEnterpriseMembersOptions = response.ListEnterpriseMembersOptions

// SearchEnterpriseMemberOptions is an alias for response.SearchEnterpriseMemberOptions.
type SearchEnterpriseMemberOptions = response.SearchEnterpriseMemberOptions

// AddEnterpriseMemberOptions is an alias for response.AddEnterpriseMemberOptions.
type AddEnterpriseMemberOptions = response.AddEnterpriseMemberOptions

// UpdateEnterpriseMemberOptions is an alias for response.UpdateEnterpriseMemberOptions.
type UpdateEnterpriseMemberOptions = response.UpdateEnterpriseMemberOptions

// ListEnterpriseReposOptions is an alias for response.ListEnterpriseReposOptions.
type ListEnterpriseReposOptions = response.ListEnterpriseReposOptions

// ListEnterprisePullRequestsOptions is an alias for response.ListEnterprisePullRequestsOptions.
type ListEnterprisePullRequestsOptions = response.ListEnterprisePullRequestsOptions

// ListEnterprises lists enterprises for the authenticated user.
func (c *Client) ListEnterprises(opts ListEnterprisesOptions) ([]EnterpriseBasic, error) {
	var enterprises []EnterpriseBasic
	err := c.doFromEndpointWithQuery(Enterprises.List, nil, buildListEnterprisesQuery(opts), nil, &enterprises)
	if err != nil {
		return nil, err
	}
	return enterprises, nil
}

// GetEnterprise gets a single enterprise.
func (c *Client) GetEnterprise(enterprise string) (*EnterpriseBasic, error) {
	var result EnterpriseBasic
	err := c.DoFromEndpoint(Enterprises.Get, []interface{}{enterprise}, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListEnterpriseMembers lists members in an enterprise.
func (c *Client) ListEnterpriseMembers(enterprise string, opts ListEnterpriseMembersOptions) ([]EnterpriseMember, error) {
	var members []EnterpriseMember
	err := c.doFromEndpointWithQuery(Enterprises.Members, []interface{}{enterprise}, buildEnterpriseMembersQuery(opts), nil, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// AddEnterpriseMember invites or adds a member to an enterprise.
func (c *Client) AddEnterpriseMember(enterprise string, opts AddEnterpriseMemberOptions) error {
	return c.DoFromEndpoint(Enterprises.Create, []interface{}{enterprise}, opts, nil)
}

// SearchEnterpriseMember searches an enterprise member by username or email.
func (c *Client) SearchEnterpriseMember(enterprise string, opts SearchEnterpriseMemberOptions) (*EnterpriseMember, error) {
	var member EnterpriseMember
	query := util.BuildQuery("query_type", opts.QueryType, "query_value", opts.QueryValue)
	err := c.doFromEndpointWithQuery(Enterprises.SearchMembers, []interface{}{enterprise}, query, nil, &member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetEnterpriseMember gets a member in an enterprise.
func (c *Client) GetEnterpriseMember(enterprise, username string) (*EnterpriseMember, error) {
	var member EnterpriseMember
	err := c.DoFromEndpoint(Enterprises.Member, []interface{}{enterprise, username}, nil, &member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// UpdateEnterpriseMember updates an enterprise member.
func (c *Client) UpdateEnterpriseMember(enterprise, username string, opts UpdateEnterpriseMemberOptions) (*EnterpriseMember, error) {
	var member EnterpriseMember
	err := c.DoFromEndpoint(Enterprises.Update, []interface{}{enterprise, username}, opts, &member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// RemoveEnterpriseMember removes a member from an enterprise.
func (c *Client) RemoveEnterpriseMember(enterprise, username string) error {
	return c.DoFromEndpoint(Enterprises.Delete, []interface{}{enterprise, username}, nil, nil)
}

// ListEnterpriseRepos lists repositories in an enterprise.
func (c *Client) ListEnterpriseRepos(enterprise string, opts ListEnterpriseReposOptions) ([]Repository, error) {
	var repos []Repository
	err := c.doFromEndpointWithQuery(Enterprises.Repos, []interface{}{enterprise}, buildEnterpriseReposQuery(opts), nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// ListEnterprisePullRequests lists pull requests in an enterprise.
func (c *Client) ListEnterprisePullRequests(enterprise string, opts ListEnterprisePullRequestsOptions) ([]PullRequest, error) {
	var prs []PullRequest
	err := c.doFromEndpointWithQuery(Enterprises.PullRequests, []interface{}{enterprise}, buildEnterprisePullRequestsQuery(opts), nil, &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func buildListEnterprisesQuery(opts ListEnterprisesOptions) string {
	params := paginationParams(opts.Page, opts.PerPage)
	if opts.Admin != nil {
		params = append(params, "admin", strconv.FormatBool(*opts.Admin))
	}
	return util.BuildQuery(params...)
}

func buildEnterpriseMembersQuery(opts ListEnterpriseMembersOptions) string {
	params := []string{"role", opts.Role}
	params = append(params, paginationParams(opts.Page, opts.PerPage)...)
	return util.BuildQuery(params...)
}

func buildEnterpriseReposQuery(opts ListEnterpriseReposOptions) string {
	params := []string{"search", opts.Search, "type", opts.Type}
	if opts.Direct != nil {
		params = append(params, "direct", strconv.FormatBool(*opts.Direct))
	}
	params = append(params, paginationParams(opts.Page, opts.PerPage)...)
	return util.BuildQuery(params...)
}

func buildEnterprisePullRequestsQuery(opts ListEnterprisePullRequestsOptions) string {
	params := []string{
		"issue_number", opts.IssueNumber,
		"repo", opts.Repo,
		"state", opts.State,
		"head", opts.Head,
		"base", opts.Base,
		"sort", opts.Sort,
		"since", opts.Since,
		"direction", opts.Direction,
		"labels", opts.Labels,
	}
	if opts.ProgramID > 0 {
		params = append(params, "program_id", strconv.Itoa(opts.ProgramID))
	}
	if opts.MilestoneNumber > 0 {
		params = append(params, "milestone_number", strconv.Itoa(opts.MilestoneNumber))
	}
	params = append(params, paginationParams(opts.Page, opts.PerPage)...)
	return util.BuildQuery(params...)
}
