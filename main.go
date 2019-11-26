package main

import (
	"flag"
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
		for {
			server.Schedule(*ip, *port)
			time.Sleep(3 * time.Minute)
		}
	}
}
