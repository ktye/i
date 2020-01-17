package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type c = byte

var W bool // wasm64
var line int


func main() {
	var text bool
	flag.BoolVar(&text, "text", false, "text output")
	flag.BoolVar(&W, "64", false, "wasm64")
	flag.Parse()
	fmt.Println("text", text)
	scanner := bufio.NewScanner(os.Stdin)
	line = 1
	for scanner.Scan() {
		O(run([]c(scanner.Text())), text)
		if text {
			fmt.Println()
		}
		line++
	}
	fatal(scanner.Err())
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
func hxb(x c) (c, c) { h := "0123456789abcdef"; return h[x>>4], h[x&0x0F] }
func O(b []c, text bool) {
	if !text {
		os.Stdout.Write(b)
	} else {
		for _, z := range b {
			x, y := hxb(z)
			os.Stdout.Write([]c{x, y})
		}
	}
}
func run(b []c) (r []c) {
	p = 0
	for {
		p += white(b[p:])
		if p >= len(b)-1 {
			return r
		}
		if n, a := pAsm(b[p:]); n > 0 {
			r = append(r, a)
			b = b[n:]
		} else {
			panic(fmt.Errorf("%d:%d: parse", line, p+1))
		}
	}
}
func pAsm(b []c) (n int, r []c) {
	if b[0] != '.' {
		return 0, nil
	}
	for i := 1; i<len(b); i+=2 {
		if crHx(b[i]) && len(b) > i+1 && crHx(b[i+1]) {
			r = append(r, (xtoc(b[i]) << 4) | xtoc(b[i+1]))
		} else {
			break
		}
	}
	return 1 + len(r), r
}

func cr09(c c) bool     { return c >= '0' && c <= '9' }
func craZ(c c) bool     { return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') }
func cr0Z(c c) bool     { return cr09(c) || craZ(c) }
func crHx(c c) bool     { return cr09(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') }
func xtoc(x c) c {
	switch {
	case x < ':':
		return x - '0'
	case x < 'G':
		return 10 + x - 'A'
	default:
		return 10 + x - 'a'
	}
}
func white(b []c) int {
	for i := range b {
		if b != ' ' && b != '\t' && b != ';' {
			return i
		}
	}
	return len(b)
}

/*
func cat(x, y []b) b { return append(x, y...) }
func lebu(b []b, v uint64) []byte { // encode leb128
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}
func lebs(b []b, v int64) []b { // encode signed leb128
	for {
		c := uint8(v & 0x7f)
		s := uint8(v & 0x40)
		v >>= 7
		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}
*/
