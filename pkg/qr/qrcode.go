package qr

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

// Generate a QR code and save it to a file
func Generate(data string, filename string) error {
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
	file, err := os.Create("./" + filename)
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

func ShowQrCode(url string) {
	qrCode, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		fmt.Println("Error creating QR code for terminal:", err)
		return
	}
	fmt.Println("QR code for", url, ":")
	fmt.Println(qrCode.ToSmallString(false))
}