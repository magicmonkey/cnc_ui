package gcode

import (
	"fmt"
	"io"
	"time"

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

	fmt.Println("Serial port open")

	// Start the serial port read loop
	go serial_read_loop()

	// Start the status request loop
	go status_request_loop()
}

func status_request_loop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			SendGcode("M408 S0")
		}
	}
}

func Close() {
	fmt.Println("Closing serial port...")
	SerialPort.Close()
	fmt.Println("Closed serial port")
}
