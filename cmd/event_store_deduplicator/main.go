package main

import (
	"flag"
	"github.com/initialed85/eds-cctv-system/pkg/event_store_deduplicator"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	path := flag.String("path", "", "Path to serve")

	flag.Parse()

	if *path == "" {
		log.Fatal("-path cannot be empty")
	}

	d := event_store_deduplicator.New(*path)

	err := d.Deduplicate()
	if err != nil {
		log.Fatal(err)
	}
}
