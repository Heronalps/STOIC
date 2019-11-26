package server

import (
	"log"
	"net"
	"strconv"
	"strings"
)

/*
SocketServer sends wtb task request to the client socket
*/
func SocketServer(ip string, port int, runtime string, imageNum int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	StopCharacter := " #"
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer conn.Close()
	// TODO: Add lock to avoid race condition on kubeless function on the client
	// Server socket sends message only if it obtains the lock, otherwise being blocked
	message := runtime + " " + strconv.Itoa(imageNum)
	conn.Write([]byte(message))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", message)
}
