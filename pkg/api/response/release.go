package response

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

// CreateReleaseOptions contains options for creating a release
type CreateReleaseOptions struct {
	TagName         string `json:"tag_name"`
	Name            string `json:"name"`
	Body            string `json:"body"`
	TargetCommitish string `json:"target_commitish,omitempty"`
	Prerelease      bool   `json:"prerelease,omitempty"`
}

// UpdateReleaseOptions contains options for updating a release
type UpdateReleaseOptions struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	Prerelease *bool  `json:"prerelease,omitempty"`
}
