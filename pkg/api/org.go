package api

import (
	"strconv"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// Org is an alias for response.Org
type Org = response.Org

// OrgMember is an alias for response.OrgMember
type OrgMember = response.OrgMember

// ListOrgMembersOptions is an alias for response.ListOrgMembersOptions
type ListOrgMembersOptions = response.ListOrgMembersOptions

// ListOrgReposOptions is an alias for response.ListOrgReposOptions
type ListOrgReposOptions = response.ListOrgReposOptions

// ListOrgs lists organizations for the current user
func (c *Client) ListOrgs() ([]Org, error) {
	var orgs []Org
	err := c.DoFromEndpoint(UserOrgs.List, nil, nil, &orgs)
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
	if opts.Page > 0 {
		params = append(params, "page", strconv.Itoa(opts.Page))
	}
	return util.BuildQuery(params...)
}

func buildListOrgReposQuery(opts ListOrgReposOptions) string {
	params := []string{}
	if opts.Type != "" {
		params = append(params, "type", opts.Type)
	}
	if opts.Page > 0 {
		params = append(params, "page", strconv.Itoa(opts.Page))
	}
	return util.BuildQuery(params...)
}
