package response

// IssueState represents the state of an issue.
type IssueState string

const (
	IssueStateOpen        IssueState = "open"
	IssueStateClosed      IssueState = "closed"
	IssueStateProgressing IssueState = "progressing"
	IssueStateRejected    IssueState = "rejected"
	IssueStateAll         IssueState = "all"
)

// Issue represents a Gitee issue
type Issue struct {
	ID      int64  `json:"id"`
	Number  int64  `json:"number"`
	State   string `json:"state"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	User    struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
	} `json:"user"`
	Labels      []Label `json:"labels"`
	Assignee    *User   `json:"assignee"`
	Assignees   []User  `json:"assignees"`
	Comments    int     `json:"comments"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	ClosedAt    string  `json:"closed_at"`
	PullRequest *struct {
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
	} `json:"pull_request"`
}

// Label represents a Gitee label
type Label struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// User represents a Gitee user
type User struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// IssueComment represents a comment on an issue
type IssueComment struct {
	ID        int64  `json:"id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

// ListIssuesOptions contains the optional parameters for ListIssues
type ListIssuesOptions struct {
	State        string
	Labels       string
	Sort         string // created, updated
	Direction    string // asc, desc
	Milestone    string
	Assignee     string
	Creator      string
	Program      string
	Q            string
	SecurityHole *bool
	Page         int
	PerPage      int
}

// CreateIssueRequest is the request body for creating an issue
type CreateIssueRequest struct {
	Title         string `json:"title"`
	Body          string `json:"body,omitempty"`
	Repo          string `json:"repo,omitempty"`
	IssueType     string `json:"issue_type,omitempty"`
	Assignee      string `json:"assignee,omitempty"`
	Collaborators string `json:"collaborators,omitempty"`
	Milestone     int    `json:"milestone,omitempty"`
	Labels        string `json:"labels,omitempty"`
	Program       string `json:"program,omitempty"`
	SecurityHole  bool   `json:"security_hole,omitempty"`
	Branch        string `json:"branch,omitempty"`
}

// UpdateIssueRequest is the request body for updating an issue
type UpdateIssueRequest struct {
	Repo          string `json:"repo,omitempty"`
	Title         string `json:"title,omitempty"`
	Body          string `json:"body,omitempty"`
	State         string `json:"state,omitempty"`
	Assignee      string `json:"assignee,omitempty"`
	Collaborators string `json:"collaborators,omitempty"`
	Milestone     int    `json:"milestone,omitempty"`
	Labels        string `json:"labels,omitempty"`
	Program       string `json:"program,omitempty"`
	SecurityHole  bool   `json:"security_hole,omitempty"`
	Branch        string `json:"branch,omitempty"`
}
