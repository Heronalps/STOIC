package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/heronalps/STOIC/client"
	"github.com/heronalps/STOIC/server"
)

func main() {
	// Command-line args: client/server, ip, port
	runtime := flag.String("runtime", "client", "client/server")
	ip := flag.String("ip", "127.0.0.1", "IP address of client")
	port := flag.Int("port", 5001, "Port of client")
	flag.Parse()

	if *runtime == "client" {
		client.SocketClient(*port)
	} else if *runtime == "server" {
		var (
			imageNum int
			batches  int
		)

		for {
			imageNum += server.Schedule(*ip, *port)
			batches++
			fmt.Printf("%d images has been inferenced...\n", imageNum)
			fmt.Printf("%d batches has been processed...\n", batches)
			time.Sleep(2 * time.Second)
		}
	}
}
