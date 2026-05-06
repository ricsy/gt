package api

// Org represents a Gitee organization
type Org struct {
	Login    string `json:"login"`
	Name     string `json:"name"`
	Blog     string `json:"blog"`
	Email    string `json:"email"`
	HtmlUrl  string `json:"html_url"`
	Location string `json:"location"`
}

// ListOrgs lists organizations for the current user
func (c *Client) ListOrgs() ([]Org, error) {
	var orgs []Org
	err := c.DoFromEndpoint(UserOrgs.List, nil, nil, &orgs)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

// GetOrg gets an organization by login
func (c *Client) GetOrg(login string) (*Org, error) {
	var org Org
	err := c.DoFromEndpoint(Orgs.List, []interface{}{login}, nil, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}
