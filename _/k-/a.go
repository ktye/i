package main

import "syscall"

func main() {
	ini(make([]f, 1<<13))
	table[21] = red
	table[40] = exi
	table[21+dyad] = wrt
	for {
		evp(red(wrt(mku(0), enl(mkc(' '))))) // r: 1: (` 1: ," ")
	}
}
func grw() { m.f = append(m.f, make([]f, len(m.f))...) }
func red(x k) (r k) { // 1:x
	var a [1024]c // don't write longer lines than this
	b := a[:]
	if n, err := syscall.Read(syscall.Stdin, b); err != nil {
		panic(err)
	} else if n > 1 {
		return dex(x, mkb(b[:n-1]))
	} else if n == 1 {
		return dex(x, mk(C, 0))
	} else {
		exi(0)
	}
	return 0
}
func wrt(x, y k) (r k) { // x 1:y
	t, n := typ(y)
	if t != C || n == atom {
		panic("type")
	}
	p := 8 + y<<2
	print(s(m.c[p : p+n])) // stderr only
	return dex(y, x)
}
func exi(x k) (r k) { panic("ciao") }
