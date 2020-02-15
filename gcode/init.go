package gcode

import (
	"fmt"
	"io"

	"github.com/jacobsa/go-serial/serial"
)

var SerialPort io.ReadWriteCloser

func Initialise() {
	fmt.Println("Opening serial port...")

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

	// Start the status request loop
	//go status_request_loop()
}

func Close() {
	fmt.Println("Closing serial port...")
	SerialPort.Close()
	fmt.Println("Closed serial port")
}
