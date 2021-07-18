package gpio

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// InputPin allows to read the value of the mapped GPIO pin.
// GetValue() can be used for this.
type InputPin struct {
	valueFile string
}

// OutputPin allows to read the value of the mapped GPIO pin or write it.
// The respective functions are Enable() and Disable().
type OutputPin struct {
	*InputPin
}

// Enable writes "1" into the value. This can for example be used to turn on
// an LED.
func (outputPin *OutputPin) Enable() error {
	return os.WriteFile(outputPin.valueFile, []byte("1"), 0)
}

// Disable writes "0" into the value. This can for example be used to turn off
// an LED.
func (outputPin *OutputPin) Disable() error {
	return os.WriteFile(outputPin.valueFile, []byte("0"), 0)
}

// GetValue return the current value the port is set to.
func (inputPin *InputPin) GetValue() (bool, error) {
	data, readError := os.ReadFile(inputPin.valueFile)
	if readError != nil {
		return false, readError
	}
	return len(data) > 0 && data[0] == '1', nil
}

// SetupGPIOOutputPort exports the desired port and sets the direction to
// output. The accepted port numbers differ between Raspberry Pi models.
func SetupGPIOOutputPort(port uint64) (*OutputPin, error) {
	setupError := setupGPIOPort(port, "out")
	if setupError != nil {
		return nil, setupError
	}

	outputPin := &OutputPin{
		InputPin: &InputPin{
			valueFile: fmt.Sprintf("/sys/class/gpio/gpio%d/value", port),
		},
	}

	//Initially disable port.
	gpioPortDisableError := outputPin.Disable()
	if gpioPortDisableError != nil {
		return nil, fmt.Errorf("error disabling output for port %d: %s", port, gpioPortDisableError)
	}

	return outputPin, nil
}

// SetupGPIOInputPort exports the desired port and sets the direction to
// input. The accepted port numbers differ between Raspberry Pi models.
func SetupGPIOInputPort(port uint64) (*InputPin, error) {
	setupError := setupGPIOPort(port, "in")
	if setupError != nil {
		return nil, setupError
	}

	return &InputPin{
		valueFile: fmt.Sprintf("/sys/class/gpio/gpio%d/value", port),
	}, nil
}

func setupGPIOPort(port uint64, direction string) error {
	//The raspberry creates a separate folder for each exported GPIO port.
	//If a port has been exported once, it can't be unexported, therefore
	//the existence of this folder implies that the port is exported.
	_, targetGPIOFolderErr := os.Stat(fmt.Sprintf("/sys/class/gpio/gpio%d", port))
	if os.IsNotExist(targetGPIOFolderErr) {
		log.Printf("Exporting GPIO port %d\n", port)
		gpioPortExportErr := os.WriteFile("/sys/class/gpio/export", []byte(strconv.FormatUint(port, 10)), os.ModeAppend)
		if gpioPortExportErr != nil {
			return fmt.Errorf("error exporting port %d: %s", port, gpioPortExportErr)
		}
	} else if targetGPIOFolderErr != nil {
		return fmt.Errorf("unexpected error checking export status of GPIO port %d: %s", port, targetGPIOFolderErr)
	}

	var gpioPortDirectionError error
	directionPath := fmt.Sprintf("/sys/class/gpio/gpio%d/direction", port)
	//Seemingly the raspberry seems to initially create the folder and files
	//required for accessing the port, but initally sets incorrect permissions.
	//Here we are waiting until we are able to change the direction
	//successfully. If we can't do this within 200ms, we stop.
	tryTicker := time.NewTicker(20 * time.Millisecond)
	defer tryTicker.Stop()
	for i := 0; i < 10; i++ {
		<-tryTicker.C
		gpioPortDirectionError = os.WriteFile(directionPath, []byte(direction), os.ModeAppend)
		//Try until success
		if gpioPortDirectionError == nil {
			break
		}
	}

	if gpioPortDirectionError != nil {
		return fmt.Errorf("error setting direction to %s for port %d: %s", direction, port, gpioPortDirectionError)
	}

	return nil
}
