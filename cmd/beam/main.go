package main

import (
	"flag"

	"github.com/tachRoutine/beamdrop-go/beam"
	"github.com/tachRoutine/beamdrop-go/pkg/logger"
)

func main() {
	logger.Info("Starting beamdrop application")
	
	sharedDir := flag.String("dir", ".", "Directory to share files from")
	help := flag.Bool("h", false, "Show help message")
	flag.Parse()
	
	if flag.NArg() > 0 {
		logger.Debug("Extra arguments provided, showing help")
		PrintHelp()
		return
	}
	if *sharedDir == "" {
		logger.Error("Shared directory is required")
		return
	}
	if *help {
		logger.Debug("Help flag provided, showing help")
		PrintHelp()
		return
	}

	logger.Info("Starting server with shared directory: %s", *sharedDir)
	beam.StartServer(*sharedDir)
}
