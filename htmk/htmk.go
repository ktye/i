package main

import (
	"fmt"
	"bytes"
	_ "embed"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// get put post patch delete
// get["/table/t?row=i&col=i"]{x..}

var mu sync.Mutex

func main() {
	kinit()
	r := func(s string, i uint64, a int32) {
		x := Asn(Ky(s), ti(14, int32(l2(i, KC([]byte(s))))))
		SetI32(int32(x)-12, a)
	}
	r("get", 0, 2)
	r("put", 1, 2)
	r("post", 2, 2)
	r("delete", 3, 2)
	r("patch", 4, 2)
	r("png", 5, 1)
	for _, f := range os.Args[1:] {
		b, e := os.ReadFile(f)
		fatal(e)
		dx(Val(KC(b)))
	}
	http.Handle("/", http.FileServer(http.Dir(".")))
	P := os.Getenv("P")
	if P == "" {
		P = ":3001"
	}
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Println("request", r.Method, r.URL) })
	fmt.Println("http://localhost"+P)
	http.ListenAndServe(P, nil)
}
func Ky(s string) uint64 { return sc(KC([]byte(s))) }
func KC(b []byte) uint64 { r := mk(Ct, int32(len(b))); copy(Bytes[int32(r):], b); return r }
func CK(x uint64) string { r := string(Bytes[int32(x) : int32(x)+nn(x)]); dx(x); return r }
func sK(x uint64) string { return CK(cs(x)) }
func SK(x uint64) []string {
	r := make([]string, nn(x))
	for i := range r {
		r[i] = CK(rx(uint64(I32(0) + I32(int32(x)+4*int32(i)))))
	}
	dx(x)
	return r
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

type H struct {
	m, p string
	q    map[string]int32
	f    uint64
}

func Native(x, y int64) int64 {
	if x == 5 {
		return int64(pngx(Fst(uint64(x))))
	}
	h := hparse(CK(x0(uint64(y))))
	h.m = []string{"GET", "PUT", "POST", "DELETE", "PATCH"}[x]
	h.f = r1(uint64(y))
	fmt.Println("register", h.m+" "+h.p)
	//http.Handle(h.m+" "+h.p, h)
	http.Handle(h.p, h)
	return 0
}
func hparse(s string) (h H) {
	var q string
	h.p, q, _ = strings.Cut(s, "?")
	h.q = make(map[string]int32)
	v := strings.Split(q, "&")
	m := map[string]int32{"i": it, "s": st, "f": ft} //def:Ct
	for _, s := range v {
		k, t, _ := strings.Cut(s, "=")
		h.q[k] = m[t]
	}
	return h
}
func (h H) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serve html", h.m, h.p, h.f)
	mu.Lock()
	q := r.URL.Query()
	for k, t := range h.q {
		dx(Asn(Ky(k), Kts(t, q.Get(k))))
	}
	b, _ := io.ReadAll(r.Body)
	x := lambda(h.f, KC(b))
	if tp(x) != Ct {
		x = Kst(x)
	}
	w.Write(Bytes[int32(x) : int32(x)+nn(x)])
	dx(x)
	mu.Unlock()
}
func Kts(t int32, s string) (r uint64) {
	if t == it {
		j, _ := strconv.Atoi(s)
		r = Ki(int32(j))
	} else if t == st {
		r = sc(KC([]byte(s)))
	} else if t == ft {
		f, _ := strconv.ParseFloat(s, 64)
		r = Kf(f)
	} else {
		r = KC([]byte(s))
	}
	return r
}
func pngx(x uint64) uint64 {
	n := 4 * int(nn(x))
	m := nn(Fst(rx(x)))
	im := image.NewRGBA(image.Rect(0, int(m), 0, n/4))
	p := 0
	for i := int32(0); i < m; i++ {
		y := I32(int32(x) + 8*i)
		copy(im.Pix[p:p+n], Bytes[y:y+int32(n)])
	}
	for i := 3; i < n; i += 4 {
		im.Pix[i] = 255
	}
	var buf bytes.Buffer
	png.Encode(&buf, im)
	b := buf.Bytes()
	r := mk(Ct, int32(len(b)))
	copy(Bytes[int32(r):], b)
	dx(x)
	return r
}
