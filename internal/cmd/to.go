package cmd

import (
	"fmt"
	"strings"

	"github.com/pkg/browser"
	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

var toOpts struct {
	Repo bool
}

var resolveCurrentUser = auth.CurrentUser
var openBrowserURL = browser.OpenURL
var verifyToRepoExists = func(owner, repo string) error {
	client := newCommandAPIClient(resolveCommandHost(), "")
	if _, err := client.GetRepo(owner, repo); err != nil {
		if isToNotFoundError(err) {
			return fmt.Errorf("repository not found: %s/%s", owner, repo)
		}
		return fmt.Errorf("failed to verify repository %s/%s: %w", owner, repo, err)
	}
	return nil
}
var verifyToNamespaceExists = func(namespace string) error {
	client := newCommandAPIClient(resolveCommandHost(), "")
	if _, err := client.GetUser(namespace); err == nil {
		return nil
	} else if !isToNotFoundError(err) {
		return fmt.Errorf("failed to verify user %s: %w", namespace, err)
	}

	if _, err := client.GetOrg(namespace); err == nil {
		return nil
	} else if !isToNotFoundError(err) {
		return fmt.Errorf("failed to verify organization %s: %w", namespace, err)
	}

	return fmt.Errorf("namespace not found: %s", namespace)
}

type toTargetKind string

const (
	toTargetKindRepo      toTargetKind = "repo"
	toTargetKindNamespace toTargetKind = "namespace"
)

type resolvedToTarget struct {
	Kind      toTargetKind
	URL       string
	Namespace string
	Owner     string
	Repo      string
}

var toCmd = &cobra.Command{
	Use:   "to [target]",
	Short: "Open a repository or namespace in the browser",
	Args:  cobra.RangeArgs(0, 1),
	RunE:  toCommand,
}

func init() {
	rootCmd.AddCommand(toCmd)
	toCmd.Flags().BoolVarP(&toOpts.Repo, "repo", "r", false, "Treat a single-segment target as a repository owned by the authenticated user")
}

func toCommand(cmd *cobra.Command, args []string) error {
	target := ""
	if len(args) == 1 {
		target = args[0]
	}

	resolvedTarget, err := resolveToTarget(target, toOpts.Repo)
	if err != nil {
		return err
	}

	if err := validateToTarget(resolvedTarget); err != nil {
		return err
	}

	if err := openBrowserURL(resolvedTarget.URL); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	cmd.Printf("Opened: %s\n", resolvedTarget.URL)
	return nil
}

func resolveToTargetURL(target string, treatAsRepo bool) (string, error) {
	resolvedTarget, err := resolveToTarget(target, treatAsRepo)
	if err != nil {
		return "", err
	}
	return resolvedTarget.URL, nil
}

func resolveToTarget(target string, treatAsRepo bool) (*resolvedToTarget, error) {
	host := resolveCommandHost()
	baseURL := config.WebURL(host)
	target = strings.TrimSpace(target)

	if target == "" {
		owner, repoName, err := resolveRepoFlag("")
		if err != nil {
			return nil, err
		}
		return &resolvedToTarget{
			Kind:  toTargetKindRepo,
			URL:   fmt.Sprintf("%s/%s/%s", baseURL, owner, repoName),
			Owner: owner,
			Repo:  repoName,
		}, nil
	}

	if strings.Contains(target, "/") {
		owner, repoName, err := ResolveRepo(target)
		if err != nil {
			return nil, err
		}
		return &resolvedToTarget{
			Kind:  toTargetKindRepo,
			URL:   fmt.Sprintf("%s/%s/%s", baseURL, owner, repoName),
			Owner: owner,
			Repo:  repoName,
		}, nil
	}

	if treatAsRepo {
		currentUser, err := resolveCurrentUser(host)
		if err != nil {
			return nil, fmt.Errorf("authenticated user is required to resolve repo target %q: %w", target, err)
		}
		return &resolvedToTarget{
			Kind:  toTargetKindRepo,
			URL:   fmt.Sprintf("%s/%s/%s", baseURL, currentUser, target),
			Owner: currentUser,
			Repo:  target,
		}, nil
	}

	return &resolvedToTarget{
		Kind:      toTargetKindNamespace,
		URL:       fmt.Sprintf("%s/%s", baseURL, target),
		Namespace: target,
	}, nil
}

func validateToTarget(target *resolvedToTarget) error {
	if target == nil {
		return fmt.Errorf("target is required")
	}

	switch target.Kind {
	case toTargetKindRepo:
		return verifyToRepoExists(target.Owner, target.Repo)
	case toTargetKindNamespace:
		return verifyToNamespaceExists(target.Namespace)
	default:
		return fmt.Errorf("unsupported target kind: %s", target.Kind)
	}
}

func isToNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "unexpected status code: 404")
}
