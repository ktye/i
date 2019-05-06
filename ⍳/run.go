package main

import (
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ktye/i"
)

type s = string
type v = interface{}
type l = []v

var fmt func(v) v
var jon func(v, v) v
var num func(v) v

func run(t string, a map[v]v) (r interface{}) {
	defer func() {
		if c := recover(); c != nil {
			rs := ""
			for _, s := range strings.Split(string(debug.Stack()), "\n") {
				if strings.HasPrefix(s, "\t") {
					rs += "\n" + s[1:]
				}
			}
			if s, o := c.(string); o {
				rs += "\n" + s
			} else if e, o := c.(error); o {
				rs += "\n" + e.Error()
			}
			r = rs
		}
	}()
	pr := i.P(t)
	return i.E(pr, a)
}

func kinit() map[v]v {
	a := make(map[v]v)
	i.E(l{}, a)
	fmt = a["$:"].(func(x v) v)
	jon = a["jon"].(func(x, y v) v)
	num = a["num"].(func(x v) v)

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
	a["x"] = T(1.0)
	a["q"] = func(x v) v { os.Exit(1); return 9876543210 } // a break point
	return a
}

type T float64

func (t T) Inc() T    { return t + 1.0 }
func (t T) Add(b T) T { return t + b }
