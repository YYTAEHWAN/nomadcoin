package main

import (
	"learngo/github.com/nomadcoders/cli"
	"learngo/github.com/nomadcoders/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
