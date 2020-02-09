package cnc

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
