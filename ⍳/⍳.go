// +build !ui

// ⍳ interpret
package main

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/ktye/i"
)

// args:
// 0: read from stdin, execute each line, continue on error
// filename: execute file, exit on error
// else: exec argv

func main() {
	a := kinit()
	a["print"] = p

	if len(os.Args) > 1 {
		if b, err := ioutil.ReadFile(os.Args[1]); err == nil {
			i.E(i.P(string(b)), a)
		} else {
			p(i.E(i.P(jon(" ", os.Args[1:]).(string)), a))
		}
		return
	}

	r := bufio.NewScanner(os.Stdin)
	for r.Scan() {
		p(run(r.Text(), a))
	}
}
func p(x v) { // print
	if x == nil {
		return
	}
	s, o := x.(string)
	if !o {
		s = fmt(x).(string)
	}
	// s = sxl(s) convert "data:image/png;base64..." to sixel (see github.com/ktye/ui/examples/interpret/sixel.go)
	println(s)
}
