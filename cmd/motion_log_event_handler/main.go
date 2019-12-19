package main

import (
	"flag"
	"github.com/initialed85/eds-cctv-system/internal/event_store"
	"github.com/initialed85/eds-cctv-system/pkg/motion_log_event_handler"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	flag.Parse()

	if *filePath == "" {
		log.Fatal("filePath cannot be empty")
	}

	if *jsonLinesPath == "" {
		log.Fatal("jsonLinesPath cannot be empty")
	}

	callback := func(timestamp time.Time, highResImagePath, lowResImagePath, highResVideoPath, lowResVideoPath string) error {
		event := event_store.Event{
			Timestamp:        timestamp,
			HighResImagePath: highResImagePath,
			LowResImagePath:  lowResImagePath,
			HighResVideoPath: highResVideoPath,
			LowResVideoPath:  lowResVideoPath,
		}

		err := event_store.WriteJSONLine(event, *jsonLinesPath)
		if err != nil {
			log.Printf("failed to write JSONLine because: %v", err)
		}

		return nil
	}

	log.Printf("creating")

	m, err := motion_log_event_handler.New(*filePath, callback)
	if err != nil {
		log.Fatalf("failed to create MotionLogEventHandler because: %v", err)
	}

	log.Printf("starting")

	m.Start()

	log.Printf("running")

	waitForCtrlC()

	log.Printf("stopping")

	err = m.Stop()
	if err != nil {
		log.Fatalf("failed to stop MotionLogEventHandler because: %v", err)
	}

	log.Printf("stopped")
}
