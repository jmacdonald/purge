package main

import (
	"github.com/jmacdonald/liberator/filesystem/directory"
	"github.com/jmacdonald/liberator/input"
	"github.com/jmacdonald/liberator/view"
	"github.com/nsf/termbox-go"
	"os"
	"runtime"
)

func main() {
	// Use all available "logical CPUs", as reported by the machine.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize a navigator in the current directory.
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	nav := directory.NewNavigator(currentPath)

	// Set up the terminal screen.
	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Do an initial render.
	view.Render(nav)

	// main application loop
	for {
		// Read a character from STDIN.
		character := input.Read(os.Stdin)

		// Invoke the correspoding navigator action,
		// and exit the main loop if it returns true (exit request).
		if input.Map(character, nav) {
			break
		}

		// Render the updated state.
		view.Render(nav)
	}
}
