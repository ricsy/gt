package api

import (
	"strconv"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// Event is an alias for response.Event.
type Event = response.Event

// ListActivityOptions is an alias for response.ListActivityOptions.
type ListActivityOptions = response.ListActivityOptions

// ListRepoEvents lists repository events.
func (c *Client) ListRepoEvents(owner, repo string, opts ListActivityOptions) ([]Event, error) {
	var events []Event
	path := Activity.List.Build(owner, repo)
	err := c.Do("GET", path+buildActivityQuery(opts), nil, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListNetworkEvents lists public repository network events.
func (c *Client) ListNetworkEvents(owner, repo string, opts ListActivityOptions) ([]Event, error) {
	var events []Event
	path := Activity.NetworkEvents.Build(owner, repo)
	err := c.Do("GET", path+buildActivityQuery(opts), nil, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListOrgEvents lists organization public events.
func (c *Client) ListOrgEvents(org string, opts ListActivityOptions) ([]Event, error) {
	var events []Event
	path := Activity.OrgEvents.Build(org)
	err := c.Do("GET", path+buildActivityQuery(opts), nil, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListUserEvents lists user events.
func (c *Client) ListUserEvents(username string, publicOnly bool, opts ListActivityOptions) ([]Event, error) {
	var events []Event
	endpoint := Activity.Events
	if publicOnly {
		endpoint = Activity.PublicEvents
	}
	path := endpoint.Build(username)
	err := c.Do("GET", path+buildActivityQuery(opts), nil, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListUserReceivedEvents lists events received by a user.
func (c *Client) ListUserReceivedEvents(username string, publicOnly bool, opts ListActivityOptions) ([]Event, error) {
	var events []Event
	endpoint := Activity.ReceivedEvents
	if publicOnly {
		endpoint = Activity.ReceivedPublicEvents
	}
	path := endpoint.Build(username)
	err := c.Do("GET", path+buildActivityQuery(opts), nil, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListSubscriptions lists repositories watched by the authenticated user.
func (c *Client) ListSubscriptions(opts ListActivityOptions) ([]Repository, error) {
	var repos []Repository
	err := c.Do("GET", Activity.Subscriptions.Path+buildActivityQuery(opts), nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// ListUserSubscriptions lists repositories watched by a user.
func (c *Client) ListUserSubscriptions(username string, opts ListActivityOptions) ([]Repository, error) {
	var repos []Repository
	path := Activity.Get.Build(username)
	err := c.Do("GET", path+buildActivityQuery(opts), nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// CheckSubscription checks whether the authenticated user watches a repository.
func (c *Client) CheckSubscription(owner, repo string) error {
	return c.DoFromEndpoint(Activity.Subscription, []interface{}{owner, repo}, nil, nil)
}

// WatchRepo watches a repository.
func (c *Client) WatchRepo(owner, repo string) error {
	req := map[string]string{"watch_type": "watching"}
	return c.DoFromEndpoint(Activity.StarPut, []interface{}{owner, repo}, req, nil)
}

// UnwatchRepo unwatches a repository.
func (c *Client) UnwatchRepo(owner, repo string) error {
	return c.DoFromEndpoint(Activity.StarDel, []interface{}{owner, repo}, nil, nil)
}

// ListSubscribers lists users watching a repository.
func (c *Client) ListSubscribers(owner, repo string, opts ListActivityOptions) ([]User, error) {
	var users []User
	path := Activity.Subscribers.Build(owner, repo)
	err := c.Do("GET", path+buildActivityQuery(opts), nil, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func buildActivityQuery(opts ListActivityOptions) string {
	var params []string
	if opts.Page > 0 {
		params = append(params, "page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params = append(params, "per_page", strconv.Itoa(opts.PerPage))
	}
	if len(params) == 0 {
		return ""
	}
	return "?" + util.BuildQuery(params...)
}
