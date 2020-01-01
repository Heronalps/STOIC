package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/heronalps/STOIC/client"
)

/*
SocketServer sends wtb task request to the client socket
*/
func SocketServer(ip string, port int, imageNum int) float64 {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	StopCharacter := " #"
	if err != nil {
		log.Fatalln(err)
		return 0
	}

	defer conn.Close()
	// TODO: Add lock to avoid race condition on kubeless function on the client
	// Server socket sends message only if it obtains the lock, otherwise being blocked
	message := fmt.Sprintf("%s", strconv.Itoa(imageNum))
	conn.Write([]byte(message))
	conn.Write([]byte(StopCharacter))
	log.Printf("Sent: %s \n", message)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Println("Received output from client...")
	elapsed := client.ParseElapsed(buff[:n])
	if err != nil {
		log.Println(err.Error())
	}
	return elapsed
}
