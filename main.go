package main

import (
	"bytes"
	"fmt"
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
		buttonValue, buttonError := os.ReadFile("/sys/class/gpio/gpio17/value")
		if buttonError == nil {
			if buttonValue[0] == '1' {
				fmt.Println("BUTTON WAS PRESSED")
				barcode, barcodeError := readBarcode()
				if barcodeError == nil {
					fmt.Println("The barcode is", barcode)
				} else {
					fmt.Println(barcodeError)	
					//TODO: was wenn barcode nicht gelesen wrden konnte?
				}
				
			}
		} else {
			fmt.Println(buttonError)
			os.Exit(-1)
		}

	}
}

func readBarcode() (string,error) {
	imageTempPath := "/run/user/1000/barcode.png"
	captureError := exec.Command("raspistill", "--timeout", "1" , "--output", imageTempPath, "--nopreview", "--quality", "10").Run()
	
	if captureError != nil {
		return "", captureError
	}

	zBarCommand := exec.Command("zbarimg", "-Sdisable", "-Sean13.enable", "--raw", "--quiet", imageTempPath)
	var barcodeBuffer bytes.Buffer
	zBarCommand.Stdout = &barcodeBuffer
	zBarError := zBarCommand.Run()
	if zBarError != nil {
		return "", zBarError
	}
	return barcodeBuffer.String(), nil
}
