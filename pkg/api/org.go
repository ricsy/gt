package api

import "github.com/ricsy/gt/pkg/api/response"

// Org is an alias for response.Org
type Org = response.Org

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
