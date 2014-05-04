package main

import (
	"github.com/jmacdonald/purge/filesystem/directory"
	"github.com/jmacdonald/purge/input"
	"github.com/jmacdonald/purge/view"
	"os"
	"runtime"
)

func main() {
	// Use all available "logical CPUs", as reported by the machine.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize (and schedule cleanup for) the view.
	view.Initialize()
	defer view.Close()

	// Initialize a navigator in the current directory.
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	nav := directory.NewNavigator(currentPath)

	// Create a buffer channel with which we'll
	// communicate with the view goroutine.
	viewBuffers := make(chan *view.Buffer)

	// Start the view in a goroutine.
	go view.New(viewBuffers, exit, complete)

	// main application loop
	for {
		// Read a character from STDIN.
		character := input.Read(os.Stdin)

		// Invoke the correspoding navigator action,
		// and exit the main loop if it returns true (exit request).
		if input.Map(character, nav) {
			// Break out of the main application loop.
			break
		}

		// Render the updated state.
		viewBuffers <- nav.View(view.Height())
	}
}
