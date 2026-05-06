package cmd

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api"
	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

func newOrgCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List organizations for the current user",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := auth.GetToken(config.DefaultHost)
			if err != nil {
				return fmt.Errorf("authentication required: %w", err)
			}

			client := api.NewClient(config.DefaultHost, token)
			orgs, err := client.ListOrgs()
			if err != nil {
				return err
			}

			for _, org := range orgs {
				cmd.Printf("%s - %s\n", org.Login, org.Name)
			}
			return nil
		},
	}

	viewCmd := &cobra.Command{
		Use:   "view <org>",
		Short: "View organization details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orgLogin := args[0]

			token, err := auth.GetToken(config.DefaultHost)
			if err != nil {
				return fmt.Errorf("authentication required: %w", err)
			}

			client := api.NewClient(config.DefaultHost, token)
			org, err := client.GetOrg(orgLogin)
			if err != nil {
				return err
			}

			cmd.Printf("Login: %s\nName: %s\nLocation: %s\nBlog: %s\nEmail: %s\nURL: %s\n",
				org.Login, org.Name, org.Location, org.Blog, org.Email, org.HtmlUrl)
			return nil
		},
	}

	cmd := &cobra.Command{
		Use:   "org",
		Short: "Manage organizations",
	}
	cmd.AddCommand(listCmd, viewCmd)

	return cmd
}

func init() {
	rootCmd.AddCommand(newOrgCmd())
}
