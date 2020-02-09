package main

import (
	"fmt"
	"github.com/magicmonkey/cnc/gamepad/cnc"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("Starting...")
	cnc.Initialise()

	// Look for ctrl-C
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Caught ctrl-C, exiting...")
		cleanup()
		os.Exit(0)
	}()

	fmt.Println("Running...")
	cnc.Run()
}

func cleanup() {
	cnc.Close()
}
