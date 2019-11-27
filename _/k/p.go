package main

import (
	"image"

	"github.com/ktye/plot"
)

// plot
// before a pure k plot library is done, ktye/plot is used as a backend
//
// plot y
//    F        /single line in x-y plot with default x axis !#F
//    (F;F;..) /multiple lines
//    Z        /polar diagram with points of a single data set
//    (Z;Z;..) /polar diagram with points for multiple data sets
//
// x plot y
//  x same size as y or a single vector that extends for every line in y
//  F plot F  / x-y plot
//  F plot Z  / amplitude/phase-over-x plot
//
// return
//  plot by default rasterizes the image and writes it to 9: (which displays on the root window)
//  plot `!x     and  x plot `!y     return a dict with the lowlevel interface
//  plot d       plots dict with the low level interface
//  ` plot d     return as image (w;I)
//  plot `x      and `x plot`y   lookup the symbols x and y and use the names as axis labels
//
// low level interface (dict) (simplified version of github.com/ktye/plot/blob/master/plot.go)
//  d.type:        `xy|`foto|`polar|`ring|`ampang|`raster
//  d.lim:         xmin xmax ymin ymax zmin zmax (F length 6 or 4(zmin:zmax:0) or 3(ymin:0) or 2(xmin;xmax) or 1(ymax))
//  d[`x`y`t]:     xlabel, ylabel, title
//  d[`w`y]:       image width, height in pixels
//  d.l            table each row a line dict (l)
//  d.l[0].x       x vector F
//  d.l[0].y       y vector F or Z
//  d.l[0].w       line width / marker size
//  d.l[0].c       color index
//  d.l[0].m       0|1|2 line|marker|both
//
// examples
//  plot rand 1000a
//  (.1*!10)plot 10 10#100 rand 5

