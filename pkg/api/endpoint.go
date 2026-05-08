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
	Latest                  Endpoint
	Create                  Endpoint
	Update                  Endpoint
	Add                     Endpoint
	Remove                  Endpoint
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
	Followers               Endpoint
	Following               Endpoint
	CheckFollowing          Endpoint
	Follow                  Endpoint
	Unfollow                Endpoint
	Keys                    Endpoint
	Key                     Endpoint
	Namespaces              Endpoint
	Namespace               Endpoint
	Members                 Endpoint
	Member                  Endpoint
	SearchMembers           Endpoint
	Repos                   Endpoint
	PullRequests            Endpoint
	Protection              Endpoint
	Blob                    Endpoint
	Tree                    Endpoint
	Metrics                 Endpoint
	Events                  Endpoint
	PublicEvents            Endpoint
	ReceivedEvents          Endpoint
	ReceivedPublicEvents    Endpoint
	OrgEvents               Endpoint
	NetworkEvents           Endpoint
	Subscriptions           Endpoint
	Subscription            Endpoint
	Subscribers             Endpoint
	OrgMembers              Endpoint
	OrgRepos                Endpoint
	GetUserOrgsByUsername   Endpoint
}

// Repo endpoints
var Repo = EndpointGroup{
	Get:    Endpoint{GET, "/repos/%s/%s"},
	Update: Endpoint{PATCH, "/repos/%s/%s"},
	Delete: Endpoint{DELETE, "/repos/%s/%s"},
}

// Branches repository branch endpoints
var Branches = EndpointGroup{
	List:       Endpoint{GET, "/repos/%s/%s/branches"},
	Create:     Endpoint{POST, "/repos/%s/%s/branches"},
	Get:        Endpoint{GET, "/repos/%s/%s/branches/%s"},
	Protection: Endpoint{PUT, "/repos/%s/%s/branches/%s/protection"},
	Delete:     Endpoint{DELETE, "/repos/%s/%s/branches/%s/protection"},
}

// Orgs Org endpoints
var Orgs = EndpointGroup{
	Get:                   Endpoint{GET, "/orgs/%s"},
	Members:               Endpoint{GET, "/orgs/%s/members"},
	Repos:                 Endpoint{GET, "/orgs/%s/repos"},
	GetUserOrgsByUsername: Endpoint{GET, "/users/%s/orgs"},
}

// Collaborators collaborator endpoints
var Collaborators = EndpointGroup{
	List:       Endpoint{GET, "/repos/%s/%s/collaborators"},
	Get:        Endpoint{GET, "/repos/%s/%s/collaborators/%s"},
	Protection: Endpoint{GET, "/repos/%s/%s/collaborators/%s/permission"},
	Add:        Endpoint{PUT, "/repos/%s/%s/collaborators/%s"},
	Remove:     Endpoint{DELETE, "/repos/%s/%s/collaborators/%s"},
}

// RepoForks fork endpoints
var RepoForks = EndpointGroup{
	List: Endpoint{GET, "/repos/%s/%s/forks"},
	Fork: Endpoint{POST, "/repos/%s/%s/forks"},
}

// UserOrgs UserOrg endpoints
var UserOrgs = EndpointGroup{
	List: Endpoint{GET, "/user/orgs"},
}

// UserRepos UserRepo endpoints
var UserRepos = EndpointGroup{
	List:   Endpoint{GET, "/user/repos"},
	Create: Endpoint{POST, "/user/repos"},
}

// Issues Issue endpoints
var Issues = EndpointGroup{
	List:   Endpoint{GET, "/repos/%s/%s/issues"},
	Get:    Endpoint{GET, "/repos/%s/%s/issues/%s"},
	Create: Endpoint{POST, "/repos/%s/issues"},
	Update: Endpoint{PATCH, "/repos/%s/issues/%s"},
}

// PRs PR endpoints
var PRs = EndpointGroup{
	List:    Endpoint{GET, "/repos/%s/%s/pulls"},
	Get:     Endpoint{GET, "/repos/%s/%s/pulls/%d"},
	Create:  Endpoint{POST, "/repos/%s/%s/pulls"},
	Merge:   Endpoint{PUT, "/repos/%s/%s/pulls/%d/merge"},
	Update:  Endpoint{PATCH, "/repos/%s/%s/pulls/%d"},
	Comment: Endpoint{POST, "/repos/%s/%s/pulls/%d/comments"},
}

// PRComments PR comment endpoints
var PRComments = EndpointGroup{
	List: Endpoint{GET, "/repos/%s/%s/pulls/%d/comments"},
	Get:  Endpoint{GET, "/repos/%s/%s/pulls/%d/comments/%d"},
}

// PRCommits PR commit endpoints
var PRCommits = EndpointGroup{
	List: Endpoint{GET, "/repos/%s/%s/pulls/%d/commits"},
}

// PRFiles PR file endpoints
var PRFiles = EndpointGroup{
	List: Endpoint{GET, "/repos/%s/%s/pulls/%d/files"},
}

