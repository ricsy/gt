package auth

import (
	"os"
	"testing"

	"github.com/ricsy/gt/pkg/config"
)

var originalHome string

func TestMain(m *testing.M) {
	originalHome = os.Getenv("HOME")

	tmpDir, err := os.MkdirTemp("", "gt-auth-test")
	if err != nil {
		os.Exit(1)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	_ = os.Setenv("HOME", tmpDir)
	_ = os.Unsetenv("GITEE_TOKEN")
	config.SetConfigDirFunc(func() string { return tmpDir })

	code := m.Run()

	_ = os.Setenv("HOME", originalHome)

	os.Exit(code)
}

func TestLogin(t *testing.T) {
	err := Login("gitee.com", "test-token", "testuser")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	hosts, err := config.LoadHosts()
	if err != nil {
		t.Fatalf("LoadHosts failed: %v", err)
	}

	auth, ok := hosts["gitee.com"]
	if !ok {
		t.Fatal("gitee.com not found in hosts")
	}
	if auth.Token != "test-token" {
		t.Errorf("token = %s, want test-token", auth.Token)
	}
	if auth.User != "testuser" {
		t.Errorf("user = %s, want testuser", auth.User)
	}
}

func TestLogout(t *testing.T) {
	// Use a fresh host to avoid pollution from other tests
	host := "logout-test.example.com"

	err := Login(host, "test-token", "testuser")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	err = Logout(host)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	hosts, err := config.LoadHosts()
	if err != nil {
		t.Fatalf("LoadHosts failed: %v", err)
	}

	if _, ok := hosts[host]; ok {
		t.Fatal("host should be removed after logout")
	}
}

func TestIsLoggedIn(t *testing.T) {
	host := "isloggedin-test.example.com"

	if IsLoggedIn(host) {
		t.Fatal("should not be logged in initially")
	}

	err := Login(host, "test-token", "testuser")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if !IsLoggedIn(host) {
		t.Fatal("should be logged in after Login")
	}

	err = Logout(host)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	if IsLoggedIn(host) {
		t.Fatal("should not be logged in after Logout")
	}
}

func TestCurrentUser(t *testing.T) {
	host := "currentuser-test.example.com"

	_, err := CurrentUser(host)
	if err != ErrNotLoggedIn {
		t.Fatalf("CurrentUser on non-logged-in host = %v, want ErrNotLoggedIn", err)
	}

	err = Login(host, "test-token", "testuser")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	user, err := CurrentUser(host)
	if err != nil {
		t.Fatalf("CurrentUser failed: %v", err)
	}
	if user != "testuser" {
		t.Errorf("user = %s, want testuser", user)
	}
}

func TestGetToken(t *testing.T) {
	host := "gettoken-test.example.com"

	_, err := GetToken(host)
	if err != ErrNotLoggedIn {
		t.Fatalf("GetToken on non-logged-in host = %v, want ErrNotLoggedIn", err)
	}

	err = Login(host, "my-secret-token", "testuser")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	token, err := GetToken(host)
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}
	if token != "my-secret-token" {
		t.Errorf("token = %s, want my-secret-token", token)
	}
}

func TestGetAuth(t *testing.T) {
	host := "getauth-test.example.com"

	_, err := GetAuth(host)
	if err != ErrNotLoggedIn {
		t.Fatalf("GetAuth on non-logged-in host = %v, want ErrNotLoggedIn", err)
	}

	err = Login(host, "test-token", "testuser")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	auth, err := GetAuth(host)
	if err != nil {
		t.Fatalf("GetAuth failed: %v", err)
	}
	if auth.Token != "test-token" || auth.User != "testuser" || auth.Host != host {
		t.Errorf("GetAuth = %+v, want {Token:test-token User:testuser Host:%s}", auth, host)
	}
}

func TestLoginMultipleHosts(t *testing.T) {
	err := Login("gitee.com", "token1", "user1")
	if err != nil {
		t.Fatalf("Login gitee.com failed: %v", err)
	}
	err = Login("github.com", "token2", "user2")
	if err != nil {
		t.Fatalf("Login github.com failed: %v", err)
	}

	if !IsLoggedIn("gitee.com") {
		t.Fatal("gitee.com should be logged in")
	}
	if !IsLoggedIn("github.com") {
		t.Fatal("github.com should be logged in")
	}

	user1, _ := CurrentUser("gitee.com")
	user2, _ := CurrentUser("github.com")
	if user1 != "user1" || user2 != "user2" {
		t.Errorf("users = %s, %s, want user1, user2", user1, user2)
	}

	err = Logout("gitee.com")
	if err != nil {
		t.Fatalf("Logout gitee.com failed: %v", err)
	}

	if IsLoggedIn("gitee.com") {
		t.Fatal("gitee.com should not be logged in after logout")
	}
	if !IsLoggedIn("github.com") {
		t.Fatal("github.com should still be logged in")
	}
}
