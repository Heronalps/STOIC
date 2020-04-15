package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/heronalps/STOIC/client"
)

/*
BUFFERSIZE is the size of network buffer
*/
const BUFFERSIZE = 1024

/*
SocketServer sends wtb task request to the client socket
*/
func SocketServer(ip string, port int, zipPath string, imageNum int) float64 {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		return 0
	}

	defer conn.Close()

	file, err := os.Open(zipPath)

	if err != nil {
		log.Fatalln(err)
		return 0.0
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
		return 0.0
	}
	// base 10, patch to 10 bytes
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fileImageNum := fillString(strconv.FormatInt(int64(imageNum), 10), 10)
	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))
	conn.Write([]byte(fileImageNum))
	sendBuffer := make([]byte, BUFFERSIZE)

	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}

	recvBuff := make([]byte, 1024)
	n, _ := conn.Read(recvBuff)
	log.Println("Received output from client...")
	fmt.Println(string(recvBuff[:n]))
	elapsed := client.ParseElapsed(recvBuff[:n])
	if err != nil {
		log.Println(err.Error())
	}
	return elapsed
}

func fillString(str string, toLength int) string {
	length := len(str)
	var patch string
	if length < toLength {
		patch = strings.Repeat(":", toLength-length)
	}
	return str + patch
}
