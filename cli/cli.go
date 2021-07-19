package cli

import (
	"flag"
	"fmt"
	explorer "learngo/github.com/nomadcoders/explorer/templates"
	"learngo/github.com/nomadcoders/rest"
	"os"
)

func usage() {
	fmt.Printf("Welcome to 노마드 코인\n\n")
	fmt.Printf("Please use the following commands:\n\n")
	fmt.Printf("-port:   	 Set port of the server\n")
	fmt.Printf("-mode:       Choose among 'html' and 'rest' and 'both'\n")
	os.Exit(1)

}

func Start() {
	fmt.Println(os.Args)
	if len(os.Args) == 1 {
		usage()
	}
	port1 := flag.Int("port1", 4000, "Set port1 of the server")
	port2 := flag.Int("port2", 5000, "Set port2 of the server")
	mode := flag.String("mode", "rest", "Choose among 'html' and 'rest' and 'both")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port1)
	case "html":
		explorer.Start(*port1)
	case "both":
		go rest.Start(*port1)
		explorer.Start(*port2)
	default:
		usage()
	}

	fmt.Println(*port1, *port2, *mode)
}
