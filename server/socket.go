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

	StopCharacter := " "
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer conn.Close()

	message := runtime + " " + strconv.Itoa(imageNum)
	conn.Write([]byte(message))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", message)
}
