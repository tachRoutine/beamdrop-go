package main

import (
	"flag"
	"fmt"

	"github.com/tachRoutine/ekiliBeam-go/beam"
	"github.com/tachRoutine/ekiliBeam-go/pkg/qr"
)

func main() {
	
	sharedDir := flag.String("dir", ".", "Directory to share files from")
	flag.Parse()
	if (*sharedDir == "") {
		fmt.Println("Shared directory is required")
		return
	}
	url := beam.StartServer(*sharedDir)
	fmt.Println("Server started at", url)
	filename := "qrcode.png"
	err := qr.Generate(url, filename)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}
	fmt.Println("QR code generated and saved to", filename)
}
