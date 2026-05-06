package response

// PRState represents the state of a pull request.
type PRState string

const (
	PRStateOpen   PRState = "open"
	PRStateClosed PRState = "closed"
	PRStateMerged PRState = "merged"
	PRStateAll    PRState = "all"
)

// PullRequest represents a Gitee pull request
type PullRequest struct {
	ID      int64  `json:"id"`
	Number  int    `json:"number"`
	State   string `json:"state"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	User    struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
	} `json:"user"`
	Labels       []Label `json:"labels"`
	Assignee     *User   `json:"assignee"`
	Assignees    []User  `json:"assignees"`
	Comments     int     `json:"comments"`
	Commits      int     `json:"commits"`
	Additions    int     `json:"additions"`
	Deletions    int     `json:"deletions"`
	FilesChanged int     `json:"files_changed"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	ClosedAt     string  `json:"closed_at"`
	MergedAt     string  `json:"merged_at"`
	Merged       bool    `json:"merged"`
	Mergeable    bool    `json:"mergeable"`
	Head         struct {
		Ref string `json:"ref"`
		Sha string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
		Sha string `json:"sha"`
	} `json:"base"`
}

// CreatePRRequest is the request body for creating a PR
type CreatePRRequest struct {
	Title                 string `json:"title"`
	Body                  string `json:"body,omitempty"`
	Head                  string `json:"head"`
	Base                  string `json:"base"`
	MilestoneNumber       int    `json:"milestone_number,omitempty"`
	Labels                string `json:"labels,omitempty"`
	Issue                 string `json:"issue,omitempty"`
	Assignees             string `json:"assignees,omitempty"`
	Testers               string `json:"testers,omitempty"`
	AssigneesNumber       int    `json:"assignees_number,omitempty"`
	TestersNumber         int    `json:"testers_number,omitempty"`
	RefPullRequestNumbers string `json:"ref_pull_request_numbers,omitempty"`
	PruneSourceBranch     bool   `json:"prune_source_branch,omitempty"`
	CloseRelatedIssue     bool   `json:"close_related_issue,omitempty"`
	Draft                 bool   `json:"draft,omitempty"`
	Squash                bool   `json:"squash,omitempty"`
	SecurityHole          bool   `json:"security_hole,omitempty"`
}

// MergePRRequest is the request body for merging a PR
type MergePRRequest struct {
	MergeMethod       string `json:"merge_method,omitempty"`
	PruneSourceBranch bool   `json:"prune_source_branch,omitempty"`
	CloseRelatedIssue bool   `json:"close_related_issue,omitempty"`
	Title             string `json:"title,omitempty"`
	Description       string `json:"description,omitempty"`
}

// UpdatePRRequest is the request body for updating a PR
type UpdatePRRequest struct {
	Title                 string `json:"title,omitempty"`
	Body                  string `json:"body,omitempty"`
	State                 string `json:"state,omitempty"`
	MilestoneNumber       int    `json:"milestone_number,omitempty"`
	Labels                string `json:"labels,omitempty"`
	AssigneesNumber       int    `json:"assignees_number,omitempty"`
	TestersNumber         int    `json:"testers_number,omitempty"`
	RefPullRequestNumbers string `json:"ref_pull_request_numbers,omitempty"`
	CloseRelatedIssue     bool   `json:"close_related_issue,omitempty"`
	Draft                 bool   `json:"draft,omitempty"`
	Squash                bool   `json:"squash,omitempty"`
	SecurityHole          bool   `json:"security_hole,omitempty"`
}
