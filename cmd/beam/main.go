package main

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

func main() {
	
}

func generateQRCode(data string, filename string) error {
	// Generate a QR code
	qrCode, err := qrcode.New("ekilie.com", qrcode.Medium)
	if err != nil {
		return err
	}

	// Print the QR code
	pngData, err := qrCode.PNG(256)
	if err != nil {
		return err
	}
	// Write PNG data to a file
	file, err := os.Create("qrcode.png")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(pngData)
	if err != nil {
		return err
	}
	return nil
}
