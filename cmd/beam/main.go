package main

import (
	"fmt"

	"github.com/tachRoutine/ekiliBeam-go/pkg/qr"
)

func main() {
	sharedDir := "./"
	data := "ekilie.com"
	filename := "qrcode.png"
	err := qr.Generate(data, filename)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}
	fmt.Println("QR code generated and saved to", filename)

}
