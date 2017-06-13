// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Run with "web" command-line argument for web server.
// See page 13.
//!+main

// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//!-main
// Packages not needed by version in book.

//!+main

var palette = []color.Color{color.Black, color.RGBA{R: 30, G: 197, B: 3, A: 1}}

const (
	bgIndex = 0
	fgIndex = 1
)

var (
	cycles  float64 // number of complete x oscillator revolutions
	res     float64 // angular resolution
	size    float64 // image canvas covers [-size..+size]
	nframes int     // number of animation frames
	delay   int     // delay between frames in 10ms units
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "web" {
		handler := func(w http.ResponseWriter, r *http.Request) {
			cycles = 5
			res = 0.001
			size = 100
			nframes = 64
			delay = 8

			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed parsing form: "+err.Error(), http.StatusInternalServerError)
				return
			}

			fCycles := r.FormValue("cycles")
			if fCycles != "" {
				intFCycles, err := strconv.Atoi(fCycles)
				if err != nil {
					http.Error(w, "Failed parsing form: "+err.Error(), http.StatusBadRequest)
				}

				cycles = float64(intFCycles)
			}

			fRes := r.FormValue("res")
			if fRes != "" {
				intFRes, err := strconv.Atoi(fRes)
				if err != nil {
					http.Error(w, "Failed parsing form: "+err.Error(), http.StatusBadRequest)
				}

				res = float64(intFRes)
			}

			size := r.FormValue("size")
			if size != "" {
				intSize, err := strconv.Atoi(size)
				if err != nil {
					http.Error(w, "Failed parsing form: "+err.Error(), http.StatusBadRequest)
				}

				cycles = float64(intSize)
			}

			fNframes := r.FormValue("nframes")
			if fNframes != "" {
				intFNframes, err := strconv.Atoi(fNframes)
				if err != nil {
					http.Error(w, "Failed parsing form: "+err.Error(), http.StatusBadRequest)
				}

				nframes = intFNframes
			}

			fDelay := r.FormValue("delay")
			if fDelay != "" {
				intfDelay, err := strconv.Atoi(fDelay)
				if err != nil {
					http.Error(w, "Failed parsing form: "+err.Error(), http.StatusBadRequest)
				}

				delay = intfDelay
			}

			lissajous(w)
		}
		http.HandleFunc("/", handler)

		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}

	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, int(2*size+1), int(2*size+1))
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)

			img.SetColorIndex(int(size+x*size+0.5), int(size+y*size+0.5),
				fgIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

//!-main
