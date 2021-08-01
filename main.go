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
	blink(GREEN, 100*time.Millisecond, 3)
	log.Println("READY - WAITING FOR USER INPUT")

	for {
		select {
		case <-onExit:
			{
				//Visual indicator for shutdown
				blink(RED, 100*time.Millisecond, 3)

				log.Println("Disabling LEDs")
				//Make sure that on unexpected / expected shutdowns the LEDs are
				//properly turned off.
				setLEDColor(NONE)
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

		blueLEDPin.Enable()
		captureError := captureImage()
		if captureError != nil {
			//Valid case, as the picture can be blurry or there's no barcode present.
			log.Println("Error capturing barcode image:", captureError)
			blueLEDPin.Disable()
			redLEDPin.Enable()
			return nil
		}
		blueLEDPin.Disable()

		authBarcode, parseError := parseBarcodeImage()
		if parseError != nil {
			//Valid case, as the picture can be blurry or there's no barcode present.
			log.Println("Error parsing barcode:", parseError)
			redLEDPin.Enable()
			return nil
		}

		blueLEDPin.Enable()
		redLEDPin.Enable()

		captureError = captureImage()
		if captureError != nil {
			//Valid case, as the picture can be blurry or there's no barcode present.
			log.Println("Error capturing barcode image:", captureError)
			blueLEDPin.Disable()
			redLEDPin.Enable()
			return nil
		}

		blueLEDPin.Disable()
		blueLEDPin.Disable()

		bottleBarcode, parseError := parseBarcodeImage()
		if parseError != nil {
			//Valid case, as the picture can be blurry or there's no barcode present.
			log.Println("Error parsing barcode:", parseError)
			redLEDPin.Enable()
			return nil
		}

		log.Printf("Authbarcode: %s; Bottlebarcode: %s\n", authBarcode, bottleBarcode)
	}

	return nil
}

type LEDColor int

const (
	NONE LEDColor = iota
	RED
	GREEN
	BLUE
	YELLOW
	CYAN
	MAGENTA
	WHITE
)

func setLEDColor(color LEDColor) {
	switch color {
	case NONE:
		redLEDPin.Disable()
		greenLEDPin.Disable()
		blueLEDPin.Disable()
	case RED:
		redLEDPin.Enable()
		greenLEDPin.Disable()
		blueLEDPin.Disable()
	case GREEN:
		redLEDPin.Disable()
		greenLEDPin.Enable()
		blueLEDPin.Disable()
	case BLUE:
		redLEDPin.Disable()
		greenLEDPin.Disable()
		blueLEDPin.Enable()
	case YELLOW:
		redLEDPin.Enable()
		greenLEDPin.Enable()
		blueLEDPin.Disable()
	case CYAN:
		redLEDPin.Disable()
		greenLEDPin.Enable()
		blueLEDPin.Enable()
	case MAGENTA:
		redLEDPin.Enable()
		greenLEDPin.Disable()
		blueLEDPin.Enable()
	case WHITE:
		redLEDPin.Enable()
		greenLEDPin.Enable()
		blueLEDPin.Enable()
	}
}

func blink(led LEDColor, interval time.Duration, count int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for i := 0; i < count; i++ {
		<-ticker.C
		setLEDColor(led)
		<-ticker.C
		setLEDColor(NONE)
	}
}

const imageTempPath = "/run/user/1000/barcode.jpg"

func captureImage() error {
	raspistillCommand := exec.Command("raspistill", "--timeout", "1500", "--encoding", "jpg", "--output", imageTempPath, "--nopreview", "--quality", "10", "--rotation", "180", "--width", "1200", "--height", "1200")
	log.Println("Executing:", raspistillCommand)
	captureError := raspistillCommand.Run()
	if captureError != nil {
		return fmt.Errorf("error capturing image: %s", captureError)
	}

	return nil
}

func parseBarcodeImage() (string, error) {
	zBarCommand := exec.Command("zbarimg", "-Sdisable", "-Sean13.enable", "-Sposition=false", "-Sx-density=2", "-Sy-density=2", "--raw", "--quiet", imageTempPath)
	log.Println("Executing:", zBarCommand)
	var barcodeBuffer bytes.Buffer
	zBarCommand.Stdout = &barcodeBuffer
	zBarError := zBarCommand.Run()

	if zBarError != nil {
		return "", fmt.Errorf("error parsing image:%s", zBarError)
	}

	//zbar output contains trailing whitespace, therefore we always trim.
	return strings.TrimSpace(barcodeBuffer.String()), nil
}
