// +build ignore

package main

import (
	"fmt"
	"strconv"
	"image/png"
	"bytes"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic("args")
	}
	b, e := os.ReadFile(os.Args[1])
	fatal(e)
	m, e := png.Decode(bytes.NewReader(b))
	fatal(e)

	hist := make(map[string]int)

	rect := m.Bounds()
	for x := rect.Min.X; x<rect.Max.X; x++ {
		for y := rect.Min.Y; y<rect.Max.Y; y++ {
			r, g, b, a := m.At(x, y).RGBA()
			//fmt.Println(x, y, r, g, b, a)
			hist[si(r)+"."+si(g)+"."+si(b)+"."+si(a)]++
		}
	}
	for k, v := range hist {
		fmt.Println(k,v)
	}
}
func si(i uint32) string {
	return strconv.Itoa(int(i))
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
