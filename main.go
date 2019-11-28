package main

import (
	"flag"
	"fmt"

	"github.com/heronalps/STOIC/client"
	"github.com/heronalps/STOIC/server"
)

func main() {
	// Command-line args: client/server, ip, port
	node := flag.String("node", "client", "client/server")
	ip := flag.String("ip", "127.0.0.1", "IP address of client")
	port := flag.Int("port", 5001, "Port of client")
	runtime := flag.String("runtime", "", "Runtime of WTB tasks: euca/cpu/gpu1/gpu2")
	flag.Parse()

	if *node == "client" {
		client.SocketClient(*port)
	} else if *node == "server" {
		var (
			imageNum  int
			batches   int
			totalTime float64
		)

		for {
			images, elapsed := server.Schedule(*ip, *port, *runtime)
			if elapsed == 0.0 {
				fmt.Println("The task was not executed...")
				fmt.Println("continue...")
				fmt.Println("==================Next Batch===========================")
				continue
			}
			imageNum += images
			totalTime += elapsed
			batches++
			fmt.Printf("%d images has been inferenced...\n", imageNum)
			fmt.Printf("%d batches has been processed...\n", batches)
			fmt.Printf("%f seconds has elapsed...\n", totalTime)
			fmt.Println("==================Next Batch===========================")
		}
	}
}
