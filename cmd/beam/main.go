package main

import (
	"flag"
	"fmt"

	"github.com/tachRoutine/beamdrop-go/beam"
	"github.com/tachRoutine/beamdrop-go/pkg/qr"
)

func main() {
	
	sharedDir := flag.String("dir", ".", "Directory to share files from")
	help := flag.Bool("h", false, "Show help message")
	flag.Parse()
	// if flag.NArg() > 0 {
	// 	PrintHelp()
	// 	return
	// }
	if (*sharedDir == "") {
		fmt.Println("Shared directory is required")
		return
	}
	if *help {
		PrintHelp()
		return
	}
	url := beam.StartServer(*sharedDir)
	filename := "qrcode.png"
	err := qr.Generate(url, filename)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}
	fmt.Println("QR code generated and saved to", filename)
}