// Releases Release endpoints
var Releases = EndpointGroup{
	List:    Endpoint{GET, "/repos/%s/%s/releases"},
	Create:  Endpoint{POST, "/repos/%s/%s/releases"},
	Get:     Endpoint{GET, "/repos/%s/%s/releases/tags/%s"},
	GetByID: Endpoint{GET, "/repos/%s/%s/releases/%d"},
	Delete:  Endpoint{DELETE, "/repos/%s/%s/releases/%d"},
	Update:  Endpoint{PATCH, "/repos/%s/%s/releases/%d"},
	Latest:  Endpoint{GET, "/repos/%s/%s/releases/latest"},
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

// IssueComments Issue comments endpoints
var IssueComments = EndpointGroup{
	List:   Endpoint{GET, "/repos/%s/%s/issues/%s/comments"},
	Get:    Endpoint{GET, "/repos/%s/%s/issues/%s/comments/%d"},
	Create: Endpoint{POST, "/repos/%s/%s/issues/%s/comments"},
	Update: Endpoint{PATCH, "/repos/%s/%s/issues/%s/comments/%d"},
	Delete: Endpoint{DELETE, "/repos/%s/%s/issues/%s/comments/%d"},
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

// Enterprises Enterprise endpoints
var Enterprises = EndpointGroup{
	List:          Endpoint{GET, "/user/enterprises"},
	Get:           Endpoint{GET, "/enterprises/%s"},
	Members:       Endpoint{GET, "/enterprises/%s/members"},
	Create:        Endpoint{POST, "/enterprises/%s/members"},
	SearchMembers: Endpoint{GET, "/enterprises/%s/members/search"},
	Member:        Endpoint{GET, "/enterprises/%s/members/%s"},
	Update:        Endpoint{PUT, "/enterprises/%s/members/%s"},
	Delete:        Endpoint{DELETE, "/enterprises/%s/members/%s"},
	Repos:         Endpoint{GET, "/enterprises/%s/repos"},
	PullRequests:  Endpoint{GET, "/enterprises/%s/pull_requests"},
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

// Users user endpoints
var Users = EndpointGroup{
	Get:            Endpoint{GET, "/user"},
	Update:         Endpoint{PATCH, "/user"},
	Followers:      Endpoint{GET, "/user/followers"},
	Following:      Endpoint{GET, "/user/following"},
	CheckFollowing: Endpoint{GET, "/user/following/%s"},
	Follow:         Endpoint{PUT, "/user/following/%s"},
	Unfollow:       Endpoint{DELETE, "/user/following/%s"},
	Keys:           Endpoint{GET, "/user/keys"},
	Create:         Endpoint{POST, "/user/keys"},
	Key:            Endpoint{GET, "/user/keys/%d"},
	Delete:         Endpoint{DELETE, "/user/keys/%d"},
	Namespaces:     Endpoint{GET, "/user/namespaces"},
	Namespace:      Endpoint{GET, "/user/namespace"},
}

// PublicUsers public user endpoints
var PublicUsers = EndpointGroup{
	Get:            Endpoint{GET, "/users/%s"},
	Followers:      Endpoint{GET, "/users/%s/followers"},
	Following:      Endpoint{GET, "/users/%s/following"},
	CheckFollowing: Endpoint{GET, "/users/%s/following/%s"},
	Keys:           Endpoint{GET, "/users/%s/keys"},
}

// GitData git data endpoints
var GitData = EndpointGroup{
	Blob:    Endpoint{GET, "/repos/%s/%s/git/blobs/%s"},
	Tree:    Endpoint{GET, "/repos/%s/%s/git/trees/%s"},
	Metrics: Endpoint{GET, "/repos/%s/%s/git/gitee_metrics"},
}

// Activity event and subscription endpoints
var Activity = EndpointGroup{
	Events:               Endpoint{GET, "/users/%s/events"},
	PublicEvents:         Endpoint{GET, "/users/%s/events/public"},
	ReceivedEvents:       Endpoint{GET, "/users/%s/received_events"},
	ReceivedPublicEvents: Endpoint{GET, "/users/%s/received_events/public"},
	OrgEvents:            Endpoint{GET, "/orgs/%s/events"},
	NetworkEvents:        Endpoint{GET, "/networks/%s/%s/events"},
	List:                 Endpoint{GET, "/repos/%s/%s/events"},
	Subscriptions:        Endpoint{GET, "/user/subscriptions"},
	Subscription:         Endpoint{GET, "/user/subscriptions/%s/%s"},
	StarPut:              Endpoint{PUT, "/user/subscriptions/%s/%s"},
	StarDel:              Endpoint{DELETE, "/user/subscriptions/%s/%s"},
	Subscribers:          Endpoint{GET, "/repos/%s/%s/subscribers"},
	Get:                  Endpoint{GET, "/users/%s/subscriptions"},
}
