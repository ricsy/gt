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
	ID            int64  `json:"id"`
	Number        int64  `json:"number"`
	State         string `json:"state"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	BodyHTML      string `json:"body_html"`
	HTMLURL       string `json:"html_url"`
	URL           string `json:"url"`
	RepositoryURL string `json:"repository_url"`
	LabelsURL     string `json:"labels_url"`
	CommentsURL   string `json:"comments_url"`
	ParentURL     string `json:"parent_url"`
	ParentID      int64  `json:"parent_id"`
	Depth         int    `json:"depth"`
	User          struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
	} `json:"user"`
	Labels        []Label     `json:"labels"`
	Assignee      *UserBasic  `json:"assignee"`
	Assignees     []UserBasic `json:"assignees"`
	Collaborators []UserBasic `json:"collaborators"`
	Milestone     *Milestone  `json:"milestone"`
	Repository    *Project    `json:"repository"`
	Comments      int         `json:"comments"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	ClosedAt      string      `json:"closed_at"`
	PullRequest   *struct {
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

// UserBasic represents basic user info
type UserBasic struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// Milestone represents a Gitee milestone
type Milestone struct {
	ID           int64      `json:"id"`
	Number       int        `json:"number"`
	State        string     `json:"state"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Creator      *UserBasic `json:"creator"`
	OpenIssues   int        `json:"open_issues"`
	ClosedIssues int        `json:"closed_issues"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
	DueOn        string     `json:"due_on"`
}

// Project represents a Gitee project/repository
type Project struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
	} `json:"owner"`
	Description     string `json:"description"`
	Private         bool   `json:"private"`
	Fork            bool   `json:"fork"`
	HTMLURL         string `json:"html_url"`
	SSHURL          string `json:"ssh_url"`
	CloneURL        string `json:"clone_url"`
	DefaultBranch   string `json:"default_branch"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	PushedAt        string `json:"pushed_at"`
	Homepage        string `json:"homepage"`
	Language        string `json:"language"`
	StarCount       int    `json:"stargazers_count"`
	WatchCount      int    `json:"watchers_count"`
	ForksCount      int    `json:"forks_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
}

// IssueComment represents a comment on an issue
type IssueComment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	BodyHTML  string    `json:"body_html"`
	CreatedAt string    `json:"created_at"`
	User      UserBasic `json:"user"`
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
	Since        string // ISO 8601 format
	Schedule     string
	Deadline     string
	CreatedAt    string
	FinishedAt   string
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
	CVEID         string `json:"cve_id,omitempty"`
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
