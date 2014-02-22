package main

import "github.com/nsf/termbox-go"
import "time"
import "os"

func main() {
	// Initialize a navigator in the current directory.
	nav := directory.NewNavigator(os.Getwd())

	// main application loop
	for {
		// Read a character from STDIN.
		character := input.Read(os.Stdin)

		// Invoke the correspoding navigator action
		input.Map(character, nav)
	}
}
