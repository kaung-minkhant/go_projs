package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, 2, 2
		width, height          = 1024, 1024
		subPixels              = 4
		epsX                   = (xmax - xmin) / width
		epsY                   = (ymax - ymin) / height
	)

	offX := []float64{-epsX, epsX}
	offY := []float64{-epsY, epsY}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			subPixels := make([]color.Color, 0, 4)
			for i := 0; i < 2; i++ {
				for j := 0; j < 2; j++ {
					z := complex(x+offX[i], y+offY[j])
					subPixels = append(subPixels, mandelbrot(z))
				}
			}
			img.Set(px, py, avg(subPixels))
		}
	}
	png.Encode(os.Stdout, img)
}

func avg(colors []color.Color) color.Color {
	var r, g, b, a uint32 = 0, 0, 0, 0
	for _, c := range colors {
		_r, _g, _b, _a := c.RGBA()
		r += _r
		g += _g
		b += _b
		a += _a
	}
	n := uint32(len(colors))
	return color.RGBA{uint8(r / n), uint8(g / n), uint8(b / n), uint8(a / n)}

}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128

	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.RGBA{contrast * n % 255, contrast * n * contrast * n % 255, 255 - contrast*n, 255 - contrast*n}
			// return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
