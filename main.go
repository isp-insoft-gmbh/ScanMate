// +build !test

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isp-insoft-gmbh/scanmate/barcode"
	"github.com/isp-insoft-gmbh/scanmate/gpio"
	"github.com/isp-insoft-gmbh/scanmate/led"
)

var (
	buttonPin gpio.InputPin
)

func main() {
	log.Println("START")
	defer log.Println("END")

	//To supporting reacting to Ctrl-C, to allow gracefully killing the application.
	onExit := make(chan os.Signal, 1)
	signal.Notify(onExit, syscall.SIGINT, syscall.SIGTERM)

	var setupError error
	buttonPin, setupError = gpio.SetupGPIOInputPort(17)
	if setupError != nil {
		log.Fatalf("Error setting up pin 17: %s\n", setupError)
	}

	led.InitLEDs(7, 8, 25)

	//Visual indicator for being ready
	led.Blink(led.GREEN, 100*time.Millisecond, 3)
	log.Println("READY - WAITING FOR USER INPUT")

	for {
		select {
		case <-onExit:
			{
				//Visual indicator for shutdown
				led.Blink(led.RED, 100*time.Millisecond, 3)

				log.Println("Disabling LEDs")
				//Make sure that on unexpected / expected shutdowns the LEDs are
				//properly turned off.
				led.SetLEDColor(led.NONE)
				log.Println("END")
				os.Exit(0)
			}
		default:
			{
				err := logic()
				if err != nil {
					log.Printf("Fatal error: %s\n", err)
					os.Exit(1)
				}
			}
		}
	}
}

// logic contains the application logic requried for the workflow of buying
// a drink. An error is only returned, if a non recoverable situation has
// been encountered.
func logic() error {
	buttonValue, buttonError := buttonPin.GetValue()
	if buttonError != nil {
		return buttonError
	}

	//Values can be 0 and 1.
	if buttonValue {
		log.Println("BUTTON WAS PRESSED")
		log.Println("Waiting for authentication barcode.")

		led.SetLEDColor(led.BLUE)
		captureError := barcode.CaptureImage()
		if captureError != nil {
			//Valid case, as the picture can be blurry or there's no barcode present.
			log.Println("Error capturing barcode image:", captureError)
			led.SetLEDColor(led.RED)
			return nil
		}
		led.SetLEDColor(led.NONE)

		authBarcode, parseError := barcode.ParseBarcodeImage()
		if parseError != nil {
			//Valid case, as the picture can be blurry or there's no barcode present.
			log.Println("Error parsing barcode:", parseError)
			led.SetLEDColor(led.RED)
			return nil
		}

		led.SetLEDColor(led.MAGENTA)

		captureError = barcode.CaptureImage()
		if captureError != nil {
			//Valid case, as the picture can be blurry or there's no barcode present.
			log.Println("Error capturing barcode image:", captureError)
			led.SetLEDColor(led.RED)
			return nil
		}

		led.SetLEDColor(led.NONE)

		bottleBarcode, parseError := barcode.ParseBarcodeImage()
		if parseError != nil {
			//Valid case, as the picture can be blurry or there's no barcode present.
			log.Println("Error parsing barcode:", parseError)
			led.SetLEDColor(led.RED)
			return nil
		}

		log.Printf("Authbarcode: %s; Bottlebarcode: %s\n", authBarcode, bottleBarcode)
	}

	return nil
}
