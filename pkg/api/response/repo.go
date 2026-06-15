package response

// BranchCommitActor 表示分支提交元数据中的作者或提交者信息。
type BranchCommitActor struct {
	Name  string `json:"name"`
	Date  string `json:"date"`
	Email string `json:"email"`
}

// BranchTree 表示分支提交关联的树对象。
type BranchTree struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

// BranchCommitDetail 表示 branch 响应中的嵌套 commit 明细。
type BranchCommitDetail struct {
	Author    BranchCommitActor `json:"author"`
	Committer BranchCommitActor `json:"committer"`
	Message   string            `json:"message"`
	URL       string            `json:"url"`
	Tree      BranchTree        `json:"tree"`
}

// BranchUserRef 表示 branch 响应中附带的用户引用。
type BranchUserRef struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	URL       string `json:"url"`
	AvatarURL string `json:"avatar_url"`
}

// BranchCommit 表示分支接口返回的提交摘要。
// live API 返回的是对象，而不是旧实现假定的 SHA 字符串。
type BranchCommit struct {
	SHA       string             `json:"sha"`
	URL       string             `json:"url"`
	Commit    BranchCommitDetail `json:"commit"`
	Author    BranchUserRef      `json:"author"`
	Committer BranchUserRef      `json:"committer"`
	Parents   []any              `json:"parents"`
}

// BranchLinks 表示单分支详情返回的链接集合。
type BranchLinks struct {
	HTML string `json:"html"`
	Self string `json:"self"`
}

// Repository represents a Gitee repository
type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
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
	Name          string       `json:"name"`
	Commit        BranchCommit `json:"commit"`
	Protected     bool         `json:"protected"`
	ProtectionURL string       `json:"protection_url"`
}

// CompleteBranch represents a detailed repository branch.
type CompleteBranch struct {
	Name          string       `json:"name"`
	Commit        BranchCommit `json:"commit"`
	Links         BranchLinks  `json:"_links"`
	Protected     bool         `json:"protected"`
	ProtectionURL string       `json:"protection_url"`
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

type ListRepoCommitsOptions struct {
	SHA     string `json:"sha,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}

// CreateRepoOptions contains options for creating a repository
type CreateRepoOptions struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	Homepage           string `json:"homepage,omitempty"`
	Private            bool   `json:"private"`
	AutoInit           bool   `json:"auto_init"`
	HasIssues          *bool  `json:"has_issues,omitempty"`
	HasWiki            *bool  `json:"has_wiki,omitempty"`
	CanComment         *bool  `json:"can_comment,omitempty"`
	GitignoreTemplate  string `json:"gitignore_template,omitempty"`
	LicenseTemplate    string `json:"license_template,omitempty"`
	Path               string `json:"path,omitempty"`
	Namespace          string `json:"namespace,omitempty"`
	Public             *bool  `json:"public,omitempty"`
	Outsourced         *bool  `json:"outsourced,omitempty"`
	ProjectCreator     string `json:"project_creator,omitempty"`
	Members            string `json:"members,omitempty"`
	TemplateApplyScope int    `json:"template_apply_scope,omitempty"`
	Model              string `json:"model,omitempty"`
	Enterprise         string `json:"enterprise,omitempty"`
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
	UserBasic
	Remark       string `json:"remark"`
	FollowersURL string `json:"followers_url"`
	FollowingURL string `json:"following_url"`
	GistsURL     string `json:"gists_url"`
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
