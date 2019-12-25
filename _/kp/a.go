package main

import "syscall"

func main() {
	ini(make([]f, 1<<13))
	nm = pnm
	ns = pns
	msrt = pmsrt
	table['!'] = ptil
	table['1'] = red
	table['q'] = exi
	table[dy+'1'] = wrt
	for {
		evp(red(wrt(inc(nans), enl(mkc(' '))))) // r: 1: (` 1: ," ")
	}
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
