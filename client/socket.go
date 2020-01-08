// Client socket listening to Server socket for relaying the wtb task to Nautilus

package client

import (
	"bufio"
	"bytes"
	"fmt"
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
func SocketClient(port int, runtime string, app string, version string) {
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
		go handler(conn, runtime, app, version)
	}
}

func handler(conn net.Conn, runtime string, app string, version string) {
	defer conn.Close()
	var (
		reader   = bufio.NewReader(conn)
		writer   = bufio.NewWriter(conn)
		buf      = make([]byte, 1024)
		imageNum int
		data     bytes.Buffer
	)
ILOOP:
	for {
		n, err := reader.Read(buf)
		data.Write(buf[:n])
		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if isTransportOver(data.String()) {
				// fmt.Printf("data :%sEOF \n", data.String())
				dataSlice := strings.Split(data.String(), " ")
				imageNum, err = strconv.Atoi(dataSlice[0])

				// fmt.Printf("runtime: %s \n", runtime)
				// fmt.Printf("imageNum: %d \n", imageNum)
				break ILOOP
			}
		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}
	}
	// TODO : Add lock to avoid race condition on kubeless function
	// It requires checking kubeless process and
	// write back to server socket if the kubeless function is available

	output := Schedule(runtime, imageNum, app, version)

	writer.Write(output)
	writer.Flush()
	fmt.Println("Sent output to server...")
}

func isTransportOver(data string) (over bool) {
	StopCharacter := " #"
	over = strings.HasSuffix(data, StopCharacter)
	return
}
