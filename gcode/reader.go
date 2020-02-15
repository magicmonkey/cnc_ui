package gcode

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/magicmonkey/cnc/gamepad/display"
)

type StatusResponse struct {
	Status   string    `json:"status"`
	Position []float32 `json:"pos"`
	Babystep float32   `json:"babystep"`
	Homed    []int     `json:"homed"`
}

func serial_read_loop() {
	buf := make([]byte, 32)
	for {
		n, err := SerialPort.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {
			buf = buf[:n]
			fmt.Println("Rx: ", hex.EncodeToString(buf))
		}
	}
}

func read() {
	// Setup reader in a loop, mostly for M408 responses
	b := []byte(`{"status":"I","heaters":[0.0],"active":[0.0],"standby":[0.0],"hstat":[0],"pos":[646.000,510.000,85.500],"machine":[646.000,510.000,85.500],"sfactor":100.00,"efactor":[],"babystep":0.000,"tool":-1,"probe":"0","fanPercent":[0.0,0.0,100.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0],"fanRPM":0,"homed":[1,1,1],"msgBox.mode":-1}`)
	var s StatusResponse
	err := json.Unmarshal(b, &s)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", s)
	return

	// Check for status message
	// status: I=idle, P=printing from SD card, S=stopped (i.e. needs a reset), C=running config file (i.e starting up), A=paused, D=pausing, R=resuming from a pause, B=busy (e.g. running a macro), F=performing firmware update
	switch s.Status {
	case "I":
		display.ShowBackground("")
	case "B":
		display.ShowBackground("Busy")
	case "P":
		display.ShowBackground("Exec")
	case "S":
		display.ShowBackground("Stop")
	case "C":
		display.ShowBackground("Strt")
	case "A":
		display.ShowBackground("Paus")
	case "D":
		display.ShowBackground("Paus")
	case "R":
		display.ShowBackground("Resm")
	case "F":
		display.ShowBackground("FwUp")
	}
}
