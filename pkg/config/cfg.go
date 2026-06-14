package config

// Config holds user preferences
type Config struct {
	DefaultRepo        string `yaml:"default_repo"`
	DefaultOwner       string `yaml:"default_owner"`
	DefaultHost        string `yaml:"default_host"`
	RepoScopeMode      string `yaml:"repo_scope_mode,omitempty"`
	RepoScopeNamespace string `yaml:"repo_scope_namespace,omitempty"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		DefaultHost: DefaultHost,
	}
}
