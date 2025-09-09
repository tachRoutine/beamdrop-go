package main

import "github.com/tachRoutine/beamdrop-go/pkg/logger"

func Help() string {
	return `beamdrop - A simple file sharing tool

NOTE: YOU NEED TO BE IN THE SAME NETWORK AS THE RECEIVER

Usage:
  beam [options]

Options:
  -dir string
		Directory to share files from (default ".")
  -h, --help
  --no-qr 
  		Disable QR code generation`
}

func PrintHelp() {
	logger.Info(Help())
}
