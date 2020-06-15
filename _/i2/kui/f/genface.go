package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Generate a go source file that implements font.Face

type font struct {
	w, h int
	name string
}

func main() {
	fonts := []font{
		font{10, 20, "10x20"},
		font{16, 32, "16x32"},
	}
	fmt.Println(head)
	for _, f := range fonts {
		b, err := ioutil.ReadFile(f.name + ".png")
		if err != nil {
			panic(err)
		}
		m, err := png.Decode(bytes.NewReader(b))
		if err != nil {
			panic(err)
		}
		ww, hh := 41+42*f.w, 3+4*f.h
		if m.Bounds() != image.Rect(0, 0, ww, hh) {
			fmt.Printf("%dx%d != %v\n", ww, hh, m.Bounds())
			panic("wrong image size")
		}
		chars := []string{"0123456789ABCDEF", ":+-*%&|<>=!~,^#_$?@.0123456789'/\\;`\"(){}[]", "abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"}
		mask := make(map[rune][]bool)
		y := 0
		for row := range chars {
			x := 0
			cols := len(chars[row])
			for col := 0; col < cols; col++ {
				s, c := glyph(m, x, y, f.w, f.h), chars[row][col]
				mask[rune(c)] = s
				x += 1 + f.w
			}
			y += 1 + f.h
		}
		mask[32] = make([]bool, f.w*f.h)  // space
		mask[127] = make([]bool, f.w*f.h) // replacement
		s := fACe
		s = strings.Replace(s, "NN", f.name, -1)
		s = strings.Replace(s, "WW", strconv.Itoa(f.w), -1)
		s = strings.Replace(s, "HH", strconv.Itoa(f.h), -1)
		fmt.Println(s)

		s = mASk
		s = strings.Replace(s, "NN", f.name, -1)
		s = strings.Replace(s, "WW", strconv.Itoa(f.w), -1)
		s = strings.Replace(s, "HH", strconv.Itoa(f.h), -1)
		fmt.Print(s)

		encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
		for i := rune(0x20); i <= 127; i++ {
			b := mask[i]
			if b == nil {
				fmt.Println("missing rune", i)
				panic("x")
			}
			encoder.Write(bb(b))
		}
		fmt.Println(`")}`)
	}
}
func glyph(m image.Image, x, y, w, h int) []bool {
	b, k := make([]bool, w*h), 0
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if rr, gg, bb, _ := m.At(x+i, y+j).RGBA(); rr == 0 && gg == 0 && bb == 0 {
				b[k] = true
			}
			k++
		}
	}
	return b
}
func hxb(x byte) (byte, byte) { h := "0123456789abcdef"; return h[x>>4], h[x&0x0F] }
func bb(v []bool) (r []byte) {
	p := 0
	b := byte(0)
	for _, t := range v {
		if t {
			b |= 1 << p
		}
		p++
		if p == 8 {
			r = append(r, b)
			b = 0
			p = 0
		}
	}
	if p != 0 {
		r = append(r, b)
	}
	return r
}

const fACe = `
var FaceNN = &Face{
	Advance: WW,
	Width:   WW,
	Height:  HH,
	Ascent:  HH,
	Descent: 0,
	Mask:    maskNN,
	Ranges: []Range{
		{'\u0020', '\u007f', 0},
		{'\ufffd', '\ufffe', 95},
	},
}
`

const mASk = `
var maskNN = &image.Alpha{
	Stride: WW,
	Rect:   image.Rectangle{Max: image.Point{WW, 96 * HH}},
	Pix: dec(WW, HH, "`

// derived from: image/master/font/basicfont/basicfont.go
const head = `
package rasterfont

import (
	"image"
	"encoding/base64"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func dec(w, h int, s string) (r []byte) {
	n := w*h*96
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	p, k := uint(0), 0
	c := b[0]
	r = make([]byte, n)
	for i := range r {
		if 1&(c>>p) != 0 {
			r[i] = 255
		}
		p++
		if p == 8 {
			p = 0
			k++
			if k < len(b) {
				c = b[k]
			}
		}
	}
	return r
}

type Range struct {
	Low, High rune
	Offset    int
}
type Face struct {
	Advance int
	Width   int
	Height  int
	Ascent  int
	Descent int
	Left    int
	Mask    image.Image
	Ranges  []Range
}

func (f *Face) Close() error                   { return nil }
func (f *Face) Kern(r0, r1 rune) fixed.Int26_6 { return 0 }
func (f *Face) Metrics() font.Metrics {
	return font.Metrics{
		Height:     fixed.I(f.Height),
		Ascent:     fixed.I(f.Ascent),
		Descent:    fixed.I(f.Descent),
		XHeight:    fixed.I(f.Ascent),
		CapHeight:  fixed.I(f.Ascent),
		CaretSlope: image.Point{X: 0, Y: 1},
	}
}
func (f *Face) Glyph(dot fixed.Point26_6, r rune) (
	dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
loop:
	for _, rr := range [2]rune{r, '\ufffd'} {
		for _, rng := range f.Ranges {
			if rr < rng.Low || rng.High <= rr {
				continue
			}
			maskp.Y = (int(rr-rng.Low) + rng.Offset) * (f.Ascent + f.Descent)
			ok = true
			break loop
		}
	}
	if !ok {
		return image.Rectangle{}, nil, image.Point{}, 0, false
	}
	x := int(dot.X+32)>>6 + f.Left
	y := int(dot.Y+32) >> 6
	dr = image.Rectangle{
		Min: image.Point{
			X: x,
			Y: y - f.Ascent,
		},
		Max: image.Point{
			X: x + f.Width,
			Y: y + f.Descent,
		},
	}
	return dr, f.Mask, maskp, fixed.I(f.Advance), true
}
func (f *Face) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	return fixed.R(0, -f.Ascent, f.Width, +f.Descent), fixed.I(f.Advance), true
}
func (f *Face) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	return fixed.I(f.Advance), true
}
`
