// +build ignore

package main

// go run gen.go 10x20.png 10 20

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		panic("args: file.png width height")
	}
	w, err := strconv.Atoi(os.Args[2])
	p(err)
	h, err := strconv.Atoi(os.Args[3])
	p(err)
	f, err := os.Open(os.Args[1])
	p(err)
	defer f.Close()
	m, err := png.Decode(f)
	p(err)
	ww, hh := 41+42*w, 3+4*h
	if m.Bounds() != image.Rect(0, 0, ww, hh) {
		p(fmt.Errorf("wrong image size: expected %d x %d\n", ww, hh))
	}

	chars := []string{"0123456789ABCDE", ":+-*%&|<>=!~,^#_$?@.0123456789'/\\;`\"(){}[]", "abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"}
	y := 0
	fmt.Printf("font:(")
	cls, nl := "", ""
	for row := range chars {
		x := 0
		cols := len(chars[row])
		for col := 0; col < cols; col++ {
			s, c := glyph(m, x, y, w, h), chars[row][col]
			if row == 0 && col <= 16 {
				fmt.Printf("%s%s /%c%c", nl, s, c, c)
			} else {
				if row == len(chars)-1 && col == cols-1 {
					cls = ")"
				}
				fmt.Printf("\n %s%s /%c", s, cls, chars[row][col])
			}
			nl = "\n "
			x += 1 + w
		}
		y += 1 + h
	}
	fmt.Printf("\n")
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
	c := make([]byte, len(b)/4)
	for i := 0; i < len(c); i += 2 {
		c[i], c[i+1] = hxb(char(b[4*i : 4*i+8]))
	}
	return "0x" + string(c)
}
func hxb(x byte) (byte, byte) { h := "0123456789abcdef"; return h[x>>4], h[x&0x0F] }
func char(b []bool) (c byte) {
	for i := 0; i < 8; i++ {
		if b[i] {
			c |= 1 << i
		}
	}
	return c
}
func p(err error) {
	if err != nil {
		panic(err)
	}
}
