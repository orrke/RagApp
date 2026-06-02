package main

import (
	"flag"
)

// Args is a struct that holds the parameters for the server address and port
type Args struct {
	address  string
	port     string
	docsPath *string
	model    *string
}

// GetArgs fetches the arguments for the Args struct so the server can be put at the right address and port, and process the documents at the right
func GetArgs() Args {
	//setup the different flags the program supports
	address := flag.String("address", "0.0.0.0", "Server address")
	port := flag.String("port", "5051", "Server port")
	docsPath := flag.String("docs", "", "Documents path")
	model := flag.String("model", "", "Documents model path")

	//parse the flags given to the program
	flag.Parse()

	//convert the pointer to nil if the path is empty
	if *docsPath == "" {
		docsPath = nil
	}

	if *model == "" {
		model = nil
	}

	return Args{
		address:  *address,
		port:     *port,
		docsPath: docsPath,
		model:    model,
	}
}
