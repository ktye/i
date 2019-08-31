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
	"sync"
)

var door sync.Mutex
var stdout *bytes.Buffer
var dpng []byte
var ee, ss, dd k

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":2019", nil)
	println(err.Error())
}
func hdr(r *http.Request) string { // TODO rm
	var s string
	for _, h := range []string{"n", "a", "b", "w", "h", "k"} {
		s += " " + h + "=" + r.Header.Get(h)
	}
	return s
}
func handler(w http.ResponseWriter, r *http.Request) {
	door.Lock()
	defer door.Unlock()
	println("url", r.URL.Path, r.Method, r.URL.RawQuery, hdr(r))
	defer r.Body.Close()
	switch r.URL.Path {
	case "/.d":
		sendImage(w)
		return
	case "/.e":
		e := atx(lup(mku(0)), mks("e"))
		n, p := m.k[e]&atom, 8+e<<2
		w.Header().Set("Content-Type", "text/plain")
		w.Write(m.c[p : p+n])
		dec(e)
		return
	case "/k":
		if r.Method == "GET" {
			kinit()
			w.Write([]c(h)) // send a new front-end
			return
		}
	case "/ws": // workspace
		sendFile(w, m.c, "ws")
		return
	default: // file server
		name := r.URL.Path
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
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
		sendFile(w, m.c[p:p+n], name)
		return
	} // POST /k:
	stdout = bytes.NewBuffer(nil)
	t, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		return
	}
	if n := r.Header.Get("n"); n == ".e" { // store editor state
		dec(amd(mku(0), mks("e"), inc(null), mkb(t)))                      // @[`;`e;;"line1\nline2"]
		dec(amd(mku(0), mks("s"), inc(null), cat(hi(r, "a"), hi(r, "b")))) // @[`;`s;;7 10]
	} else if n != "" { // file upload
		println("write to n:", n)
		out(wrt(mkb([]c(n)), mkb(t)))
		return
	}
	ee = decr(ee, atx(lup(mku(0)), mks("e")))
	ss = decr(ss, atx(lup(mku(0)), mks("s")))
	dd = decr(dd, kxy(mks(".rsz"), hi(r, "h"), hi(r, "w")))
	if !try(r.Header.Get("k")) { // eval k expr
		w.Write(stdout.Bytes())
		return
	}
	e := atx(lup(mku(0)), mks("e"))
	if ee != null && !match(e, ee) {
		println(".e`update")
		w.Header().Set("n", ".e")
	}
	dec(e)
	s := atx(lup(mku(0)), mks("s"))
	if ss != null && !match(s, ss) {
		st, sst, sn, ssn := typs(s, ss)
		println(".s`update", st, sst == A, sn, ssn)
		if st, sn := typ(s); st == I {
			w.Header().Set("n", ".e")
			if sn == atom || sn == 1 {
				hs(w, "a", 2+s)
				hs(w, "b", 2+s)
			} else if sn == 2 {
				hs(w, "a", 2+s)
				hs(w, "b", 3+s)
			}
		}
	}
	dec(s)
	d := atx(lup(mku(0)), mks("d"))
	if dd != null && !match(d, dd) {
		println("update .d")
		w.Write([]c(setImage(d)))
	}
	dec(d)

	b := stdout.Bytes()
	println("write ", len(b), "bytes")
	w.Write(b)
}
func printk(x k, a s) {
	if m.k[x]>>28 == N {
		println(a, "NULL")
	}
	c := kst(inc(x))
	defer dec(c)
	t, n := typ(c)
	if t != C {
		println(a, "type:", t, n)
		return
	}
	p := 8 + c<<2
	println(a, s(m.c[p:p+n]))
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
	table[21] = red      // 0:x
	table[21+dyad] = wrt // x 0:y
	ee, ss, dd = inc(null), inc(null), inc(null)
	dec(evl(prs(mkb([]c(`.e:"";.s:0 0;.d:,,0`)))))
	mkk(".rsz", "{$[(x*y)~+/#:'.d;.d;.d::(y;x)#0]}")          // resize[.d;h w]
	dec(asn(mks(".f"), key(mk(L, 0), mk(L, 0)), mk(N, atom))) // (memfs) `.f:("file1","file2")!(0x1234;0x5678..)
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
	if xt == S {
		x = str(x)
		xt, xn = typ(x)
	}
	if xt != C || xn == atom {
		panic("type")
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
	kk, vv := m.k[2+fs], m.k[3+fs]
	n := m.k[kk] & atom
	for i := k(0); i < n; i++ {
		if match(m.k[2+kk+i], x) {
			return decr(fs, inc(m.k[2+vv+i])), i
		}
	}
	return decr(fs, 0), n
}
func try(s s) (o bool) {
	defer func() {
		if r := recover(); r != nil { // TODO: restore ws
			// TODO print error
			//a, b := stack(r)
			//dec(asn(mks(".stk"), mkb([]c(a)), mk(N, atom))) // print with: \s
			//dec(wrt(mku(0), ano(m.k[srcp], mkb([]c(b)))))
		}
	}()
	evp(mkb([]c(s)))
	return true
}
func setImage(d k) s {
	defer dec(d)
	t, n := typ(d)
	if d == 0 || t != L || n == atom {
		return "bad`.d\n"
	}
	w, h, cp := m.k[m.k[2+d]]&atom, n, 0
	if w == atom || h == atom || w == 0 || h == 0 {
		return "bad`.d\n"
	}
	kp := 2 + m.k[2+d]
	im, b, col, row := image.NewRGBA(image.Rectangle{Max: image.Point{int(w), int(h)}}), bytes.NewBuffer(nil), k(0), k(0)
	for i := k(0); i < w*h; i++ { // see golang.org/src/image/image.go for the format of *RGBA
		p := m.k[kp+col]
		im.Pix[cp+0] = c(p & 0xFF0000 >> 16)    // r
		im.Pix[cp+1] = c(p & 0x00FF00 >> 8)     // g
		im.Pix[cp+2] = c(p & 0xFF)              // b
		im.Pix[cp+3] = c(^p & 0xFF000000 >> 24) // a
		cp += 4
		col++
		if col == w {
			col, row = 0, row+1
			kp = 2 + m.k[2+d+row]
			if t, n := typ(m.k[kp-2]); t != I || n != w {
				return "bad`.d\n"
			}
		}
	}
	if e := png.Encode(b, im); e != nil {
		return "png:" + e.Error()
	}
	buf := bytes.NewBuffer(nil)
	buf.Write([]c("data:image/png;base64,"))
	e := base64.NewEncoder(base64.StdEncoding, buf) // TODO `b64@
	io.Copy(e, b)
	e.Close()
	dpng = buf.Bytes()
	return ""
}
func hi(r *http.Request, s s) k           { i, _ := strconv.Atoi(r.Header.Get(s)); return mki(k(i)) }
func hs(w http.ResponseWriter, s s, x k)  { w.Header().Set(s, strconv.Itoa(int(m.k[x]))) }
func hss(w http.ResponseWriter, s s, x s) { w.Header().Set(s, s) }
func sendFile(w http.ResponseWriter, b []c, n s) {
	hss(w, "Content-Type", "application/octet-stream")
	hss(w, "Content-Disposition", "attachment;filename=\""+n+"\"")
	w.Write(b)
}
func stack(c interface{}) (stk, err s) { // compact stack trace (go runtime)
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
