package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

var pallete = []color.Color{color.Black, color.RGBA{0xff, 0x00, 0x00, 0xff}, color.RGBA{0x00, 0xff, 0x00, 0xff}, color.RGBA{0x00, 0x00, 0xff, 0xff}}

const (
	whiteIndex = 0
	blackIndex = 1
)

func handler(w http.ResponseWriter, r *http.Request) {
  mu.Lock()
  counter++
  mu.Unlock()
  fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL.Path, r.Proto)
  for k, v := range r.Header {
    fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
  }
  fmt.Fprintf(w, "Host = %q\n", r.Host)
  fmt.Fprintf(w, "Remote Address = %q\n", r.RemoteAddr)

  if err := r.ParseForm(); err != nil {
    log.Print(err)
  }
  for k, v := range r.Form {
    fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
  }
}

var counter int = 0
var mu sync.Mutex
func countHandler(w http.ResponseWriter, _ *http.Request) {
  mu.Lock()
  fmt.Fprintf(w, "Request count is %d\n", counter)
  mu.Unlock()
}

func lissajousHandler(w http.ResponseWriter, r *http.Request) {
  mu.Lock()
  counter++
  mu.Unlock()
  query := r.URL.Query()
  cyclesString := query.Get("cycle")
  if cyclesString == "" {
    cyclesString = "5"
  }
  cycles, err := strconv.Atoi(cyclesString)
  if err != nil {
    fmt.Fprintf(w, "Cannot get cycles: %v\n", err)
    return
  }
  lissajous(w, cycles)
}

func main() {
  http.HandleFunc("/", handler)
  http.HandleFunc("/count", countHandler)
  http.HandleFunc("/lissa", lissajousHandler)

  log.Fatal(http.ListenAndServe(":8080", nil))
}

func lissajous(out io.Writer, cycles int) {
	const (
		res     = 0.001
		size    = 500
		nframes = 64
		delay   = 8
	)
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0e0

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, pallete)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(freq*t + phase)
			img.SetColorIndex(size+int(x*size), size+int(y*size), uint8(i % len(pallete) + 1))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
