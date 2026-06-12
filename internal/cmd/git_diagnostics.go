package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
)

type credentialTarget struct {
	Target string `json:"target"`
	User   string `json:"user,omitempty"`
}

type authDoctorReport struct {
	Host                    string             `json:"host"`
	StoredAuth              config.HostAuth    `json:"stored_auth"`
	StoredAuthErr           string             `json:"stored_auth_error,omitempty"`
	CredentialHelper        string             `json:"credential_helper,omitempty"`
	CredentialHelperErr     string             `json:"credential_helper_error,omitempty"`
	GitCredential           config.HostAuth    `json:"git_credential"`
	GitCredentialErr        string             `json:"git_credential_error,omitempty"`
	CredentialTargets       []credentialTarget `json:"credential_targets,omitempty"`
	CredentialTargetsErr    string             `json:"credential_targets_error,omitempty"`
	CurrentDirectoryIsRepo  bool               `json:"current_directory_is_repo"`
	CurrentDirectoryRepoErr string             `json:"current_directory_repo_error,omitempty"`
	OriginURL               string             `json:"origin_url,omitempty"`
	OriginURLErr            string             `json:"origin_url_error,omitempty"`
	OriginAccessOK          bool               `json:"origin_access_ok"`
	OriginAccessErr         string             `json:"origin_access_error,omitempty"`
	AuthUserMatchesGit      bool               `json:"auth_user_matches_git"`
	ReadyForHTTPSGit        bool               `json:"ready_for_https_git"`
}

var gitCredentialApprove = func(host, username, token string) error {
	payload := strings.Join([]string{
		"protocol=https",
		fmt.Sprintf("host=%s", host),
		fmt.Sprintf("username=%s", username),
		fmt.Sprintf("password=%s", token),
		"",
	}, "\n")

	cmd := exec.Command("git", "credential", "approve")
	cmd.Stdin = strings.NewReader(payload)
	output, err := cmd.CombinedOutput()
	if err != nil {
		trimmed := strings.TrimSpace(string(output))
		if trimmed == "" {
			return fmt.Errorf("git credential approve failed: %w", err)
		}
		return fmt.Errorf("git credential approve failed: %w: %s", err, trimmed)
	}
	return nil
}

var gitCredentialFill = func(host string) (config.HostAuth, error) {
	output, err := runGitCommand(
		strings.NewReader(strings.Join([]string{
			"protocol=https",
			fmt.Sprintf("host=%s", host),
			"",
		}, "\n")),
		"credential",
		"fill",
	)
	if err != nil {
		return config.HostAuth{}, fmt.Errorf("git credential fill failed: %w", err)
	}

	result := config.HostAuth{}
	for _, line := range strings.Split(output, "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch key {
		case "username":
			result.User = value
		case "password":
			result.Token = value
		}
	}

	if result.User == "" || result.Token == "" {
		return config.HostAuth{}, fmt.Errorf("git credential fill returned incomplete credentials")
	}

	return result, nil
}

var gitCredentialReject = func(host, username string) error {
	lines := []string{
		"protocol=https",
		fmt.Sprintf("host=%s", host),
	}
	if username != "" {
		lines = append(lines, fmt.Sprintf("username=%s", username))
	}
	lines = append(lines, "")

	_, err := runGitCommand(strings.NewReader(strings.Join(lines, "\n")), "credential", "reject")
	return err
}

var gitCredentialHelperGet = func() (string, error) {
	output, err := runGitCommand(nil, "config", "--get", "credential.helper")
	if err != nil {
		return "", err
	}
	return output, nil
}

