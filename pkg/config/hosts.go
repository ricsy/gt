package config

// HostAuth holds authentication info for a host
type HostAuth struct {
	Token string `yaml:"token"`
	User  string `yaml:"user"`
	Host  string `yaml:"host"`
}

// DefaultHost is the default host (gitee.com)
const DefaultHost = "gitee.com"
