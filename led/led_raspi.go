// +build !dummy,!test

package led

// SetLEDColor will change the LEDs color to the passed value. Manual combining
// of colors isn't possible. For example to get white, you'll have to use the
// constant WHITE, instead of combining RED, BLUE and GREEN.
func SetLEDColor(color LEDColor) {
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
