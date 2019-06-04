package w

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type z = complex128
type l = []interface{}
type d = [2]interface{}

func TestIni(t *testing.T) {
	ini()
	st := Stats()
	if st.UsedBlocks() != 1 {
		t.Fatal()
	}
	println("used blocks", st.UsedBlocks())
	//mk(1, 9000)
	//pfl()
	//xxd()
}

func TestNumMonad(t *testing.T) {
	ini()
	xv := []interface{}{c(3), []c{3, 5}, -5, []int{3, -9}, 3.2, []f{-3.5, 2.9, 0}, 2 - 4i, []z{4 - 2i, 3 + 4i}}
	testCases := []struct {
		f func(k) k
		s string
		r []interface{}
	}{
		{neg, "-", []interface{}{c(253), []c{253, 251}, 5, []int{-3, 9}, -3.2, []f{3.5, -2.9, -0}, -2 + 4i, []z{-4 + 2i, -3 - 4i}}},
		{fst, "*", []interface{}{c(3), c(3), -5, 3, 3.2, -3.5, 2 - 4i, 4 - 2i}},
	}
	for j, tc := range testCases {
		for i := range xv {
			x := K(xv[i])
			if x == 0 {
				t.Fatalf("cannot import go type %T", xv[i])
			}
			y := tc.f(x)
			r := Go(y)
			fmt.Printf("%s[%v] = %v\n", tc.s, xv[i], r)
			if !reflect.DeepEqual(r, tc.r[i]) {
				t.Fatalf("[%d/%d]: expected: %v got %v (@%d)\n", j, i, tc.r[i], r, y)
			}
			xdec(x)
			xdec(y)
			if m.k[x]>>28 != 0 || m.k[y]>>28 != 0 {
				panic("x|y is not free")
			}
			if u := Stats().UsedBlocks(); u != 1 {
				t.Fatalf("leak")
			}
		}
	}
}
func TestMonad(t *testing.T) {
	ini()
	testCases := []struct {
		f    func(k) k
		s    string
		x, r interface{}
	}{
		{til, "!", 3, []int{0, 1, 2}},
		// TODO !overloads
		{fst, "*", l{3, 4, 5}, 3},
		{fst, "*", "alpha", "alpha"},
		{fst, "*", l{"alpha"}, "alpha"},
		{fst, "*", d{l{"x", "y"}, l{[]int{5, 3}, 4}}, []int{5, 3}},
		{fst, "*", d{[]string{"x", "y"}, []int{7, 2}}, 7},
		// TODO fst func
	}
	for j, tc := range testCases {
		println("NEWTC")
		x := K(tc.x)
		if x == 0 {
			t.Fatalf("cannot import go type %T", tc.x)
		}
		fmt.Println("tc.x", tc.x)
		xxd()
		y := tc.f(x)
		r := Go(y)
		if !reflect.DeepEqual(r, tc.r) {
			xxd()
			t.Fatalf("monad[%d]: expected: %v got %v (@%d)\n", j, tc.r, r, y)
		}
		xxd()
		println("xdec")
		xdec(x)
		xxd()
		println("ydec")
		xdec(y)
		xxd()
		if m.k[x]>>28 != 0 || m.k[y]>>28 != 0 {
			panic("x|y is not free")
		}
		if u := Stats().UsedBlocks(); u != 1 {
			t.Fatalf("leak: %d", u)
		}
	}
}
func pfl() {
	for i := 4; i < 32; i++ {
		println(i, strconv.FormatUint(uint64(m.k[i]), 16), strconv.FormatUint(uint64(m.k[i]<<2), 16))
	}
}
func xdec(x k) {
	if m.k[x]>>28 != 0 {
		dec(x)
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

type Bucket struct {
	Type       uint32
	Used, Free uint32 // num blocks
	Net        uint32
}
type MemStats map[uint32]Bucket

func (b Bucket) Overhead() uint32 {
	return b.Used*uint32(1<<b.Type) - b.Net
}
func (s MemStats) UsedBlocks() (t uint32) {
	for _, b := range s {
		t += b.Used
	}
	return t
}
func Stats() MemStats {
	st := make(MemStats)
	a := uint32(0)
	o := uint32(0)
	for a < 1<<(m.k[2]-2) {
		//ax := strconv.FormatUint(uint64(a<<2), 16)
		tp := m.k[a] >> 28
		//print("block ", a, " 0x", ax, " t=", tp)
		if tp == 0 {
			t := m.k[a]
			//print(" free bt=", t)
			if t < 4 || t > 31 {
				println(a, t)
				panic("size")
			}
			b := st[t]
			b.Type = t
			b.Free++
			st[t] = b
			o = 1 << (t - 2)
		} else {
			tt, n := typ(a)
			t := bk(tt, n)
			//print(" used type=", tt, " num=", n, " bt=", t)
			if t < 4 || t > 31 {
				println(a, t)
				panic("size")
			}
			b := st[t]
			b.Type = t
			b.Used++
			if n == atom {
				n = 1
			}
			b.Net += n * lns[tp]
			st[t] = b
			o = 1 << (t - 2)
		}
		a += o
		//print(" +", o, "\n")
	}
	return st
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
func Go(x k) interface{} { // convert k value to go type (returns nil on error)
	str := func(x k, j int) string {
		buf := make([]byte, 8)
		n := 0
		for i := range buf {
			if v := m.c[8+8*j+int(x<<2)+i]; v != 0 {
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
