package main

import (
	"fmt"

	"github.com/heronalps/STOIC/server"
)

func main() {
	// fmt.Println("This is main function!")
	// fmt.Println(server.GetBandWidth())
	// fmt.Println(server.Extrapolate("cpu", 1.234))
	for i := 0; i < 10; i++ {
		fmt.Println(server.ImageCache())
	}

}
