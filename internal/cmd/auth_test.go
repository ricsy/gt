package cmd

import (
	"bytes"
	"testing"

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

	expectedCmds := []string{"login", "logout", "status", "token"}
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
