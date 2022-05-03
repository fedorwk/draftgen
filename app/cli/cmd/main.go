package main

import (
	"log"
	"os"

	"github.com/fedorwk/draftgen/app/cli"
)

func main() {
	err := cli.Run(os.Args[1:])
	if err != nil {
		log.Fatalln(err)
	}
}
