package main

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
	"github.com/tachRoutine/ekiliBeam-go.git/pkg/qr"
)

func main() {
	data := "ekilie.com"
	filename := "qrcode.png"
	err := qr.Generate(data, filename)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}
	fmt.Println("QR code generated and saved to", filename)
	
}
