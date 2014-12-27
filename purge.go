package main

import (
	"os"
	"runtime"
	"fmt"

	"github.com/jmacdonald/purge/filesystem/directory/navigator"
	"github.com/jmacdonald/purge/input"
	"github.com/jmacdonald/purge/view"
)

func main() {
	// Use all available "logical CPUs", as reported by the machine.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Determine in which directory to start,
	// validating the path if passed by the user.
	var startingPath string
	if len(os.Args) > 1 {
		startingPath = os.Args[1]

		// Check that the specified directory exists.
		path, error := os.Stat(startingPath)
		if !(error == nil && path.IsDir()) {
			fmt.Println("Can't open the specified directory.")
			return
		}
	} else {
		// Fall back to the current directory.
		var err error
		startingPath, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	// Initialize (and schedule cleanup for) the view.
	view.Initialize()
	defer view.Close()

	// Create a command channel that we'll use
	// to communicate with the navigator.
	nav := make(chan string)

	// Create a buffer channel that the navigator will
	// use to push updates to the view after state changes.
	buffers := make(chan *view.Buffer)

	// Start the view in a goroutine.
	go view.New(buffers)

	// Start the navigator in the starting directory.
	go navigator.NewNavigator(startingPath, nav, buffers)

	// Listen for user input, relaying the
	// appropriate commands to the navigator.
	for {
		// Read a character from STDIN.
		character := input.Read(os.Stdin)

		// Map the character to its corresponding command.
		command := input.Map[character]

		// Don't pass the quit command along, just exit the application loop.
		if command == "Quit" {
			break
		}

		// Send the command along to the navigator.
		nav <- command
	}
}
