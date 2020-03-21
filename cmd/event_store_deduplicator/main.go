package main

import (
	"flag"
	"github.com/initialed85/eds-cctv-system/pkg/event_store_deduplicator"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	sourcePath := flag.String("sourcePath", "", "source event store")
	destinationPath := flag.String("destinationPath", "", "destination event store")

	flag.Parse()

	if *sourcePath == "" {
		log.Fatal("-sourcePath cannot be empty")
	}

	if *destinationPath == "" {
		log.Fatal("-destinationPath cannot be empty")
	}

	d := event_store_deduplicator.New(*sourcePath, *destinationPath)

	err := d.Deduplicate()
	if err != nil {
		log.Fatal(err)
	}
}
