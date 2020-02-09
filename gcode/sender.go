package gcode

import (
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
)

var SerialPort io.ReadWriteCloser

func SendGcode(gcode string) {
	fmt.Println(gcode)
	_, err := SerialPort.Write([]byte(gcode + "\n"))
	if err != nil {
		panic(err)
	}
}

func Initialise() {
	fmt.Println("Opening serial port...")
	read()

	options := serial.OpenOptions{
		PortName:        "/dev/serial0",
		BaudRate:        57600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}
	var err error
	SerialPort, err = serial.Open(options)
	if err != nil {
		panic(err)
	}
	defer SerialPort.Close()

	fmt.Println("Serial port open")
}

func Close() {
	fmt.Println("Closing serial port...")
	SerialPort.Close()
	fmt.Println("Closed serial port")
}
