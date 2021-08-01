// +build dummy test

package led

import (
	"log"
)

//This file contains dummy-implementations of the functions that interact with
//LEDs, and therefore the GPIO pins of a raspberry pi. This helps quickly
//prototyping and testing new workflows and prevents requiring access to a
//raspberry pi at all times.

// SetLEDColor will change the LEDs color to the passed value. Manual combining
// of colors isn't possible. For example to get white, you'll have to use the
// constant WHITE, instead of combining RED, BLUE and GREEN.
func SetLEDColor(color LEDColor) {
	switch color {
	case NONE:
		log.Println("----- LED: OFF     -----")
	case RED:
		log.Println("----- LED: RED     -----")
	case GREEN:
		log.Println("----- LED: GREEN   -----")
	case BLUE:
		log.Println("----- LED: BLUE    -----")
	case YELLOW:
		log.Println("----- LED: YELLOW  -----")
	case CYAN:
		log.Println("----- LED: CYAN    -----")
	case MAGENTA:
		log.Println("----- LED: MAGENTA -----")
	case WHITE:
		log.Println("----- LED: WHITE   -----")
	}
}
