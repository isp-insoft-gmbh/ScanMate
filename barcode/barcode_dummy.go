// +build dummy test

package barcode

import (
	"bufio"
	"errors"
	"log"
	"math/rand"
	"os"
	"strings"
)

func CaptureImage() error {
	log.Println("Return Error on 'CaptureImage' (y/n)")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		scanErr := scanner.Err()
		panic(scanErr)
	}

	if strings.EqualFold(scanner.Text(), "y") {
		return errors.New("error capturing image")
	}

	return nil
}

func ParseBarcodeImage() (string, error) {
	log.Println("Error type for 'ParseBarcodeImage'? (b(arcodenotfound)/f(atal)/n(one))")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		scanErr := scanner.Err()
		panic(scanErr)
	}

	text := scanner.Text()
	if strings.EqualFold(text, "nobarcode") {
		return "", ErrNoBarcodeFound
	}

	if strings.EqualFold(text, "fatal") {
		return "", errors.New("fatal error")
	}

	log.Println("Enter desired barcode or leave empty for random barcode.")
	if !scanner.Scan() {
		scanErr := scanner.Err()
		panic(scanErr)
	}

	barcode := strings.TrimSpace(scanner.Text())
	if barcode != "" {
		return barcode, nil
	}

	var barcodeBuilder strings.Builder
	barcodeChars := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	for i := 0; i < 13; i++ {
		barcodeBuilder.WriteRune(barcodeChars[rand.Intn(len(barcodeChars))])
	}
	return barcodeBuilder.String(), nil
}
