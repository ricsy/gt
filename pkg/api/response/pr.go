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
	ID                int64  `json:"id"`
	Number            int    `json:"number"`
	State             string `json:"state"`
	Title             string `json:"title"`
	Body              string `json:"body"`
	BodyHTML          string `json:"body_html,omitempty"`
	HTMLURL           string `json:"html_url"`
	URL               string `json:"url"`
	DiffURL           string `json:"diff_url"`
	PatchURL          string `json:"patch_url"`
	IssueURL          string `json:"issue_url"`
	CommitsURL        string `json:"commits_url"`
	ReviewCommentsURL string `json:"review_comments_url"`
	ReviewCommentURL  string `json:"review_comment_url"`
	CommentsURL       string `json:"comments_url"`
	User              struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
	} `json:"user"`
	Labels             []Label     `json:"labels"`
	Assignee           *UserBasic  `json:"assignee"`
	Assignees          []UserBasic `json:"assignees"`
	Testers            []UserBasic `json:"testers"`
	CloseRelatedIssue  int         `json:"close_related_issue"`
	PruneBranch        bool        `json:"prune_branch"`
	AssigneesNumber    int         `json:"assignees_number"`
	TestersNumber      int         `json:"testers_number"`
	APIReviewersNumber int         `json:"api_reviewers_number"`
	APIReviewers       []UserBasic `json:"api_reviewers"`
	Milestone          *Milestone  `json:"milestone"`
	Locked             bool        `json:"locked"`
	Comments           int         `json:"comments"`
	Commits            int         `json:"commits"`
	Additions          int         `json:"additions"`
	Deletions          int         `json:"deletions"`
	FilesChanged       int         `json:"files_changed"`
	CreatedAt          string      `json:"created_at"`
	UpdatedAt          string      `json:"updated_at"`
	ClosedAt           string      `json:"closed_at"`
	MergedAt           string      `json:"merged_at"`
	Merged             bool        `json:"merged"`
	Mergeable          bool        `json:"mergeable"`
	CanMergeCheck      bool        `json:"can_merge_check,omitempty"`
	Head               struct {
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

// ListPRsOptions contains the optional parameters for ListPRs
type ListPRsOptions struct {
	State           string
	Head            string
	Base            string
	Sort            string // created, updated, popularity, long-running
	Since           string
	Direction       string // asc, desc
	MilestoneNumber int
	Labels          string
	Page            int
	PerPage         int
	Author          string
	Assignee        string
	Tester          string
}

// PullRequestComment represents a comment on a pull request
type PullRequestComment struct {
	ID          int64     `json:"id"`
	Body        string    `json:"body"`
	BodyHTML    string    `json:"body_html"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	User        UserBasic `json:"user"`
	Position    int       `json:"position,omitempty"`
	Line        int       `json:"line,omitempty"`
	TreeID      string    `json:"tree_id,omitempty"`
	InReplyToID int64     `json:"in_reply_to_id,omitempty"`
}

// PullRequestCommit represents a commit in a pull request
type PullRequestCommit struct {
	Sha    string     `json:"sha"`
	Commit CommitInfo `json:"commit"`
}

// CommitInfo contains commit details
type CommitInfo struct {
	Message string       `json:"message"`
	Author  CommitAuthor `json:"author"`
}

// CommitAuthor represents a commit author
type CommitAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

// PullRequestFile represents a file changed in a pull request
type PullRequestFile struct {
	Sha       string `json:"sha"`
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
	Patch     string `json:"patch,omitempty"`
}
