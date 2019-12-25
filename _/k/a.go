package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"
	"sync"
)

var rd func() []c
var dr sync.Mutex
var ws []f

func main() {
	ini(make([]f, 1<<13))
	ws := make([]f, 1<<13)
	save(ws)
	table['1'] = red
	table['q'] = exi
	// table[39] = trp
	table[139-dy] = plo
	table[dy+'1'] = wrt
	table[dy+'9'] = drw
	table[139] = plt
	args := os.Args[1:]
	if len(args) == 1 && args[0] == "-kwac" {
		inikwac()
		return
	}
	if len(args) > 0 {
		defer rest(nil)
		rd = read
		zx := mk(L, k(len(args))) // .z.x: args
		for i, a := range args {
			m.k[2+k(i)+zx] = mkb([]c(a))
		}
		asn(mks(".z.x"), inc(zx), 0)
		lod(inc(m.k[2+zx]))
		dec(zx)
	}
	rd = readline(bufio.NewScanner(os.Stdin)) // 0:` or 1:` read a single line in interactive mode
	for {
		try()
	}
}
func try() {
	defer rest(ws)
	evp(red(wrt(inc(nans), enl(mkc(' '))))) // r: 1: ("" 1: ," ")
	if len(ws) < len(m.f) {
		ws = make([]f, len(m.f))
	}
	save(ws)
}

/* todo
func trp(x, y k) (r k) {
	defer func() {
		if rc := recover(); rc != nil {
			_, b := stack(rc)
			r = mkb([]c(b))
		}
	}()
	return cal(x, enlist(y))
}
*/
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
		if b == nil {
			exi(mki(0))
		}
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
	return dex(x, r)
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
			return nil
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
	return dex(y, x)
}
func exi(x k) (r k) { // exit built-in
	t, n := typ(x)
	if t == I && n == atom {
		os.Exit(int(m.k[2+x]))
	}
	os.Exit(1)
	return 0
}
func rest(mem []f) {
	if r := recover(); r != nil {
		a, b := stack(r)
		if mem != nil { // interactive
			swap(mem)
			dec(asn(mks(".stk"), mkb([]byte(a)), 0)) // stack trace: \s
		} else {
			println(a + "\n")
		}
		dec(wrt(inc(nans), ano(m.k[srcp], mkb([]byte(b)))))
	}
}
func stack(c interface{}) (stk, err string) {
	h := false
	for _, s := range strings.Split(string(debug.Stack()), "\n") {
		if h && strings.HasPrefix(s, "\t") {
			if i := strings.Index(s, "/i/"); i > 0 {
				s = strings.TrimSpace(s[i:])
			}
			if len(s) > 0 {
				stk += "\n" + s
			}
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
func inikwac() { // write initial memory as data section
	skip := 0
	fmt.Printf("(0;0x")
	for i, c := range m.c {
		if c == 0 {
			skip++
		} else {
			if skip < 8 {
				for i := 0; i < skip; i++ {
					fmt.Printf("00")
				}
			} else if skip != 0 {
				fmt.Printf(";%d;0x", i)
			}
			fmt.Printf("%02x", c)
			skip = 0
		}
	}
	fmt.Println(")")
}
func pr(x k, a ...interface{}) {
	fmt.Printf(":%x ", x)
	r := kst(inc(x))
	_, n := typ(r)
	s := s(m.c[8+r<<2 : 8+n+r<<2])
	dec(r)
	fmt.Println(a, s)
}
func fatal(s string) { println(s); os.Exit(1) }
