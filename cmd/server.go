package cmd

import (
	"fmt"
	"os"

	"github.com/heronalps/STOIC/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var (
	ip             string
	presetImageNum int
	imageNum       int
	totalBatch     int
	preset         bool
	index          int
	zipPath        string

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run STOIC Server",
		Long:  `Run STOIC Socket Server`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				totalImage int
				batches    int = 1
				totalTime  float64
			)

			for {
				// Listens for IoT Camera Traps

				// Server is triggered by Camera traps to generate a batch
				fmt.Printf("Batches : %d \n", batches)
				zipPath, imageNum = server.GenerateBatch(presetImageNum, batches)

				elapsed := server.SocketServer(ip, port, zipPath, imageNum)
				if elapsed == 0.0 {
					fmt.Println("The task was not executed...")
					fmt.Println("continue...")
					fmt.Println("==================Next Batch===========================")
					continue
				}
				totalImage += imageNum
				totalTime += elapsed
				batches++
				fmt.Printf("After Batches : %d", batches)
				fmt.Printf("%d images has been inferenced...\n", totalImage)
				fmt.Printf("%d batches has been processed...\n", batches)
				fmt.Printf("%f seconds for this batch...\n", elapsed)
				fmt.Printf("%f seconds has elapsed...\n", totalTime)
				// When batch is 0, it marks infinitely execution
				if totalBatch != 0 && totalBatch == batches {
					os.Exit(0)
				}
				fmt.Println("==================Next Batch===========================")
			}
		},
	}
)

func init() {
	runCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&ip, "ip", "127.0.0.1", "The IP address of client")
	serverCmd.Flags().IntVarP(&presetImageNum, "image", "n", 0, "Image number in one batch")
	serverCmd.Flags().IntVarP(&totalBatch, "batch", "b", 0, "Total batches of image")
	serverCmd.Flags().BoolVarP(&preset, "preset", "s", false, "If the batch size is preset")
}
