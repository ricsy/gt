package api

import (
	"strconv"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// UserDetail is an alias for response.UserDetail.
type UserDetail = response.UserDetail

// UpdateUserOptions is an alias for response.UpdateUserOptions.
type UpdateUserOptions = response.UpdateUserOptions

// ListUsersOptions is an alias for response.ListUsersOptions.
type ListUsersOptions = response.ListUsersOptions

// SSHKey is an alias for response.SSHKey.
type SSHKey = response.SSHKey

// CreateSSHKeyOptions is an alias for response.CreateSSHKeyOptions.
type CreateSSHKeyOptions = response.CreateSSHKeyOptions

// Namespace is an alias for response.Namespace.
type Namespace = response.Namespace

// ListNamespacesOptions is an alias for response.ListNamespacesOptions.
type ListNamespacesOptions = response.ListNamespacesOptions

// GetAuthenticatedUser gets the authenticated user's profile.
func (c *Client) GetAuthenticatedUser() (*UserDetail, error) {
	var user UserDetail
	err := c.DoFromEndpoint(Users.Get, nil, nil, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateAuthenticatedUser updates the authenticated user's profile.
func (c *Client) UpdateAuthenticatedUser(opts UpdateUserOptions) (*User, error) {
	var user User
	err := c.DoFromEndpoint(Users.Update, nil, opts, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListFollowers lists followers for the authenticated user.
func (c *Client) ListFollowers(opts ListUsersOptions) ([]User, error) {
	var users []User
	err := c.Do("GET", Users.Followers.Path+buildUsersQuery(opts), nil, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ListFollowing lists users followed by the authenticated user.
func (c *Client) ListFollowing(opts ListUsersOptions) ([]User, error) {
	var users []User
	err := c.Do("GET", Users.Following.Path+buildUsersQuery(opts), nil, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CheckFollowing checks whether the authenticated user follows username.
func (c *Client) CheckFollowing(username string) error {
	return c.DoFromEndpoint(Users.CheckFollowing, []interface{}{username}, nil, nil)
}

// FollowUser follows username as the authenticated user.
func (c *Client) FollowUser(username string) error {
	return c.DoFromEndpoint(Users.Follow, []interface{}{username}, nil, nil)
}

// UnfollowUser unfollows username as the authenticated user.
func (c *Client) UnfollowUser(username string) error {
	return c.DoFromEndpoint(Users.Unfollow, []interface{}{username}, nil, nil)
}

// ListSSHKeys lists SSH keys for the authenticated user.
func (c *Client) ListSSHKeys(opts ListUsersOptions) ([]SSHKey, error) {
	var keys []SSHKey
	err := c.Do("GET", Users.Keys.Path+buildUsersQuery(opts), nil, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// CreateSSHKey creates an SSH key for the authenticated user.
func (c *Client) CreateSSHKey(opts CreateSSHKeyOptions) (*SSHKey, error) {
	var key SSHKey
	err := c.DoFromEndpoint(Users.Create, nil, opts, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// GetSSHKey gets an SSH key for the authenticated user.
func (c *Client) GetSSHKey(id int64) (*SSHKey, error) {
	var key SSHKey
	err := c.DoFromEndpoint(Users.Key, []interface{}{id}, nil, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// DeleteSSHKey deletes an SSH key for the authenticated user.
func (c *Client) DeleteSSHKey(id int64) error {
	return c.DoFromEndpoint(Users.Delete, []interface{}{id}, nil, nil)
}

// ListNamespaces lists namespaces for the authenticated user.
func (c *Client) ListNamespaces(opts ListNamespacesOptions) ([]Namespace, error) {
	var namespaces []Namespace
	err := c.Do("GET", Users.Namespaces.Path+buildNamespacesQuery(opts), nil, &namespaces)
	if err != nil {
		return nil, err
	}
	return namespaces, nil
}

// GetNamespace gets a namespace by path for the authenticated user.
func (c *Client) GetNamespace(path string) (*Namespace, error) {
	var namespace Namespace
	query := "?" + util.BuildQuery("path", path)
	err := c.Do("GET", Users.Namespace.Path+query, nil, &namespace)
	if err != nil {
		return nil, err
	}
	return &namespace, nil
}

// GetUser gets a public user profile.
func (c *Client) GetUser(username string) (*response.UserInfo, error) {
	var user response.UserInfo
	err := c.DoFromEndpoint(PublicUsers.Get, []interface{}{username}, nil, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListUserFollowers lists followers for username.
func (c *Client) ListUserFollowers(username string, opts ListUsersOptions) ([]User, error) {
	var users []User
	path := PublicUsers.Followers.Build(username)
	err := c.Do("GET", path+buildUsersQuery(opts), nil, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ListUserFollowing lists users followed by username.
func (c *Client) ListUserFollowing(username string, opts ListUsersOptions) ([]User, error) {
	var users []User
	path := PublicUsers.Following.Build(username)
	err := c.Do("GET", path+buildUsersQuery(opts), nil, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CheckUserFollowing checks whether username follows targetUser.
func (c *Client) CheckUserFollowing(username, targetUser string) error {
	return c.DoFromEndpoint(PublicUsers.CheckFollowing, []interface{}{username, targetUser}, nil, nil)
}

// ListUserSSHKeys lists public SSH keys for username.
func (c *Client) ListUserSSHKeys(username string, opts ListUsersOptions) ([]response.SSHKeyBasic, error) {
	var keys []response.SSHKeyBasic
	path := PublicUsers.Keys.Build(username)
	err := c.Do("GET", path+buildUsersQuery(opts), nil, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func buildUsersQuery(opts ListUsersOptions) string {
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

func buildNamespacesQuery(opts ListNamespacesOptions) string {
	var params []string
	if opts.Mode != "" {
		params = append(params, "mode", opts.Mode)
	}
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
