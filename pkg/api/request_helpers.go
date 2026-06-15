package api

import (
	"strconv"

	"github.com/ricsy/gt/pkg/util"
)

func (c *Client) doFromEndpointWithQuery(e Endpoint, pathArgs []interface{}, query string, body interface{}, response interface{}) error {
	if query != "" {
		e.Path += "?" + query
	}
	return c.DoFromEndpoint(e, pathArgs, body, response)
}

func (c *Client) doGetWithQuery(path, query string, response interface{}) error {
	return c.Do("GET", path+query, nil, response)
}

func buildOptionalQuery(params ...string) string {
	query := util.BuildQuery(params...)
	if query == "" {
		return ""
	}
	return "?" + query
}

func paginationParams(page, perPage int) []string {
	var params []string
	if page > 0 {
		params = append(params, "page", strconv.Itoa(page))
	}
	if perPage > 0 {
		params = append(params, "per_page", strconv.Itoa(perPage))
	}
	return params
}
