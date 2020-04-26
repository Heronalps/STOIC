package cmd

import (
	"github.com/heronalps/STOIC/client"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var (
	runtime    string
	app        string
	version    string
	allRuntime bool
	imageNum   int
	batches    int
	clientCmd  = &cobra.Command{
		Use:   "client",
		Short: "Run STOIC client",
		Long:  `Run STOIC socket client`,
		Run: func(cmd *cobra.Command, args []string) {
			// go client.StartInquisitor(winSizeInterval, inqApp, interval)
			client.SetupRegression(app, version)
			// client.SocketClient(port, runtime, app, version, allRuntime, imageNum, batches)
		},
	}
)

func init() {
	runCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&runtime, "runtime", "r", "", "Runtimes of WTB task: edge/cpu/gpu1/gpu2")
	clientCmd.Flags().StringVar(&app, "app", "image-clf-inf", "The ML application")
	clientCmd.Flags().StringVar(&version, "version", "1.0", "The version of application")
	clientCmd.Flags().BoolVar(&allRuntime, "all", false, "Send request to all runtime for experiment")
	clientCmd.Flags().IntVarP(&imageNum, "image", "i", 0, "Preset image num")
	clientCmd.Flags().IntVarP(&batches, "batches", "b", 0, "Preset batch number")
}
