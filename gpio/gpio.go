package gpio

// InputPin allows to read the value of the mapped GPIO p^in.
type InputPin interface {
	//Enable returns the GPIO pins value (1 or 0)
	GetValue() (bool, error)
}

// OutputPin allows to read the value of the mapped GPIO pin or write it.
// The respective functions are Enable() and Disable().
type OutputPin interface {
	InputPin

	//Enable sets the GPIO pins value to 1 (true)
	Enable() error
	//Enable sets the GPIO pins value to 0 (false)
	Disable() error
}
