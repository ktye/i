package main

import "sync"

const np = k(1024)

func ptil(x k) (r k) { // !x
	t, n := typ(x)
	if n != atom || t > Z || n < np {
		return til(x)
	}
	nn := idx(x, t)
	if nn <= 0 {
		return til(x)
	}
	return dex(x, pto(pjota(k(nn)), t))
}
func pjota(n k) (r k) {
	var wg sync.WaitGroup
	r = mk(I, n)
	for i := k(0); i < n; i += np {
		wg.Add(1)
		j := np
		if i+j > n {
			j = n - i
		}
		go func(x, o, n k) {
			for i := k(0); i < n; i++ {
				m.k[x+i] = o + i
			}
			wg.Done()
		}(2+r+i, i, j)
	}
	wg.Wait()
	return r
}

func pfkk(xp, yp, n k, f f1) {
	var wg sync.WaitGroup
	for i := k(0); i < n; i += np {
		wg.Add(1)
		j := np
		if i+j > n {
			j = n - i
		}
		go func(x, y, n k) {
			for i := k(0); i < n; i++ {
				f(x+i, y+i)
			}
			wg.Done()
		}(xp+i, yp+i, j)
	}
	wg.Wait()
}
func pfkkk(rp, xp, yp, dx, dy, n k, f f2) {
	var wg sync.WaitGroup
	for i := k(0); i < n; i += np {
		wg.Add(1)
		j := np
		if i+j > n {
			j = n - i
		}
		go func(r, x, y, dx, dy, n k) {
			for i := k(0); i < n; i++ {
				f(r+i, x, y)
				x += dx
				y += dy
			}
			wg.Done()
		}(rp+i, xp+i, yp+i, dx, dy, j)
	}
	wg.Wait()
}

func pto(x, rt k) (r k) {
	if rt == 0 || rt >= L {
		return x
	}
	t, n := typ(x)
	if rt == t {
		return x
	} else if t == L {
		r = mk(L, n)
		for i := k(0); i < n; i++ {
			m.k[2+i+r] = pto(inc(m.k[2+x+i]), rt)
		}
		return dex(x, uf(r))
	}
	var g func(k, k)
	if t == S && rt == I { // for symbol conversion to bool
		g = func(r, x k) {
			if m.k[x] == 0 {
				m.k[r] = 0
			} else {
				m.k[r] = 1
			}
		}
	} else {
		g = tox[4*(t-1)+rt-1]
	}
	r = mk(rt, n)
	n = atm1(n)
	xp, rp := ptr(x, t), ptr(r, rt)
	pfkk(rp, xp, n, g)
	return dex(x, r)
}

func pnm(x, rt k, fx []f1) (r k) { // numeric monad (parallel)
	t, n := typ(x)
	if n < np {
		return snm(x, rt, fx)
	}
	min := C
	if fx[C] == nil {
		min = I
	}
	if fx[I] == nil {
		min = F
	} // TODO: Z only for ff == nil ?
	if min > t { // uptype x
		x, t = to(x, min), min
	}
	if t == Z && fx[Z] == nil { // e.g. real functions
		x, t = to(x, F), F
	}
	if m.k[1+x] == 1 && t < L {
		r = inc(x)
	} else {
		r = mk(t, n)
	}
	n = atm1(n)
	switch t {
	case L:
		for j := k(0); j < k(n); j++ {
			m.k[r+2+j] = snm(inc(m.k[j+2+x]), rt, fx)
		}
	case A:
		if r != x {
			m.k[2+r] = inc(m.k[2+x])
		}
		m.k[3+r] = snm(m.k[3+x], rt, fx)
	case C, I, F, Z:
		rp, xp, f := ptr(r, t), ptr(x, t), fx[t]
		pfkk(rp, xp, n, f)
	default:
		panic("type")
	}
	if rt != 0 && t > rt && t < L { // only down-type
		r = pto(r, rt)
	}
	return dex(x, r)
}
func pns(rp, xp, yp, t, xn, yn, c k, f f2) { // parallel numeric scalar dyads
	var dx, dy, n k
	switch c {
	case 0: // v f v
		dx, dy, n = 1, 1, xn
	case 1: // a f v
		dx, dy, n = 0, 1, yn
	case 2: // v f a
		dx, dy, n = 1, 0, xn
	default:
		panic("assert")
	}
	if n < np {
		sns(rp, xp, yp, t, xn, yn, c, f)
		return
	}
	pfkkk(rp, xp, yp, dx, dy, n, f)
}
