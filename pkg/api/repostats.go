package api

import (
	"github.com/ricsy/gt/pkg/api/response"
)

// GetTrafficDataOptions is an alias for response.GetTrafficDataOptions
type GetTrafficDataOptions = response.GetTrafficDataOptions

// GetTrafficData gets repository traffic data
func (c *Client) GetTrafficData(owner, repo string, opts GetTrafficDataOptions) (*response.TrafficData, error) {
	var trafficData response.TrafficData
	err := c.DoFromEndpoint(RepoStats.Get, []interface{}{owner, repo}, opts, &trafficData)
	if err != nil {
		return nil, err
	}
	return &trafficData, nil
}

// GetRepoLanguages gets repository language statistics
func (c *Client) GetRepoLanguages(owner, repo string) (*response.Languages, error) {
	var languages response.Languages
	err := c.DoFromEndpoint(RepoLanguages.Get, []interface{}{owner, repo}, nil, &languages)
	if err != nil {
		return nil, err
	}
	return &languages, nil
}

// GetRepoContributors gets repository contributors
func (c *Client) GetRepoContributors(owner, repo string) ([]response.Contributor, error) {
	var contributors []response.Contributor
	err := c.DoFromEndpoint(RepoContributors.Get, []interface{}{owner, repo}, nil, &contributors)
	if err != nil {
		return nil, err
	}
	return contributors, nil
}
