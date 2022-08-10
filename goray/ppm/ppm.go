// Package ppm implements a simple PPM image encoder.
// it supports only the P6 format.
package ppm

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// PpmImage is the basic Ppm image structure.
type PpmImage struct {
	Name string
	W, H uint
	header []byte
	pixels *bytes.Buffer
}

// New creates a new PpmImage.
func New(name string, width, height uint) *PpmImage {
	return &PpmImage{
		Name: name,
		W: width,
		H: height,
		header: []byte(fmt.Sprintf("P6\n%d %d\n255\n", width, height)),
		pixels: bytes.NewBuffer([]byte{}),
	}
}

// Save saves the image to disk (at the specified path)
// It uses the Name field as filename
func (img *PpmImage) Save(path string) error {
	imgPath := filepath.Join(path, img.Name)

	// Open the file
	file, err := os.OpenFile(imgPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return err
    }
    defer file.Close()

	// Write the header
	_, err = file.Write(img.header)
	if err != nil {
		return err
	}

	// Write the pixels
	_, err = img.pixels.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}

// WritePixel takes 3 bytes as an rgb value and writes it into the
// pixel buffer sequentially. Will panic if passing other than 3 bytes 
func (img *PpmImage) WritePixel(pixel []byte) {
	if len(pixel) != 3 {
		panic("WritePixel takes only 3 bytes")
	}
	
	img.pixels.Write(pixel)
}

