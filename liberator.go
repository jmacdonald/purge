package main

import "os"
import "fmt"
import "github.com/jmacdonald/liberator/filesystem/directory"
import "github.com/jmacdonald/liberator/input"

func main() {
	currentPath, _ := os.Getwd()
	navigator := directory.NewNavigator(currentPath)

	for {
		input.Read(os.Stdin, navigator)
	}
}
