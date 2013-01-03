Pertelian Driver for Go.
------------------------

This is a simple driver to drive the Pertelian X2040 USB LCD display screen.


Usage:
------

	package main

	import (
		"time"

		"github.com/AmandaCameron/go.pert"
	)


	func main() {
		screen, err := pert.NewScreen("/dev/ttyUSB0") // Your device here.
		if err != nil {
			panic(err)
		}

		defer screen.Close()

		screen.Put(pert.LINE_ONE, "Hello World!")
		screen.Put(pert.LINE_TWO, "How are you?")
		screen.Put(pert.LINE_THREE, "I'm doing fine, myself.")
		screen.Put(pert.LINE_FOUR, "Hey, do you want pie?")

		screen.SetBacklight(false)
		time.Sleep(10 * time.Second)
		screen.SetBacklight(true)
	}



TODO:
-----
  * Custom Glyph Drawing
  * Further testing to make sure it doesn't break the display.