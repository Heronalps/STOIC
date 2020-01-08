package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "All commands related to DB operation",
	Long:  `All commands related to DB operation`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify subcommand (init)..")
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
