package api

import "fmt"

// Release represents a Gitee release
type Release struct {
	ID           int64  `json:"id"`
	TagName      string `json:"tag_name"`
	Name         string `json:"name"`
	Body         string `json:"body"`
	TargetCommit string `json:"target_commitish"`
	CreatedAt    string `json:"created_at"`
	PublishedAt  string `json:"published_at"`
	HtmlUrl      string `json:"html_url"`
}

// ListReleases lists releases for a repository
func (c *Client) ListReleases(owner, repo string) ([]Release, error) {
	var releases []Release
	path := fmt.Sprintf(apiPathReleases, owner, repo)
	err := c.Do("GET", path, nil, &releases)
	if err != nil {
		return nil, err
	}
	return releases, nil
}

// GetRelease gets a release by tag
func (c *Client) GetRelease(owner, repo, tag string) (*Release, error) {
	var release Release
	path := fmt.Sprintf(apiPathReleases+"/tags/%s", owner, repo, tag)
	err := c.Do("GET", path, nil, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

// CreateReleaseOptions contains options for creating a release
type CreateReleaseOptions struct {
	TagName         string `json:"tag_name"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	TargetCommitish string `json:"target_commitish,omitempty"`
	Prerelease      bool   `json:"prerelease,omitempty"`
}

// CreateRelease creates a new release
func (c *Client) CreateRelease(owner, repo string, opts CreateReleaseOptions) (*Release, error) {
	var release Release
	path := fmt.Sprintf(apiPathReleases, owner, repo)
	err := c.Do("POST", path, opts, &release)
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
	path := fmt.Sprintf(apiPathReleases+"/%d", owner, repo, release.ID)
	return c.Do("DELETE", path, nil, nil)
}

// UpdateReleaseOptions contains options for updating a release
type UpdateReleaseOptions struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	Prerelease *bool  `json:"prerelease,omitempty"`
}

// UpdateRelease updates a release
func (c *Client) UpdateRelease(owner, repo string, releaseID int64, opts UpdateReleaseOptions) (*Release, error) {
	var release Release
	path := fmt.Sprintf(apiPathReleases+"/%d", owner, repo, releaseID)
	err := c.Do("PATCH", path, opts, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}
