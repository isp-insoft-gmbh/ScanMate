package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/isp-insoft-gmbh/scanmate/gpio"
)

var (
	buttonPin   *gpio.InputPin
	redLEDPin   *gpio.OutputPin
	greenLEDPin *gpio.OutputPin
	blueLEDPin  *gpio.OutputPin
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
	redLEDPin, setupError = gpio.SetupGPIOOutputPort(7)
	if setupError != nil {
		log.Fatalf("Error setting up pin 7: %s\n", setupError)
	}
	greenLEDPin, setupError = gpio.SetupGPIOOutputPort(8)
	if setupError != nil {
		log.Fatalf("Error setting up pin 8: %s\n", setupError)
	}
	blueLEDPin, setupError = gpio.SetupGPIOOutputPort(25)
	if setupError != nil {
		log.Fatalf("Error setting up pin 25: %s\n", setupError)
	}

	//Visual indicator for being ready
	blink(greenLEDPin, 100*time.Millisecond, 3)
	log.Println("READY - WAITING FOR USER INPUT")

	for {
		select {
		case <-onExit:
			{
				//Visual indicator for shutdown
				blink(redLEDPin, 100*time.Millisecond, 3)

				log.Println("Disabling LEDs")
				//Make sure that on unexpected / expected shutdowns the LEDs are
				//properly turned off.
				disableLEDs()
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

func disableLEDs() {
	blueLEDPin.Disable()
	redLEDPin.Disable()
	greenLEDPin.Disable()
}

func blink(led *gpio.OutputPin, interval time.Duration, count int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for i := 0; i < count; i++ {
		<-ticker.C
		led.Enable()
		<-ticker.C
		led.Disable()
	}
}

// readBarcode takes an image via raspistill and tries parsing it with zbarimg.
func readBarcode() (string, error) {
	blueLEDPin.Enable()

	imageTempPath := "/run/user/1000/barcode.jpg"

	raspistillCommand := exec.Command("raspistill", "--timeout", "1500", "--encoding", "jpg", "--output", imageTempPath, "--nopreview", "--quality", "10", "--rotation", "180")
	log.Println("Executing:", raspistillCommand)
	captureError := raspistillCommand.Run()
	if captureError != nil {
		return "", fmt.Errorf("error capturing image: %s", captureError)
	}

	blueLEDPin.Disable()

	zBarCommand := exec.Command("zbarimg", "-Sdisable", "-Sean13.enable", "--raw", "--quiet", imageTempPath)
	log.Println("Executing:", zBarCommand)
	var barcodeBuffer bytes.Buffer
	zBarCommand.Stdout = &barcodeBuffer
	zBarError := zBarCommand.Run()

	//Disable both red and green, since we are too lazy to write two code pieces.
	//FIXME This will cause unnecessary blocking, making the actual purchasing process slower.
	//However, just making a background thread will cause multi threading problems.
	defer func() {
		<-time.NewTimer(3 * time.Second).C
		greenLEDPin.Disable()
		redLEDPin.Disable()
	}()

	if zBarError != nil {
		redLEDPin.Enable()
		return "", fmt.Errorf("error parsing image:%s", zBarError)
	}

	greenLEDPin.Enable()
	//zbar output contains trailing whitespace, therefore we always trim.
	return strings.TrimSpace(barcodeBuffer.String()), nil
}
