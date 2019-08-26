package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

var stdout op
var ee, ss, dd bool

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":2019", nil)
	println(err.Error())
}
func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.URL.Path != "/k" {
		http.Error(w, "∄"+r.URL.Path, 404)
		return
	}
	println("url", r.URL.Path, r.Method)
	if method := r.Method; method == "GET" { // send a new front-end
		kinit()
		w.Write([]c(h))
		return
	}

	e, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		return
	}
	table[dyad] = nil // unset trigger

	if n := r.Header.Get("n"); n == ".e" {
		s, a, b := r.Header.Get("s"), strconv.Atoi(r.Header.Get("a")), strconv.Atoi(r.Header.Get("b"))
		dec(asn(mku(0x2e65000000000000), mkb(e), mk(N, atom)))                                // `.e:"line1\nline2"
		dec(asn(mku(0x2e73000000000000), cat(cat(mkb([]c(s)), mki(a)), mki(b)), mk(N, atom))) // `.s:("sel";7;10)
	} else {
		dec(kxy(mks(".fput"), mkb([]c(n)), mkb(e)))
	}

	// xxx

	table[dyad] = trg
	if name := r.Header.Get("file"); name != "" { // upload
		files[name] = body
		w.Write([]c("w " + name + "\n"))
		return
	}
	draw = false
	stdout.Buffer = bytes.NewBuffer(nil)
	dec(kxy(mks(".rsz"), mki(strconv.Atoi(r.Header.Get("w"))), mki(strconv.Atoi(r.Header.Get("h")))))
	try(r.Header.Get("k"))
	if draw {
		sendImage(w)
		println("send image")
		return
	}
	println("send response")
	w.Write(stdout.Bytes())
}
func kinit() { // each GET /k (e.g. page reload)
	println("kinit")
	ini()
	table[21] = red                                           // 0:x
	table[40] = kinit                                         // \\
	table[dyad] = trg                                         // trigger `.e`.s`.d
	table[21+dyad] = wrt                                      // x 0:y
	mkk(".rsz", "{$[~(x*y)~+/#:'`.d;`.d:y x#0;0]}")           // resize(w,h)
	mkk(".fput", "{`.f[`$$[8<#x:8#x:x]]:y}")                  // store file
	dec(asn(mks(".f"), key(mk(S, 0), mk(L, 0)), mk(N, atom))) // memfs `.f:(0#`)!()
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
	if n == 0 {
		panic("stdin")
	}
	xp := ptr(x, C)
	name := s(m.c[xp : xp+n])
	if b, o := files[name]; o == false {
		panic("∄ " + name)
	} else {
		return decr(x, mkb(b))
	}
}
func wrt(x, y k) (r k) { // x 1:y
	t, n := typ(x)
	if t == S {
		x = str(x)
		t, n = typ(x)
	}
	if t != C {
		panic("type")
	}
	if n != 0 {
		panic("nyi: write to virtual file")
	}
	t, n = typ(y)
	if t != C || n == atom {
		panic("type")
	}
	yp := 8 + y<<2
	stdout.Write(m.c[yp : yp+n])
	return decr(y, x)
}
func trg(x, y, f k) (r k) { // trigger assignments
	switch sym(8 + x<<2) {
	case 0x2e65000000000000: // `.e
		ee = true
	case 0x2e73000000000000: // `.s
		ss = true
	case 0x2e64000000000000: // `.d
		dd = true
	}
	return decr2(y, f, x)
}
func try(b []c) {
	defer stk()
	evp(mkb(b))
}
func sendImage(w http.ResponseWriter) {
	wk := lupo(mku(0x7700000000000000)) // `w
	if t, n := typ(wk); wk == 0 || t != I || n != atom {
		w.Write([]c("w is overwritten"))
		return
	}
	defer dec(wk)
	dpy := lupo(mku(0x6470790000000000)) // `dpy
	if dpy == 0 {
		w.Write([]c("no dpy"))
		return
	}
	t, n := typ(dpy)
	if t != I || n == atom || n == 0 {
		w.Write([]c("wrong dpy type"))
		return
	}
	defer dec(dpy)

	p, ww, hh, cp, kp := ptr(dpy, I), m.k[2+wk], n/m.k[2+wk], 0, k(0)
	im, b := image.NewRGBA(image.Rectangle{Max: image.Point{int(ww), int(hh)}}), bytes.NewBuffer(nil)
	for i := k(0); i < n; i++ { // see: golang.org/src/image/image.go
		kp = m.k[p]
		im.Pix[cp+0] = c(kp & 0xFF0000 >> 16)    // r
		im.Pix[cp+1] = c(kp & 0x00FF00 >> 8)     // g
		im.Pix[cp+2] = c(kp & 0xFF)              // b
		im.Pix[cp+3] = c(^kp & 0xFF000000 >> 24) // a
		cp += 4
		p++
	}
	if e := png.Encode(b, im); e != nil {
		w.Write([]c("png:" + e.Error()))
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write([]c("data:image/png;base64,"))
	e := base64.NewEncoder(base64.StdEncoding, w)
	io.Copy(e, b)
	e.Close()
}
func stk() {
	if r := recover(); r != nil {
		draw = false
		a, b := stack(r)
		dec(asn(mks(".stk"), mkb([]c(a)), mk(N, atom))) // stack trace: \s
		dec(wrt(mku(0), ano(m.k[srcp], mkb([]c(b)))))
	}
}
func stack(c interface{}) (stk, err s) {
	println(s(debug.Stack()))
	h := false
	for _, s := range strings.Split(s(debug.Stack()), "\n") {
		if h && strings.HasPrefix(s, "\t") {
			if i := strings.Index(s, "/ktye/i/w/"); i > 0 {
				s = strings.TrimSpace(s[i+7:])
				if len(s) > 0 {
					stk += "\n" + s
				}
			}
		}
		if strings.Index(s, "panic.go") > 0 { // skip first lines
			h = true
		}
	}
	err = "?"
	if s, o := c.(s); o {
		err = s
	} else if e, o := c.(error); o {
		err = e.Error()
	}
	return stk, err
}
func (w op) Write(b []c) (int, error) { // debug wrapper
	println(s(b))
	return w.Buffer.Write(b)
}
