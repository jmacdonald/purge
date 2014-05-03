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

	// Initialize a navigator in the current directory.
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	nav := directory.NewNavigator(currentPath)

	// Create a buffer channel with which we'll
	// communicate with the view goroutine.
	viewBuffers := make(chan *view.Buffer)

	// Create an exit channel with which we'll
	// tell the view to clean up prior to an exit.
	exit := make(chan bool)

	// Create a complete channel with which the view will use
	// to tell us an operation we're blocking on is complete.
	complete := make(chan bool)

	// Start the view in a goroutine.
	go view.New(viewBuffers, exit, complete)

	// Wait for the view to initialize, and then do an initial render.
	<-complete
	viewBuffers <- nav.View(view.Height())

	// main application loop
	for {
		// Read a character from STDIN.
		character := input.Read(os.Stdin)

		// Invoke the correspoding navigator action,
		// and exit the main loop if it returns true (exit request).
		if input.Map(character, nav) {
			// Signal the view to clean up.
			exit <- true

			// Wait until the view signals that it's complete.
			<-complete

			// Break out of the main application loop.
			break
		}

		// Render the updated state.
		viewBuffers <- nav.View(view.Height())
	}
}
