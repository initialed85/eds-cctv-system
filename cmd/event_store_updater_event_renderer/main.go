package main

import (
	"flag"
	"github.com/initialed85/eds-cctv-system/pkg/event_store_updater_event_renderer"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func waitForCtrlC() {
	sig := make(chan os.Signal, 2)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	storePath := flag.String("storePath", "", "Path to event store file")
	renderPath := flag.String("renderPath", "", "Path to render folder")

	flag.Parse()

	if *storePath == "" {
		log.Fatal("-storePath cannot be empty")
	}

	if *renderPath == "" {
		log.Fatal("-renderPath cannot be empty")
	}

	log.Printf("creating")

	e, err := event_store_updater_event_renderer.New(*storePath, *renderPath)
	if err != nil {
		log.Fatalf("failed to create Renderer because: %v", err)
	}

	log.Printf("starting")

	e.Start()

	log.Printf("running")

	waitForCtrlC()

	log.Printf("stopping")

	e.Stop()

	log.Printf("stopped")
}
