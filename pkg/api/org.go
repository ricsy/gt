package api

import (
	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// Org is an alias for response.Org
type Org = response.Org

// OrgMember is an alias for response.OrgMember
type OrgMember = response.OrgMember

// ListOrgMembersOptions is an alias for response.ListOrgMembersOptions
type ListOrgMembersOptions = response.ListOrgMembersOptions

// ListOrgsOptions is an alias for response.ListOrgsOptions
type ListOrgsOptions = response.ListOrgsOptions

// ListOrgReposOptions is an alias for response.ListOrgReposOptions
type ListOrgReposOptions = response.ListOrgReposOptions

// ListOrgs lists organizations for the current user
func (c *Client) ListOrgs(opts ListOrgsOptions) ([]Org, error) {
	var orgs []Org
	err := c.doFromEndpointWithQuery(UserOrgs.List, nil, buildListOrgsQuery(opts), nil, &orgs)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

// GetUserOrgsByUsername gets organizations for a user by username
func (c *Client) GetUserOrgsByUsername(username string) ([]Org, error) {
	var orgs []Org
	err := c.DoFromEndpoint(Orgs.GetUserOrgsByUsername, []interface{}{username}, nil, &orgs)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

// GetOrg gets an organization by login
func (c *Client) GetOrg(login string) (*Org, error) {
	var org Org
	err := c.DoFromEndpoint(Orgs.Get, []interface{}{login}, nil, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// ListOrgMembers lists organization members.
func (c *Client) ListOrgMembers(org string, opts ListOrgMembersOptions) ([]OrgMember, error) {
	var members []OrgMember
	err := c.doFromEndpointWithQuery(Orgs.Members, []interface{}{org}, buildListOrgMembersQuery(opts), nil, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// ListOrgRepos lists organization repositories.
func (c *Client) ListOrgRepos(org string, opts ListOrgReposOptions) ([]Repository, error) {
	var repos []Repository
	err := c.doFromEndpointWithQuery(Orgs.Repos, []interface{}{org}, buildListOrgReposQuery(opts), nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func buildListOrgMembersQuery(opts ListOrgMembersOptions) string {
	params := []string{}
	if opts.Role != "" {
		params = append(params, "role", opts.Role)
	}
	params = append(params, paginationParams(opts.Page, opts.PerPage)...)
	return util.BuildQuery(params...)
}

func buildListOrgReposQuery(opts ListOrgReposOptions) string {
	params := []string{}
	if opts.Type != "" {
		params = append(params, "type", opts.Type)
	}
	params = append(params, paginationParams(opts.Page, 0)...)
	return util.BuildQuery(params...)
}

func buildListOrgsQuery(opts ListOrgsOptions) string {
	params := []string{}
	if opts.Admin {
		params = append(params, "admin", "true")
	}
	params = append(params, paginationParams(opts.Page, opts.PerPage)...)
	return util.BuildQuery(params...)
}
