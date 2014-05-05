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

	// Create a command channel that we'll use to
	// communicate commands to the navigator.
	commands := make(chan string)

	// Create a buffer channel that the view will
	// use to push updates to the view after state changes.
	buffers := make(chan *view.Buffer)

	// Start the view in a goroutine.
	go view.New(buffers)

	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	go directory.NewNavigator(currentPath, commands, buffers)

	// main application loop
	for {
		// Read a character from STDIN.
		character := input.Read(os.Stdin)

		// Map the character to its corresponding command.
		command := input.Map[character]

		if command == "Quit" {
			break
		}

		// Send the command along to the navigator.
		commands <- command
	}
}
