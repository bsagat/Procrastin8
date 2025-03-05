package models

import "flag"

var (
	HelpFlag        = flag.Bool("help", false, "Prints information about program")
	Port            = flag.String("port", "8080", "Deafult port number")
	HelpInformation = `Here is Help information
`
)
