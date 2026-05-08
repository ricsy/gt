package api

import "github.com/ricsy/gt/pkg/api/response"

// Release is an alias for response.Release
type Release = response.Release

// CreateReleaseOptions is an alias for response.CreateReleaseOptions
type CreateReleaseOptions = response.CreateReleaseOptions

// UpdateReleaseOptions is an alias for response.UpdateReleaseOptions
type UpdateReleaseOptions = response.UpdateReleaseOptions

// ListReleases lists releases for a repository
func (c *Client) ListReleases(owner, repo string) ([]Release, error) {
	var releases []Release
	err := c.DoFromEndpoint(Releases.List, []interface{}{owner, repo}, nil, &releases)
	if err != nil {
		return nil, err
	}
	return releases, nil
}

// GetRelease gets a release by tag
func (c *Client) GetRelease(owner, repo, tag string) (*Release, error) {
	var release Release
	err := c.DoFromEndpoint(Releases.Get, []interface{}{owner, repo, tag}, nil, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// GetReleaseByID gets a release by ID
func (c *Client) GetReleaseByID(owner, repo string, id int64) (*Release, error) {
	var release Release
	err := c.DoFromEndpoint(Releases.GetByID, []interface{}{owner, repo, id}, nil, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// CreateRelease creates a new release
func (c *Client) CreateRelease(owner, repo string, opts CreateReleaseOptions) (*Release, error) {
	var release Release
	err := c.DoFromEndpoint(Releases.Create, []interface{}{owner, repo}, opts, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// DeleteRelease deletes a release by tag
func (c *Client) DeleteRelease(owner, repo, tag string) error {
	release, err := c.GetRelease(owner, repo, tag)
	if err != nil {
		return err
	}
	return c.DoFromEndpoint(Releases.Delete, []interface{}{owner, repo, release.ID}, nil, nil)
}

// UpdateRelease updates a release
func (c *Client) UpdateRelease(owner, repo string, releaseID int64, opts UpdateReleaseOptions) (*Release, error) {
	var release Release
	err := c.DoFromEndpoint(Releases.Update, []interface{}{owner, repo, releaseID}, opts, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// GetLatestRelease gets the latest release for a repository
func (c *Client) GetLatestRelease(owner, repo string) (*Release, error) {
	var release Release
	err := c.DoFromEndpoint(Releases.Latest, []interface{}{owner, repo}, nil, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}
