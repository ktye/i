// ⍳ interpret
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ktye/i"
)

// args:
// 0: read from stdin, continue on error
// filename: execute file, exit on error
// else: execute argv
var file string
var line int

type v = interface{}
type kt = map[v]v

func main() {
	var a kt

	var r io.Reader
	if len(os.Args) < 2 {
		r = os.Stdin
	} else {
		if f, err := os.Open(os.Args[1]); err == nil {
			defer f.Close()
			r = f
			file = os.Args[1]
		} else {
			r = strings.NewReader(strings.Join(os.Args[1:], " "))
		}
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		line++
		fmt.Println(run(a, s.Text()))
	}
}

func run(a kt, t string) (r interface{}) {
	defer func() {
		if c := recover(); c != nil {
			debug.PrintStack()
			r = c
			if file != "" {
				fmt.Printf("%s:%d: %v\n", file, line, r)
				os.Exit(1)
			}
		}
	}()
	r = i.E(i.P(t), a)
	return
}
