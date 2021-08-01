package led

import (
	"log"
	"time"

	"github.com/isp-insoft-gmbh/scanmate/gpio"
)

// LEDColor is type to represents all possible LED states.
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

var (
	redLEDPin   gpio.OutputPin
	greenLEDPin gpio.OutputPin
	blueLEDPin  gpio.OutputPin
)

// InitLEDs sets up the led package for usage. Calling any package functions
// before initilization will result in errors.
func InitLEDs(gpioPinRed, gpioPinGreen, gpioPinBlue uint64) {
	var setupError error
	redLEDPin, setupError = gpio.SetupGPIOOutputPort(gpioPinRed)
	if setupError != nil {
		log.Fatalf("Error setting up pin 7: %s\n", setupError)
	}
	greenLEDPin, setupError = gpio.SetupGPIOOutputPort(gpioPinGreen)
	if setupError != nil {
		log.Fatalf("Error setting up pin 8: %s\n", setupError)
	}
	blueLEDPin, setupError = gpio.SetupGPIOOutputPort(gpioPinBlue)
	if setupError != nil {
		log.Fatalf("Error setting up pin 25: %s\n", setupError)
	}
}

// Blink lets the LED flash up in the desired color for a specified amount
// of times.
// For example, calling this:
//     blink(RED, 100 * time.Millisecond, 2)
// Will result in:
//     LEAVE CURRENT LED STATE UNCHANGED
//     WAIT(100ms)
//     LED RED
//     WAIT(100ms)
//     LED OFF
//     WAIT(100ms)
//     LED RED
//     WAIT(100ms)
//     LED OFF
func Blink(led LEDColor, interval time.Duration, count int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for i := 0; i < count; i++ {
		<-ticker.C
		SetLEDColor(led)
		<-ticker.C
		SetLEDColor(NONE)
	}
}
