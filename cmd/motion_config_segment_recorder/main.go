package main

import (
	"flag"
	"github.com/initialed85/eds-cctv-system/pkg/motion_config_segment_recorder"
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

	configPath := flag.String("configPath", "", "path to look for motion configs")
	destinationPath := flag.String("destinationPath", "", "path to output video segments to")
	duration := flag.Int("duration", 0, "")

	flag.Parse()

	if *configPath == "" {
		log.Fatal("-configPath cannot be empty")
	}

	if *destinationPath == "" {
		log.Fatal("-destinationPath cannot be empty")
	}

	if *duration <= 0 {
		log.Fatal("-duration cannot be empty or negative")
	}

	m, err := motion_config_segment_recorder.New(*configPath, *destinationPath, *duration)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Start()
	if err != nil {
		log.Fatal(err)
	}

	waitForCtrlC()

	m.Stop()
}
