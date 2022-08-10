package scene

import "math"

type Image struct {
	aspect_ratio float64
	Width, Height int
}

func NewImage(aspect_ratio float64, width int) Image {
	return Image{aspect_ratio, width, int(math.Floor(float64(width) / aspect_ratio))}
}