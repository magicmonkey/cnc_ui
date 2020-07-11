# Control your Ooznest Workbee CNC with a Nintendo-style gamepad controller

This project was created because I wanted a more tangible UI for my CNC than the webpage interface.  After hitting the wrong movement button on my touchscreen phone web browser a few times, I decided that real buttons would be best.  Initially, I made an arcade-machine style UI, but then refined it to use a gamepad instead.

The gamepad used is the 8BitDo SF30 Pro ([https://www.amazon.co.uk/8Bitdo-SF30-PRO-Bluetooth-Gamepad/dp/B0748S3GXG/]) in a particular mode (_// TODO add detalis of the mode here //_)

The project uses a Raspberry Pi Zero with the gamepad plugged in over USB, and a serial cable joining it to the socket on the Duet which is designed for the "PanelDue" interface.  This port accepts GCode over serial, so can fully control the CNC.

Alternatives which I discounted:
* send the GCode over the web interface, however I would be concerned about the reliability of the "emergency stop" button with this
* attach the gamepad to the Raspberry Pi over bluetooth instead of by USB, however the 8BitDo gamepad which I'm using goes into a "sleep" mode if it isn't used for 15 minutes and then has to reconnect, which is a tedious delay


// TODO add photos and more write-up here

