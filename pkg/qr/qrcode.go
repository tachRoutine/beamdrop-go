package qr

import (
	"os"

	"github.com/skip2/go-qrcode"
)

func Generate(data string, filename string) error {
	// Generate a QR code
	qrCode, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return err
	}

	// Print the QR code
	pngData, err := qrCode.PNG(256)
	if err != nil {
		return err
	}
	// Write PNG data to a file
	file, err := os.Create(filename)
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