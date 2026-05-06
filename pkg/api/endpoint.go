package api

import "fmt"

// HTTPMethod represents an HTTP method
type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PATCH  HTTPMethod = "PATCH"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)

// Endpoint represents an API endpoint configuration
type Endpoint struct {
	Method HTTPMethod
	Path   string
}

// Build builds the path with given args
func (e Endpoint) Build(args ...interface{}) string {
	return fmt.Sprintf(e.Path, args...)
}

// EndpointGroup is a collection of endpoints for a resource
type EndpointGroup struct {
	List                    Endpoint
	Get                     Endpoint
	GetByID                 Endpoint
	Create                  Endpoint
	Update                  Endpoint
	Replace                 Endpoint
	Delete                  Endpoint
	Test                    Endpoint
	Merge                   Endpoint
	Comment                 Endpoint
	SearchRepos             Endpoint
	SearchIssues            Endpoint
	SearchUsers             Endpoint
	ListLicenses            Endpoint
	GetLicense              Endpoint
	GetLicenseRaw           Endpoint
	GetRepoLicense          Endpoint
	ListGitignoreTemplates  Endpoint
	GetGitignoreTemplate    Endpoint
	GetGitignoreTemplateRaw Endpoint
	RenderMarkdown          Endpoint
	Starred                 Endpoint
	Star                    Endpoint
	StarPut                 Endpoint
	StarDel                 Endpoint
	Forks                   Endpoint
	Fork                    Endpoint
	Commits                 Endpoint
	Comments                Endpoint
	CreateComment           Endpoint
	UpdateComment           Endpoint
	DeleteComment           Endpoint
	GetAnnotations          Endpoint
	GetCommitCheckRuns      Endpoint
	Patch                   Endpoint
}

// Repo endpoints
var Repo = EndpointGroup{
	Get: Endpoint{GET, "/repos/%s/%s"},
}

// Orgs Org endpoints
var Orgs = EndpointGroup{
	Get: Endpoint{GET, "/orgs/%s"},
}

// UserOrgs UserOrg endpoints
var UserOrgs = EndpointGroup{
	List: Endpoint{GET, "/user/orgs"},
}

// UserRepos UserRepo endpoints
var UserRepos = EndpointGroup{
	List: Endpoint{GET, "/user/repos"},
}

// Issues Issue endpoints
var Issues = EndpointGroup{
	List:   Endpoint{GET, "/repos/%s/%s/issues"},
	Create: Endpoint{POST, "/repos/%s/issues"},
	Update: Endpoint{PATCH, "/repos/%s/issues/%s"},
}

// PRs PR endpoints
var PRs = EndpointGroup{
	List:    Endpoint{GET, "/repos/%s/%s/pulls"},
	Create:  Endpoint{POST, "/repos/%s/%s/pulls"},
	Merge:   Endpoint{PUT, "/repos/%s/%s/pulls/%d/merge"},
	Update:  Endpoint{PATCH, "/repos/%s/%s/pulls/%d"},
	Comment: Endpoint{POST, "/repos/%s/%s/pulls/%d/comments"},
}

// Releases Release endpoints
var Releases = EndpointGroup{
	List:    Endpoint{GET, "/repos/%s/%s/releases"},
	Create:  Endpoint{POST, "/repos/%s/%s/releases"},
	Get:     Endpoint{GET, "/repos/%s/%s/releases/tags/%s"},
	GetByID: Endpoint{GET, "/repos/%s/%s/releases/%d"},
	Delete:  Endpoint{DELETE, "/repos/%s/%s/releases/%d"},
	Update:  Endpoint{PATCH, "/repos/%s/%s/releases/%d"},
}

// UserReposByName User repos by username endpoints
var UserReposByName = EndpointGroup{
	List: Endpoint{GET, "/users/%s/repos"},
}

// Webhooks Webhook endpoints
var Webhooks = EndpointGroup{
	List:   Endpoint{GET, "/repos/%s/%s/hooks"},
	Get:    Endpoint{GET, "/repos/%s/%s/hooks/%d"},
	Create: Endpoint{POST, "/repos/%s/%s/hooks"},
	Update: Endpoint{PATCH, "/repos/%s/%s/hooks/%d"},
	Delete: Endpoint{DELETE, "/repos/%s/%s/hooks/%d"},
	Test:   Endpoint{POST, "/repos/%s/%s/hooks/%d/tests"},
}

// Search endpoints
var Search = EndpointGroup{
	SearchRepos:  Endpoint{GET, "/search/repositories"},
	SearchIssues: Endpoint{GET, "/search/issues"},
	SearchUsers:  Endpoint{GET, "/search/users"},
}

// Miscellaneous endpoints
var Miscellaneous = EndpointGroup{
	ListLicenses:            Endpoint{GET, "/licenses"},
	GetLicense:              Endpoint{GET, "/licenses/%s"},
	GetLicenseRaw:           Endpoint{GET, "/licenses/%s/raw"},
	GetRepoLicense:          Endpoint{GET, "/repos/%s/%s/license"},
	ListGitignoreTemplates:  Endpoint{GET, "/gitignore/templates"},
	GetGitignoreTemplate:    Endpoint{GET, "/gitignore/templates/%s"},
	GetGitignoreTemplateRaw: Endpoint{GET, "/gitignore/templates/%s/raw"},
	RenderMarkdown:          Endpoint{POST, "/markdown"},
}

