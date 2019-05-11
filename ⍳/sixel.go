package main

import (
	"bytes"
	"encoding/base64"
	_fmt "fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
	"sort"
	"strings"
)

// Values that are url-encoded images "data:image/png;base64,..." are displayed as images
// instead of text, depending on the terminal type.
// By default, images are converted to sixel format, which is supported by xterm and mintty.
//
// The plot package encodes plots in this format with the default stringer.

const pngPrefix = "data:image/png;base64,"

func sxl(s string) string {
	if !strings.HasPrefix(s, pngPrefix) {
		return s
	}
	b, err := base64.StdEncoding.DecodeString(s[len(pngPrefix):])
	if err != nil {
		return err.Error()
	}
	r := bytes.NewReader(b)
	img, err := png.Decode(r)
	if err != nil {
		return err.Error()
	}
	var buf bytes.Buffer
	newEncoder(&buf).encode(img)
	return string(buf.Bytes())
}

func guessPalette(img image.Image) color.Palette {
	p := make(map[color.RGBA]uint)
	rect := img.Bounds()
	for i := rect.Min.X; i < rect.Max.X; i++ {
		for k := rect.Min.Y; k < rect.Max.Y; k++ {
			r, g, b, a := img.At(i, k).RGBA()
			c := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
			p[c]++
		}
	}
	colors := make([]color.RGBA, len(p))
	n := 0
	for c := range p {
		colors[n] = c
		n++
	}
	sort.Slice(colors, func(i, j int) bool { return p[colors[i]] < p[colors[j]] })
	pal := make(color.Palette, 254)
	for i := range pal {
		pal[i] = color.Black
		if i < len(colors) {
			pal[i] = colors[i]
		}
	}
	return pal
}

// The following is a copy from package sixel: github.com/mattn/go-sixel/master/sixel.go (License MIT)
// Changes: decoder remove, use guessed palette instead of the quantiziser "github.com/soniakeys/quant/median".

// encoder encode image to sixel format
type encoder struct {
	w      io.Writer
	Width  int
	Height int
}

// NewEncoder return new instance of Encoder
func newEncoder(w io.Writer) *encoder {
	return &encoder{w: w}
}

const (
	specialChNr = byte(0x6d)
	specialChCr = byte(0x64)
)

