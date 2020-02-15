package display

import (
	"fmt"

	flp "github.com/adrianh-za/go-fourletterphat-rpi"
	i2c "github.com/d2r2/go-i2c"
)

var i2c_port *i2c.I2C
var foregroundText string
var backgroundText string

func Close() {
	fmt.Println("Closing display...")
	flp.ClearChars(i2c_port)
	i2c_port.Close()
	fmt.Println("Closed display")
}

func ShowForeground(t string) {
	foregroundText = t
	updateDisplay()
}

func ShowBackground(t string) {
	backgroundText = t
	updateDisplay()
}

func updateDisplay() {
	if foregroundText != "" {
		fmt.Println("***", foregroundText)
		flp.WriteCharacters(i2c_port, foregroundText)
	} else {
		fmt.Println("***", backgroundText)
		flp.WriteCharacters(i2c_port, backgroundText)
	}
}

func Initialise() {
	fmt.Println("Opening display...")
	var err error
	i2c_port, err = i2c.NewI2C(flp.AddressDefault, 1)
	if err != nil {
		panic(err)
	}

	// Initialize the LED display
	err = flp.Initialize(i2c_port) // Will set brightness to 15, will switch of blink, clears display
	if err != nil {
		panic(err)
	}
	fmt.Println("Display open")
}
