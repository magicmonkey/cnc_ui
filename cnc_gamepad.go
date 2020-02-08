package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/jacobsa/go-serial/serial"
	"github.com/jinzhu/copier"
	"github.com/karalabe/hid"

	flp "github.com/adrianh-za/go-fourletterphat-rpi"
	i2c "github.com/d2r2/go-i2c"
)

type Stick_s struct {
	LeftRight int
	UpDown    int
	Press     bool
	Left      bool
	Right     bool
	Up        bool
	Down      bool
}

type ButtonState struct {
	Dpad struct {
		Left  bool
		Right bool
		Up    bool
		Down  bool
	}
	Buttons struct {
		X      bool
		Y      bool
		A      bool
		B      bool
		Select bool
		Start  bool
		Star   bool
		Home   bool
	}
	Shoulder struct {
		L1 bool
		L2 bool
		R1 bool
		R2 bool
	}
	Sticks struct {
		Left  Stick_s
		Right Stick_s
	}
}

func (b *ButtonState) InitFromRaw(raw []byte) {
	// D-pad
	switch raw[2] {
	case 0:
		b.Dpad.Up = true
		b.Dpad.Right = false
		b.Dpad.Down = false
		b.Dpad.Left = false
		break
	case 1:
		b.Dpad.Up = true
		b.Dpad.Right = true
		b.Dpad.Down = false
		b.Dpad.Left = false
		break
	case 2:
		b.Dpad.Up = false
		b.Dpad.Right = true
		b.Dpad.Down = false
		b.Dpad.Left = false
		break
	case 3:
		b.Dpad.Up = false
		b.Dpad.Right = true
		b.Dpad.Down = true
		b.Dpad.Left = false
		break
	case 4:
		b.Dpad.Up = false
		b.Dpad.Right = false
		b.Dpad.Down = true
		b.Dpad.Left = false
		break
	case 5:
		b.Dpad.Up = false
		b.Dpad.Right = false
		b.Dpad.Down = true
		b.Dpad.Left = true
		break
	case 6:
		b.Dpad.Up = false
		b.Dpad.Right = false
		b.Dpad.Down = false
		b.Dpad.Left = true
		break
	case 7:
		b.Dpad.Up = true
		b.Dpad.Right = false
		b.Dpad.Down = false
		b.Dpad.Left = true
		break
	default:
		b.Dpad.Left = false
		b.Dpad.Right = false
		b.Dpad.Up = false
		b.Dpad.Down = false
		break
	}
	if (raw[0] & 0b00000001) > 0 {
		b.Buttons.A = true
	} else {
		b.Buttons.A = false
	}
	if (raw[0] & 0b00000010) > 0 {
		b.Buttons.B = true
	} else {
		b.Buttons.B = false
	}
	if (raw[0] & 0b00001000) > 0 {
		b.Buttons.X = true
	} else {
		b.Buttons.X = false
	}
	if (raw[0] & 0b00010000) > 0 {
		b.Buttons.Y = true
	} else {
		b.Buttons.Y = false
	}
	if (raw[1] & 0b00000100) > 0 {
		b.Buttons.Select = true
	} else {
		b.Buttons.Select = false
	}
	if (raw[1] & 0b00001000) > 0 {
		b.Buttons.Start = true
	} else {
		b.Buttons.Start = false
	}
	if (raw[0] & 0b00000100) > 0 {
		b.Buttons.Home = true
	} else {
		b.Buttons.Home = false
	}
	if (raw[0] & 0b01000000) > 0 {
		b.Shoulder.L1 = true
	} else {
		b.Shoulder.L1 = false
	}
	if (raw[1] & 0b00000001) > 0 {
		b.Shoulder.L2 = true
	} else {
		b.Shoulder.L2 = false
	}
	if (raw[0] & 0b10000000) > 0 {
		b.Shoulder.R1 = true
	} else {
		b.Shoulder.R1 = false
	}
	if (raw[1] & 0b00000010) > 0 {
		b.Shoulder.R2 = true
	} else {
		b.Shoulder.R2 = false
	}
	if (raw[1] & 0b00100000) > 0 {
		b.Sticks.Left.Press = true
	} else {
		b.Sticks.Left.Press = false
	}
	if (raw[1] & 0b01000000) > 0 {
		b.Sticks.Right.Press = true
	} else {
		b.Sticks.Right.Press = false
	}
	b.Sticks.Left.LeftRight = int(raw[3]) - 128
	if b.Sticks.Left.LeftRight < -100 {
		b.Sticks.Left.Left = true
		b.Sticks.Left.Right = false
	} else if b.Sticks.Left.LeftRight > 100 {
		b.Sticks.Left.Left = false
		b.Sticks.Left.Right = true
	} else {
		b.Sticks.Left.Left = false
		b.Sticks.Left.Right = false
	}
	b.Sticks.Left.UpDown = 128 - int(raw[4])
	if b.Sticks.Left.UpDown < -100 {
		b.Sticks.Left.Down = true
		b.Sticks.Left.Up = false
	} else if b.Sticks.Left.UpDown > 100 {
		b.Sticks.Left.Down = false
		b.Sticks.Left.Up = true
	} else {
		b.Sticks.Left.Down = false
		b.Sticks.Left.Up = false
	}
	b.Sticks.Right.LeftRight = int(raw[5]) - 128
	if b.Sticks.Right.LeftRight < -100 {
		b.Sticks.Right.Left = true
		b.Sticks.Right.Right = false
	} else if b.Sticks.Right.LeftRight > 100 {
		b.Sticks.Right.Left = false
		b.Sticks.Right.Right = true
	} else {
		b.Sticks.Right.Left = false
		b.Sticks.Right.Right = false
	}
	b.Sticks.Right.UpDown = 128 - int(raw[6])
	if b.Sticks.Right.UpDown < -100 {
		b.Sticks.Right.Down = true
		b.Sticks.Right.Up = false
	} else if b.Sticks.Right.UpDown > 100 {
		b.Sticks.Right.Down = false
		b.Sticks.Right.Up = true
	} else {
		b.Sticks.Right.Down = false
		b.Sticks.Right.Up = false
	}
}

