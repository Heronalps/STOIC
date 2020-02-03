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
	batch          int
	preset         bool
	index          int = 222
	// randomSize     []int = []int{33, 20, 59, 10, 75, 17, 37, 132, 26, 49, 10, 93,
	// 	20, 47, 66, 62, 23, 35, 63, 18, 132, 24, 75, 22}

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run STOIC Server",
		Long:  `Run STOIC Socket Server`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				totalImage int
				batches    int
				totalTime  float64
			)

			for {
				// image flag has high precedence than preset
				if presetImageNum == 0 {
					if preset {
						fmt.Printf("index : %v ..\n", index)
						imageNum = server.Workload[index]
						index++
						// if index >= len(server.Workload)/30 {
						// 	os.Exit(0)
						// }
					} else {
						imageNum = server.ImageCache()
					}
				} else {
					imageNum = presetImageNum
				}

				elapsed := server.SocketServer(ip, port, imageNum)
				if elapsed == 0.0 {
					fmt.Println("The task was not executed...")
					fmt.Println("continue...")
					fmt.Println("==================Next Batch===========================")
					continue
				}
				totalImage += imageNum
				totalTime += elapsed
				batches++
				fmt.Printf("%d images has been inferenced...\n", totalImage)
				fmt.Printf("%d batches has been processed...\n", batches)
				fmt.Printf("%f seconds for this batch...\n", elapsed)
				fmt.Printf("%f seconds has elapsed...\n", totalTime)
				// When batch is 0, it marks infinitely execution
				if batch != 0 && batch == batches {
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
	serverCmd.Flags().IntVarP(&batch, "batch", "b", 0, "Batches of image")
	serverCmd.Flags().BoolVarP(&preset, "preset", "s", false, "If the batch size is preset")
}
