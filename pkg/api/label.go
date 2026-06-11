package api

import (
	"github.com/ricsy/gt/pkg/api/response"
)

// CreateLabelOptions is an alias for response.CreateLabelOptions
type CreateLabelOptions = response.CreateLabelOptions

// UpdateLabelOptions is an alias for response.UpdateLabelOptions
type UpdateLabelOptions = response.UpdateLabelOptions

// ListLabels lists labels for a repo
func (c *Client) ListLabels(owner, repo string) ([]response.Label, error) {
	var labels []response.Label
	err := c.DoFromEndpoint(Labels.List, []interface{}{owner, repo}, nil, &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// GetLabel gets a single label
func (c *Client) GetLabel(owner, repo, name string) (*response.Label, error) {
	var label response.Label
	err := c.DoFromEndpoint(Labels.Get, []interface{}{owner, repo, name}, nil, &label)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

// CreateLabel creates a new label
func (c *Client) CreateLabel(owner, repo string, opts response.CreateLabelOptions) (*response.Label, error) {
	var label response.Label
	err := c.DoFromEndpoint(Labels.Create, []interface{}{owner, repo}, opts, &label)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

// UpdateLabel updates a label
func (c *Client) UpdateLabel(owner, repo, originalName string, opts response.UpdateLabelOptions) (*response.Label, error) {
	var label response.Label
	err := c.DoFromEndpoint(Labels.Update, []interface{}{owner, repo, originalName}, opts, &label)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

// DeleteLabel deletes a label
func (c *Client) DeleteLabel(owner, repo, name string) error {
	return c.DoFromEndpoint(Labels.Delete, []interface{}{owner, repo, name}, nil, nil)
}

// ListIssueLabels lists labels for an issue
func (c *Client) ListIssueLabels(owner, repo, number string) ([]response.Label, error) {
	var labels []response.Label
	err := c.DoFromEndpoint(IssueLabels.List, []interface{}{owner, repo, number}, nil, &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// AddIssueLabels adds labels to an issue
func (c *Client) AddIssueLabels(owner, repo, number string, names []string) ([]response.Label, error) {
	var labels []response.Label
	err := c.DoFromEndpoint(IssueLabels.Create, []interface{}{owner, repo, number}, names, &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// ReplaceIssueLabels replaces all labels for an issue
func (c *Client) ReplaceIssueLabels(owner, repo, number string, names []string) ([]response.Label, error) {
	var labels []response.Label
	err := c.DoFromEndpoint(IssueLabels.Replace, []interface{}{owner, repo, number}, names, &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// DeleteAllIssueLabels deletes all labels from an issue
func (c *Client) DeleteAllIssueLabels(owner, repo, number string) error {
	return c.DoFromEndpoint(IssueLabels.Delete, []interface{}{owner, repo, number}, nil, nil)
}

// DeleteIssueLabel deletes a specific label from an issue
func (c *Client) DeleteIssueLabel(owner, repo, number, name string) error {
	return c.DoFromEndpoint(IssueLabels.Remove, []interface{}{owner, repo, number, name}, nil, nil)
}

// ListProjectLabels lists project labels for a repo
func (c *Client) ListProjectLabels(owner, repo string) ([]response.ProjectLabel, error) {
	var labels []response.ProjectLabel
	err := c.DoFromEndpoint(ProjectLabels.List, []interface{}{owner, repo}, nil, &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// AddProjectLabels adds labels to a project
func (c *Client) AddProjectLabels(owner, repo string, names []string) ([]response.ProjectLabel, error) {
	var labels []response.ProjectLabel
	err := c.DoFromEndpoint(ProjectLabels.Create, []interface{}{owner, repo}, names, &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// ReplaceProjectLabels replaces all project labels
func (c *Client) ReplaceProjectLabels(owner, repo string, names []string) ([]response.ProjectLabel, error) {
	var labels []response.ProjectLabel
	err := c.DoFromEndpoint(ProjectLabels.Replace, []interface{}{owner, repo}, names, &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// DeleteProjectLabels deletes project labels
func (c *Client) DeleteProjectLabels(owner, repo string, names []string) error {
	return c.DoFromEndpoint(ProjectLabels.Delete, []interface{}{owner, repo}, names, nil)
}

// ListEnterpriseLabels lists labels for an enterprise
func (c *Client) ListEnterpriseLabels(enterprise string) ([]response.Label, error) {
	var labels []response.Label
	err := c.DoFromEndpoint(EnterpriseLabels.List, []interface{}{enterprise}, nil, &labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

// GetEnterpriseLabel gets a single enterprise label
func (c *Client) GetEnterpriseLabel(enterprise, name string) (*response.Label, error) {
	var label response.Label
	err := c.DoFromEndpoint(EnterpriseLabels.Get, []interface{}{enterprise, name}, nil, &label)
	if err != nil {
		return nil, err
	}
	return &label, nil
}