var SerialPort io.ReadWriteCloser
var i2c_port *i2c.I2C

func main() {

	fmt.Println("Starting...")

	fmt.Println("Opening serial port...")
	// Serial port stuff
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

	fmt.Println("Opening display...")
	i2c_port, err = i2c.NewI2C(flp.AddressDefault, 1)
	if (err != nil) {
		panic(err)
	}

	// Initialize the LED display
	flp.Initialize(i2c_port) // Will set brightness to 15, will switch of blink, clears display
	fourletter("CNC")
	fmt.Println("Display open")

	for {
		fmt.Println("Opening USB port...")
		processDevice()
		time.Sleep(1 * time.Second)
	}
}

func processDevice() {
	//devs := hid.Enumerate(0x057e, 0x2009)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()



	devs := hid.Enumerate(0x2dc8, 0x6000)
	if len(devs) == 0 {
		fmt.Println("No USB devices")
		return
	}

	dev, err := devs[0].Open()
	if err != nil {
		fmt.Println("Error opening USB device - permissions maybe?")
		panic(err)
		return
	}

	/*
		_, err = dev.Write([]byte{0})
		if err != nil {
			panic(err)
		}
	*/

	buffer := make([]byte, 7)

	buttons := new(ButtonState)
	prev_buttons := new(ButtonState)

	fmt.Println("Reading from USB...")

	for {
		_, err = dev.Read(buffer)
		if err != nil {
			panic(err)
		}
		buttons.InitFromRaw(buffer)
		processButtonPress(buttons, prev_buttons)
		copier.Copy(prev_buttons, buttons)
	}

}

