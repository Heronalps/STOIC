// Client socket listening to Server socket for relaying the wtb task to Nautilus

package client

import (
	"bufio"
	"bytes"
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
		reader   = bufio.NewReader(conn)
		buf      = make([]byte, 1024)
		runtime  string
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
				runtime = dataSlice[0]
				imageNum, err = strconv.Atoi(dataSlice[1])
				// fmt.Printf("runtime: %s \n", runtime)
				// fmt.Printf("imageNum: %d \n", imageNum)
				break ILOOP
			}
		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}
	}
	Request(runtime, imageNum)
}

func isTransportOver(data string) (over bool) {
	StopCharacter := " #"
	over = strings.HasSuffix(data, StopCharacter)
	return
}