var gitCredentialTargetsList = func(host string) ([]credentialTarget, error) {
	if runtime.GOOS != "windows" {
		return nil, nil
	}

	cmd := exec.Command("cmdkey", "/list")
	output, err := cmd.CombinedOutput()
	trimmed := strings.TrimSpace(string(output))
	if err != nil {
		if trimmed == "" {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %s", err, trimmed)
	}

	var targets []credentialTarget
	var current *credentialTarget
	for _, rawLine := range strings.Split(string(output), "\n") {
		line := strings.TrimSpace(strings.ReplaceAll(rawLine, "\r", ""))
		if strings.HasPrefix(line, "Target:") {
			target := strings.TrimSpace(strings.TrimPrefix(line, "Target:"))
			if strings.Contains(target, host) {
				targets = append(targets, credentialTarget{Target: target})
				current = &targets[len(targets)-1]
			} else {
				current = nil
			}
			continue
		}
		if current != nil && strings.HasPrefix(line, "User:") {
			current.User = strings.TrimSpace(strings.TrimPrefix(line, "User:"))
		}
	}

	return targets, nil
}

var gitCredentialTargetDelete = func(target string) error {
	if runtime.GOOS != "windows" {
		return nil
	}
	cmd := exec.Command("cmdkey", "/delete:"+target)
	output, err := cmd.CombinedOutput()
	trimmed := strings.TrimSpace(string(output))
	if err != nil {
		if trimmed == "" {
			return err
		}
		return fmt.Errorf("%w: %s", err, trimmed)
	}
	return nil
}

var gitLsRemote = func(target string) error {
	_, err := runGitCommand(nil, "ls-remote", target)
	return err
}

var gitRemoteGetURL = func(name string) (string, error) {
	return runGitCommand(nil, "remote", "get-url", name)
}

var gitRemoteSetURL = func(name, target string) error {
	_, err := runGitCommand(nil, "remote", "set-url", name, target)
	return err
}

var gitRemoteAdd = func(name, target string) error {
	_, err := runGitCommand(nil, "remote", "add", name, target)
	return err
}

var gitCurrentBranch = func() (string, error) {
	return runGitCommand(nil, "branch", "--show-current")
}

var gitPushUpstream = func(remote, branch string) error {
	_, err := runGitCommand(nil, "push", "-u", remote, branch)
	return err
}

var gitIsInsideWorkTree = func() (bool, error) {
	output, err := runGitCommand(nil, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		if strings.Contains(err.Error(), "not a git repository") {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(output) == "true", nil
}

func (r authDoctorReport) ToJSON() (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func runAuthDoctor(host string) authDoctorReport {
	report := authDoctorReport{Host: host}

	if stored, err := auth.GetAuth(host); err != nil {
		report.StoredAuthErr = err.Error()
	} else {
		report.StoredAuth = stored
	}

	if helper, err := gitCredentialHelperGet(); err != nil {
		report.CredentialHelperErr = err.Error()
	} else {
		report.CredentialHelper = helper
	}

	if targets, err := gitCredentialTargetsList(host); err != nil {
		report.CredentialTargetsErr = err.Error()
	} else {
		report.CredentialTargets = targets
	}

	if credential, err := gitCredentialFill(host); err != nil {
		report.GitCredentialErr = err.Error()
	} else {
		report.GitCredential = credential
	}

	if inside, err := gitIsInsideWorkTree(); err != nil {
		report.CurrentDirectoryRepoErr = err.Error()
	} else {
		report.CurrentDirectoryIsRepo = inside
	}

	if report.CurrentDirectoryIsRepo {
		if origin, err := gitRemoteGetURL("origin"); err != nil {
			report.OriginURLErr = err.Error()
		} else {
			report.OriginURL = origin
			if err := gitLsRemote(origin); err != nil {
				report.OriginAccessErr = err.Error()
			} else {
				report.OriginAccessOK = true
			}
		}
	}

	report.AuthUserMatchesGit = report.StoredAuthErr == "" &&
		report.GitCredentialErr == "" &&
		report.StoredAuth.User != "" &&
		report.StoredAuth.User == report.GitCredential.User

	report.ReadyForHTTPSGit =
		report.StoredAuthErr == "" &&
			report.CredentialHelperErr == "" &&
			report.GitCredentialErr == "" &&
			report.AuthUserMatchesGit

	return report
}

func runGitCommand(stdin *strings.Reader, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if stdin != nil {
		cmd.Stdin = stdin
	}
	output, err := cmd.CombinedOutput()
	trimmed := strings.TrimSpace(string(output))
	if err != nil {
		if trimmed == "" {
			return "", err
		}
		return "", fmt.Errorf("%w: %s", err, trimmed)
	}
	return trimmed, nil
}
