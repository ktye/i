// +build ignore

package main

// create png files for font in bdf and plan9 formats (assuming monospace).
// run with
//  go run png.go fontname.bdf  # creates fontname.png

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/zachomedia/go-bdf"
	"golang.org/x/image/font"
	"golang.org/x/image/font/plan9font"
	"golang.org/x/image/math/fixed"
)

func main() {
	if len(os.Args) != 2 {
		panic("args")
	}
	var face font.Face
	var off [2]int // dot position in character rectangle
	var w, h int
	name := os.Args[1]
	if strings.HasSuffix(name, ".bdf") {
		name = strings.TrimSuffix(name, ".bdf")
		b, err := ioutil.ReadFile(name + ".bdf")
		p(err)
		f, err := bdf.Parse(b)
		face = f.NewFace()
		p(err)
		h = f.Size
		_, advance, _ := face.GlyphBounds('M')
		w = advance.Floor()
		off = f.Encoding['M'].LowerPoint
	} else if strings.HasSuffix(name, ".font") {
		data, err := ioutil.ReadFile(name)
		p(err)
		name = strings.TrimSuffix(name, ".font")
		face, err = plan9font.ParseFont(data, func(rel string) ([]byte, error) { return ioutil.ReadFile(rel) })
		p(err)
		_, advance, _ := face.GlyphBounds('M')
		w = advance.Floor()
		h = face.Metrics().Height.Floor()
		off[1] = -face.Metrics().Descent.Floor()
		fmt.Printf("metrics: %v\n", face.Metrics())
	} else {
		panic("format")
	}
	ww, hh := 41+42*w, 3+4*h
	println(w, h, ww, hh)

	filename := name + "_" + strconv.Itoa(w) + "x" + strconv.Itoa(h) + ".png"
	file, err := os.Create(filename)
	p(err)
	defer file.Close()
	m := image.NewRGBA(image.Rect(0, 0, ww, hh))
	red := color.RGBA{255, 0, 0, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{red}, image.ZP, draw.Src)

	s := make([]byte, 42)
	copy(s, "????????????????")
	copy(s[42-len(filename):], filename)
	chars := []string{string(s), ":+-*%&|<>=!~,^#_$?@.0123456789'/\\;`\"(){}[]", "abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"}
	y := 0
	for row := range chars {
		x := 0
		for col := range chars[row] {
			glyph(m, face, chars[row][col], off, x, y, w, h)
			x += 1 + w
		}
		y += 1 + h
	}

	p(png.Encode(file, m))
}

func glyph(m draw.Image, f font.Face, r byte, off [2]int, x, y, w, h int) {
	draw.Draw(m, image.Rect(x, y, x+w, y+h), &image.Uniform{color.White}, image.ZP, draw.Src)
	dot := fixed.P(x+off[0], y+h+off[1])
	dr, mask, mp, _, _ := f.Glyph(dot, rune(r))
	draw.DrawMask(m, dr, &image.Uniform{color.Black}, image.ZP, mask, mp, draw.Src)
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}
