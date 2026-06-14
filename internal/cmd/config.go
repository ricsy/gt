package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage gt configuration",
		Long:  `Manage gt configuration settings.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}
			key := args[0]
			switch key {
			case "default_repo":
				fmt.Println(cfg.DefaultRepo)
			case "default_owner":
				fmt.Println(cfg.DefaultOwner)
			case "default_host":
				fmt.Println(cfg.DefaultHost)
			case "repo_scope_mode":
				fmt.Println(cfg.RepoScopeMode)
			case "repo_scope_namespace":
				fmt.Println(cfg.RepoScopeNamespace)
			default:
				return fmt.Errorf("unknown key: %s", key)
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}
			key, value := args[0], args[1]
			switch key {
			case "default_repo":
				cfg.DefaultRepo = value
			case "default_owner":
				cfg.DefaultOwner = value
			case "default_host":
				cfg.DefaultHost = value
			case "repo_scope_mode":
				cfg.RepoScopeMode = value
			case "repo_scope_namespace":
				cfg.RepoScopeNamespace = value
			default:
				return fmt.Errorf("unknown key: %s", key)
			}
			return config.SaveConfig(cfg)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "default_repo: %s\n", cfg.DefaultRepo)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "default_owner: %s\n", cfg.DefaultOwner)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "default_host: %s\n", cfg.DefaultHost)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "repo_scope_mode: %s\n", cfg.RepoScopeMode)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "repo_scope_namespace: %s\n", cfg.RepoScopeNamespace)
			return nil
		},
	})

	return cmd
}

func init() {
	rootCmd.AddCommand(newConfigCmd())
}
