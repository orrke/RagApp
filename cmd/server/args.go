package main

import "os"

// Args is a struct that holds the parameters for the server address and port
type Args struct {
	address string
	port    string
}

// GetArgs fetches the arguments for the Args struct so the server can be put at the right address and port
func GetArgs() Args {
	args := os.Args[1:]

	address := "localhost"
	port := "5051"

	for i, arg := range args {
		if arg == "--address" {
			address = args[i+1]
		} else if arg == "--port" {
			port = args[i+1]
		}
	}

	return Args{address, port}
}
