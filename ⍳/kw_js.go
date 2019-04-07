// wasm version
// build with
//	GOOS=js GOARCH=wasm go build
package main

import (
	"runtime/debug"
	"syscall/js"

	"github.com/ktye/i"
)

type v = interface{}
type s = string
type j = js.Value

var a map[v]v
var fmt func(v) v

var ex = js.FuncOf(func(_ j, x []j) v {
	return e(x[0].String())
})

func e(x s) (g s) {
	defer func() {
		if c := recover(); c != nil {
			g = string(debug.Stack())
		}
	}()
	r := i.E(i.P(x), a)
	t, o := r.(s)
	if o {
		return t
	}
	return fmt(r).(s)
}
func main() {
	c := make(chan bool)
	js.Global().Set("e", ex)
	println("i am wasm")
	<-c
}
func init() {
	a = make(map[v]v)
	i.E([]v{}, a)
	fmt = a["$:"].(func(x v) v)
}
