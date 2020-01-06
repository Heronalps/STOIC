package main

import (
	"github.com/heronalps/STOIC/client"
	"github.com/heronalps/STOIC/cmd"
)

func main() {
	cmd.Execute()
}

func init() {
	//Set up a HashSet for all Nautilus Runtimes
	// The benefit of using struct{} is saving space in the value field
	client.NautilusRuntimes = make(map[string]struct{})
	client.NautilusRuntimes["cpu"] = struct{}{}
	client.NautilusRuntimes["gpu1"] = struct{}{}
	client.NautilusRuntimes["gpu2"] = struct{}{}
}
