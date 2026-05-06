package cmd

import (
	"bytes"
	"testing"

	"github.com/ricsy/gt/pkg/util"
)

func TestPRCommandHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	prCmd.SetOut(buf)
	prCmd.SetArgs([]string{"--help"})

	err := prCmd.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestPRListCmdFlags(t *testing.T) {
	if prListCmd.Flags().Lookup("repo") == nil {
		t.Error("prListCmd should have --repo flag")
	}
	if prListCmd.Flags().Lookup("state") == nil {
		t.Error("prListCmd should have --state flag")
	}
}

func TestPRViewCmdArgs(t *testing.T) {
	if prViewCmd.Args == nil {
		t.Error("prViewCmd should have Args defined")
	}
}

func TestPRCreateCmdFlags(t *testing.T) {
	if prCreateCmd.Flags().Lookup("title") == nil {
		t.Error("prCreateCmd should have --title flag")
	}
	if prCreateCmd.Flags().Lookup("body") == nil {
		t.Error("prCreateCmd should have --body flag")
	}
	if prCreateCmd.Flags().Lookup("head") == nil {
		t.Error("prCreateCmd should have --head flag")
	}
	if prCreateCmd.Flags().Lookup("base") == nil {
		t.Error("prCreateCmd should have --base flag")
	}
}

func TestPRMergeCmd(t *testing.T) {
	if prMergeCmd.Use != "merge <number>" {
		t.Errorf("prMergeCmd.Use = %s, want 'merge <number>'", prMergeCmd.Use)
	}
}

func TestPRCloseCmd(t *testing.T) {
	if prCloseCmd.Use != "close <number>" {
		t.Errorf("prCloseCmd.Use = %s, want 'close <number>'", prCloseCmd.Use)
	}
}

func TestPRCommentCmdFlags(t *testing.T) {
	if prCommentCmd.Flags().Lookup("body") == nil {
		t.Error("prCommentCmd should have --body flag")
	}
}

func TestSplitOwnerRepo(t *testing.T) {
	owner, repo := util.SplitOwnerRepo("owner/repo")
	if owner != "owner" || repo != "repo" {
		t.Errorf("SplitOwnerRepo = (%s, %s), want (owner, repo)", owner, repo)
	}

	owner, repo = util.SplitOwnerRepo("onlyname")
	if repo != "onlyname" || owner != "" {
		t.Errorf("SplitOwnerRepo single = (%s, %s), want (, onlyname)", owner, repo)
	}
}
