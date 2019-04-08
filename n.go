package i

import (
	"reflect"
	"sort"
)

type fz1 func(z) z
type fz2 func(z, z) z

func zter(c bool, a, b z) z {
	if c {
		return a
	}
	return b
}
func nm(x v, fz fz1) v {
	if d, ok := md(x); ok {
		for i, v := range d.v {
			d.v[i] = nm(v, fz)
		}
		return d.mp()
	}
	if iv, ok := x.(l); ok {
		r := make(l, len(iv))
		for i := range r {
			r[i] = nm(iv[i], fz)
		}
		u, _ := uf(r)
		return u
	}
	z, vec, t := nv(x)
	for i, x := range z {
		z[i] = fz(x)
	}
	return vn(z, vec, t)
}

func nd(x, y v, fz fz2) v {
	if r, ok := ndic(x, y, fz); ok {
		return r
	}

	xn, yn := ln(x), ln(y)
	switch {
	case xn >= 0 && yn >= 0 && xn != yn:
		return e("length")
	case xn < 0 && yn >= 0:
		x, xn = rsh(zi(yn), x), yn
	case yn < 0 && xn >= 0:
		y, yn = rsh(zi(xn), y), xn
	}
	xl, yl := false, false
	if xn >= 0 && rtyp(x).Elem().Kind() == reflect.Interface {
		xl = true
	}
	if yn >= 0 && rtyp(y).Elem().Kind() == reflect.Interface {
		yl = true
	}
	if xl || yl {
		r := make(l, xn) // TODO: make custom interface type, if both have the same type
		for i := range r {
			r[i] = nd(at(x, i), at(y, i), fz)
		}
		u, _ := uf(r)
		return u
	}

	xz, xvec, xt := nv(x)
	yz, yvec, yt := nv(y)
	n1, n2 := len(xz), len(yz)
	if n1 == 0 || n2 == 0 {
		if xt == yt && xt != nil {
			return ms(xt, 0).Interface()
		}
		return l{}
	}
	if n1 < 0 && n2 > 1 {
		xz, n1 = zvn(xz[0], n2), n2
	} else if n1 > 1 && n2 < 0 {
		yz, n2 = zvn(yz[0], n1), n1
	}
	if n1 != n2 {
		e("length")
	}
	for i := range xz {
		xz[i] = fz(xz[i], yz[i])
	}
	vec := false
	if xvec || yvec {
		vec = true
	}
	if xt == yt {
		return vn(xz, vec, xt)
	}
	return vn(xz, vec, nil)
}
func ndic(x, y v, fz fz2) (v, bool) {
	xd, isx := md(x)
	yd, isy := md(y)
	if isx == false && isy == false {
		return nil, false
	}
	if isx && isy {
		// That could be an identity element, depending on the verb.
		// oK just fills the other value without applying the function.
		zero := 0.0 // If there is no agreement, i can do what i want.
		for i, k := range xd.k {
			yi, v := yd.at(k)
			if yi < 0 {
				xd.v[i] = nd(xd.v[i], zero, fz)
			} else {
				xd.v[i] = nd(xd.v[i], v, fz)
			}
		}
		for i, k := range yd.k {
			if idx, _ := xd.at(k); idx < 0 {
				xd.k = append(xd.k, k)
				xd.v = append(xd.v, nd(zero, yd.v[i], fz))
			}
		}
		if xd.t != yd.t {
			xd.t = nil
		}
		return xd.mp(), true
	}
	// d+v is not allowed, but d+a is.
	d := xd
	a := y
	flip := false
	if isx == false {
		d = yd
		a = x
		flip = true
	}
	if ln(a) >= 0 {
		e("type") // d+v is not allowed
	}
	for i, _ := range d.k {
		x, y = d.v[i], a
		if flip {
			x, y = y, x
		}
		d.v[i] = nd(x, y, fz)
	}
	return xd.mp(), true
}

