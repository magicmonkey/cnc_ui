package gcode

import (
	"encoding/json"
	"fmt"
)

type StatusResponse struct {
	Status   string    `json:"status"`
	Position []float32 `json:"pos"`
	Babystep float32   `json:"babystep"`
	Homed    []int     `json:"homed"`
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
}
