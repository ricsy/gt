package response

// Repository represents a Gitee repository
type Repository struct {
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

// Branch represents a repository branch.
type Branch struct {
	Name          string `json:"name"`
	Commit        string `json:"commit"`
	Protected     bool   `json:"protected"`
	ProtectionURL string `json:"protection_url"`
}

// CompleteBranch represents a detailed repository branch.
type CompleteBranch struct {
	Name          string `json:"name"`
	Commit        string `json:"commit"`
	Links         string `json:"_links"`
	Protected     bool   `json:"protected"`
	ProtectionURL string `json:"protection_url"`
}

// ListBranchesOptions contains options for listing repository branches.
type ListBranchesOptions struct {
	Sort      string
	Direction string
	Page      int
	PerPage   int
}

// CreateBranchOptions contains options for creating a repository branch.
type CreateBranchOptions struct {
	Refs       string `json:"refs"`
	BranchName string `json:"branch_name"`
}

// CreateRepoOptions contains options for creating a repository
type CreateRepoOptions struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
}

// UpdateRepoOptions contains options for updating a repository
type UpdateRepoOptions struct {
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	Homepage             string `json:"homepage,omitempty"`
	HasIssues            *bool  `json:"has_issues,omitempty"`
	HasWiki              *bool  `json:"has_wiki,omitempty"`
	CanComment           *bool  `json:"can_comment,omitempty"`
	IssueComment         *bool  `json:"issue_comment,omitempty"`
	SecurityHoleEnabled  *bool  `json:"security_hole_enabled,omitempty"`
	Private              *bool  `json:"private,omitempty"`
	Path                 string `json:"path,omitempty"`
	DefaultBranch        string `json:"default_branch,omitempty"`
	PullRequestsEnabled  *bool  `json:"pull_requests_enabled,omitempty"`
	OnlineEditEnabled    *bool  `json:"online_edit_enabled,omitempty"`
	LightweightPrEnabled *bool  `json:"lightweight_pr_enabled,omitempty"`
	MergeEnabled         *bool  `json:"merge_enabled,omitempty"`
	SquashEnabled        *bool  `json:"squash_enabled,omitempty"`
	RebaseEnabled        *bool  `json:"rebase_enabled,omitempty"`
	DefaultMergeMethod   string `json:"default_merge_method,omitempty"`
	IssueTemplateSource  string `json:"issue_template_source,omitempty"`
}

// Collaborator represents a repository collaborator.
type Collaborator struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// CollaboratorPermission represents a collaborator's permission level.
type CollaboratorPermission struct {
	Permission string `json:"permission"`
}

// ForkRepository represents a forked repository.
type ForkRepository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
	} `json:"owner"`
	Fork        bool   `json:"fork"`
	HTMLURL     string `json:"html_url"`
	CloneURL    string `json:"clone_url"`
	Description string `json:"description"`
	StarCount   int    `json:"stargazers_count"`
	ForksCount  int    `json:"forks_count"`
}

// ListForksOptions contains options for listing forks.
type ListForksOptions struct {
	Sort    string `json:"sort,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}
