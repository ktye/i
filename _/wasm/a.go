package main

import "syscall/js"

//go:export run
func run() {
	document := js.Global().Get("document")
	s := document.Call("getElementById", "a").Get("value").String()
	println("run:", s)
	r := evl(prs(mkb([]c(s))))
	if m.k[r]>>28 == N {
		dec(r)
		println("null result")
		document.Call("getElementById", "b").Set("value", "")
	}
	r = kst(r)
	n, p := m.k[r]&atom, 8+r<<2
	s = string(m.c[p:p+n])
	document.Call("getElementById", "b").Set("value", s)
	println("result",s)
	dec(r)
}

func main() { ini() }
func red(x k) k
