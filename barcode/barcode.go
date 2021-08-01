package barcode

import (
	"errors"
)

// ErrNoBarcodeFound implies that the image was correctly parsed, but did not
// contain anything that was parsable as a barcode.
var ErrNoBarcodeFound = errors.New("the image doesn't contain a barcode or isn't sharp enough")