func processButtonPress(curr_buttons *ButtonState, prev_buttons *ButtonState) {
	// Red button (A) is emergency stop
	if curr_buttons.Buttons.A && curr_buttons.Buttons.A != prev_buttons.Buttons.A {
		send_gcode("M112")
		send_gcode("M999")
	}

	// - Home (with L2)
	// - Set zero (with L1)
	// - Return to zero (with R1)
	// - Probe (with R2)

	if curr_buttons.Shoulder.L2 != prev_buttons.Shoulder.L2 {
		if curr_buttons.Shoulder.L2 {
			fourletter("Home")
		} else {
			fourletter("")
		}
	}

	if curr_buttons.Shoulder.L1 != prev_buttons.Shoulder.L1 {
		if curr_buttons.Shoulder.L1 {
			fourletter("SetZ")
		} else {
			fourletter("")
		}
	}

	if curr_buttons.Shoulder.R1 != prev_buttons.Shoulder.R1 {
		if curr_buttons.Shoulder.R1 {
			fourletter("RetZ")
		} else {
			fourletter("")
		}
	}

	if curr_buttons.Shoulder.R2 != prev_buttons.Shoulder.R2 {
		if curr_buttons.Shoulder.R2 {
			fourletter("Prob")
		} else {
			fourletter("")
		}
	}

	// Yellow button (B) is z-axis operations:
	if curr_buttons.Buttons.B && curr_buttons.Buttons.B != prev_buttons.Buttons.B {
		if curr_buttons.Shoulder.L2 {
			send_gcode("G28 Z")
		}
		if curr_buttons.Shoulder.L1 {
			send_gcode("G10 P2 L20 Z0")
			send_gcode("M500")
		}
		if curr_buttons.Shoulder.R1 {
			send_gcode("G0 Z0")
		}
		if curr_buttons.Shoulder.R2 {
			// TODO
		}
	}

	// Green button (Y) is x-axis operations:
	if curr_buttons.Buttons.Y && curr_buttons.Buttons.Y != prev_buttons.Buttons.Y {
		if curr_buttons.Shoulder.L2 {
			send_gcode("G28 X")
		}
		if curr_buttons.Shoulder.L1 {
			send_gcode("G10 P2 L20 X0")
			send_gcode("M500")
		}
		if curr_buttons.Shoulder.R1 {
			send_gcode("G0 X0")
		}
		if curr_buttons.Shoulder.R2 {
			// TODO
		}
	}

	// Blue button (X) is y-axis operations:
	if curr_buttons.Buttons.X && curr_buttons.Buttons.X != prev_buttons.Buttons.X {
		if curr_buttons.Shoulder.L2 {
			send_gcode("G28 Y")
		}
		if curr_buttons.Shoulder.L1 {
			send_gcode("G10 P2 L20 Y0")
			send_gcode("M500")
		}
		if curr_buttons.Shoulder.R1 {
			send_gcode("G0 Y0")
		}
		if curr_buttons.Shoulder.R2 {
			// TODO
		}
	}

	// Home button (below yellow) is "home x and y"
	if curr_buttons.Buttons.Home && curr_buttons.Buttons.Home != prev_buttons.Buttons.Home {
		send_gcode("G0 X0 Y0")
	}

	// D-pad up is plus Y
	if curr_buttons.Dpad.Up && curr_buttons.Dpad.Up != prev_buttons.Dpad.Up {
		distance := 0.0
		if curr_buttons.Shoulder.L2 {
			distance = 0.5
		}
		if curr_buttons.Shoulder.L1 {
			distance = 1
		}
		if curr_buttons.Shoulder.R1 {
			distance = 10
		}
		if curr_buttons.Shoulder.R2 {
			distance = 50
		}
		if distance == 0.0 {
			return
		}
		send_gcode("M120")
		send_gcode("G91")
		cmd := fmt.Sprintf("G1 Y%.1f F2000", distance)
		send_gcode(cmd)
		send_gcode("M121")
	}

	// D-pad down is minus Y
	if curr_buttons.Dpad.Down && curr_buttons.Dpad.Down != prev_buttons.Dpad.Down {
		distance := 0.0
		if curr_buttons.Shoulder.L2 {
			distance = 0.5
		}
		if curr_buttons.Shoulder.L1 {
			distance = 1
		}
		if curr_buttons.Shoulder.R1 {
			distance = 10
		}
		if curr_buttons.Shoulder.R2 {
			distance = 50
		}
		if distance == 0.0 {
			return
		}
		send_gcode("M120")
		send_gcode("G91")
		cmd := fmt.Sprintf("G1 Y-%.1f F2000", distance)
		send_gcode(cmd)
		send_gcode("M121")
	}

	// D-pad right is plus X
	if curr_buttons.Dpad.Right && curr_buttons.Dpad.Right != prev_buttons.Dpad.Right {
		distance := 0.0
		if curr_buttons.Shoulder.L2 {
			distance = 0.5
		}
		if curr_buttons.Shoulder.L1 {
			distance = 1
		}
		if curr_buttons.Shoulder.R1 {
			distance = 10
		}
		if curr_buttons.Shoulder.R2 {
			distance = 50
		}
		if distance == 0.0 {
			return
		}
		send_gcode("M120")
		send_gcode("G91")
		cmd := fmt.Sprintf("G1 X%.1f F2000", distance)
		send_gcode(cmd)
		send_gcode("M121")
	}

	// D-pad left is minus X
	if curr_buttons.Dpad.Left && curr_buttons.Dpad.Left != prev_buttons.Dpad.Left {
		distance := 0.0
		if curr_buttons.Shoulder.L2 {
			distance = 0.5
		}
		if curr_buttons.Shoulder.L1 {
			distance = 1
		}
		if curr_buttons.Shoulder.R1 {
			distance = 10
		}
		if curr_buttons.Shoulder.R2 {
			distance = 50
		}
		if distance == 0.0 {
			return
		}
		send_gcode("M120")
		send_gcode("G91")
		cmd := fmt.Sprintf("G1 X-%.1f F2000", distance)
		send_gcode(cmd)
		send_gcode("M121")
	}

	// Joystick Left pushed up is plus Z
	if curr_buttons.Sticks.Left.Up && curr_buttons.Sticks.Left.Up != prev_buttons.Sticks.Left.Up {
		distance := 0.0
		if curr_buttons.Shoulder.L2 {
			distance = 0.1
		}
		if curr_buttons.Shoulder.L1 {
			distance = 0.5
		}
		if curr_buttons.Shoulder.R1 {
			distance = 1
		}
		if curr_buttons.Shoulder.R2 {
			distance = 10
		}
		if distance == 0.0 {
			return
		}
		send_gcode("M120")
		send_gcode("G91")
		cmd := fmt.Sprintf("G1 Z%.1f F2000", distance)
		send_gcode(cmd)
		send_gcode("M121")
	}

	// Joystick Left pushed down is minus Z
	if curr_buttons.Sticks.Left.Down && curr_buttons.Sticks.Left.Down != prev_buttons.Sticks.Left.Down {
		distance := 0.0
		if curr_buttons.Shoulder.L2 {
			distance = 0.1
		}
		if curr_buttons.Shoulder.L1 {
			distance = 0.5
		}
		if curr_buttons.Shoulder.R1 {
			distance = 1
		}
		if curr_buttons.Shoulder.R2 {
			distance = 10
		}
		if distance == 0.0 {
			return
		}
		send_gcode("M120")
		send_gcode("G91")
		cmd := fmt.Sprintf("G1 Z-%.1f F2000", distance)
		send_gcode(cmd)
		send_gcode("M121")
	}

}

func send_gcode(gcode string) {
	fmt.Println(gcode)
	_, err := SerialPort.Write([]byte(gcode + "\n"))
	if err != nil {
		panic(err)
	}
}

func cleanup() {
	fmt.Println("Closing display...")
	SerialPort.Close()
	flp.ClearChars(i2c_port)
	i2c_port.Close()
	fmt.Println("Closed display")
}

func fourletter(t string) {
	flp.WriteCharacters(i2c_port, t)
}
