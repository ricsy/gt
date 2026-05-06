package cmd

import (
	"github.com/spf13/cobra"
)

func newOrgCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List organizations for the current user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient()
			if err != nil {
				return err
			}

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

			client, err := getClient()
			if err != nil {
				return err
			}

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
