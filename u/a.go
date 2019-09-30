package main

import (
	"io"
	"io/ioutil"
)

func kinit() {
	ini()
	table[21] = red
}
func exi(x k) {
	println("exit?")
}

var stdin io.Reader

func red(x k) (r k) { // 1:x
	t, n := typ(x)
	if t == S {
		x = str(x)
		t, n = typ(x)
	}
	if t != C {
		panic("type")
	}
	var b []c
	if n == 0 {
		if c, err := ioutil.ReadAll(stdin); err != nil {
			panic(err)
		} else {
			b = c
		}
		if b == nil {
			exi(mki(0))
		}
	} else {
		panic("no filesystem")
		/*
			xp := 8 + x<<2
			p, err := ioutil.ReadFile(string(m.c[xp : xp+n]))
			if err != nil {
				panic(err)
			}
			b = p
		*/
	}
	n = k(len(b))
	r = mk(C, n)
	rp := 8 + r<<2
	copy(m.c[rp:rp+n], b)
	return r
}
