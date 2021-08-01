// +build test

package gpio

import (
	"errors"
	"testing"
)

func Test_InputPinDummy(t *testing.T) {
	//Dummy-Check to see if interface is implemented
	var _ InputPin = new(InputPinDummy)

	inputPin, err := SetupGPIOInputPort(1)
	if err != nil {
		t.Errorf("Failed setting up pin: %s", err)
	}

	if inputPin == nil {
		t.Error("input pin was nil")
	}

	if inputPin.Value {
		t.Error("Initial input pin value should be false")
	}

	inputPin.Value = true

	value, err := inputPin.GetValue()
	if err != nil {
		t.Errorf("Error getting value: %s", err)
	}
	if !value {
		t.Error("input pin value after update should be true")
	}

	expectedError := errors.New("err")
	inputPin.NextGetValueError = expectedError

	value, err = inputPin.GetValue()
	if err != expectedError {
		t.Errorf("Expected error '%s', but got: %s", expectedError, err)
	}

	if value {
		t.Error("Since an error was returned, the value should be false34")
	}
}

func Test_OutputPinDummy(t *testing.T) {
	//Dummy-Check to see if interface is implemented
	var _ OutputPin = new(OutputPinDummy)
}
