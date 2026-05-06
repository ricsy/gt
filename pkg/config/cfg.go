package config

// Config holds user preferences
type Config struct {
	DefaultRepo  string `yaml:"default_repo"`
	DefaultOwner string `yaml:"default_owner"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{}
}
