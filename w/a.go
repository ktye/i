package main

import (
	// to be removed when prs is implemented
	"runtime/debug"
	"syscall"
)

// To read command line arguments without importing package os, see go/src/os/{proc.go,exec_windows.go}
// Windows needs special handling, see: exec_windows.go (init)
// 	p := syscall.GetCommandLine()
// 	cmd := syscall.UTF16ToString((*[0xffff]uint16)(unsafe.Pointer(p))[:])
//	cmd is a single string and might need special splitting

func main() {
	ini()
	var buf [1024]byte
	p := buf[:]
	for {
		// Read from stdin without os.Stdin
		// syscall.Stdin is not 0 on windows, but a call to GetStdHandle(-10)
		n, err := syscall.Read(syscall.Stdin, p)
		if err != nil {
			panic(err)
		}
		if n > 0 {
			do(p[:n])
		}
	}
}
func do(s []byte) {
	defer func() {
		if c := recover(); c != nil {
			println(string(debug.Stack()))
			if s, o := c.(string); o {
				println(s)
			} else if e, o := c.(error); o {
				println(e.Error())
			}
		}
	}()
	r := kst(evl(prs(mk(C, k(len(s))))))
	n := m.k[r] & atom
	println(string(m.c[8+r<<2 : 8+n+r<<2]))
	dec(r)
}

// type conversions between go and k (used here and in k_test.go)

func K(x interface{}) k { // convert go value to k type, returns 0 on error
	kstr := func(dst k, s string) {
		u, n := uint64(0), len(s)
		if n > 8 {
			n = 8
		}
		for i := 0; i < n; i++ { // byte order independend
			u |= uint64(s[i]) << (8 * c(7-i))
		}
		mys(dst, u)
	}
	var r k
	switch a := x.(type) {
	case bool:
		r = mk(C, atom)
		m.c[8+r<<2] = 0
		if a {
			m.c[8+r<<2] = 1
		}
	case byte:
		r = mk(C, atom)
		m.c[8+r<<2] = a
	case int:
		r = mk(I, atom)
		m.k[2+r] = k(a)
	case float64:
		r = mk(F, atom)
		m.f[1+r>>1] = a
	case complex128:
		r = mk(Z, atom)
		m.z[1+r>>2] = a
	case string:
		r = mk(S, atom)
		kstr(8+r<<2, a)
	case []bool:
		buf := make([]byte, len(a))
		for i, v := range a {
			if v {
				buf[i] = 1
			}
		}
		return K(buf)
	case []byte:
		r = mk(C, k(len(a)))
		for i, v := range a {
			m.c[8+i+int(r<<2)] = v
		}
	case []int:
		r = mk(I, k(len(a)))
		for i, v := range a {
			m.k[2+i+int(r)] = k(v)
		}
	case []float64:
		r = mk(F, k(len(a)))
		for i, v := range a {
			m.f[1+i+int(r>>1)] = v
		}
	case []complex128:
		r = mk(Z, k(len(a)))
		for i, v := range a {
			m.z[1+i+int(r>>2)] = v
		}
	case []string:
		r = mk(S, k(len(a)))
		for i := range a {
			kstr(8+8*k(i)+r<<2, a[i])
		}
	case []interface{}:
		if len(a) == 1 { // collapse list of atom to single element vector
			rr := K(a[0])
			t, n := typ(rr)
			if n == atom { // TODO: allow ,d?
				r = rr
				m.k[r] = t<<28 | 1
				return r
			} else {
				dec(rr)
			}
		}
		r = mk(L, k(len(a)))
		for i, v := range a {
			u := K(v)
			m.k[2+i+int(r)] = u
		}
	case [2]interface{}:
		key := K(a[0])
		val := K(a[1])
		_, nk := typ(key)
		_, nv := typ(val)
		if nk != nv {
			return 0
		}
		r = mk(D, atom)
		m.k[2+r] = key
		m.k[3+r] = val
	}
	return r
}
func G(x k) interface{} { // convert k value to go type (returns nil on error)
	t, n := typ(x)
	if n == atom {
		switch t {
		case C:
			return c(m.c[8+x<<2])
		case I:
			return int(i(m.k[2+x]))
		case F:
			return m.f[1+x>>1]
		case Z:
			return m.z[1+x>>2]
		case S:
			return ustr(sym(8 + x<<2))
		case D:
			return [2]interface{}{G(m.k[2+x]), G(m.k[3+x])}
		}
	} else {
		switch t {
		case C:
			r := make([]byte, n)
			for i := range r {
				r[i] = c(m.c[8+i+int(x<<2)])
			}
			return r
		case I:
			r := make([]int, n)
			for i := range r {
				r[i] = int(int32(m.k[2+i+int(x)]))
			}
			return r
		case F:
			r := make([]f, n)
			for i := range r {
				r[i] = m.f[1+i+int(x>>1)]
			}
			return r
		case Z:
			r := make([]complex128, n)
			for i := range r {
				r[i] = m.z[1+i+int(x>>2)]
			}
			return r
		case S:
			r := make([]string, n)
			for i := range r {
				r[i] = ustr(sym(8 + 8*k(i) + x<<2))
			}
			return r
		case L:
			r := make([]interface{}, n)
			for i := range r {
				r[i] = G(m.k[2+i+int(x)])
			}
			return r
		}
	}
	return nil
}
