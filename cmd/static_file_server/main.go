package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("port", 0, "Port to listen on")
	path := flag.String("path", "", "Path to serve")

	flag.Parse()

	if *port <= 0 {
		log.Fatal("-port cannot be empty or negative")
	}

	if *path == "" {
		log.Fatal("-path cannot be empty")
	}

	err := http.ListenAndServe(fmt.Sprintf(":%v", *port), http.FileServer(http.Dir(*path)))
	if err != nil {
		log.Fatal(err)
	}
}
