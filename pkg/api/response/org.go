package response

// Org represents a Gitee organization
type Org struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Blog     string `json:"blog"`
	Email    string `json:"email"`
	HtmlUrl  string `json:"html_url"`
	Location string `json:"location"`
}

// OrgMember represents an organization member.
type OrgMember struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ListOrgMembersOptions contains options for listing org members.
type ListOrgMembersOptions struct {
	Page int
}

// ListOrgReposOptions contains options for listing org repos.
type ListOrgReposOptions struct {
	Type string
	Page int
}