// Milestones Milestone endpoints
var Milestones = EndpointGroup{
	List:   Endpoint{GET, "/repos/%s/%s/milestones"},
	Create: Endpoint{POST, "/repos/%s/%s/milestones"},
	Get:    Endpoint{GET, "/repos/%s/%s/milestones/%d"},
	Update: Endpoint{PATCH, "/repos/%s/%s/milestones/%d"},
	Delete: Endpoint{DELETE, "/repos/%s/%s/milestones/%d"},
}

// Labels Label endpoints
var Labels = EndpointGroup{
	List:   Endpoint{GET, "/repos/%s/%s/labels"},
	Create: Endpoint{POST, "/repos/%s/%s/labels"},
	Get:    Endpoint{GET, "/repos/%s/%s/labels/%s"},
	Update: Endpoint{PATCH, "/repos/%s/%s/labels/%s"},
	Delete: Endpoint{DELETE, "/repos/%s/%s/labels/%s"},
}

// IssueLabels Issue labels endpoints
var IssueLabels = EndpointGroup{
	List:    Endpoint{GET, "/repos/%s/%s/issues/%s/labels"},
	Create:  Endpoint{POST, "/repos/%s/%s/issues/%s/labels"},
	Replace: Endpoint{PUT, "/repos/%s/%s/issues/%s/labels"},
	Delete:  Endpoint{DELETE, "/repos/%s/%s/issues/%s/labels"},
}

// ProjectLabels Project labels endpoints
var ProjectLabels = EndpointGroup{
	List:    Endpoint{GET, "/repos/%s/%s/project_labels"},
	Create:  Endpoint{POST, "/repos/%s/%s/project_labels"},
	Replace: Endpoint{PUT, "/repos/%s/%s/project_labels"},
	Delete:  Endpoint{DELETE, "/repos/%s/%s/project_labels"},
}

// EnterpriseLabels Enterprise labels endpoints
var EnterpriseLabels = EndpointGroup{
	List: Endpoint{GET, "/enterprises/%s/labels"},
	Get:  Endpoint{GET, "/enterprises/%s/labels/%s"},
}

// Emails Email endpoints
var Emails = EndpointGroup{
	List: Endpoint{GET, "/emails"},
}

// RepoStats Repo statistics endpoints
var RepoStats = EndpointGroup{
	Get: Endpoint{POST, "/repos/%s/%s/traffic-data"},
}

// RepoLanguages Language statistics endpoints
var RepoLanguages = EndpointGroup{
	Get: Endpoint{GET, "/repos/%s/%s/languages"},
}

// RepoContributors Contributor endpoints
var RepoContributors = EndpointGroup{
	Get: Endpoint{GET, "/repos/%s/%s/contributors"},
}

// Gists Gist endpoints
var Gists = EndpointGroup{
	List:          Endpoint{GET, "/gists"},
	Create:        Endpoint{POST, "/gists"},
	Get:           Endpoint{GET, "/gists/%s"},
	Update:        Endpoint{PATCH, "/gists/%s"},
	Delete:        Endpoint{DELETE, "/gists/%s"},
	Starred:       Endpoint{GET, "/gists/starred"},
	Star:          Endpoint{GET, "/gists/%s/star"},
	StarPut:       Endpoint{PUT, "/gists/%s/star"},
	StarDel:       Endpoint{DELETE, "/gists/%s/star"},
	Forks:         Endpoint{GET, "/gists/%s/forks"},
	Fork:          Endpoint{POST, "/gists/%s/forks"},
	Commits:       Endpoint{GET, "/gists/%s/commits"},
	Comments:      Endpoint{GET, "/gists/%s/comments"},
	Comment:       Endpoint{GET, "/gists/%s/comments/%d"},
	CreateComment: Endpoint{POST, "/gists/%s/comments"},
	UpdateComment: Endpoint{PATCH, "/gists/%s/comments/%d"},
	DeleteComment: Endpoint{DELETE, "/gists/%s/comments/%d"},
}

// CheckRuns Check runs endpoints
var CheckRuns = EndpointGroup{
	Create:             Endpoint{POST, "/repos/%s/%s/check-runs"},
	Get:                Endpoint{GET, "/repos/%s/%s/check-runs/%d"},
	Update:             Endpoint{PATCH, "/repos/%s/%s/check-runs/%d"},
	GetAnnotations:     Endpoint{GET, "/repos/%s/%s/check-runs/%d/annotations"},
	GetCommitCheckRuns: Endpoint{GET, "/repos/%s/%s/commits/%s/check-runs"},
}

// RepoNotifications repo notification endpoints
var RepoNotifications = EndpointGroup{
	List:   Endpoint{GET, "/repos/%s/%s/notifications"},
	Update: Endpoint{PUT, "/repos/%s/%s/notifications"},
}

// NotificationThreads notification threads endpoints
var NotificationThreads = EndpointGroup{
	List:   Endpoint{GET, "/notifications/threads"},
	Update: Endpoint{PUT, "/notifications/threads"},
	Get:    Endpoint{GET, "/notifications/threads/%s"},
	Patch:  Endpoint{PATCH, "/notifications/threads/%s"},
}

// NotificationCount notification count endpoints
var NotificationCount = EndpointGroup{
	Get: Endpoint{GET, "/notifications/count"},
}

// Messages message endpoints
var Messages = EndpointGroup{
	List:   Endpoint{GET, "/notifications/messages"},
	Create: Endpoint{POST, "/notifications/messages"},
	Update: Endpoint{PUT, "/notifications/messages"},
	Get:    Endpoint{GET, "/notifications/messages/%s"},
	Patch:  Endpoint{PATCH, "/notifications/messages/%s"},
}
