package response

// UserBasic represents basic user info
type UserBasic struct {
	ID                int64  `json:"id"`
	Login             string `json:"login"`
	Name              string `json:"name"`
	AvatarURL         string `json:"avatar_url"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	Remark            string `json:"remark"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	MemberRole        string `json:"member_role"`
	Blog              string `json:"blog"`
	Weibo             string `json:"weibo"`
	Bio               string `json:"bio"`
	PublicRepos       int    `json:"public_repos"`
	PublicGists       int    `json:"public_gists"`
	Followers         int    `json:"followers"`
	Following         int    `json:"following"`
	Stared            int    `json:"stared"`
	Watched           int    `json:"watched"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	Email             string `json:"email"`
	Company           string `json:"company"`
	Profession        string `json:"profession"`
	Wechat            string `json:"wechat"`
	QQ                string `json:"qq"`
	Linkedin          string `json:"linkedin"`
}

// User is an alias for UserBasic
type User = UserBasic

// UserDetail represents the authenticated user's profile.
type UserDetail struct {
	UserBasic
	BotInfo *BotInfo `json:"bot_info,omitempty"`
}

// BotInfo represents bot metadata attached to a user.
type BotInfo struct {
	Name string `json:"name"`
}

// UpdateUserOptions contains options for updating the authenticated user.
type UpdateUserOptions struct {
	Name  *string `json:"name,omitempty"`
	Blog  *string `json:"blog,omitempty"`
	Weibo *string `json:"weibo,omitempty"`
	Bio   *string `json:"bio,omitempty"`
}

// ListUsersOptions contains pagination options for user list endpoints.
type ListUsersOptions struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
}

// ListNamespacesOptions contains options for listing namespaces.
type ListNamespacesOptions struct {
	Mode    string `json:"mode,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}

// SSHKey represents a user SSH public key.
type SSHKey struct {
	ID        int64  `json:"id"`
	Key       string `json:"key"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

// SSHKeyBasic represents basic SSH key info (from GET /users/{username}/keys)
type SSHKeyBasic struct {
	ID  int64  `json:"id"`
	Key string `json:"key"`
}

// UserInfo represents public user profile from GET /users/{username}
type UserInfo struct {
	UserBasic
}

// CreateSSHKeyOptions contains options for creating an SSH key.
type CreateSSHKeyOptions struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

// Namespace represents a Gitee namespace.
type Namespace struct {
	ID           int64          `json:"id"`
	Type         string         `json:"type"`
	Name         string         `json:"name"`
	Path         string         `json:"path"`
	EnterpriseID int64          `json:"enterprise_id"`
	HTMLURL      string         `json:"html_url"`
	Parent       *NamespaceMini `json:"parent,omitempty"`
}

// NamespaceMini represents a parent namespace.
type NamespaceMini struct {
	ID           int64  `json:"id"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	EnterpriseID int64  `json:"enterprise_id"`
	HTMLURL      string `json:"html_url"`
}
