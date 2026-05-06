package api

import (
	"strconv"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// ListMilestonesOptions is an alias for response.ListMilestonesOptions
type ListMilestonesOptions = response.ListMilestonesOptions

// CreateMilestoneOptions is an alias for response.CreateMilestoneOptions
type CreateMilestoneOptions = response.CreateMilestoneOptions

// UpdateMilestoneOptions is an alias for response.UpdateMilestoneOptions
type UpdateMilestoneOptions = response.UpdateMilestoneOptions

// ListMilestones lists milestones for a repo
func (c *Client) ListMilestones(owner, repo string, opts ListMilestonesOptions) ([]response.Milestone, error) {
	var milestones []response.Milestone
	path := Milestones.List.Build(owner, repo)
	query := buildMilestonesQuery(opts)
	err := c.Do("GET", path+query, nil, &milestones)
	if err != nil {
		return nil, err
	}
	return milestones, nil
}

// GetMilestone gets a single milestone
func (c *Client) GetMilestone(owner, repo string, number int) (*response.Milestone, error) {
	var milestone response.Milestone
	err := c.DoFromEndpoint(Milestones.Get, []interface{}{owner, repo, number}, nil, &milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

// CreateMilestone creates a new milestone
func (c *Client) CreateMilestone(owner, repo string, opts CreateMilestoneOptions) (*response.Milestone, error) {
	var milestone response.Milestone
	err := c.DoFromEndpoint(Milestones.Create, []interface{}{owner, repo}, opts, &milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

// UpdateMilestone updates a milestone
func (c *Client) UpdateMilestone(owner, repo string, number int, opts UpdateMilestoneOptions) (*response.Milestone, error) {
	var milestone response.Milestone
	err := c.DoFromEndpoint(Milestones.Update, []interface{}{owner, repo, number}, opts, &milestone)
	if err != nil {
		return nil, err
	}
	return &milestone, nil
}

// DeleteMilestone deletes a milestone
func (c *Client) DeleteMilestone(owner, repo string, number int) error {
	return c.DoFromEndpoint(Milestones.Delete, []interface{}{owner, repo, number}, nil, nil)
}

func buildMilestonesQuery(opts ListMilestonesOptions) string {
	var params []string
	if opts.State != "" {
		params = append(params, "state", opts.State)
	}
	if opts.Sort != "" {
		params = append(params, "sort", opts.Sort)
	}
	if opts.Direction != "" {
		params = append(params, "direction", opts.Direction)
	}
	if opts.Page > 0 {
		params = append(params, "page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params = append(params, "per_page", strconv.Itoa(opts.PerPage))
	}
	return util.BuildQuery(params...)
}
