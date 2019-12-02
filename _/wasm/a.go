//+build !tinygo

package main

import "syscall/js"

var ctx, u8a, img js.Value
var w, h k

func main() {
	ini(make([]f, 1<<13))
	table[21] = red
	table[21+dyad] = wrt
	table[39] = trp
	js.Global().Set("ui", js.FuncOf(ui))
	select {}
}
func grw(c k) {
	if 2*len(m.f) <= cap(m.f) {
		m.f = m.f[:2*len(m.f)]
	} else {
		x := make([]f, 2*len(m.f), c/4)
		copy(x, m.f)
		m.f = x
	}
}
func ui(_ js.Value, args []js.Value) interface{} {
	var r k
	a := func(n int) k { return k(args[n].Int()) }
	switch a(0) {
	case 0:
		r = cal(lup(mks("uk")), l2(mki(a(1)), modifiers(a(2), a(3), a(4))))
	case 1:
		// TODO mouse
		return 0
	case 2:
		println("size", a(1), a(2))
		if w != a(1) || h != a(2) {
			w, h = a(1), a(2)
			ctx = js.Global().Get("ctx")
			u8a = js.Global().Get("Uint8Array").New(int(4 * w * h))
			img = ctx.Call("createImageData", w, h)
		}
		r = cal2(lup(mks("us")), mki(w), mki(h))
	case 339:
		return dex(asn(mks("ut"), enlist(mkb([]c(" "+args[1].String()))), inc(null)), 0)
	case 416:
		return dex(evl(prs(red(mkb([]c(args[1].String()))))), 0)
	case 7088:
		r = kst(evl(prs(mkb([]c(args[1].String())))))
		t, n := typ(r)
		if t != C || n == atom {
			return dex(r, 0)
		}
		p := ptr(r, C)
		println(string(m.c[p : p+n]))
		return dex(r, 0)
	default:
		return nil
	}
	return dex(draw(r), 0)
}
func modifiers(a, b, c k) (r k) {
	r = mk(I, 3)
	m.k[2+r] = k(a)
	m.k[3+r] = k(b)
	m.k[4+r] = k(c)
	return r
}
func draw(x k) k {
	t, n := typ(x)
	if t != I || n == atom || n == 0 {
		return x
	}
	for i := k(0); i < n; i++ {
		m.k[2+i+x] |= 0xFF000000
	}
	p := ptr(x, C)
	js.CopyBytesToJS(u8a, m.c[p:p+4*n])
	img.Get("data").Call("set", u8a)
	ctx.Call("putImageData", img, 0, 0)
	return x
}
func red(x k) (r k) { // 1:x read from js f[file]
	t, n := typ(x)
	if t == S && n == atom {
		x = str(x)
		t, n = typ(x)
	}
	if t != C || n == atom {
		panic("type")
	}
	p := ptr(x, C)
	a := js.Global().Get("f").Get(string(m.c[p : p+n]))
	ln := a.Get("length")
	r = mk(C, k(ln.Int()))
	js.CopyBytesToGo(m.c[8+r<<2:], a)
	return dex(x, r)
}
func wrt(x, y k) (r k) { // x 1:y
	xt, yt, xn, yn := typs(x, y)
	if yt != C || yn == atom {
		panic("type")
	}
	if xt == S {
		x = str(x)
		xt, xn = typ(x)
	}
	if xt != C || xn == atom {
		panic("type")
	}
	if xn != 0 {
		xp, yp := ptr(x, C), ptr(y, C)
		s := string(m.c[xp : xp+xn])
		u := js.Global().Get("Uint8Array").New(int(yn))
		js.CopyBytesToJS(u, m.c[yp:yp+yn])
		js.Global().Get("f").Set(s, u)
		return dex(y, x)
	}
	cat := mk(N+2, atom)
	m.k[cat+2] = 12 + dyad
	return dex(asn(mks("ut"), spl(mkc('\n'), y), cat), x)
}
func trp(x, y k) (r k) {
	if t, n := typ(y); t == C && n != atom && n < 512 {
		p := ptr(y, C)
		js.Global().Call("hash", string(m.c[p:p+n]))
	}
	defer func() {
		if rc := recover(); rc != nil {
			s := "?"
			if t, o := rc.(string); o {
				s = t
			} else if t, o := rc.(error); o {
				s = t.Error()
			}
			r = mkb([]c(s))
		}
	}()
	return cal(x, enlist(y))
}

/* tinygo
// max memory, heap size must be set explicitly when compiling with tinygo
const maxmem = 1 << 21 // number of floats (16 MB)
var mem []float64
func main() {
	mem = make([]float64, maxmem)
	ini(mem[:1<<13])
	table[21] = red
}
func grw() {
	if len(m.f) >= maxmem {
		panic("mem")
	}
	m.f = mem[:2*len(m.f)]
}
*/
