package main

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

func main() {
	data := "ekilie.com"
	filename := "qrcode.png"
	err := generateQRCode(data, filename)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}
	fmt.Println("QR code generated and saved to", filename)
	
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
