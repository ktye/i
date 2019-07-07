package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"
)

var rd func() []c

func main() {
	ini()
	table[21] = red
	table[40] = exi
	table[21+dyad] = wrt
	if len(os.Args) < 2 {
		rd = readline(bufio.NewScanner(os.Stdin)) // 0:` or 1:` read a single line in interactive mode
		rpl()
	} else {
		defer stk()
		rd = read
		args := os.Args[1:]
		zx := mk(L, k(len(args))) // .z.x: args
		for i, a := range args {
			m.k[2+k(i)+zx] = mkb([]c(a))
		}
		asn(mks(".z.x"), inc(zx))
		lod(inc(m.k[2+zx]))
		dec(zx)
	}
}
func rpl() {
	for {
		try()
	}
}
func try() {
	defer stk()
	r := red(wrt(mku(0), enl(mkc(' ')))) // r: 1: ("" 1: ," ")
	if m.k[r]&atom == 0 {
		exi(mki(0))
	}
	evp(r)
}
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
		b = rd()
	} else {
		xp := 8 + x<<2
		p, err := ioutil.ReadFile(string(m.c[xp : xp+n]))
		if err != nil {
			panic(err)
		}
		b = p
	}
	n = k(len(b))
	r = mk(C, n)
	rp := 8 + r<<2
	copy(m.c[rp:rp+n], b)
	return r
}
func read() []c { // read all from stdin (non-interactive)
	b, err := ioutil.ReadAll(os.Stdin)
	if err == nil {
		return b
	}
	return []c{}
}
func readline(sc *bufio.Scanner) func() []c { // read single line (interactive)
	return func() []c {
		if sc.Scan() == false {
			return []c{}
		}
		return sc.Bytes()
	}
}
func wrt(x, y k) k { // x 1:y
	t, n := typ(x)
	if t == S {
		x = str(x)
		t, n = typ(x)
	}
	if t != C {
		panic("type")
	}
	if n != 0 {
		panic("nyi") // write to a file
	}
	t, n = typ(y)
	if t != C || n == atom {
		panic("type")
	}
	yp := 8 + y<<2
	w := bufio.NewWriter(os.Stdout)
	w.Write(m.c[yp : yp+n])
	w.Flush()
	return decr(y, x)
}
func exi(x k) (r k) { // exit built-in
	t, n := typ(x)
	if t == I && n == atom {
		os.Exit(int(m.k[2+x]))
	}
	os.Exit(1)
	return mk(N, atom)
}
func stk() {
	if r := recover(); r != nil {
		a, b := stack(r)
		println(a)
		s := cat(lup(mks(".f")), mkc(':'))
		s = cat(s, str(lup(mks(".n"))))
		s = cat(s, mkb([]c{':', ' '}))
		s = cat(s, lup(mks(".l")))
		s = cat(s, mkc('\n'))
		s = cat(s, mkb([]c(b+"\n")))
		dec(wrt(mku(0), s)) // file:n: "line"\nerror
	}
}
func stack(c interface{}) (stk, err string) {
	h := false
	for _, s := range strings.Split(string(debug.Stack()), "\n") {
		if h && strings.HasPrefix(s, "\t") {
			if i := strings.Index(s, "/ktye/i/"); i > 0 {
				s = s[i+7:]
			}
			stk += "\n" + s
		}
		if strings.Index(s, "panic.go") > 0 { // skip first lines
			h = true
		}
	}
	err = "?"
	if s, o := c.(string); o {
		err = s
	} else if e, o := c.(error); o {
		err = e.Error()
	}
	return stk, err
}
func fatal(s string) { println(s); os.Exit(1) }
