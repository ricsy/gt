package response

// Org represents a Gitee organization (Group schema)
type Org struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	Blog        string `json:"blog"`
	Email       string `json:"email"`
	HtmlUrl     string `json:"html_url"`
	Location    string `json:"location"`
	AvatarUrl   string `json:"avatar_url"`
	ReposUrl    string `json:"repos_url"`
	EventsUrl   string `json:"events_url"`
	MembersUrl  string `json:"members_url"`
	Description string `json:"description"`
	FollowCount int    `json:"follow_count"`
}

// OrgDetail represents detailed org info (GroupDetail schema)
type OrgDetail struct {
	Org
	CreatedAt    string `json:"created_at"`
	Public       bool   `json:"public"`
	Enterprise   string `json:"enterprise,omitempty"`
	Members      int    `json:"members"`
	PublicRepos  int    `json:"public_repos"`
	PrivateRepos int    `json:"private_repos"`
	Owner        string `json:"owner"`
}

// OrgMember represents an organization member.
type OrgMember struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
	URL       string `json:"url"`
}

// ListOrgMembersOptions contains options for listing org members.
type ListOrgMembersOptions struct {
	Role    string
	Page    int
	PerPage int
}

// ListOrgsOptions contains options for listing orgs.
type ListOrgsOptions struct {
	Admin   bool
	Page    int
	PerPage int
}

// ListOrgReposOptions contains options for listing org repos.
type ListOrgReposOptions struct {
	Type string
	Page int
}
