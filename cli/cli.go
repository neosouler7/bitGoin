package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/neosouler7/bitGoin/explorer"
	"github.com/neosouler7/bitGoin/rest"
)

func usage() {
	fmt.Printf("Welcome to bitGoin :)\n")
	fmt.Printf("Please use the following commands:\n\n")
	fmt.Printf("-port:    Set the PORT of the server\n")
	fmt.Printf("-mode:    Choose between 'html' and 'rest'\n")
	os.Exit(1)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
}
