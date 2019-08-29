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

var stdout *bytes.Buffer
var ee, ss, dd bool
var dpng []byte

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":2019", nil)
	println(err.Error())
}
func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.URL.Path {
	case "/d.png":
		sendImage(w)
		return
	case "/k":
	case "/ws": // workspace
		sendFile(w, m.c, "ws")
	default: // file server
		name := r.URL.Path
		if len(name) > 0 && name[0] == '/' {
			name := name[1:]
		}
		if len(name) == 0 {
			return
		}
		rr, _ := lupf(mkb([]c(name)))
		if rr == 0 {
			http.Error(w, "∄"+r.URL.Path, 404)
			return
		}
		n, p := m.k[rr]&atom, 8+rr<<2
		sendFile(w, m.c[rr:rr+n], name)
	}
	println("url", r.URL.Path, r.Method)
	if method := r.Method; method == "GET" { // send a new front-end
		kinit()
		w.Write([]c(h))
		return
	}
	table[dyad] = nil // unset trigger
	stdout = bytes.NewBuffer(nil)
	e, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		return
	}
	ee, ss, dd = false, false, false
	if n := r.Header.Get("n"); n == ".e" {
		i := mk(I, 2)
		m.k[2+i], m.k[3+i] = hi(r, "a"), hi(r, "b")
		dec(asn(mku(0x2e65000000000000), mkb(e), mk(N, atom))) // `.e:"line1\nline2"
		dec(asn(mku(0x2e73000000000000), i, mk(N, atom)))      // `.s:7 10
	} else {
		out(wrt(mkb([]c(n)), mkb(e)))
		return
	}
	dec(kxy(mks(".rsz"), mki(hi(r, "w")), mki(hi(r, "h"))))
	table[dyad] = trg
	try(r.Header.Get("k"))
	if ss {
		s := lup(mku(0x2e73000000000000))
		if st, sn := typ(s); st == I {
			if sn == atom || sn < 2 {
				hs(w, "a", 2+s)
				hs(w, "b", 2+s)
			} else if sn > 1 {
				hs(w, "a", 2+s)
				hs(w, "b", 3+s)
			}
		}
		dec(s)
	}
	if ee {
		w.Header().Set("n", ".e")
	}
	if dd {
		w.Write([]c(setImage()))
	}
	w.Write(stdout.Bytes())
}
func sendImage(w http.ResponseWriter) {
	if dpng == nil {
		http.Error(w, "no image", 404)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(dpng)
}
func kinit() { // each GET /k (e.g. page reload)
	println("kinit")
	ini()
	table[21] = red                                           // 0:x
	table[40] = kinit                                         // \\
	table[dyad] = trg                                         // trigger `.e`.s`.d
	table[21+dyad] = wrt                                      // x 0:y
	mkk(".rsz", "{$[~(x*y)~+/#:'`.d;`.d:y x#0;0]}")           // resize(w,h)
	dec(asn(mks(".f"), key(mk(L, 0), mk(L, 0)), mk(N, atom))) // memfs `.f:("file1","file2")!(0x1234;0x5678..)
}
func red(x k) (r k) { // 1:x
	t, n := typ(x)
	if t == S {
		x = str(x)
		t, n = typ(x)
	}
	if t != C || n == atom {
		panic("type")
	}
	if n == 0 {
		panic("stdin")
	}
	xp := ptr(x, C)
	name := s(m.c[xp : xp+n])
	r, _ = lupf(x)
	if r == 0 {
		panic("∄ " + name)
	}
	return r
}
func wrt(x, y k) (r k) { // x 1:y
	xt, yt, xn, yn := typs(x, y)
	if xt != C || xn == atom {
		panic("type")
	} else if xt == S {
		x = str(x)
		xt, xn = typ(x)
	}
	if yt != C || yn == atom {
		panic("type")
	}
	if xn != 0 {
		b, ii := lupf(inc(x))
		dec(amd(mks(".f"), mki(ii), y, mk(N, atom)))
		return decr(b, x)
	}
	yp := 8 + y<<2
	stdout.Write(m.c[yp : yp+yn])
	return decr(y, x)
}
func lupf(x k) (r, j k) {
	fs := lup(mks(".f"))
	n := m.k[fs] & atom
	for i := k(0); i < n; i++ {
		if match(m.k[2+fs+i], x) {
			return decr(fs, atx(fs, x)), i
		}
	}
	return decr(fs, 0), n
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
func try(s s) {
	defer stk()
	evp(mkb([]c(s)))
}
func setImage() s {
	d := lupo(mku(0x2e64000000000000)) // `.d
	t, n := typ(d)
	if d == 0 || t != I || n == atom {
		return "bad`.d\n"
	}
	defer dec(d)

	// TODO ...

	p, ww, hh, cp, kp := ptr(d, I), m.k[2+wk], n/m.k[2+wk], 0, k(0)
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
	w := bytes.NewBuffer()
	w.Write([]c("data:image/png;base64,"))
	e := base64.NewEncoder(base64.StdEncoding, w)
	io.Copy(e, b)
	e.Close()
	dpng = w.Bytes()
}
func hi(r *http.Request, s s) k           { i, _ := strconv.Atoi(r.Header.Get(s)); return k(i) }
func hs(w http.ResponseWriter, s s, x k)  { w.Header().Set(s, strconv.Itoa(int(m.k[x]))) }
func hss(w http.ResponseWriter, s s, x s) { w.Header().Set(s, s) }
func sendFile(w http.ResponseWriter, b []c, n s) {
	hss(w, "Content-Type", "application/octet-stream")
	hss(w, "Content-Disposition", "attachment;filename=\""+n+"\"")
	w.Write(b)
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
