package main

import (
	"flag"
	"fmt"

	"github.com/tachRoutine/beamdrop-go/beam"
)

func main() {

	sharedDir := flag.String("dir", ".", "Directory to share files from")
	help := flag.Bool("h", false, "Show help message")
	flag.Parse()
	// if flag.NArg() > 0 {
	// 	PrintHelp()
	// 	return
	// }
	if *sharedDir == "" {
		fmt.Println("Shared directory is required")
		return
	}
	if *help {
		PrintHelp()
		return
	}

	beam.StartServer(*sharedDir)
}
