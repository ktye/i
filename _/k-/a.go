package main

import "syscall"

func main() {
	ini()
	var buf [1024]c // don't write longer lines than this
	p := buf[:]
	for {
		n, err := syscall.Read(syscall.Stdin, p)
		if err != nil {
			panic(err)
		}
		if n > 1 {
			do(p[:n-1])
		}
	}
}
func do(c []c) {
	x := mk(C, k(len(c)))
	copy(m.c[8+x<<2:], c)
	x = kst(evl(prs(x)))
	p, n := 8+x<<2, m.k[x]&atom
	println(s(m.c[p : p+n]))
}
func red(x k) (r k) { panic("fs") } // 1:x
