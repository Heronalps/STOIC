package cmd

import (
	"fmt"
	"os"

	"github.com/heronalps/STOIC/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var (
	ip       string
	runtime  string
	imageNum int
	batch    int
	preset   bool

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run STOIC Server",
		Long:  `Run STOIC Socket Server`,
		Run: func(cmd *cobra.Command, args []string) {
			randomSize := []int{33, 20, 59, 10, 75, 17, 37, 132, 26, 49, 10, 93, 20, 47, 66, 62, 23, 35, 63, 18, 132, 24, 75, 22}
			slice := randomSize[21:]
			var (
				totalImage int
				batches    int
				totalTime  float64
			)

			for i := 0; i < len(slice); i++ {
				if preset {
					imageNum = slice[i]
				}
				images, elapsed := server.Schedule(ip, port, runtime, imageNum)
				if elapsed == 0.0 {
					fmt.Println("The task was not executed...")
					fmt.Println("continue...")
					fmt.Println("==================Next Batch===========================")
					continue
				}
				totalImage += images
				totalTime += elapsed
				batches++
				fmt.Printf("%d images has been inferenced...\n", totalImage)
				fmt.Printf("%d batches has been processed...\n", batches)
				fmt.Printf("%f seconds for this batch...\n", elapsed)
				fmt.Printf("%f seconds has elapsed...\n", totalTime)
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
	serverCmd.Flags().StringVarP(&runtime, "runtime", "r", "", "Runtimes of WTB task: edge/cpu/gpu1/gpu2")
	serverCmd.Flags().IntVarP(&imageNum, "image", "n", 0, "Image number in one batch")
	serverCmd.Flags().IntVarP(&batch, "batch", "b", 0, "Batches of image")
	serverCmd.Flags().BoolVarP(&preset, "preset", "s", false, "If the batch size is preset")
}