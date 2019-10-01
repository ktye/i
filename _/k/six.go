// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	w, h := 10, 10
	if len(os.Args) == 3 {
		w, _ = strconv.Atoi(os.Args[1])
		h, _ = strconv.Atoi(os.Args[2])
	}
	img := make([]int32, w*h)
	d, y := 0.0047/float64(w), 0.6535
	for j := 0; j < h; j++ {
		x := 0.002
		for i := 0; i < w; i++ {
			x += d
			img[j+w*i] = cmap(mandel(x, y))
		}
		y += d
	}
	var b bytes.Buffer
	p, v := paletted(img)
	sixel(p, v, w, &b)
	io.Copy(os.Stdout, &b)
}
func mandel(x, y float64) byte {
	a, b := x, y
	n := 0
	for ; n <= 255; n++ {
		if 4.0 <= a*a+b*b {
			break
		}
		a, b = a*a-b*b+x, 2.0*a*b+y
	}
	return byte(n)
}
func cmap(n byte) int32 {
	r, g, b := n, n, 255-255*(n-127)*(n-128)/128/128
	return int32(r)<<16 + int32(g)<<8 + int32(b)
}
func paletted(img []int32) (p []byte, palette []int32) {
	m := make(map[int32]byte)
	p = make([]byte, len(img))
	palette = make([]int32, 0, 256)
	for i, c := range img {
		if n, o := m[c]; !o {
			n = byte(len(palette))
			palette = append(palette, c)
			if len(palette) > 256 {
				panic("too many colors")
			}
			m[c] = n
			p[i] = n
		} else {
			p[i] = n
		}
	}
	return p, palette
}
func sixel(img []byte, pal []int32, width int, w io.Writer) error {
	height := len(img) / width
	w.Write([]byte{0x1b, 0x50, 0x30, 0x3b, 0x30, 0x3b, 0x38, 0x71, 0x22, 0x31, 0x3b, 0x31})
	for n, v := range pal {
		r, g, b := v>>16, (v&0xFF00)>>8, v&0xFF
		r = r * 100 / 0xFF
		g = g * 100 / 0xFF
		b = b * 100 / 0xFF
		fmt.Fprintf(w, "#%d;2;%d;%d;%d", n, r, g, b)
	}
	for z := 0; z < (height+5)/6; z++ {
		buf := make([]byte, width)
		if z > 0 {
			w.Write([]byte{'-'}) // new line
		}
		for n := byte(0); n < byte(len(pal)); n++ {
			any := false
			for x := 0; x < width; x++ {
				buf[x] = 0
				for p := 0; p < 6; p++ {
					y := z*6 + p
					for x := 0; x < width; x++ {
						if img[width*y+x] == n {
							buf[x] |= 1 << byte(p)
						}

					}
				}
				if buf[x] != 0 {
					any = true
				}
			}
			if any {
				fmt.Fprintf(w, "#%d", n)
				for _, c := range buf {
					w.Write([]byte{'?' + c})
				}
				w.Write([]byte{'$'})
			}
		}
	}
	w.Write([]byte{0x1b, 0x5c})
	return nil
}
func fatal(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
