package cmd

import "testing"

func TestUserCommand(t *testing.T) {
	if userCmd.Use != "user" {
		t.Errorf("expected use 'user', got %s", userCmd.Use)
	}
	if len(userCmd.Commands()) != 9 {
		t.Errorf("expected 9 subcommands, got %d", len(userCmd.Commands()))
	}
}

func TestUserKeyCommand(t *testing.T) {
	if userKeyCmd.Use != "key" {
		t.Errorf("expected use 'key', got %s", userKeyCmd.Use)
	}
	if len(userKeyCmd.Commands()) != 4 {
		t.Errorf("expected 4 subcommands, got %d", len(userKeyCmd.Commands()))
	}
}

func TestUserCommandFlags(t *testing.T) {
	if userUpdateCmd.Flags().Lookup("name") == nil {
		t.Error("expected --name flag")
	}
	if userFollowersCmd.Flags().Lookup("page") == nil {
		t.Error("expected --page flag")
	}
	if userKeyCreateCmd.Flags().Lookup("title") == nil {
		t.Error("expected --title flag")
	}
	if userKeyCreateCmd.Flags().Lookup("key") == nil {
		t.Error("expected --key flag")
	}
	if userNamespaceViewCmd.Flags().Lookup("path") == nil {
		t.Error("expected --path flag")
	}
}
