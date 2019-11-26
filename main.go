package main

import "github.com/heronalps/STOIC/server"

func main() {
	// fmt.Println("This is main function!")
	// fmt.Println(server.GetBandWidth())
	// fmt.Println(server.Extrapolate("cpu", 1.234))
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(server.ImageCache())
	// }
	// client.SocketClient(5001)
	// server.SocketServer("127.0.0.1", 5001, "cpu", 10)
	server.Schedule("127.0.0.1", 5001)
	// client.Deploy("racelab", "image-clf-inf", 0)

}
