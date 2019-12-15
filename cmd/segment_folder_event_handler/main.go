package main

import (
	"eds-cctv-system/pkg/segment_folder_event_handler"
	"flag"
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

	flag.Parse()

	if *folderPath == "" {
		log.Fatal("folderPath cannot be empty")
	}

	log.Printf("creating")

	s, err := segment_folder_event_handler.New(*folderPath)
	if err != nil {
		log.Fatalf("failed to create SegmentFolderEventHandler because: %v", err)
	}

	log.Printf("starting")

	s.Start()

	log.Printf("running")

	waitForCtrlC()

	log.Printf("stopping")

	err = s.Stop()
	if err != nil {
		log.Fatalf("failed to stop SegmentFolderEventHandler because: %v", err)
	}

	log.Printf("stopped")
}
