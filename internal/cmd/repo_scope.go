package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
)

const (
	repoScopeModeNone     = ""
	repoScopeModePersonal = "personal"
	repoScopeModeOrg      = "org"
)

type repoScope struct {
	Mode      string
	Namespace string
}

func loadRepoScope() (repoScope, error) {
	if envMode := strings.TrimSpace(os.Getenv("GT_REPO_SCOPE_MODE")); envMode != "" {
		scope := repoScope{
			Mode:      envMode,
			Namespace: strings.TrimSpace(os.Getenv("GT_REPO_SCOPE_NAMESPACE")),
		}
		switch scope.Mode {
		case repoScopeModePersonal:
			scope.Namespace = ""
			return scope, nil
		case repoScopeModeOrg:
			if scope.Namespace == "" {
				return repoScope{}, fmt.Errorf("repo scope mode org requires GT_REPO_SCOPE_NAMESPACE")
			}
			return scope, nil
		default:
			return repoScope{}, fmt.Errorf("invalid repo scope mode: %s", scope.Mode)
		}
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return repoScope{}, err
	}

	scope := repoScope{
		Mode:      strings.TrimSpace(cfg.RepoScopeMode),
		Namespace: strings.TrimSpace(cfg.RepoScopeNamespace),
	}

	switch scope.Mode {
	case repoScopeModeNone, repoScopeModePersonal:
		scope.Namespace = ""
		return scope, nil
	case repoScopeModeOrg:
		if scope.Namespace == "" {
			return repoScope{}, fmt.Errorf("repo scope mode org requires repo_scope_namespace")
		}
		return scope, nil
	default:
		return repoScope{}, fmt.Errorf("invalid repo scope mode: %s", scope.Mode)
	}
}

func resolvePersonalRepoOwner() (string, error) {
	authInfo, err := auth.GetAuth(resolveCommandHost())
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(authInfo.User) == "" {
		return "", fmt.Errorf("authenticated user is empty")
	}
	return strings.TrimSpace(authInfo.User), nil
}

func resolveScopedRepoOwner() (string, error) {
	scope, err := loadRepoScope()
	if err != nil {
		return "", err
	}

	switch scope.Mode {
	case repoScopeModeNone:
		return "", nil
	case repoScopeModePersonal:
		return resolvePersonalRepoOwner()
	case repoScopeModeOrg:
		return scope.Namespace, nil
	default:
		return "", fmt.Errorf("invalid repo scope mode: %s", scope.Mode)
	}
}

func resolveDefaultRepoOwner() (string, error) {
	scopedOwner, err := resolveScopedRepoOwner()
	if err != nil {
		return "", err
	}
	if scopedOwner != "" {
		return scopedOwner, nil
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(cfg.DefaultOwner), nil
}

func enforceRepoScopeOwner(owner string) error {
	scopedOwner, err := resolveScopedRepoOwner()
	if err != nil {
		return err
	}
	if scopedOwner == "" || owner == "" || owner == scopedOwner {
		return nil
	}

	return fmt.Errorf("repo scope is locked to %s; requested owner %s is outside the active scope", scopedOwner, owner)
}

func resolveRepoCreateNamespace(explicitNamespace string) (string, error) {
	namespace := strings.TrimSpace(explicitNamespace)
	scope, err := loadRepoScope()
	if err != nil {
		return "", err
	}

	switch scope.Mode {
	case repoScopeModeNone:
		return namespace, nil
	case repoScopeModePersonal:
		if namespace != "" {
			personalOwner, err := resolvePersonalRepoOwner()
			if err != nil {
				return "", err
			}
			if namespace != personalOwner {
				return "", fmt.Errorf("repo scope is personal; namespace %s is outside the active scope", namespace)
			}
		}
		return namespace, nil
	case repoScopeModeOrg:
		if namespace != "" && namespace != scope.Namespace {
			return "", fmt.Errorf("repo scope is locked to org %s; namespace %s is outside the active scope", scope.Namespace, namespace)
		}
		return scope.Namespace, nil
	default:
		return "", fmt.Errorf("invalid repo scope mode: %s", scope.Mode)
	}
}

func describeRepoScope(scope repoScope) string {
	switch scope.Mode {
	case repoScopeModePersonal:
		return repoScopeModePersonal
	case repoScopeModeOrg:
		return fmt.Sprintf("%s:%s", repoScopeModeOrg, scope.Namespace)
	default:
		return "none"
	}
}
