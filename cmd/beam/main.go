package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tachRoutine/beamdrop-go/beam"
)

func main() {
	sharedDir := flag.String("dir", ".", "Directory to share files from")
	port := flag.Int("port", 0, "Port to run the server on (default from config)")
	password := flag.String("password", "", "Optional password to protect the server")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	noQR := flag.Bool("no-qr", false, "Disable QR code display")
	help := flag.Bool("h", false, "Show help message")

	flag.Parse()

	if flag.NArg() > 0 {
		PrintHelp()
		return
	}
	if *sharedDir == "" {
		fmt.Println("Shared directory is required")
		return
	}
	if *help {
		PrintHelp()
		return
	}

	// Check if directory exists
	if _, err := os.Stat(*sharedDir); os.IsNotExist(err) {
		fmt.Printf("Error: Directory '%s' does not exist\n", *sharedDir)
		return
	}

	config := beam.ServerConfig{
		SharedDir: *sharedDir,
		Port:      *port,
		Password:  *password,
		Verbose:   *verbose,
		NoQR:      *noQR,
	}

	beam.StartServer(config)
}
