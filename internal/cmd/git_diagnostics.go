package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
)

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

var gitCredentialHelperGet = func() (string, error) {
	output, err := runGitCommand(nil, "config", "--get", "credential.helper")
	if err != nil {
		return "", err
	}
	return output, nil
}

var gitLsRemote = func(target string) error {
	_, err := runGitCommand(nil, "ls-remote", target)
	return err
}

var gitRemoteGetURL = func(name string) (string, error) {
	return runGitCommand(nil, "remote", "get-url", name)
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

type authDoctorReport struct {
	Host                string
	StoredAuth          config.HostAuth
	StoredAuthErr       error
	CredentialHelper    string
	CredentialHelperErr error
	GitCredential       config.HostAuth
	GitCredentialErr    error
}

func (r authDoctorReport) ReadyForHTTPSGit() bool {
	return r.StoredAuthErr == nil && r.CredentialHelperErr == nil && r.GitCredentialErr == nil
}

func runAuthDoctor(host string) authDoctorReport {
	report := authDoctorReport{Host: host}
	report.StoredAuth, report.StoredAuthErr = auth.GetAuth(host)
	report.CredentialHelper, report.CredentialHelperErr = gitCredentialHelperGet()
	report.GitCredential, report.GitCredentialErr = gitCredentialFill(host)
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
