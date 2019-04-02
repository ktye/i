// ⍳ interpret
package main

import (
	"bufio"
	//"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ktye/i"
)

// args:
// 0: read from stdin, execute each line, continue on error
// filename: execute file, exit on error
// else: exec argv
var file string
var line int

type v = interface{}
type l = []v

func main() {
	a := make(map[v]v)
	i.E(nil, a)

	if len(os.Args) > 1 {
		if b, err := ioutil.ReadFile(os.Args[1]); err == nil {
			i.E(i.P(string(b)), a)
		} else {
			p(i.E(i.P(strings.Join(os.Args[1:], " ")), a))
		}
		return
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line++
		p(run(s.Text(), a))
	}
}

func run(t string, a map[v]v) (r interface{}) {
	defer func() {
		if c := recover(); c != nil {
			debug.PrintStack()
			r = c
			if file != "" {
				//fmt.Printf("%s:%d: %v\n", file, line, r)
				os.Exit(1)
			}
		}
	}()
	return i.E(i.P(t), a)
}

func p(v v) {
	println(fmt(v).(string))
}

var fmt func(v v) v

func init() {
	a := make(map[v]v)
	i.E(l{}, a)
	fmt = a["$:"].(func(v v) v)
}
