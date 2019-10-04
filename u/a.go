package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"
)

func kinit(args []s) {
	ini()
	table[21] = red
	table[21+dyad] = wrt
	evl(prs(mkb([]c(tk)))) // load built-in t.k
	for _, a := range args {
		evp(red(mkb([]c(a))))
	}
}
func call(s s, x, y k) (r k) { return cal2(evl(mks("u"+s)), x, y) }

func exi(x k) { os.Exit(0) } // \\

var stdin io.Reader

func red(x k) (r k) { // 1:x
	t, n := typ(x)
	if t == S {
		x = str(x)
		t, n = typ(x)
	}
	if t != C {
		println(t)
		panic("type")
	}
	var b []c
	if n == 0 {
		panic("read stdin") // TODO
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
		panic("fs") // write to a file
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
