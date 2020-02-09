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

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	fmt.Println("Running...")
	cnc.Run()
}

func cleanup() {
	cnc.Close()
}
