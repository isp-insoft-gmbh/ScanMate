// +build test

package gpio

// InputPinDummy implements gpio.InputPin.
type InputPinDummy struct {
	port              uint64
	NextGetValueError error
	Value             bool
}

// OutputPinDummy implements gpio.OutputPin.
type OutputPinDummy struct {
	*InputPinDummy
	NextEnableError  error
	NextDisableError error
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

// Enable writes "1" into the value. This can for example be used to turn on
// an LED.
func (outputPin *OutputPinDummy) Enable() error {
	nextError := outputPin.NextEnableError
	if nextError != nil {
		outputPin.NextEnableError = nil
		return nextError
	}

	outputPin.Value = true
	return nil
}

// Disable writes "0" into the value. This can for example be used to turn off
// an LED.
func (outputPin *OutputPinDummy) Disable() error {
	nextError := outputPin.NextDisableError
	if nextError != nil {
		outputPin.NextDisableError = nil
		return nextError
	}

	outputPin.Value = false
	return nil
}

// GetValue return the current value the port is set to.
func (inputPin *InputPinDummy) GetValue() (bool, error) {
	nextError := inputPin.NextGetValueError
	if nextError != nil {
		inputPin.NextGetValueError = nil
		return false, nextError
	}

	return inputPin.Value, nil
}
