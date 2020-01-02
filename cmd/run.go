package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var (
	port   int
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run STOIC client / server ",
		Long:  `Run STOIC client / server to schedule tasks`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Please specify subcommand (client/server/inquisitor)..")
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().IntVarP(&port, "port", "p", 5001, "Port of Client")
}
