package response

// EnterpriseBasic represents a Gitee enterprise.
type EnterpriseBasic struct {
	ID        int64  `json:"id"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	AvatarURL string `json:"avatar_url"`
}

// EnterpriseMember represents a member of an enterprise.
type EnterpriseMember struct {
	URL        string           `json:"url"`
	Active     bool             `json:"active"`
	Remark     string           `json:"remark"`
	Role       string           `json:"role"`
	Outsourced bool             `json:"outsourced"`
	Enterprise *EnterpriseBasic `json:"enterprise"`
	User       *UserBasic       `json:"user"`
}

// ListEnterprisesOptions contains options for listing authenticated user's enterprises.
type ListEnterprisesOptions struct {
	Page    int
	PerPage int
	Admin   *bool
}

// ListEnterpriseMembersOptions contains options for listing enterprise members.
type ListEnterpriseMembersOptions struct {
	Role    string
	Page    int
	PerPage int
}

// SearchEnterpriseMemberOptions contains options for searching an enterprise member.
type SearchEnterpriseMemberOptions struct {
	QueryType  string
	QueryValue string
}

// AddEnterpriseMemberOptions contains options for inviting an enterprise member.
type AddEnterpriseMemberOptions struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
	Name     string `json:"name,omitempty"`
}

// UpdateEnterpriseMemberOptions contains options for updating an enterprise member.
type UpdateEnterpriseMemberOptions struct {
	Role   string `json:"role,omitempty"`
	Active *bool  `json:"active,omitempty"`
	Name   string `json:"name,omitempty"`
}

// ListEnterpriseReposOptions contains options for listing enterprise repositories.
type ListEnterpriseReposOptions struct {
	Search  string
	Type    string
	Direct  *bool
	Page    int
	PerPage int
}

// ListEnterprisePullRequestsOptions contains options for listing enterprise pull requests.
type ListEnterprisePullRequestsOptions struct {
	IssueNumber     string
	Repo            string
	ProgramID       int
	State           string
	Head            string
	Base            string
	Sort            string
	Since           string
	Direction       string
	MilestoneNumber int
	Labels          string
	Page            int
	PerPage         int
}
