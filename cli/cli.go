package cli

import (
	"elize/restapi"
	"elize/sample_explorer"
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("Hi, This is ELIZE.\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port=Port_number\n\n")
	fmt.Printf("-mode=html (or rest)\n\n")
	fmt.Printf("\n\n")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}
	port := flag.Int("port", 4000, "Set Port")
	mode := flag.String("mode", "rest", "Choose 'rest' or 'html'")
	flag.Parse()
	switch *mode {
	case "rest":
		restapi.Start(*port)
	case "html":
		sample_explorer.Start(*port)
	default:
		usage()
	}
}
