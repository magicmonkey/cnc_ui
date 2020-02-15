package gcode

import (
	"fmt"
	"io"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

var SerialPort io.ReadWriteCloser

var FilesList map[string][]string

func Initialise() {
	fmt.Println("Opening serial port...")

	options := serial.OpenOptions{
		PortName:              "/dev/serial0",
		BaudRate:              57600,
		DataBits:              8,
		StopBits:              1,
		InterCharacterTimeout: 1000,
	}
	var err error
	SerialPort, err = serial.Open(options)
	if err != nil {
		panic(err)
	}

	fmt.Println("Serial port open")

	FilesList = make(map[string][]string)

	// Start the serial port read loop
	go serial_read_loop()

	// Start the files request loop
	go gcode_send_loop(1*time.Second, "M408 S0")

	// Start the status request loop
	go gcode_send_loop(3*time.Second, "M20 S2")
}

func gcode_send_loop(interval time.Duration, gcode string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			SendGcode(gcode)
		}
	}
}

func Close() {
	fmt.Println("Closing serial port...")
	SerialPort.Close()
	fmt.Println("Closed serial port")
}
