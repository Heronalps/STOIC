package cmd

import (
	"fmt"
	"time"

	"github.com/heronalps/STOIC/client"
	"github.com/spf13/cobra"
)

// inquisitorCmd represents the inquisitor command
var (
	interval      int
	inquisitorCmd = &cobra.Command{
		Use:   "inquisitor",
		Short: "Inquisitor keeps probing Nautilus for deployment time of runtimes",
		Long:  `Inquisitor keeps probing Nautilus for deployment tiem of runtimes`,
		Run: func(cmd *cobra.Command, args []string) {
			for {
				client.UpdateDeploymentTimeTable()
				fmt.Println("Waiting for next round ...")
				time.Sleep(time.Second * time.Duration(interval))
			}
		},
	}
)

func init() {
	runCmd.AddCommand(inquisitorCmd)
	inquisitorCmd.Flags().IntVarP(&interval, "interval", "i", 600, "The interval of inquire deployment time on Nautilus")
}
