package gcode

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/magicmonkey/cnc/gamepad/display"
)

type StatusResponse struct {

	// Response to M408 (status)
	Status   string    `json:"status"`
	Position []float32 `json:"pos"`
	Babystep float32   `json:"babystep"`
	Homed    []int     `json:"homed"`

	// Response to M20 (file browser)
	Dir   string   `json:"dir"`
	First int      `json:"first"`
	Files []string `json:"files"`
	Next  int      `json:"next"`
	Err   int      `json:"err"`
}

// {"dir":"0:/gcodes/","first":0,"files":["*Phone stand","railway straight 12mm 160mm.nc","railway straight 12mm 160mm v2.nc","railway curve right 12mm.nc","spoilboard slot 1.nc","Surfacing pass 130x140.nc","spoilboard slot 2.nc","*xyz-probe","8mm dowel hole with eighth inch bit.nc","*House number","*celtic coaster","*ui","*emmie heart sign","*surfacing","*Thread storage","*Airing cupboard shelves"],"next":0,"err":0}

func serial_read_loop() {
	for {
		dec := json.NewDecoder(SerialPort)
		for {
			var m StatusResponse
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			fmt.Printf("%#v\n", m)
			interpret(m)
		}
	}
}

func interpret(s StatusResponse) {
	// Check for status message

	if s.Status != "" {
		// It's an M408 response
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
	} else if s.Dir != "" {
		// It's an M20 response
		FilesList[s.Dir] = s.Files
	}
}
