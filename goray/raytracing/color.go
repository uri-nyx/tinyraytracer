package raytracing

import (
	"math"
)

type Color = V3

func NewColor(r, g, b float64) Color {
	return Color{r, g, b}
}

func bound(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func (c Color) DenormalizeSampledColor(samplesPerPixel int, gamma float64) []byte {
	r, g, b := c.X(), c.Y(), c.Z()

	scale := 1.0 / float64(samplesPerPixel)
	r = math.Pow(r * scale, 1 / gamma)
	g = math.Pow(g * scale, 1 / gamma)
	b = math.Pow(b * scale, 1 / gamma)

	return []byte{
		byte(math.Floor(bound(r, 0, 0.999) * 256)),
		byte(math.Floor(bound(g, 0, 0.999) * 256)),
		byte(math.Floor(bound(b, 0, 0.999) * 256))}
}