func nv(x v) (zv, bool, rT) { // import any number or numeric vector types
	switch t := x.(type) {
	case z:
		return zv{t}, false, rTz
	case zv:
		r := make(zv, len(t))
		copy(r, t)
		return t, true, rTz
	case s:
		e("type")
	}
	v := rval(x)
	if v.Kind() == reflect.Slice {
		n := v.Len()
		z := reflect.Zero(v.Type().Elem())
		if z.Type().ConvertibleTo(rTf) {
			r := make(zv, n)
			for i := range r {
				r[i] = complex(v.Index(i).Convert(rTf).Float(), 0)
			}
			return r, true, z.Type()
		} else if z.Type().ConvertibleTo(rTz) {
			r := make(zv, n)
			for i := range r {
				r[i] = v.Index(i).Convert(rTz).Complex()
			}
			return r, true, z.Type()
		} else if z.Type().ConvertibleTo(rTb) {
			r := make(zv, n)
			for i := range r {
				b := v.Index(i).Convert(rTb).Bool()
				if b {
					r[i] = 1
				}
			}
			return r, true, z.Type()
		}
		e("type")
	}

	if v.Type().ConvertibleTo(rTf) {
		return zv{complex(v.Convert(rTf).Float(), 0)}, false, v.Type()
	}
	if v.Type().ConvertibleTo(rTz) {
		return zv{v.Convert(rTz).Complex()}, false, v.Type()
	}
	if v.Type().ConvertibleTo(rTb) {
		b := v.Convert(rTb).Bool()
		r := complex(0, 0)
		if b {
			r = 1
		}
		return zv{r}, false, v.Type()
	}
	e("type")
	return nil, false, rTf
}
func vn(x zv, vec bool, t rT) v { // convert numbers back to original type
	if t == rTz || t == nil {
		if vec {
			return x
		}
		return x[0]
	}
	scalar := func(x z) rV {
		if t.Kind() == reflect.Bool {
			b := false
			if x != 0 {
				b = true
			}
			return rval(b).Convert(t)
		} else if t.Kind() == reflect.Complex128 {
			return rval(x).Convert(t)
		}
		return rval(real(x)).Convert(t)
	}
	if vec == false {
		return scalar(x[0]).Interface()
	}
	n := len(x)
	r := ms(t, n)
	for i := range x {
		r.Index(i).Set(scalar(x[i]))
	}
	return r.Interface()
}
func sn(v v) (zv, bool, bool) { // import strings as numbers; for =<>
	s, n, _, o := sy(v)
	if o == false {
		return nil, false, false
	}
	if n < 0 {
		return zv{0}, false, true
	}
	m := strmap(s)
	r := make(zv, n)
	for i := range s {
		r[i] = m[s[i]]
	}
	return r, true, true
}
func sn2(x, y v) (v, v) { // map strings to floats
	sx, nx, _, o := sy(x)
	if o == false {
		return x, y
	}
	sy, ny, _, o := sy(y)
	if o == false {
		return x, y
	}
	vec := true
	if nx < 0 && ny < 0 {
		vec, nx, ny = false, 1, 1
	} else if nx < 0 {
		sx, nx = rsh(zi(ny), sx).(sv), ny
	} else if ny < 0 {
		sy, ny = rsh(zi(nx), sy).(sv), nx
	} else if nx != ny {
		e("length")
	}
	b := make(sv, nx+ny)
	copy(b, sx)
	copy(b[nx:], sy)
	m := strmap(b)
	rx := make(zv, nx)
	for i := range sx {
		rx[i] = m[sx[i]]
	}
	ry := make(zv, ny)
	for i := range sy {
		ry[i] = m[sy[i]]
	}
	if !vec {
		return rx[0], ry[0]
	}
	return rx, ry
}
func strmap(x sv) map[s]z { // map s to f uniq and comparable
	idx := make(zv, len(x))
	for i := 0; i < len(x); i++ {
		idx[i] = zi(i)
	}
	c := cp(x).(sv)
	u := grades{sort.StringSlice(c), idx}
	sort.Sort(u)
	m := make(map[s]z)
	w := complex(0, 0)
	for i := range u.idx {
		if i == 0 || c[i] != c[i-1] {
			m[c[i]] = w
			w += 1.0
		}
	}
	return m
}
func zvn(x z, n int) zv {
	r := make(zv, n)
	for i := range r {
		r[i] = x
	}
	return r
}
func zi(n int) z { return complex(float64(n), 0) }
func nl(x v) l {
	v := rval(x)
	if v.Kind() == reflect.Slice {
		r := make(l, v.Len())
		for i := range r {
			rval(x).Index(i).Set(v.Index(i))
		}
		return r
	} else {
		return l{x}
	}
}
