package cmd

import (
	"github.com/heronalps/STOIC/client"
	"github.com/spf13/cobra"
)

// inquisitorCmd represents the inquisitor command
var (
	interval        int
	winSizeInterval int

	inquisitorCmd = &cobra.Command{
		Use:   "inquisitor",
		Short: "Inquisitor keeps probing Nautilus for deployment time of runtimes",
		Long:  `Inquisitor keeps probing Nautilus for deployment tiem of runtimes`,
		Run: func(cmd *cobra.Command, args []string) {
			client.StartInquisitor(winSizeInterval, app, interval)
		},
	}
)

func init() {
	runCmd.AddCommand(inquisitorCmd)
	inquisitorCmd.Flags().IntVarP(&interval, "interval", "i", 60, "The interval of inquire deployment time on Nautilus")
	inquisitorCmd.Flags().StringVarP(&app, "app", "a", "image-clf-inf", "The application of deployment")
	inquisitorCmd.Flags().IntVarP(&winSizeInterval, "winInt", "w", 100, "The number of deployments between every two window size calibrations")
}
