package main

import (
	"log"
)

func main() {
	server := initSWHandler()
	defer server.Shutdown()
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
