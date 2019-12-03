// +build ignore

// generate k fonts from png file.
// go run gen.go 10x20.png k 10 20 > f2.k
//
// convert png image to psf text format (www.seasip.info/Unix/PSF/)
// psftools can then be used to convert to binary psf or bdf format.
//
// go run gen.go 16x32.png psf 16 32 | txt2psf > 16x32.psf
// go run gen.go 16x32.png psf 16 32 | txt2psf | psf2bdf > 16x32.psf
package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
)

var φ string

func main() {
	if len(os.Args) != 5 {
		panic("args: file.png format width height")
	}
	f, err := os.Open(os.Args[1])
	p(err)
	φ = os.Args[2]
	if φ != "k" && φ != "psf" {
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
}

func glyph(m image.Image, x, y, w, h int) string {
	b, k := make([]bool, w*h), 0
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if rr, gg, bb, _ := m.At(x+i, y+j).RGBA(); rr == 0 && gg == 0 && bb == 0 {
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
	}
	panic("φ")
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
func p(err error) {
	if err != nil {
		panic(err)
	}
}
