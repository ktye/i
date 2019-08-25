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

type op struct{ *bytes.Buffer }

var stdout op
var lastImg []k
var files map[s][]c

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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		return
	}
	if name := r.Header.Get("file"); name != "" {
		files[name] = body
		w.Write([]c("w " + name + "\n"))
		return
	}
	stdout.Buffer = bytes.NewBuffer(nil)
	setSize(r.Header.Get("width"), r.Header.Get("height"))
	try(body)
	if sendImage(w) {
		println("send image")
		return
	}
	println("send response")
	w.Write(stdout.Bytes())
}
func kinit() {
	println("kinit")
	files = make(map[s][]c)
	lastImg = nil
	ini()
	table[21] = red      // 0:x
	table[40] = kinit    // \\
	table[21+dyad] = wrt // x 0:y
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
func try(b []c) {
	defer stk()
	evp(mkb(b))
}
func setSize(width, height s) {
	wi, _ := strconv.Atoi(width)
	hi, _ := strconv.Atoi(height)
	dec(asn(mku(0x7700000000000000), mki(k(wi)), mk(N, atom)))      // `w
	dec(asn(mku(0x6800000000000000), mki(k(hi)), mk(N, atom)))      // `h
	dec(asn(mku(0x6470790000000000), mk(I, k(wi*hi)), mk(N, atom))) // `dpy
}
func sendImage(w http.ResponseWriter) bool {
	wk := lupo(mku(0x7700000000000000)) // `w
	if t, n := typ(wk); wk == 0 || t != I || n != atom {
		println("no w")
		return false
	}
	defer dec(wk)
	dpy := lupo(mku(0x6470790000000000)) // `dpy
	if dpy == 0 {
		println("no dpy")
		return false
	}
	t, n := typ(dpy)
	if t != I || n == atom || n == 0 {
		println("wrong dpy type", t, n)
		return false
	}
	defer dec(dpy)

	p := ptr(dpy, I)
	if lastImg != nil && n == k(len(lastImg)) {
		same := true
		for i := k(0); i < n; i++ {
			if m.k[p+i] != lastImg[i] {
				println("image differs at", i, "of", n)
				same = false
				break
			}
		}
		if same {
			println("same as last time")
			return false
		}
	} else {
		println("new image", n, "last==nil?", lastImg == nil)
		lastImg = make([]k, n)

	}
	copy(lastImg, m.k[p:p+n])

	ww, hh := m.k[2+wk], n/m.k[2+wk]
	println("w/h", ww, hh)
	im, b, cp, kp := image.NewRGBA(image.Rectangle{Max: image.Point{int(ww), int(hh)}}), bytes.NewBuffer(nil), 0, k(0)
	for i := k(0); i < n; i++ { // see: golang.org/src/image/image.go
		kp = m.k[p]
		im.Pix[cp+0] = c(kp & 0xFF0000 >> 16)    // r
		im.Pix[cp+1] = c(kp & 0x00FF00 >> 8)     // g
		im.Pix[cp+2] = c(kp & 0xFF)              // b
		im.Pix[cp+3] = c(^kp & 0xFF000000 >> 24) // a
		cp += 4
		p++
	}
	if png.Encode(b, im) != nil {
		println("pngEncode")
		return false
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write([]c("data:image/png;base64,"))
	e := base64.NewEncoder(base64.StdEncoding, w)
	io.Copy(e, b)
	e.Close()
	return true
}
func stk() {
	if r := recover(); r != nil {
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
func (w op) Write(b []c) (int, error) {
	println(s(b))
	return w.Buffer.Write(b)
}
