package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync"
)

var rd func() []c
var dr sync.Mutex

func main() {
	ini()
	table[21] = red
	table[40] = exi
	table[21+dyad] = wrt
	args, addr := os.Args[1:], ""
	if len(args) == 1 && args[0] == "-kwac" {
		inikwac()
		return
	} else if len(args) > 1 && args[0] == "-p" {
		addr = args[1]
		if _, o := atoi([]c(args[1])); o {
			addr = ":" + args[1]
		}
		args = args[2:]
	} else if len(args) > 0 && args[0] == "-u" {
		addr = ":2019"
		dec(evl(prs(mkb([]c(ui)))))
		args = args[1:]
	}
	if len(args) > 0 {
		defer stk(false)
		rd = read
		zx := mk(L, k(len(args))) // .z.x: args
		for i, a := range args {
			m.k[2+k(i)+zx] = mkb([]c(a))
		}
		asn(mks(".z.x"), inc(zx), mk(N, atom))
		lod(inc(m.k[2+zx]))
		dec(zx)
	}
	if addr != "" {
		go http.ListenAndServe(addr, http.HandlerFunc(srv))
	}
	rd = readline(bufio.NewScanner(os.Stdin)) // 0:` or 1:` read a single line in interactive mode
	for {
		try()
	}
}
func try() {
	defer stk(true)
	evp(red(wrt(mku(0), enl(mkc(' '))))) // r: 1: ("" 1: ," ")
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
func stk(hide bool) {
	if r := recover(); r != nil {
		a, b := stack(r)
		if hide { // interactive
			dec(asn(mks(".stk"), mkb([]byte(a)), mk(N, atom))) // stack trace: \s
		} else {
			println(a + "\n")
		}
		dec(wrt(mku(0), ano(m.k[srcp], mkb([]byte(b)))))
	}
}
func stack(c interface{}) (stk, err string) {
	h := false
	for _, s := range strings.Split(string(debug.Stack()), "\n") {
		if h && strings.HasPrefix(s, "\t") {
			if i := strings.Index(s, "/ktye/i/"); i > 0 {
				s = strings.TrimSpace(s[i+7:])
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
func srv(w http.ResponseWriter, r *http.Request) {
	dr.Lock()
	defer dr.Unlock()
	buf := bytes.NewBuffer(nil)
	defer func() {
		w.Write(buf.Bytes())
		r.Body.Close()
	}()
	defer func() {
		if rec := recover(); rec != nil {
			a, b := stack(rec)
			println(a)
			buf = bytes.NewBuffer(nil)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(500)
			dec(wrt(mku(0), mkb([]byte(b))))
		}
	}()

	z, f := lupo(mku(0)), k(0)
	if z == 0 {
		return
	}
	z = atx(z, mks("Z"))
	if r.Method == "GET" { // .Z.G
		println("GET", r.URL.Path)
		f = atx(z, mks("G"))
	} else if r.Method == "POST" { // .Z.P
		println("POST", r.URL.Path)
		f = atx(z, mks("P"))
	} else {
		dec(z)
		return
	}
	if m.k[f]>>28 != N+1 {
		dec(f)
		return
	}
	hk, hv := mk(S, k(len(r.Header))), mk(L, k(len(r.Header)))
	kp, j := 8+hk<<2, k(0)
	for key := range r.Header {
		kv := key
		if len(kv) > 8 {
			kv = kv[:8]
		}
		mys(kp, btou([]c(kv)))
		m.k[2+j+hv] = mkb([]c(r.Header.Get(key)))
		kp, j = kp+8, j+1
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	l := mk(L, 3)                   // TODO? Query(dict)?
	m.k[2+l] = mkb([]c(r.URL.Path)) // ?
	m.k[3+l] = key(hk, hv)          // ?
	m.k[4+l] = mkb(b)               // ?
	println("call .Z.GP")
	y := cal(f, enlist(l)) // ? assume ( hdr(`a); body(`c) )
	t, n := typ(y)
	if t == L && n == 2 && m.k[m.k[2+y]]>>28 == A { // (hdr;"body") or "body"
		hdr := fst(inc(y))
		keys, vals := str(inc(m.k[2+hdr])), m.k[3+hdr]
		for i := k(0); i < atm1(m.k[keys]&atom); i++ {
			v := str(atx(inc(vals), mki(i)))
			kp, vp := ptr(m.k[2+i+keys], C), ptr(v, C)
			kn, vn := m.k[2+i+keys]&atom, m.k[v]&atom
			w.Header().Set(string(m.c[kp:kp+kn]), string(m.c[vp:vp+vn]))
		}
		y = drop(1, y)
		y, n = typ(y)
	}
	if t != C || n == atom {
		panic("type")
	} else if n > 0 {
		bp := ptr(m.k[3+y], C)
		n = atm1(m.k[m.k[3+y]] & atom)
		buf.Write(m.c[bp : bp+n])
	}
	dec(y)
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

const ui = ` \"ui:localhost:2019"
.Z.G:{ \x;"welcome"}`
