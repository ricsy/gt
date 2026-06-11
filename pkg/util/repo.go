package util

import (
	"net/url"
	"strings"
)

// SplitOwnerRepo splits "owner/repo" into owner and repo parts.
func SplitOwnerRepo(repo string) (owner, repoName string) {
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", repo
}

// ExtractDigits returns only digit characters from input.
func ExtractDigits(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
}

// BuildQuery builds a query string from key-value pairs.
func BuildQuery(params ...string) string {
	if len(params)%2 != 0 {
		return ""
	}

	values := url.Values{}
	for i := 0; i < len(params); i += 2 {
		key, value := params[i], params[i+1]
		if value == "" {
			continue
		}
		values.Set(key, value)
	}
	return values.Encode()
}
