package utils

import (
	"fmt"
	"os"
)

func PrintHelp() {
	helpInformation := `./TodoApp --help
Usage:
  	TodoApp [--port <N>]
  	TodoApp --help
Options:
  	--port N     Port number
	--help 		 Prints help information`
	fmt.Println(helpInformation)
	os.Exit(0)
}
