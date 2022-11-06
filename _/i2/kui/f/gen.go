// +build ignore

// generate k fonts from png file.
// go run gen.go 10x20.png k 10 20 > f2.k
//
// convert png image to psf text format (www.seasip.info/Unix/PSF/)
// psftools can then be used to convert to binary psf or bdf format.
//
// go run gen.go 16x32.png psf 16 32 | txt2psf > 16x32.psf
// go run gen.go 16x32.png psf 16 32 | txt2psf | psf2bdf > 16x32.psf
//
// convert to decker fonts: github.com/JohnEarnest/Decker/blob/main/docs/format.md#data-blocks
// go run gen.go 16x32.png decker 16 32   (to stdout)
package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
)

type dekfnt map[rune]string

var φ string

func main() {
	if len(os.Args) != 5 {
		panic("args: file.png format width height")
	}
	f, err := os.Open(os.Args[1])
	p(err)
	φ = os.Args[2]
	if φ != "k" && φ != "psf" && φ != "decker" {
		panic("format")
	}
	w, err := strconv.Atoi(os.Args[3])
	p(err)
	h, err := strconv.Atoi(os.Args[4])
	p(err)

	defer f.Close()
	m, err := png.Decode(f)
	p(err)
	ww, hh := 41+42*w, 3+4*h
	if m.Bounds() != image.Rect(0, 0, ww, hh) {
		p(fmt.Errorf("wrong image size: expected %d x %d, got %s\n", ww, hh, m.Bounds()))
	}
	nn := 42 + 26 + 26
	printp("%%PSF2\nVersion: 0\nFlags: 1\nLength: %d\nWidth: %d\nHeight: %d\n\n", nn, w, h)

	dek := make(dekfnt)
	chars := []string{"0123456789ABCDEF", ":+-*%&|<>=!~,^#_$?@.0123456789'/\\;`\"(){}[]", "abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"}
	y, n := 0, 0
	printk("font:(")
	cls, nl := "", ""
	for row := range chars {
		x := 0
		cols := len(chars[row])
		for col := 0; col < cols; col++ {
			s, c := glyph(m, x, y, w, h), chars[row][col]
			if row == 0 && col <= 16 {
				printk("%s%s /%c%c", nl, s, c, c)
			} else {
				if row == len(chars)-1 && col == cols-1 {
					cls = ")"
				}
				printk("\n %s%s /%c", s, cls, chars[row][col])
				printp("%%\nBitmap: %s\nUnicode: [00%X];\n", s, c)
				dek[rune(chars[row][col])] = s
				n++
			}
			nl = "\n "
			x += 1 + w
		}
		y += 1 + h
	}
	printk("\n")
	if n != nn {
		println(nn, n)
		panic("numchars")
	}

	// (32-126; 95 glyphs total), an extra glyph representing an ellipsis
	// - 1 unsigned byte: the maximum width of each glyph in the font, in pixels.
	// - 1 unsigned byte: the height of glyphs, in pixels, including all vertical padding.
	// - 1 unsigned byte: the number of horizontal pixels to advance between characters.
	// - 96 glyph records
	{
		printd("%%%%FNT0")
		var d []byte
		d = append(d, byte(w), byte(h), 0)
		d = append(d, []byte(dekglyph(make([]bool, w*h), w, h))...) //space
		for i := 0; i < 94; i++ {
			if len(dek[33+rune(i)]) == 0 {
				panic("missing " + strconv.Itoa(i))
			}
			d = append(d, []byte(dek[32+rune(i)])...)
		}
		d = append(d, []byte(dek['.'])...) //todo "…"
		printd(base64.StdEncoding.EncodeToString(d))
		printd("\n")
	}
}

func glyph(m image.Image, x, y, w, h int) string {
	b, k := make([]bool, w*h), 0
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if rr, gg, bb, aa := m.At(x+i, y+j).RGBA(); rr == 0 && gg == 0 && bb == 0 && aa == 65535 {
				b[k] = true
			}
			k++
		}
	}
	if φ == "k" {
		c := make([]byte, len(b)/4)
		for i := 0; i < len(c); i += 2 {
			c[i], c[i+1] = hxb(char(b[4*i : 4*i+8]))
		}
		return "0x" + string(c)
	} else if φ == "psf" {
		h := make([]byte, w*h)
		for i, t := range b {
			h[i] = '-'
			if t {
				h[i] = '#'
			}
		}
		return string(h)
	} else if φ == "decker" {
		return dekglyph(b, w, h)
	}
	panic("φ")
}
func dekglyph(b []bool, w, h int) string {
	// - 1 unsigned byte giving the true width of the glyph: how many pixels to advance horizontally after drawing the glyph.
	//   - (width/8)*height bytes of packed image data, in which each byte represents 8 horizontally adjacent pixels.
	//   - Glyphs with a width that is not evenly divisible by 8 will be padded with 0 bits. Note the similarity to IMG0.

	w8 := w / 8
	if 8*w8 < w {
		w8++
	}
	u := append([]byte{}, byte(w))
	for i := 0; i < h; i++ {
		a := make([]bool, 8*w8)
		copy(a, b[i*w:w+i*w])
		for j := 0; j < 8*w8; j += 8 {
			u = append(u, char(a[j:j+8]))
		}
	}
	return string(u)
}
func hxb(x byte) (byte, byte) { h := "0123456789abcdef"; return h[x>>4], h[x&0x0F] }
func char(b []bool) (c byte) {
	for i := 0; i < 8; i++ {
		if b[i] {
			c |= 1 << uint(7-i)
		}
	}
	return c
}
func printk(f string, a ...interface{}) {
	if φ == "k" {
		fmt.Printf(f, a...)
	}
}
func printp(f string, a ...interface{}) {
	if φ == "psf" {
		fmt.Printf(f, a...)
	}
}
func printd(f string, a ...interface{}) {
	if φ == "decker" {
		fmt.Printf(f, a...)
	}
}
func p(err error) {
	if err != nil {
		panic(err)
	}
}
