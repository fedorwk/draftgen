package main

import (
	"log"

	"github.com/fedorwk/draftgen/app/server"
)

func main() {
	log.Fatalln(server.Run())
}
