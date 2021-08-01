// +build !dummy,!test

package barcode

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const imageTempPath = "/run/user/1000/barcode.jpg"

func CaptureImage() error {
	raspistillCommand := exec.Command("raspistill", "--timeout", "1500", "--encoding", "jpg", "--output", imageTempPath, "--nopreview", "--quality", "10", "--rotation", "180", "--width", "1200", "--height", "1200")
	log.Println("Executing:", raspistillCommand)
	captureError := raspistillCommand.Run()
	if captureError != nil {
		return fmt.Errorf("error capturing image: %s", captureError)
	}

	return nil
}

func ParseBarcodeImage() (string, error) {
	zBarCommand := exec.Command("zbarimg", "-Sdisable", "-Sean13.enable", "-Sposition=false", "-Sx-density=2", "-Sy-density=2", "--raw", "--quiet", imageTempPath)
	log.Println("Executing:", zBarCommand)
	var barcodeBuffer bytes.Buffer
	zBarCommand.Stdout = &barcodeBuffer
	zBarError := zBarCommand.Run()

	if zBarError != nil {
		//4. No barcode was detected in one or more of the images. No other errors occurred.
		if zBarCommand.ProcessState.ExitCode() == 4 {
			return "", ErrNoBarcodeFound
		}

		return "", fmt.Errorf("error parsing image: %w", zBarError)
	}

	//zbar output contains trailing whitespace, therefore we always trim.
	return strings.TrimSpace(barcodeBuffer.String()), nil
}
