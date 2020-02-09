package cnc

import (
	"fmt"
	"io"
	"time"

	"github.com/jinzhu/copier"
	"github.com/karalabe/hid"
	"github.com/magicmonkey/cnc/gamepad/display"
	"github.com/magicmonkey/cnc/gamepad/gcode"
)

var SerialPort io.ReadWriteCloser

func processDevice() {
	//devs := hid.Enumerate(0x057e, 0x2009)

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
		gcode.SendGcode("M112")
		gcode.SendGcode("M999")
		return
	}

	// - Home (with L2)
	// - Set zero (with L1)
	// - Return to zero (with R1)
	// - Probe (with R2)

	if curr_buttons.Shoulder.L2 != prev_buttons.Shoulder.L2 {
		if curr_buttons.Shoulder.L2 {
			display.Show("Home")
		} else {
			display.Show("")
		}
	}

	if curr_buttons.Shoulder.L1 != prev_buttons.Shoulder.L1 {
		if curr_buttons.Shoulder.L1 {
			display.Show("SetZ")
		} else {
			display.Show("")
		}
	}

	if curr_buttons.Shoulder.R1 != prev_buttons.Shoulder.R1 {
		if curr_buttons.Shoulder.R1 {
			display.Show("RetZ")
		} else {
			display.Show("")
		}
	}

	if curr_buttons.Shoulder.R2 != prev_buttons.Shoulder.R2 {
		if curr_buttons.Shoulder.R2 {
			display.Show("Prob")
		} else {
			display.Show("")
		}
	}

	// Yellow button (B) is z-axis operations:
	if curr_buttons.Buttons.B && curr_buttons.Buttons.B != prev_buttons.Buttons.B {
		if curr_buttons.Shoulder.L2 {
			gcode.SendGcode("G28 Z")
		}
		if curr_buttons.Shoulder.L1 {
			gcode.SendGcode("G55")
			gcode.SendGcode("G10 P2 L20 Z0")
			gcode.SendGcode("M500")
		}
		if curr_buttons.Shoulder.R1 {
			gcode.SendGcode("G0 Z0")
		}
		if curr_buttons.Shoulder.R2 {
			// TODO
		}
	}

	// Green button (Y) is x-axis operations:
	if curr_buttons.Buttons.Y && curr_buttons.Buttons.Y != prev_buttons.Buttons.Y {
		if curr_buttons.Shoulder.L2 {
			gcode.SendGcode("G28 X")
		}
		if curr_buttons.Shoulder.L1 {
			gcode.SendGcode("G55")
			gcode.SendGcode("G10 P2 L20 X0")
			gcode.SendGcode("M500")
		}
		if curr_buttons.Shoulder.R1 {
			gcode.SendGcode("G0 X0")
		}
		if curr_buttons.Shoulder.R2 {
			// TODO
		}
	}

	// Blue button (X) is y-axis operations:
	if curr_buttons.Buttons.X && curr_buttons.Buttons.X != prev_buttons.Buttons.X {
		if curr_buttons.Shoulder.L2 {
			gcode.SendGcode("G28 Y")
		}
		if curr_buttons.Shoulder.L1 {
			gcode.SendGcode("G55")
			gcode.SendGcode("G10 P2 L20 Y0")
			gcode.SendGcode("M500")
		}
		if curr_buttons.Shoulder.R1 {
			gcode.SendGcode("G0 Y0")
		}
		if curr_buttons.Shoulder.R2 {
			// TODO
		}
	}

	// Home button (below yellow) is "home x and y"
	if curr_buttons.Buttons.Home && curr_buttons.Buttons.Home != prev_buttons.Buttons.Home {
		if curr_buttons.Shoulder.L2 {
			// Home
			gcode.SendGcode("G28 X Y Z")
		}
		if curr_buttons.Shoulder.L1 {
			// Set zero
			gcode.SendGcode("G55")
			gcode.SendGcode("G10 P2 L20 X0 Y0")
			gcode.SendGcode("M500")
		}
		if curr_buttons.Shoulder.R1 {
			// Return to zero
			gcode.SendGcode("G0 X0 Y0")
		}
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
		gcode.SendGcode("M120")
		gcode.SendGcode("G91")
		cmd := fmt.Sprintf("G1 Y%.1f F2000", distance)
		gcode.SendGcode(cmd)
		gcode.SendGcode("M121")
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
		gcode.SendGcode("M120")
		gcode.SendGcode("G91")
		cmd := fmt.Sprintf("G1 Y-%.1f F2000", distance)
		gcode.SendGcode(cmd)
		gcode.SendGcode("M121")
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
		gcode.SendGcode("M120")
		gcode.SendGcode("G91")
		cmd := fmt.Sprintf("G1 X%.1f F2000", distance)
		gcode.SendGcode(cmd)
		gcode.SendGcode("M121")
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
		gcode.SendGcode("M120")
		gcode.SendGcode("G91")
		cmd := fmt.Sprintf("G1 X-%.1f F2000", distance)
		gcode.SendGcode(cmd)
		gcode.SendGcode("M121")
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
		gcode.SendGcode("M120")
		gcode.SendGcode("G91")
		cmd := fmt.Sprintf("G1 Z%.1f F2000", distance)
		gcode.SendGcode(cmd)
		gcode.SendGcode("M121")
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
		gcode.SendGcode("M120")
		gcode.SendGcode("G91")
		cmd := fmt.Sprintf("G1 Z-%.1f F2000", distance)
		gcode.SendGcode(cmd)
		gcode.SendGcode("M121")
	}

	if curr_buttons.Buttons.Start && curr_buttons.Buttons.Start != prev_buttons.Buttons.Start {
		gcode.SendGcode("M292 P0")
	}
}

func Close() {
	gcode.Close()
	display.Close()
}

func Initialise() {
	gcode.Initialise()
	display.Initialise()

	display.Show("CNC")
}

func Run() {
	for {
		fmt.Println("Opening USB port...")
		processDevice()
		time.Sleep(1 * time.Second)
	}
}
