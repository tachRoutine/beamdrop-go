package main

import (
	"fmt"

	"github.com/tachRoutine/ekiliBeam-go/beam"
	"github.com/tachRoutine/ekiliBeam-go/pkg/qr"
)

func main() {
	sharedDir := "./"
	url := beam.StartServer(sharedDir)
	fmt.Println("Server started at", url)
	filename := "qrcode.png"
	err := qr.Generate(url, filename)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}
	fmt.Println("QR code generated and saved to", filename)
}
