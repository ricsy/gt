package config

// HostAuth holds authentication info for a host
type HostAuth struct {
	Token string `yaml:"token"`
	User  string `yaml:"user"`
}

const DefaultHost = "gitee.com"

// ApiUrl returns the base API URL for a given host.
func ApiUrl(host string) string {
	return "https://" + host + "/api/v5"
}
