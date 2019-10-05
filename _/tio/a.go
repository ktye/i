package main

import "syscall"

var V string

func main() {
	ini()
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
func red(x k) (r k) { return 0 }
func wrt(x, y k) (r k) { // x 1:y
	t, n := typ(y)
	if t != C || n == atom {
		panic("type")
	}
	p := 8 + y<<2
	print(s(m.c[p : p+n]))
	return decr(y, x)
}
