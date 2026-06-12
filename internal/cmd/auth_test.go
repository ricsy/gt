package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

func TestAuthCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	authCmd.SetOut(buf)
	authCmd.SetArgs([]string{})

	err := authCmd.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestAuthSubcommands(t *testing.T) {
	subcommands := authCmd.Commands()

	expectedCmds := []string{"login", "logout", "status", "token", "setup", "doctor"}
	foundCmds := make(map[string]bool)

	for _, c := range subcommands {
		foundCmds[c.Name()] = true
	}

	for _, expected := range expectedCmds {
		if !foundCmds[expected] {
			t.Errorf("Expected subcommand %q not found", expected)
		}
	}
}

func TestLoginCmdHasFlags(t *testing.T) {
	hasHost := loginCmd.Flags().Lookup("host") != nil
	hasToken := loginCmd.Flags().Lookup("token") != nil
	hasUsername := loginCmd.Flags().Lookup("username") != nil

	if !hasHost {
		t.Error("loginCmd should have --host flag")
	}
	if !hasToken {
		t.Error("loginCmd should have --token flag")
	}
	if !hasUsername {
		t.Error("loginCmd should have --username flag")
	}
}

func TestLogoutCmdHasFlags(t *testing.T) {
	hasHost := logoutCmd.Flags().Lookup("host") != nil
	if !hasHost {
		t.Error("logoutCmd should have --host flag")
	}
}

func TestStatusCmdHasFlags(t *testing.T) {
	hasHost := statusCmd.Flags().Lookup("host") != nil
	if !hasHost {
		t.Error("statusCmd should have --host flag")
	}
}

func TestTokenCmdHasFlags(t *testing.T) {
	hasHost := tokenCmd.Flags().Lookup("host") != nil
	if !hasHost {
		t.Error("tokenCmd should have --host flag")
	}

	hasShow := tokenCmd.Flags().Lookup("show") != nil
	if !hasShow {
		t.Error("tokenCmd should have --show flag")
	}
}

func TestSetupCmdHasFlags(t *testing.T) {
	hasHost := setupCmd.Flags().Lookup("host") != nil
	if !hasHost {
		t.Error("setupCmd should have --host flag")
	}
}

func TestDoctorCmdHasFlags(t *testing.T) {
	hasHost := doctorCmd.Flags().Lookup("host") != nil
	if !hasHost {
		t.Error("doctorCmd should have --host flag")
	}
}

func TestAuthCmdExecution(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cobra.Command
		args []string
	}{
		{"login help", loginCmd, []string{"--help"}},
		{"logout help", logoutCmd, []string{"--help"}},
		{"status help", statusCmd, []string{"--help"}},
		{"token help", tokenCmd, []string{"--help"}},
		{"setup help", setupCmd, []string{"--help"}},
		{"doctor help", doctorCmd, []string{"--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tt.cmd.SetOut(buf)
			tt.cmd.SetArgs(tt.args)

			err := tt.cmd.Execute()
			if err != nil {
				t.Errorf("Execute() error = %v", err)
			}
		})
	}
}

func TestMaskToken(t *testing.T) {
	if got := maskToken("abcd"); got != "****" {
		t.Fatalf("maskToken(short) = %s, want ****", got)
	}

	if got := maskToken("abcdef1234"); got != "ab******34" {
		t.Fatalf("maskToken(long) = %s, want ab******34", got)
	}
}

func TestSetupCmdStoresGitCredential(t *testing.T) {
	originalConfigDirFunc := config.ConfigDir
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	originalApprove := gitCredentialApprove
	t.Cleanup(func() {
		gitCredentialApprove = originalApprove
	})

	if err := auth.Login("gitee.com", "test-token", "test-user"); err != nil {
		t.Fatalf("auth.Login() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = auth.Logout("gitee.com")
	})

	called := false
	gitCredentialApprove = func(host, username, token string) error {
		called = true
		if host != "gitee.com" {
			t.Fatalf("host = %q, want gitee.com", host)
		}
		if username != "test-user" {
			t.Fatalf("username = %q, want test-user", username)
		}
		if token != "test-token" {
			t.Fatalf("token = %q, want test-token", token)
		}
		return nil
	}

	if err := setupCmd.RunE(setupCmd, nil); err != nil {
		t.Fatalf("setupCmd.RunE() returned error: %v", err)
	}
	if !called {
		t.Fatal("gitCredentialApprove was not called")
	}
}

func TestSetupCmdRequiresStoredAuth(t *testing.T) {
	originalConfigDirFunc := config.ConfigDir
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	if err := auth.Logout("gitee.com"); err != nil {
		t.Fatalf("auth.Logout() returned error: %v", err)
	}

	err := setupCmd.RunE(setupCmd, nil)
	if err == nil {
		t.Fatal("setupCmd.RunE() error = nil, want non-nil without stored auth")
	}
}

func TestDoctorCmdPassesWhenStoredAuthAndGitCredentialExist(t *testing.T) {
	originalConfigDirFunc := config.ConfigDir
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	originalFill := gitCredentialFill
	originalHelper := gitCredentialHelperGet
	t.Cleanup(func() {
		gitCredentialFill = originalFill
		gitCredentialHelperGet = originalHelper
	})

	if err := auth.Login("gitee.com", "test-token", "test-user"); err != nil {
		t.Fatalf("auth.Login() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = auth.Logout("gitee.com")
	})

	gitCredentialHelperGet = func() (string, error) {
		return "manager", nil
	}
	gitCredentialFill = func(host string) (config.HostAuth, error) {
		return config.HostAuth{
			Token: "test-token",
			User:  "test-user",
		}, nil
	}

	if err := doctorCmd.RunE(doctorCmd, nil); err != nil {
		t.Fatalf("doctorCmd.RunE() returned error: %v", err)
	}
}

func TestDoctorCmdFailsWhenGitCredentialMissing(t *testing.T) {
	originalConfigDirFunc := config.ConfigDir
	configDir := t.TempDir()
	config.SetConfigDirFunc(func() string { return configDir })
	t.Cleanup(func() {
		config.SetConfigDirFunc(originalConfigDirFunc)
	})

	originalFill := gitCredentialFill
	originalHelper := gitCredentialHelperGet
	t.Cleanup(func() {
		gitCredentialFill = originalFill
		gitCredentialHelperGet = originalHelper
	})

	if err := auth.Login("gitee.com", "test-token", "test-user"); err != nil {
		t.Fatalf("auth.Login() returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = auth.Logout("gitee.com")
	})

	gitCredentialHelperGet = func() (string, error) {
		return "manager", nil
	}
	gitCredentialFill = func(host string) (config.HostAuth, error) {
		return config.HostAuth{}, errors.New("missing credential")
	}

	err := doctorCmd.RunE(doctorCmd, nil)
	if err == nil {
		t.Fatal("doctorCmd.RunE() error = nil, want non-nil when git credential lookup fails")
	}
}
