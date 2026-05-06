package response

// Project represents a Gitee project/repository
type Project struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Owner           UserBasic `json:"owner"`
	Description     string    `json:"description"`
	Private         bool      `json:"private"`
	Fork            bool      `json:"fork"`
	HTMLURL         string    `json:"html_url"`
	SSHURL          string    `json:"ssh_url"`
	CloneURL        string    `json:"clone_url"`
	DefaultBranch   string    `json:"default_branch"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	PushedAt        string    `json:"pushed_at"`
	Homepage        string    `json:"homepage"`
	Language        string    `json:"language"`
	StarCount       int       `json:"stargazers_count"`
	WatchCount      int       `json:"watchers_count"`
	ForksCount      int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
}
