package response

type ProjectBasic struct {
	ID        int64  `json:"id"`
	FullName  string `json:"full_name"`
	HumanName string `json:"human_name"`
	URL       string `json:"url"`
}

type ProjectPermission struct {
	Admin bool `json:"admin"`
	Pull  bool `json:"pull"`
	Push  bool `json:"push"`
}

type Project struct {
	ID                  int64              `json:"id"`
	FullName            string             `json:"full_name"`
	HumanName           string             `json:"human_name"`
	URL                 string             `json:"url"`
	Path                string             `json:"path"`
	Name                string             `json:"name"`
	Owner               UserBasic          `json:"owner"`
	Description         string             `json:"description"`
	Private             bool               `json:"private"`
	Public              bool               `json:"public"`
	Internal            bool               `json:"internal"`
	Fork                bool               `json:"fork"`
	HTMLURL             string             `json:"html_url"`
	SSHURL              string             `json:"ssh_url"`
	ForksURL            string             `json:"forks_url"`
	KeysURL             string             `json:"keys_url"`
	CollaboratorsURL    string             `json:"collaborators_url"`
	HooksURL            string             `json:"hooks_url"`
	BranchesURL         string             `json:"branches_url"`
	TagsURL             string             `json:"tags_url"`
	BlobsURL            string             `json:"blobs_url"`
	StargazersURL       string             `json:"stargazers_url"`
	ContributorsURL     string             `json:"contributors_url"`
	CommitsURL          string             `json:"commits_url"`
	CommentsURL         string             `json:"comments_url"`
	IssueCommentURL     string             `json:"issue_comment_url"`
	IssuesURL           string             `json:"issues_url"`
	PullsURL            string             `json:"pulls_url"`
	MilestonesURL       string             `json:"milestones_url"`
	NotificationsURL    string             `json:"notifications_url"`
	LabelsURL           string             `json:"labels_url"`
	ReleasesURL         string             `json:"releases_url"`
	Recommend           bool               `json:"recommend"`
	GVP                 bool               `json:"gvp"`
	Homepage            string             `json:"homepage"`
	Language            string             `json:"language"`
	ForksCount          int                `json:"forks_count"`
	StargazersCount     int                `json:"stargazers_count"`
	WatchersCount       int                `json:"watchers_count"`
	DefaultBranch       string             `json:"default_branch"`
	OpenIssuesCount     int                `json:"open_issues_count"`
	HasIssues           bool               `json:"has_issues"`
	HasWiki             bool               `json:"has_wiki"`
	IssueComment        bool               `json:"issue_comment"`
	CanComment          bool               `json:"can_comment"`
	PullRequestsEnabled bool               `json:"pull_requests_enabled"`
	HasPage             bool               `json:"has_page"`
	License             string             `json:"license"`
	Outsourced          bool               `json:"outsourced"`
	ProjectCreator      string             `json:"project_creator"`
	Members             []string           `json:"members"`
	PushedAt            string             `json:"pushed_at"`
	CreatedAt           string             `json:"created_at"`
	UpdatedAt           string             `json:"updated_at"`
	Paas                string             `json:"paas"`
	Stared              bool               `json:"stared"`
	Watched             bool               `json:"watched"`
	Permission          *ProjectPermission `json:"permission,omitempty"`
	Relation            string             `json:"relation"`
	AssigneesNumber     int                `json:"assignees_number"`
	TestersNumber       int                `json:"testers_number"`
	Assignee            []UserBasic        `json:"assignee"`
	Testers             []UserBasic        `json:"testers"`
	Status              string             `json:"status"`
	Programs            []map[string]any   `json:"programs"`
	ProjectLabels       []ProjectLabel     `json:"project_labels"`
	IssueTemplateSource string             `json:"issue_template_source"`
}
