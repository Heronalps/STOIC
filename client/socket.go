// Client socket listening to Server socket for relaying the wtb task to Nautilus

package client

import "fmt"

/*
BUFFERSIZE is the buffer size of the TCP socket
*/
const BUFFERSIZE = 1024

/*
SocketClient listens to the task request from the server
*/
func SocketClient(port int, runtime string, app string, version string, all bool, presetImageNum int, batches int, numThread int) {
	var (
		batch int = 1
	)
	for {
		zipPath, imageNum := GenerateBatch(presetImageNum, batch)
		output := Schedule(runtime, imageNum, zipPath, app, version, all, numThread)

		fmt.Printf("Batch %d has been processed...\n", batch)
		fmt.Printf("Output: %v \n", string(output))
		batch++
		if batch == batches {
			break
		}
	}
	fmt.Println("All batches are done processing...")
}

// func SocketClient(port int, runtime string, app string, version string, all bool) {
// 	fmt.Printf("Window Size : %v..\n", windowSizes)
// 	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
// 	if err != nil {
// 		log.Fatalf("Socket listen port %d failed,%s", port, err)
// 		os.Exit(1)
// 	}

// 	defer listen.Close()
// 	log.Printf("Begin listen port: %d", port)
// 	for {
// 		conn, err := listen.Accept()
// 		if err != nil {
// 			log.Fatalln(err)
// 			continue
// 		}
// 		go handler(conn, runtime, app, version, all)
// 	}
// }

// func handler(conn net.Conn, runtime string, app string, version string, all bool) {
// 	fmt.Println("Set up connection...")
// 	defer conn.Close()

// 	writer := bufio.NewWriter(conn)

// 	bufferFileName := make([]byte, 64)
// 	bufferFileSize := make([]byte, 10)
// 	bufferImageNum := make([]byte, 10)

// 	conn.Read(bufferFileSize)
// 	// base 10, bitsize 64 => int64
// 	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

// 	conn.Read(bufferFileName)
// 	fileName := strings.Trim(string(bufferFileName), ":")

// 	newFile, err := os.Create(fileName)

// 	if err != nil {
// 		panic(err)
// 	}

// 	conn.Read(bufferImageNum)
// 	imageNum64, _ := strconv.ParseInt(strings.Trim(string(bufferImageNum), ":"), 10, 64)
// 	imageNum := int(imageNum64)

// 	fmt.Printf("Image Num: %d \n", imageNum)

// 	defer newFile.Close()
// 	var receivedBytes int64
// 	fmt.Println("Start transmitting zip archive..")
// 	for {
// 		if (fileSize - receivedBytes) < BUFFERSIZE {
// 			// Copy the last piece of the file from the connection
// 			io.CopyN(newFile, conn, (fileSize - receivedBytes))
// 			// Flush the network connection buffer
// 			// (receivedBytes + BUFFERSIZE)-fileSize is the total bytes of the file
// 			conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
// 			break
// 		}
// 		io.CopyN(newFile, conn, BUFFERSIZE)
// 		receivedBytes += BUFFERSIZE
// 	}
// 	fmt.Println("Done receiving file...")

// 	// Parameterize Schedule with zip filename
// 	pwd, err := os.Getwd()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	zipPath := filepath.Join(pwd, fileName)
// 	output := Schedule(runtime, imageNum, zipPath, app, version, all)

// 	writer.Write(output)
// 	writer.Flush()
// 	fmt.Println("Sent output to server...")
// }
