package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var timer time.Time

func kinit(args []s) {
	ini(make([]f, 1<<13))
	table[21] = red
	table[21+dyad] = wrt
	table[27+dyad] = tim
	table[39] = trp
	if len(args) == 0 {
		evl(prs(mkb(tk))) // load bundled application, e.g. t.k
	}
	for _, a := range args {
		evp(red(mkb([]c(a))))
	}
}
func call(s s, x, y k) (r k) { return cal2(evl(mks("u"+s)), x, y) }

func exi(x k) { os.Exit(0) } // \\

var stdin io.Reader

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
	return dex(y, x)
}
func tim(x, y k) (r k) { // 0 7:.. (reset timer), "alpha"7:.. print duration
	t := m.k[x] >> 28
	if t == I && m.k[2+x] == 0 {
		timer = time.Now()
		return dex(x, y)
	}
	if timer.IsZero() {
		timer = time.Now()
	}
	x = cat(kst(x), mkb([]c(": "+time.Since(timer).String())))
	out(x)
	return y
}
func trp(x, y k) (r k) {
	defer func() {
		if rc := recover(); rc != nil {
			_, b := stack(rc)
			r = mkb([]c(b))
		}
	}()
	return cal(x, enlist(y))
}
func stk(hide bool) {
	if r := recover(); r != nil {
		a, b := stack(r)
		if hide { // interactive
			dec(asn(mks(".stk"), mkb([]byte(a)), mk(N, atom))) // stack trace: \s
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
func readAttachment(a0 s) []c {
	f, err := os.Open(a0)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Seek(-4, io.SeekEnd)
	P(err)
	kui := make([]byte, 4)
	_, err = f.Read(kui)
	P(err)
	if string(kui) != "k\\ui" {
		return []byte{}
	}
	_, err = f.Seek(-15-4, io.SeekEnd)
	P(err)
	var ln int64
	_, err = fmt.Fscanf(f, "%15d", &ln)
	P(err)
	_, err = f.Seek(-15-4-ln, io.SeekEnd)
	P(err)
	b := make([]byte, int(ln))
	_, err = f.Read(b)
	P(err)
	return b
}
func P(e error) {
	if e != nil {
		panic(e)
	}
}

var tk []byte
