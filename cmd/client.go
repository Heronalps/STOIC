package cmd

import (
	"github.com/heronalps/STOIC/client"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var (
	clientCmd = &cobra.Command{
		Use:   "client",
		Short: "Run STOIC client",
		Long:  `Run STOIC socket client`,
		Run: func(cmd *cobra.Command, args []string) {
			client.SocketClient(port)
		},
	}
)

func init() {
	runCmd.AddCommand(clientCmd)
}
