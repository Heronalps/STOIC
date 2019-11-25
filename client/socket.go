// Client socket listening to Server socket for relaying the wtb task to Nautilus

package client

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

/*
SocketClient listens to the task request from the server
*/
func SocketClient(port int) {
	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}

	defer listen.Close()

	log.Printf("Begin listen port: %d", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {

	defer conn.Close()
	var (
		reader = bufio.NewReader(conn)
		buf    = make([]byte, 1024)
	)

ILOOP:
	for {
		n, err := reader.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			dataSlice := strings.Split(data, " ")
			runtime := dataSlice[0]
			imageNum, err := strconv.Atoi(dataSlice[1])
			if err != nil {
				log.Println(err.Error())
				continue
			}
			Request(runtime, imageNum)
			if isTransportOver(data) {
				break ILOOP
			}
		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}
	}
}

func isTransportOver(data string) (over bool) {
	StopCharacter := " "
	over = strings.HasSuffix(data, StopCharacter)
	return
}
