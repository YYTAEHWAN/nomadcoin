package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("Welcome to 노마드 코인\n\n")
	fmt.Printf("Please use the following commands:\n\n")
	fmt.Printf("explorer:   Start the HTML Explorer\n")
	fmt.Printf("rest:       Start the REST API (recommanded)\n")
	os.Exit(1)

}
func main() {

	if len(os.Args) < 2 {
		usage()
	}

}
