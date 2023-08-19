package main

import (
	"os"

	"github.com/cbguder/books/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
