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
	List    Endpoint
	Get     Endpoint
	GetByID Endpoint
	Create  Endpoint
	Update  Endpoint
	Delete  Endpoint
	Test    Endpoint
	Merge   Endpoint
	Comment Endpoint
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
