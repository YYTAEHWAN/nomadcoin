package main

import (
	"github.com/nomadcoders/cli"
	"github.com/nomadcoders/db"
)

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()
}

/*
package hello

import "rsc.io/quote"

func Hello() string {
	return quote.Hello()
}
*/