// Encode do encoding
func (e *encoder) encode(img image.Image) error {
	nc := 255 // (>= 2, 8bit, index 0 is reserved for transparent key color)
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	if e.Width != 0 {
		width = e.Width
	}
	if e.Height != 0 {
		height = e.Height
	}

	/* Upstream version:
	// make adaptive palette using median cut alogrithm
	q := median.Quantizer(nc - 1)
	paletted := q.Paletted(img)
	if e.Dither {
		// copy source image to new image with applying floyd-stenberg dithering
		draw.FloydSteinberg.Draw(paletted, img.Bounds(), img, image.ZP)
	} else {
		draw.Draw(paletted, img.Bounds(), img, image.ZP, draw.Over)
	}
	*/
	paletted := image.NewPaletted(img.Bounds(), guessPalette(img))
	draw.Draw(paletted, img.Bounds(), img, image.ZP, draw.Over)

	// use on-memory output buffer for improving the performance
	var w io.Writer
	if _, ok := e.w.(*os.File); ok {
		w = bytes.NewBuffer(make([]byte, 0, 1024*32))
	} else {
		w = e.w
	}
	// DECSIXEL Introducer(\033P0;0;8q) + DECGRA ("1;1): Set Raster Attributes
	w.Write([]byte{0x1b, 0x50, 0x30, 0x3b, 0x30, 0x3b, 0x38, 0x71, 0x22, 0x31, 0x3b, 0x31})

	for n, v := range paletted.Palette {
		r, g, b, _ := v.RGBA()
		r = r * 100 / 0xFFFF
		g = g * 100 / 0xFFFF
		b = b * 100 / 0xFFFF
		// DECGCI (#): Graphics Color Introducer
		_fmt.Fprintf(w, "#%d;2;%d;%d;%d", n+1, r, g, b)
	}

	buf := make([]byte, width*nc)
	cset := make([]bool, nc)
	ch0 := specialChNr
	for z := 0; z < (height+5)/6; z++ {
		// DECGNL (-): Graphics Next Line
		if z > 0 {
			w.Write([]byte{0x2d})
		}
		for p := 0; p < 6; p++ {
			y := z*6 + p
			for x := 0; x < width; x++ {
				_, _, _, alpha := img.At(x, y).RGBA()
				if alpha != 0 {
					idx := paletted.ColorIndexAt(x, y) + 1
					cset[idx] = false // mark as used
					buf[width*int(idx)+x] |= 1 << uint(p)
				}
			}
		}
		for n := 1; n < nc; n++ {
			if cset[n] {
				continue
			}
			cset[n] = true
			// DECGCR ($): Graphics Carriage Return
			if ch0 == specialChCr {
				w.Write([]byte{0x24})
			}
			// select color (#%d)
			if n >= 100 {
				digit1 := n / 100
				digit2 := (n - digit1*100) / 10
				digit3 := n % 10
				c1 := byte(0x30 + digit1)
				c2 := byte(0x30 + digit2)
				c3 := byte(0x30 + digit3)
				w.Write([]byte{0x23, c1, c2, c3})
			} else if n >= 10 {
				c1 := byte(0x30 + n/10)
				c2 := byte(0x30 + n%10)
				w.Write([]byte{0x23, c1, c2})
			} else {
				w.Write([]byte{0x23, byte(0x30 + n)})
			}
			cnt := 0
			for x := 0; x < width; x++ {
				// make sixel character from 6 pixels
				ch := buf[width*n+x]
				buf[width*n+x] = 0
				if ch0 < 0x40 && ch != ch0 {
					// output sixel character
					s := 63 + ch0
					for ; cnt > 255; cnt -= 255 {
						w.Write([]byte{0x21, 0x32, 0x35, 0x35, s})
					}
					if cnt == 1 {
						w.Write([]byte{s})
					} else if cnt == 2 {
						w.Write([]byte{s, s})
					} else if cnt == 3 {
						w.Write([]byte{s, s, s})
					} else if cnt >= 100 {
						digit1 := cnt / 100
						digit2 := (cnt - digit1*100) / 10
						digit3 := cnt % 10
						c1 := byte(0x30 + digit1)
						c2 := byte(0x30 + digit2)
						c3 := byte(0x30 + digit3)
						// DECGRI (!): - Graphics Repeat Introducer
						w.Write([]byte{0x21, c1, c2, c3, s})
					} else if cnt >= 10 {
						c1 := byte(0x30 + cnt/10)
						c2 := byte(0x30 + cnt%10)
						// DECGRI (!): - Graphics Repeat Introducer
						w.Write([]byte{0x21, c1, c2, s})
					} else if cnt > 0 {
						// DECGRI (!): - Graphics Repeat Introducer
						w.Write([]byte{0x21, byte(0x30 + cnt), s})
					}
					cnt = 0
				}
				ch0 = ch
				cnt++
			}
			if ch0 != 0 {
				// output sixel character
				s := 63 + ch0
				for ; cnt > 255; cnt -= 255 {
					w.Write([]byte{0x21, 0x32, 0x35, 0x35, s})
				}
				if cnt == 1 {
					w.Write([]byte{s})
				} else if cnt == 2 {
					w.Write([]byte{s, s})
				} else if cnt == 3 {
					w.Write([]byte{s, s, s})
				} else if cnt >= 100 {
					digit1 := cnt / 100
					digit2 := (cnt - digit1*100) / 10
					digit3 := cnt % 10
					c1 := byte(0x30 + digit1)
					c2 := byte(0x30 + digit2)
					c3 := byte(0x30 + digit3)
					// DECGRI (!): - Graphics Repeat Introducer
					w.Write([]byte{0x21, c1, c2, c3, s})
				} else if cnt >= 10 {
					c1 := byte(0x30 + cnt/10)
					c2 := byte(0x30 + cnt%10)
					// DECGRI (!): - Graphics Repeat Introducer
					w.Write([]byte{0x21, c1, c2, s})
				} else if cnt > 0 {
					// DECGRI (!): - Graphics Repeat Introducer
					w.Write([]byte{0x21, byte(0x30 + cnt), s})
				}
			}
			ch0 = specialChCr
		}
	}
	// string terminator(ST)
	w.Write([]byte{0x1b, 0x5c})

	// copy to given buffer
	if _, ok := e.w.(*os.File); ok {
		w.(*bytes.Buffer).WriteTo(e.w)
	}

	return nil
}
