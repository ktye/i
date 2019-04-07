// +build !js

package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/ktye/i"
)

// in converts the string b (read from stdin) to data depending on x.
// The function is present if ⍳ is used in filter mode as function i:
//	cat data | ⍳ 'i 0'
// Conversions
// 	i 0          → sv{}
// 	i ""         → l{sv{}} / split at whitespace
// 	i ";"        → l{sv{}} / split at ;
// 	i`n          → as a numeric table: l{zv{...}}
// 	i`n","       → split at ","
// 	i`d`n`n`s";" → dict from csv header line / TODO
// 	i`j          → json / TODO
func in(x v, b string) v {
	if b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	lines := strings.Split(b, "\n")
	var xv []string
	if xs, o := x.(string); o {
		xv = []string{xs}
	} else {
		if xx, o := x.([]string); !o || len(xx) == 0 {
			return lines
		} else {
			xv = xx
		}
	}
	sep := ""
	if s := xv[len(xv)-1]; s == "\t" || s == ";" || s == "," || s == "|" {
		sep = s
		xv = xv[:len(xv)-1]
	}
	var r []v
	for i := range lines {
		var t []string
		if sep == "" {
			t = strings.Fields(lines[i])
		} else {
			t = strings.Split(lines[i], sep)
		}
		row := make([]v, len(t))
		for k := range row {
			row[k] = t[k]
			if (len(xv) == 1 && xv[0] == "n") || (len(xv) > i && xv[i] == "n") {
				row[k] = num(t[k])
			}
		}
		r = append(r, uf(row))
	}
	return r
}

func uf(x []v) v {
	n := make([]complex128, len(x))
	for i := range x {
		if z, o := x[i].(complex128); o {
			n[i] = z
		} else {
			return x
		}
	}
	return n
}
func setup() map[v]v {
	a := make(map[v]v)
	i.E(l{}, a)
	a["i"] = func(x v) v {
		b, _ := ioutil.ReadAll(os.Stdin)
		return in(x, string(b))
	}
	// Custom output formatters for interactive use: o$...
	a["o"] = map[v]v{
		"p": 6, // precision
		"a": 0, // polar complex degree precision
		"t": 1, // tables if possible
		"d": 1, // multiline dicts
		"m": 1, // matrix
		"l": 0, // nested list
		"q": 1, // auto quote
	}
	a["t"] = regtime()
	regplot(a)
	return a
}
