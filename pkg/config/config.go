package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	configDirName  = "gt"
	hostsFileName  = "hosts.yml"
	configFileName = "config.yml"
)

// ConfigDirImpl is the actual implementation of ConfigDir
func ConfigDirImpl() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".config", configDirName)
}

// configDirFunc is the ConfigDir implementation (can be overridden for testing)
var configDirFunc = ConfigDirImpl

// SetConfigDirFunc overrides ConfigDir for testing
func SetConfigDirFunc(fn func() string) {
	configDirFunc = fn
}

// ConfigDir returns the gt config directory path (~/.config/gt)
func ConfigDir() string {
	return configDirFunc()
}

// HostsFile returns the hosts.yml path
func HostsFile() string {
	return filepath.Join(ConfigDir(), hostsFileName)
}

// ConfigFile returns the config.yml path
func ConfigFile() string {
	return filepath.Join(ConfigDir(), configFileName)
}

// LoadHosts reads and parses hosts.yml
func LoadHosts() (map[string]HostAuth, error) {
	data, err := os.ReadFile(HostsFile())
	if err != nil {
		if os.IsNotExist(err) {
			return defaultHosts(), nil
		}
		return nil, err
	}

	var hosts map[string]HostAuth
	if err := yaml.Unmarshal(data, &hosts); err != nil {
		return nil, err
	}

	if hosts == nil {
		return defaultHosts(), nil
	}

	return hosts, nil
}

// SaveHosts writes hosts to hosts.yml
func SaveHosts(hosts map[string]HostAuth) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(hosts)
	if err != nil {
		return err
	}

	return os.WriteFile(HostsFile(), data, 0600)
}

// LoadConfig reads and parses config.yml
func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigFile())
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig writes config to config.yml
func SaveConfig(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigFile(), data, 0644)
}

// defaultHosts returns default hosts map with gitee.com
func defaultHosts() map[string]HostAuth {
	return map[string]HostAuth{
		"gitee.com": {
			Host: "gitee.com",
		},
	}
}
