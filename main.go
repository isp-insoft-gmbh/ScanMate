// +build !test

package main

import (
	"errors"
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
				if buttonPressed, buttonError := buttonPin.GetValue(); buttonError != nil {
					log.Printf("Fatal error: %s\n", buttonError)
					os.Exit(1)
				} else if buttonPressed {
					log.Println("BUTTON WAS PRESSED")
					//We wait 1 second, so the user is ready to prepare his ID-card in time.
					time.Sleep(1 * time.Second)

					if logicErr := purchaseWorkflow(); logicErr != nil &&
						//This error is somewhat expected behaviour and should't cause a crash.
						logicErr != ErrNoAttemptsLeft {

						log.Printf("Fatal error: %s\n", logicErr)
						os.Exit(1)
					}
				}
			}
		}
	}
}

var ErrNoAttemptsLeft error = errors.New("no more attempts to take an image and parse are left")

const maxTries = 10

func captureBarcode(try int, waitingColor led.LEDColor) (string, error) {
	led.SetLEDColor(waitingColor)
	//Waiting again, so the user knows the app is ready and the barcode should
	//be properly positioned.
	time.Sleep(1500 * time.Millisecond)

	captureError := barcode.CaptureImage()
	if captureError != nil {
		log.Println("Error capturing barcode image:", captureError)
		led.Blink(led.RED, 125*time.Millisecond, 4)
		return "", captureError
	}

	//Turning off LED, so the user can stop holding up their arm.
	led.SetLEDColor(led.NONE)

	parsedBarcode, parseError := barcode.ParseBarcodeImage()
	if parseError != nil {
		//Valid case, as the picture can be blurry or there's no barcode present.
		if parseError == barcode.ErrNoBarcodeFound {
			led.Blink(led.YELLOW, 125*time.Millisecond, 4)
			log.Println("Barcode couldn't be recognized")
			//After at max 10 failed tries, we give up on the user, lost cause.
			if try == maxTries {
				return "", ErrNoAttemptsLeft
			}

			log.Println("Retrying. Attempts left:", maxTries-try)
			return captureBarcode(try+1, waitingColor)
		}

		//Fatal error, since the image couldn't be parsed at all or
		//the user used up their 10 tries.
		log.Println("Error parsing barcode:", parseError)
		led.Blink(led.RED, 125*time.Millisecond, 4)
		return "", parseError
	}

	return parsedBarcode, nil
}

// purchaseWorkflow contains the application purchaseWorkflow requried for the
// workflow of buying a drink. An error is only returned, if a non recoverable
// situation has been encountered.
func purchaseWorkflow() error {
	//Make sure that no failure keeps an LED on
	defer led.SetLEDColor(led.NONE)

	log.Println("Waiting for authentication barcode.")
	authBarcode, parseError := captureBarcode(1, led.BLUE)
	if parseError != nil {
		return parseError
	}

	log.Println("Waiting for bottle barcode.")
	bottleBarcode, parseError := captureBarcode(1, led.MAGENTA)
	if parseError != nil {
		return parseError
	}

	led.SetLEDColor(led.GREEN)
	log.Printf("Authbarcode: %s; Bottlebarcode: %s\n", authBarcode, bottleBarcode)

	//The user now has 3 seconds of time to look at their success.
	time.Sleep(3 * time.Second)

	return nil
}