func plo(x k) (r k) { return plt(inc(null), x) }
func plt(x, y k) (r k) {
	xt, yt, _, yn := typs(x, y)
	xs, ys := k(0), k(0)
	x, xt, _, xs = symbolic(x)
	y, yt, yn, ys = symbolic(y)
	toImg, toDict := false, false
	if yt == A && yn == atom {
		xk := m.k[2+y]
		if m.k[xk]&atom != 1 { // plot d | ` plot d
			toImg = !match(x, null)
			return dex(x, pld(y, toImg))
		} else if m.k[xk]>>28 == S { // plot `!x  x plot `!y
			toDict = true
			y = fst(val(y))
			yt, yn = typ(y)
		} else {
			panic("type")
		}
	}
	if xt == I {
		x, xt = to(x, F), F
	} else if xt == L {
		x = toflist(x)
	}
	if yt == I {
		y, yt = to(y, F), F
	} else if yt == L {
		y = toflist(y)
	}
	dt := yt
	if dt == L {
		if m.k[2+y]&atom == 0 {
			panic("empty")
		}
		dt = m.k[m.k[2+y]] >> 28
	}
	z := false
	if dt == Z {
		z = true
	} else if dt != F {
		panic("type")
	}
	if yt == atom {
		y, yn = enl(y), 1
	}
	if yt != L {
		y, yt, yn = enlist(y), L, 1
	}
	dk := enl(mks("type"))
	dv := mk(L, 0)
	if match(x, null) && z {
		dv = lcat(dv, mks("polar"))
	} else if z {
		dv = lcat(dv, mks("ampang"))
	} else {
		dv = lcat(dv, mks("xy"))
	}
	if xs != 0 {
		dk = cat(dk, mks("x"))
		dv = lcat(dv, xs)
	}
	if ys != 0 {
		dk = cat(dk, mks("y"))
		dv = lcat(dv, ys)
	}
	dk = cat(cat(dk, mks("w")), mks("h"))
	dv = lcat(lcat(dv, mki(800)), mki(600))
	l := mk(L, yn)
	for i := k(0); i < yn; i++ {
		var xi k
		if xt == N {
			xi = inc(null)
		} else if xt == F {
			xi = inc(x)
		} else if xt == L {
			xi = to(inc(m.k[2+x+i]), F)
		}
		m.k[2+l+i] = line(i, xi, inc(m.k[2+i+y]))
	}
	l = uf(l)
	dk = cat(dk, mks("l"))
	dv = lcat(dv, l)
	if toDict {
		return key(dk, dv)
	}
	return pld(key(dk, dv), toImg)
}
func symbolic(x k) (r, n, t, s k) { // lookup symbolic data and store name as axis label
	t, n = typ(x)
	if t == S && n == atom {
		s = inc(x)
		x = lup(x)
		t, n = typ(x)
	}
	return x, t, n, s
}
func toflist(x k) (r k) {
	n := m.k[x] & atom
	r = mk(L, n)
	for i := k(0); i < n; i++ {
		if m.k[m.k[2+i+x]]>>28 == I {
			m.k[2+r+i] = to(inc(m.k[2+i+x]), F)
		} else {
			m.k[2+r+i] = inc(m.k[2+i+x])
		}
	}
	return dex(x, r)
}
func line(i, x, y k) (r k) {
	lk := enl(mks("id"))
	lv := mk(L, 1)
	m.k[2+lv] = mki(i)
	if !match(x, null) {
		lk = cat(lk, mks("x"))
		lv = lcat(lv, x)
	} else {
		dec(x)
	}
	lk = cat(lk, mks("y"))
	lv = lcat(lv, y)
	r = key(lk, lv)
	return r
}
func pld(d k, toImg bool) (r k) {
	var p plot.Plot
	p.Type = plot.PlotType(tostring(dlup(inc(d), "type")))
	p.Limits = tolimits(dlup(inc(d), "lim"))
	p.Xlabel = tostring(dlup(inc(d), "x"))
	p.Ylabel = tostring(dlup(inc(d), "y"))
	p.Title = tostring(dlup(inc(d), "t"))
	w := toint(dlup(inc(d), "w"), 800)
	h := toint(dlup(inc(d), "h"), 600)

	l := dlup(d, "l")
	lt, ln := typ(l)
	if l == 0 || lt != A || ln == atom {
		panic("type")
	}
	p.Lines = make([]plot.Line, ln)
	for i := k(0); i < ln; i++ {
		p.Lines[i] = toline(atx(inc(l), mki(i)))
	}
	dec(l)

	ip, err := plot.Plots([]plot.Plot{p}).IPlots(w, h)
	if err != nil {
		panic(err)
	}
	im := plot.Image(ip, nil, w, h).(*image.RGBA)

	n := k(im.Bounds().Dx() * im.Bounds().Dy())
	b := mk(I, n)
	c := ptr(b, C)
	copy(m.c[c:c+4*n], im.Pix[:4*n])
	if toImg {
		return l2(mki(k(w)), b)
	}
	return drw(mki(k(w)), b)
}
func tostring(x k) (r s) {
	if x == 0 {
		return ""
	}
	x = str(x)
	t, n := typ(x)
	if t != C {
		dec(x)
		return ""
	}
	p := 8 + x<<2
	r = s(m.c[p : p+n])
	dec(x)
	return r
}
func toint(x k, def int) (r int) {
	if x == 0 {
		return def
	}
	t, n := typ(x)
	if t != I || n != atom {
		dec(x)
		return def
	}
	r = int(m.k[2+x])
	dec(x)
	return r
}
func toline(d k) (l plot.Line) {
	if t, n := typ(d); t != A || n != atom {
		panic("type")
	}
	l.Id = toint(dlup(inc(d), "id"), 0)
	l.X = tofloats(dlup(inc(d), "x"))
	y := dlup(inc(d), "y")
	if t := m.k[y] >> 28; y != 0 && t == F || t == I {
		l.Y = tofloats(y)
	} else if t == Z {
		l.C = tocmplxs(y)
	} else if y != 0 {
		panic("type")
	}
	if l.X == nil && len(l.Y) > 0 {
		l.X = tofloats(jota(k(len(l.Y))))
	}
	if len(l.X) != len(l.Y) {
		panic("size")
	}
	w := toint(dlup(inc(d), "w"), 0)
	m := toint(dlup(inc(d), "m"), 0)
	if m == 1 || m == 2 {
		l.Style.Marker.Size = w
		l.Style.Marker.Marker = plot.PointMarker
	}
	if m == 2 {
		l.Style.Line.Width = w
	}
	c := toint(dlup(inc(d), "c"), 0)
	l.Style.Marker.Color = c
	l.Style.Line.Color = c
	dec(d)
	return l
}
func dlup(d k, s s) (r k) {
	kv := m.k[2+d]
	kt, kn := typ(kv)
	if kt != S {
		panic("type")
	}
	p := mks(s)
	pp := m.k[2+p]
	for i := k(0); i < kn; i++ {
		if m.k[2+kv+i] == pp {
			return decr(p, d, atx(inc(m.k[3+d]), mki(i)))
		}
	}
	return decr(p, d, 0)
}
func tofloats(x k) (r []f) {
	if x == 0 {
		return nil
	}
	t, n := typ(x)
	if t == I {
		x, t = to(x, F), F
	} else if t != F {
		dec(x)
		return nil
	}
	if n == atom {
		x, n = enl(x), 1
	}
	p := ptr(x, F)
	r = make([]f, n)
	copy(r, m.f[p:p+n])
	dec(x)
	return r
}
func tocmplxs(x k) (r []complex128) {
	n := atm1(m.k[x] & atom)
	r = make([]complex128, n)
	p := ptr(x, Z)
	copy(r, m.z[p:p+n])
	dec(x)
	return r
}
func tolimits(x k) (l plot.Limits) {
	if x == 0 {
		return l
	}
	if m.k[x]>>28 == I {
		x = to(x, F)
	}
	t, n := typ(x)
	if t != F {
		dec(x)
		return l
	}
	p := ptr(x, F)
	switch {
	case n == 1:
		l.Ymax = m.f[p+1]
	case n == 2:
		l.Ymin, l.Ymax = m.f[p], m.f[p+1]
	case n == 3:
		l.Xmin, l.Xmax, l.Ymax = m.f[p], m.f[p+1], m.f[p+2]
	case n == 4:
		l = plot.Limits{false, m.f[p], m.f[p+1], m.f[p+2], m.f[p+3], 0, 0}
	case n == 5:
		l = plot.Limits{false, m.f[p], m.f[p+1], m.f[p+2], m.f[p+3], 0, m.f[p+4]}
	case n == 6:
		l = plot.Limits{false, m.f[p], m.f[p+1], m.f[p+2], m.f[p+3], m.f[p+4], m.f[p+5]}
	}
	dec(x)
	return l
}
