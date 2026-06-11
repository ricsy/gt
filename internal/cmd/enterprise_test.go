package cmd

import "testing"

func TestEnterpriseCommand(t *testing.T) {
	if enterpriseCmd.Use != "enterprise" {
		t.Errorf("expected use 'enterprise', got %s", enterpriseCmd.Use)
	}
	if len(enterpriseCmd.Commands()) != 5 {
		t.Errorf("expected 5 subcommands, got %d", len(enterpriseCmd.Commands()))
	}
}

func TestEnterpriseMemberCommand(t *testing.T) {
	if enterpriseMemberCmd.Use != "member" {
		t.Errorf("expected use 'member', got %s", enterpriseMemberCmd.Use)
	}
	if len(enterpriseMemberCmd.Commands()) != 6 {
		t.Errorf("expected 6 subcommands, got %d", len(enterpriseMemberCmd.Commands()))
	}
}

func TestEnterpriseCommandFlags(t *testing.T) {
	if enterpriseListCmd.Flags().Lookup("admin") == nil {
		t.Error("expected --admin flag")
	}
	if enterpriseMemberListCmd.Flags().Lookup("role") == nil {
		t.Error("expected --role flag")
	}
	if enterpriseMemberSearchCmd.Flags().Lookup("query-type") == nil {
		t.Error("expected --query-type flag")
	}
	if enterpriseMemberAddCmd.Flags().Lookup("email") == nil {
		t.Error("expected --email flag")
	}
	if enterpriseMemberUpdateCmd.Flags().Lookup("active") == nil {
		t.Error("expected --active flag")
	}
	if enterpriseRepoCmd.Flags().Lookup("direct") == nil {
		t.Error("expected --direct flag")
	}
	if enterprisePRCmd.Flags().Lookup("state") == nil {
		t.Error("expected --state flag")
	}
	if enterprisePRCmd.Flags().Lookup("issue-number") == nil {
		t.Error("expected --issue-number flag")
	}
	if enterprisePRCmd.Flags().Lookup("head") == nil {
		t.Error("expected --head flag")
	}
	if enterprisePRCmd.Flags().Lookup("base") == nil {
		t.Error("expected --base flag")
	}
	if enterprisePRCmd.Flags().Lookup("since") == nil {
		t.Error("expected --since flag")
	}
	if enterprisePRCmd.Flags().Lookup("program-id") == nil {
		t.Error("expected --program-id flag")
	}
	if enterprisePRCmd.Flags().Lookup("milestone-number") == nil {
		t.Error("expected --milestone-number flag")
	}
}
