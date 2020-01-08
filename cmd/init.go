package cmd

import (
	"github.com/heronalps/STOIC/client"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize database and all necessary tables",
	Long:  `Initialize database and all necessary tables`,
	Run: func(cmd *cobra.Command, args []string) {
		client.InitDB()
	},
}

func init() {
	dbCmd.AddCommand(initCmd)
}
