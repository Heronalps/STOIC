package cmd

import (
	"github.com/heronalps/STOIC/client"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var (
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "setup 3 datapoints in Processing Time tables for regression",
		Long:  `setup 3 datapoints in Processing Time tables for regression`,
		Run: func(cmd *cobra.Command, args []string) {
			client.SetupRegression(app, version)
		},
	}
)

func init() {
	dbCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringVarP(&app, "app", "a", "image-clf-inf", "The application name")
	setupCmd.Flags().StringVarP(&version, "version", "v", "1.0", "The version of application")
}
