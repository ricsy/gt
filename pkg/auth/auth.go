package auth

import (
	"errors"
	"os"
	"sync"

	"github.com/ricsy/gt/pkg/config"
)

// ErrNotLoggedIn is returned when no authentication is found for the host
var ErrNotLoggedIn = errors.New("not logged in")

// Login saves authentication info for a host
func Login(host, token, username string) error {
	hosts, err := config.LoadHosts()
	if err != nil {
		return err
	}

	previousUser := hosts[host].User
	hosts[host] = config.HostAuth{
		Token: token,
		User:  username,
	}

	if err := config.SaveHosts(hosts); err != nil {
		return err
	}

	if previousUser != "" && previousUser != username {
		cfg, err := config.LoadConfig()
		if err != nil {
			return err
		}
		if cfg.DefaultRepo != "" {
			cfg.DefaultRepo = ""
			if err := config.SaveConfig(cfg); err != nil {
				return err
			}
		}
	}

	invalidateAuthCache()

	return nil
}

// Logout removes authentication info for a host
func Logout(host string) error {
	hosts, err := config.LoadHosts()
	if err != nil {
		return err
	}

	delete(hosts, host)

	if err := config.SaveHosts(hosts); err != nil {
		return err
	}

	invalidateAuthCache()

	return nil
}

// CurrentUser returns the logged-in username for the host
func CurrentUser(host string) (string, error) {
	auth, err := GetAuth(host)
	if err != nil {
		return "", err
	}
	return auth.User, nil
}

// IsLoggedIn checks if the host has authentication info
func IsLoggedIn(host string) bool {
	_, err := GetAuth(host)
	return err == nil
}

// GetToken returns the token for the host
func GetToken(host string) (string, error) {
	if token := os.Getenv("GITEE_TOKEN"); token != "" {
		return token, nil
	}
	auth, err := GetAuth(host)
	if err != nil {
		return "", err
	}
	return auth.Token, nil
}

var (
	authCache    map[string]config.HostAuth
	authCacheMu  sync.RWMutex
	authCacheErr error
)

// GetAuth returns the HostAuth for the host (cached)
func GetAuth(host string) (config.HostAuth, error) {
	authCacheMu.RLock()
	if authCache != nil {
		auth, ok := authCache[host]
		authCacheMu.RUnlock()
		if ok && auth.Token != "" {
			return auth, nil
		}
		return config.HostAuth{}, ErrNotLoggedIn
	}
	authCacheMu.RUnlock()

	// Load into cache
	authCacheMu.Lock()
	defer authCacheMu.Unlock()

	// Double-check after acquiring write lock
	if authCache != nil {
		auth, ok := authCache[host]
		if ok && auth.Token != "" {
			return auth, nil
		}
		return config.HostAuth{}, ErrNotLoggedIn
	}

	authCache, authCacheErr = config.LoadHosts()
	if authCacheErr != nil {
		return config.HostAuth{}, authCacheErr
	}

	auth, ok := authCache[host]
	if !ok || auth.Token == "" {
		return config.HostAuth{}, ErrNotLoggedIn
	}

	return auth, nil
}

// invalidateAuthCache clears the auth cache (call after login/logout)
func invalidateAuthCache() {
	authCacheMu.Lock()
	authCache = nil
	authCacheMu.Unlock()
}
