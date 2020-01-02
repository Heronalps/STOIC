package cmd

import (
	"github.com/heronalps/STOIC/client"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var (
	runtime   string
	app       string
	version   string
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Run STOIC client",
		Long:  `Run STOIC socket client`,
		Run: func(cmd *cobra.Command, args []string) {
			client.SetupRegression(app, version)
			client.SocketClient(port, runtime, app, version)
		},
	}
)

func init() {
	runCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&runtime, "runtime", "r", "", "Runtimes of WTB task: edge/cpu/gpu1/gpu2")
	clientCmd.Flags().StringVar(&app, "app", "wtb", "The ML application")
	clientCmd.Flags().StringVar(&version, "version", "1.0", "The version of application")
}
