//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"
)

// print dict for floating point literals that do not roundtrip
// flt.out:
//  k k.k flt.k -e

func rev(u uint64) (r uint64) {
	for i := uint64(0); i < 8; i++ {
		ui := (0xff & (u >> (8 * i)))
		j := 56 - 8*i
		r |= ui << j
	}
	return r
}
func xf(u uint64) float64 {
	return math.Float64frombits(rev(u))
}

func main() {
	m := make(map[string]float64)
	/*
		var f float64 = 3.141592653589793
		u := math.Float64bits(f)
		fmt.Printf("%x %x\n", rev(u), u)
		fmt.Println(xf(0x399d52a246df913f))
	*/
	b, e := os.ReadFile("flt.out")
	fatal(e)
	v := bytes.Split(b, []byte("\n"))
	for i := range v {
		var u uint64
		var f float64
		s := string(v[i])
		if s == "" {
			continue
		}
		s = strings.TrimPrefix(s, "0x")
		if strings.HasSuffix(s, ".") {
			s += "0"
		}
		if n, e := fmt.Sscanf(strings.TrimPrefix(s, "0x"), "%x %f", &u, &f); n != 2 || e != nil {
			fmt.Println(string(v[i]), "=>", n, e)
			panic("scanf")
		}
		g := xf(u)
		//fmt.Printf("%s g=%v u=%x rev=%x\n", s, g, u, rev(u))
		if f != g {
			s = "0x" + s[:16]
			m[s] = g
		} else {
			// fmt.Println(u, f, "roundtrip")
		}
	}
	var keys []string
	var vals []string
	for k, v := range m {
		keys = append(keys, k)
		vals = append(vals, fmt.Sprintf("\"%v\"", v))
	}
	fmt.Println("(" + strings.Join(keys, ";") + ")!(" + strings.Join(vals, ";") + ")")
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
