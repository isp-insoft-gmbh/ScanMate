package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("START")
	defer fmt.Println("END")
	onExit := make(chan os.Signal)
	go func() {
		<-onExit
		os.Exit(0)
	}()
	signal.Notify(onExit, syscall.SIGINT)
	for {
		err := logic()
		if err != nil {
			log.Printf("Fatal error: %s\n", err)
			break
		}
	}
}

// logic contains the application logic requried for the workflow of buying
// a drink. An error is only returned, if a non recoverable situation has
// been encountered.
func logic() error {
	buttonValue, buttonError := os.ReadFile("/sys/class/gpio/gpio17/value")
	if buttonError != nil {
		return buttonError
	}

	//Values can be 0 and 1.
	if buttonValue[0] == '1' {
		log.Println("BUTTON WAS PRESSED")
		barcode, barcodeError := readBarcode()
		if barcodeError == nil {
			log.Println("The barcode is", barcode)
		} else {
			//Valid case, as the picture can be blurry or there's no barcode present.
			log.Println("Error reading barcode:", barcodeError)
		}
	}

	return nil
}

// readBarcode takes an image via raspistill and tries parsing it with zbarimg.
func readBarcode() (string, error) {
	imageTempPath := "/run/user/1000/barcode.jpg"

	raspistillCommand := exec.Command("raspistill", "--timeout", "1", "--encoding", "jpg", "--output", imageTempPath, "--nopreview", "--quality", "10", "--rotation", "180")
	log.Println("Executing:", raspistillCommand)
	captureError := raspistillCommand.Run()
	if captureError != nil {
		return "", fmt.Errorf("error capturing image: %s", captureError)
	}

	zBarCommand := exec.Command("zbarimg", "-Sdisable", "-Sean13.enable", "--raw", "--quiet", imageTempPath)
	log.Println("Executing:", zBarCommand)
	var barcodeBuffer bytes.Buffer
	zBarCommand.Stdout = &barcodeBuffer
	zBarError := zBarCommand.Run()
	if zBarError != nil {
		return "", fmt.Errorf("error parsing image:%s", zBarError)
	}
	return barcodeBuffer.String(), nil
}
