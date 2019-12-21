package main

import "sync"

const np = 1024

func pfkk(xp, yp, n k, f func(x, y k)) {
	var wg sync.WaitGroup
	for i := k(0); i < n; i += np {
		wg.Add(1)
		j := i + np
		if j > n {
			j = n
		}
		go func(x, y, n k) {
			for i := k(0); i < n; i++ {
				f(x+i, y+i)
			}
		}(xp+i, yp+i, j)
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
		return nm(x, rt, fx)
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
