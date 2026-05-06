package api

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
	err := c.DoFromEndpoint(Releases.Update, []interface{}{owner, repo, releaseID}, opts, &release)
	if err != nil {
		return nil, err
	}
	return &release, nil
}
