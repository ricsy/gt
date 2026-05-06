package response

// SearchReposOptions contains optional parameters for SearchRepos
type SearchReposOptions struct {
	Q        string
	Owner    string
	Fork     *bool
	Language string
	Sort     string // last_push_at, stars_count, forks_count, watches_count
	Order    string // asc, desc
	Page     int
	PerPage  int
}

// SearchIssuesOptions contains optional parameters for SearchIssues
type SearchIssuesOptions struct {
	Q        string
	Repo     string
	Language string
	Label    string
	State    string // open, progressing, closed, rejected
	Author   string
	Assignee string
	Sort     string // created_at, updated_at, notes_count
	Order    string // asc, desc
	Page     int
	PerPage  int
}

// SearchUsersOptions contains optional parameters for SearchUsers
type SearchUsersOptions struct {
	Q       string
	Sort    string // joined_at
	Order   string // asc, desc
	Page    int
	PerPage int
}
