package main

import (
	"flag"
	"github.com/initialed85/eds-cctv-system/pkg/event_api"
	"github.com/initialed85/eds-cctv-system/pkg/segment_folder_event_handler"
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

	folderPath := flag.String("folderPath", "", "Path to segment folder")
	storePath := flag.String("storePath", "", "Path to event store file")
	port := flag.Int("port", 0, "Port to listen on")

	flag.Parse()

	if *folderPath == "" {
		log.Fatal("-folderPath cannot be empty")
	}

	if *storePath == "" {
		log.Fatal("-storePath cannot be empty")
	}

	if *port <= 0 {
		log.Fatal("-port cannot be empty or negative")
	}

	log.Printf("creating")

	a := event_api.New(*storePath, *port)

	s, err := segment_folder_event_handler.New(*folderPath, a.AddEvent)
	if err != nil {
		log.Fatalf("failed to create Handler because: %v", err)
	}

	log.Printf("starting")

	a.Start()

	s.Start()

	log.Printf("running")

	waitForCtrlC()

	log.Printf("stopping")

	a.Stop()

	err = s.Stop()
	if err != nil {
		log.Fatalf("failed to stop Handler because: %v", err)
	}

	log.Printf("stopped")
}
