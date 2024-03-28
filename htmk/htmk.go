package main

import (
	"bytes"
	_ "embed"
	"html"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

// get put post patch delete
// query[dict] png[LI]

// get["/table/t/{i:row}/{i:col}"]{x..}
// put[".."]
// post[".."]{[body]..}

var mu sync.Mutex

func main() {
	kinit()
	r2("get", 0)
	r2("put", 1)
	r2("post", 2)
	r2("delete", 3)
	r2("patch", 4)
	r1("query", 5)
	r1("png", 6)
	for _, f := range os.Args[1:] {
		b, e := os.ReadFile(f)
		fatal(e)
		dx(Val(KC(b)))
	}
	http.HandleFunc("/dev.html", dev)
	P := os.Getenv("P")
	if P == "" {
		P = ":3001"
	}
	http.ListenAndServe(P, nil)
}
func r1(s string, i uint64) {
	n1 := func(x uint64) uint64 { SetI32(int32(x)-12, 1); return x }
	Asn(Ks(s), n1(ti(14, int32(l2(i, KC(s))))))
}
func r2(s string, i uint64) { Asn(Ks(s), ti(14, int32(l2(i, KC(s))))) }
func Ks(s string) uint64    { return sc(KC(s)) }
func KC(b []byte) uint64    { r := mk(Ct, int32(len(b))); copy(Bytes[int32(r):], b); return r }
func CK(x uint64) string    { r := string(Bytes[int32(x) : int32(x)+nn(x)]); dx(x); return r }
func SK(x uint64) []string {
	r := make([]string, nn(x))
	for i := range r {
		r[i] = CK(rx(K(I32(0)) + I32(int32(x)+4*int32(i))))
	}
	dx(x)
	return r
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

var Q url.Values

func Native(x, y int64) int64 {
	if x == 5 {
		return query(Fst(x))
	} else if x == 6 {
		return pngx(Fst(x))
	}
	S := []string{"GET", "PUT", "POST", "DELETE", "PATCH"}[x]
	s, m := p(CK(x0(y)))
	f := r1(y)
	http.HandleFunc(S+" "+s, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		Q = r.URL.Query()
		for k := range m {
			dx(Asn(sK(k), m[k](r.PathValue(k))))
		}
		b, _ := io.ReadAll(r.Body)
		w.Write(b(lambda(f, KC(b))))
		mu.Unlock()
	})
}
func Kts(t T, s string) (r K) {
	if t == 3 {
		j, _ := strconv.Itoa(c)
		r = Ki(int32(j))
	} else if t == 4 {
		r = sc(KC(c))
	} else if t == 5 {
		f, _ := strconv.ParseFloat(c, 64)
		r = Kf(f)
	} else {
		r = KC(c)
	}
	return r
}
func query(x uint64) uint64 {
	k := x0(x)
	v := r1(x)
	n := nn(k)
	s := SK(rx(k))
	q := mk(Lt, nn(v))
	for i := int32(0); i < n; i++ {
		x := Ati(rx(v), i)
		SetI64(int32(q)+8*i, uf(Kts(tp(x), Q.Get(s[i]))))
		dx(x)
	}
	dx(v)
	return Key(k, q)
}
func pngx(x uint64) uint64 {
	n := 4 * nn(x)
	m := nn(Fst(rx(x)))
	im := image.NewRGBA(image.Rect(0, m, 0, n/4))
	p := 0
	for i := int32(0); i < m; i++ {
		y := I32(int32(x) + 8*i)
		copy(im.Pix[p:p+n], Bytes[y:y+n])
	}
	for i := 3; i < n; i += 4 {
		im.Pix[i] = 255
	}
	var buf bytes.Buffer
	png.Encode(&buf, im)
	b := buf.Bytes()
	r := mk(Ct, len(b))
	copy(Bytes[int32(r):], b)
	dx(x)
	return r
}
func b(x uint64) []byte {
	if tp(x) != Ct {
		x = Kst(x)
	}
	return CK(x)
}

func p(s string) (string, map[string]func(string) uint64) {
	//todo: parse "/path/{i:row}/{f:num}"
}

func dev(w http.ResponseWriter, r *http.Request) {
	var f []string
	d, _ := os.ReadDir(".")
	for _, x := range d {
		if s := x.Name(); strings.HasSuffix(s, ".html") || strings.HasSuffix(s, ".css") || strings.HasSuffix(s, ".js") || strings.HasSuffix(s, ".k") {
			f = append(f, "<span class='link'>"+html.EscapeString(s)+"</span>")
		}
	}
	w.Write([]byte(strings.Replace(devhtml, "FILES", strings.Join(f, " "), 1)))
}

const devhtml = `
<!DOCTYPE html>
<head><meta charset="utf-8"><title></title>
<style>
body{margin:0;overflow:hidden}
.link{}
.link:hover{pointer:cursor}
</style>
</head>
<body style="display:grid;grid-template-columns:auto 1fr">
<div style="height:100vh;display:flex;flex-flow:column;font-family:monospace">
<div>FILES<button>write</button></div>
<textarea style="flex-grow:1;resize:horizontal">
file content
</textarea>
</div>
<iframe style="width:100%;height:100%" src="/"></iframe></body></html>
`
