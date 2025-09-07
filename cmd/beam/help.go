package main

import "fmt"

func Help() string {
	return `beamdrop - A simple file sharing tool with enhanced features

NOTE: YOU NEED TO BE IN THE SAME NETWORK AS THE RECEIVER

Usage:
  beam [options]

Options:
  -dir string
        Directory to share files from (default ".")
  -port int
        Port to run the server on (default from config)
  -password string
        Optional password to protect the server
  -verbose
        Enable verbose logging
  -no-qr
        Disable QR code display
  -h, --help
        Show this help message

Examples:
  beam                                    # Share current directory
  beam -dir="/path/to/files"             # Share specific directory
  beam -port=8080                        # Use custom port
  beam -password="secret123"             # Password protect the server
  beam -verbose -no-qr                   # Verbose mode without QR code

Features:
  • File upload and download
  • Directory browsing
  • File preview (images, text files)
  • Multiple file upload support
  • Optional password protection
  • QR code for easy sharing
  • Real-time statistics`
}

func PrintHelp() {
	fmt.Println(Help())
}
