package utils

import (
	"fmt"
	"os"
)

func PrintHelp() {
	helpInformation := `Here is Help information`
	fmt.Println(helpInformation)
	os.Exit(0)
}
