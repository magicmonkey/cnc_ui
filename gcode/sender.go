package gcode

import "fmt"

func SendGcode(gcode string) {
	fmt.Println(gcode)
	_, err := SerialPort.Write([]byte(gcode + "\n"))
	if err != nil {
		panic(err)
	}
}
