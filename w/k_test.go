package w

import (
	"fmt"
	"reflect"
	"testing"
)

type z = complex128

func TestIni(t *testing.T) {
	ini()
	//mk(1, 9000)
	//pfl()
	//xxd()
}

func TestMonad(t *testing.T) {
	ini()
	testCases := []struct {
		f    func(k) k
		x, r interface{}
	}{
		{til, 3, []int{0, 1, 2}},
		{neg, byte(3), byte(253)},
		{neg, []byte{0, 1, 2}, []byte{0, 255, 254}},
		{neg, 2, -2},
		{neg, []int{2, 3}, []int{-2, -3}},
		{neg, 2.1, -2.1},
		{neg, []f{2.3, -3.4}, []f{-2.3, 3.4}},
		{neg, 2 + 3i, -2 - 3i},
		{neg, []z{1 + 2i, -3 + 4i}, []z{-1 - 2i, 3 - 4i}},
		{inv, 4.0, 0.25},
		{inv, 4, 0.25},
	}
	for j, tc := range testCases {
		x := K(tc.x)
		if x == 0 {
			t.Fatal()
		}
		r := Go(tc.f(x))
		fmt.Printf("f(%p)[%v] = %v\n", tc.f, tc.x, r)
		if !reflect.DeepEqual(r, tc.r) {
			t.Fatalf("monad[%d]: expected: %v got %v\n", j, tc.r, r)
		}
	}
}
func pfl() {
	for i := 4; i < 32; i++ {
		println(i, m.k[i])
	}
}
func xxd() { // memory dump
	t := [16]c{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 97, 98, 99, 100, 101, 102}
	s2 := func(x c) (c, c) { return t[x>>4], t[x&0xF] }
	l := make([]c, 49)
	for i := 0; i < len(l); i++ {
		l[i] = 32
	}
	n := 0
	u := k(0)
	e := true
	h := 0
	s := make([]c, 32)
	s[3] = '#'
	var tp, tn, rc k
	for i := 0; i < len(m.c); i += 2 {
		if i == h {
			tp, tn = typ(k(i))
			b := [8]c{'x', 'i', 'f', 'z', 's', 'g', 'd', 'l'}
			s[0] = b[tp]
			if tn != atom {
				s[0] -= 32
			}
			bt := k(0)
			if tp == 0 {
				bt = m.k[k(i)>>2]
			} else {
				bt = buk(m.k[k(i)>>2])
			}
			s[1], s[2] = s2(c(bt))
			rc = m.k[1+k(i)>>2]
			h += 1 << bt
		}
		if n == 0 {
			l[0], l[1] = s2(c(u >> 24))
			l[2], l[3] = s2(c(u >> 16))
			l[4], l[5] = s2(c(u >> 8))
			l[6], l[7] = s2(c(u))
			u += 16
			n = 8
		}
		l[n+1], l[n+2] = s2(m.c[i])
		l[n+3], l[n+4] = s2(m.c[i+1])
		if m.c[i] != 0 || m.c[i+1] != 0 {
			e = false
		}
		n += 5
		if n == 48 {
			n = 0
			if !e {
				print(string(l))
				if s[0] != 0 {
					if tn >= 0 {
						print(string(s[:4]), tn, ";", rc)
					} else {
						print(string(s[:3]), ";", rc)
					}
					s[0], tn = 0, atom
				}
				println()
			}
			e = true
		}
	}
}

// type conversions between go and k:

func K(x interface{}) k { // convert go value to k type, returns 0 on error
	kstr := func(s string) [8]byte {
		var r [8]byte
		a := []byte(s)
		for i := range r {
			if i < len(a) {
				r[i] = a[i]
			}
		}
		return r
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
		m.f[1+r>>1] = real(a)
		m.f[2+r>>1] = imag(a)
	case string:
		buf := kstr(a)
		r = mk(S, atom)
		for i := range buf {
			m.c[8+i+int(r<<2)] = buf[i]
		}
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
			m.f[1+2*i+int(r>>1)] = real(v)
			m.f[2+2*i+int(r>>1)] = imag(v)
		}
	case []string:
		r = mk(S, k(len(a)))
		for i := range a {
			buf := kstr(a[i])
			for j := range buf {
				m.c[8+8*i+j+int(r<<2)] = buf[j]
			}
		}
	case []interface{}:
		r := mk(L, k(len(a)))
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
func Go(x k) interface{} { // convert k value to go type (returns nil on error)
	str := func(x k, j int) string {
		buf := make([]byte, 8)
		n := 0
		for i := range buf {
			if v := m.c[8+8*j+int(x)+i]; v != 0 {
				buf[i] = v
				n++
			} else {
				break
			}
		}
		return string(buf[:n])
	}
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
			return complex(m.f[1+x>>1], m.f[2+x>>1])
		case S:
			return str(x, 0)
		case D:
			return [2]interface{}{Go(m.k[2+x]), Go(m.k[3+x])}
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
				r[i] = complex(m.f[1+2*i+int(x>>1)], m.f[2+2*i+int(x>>1)])
			}
			return r
		case S:
			r := make([]string, n)
			for i := range r {
				r[i] = str(x, i)
			}
			return r
		case L:
			r := make([]interface{}, n)
			for i := range r {
				r[i] = Go(m.k[2+i+int(x)])
			}
			return r
		}
	}
	return nil
}
