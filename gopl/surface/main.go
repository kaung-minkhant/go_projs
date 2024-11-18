package main

import (
	"fmt"
	"math"
)

const (
	width, height = 1000, 750
	cells         = 110
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
  fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' " + "style='stroke: grey; fill:white; stroke-width: 0.7' " + "width='%d' height='%d'>\n", width, height)

  for i := 0; i < cells; i++ {
    for j := 0; j < cells; j++ {
      ax, ay := cornor(i+1, j)
      bx, by := cornor(i, j)
      cx, cy := cornor(i, j+1)
      dx, dy := cornor(i+1, j+1)
      
      if math.IsNaN(ax) || math.IsNaN(ay) || math.IsNaN(bx) || math.IsNaN(by) || math.IsNaN(cx) || math.IsNaN(cy) ||math.IsNaN(dx) || math.IsNaN(dy) {
        continue
      }

      fmt.Printf("<polygon style='stroke: %s; fill: #222222' points='%g,%g %g,%g %g,%g %g,%g'/>\n", "#666666", ax, ay, bx, by, cx, cy, dx, dy)
    }
  }
  fmt.Println("</svg>")
}

func cornor(i, j int) (float64, float64) {
  x := xyrange * (float64(i)/cells - 0.5)
  y := xyrange * (float64(j)/cells - 0.5)

  z := f(x, y)

  sx := width/2 + (x-y)*cos30*xyscale
  sy := height/2 + (x+y)*sin30*xyscale - z*zscale

  return sx, sy
}

func f(x, y float64) float64 {
  r := math.Hypot(x, y)
  result := math.Sin(r)/r
  return result
}
