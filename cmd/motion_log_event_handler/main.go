package main

import (
	"github.com/initialed85/eds-cctv-system/pkg/motion_log_event_handler"
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

	filePath := flag.String("filePath", "", "Path to log file for Motion STDOUT and STDERR")

	flag.Parse()

	if *filePath == "" {
		log.Fatal("filePath cannot be empty")
	}

	log.Printf("creating")

	m, err := motion_log_event_handler.New(*filePath)
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
