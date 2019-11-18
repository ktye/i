package main

import "syscall"

var V string

func main() {
	ini(make([]f, 1<<13))
	table[21+dyad] = wrt
	var buf [1024]c
	p := buf[:]
	n, _ := syscall.Read(syscall.Stdin, p)
	p = p[:n-1]
	if n > 1 && p[0] == '\\' && p[1] == '\\' {
		println(V)
		return
	}
	x := mk(C, k(len(p)))
	copy(m.c[8+x<<2:], p)
	evp(x)
}
func grw(c k) {
	if 2*len(m.f) <= cap(m.f) {
		m.f = m.f[:2*len(m.f)]
	} else {
		x := make([]f, 2*len(m.f), c/4)
		copy(x, m.f)
		m.f = x
	}
}
func red(x k) (r k) { return 0 }
func wrt(x, y k) (r k) { // x 1:y
	t, n := typ(y)
	if t != C || n == atom {
		panic("type")
	}
	p := 8 + y<<2
	print(s(m.c[p : p+n]))
	return dex(y, x)
}
