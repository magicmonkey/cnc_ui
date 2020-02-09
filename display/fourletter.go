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
		flp.WriteCharacters(i2c_port, foregroundText)
	} else {
		flp.WriteCharacters(i2c_port, backgroundText)
	}
}

func Initialise() {
	// Initialize the LED display
	flp.Initialize(i2c_port) // Will set brightness to 15, will switch of blink, clears display
	fmt.Println("Display open")
}
