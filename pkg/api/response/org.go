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
