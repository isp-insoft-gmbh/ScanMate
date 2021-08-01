// +build dummy

package gpio

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// InputPinDummy implements gpio.InputPin.
type InputPinDummy struct {
	port  uint64
	Value bool
}

// OutputPinDummy implements gpio.OutputPin.
type OutputPinDummy struct {
	*InputPinDummy
}

// Enable writes "1" into the value. This can for example be used to turn on
// an LED.
func (outputPin *OutputPinDummy) Enable() error {
	log.Println("Return Error on 'OutputPin.Enable' (y/n)")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		scanErr := scanner.Err()
		panic(scanErr)
	}

	if strings.EqualFold(scanner.Text(), "y") {
		return fmt.Errorf("error enabling output pin %d", outputPin.port)
	}

	outputPin.Value = true
	return nil
}

// Disable writes "0" into the value. This can for example be used to turn off
// an LED.
func (outputPin *OutputPinDummy) Disable() error {
	log.Println("Return Error on 'OutputPin.Disable' (y/n)")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		scanErr := scanner.Err()
		panic(scanErr)
	}

	if strings.EqualFold(scanner.Text(), "y") {
		return fmt.Errorf("error disabling output pin %d", outputPin.port)
	}

	outputPin.Value = false
	return nil
}

// GetValue return the current value the port is set to.
func (inputPin *InputPinDummy) GetValue() (bool, error) {
	log.Println("Return Error on 'InputPin.GetValue' (y/n)")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		scanErr := scanner.Err()
		panic(scanErr)
	}

	if strings.EqualFold(scanner.Text(), "y") {
		return false, fmt.Errorf("error getting value from input pin %d", inputPin.port)
	}

	log.Println("Return true on 'Input.GetValue' (y/n)")
	if !scanner.Scan() {
		scanErr := scanner.Err()
		panic(scanErr)
	}

	return strings.EqualFold(scanner.Text(), "y"), nil
}

// SetupGPIOOutputPort creates a dummy output port. You can fake the current
// value and the next errors upon calling functions.
func SetupGPIOOutputPort(port uint64) (*OutputPinDummy, error) {
	return &OutputPinDummy{InputPinDummy: &InputPinDummy{port: port}}, nil
}

// SetupGPIOOutputPort creates a dummy output port. You can fake the current
// value and the next errors upon calling functions.
func SetupGPIOInputPort(port uint64) (*InputPinDummy, error) {
	return &InputPinDummy{port: port}, nil
}
