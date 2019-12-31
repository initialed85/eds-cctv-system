package main

import (
	"flag"
	"github.com/initialed85/eds-cctv-system/pkg/event_api"
	"github.com/initialed85/eds-cctv-system/pkg/motion_log_event_handler"
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

	filePath := flag.String("filePath", "", "Path to log file for Motion STDOUT and STDERR")
	jsonLinesPath := flag.String("jsonLinesPath", "", "Path to JSONLines output file")
	port := flag.Int("port", 0, "Port to listen on")

	flag.Parse()

	if *filePath == "" {
		log.Fatal("-filePath cannot be empty")
	}

	if *jsonLinesPath == "" {
		log.Fatal("-jsonLinesPath cannot be empty")
	}

	if *port <= 0 {
		log.Fatal("-port cannot be empty or negative")
	}

	log.Printf("creating")

	a := event_api.New(*jsonLinesPath, *port)

	m, err := motion_log_event_handler.New(*filePath, a.AddEvent)
	if err != nil {
		log.Fatalf("failed to create Handler because: %v", err)
	}

	log.Printf("starting")

	a.Start()

	m.Start()

	log.Printf("running")

	waitForCtrlC()

	log.Printf("stopping")

	a.Stop()

	err = m.Stop()
	if err != nil {
		log.Fatalf("failed to stop Handler because: %v", err)
	}

	log.Printf("stopped")
}
