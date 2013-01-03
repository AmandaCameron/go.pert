package pert

import (
	"os"
	"time"
)

type Screen struct {
	f *os.File
}

// The constants you pass as the line argument for Put & ClearLine
const (
	LINE_ONE   = 0x00
	LINE_TWO   = 0x40
	LINE_THREE = 0x14
	LINE_FOUR  = 0x54
)

func (s *Screen) send(data []byte) error {
	var buff []byte
	buff = make([]byte, 1, 1)
	for _, b := range data {
		buff[0] = b
		_, err := s.f.Write(buff)

		if err != nil {
			return err
		}

		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

func (s *Screen) sendCmd(cmds []byte) error {
	var buff []byte

	buff = make([]byte, len(cmds)*2)

	for i, b := range cmds {
		buff[i*2] = 0xFE
		buff[i*2+1] = b
	}

	return s.send(buff)
}

// Puts a line to the display
// line is one of the LINE_* constants.
// String is assumed truencated to 20 chars it's too long, and padded if it's too short.
func (s *Screen) Put(line int, text string) {
	s.sendCmd([]byte{byte(0x80 | line)})

	if len(text) > 20 {
		text = text[0:20]
	}

	for len(text) < 20 {
		text += " "
	}

	s.send([]byte(text))
}

// Shows a long message on the screen, this clears the screen of all other text.
func (s *Screen) Show(text string) {
	s.Clear()

	if len(text) <= 20 {
		s.Put(LINE_ONE, text)
	} else if len(text) <= 40 {
		s.Put(LINE_ONE, text[:20])
		s.Put(LINE_TWO, text[20:len(text)])
	} else if len(text) <= 60 {
		s.Put(LINE_ONE, text[:20])
		s.Put(LINE_TWO, text[20:40])
		s.Put(LINE_THREE, text[40:len(text)])
	} else if len(text) <= 80 {
		s.Put(LINE_ONE, text[:20])
		s.Put(LINE_TWO, text[20:40])
		s.Put(LINE_THREE, text[40:60])
		s.Put(LINE_FOUR, text[60:len(text)])
	} else {
		s.Put(LINE_ONE, text[:20])
		s.Put(LINE_TWO, text[20:40])
		s.Put(LINE_THREE, text[40:60])
		s.Put(LINE_FOUR, text[60:80])
	}
}

// Clears the specified line.
func (s *Screen) ClearLine(line int) {
	s.Put(line, "")
}

// Clears the whole screen.
func (s *Screen) Clear() {
	s.sendCmd([]byte{0x01})
}

// Turns the backlight on if true, otherwise turns it off.
func (s *Screen) SetBacklight(on bool) {
	if on {
		s.sendCmd([]byte{0x03})
	} else {
		s.sendCmd([]byte{0x02})
	}
}

// Clears the screen, disables the backlight, and closes the tty device.
// Should probably  be called on exit, but not a requirement.
func (s *Screen) Close() error {
	s.Clear()
	s.SetBacklight(false)

	return s.f.Close()
}

// Opens a tty device. It makes no assumptions on the device type,
// just uses os.OpenFile on whatever you give it.
func NewScreen(path string) (*Screen, error) {
	s := new(Screen)

	f, err := os.OpenFile(path, os.O_WRONLY, 0)

	if err != nil {
		return nil, err
	}

	s.f = f

	err = s.sendCmd([]byte{0x38, 0x06, 0x10, 0x0C, 0x01})

	if err != nil {
		return nil, err
	}

	s.SetBacklight(true)

	return s, nil
}
