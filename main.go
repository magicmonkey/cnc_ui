package main

import (
	"fmt"
	"github.com/magicmonkey/cnc/gamepad/cnc"
)

func main() {
	fmt.Println("Starting...")
	cnc.Initialise()

	fmt.Println("Running...")
	cnc.Run()
}
