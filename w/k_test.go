package w

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type (
	l  = []interface{}
	d  = [2]interface{}
	iv = []int
	sv = []string
)

func TestIni(t *testing.T) {
	//t.Skip()
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
	//t.Skip()
	ini()
	xv := []interface{}{c(3), []c{3, 5}, -5, iv{3, -9}, 3.2, []f{-3.5, 2.9, 0}, 2 - 4i, []z{4 - 2i, 3 + 4i}}
	c0 := c(0)
	testCases := []struct {
		f func(k) k
		s s
		r []interface{}
	}{
		{neg, "-", l{c(253), []c{253, 251}, 5, iv{-3, 9}, -3.2, []f{3.5, -2.9, -0}, -2 + 4i, []z{-4 + 2i, -3 - 4i}}},
		{fst, "*", l{c(3), c(3), -5, 3, 3.2, -3.5, 2 - 4i, 4 - 2i}},
		{rev, "|", l{c(3), []c{5, 3}, -5, iv{-9, 3}, 3.2, []f{0, 2.9, -3.5}, 2 - 4i, []z{3 + 4i, 4 - 2i}}},
		{not, "~", l{c0, []c{0, 0}, c0, []c{0, 0}, c0, []c{0, 0, 1}, c0, []c{0, 0}}},
		{enl, ",", l{[]c{3}, l{[]c{3, 5}}, iv{-5}, l{iv{3, -9}}, []f{3.2}, l{[]f{-3.5, 2.9, 0}}, []z{2 - 4i}, l{[]z{4 - 2i, 3 + 4i}}}},
		{cnt, "#", l{1, 2, 1, 2, 1, 3, 1, 2}},
		{tip, "@", l{"c", "C", "i", "I", "f", "F", "z", "Z"}},
		{evl, ".", xv},
	}
	occ := true // wrap x in inc dec
	for i := 0; i < 2; i++ {
		for j, tc := range testCases {
			for i := range xv {
				// fmt.Println("TC", xv[i])
				x := K(xv[i])
				if x == 0 {
					t.Fatalf("cannot import go type %T", xv[i])
				}
				if occ {
					inc(x)
				}
				y := tc.f(x)
				if occ {
					dec(x)
				}
				r := G(y)
				fmt.Printf("%s(%v) = %v\n", tc.s, xv[i], r)
				if !reflect.DeepEqual(r, tc.r[i]) {
					t.Fatalf("[%d/%d]: expected: %v got %v (@%d)\n", j, i, tc.r[i], r, y)
				}
				fpck("1")
				dec(y)
				if m.k[x]>>28 != 0 || m.k[y]>>28 != 0 {
					panic("x|y is not free")
				}
				if u := Stats().UsedBlocks(); u != 1 {
					t.Fatalf("leak")
				}
				fpck("2")
			}
		}
		occ = false
	}
}
func TestMonad(t *testing.T) {
	//t.Skip()
	ini()
	testCases := []struct {
		f    func(k) k
		s    s
		x, r interface{}
	}{
		{til, "!", 3, iv{0, 1, 2}},
		{til, "!", d{sv{"a", "b"}, iv{1, 2}}, sv{"a", "b"}},
		// TODO !overloads
		{fst, "*", l{3, 4, 5}, 3},
		{fst, "*", "alpha", "alpha"},
		{fst, "*", l{"alpha"}, "alpha"},
		{fst, "*", d{l{"x", "y"}, l{iv{5, 3}, 4}}, iv{5, 3}},
		{fst, "*", d{sv{"x", "y"}, iv{7, 2}}, 7},
		// TODO fst func
		{fms, "$", iv{1, 2, 3}, []c("1 2 3")},
		{fms, "$", l{1, 2, l{4, 5}}, []c("(1;2;(4;5))")},
		{fms, "$", d{l{5, 5.5}, iv{1, 2}}, []c("((5;5.5)!1 2)")},
		{rev, "|", l{}, l{}},
		{rev, "|", l{iv{3}}, l{iv{3}}},
		{rev, "|", l{1, 2}, l{2, 1}},
		{rev, "|", l{1, l{3, 4}}, l{l{3, 4}, 1}},
		{rev, "|", d{iv{1, 2}, iv{3, 4}}, d{iv{2, 1}, iv{4, 3}}},
		{rev, "|", d{sv{"alpha", "beta"}, l{3, iv{3, 5}}}, d{sv{"beta", "alpha"}, l{iv{3, 5}, 3}}},
		{wer, "&", iv{0, 0, 1, 1, 0, 1}, iv{2, 3, 5}},
		{wer, "&", iv{}, iv{}},
		{wer, "&", iv{2}, iv{0, 0}},
		{wer, "&", iv{1, 2, 3}, iv{0, 1, 1, 2, 2, 2}},
		{asc, "<", iv{1, 2, 3, 4}, iv{0, 1, 2, 3}},
		{asc, "<", iv{1, 4, 3, 2}, iv{0, 3, 2, 1}},
		{asc, "<", iv{4, 2, 3, 4}, iv{1, 2, 0, 3}},
		{asc, "<", []f{4, 1, 2}, iv{1, 2, 0}},
		{asc, "<", []c{6, 4, 2, 1}, iv{3, 2, 1, 0}},
		{asc, "<", []z{4, 1, 2}, iv{1, 2, 0}},
		{asc, "<", []z{0, 1 + 1i, 1, 2}, iv{0, 2, 1, 3}},
		{asc, "<", sv{"b", "ab", "a", "aa"}, iv{2, 3, 1, 0}},
		{dsc, ">", iv{1, 4, 3, 2}, iv{1, 2, 3, 0}},
		{enl, ",", "alpha", sv{"alpha"}},
		{enl, ",", l{1, 2, l{3, 4.5}}, l{l{1, 2, l{3, 4.5}}}},
		{enl, ",", d{iv{3, 4}, sv{"x", "y"}}, l{d{iv{3, 4}, sv{"x", "y"}}}},
		{cnt, "#", "alpha", 1},
		{cnt, "#", l{}, 0},
		{cnt, "#", l{1, 2, l{3, 4}}, 3},
		{cnt, "#", d{iv{3, 4}, sv{"x", "y"}}, 2},
		{tip, "@", l{}, ""},
		{tip, "@", d{iv{1, 2}, iv{3, 4}}, "a"},
		{evl, ".", l{"-", l{"-", 3}}, 3},
		{evl, ".", l{"-", l{"|", iv{3, 4}}}, iv{-4, -3}},
		{evl, ".", l{"-", iv{3, 4}}, iv{-3, -4}},
		{unq, "?", []c{1, 2, 43, 2}, []c{1, 2, 43}},
		{unq, "?", iv{1, 2, 3, 2}, iv{1, 2, 3}},
		{unq, "?", []f{5, 0, 0, 0, 8, 0, 0, 0, 5, 0, 0, 5}, []f{5, 0, 8}},
		{unq, "?", []z{0, 4i, 5i, 4i, 0, 3}, []z{0, 4i, 5i, 3}},
		{unq, "?", l{1, 2, 3, 1}, l{1, 2, 3}},
		{unq, "?", l{1i, l{2, sv{"a"}}, l{3, "b"}, l{2, sv{"a"}}, 1i}, l{1i, l{2, sv{"a"}}, l{3, "b"}}},
	}
	occ := true // wrap x in inc dec
	for i := 0; i < 2; i++ {
		for j, tc := range testCases {
			fmt.Println("TC", i, j, tc.s, tc.x, "occ", occ)
			x := K(tc.x)
			_ = Stats().UsedBlocks()
			if x == 0 {
				t.Fatalf("cannot import go type %T", tc.x)
			}
			if occ {
				inc(x)
			}
			y := tc.f(x)
			fpck("1")
			if occ {
				dec(x)
			}
			r := G(y)
			fmt.Printf("%s[%v] = %v\n", tc.s, tc.x, r)
			if !reflect.DeepEqual(r, tc.r) {
				t.Fatalf("monad[%d]: expected: %v got %v (@%d)\n", j, tc.r, r, y)
			}
			dec(y)

			fpck("2")
			if m.k[x]>>28 != 0 || m.k[y]>>28 != 0 {
				panic("x|y is not free")
			}
			if u := Stats().UsedBlocks(); u != 1 {
				t.Fatalf("leak: %d", u)
			}
		}
		occ = false
	}
}
func TestFms(t *testing.T) {
	//t.Skip()
	ini()
	testCases := []struct {
		x interface{}
		s s
	}{
		{[]c{}, `""`},
		{c('x'), `"x"`},
		{[]c{'x'}, `,"x"`},
		{c(28), "0x1c"},
		{[]c{28}, ",0x1c"},
		{[]c{0x1b, 0x5b, 0x5c}, "0x1b5b5c"},
		{[]c("alpha"), `"alpha"`},
		{"alpha", "`alpha"},
		{sv{"alpha"}, ",`alpha"},
		{sv{"a", "b", "c"}, "`a`b`c"},
		{1, "1"},
		{iv{}, "[-]"},
		{iv{1}, ",1"},
		{[]f{1.2, -3.5, 4}, "1.2 -3.5 4"},
		{[]f{3, 5, 4}, "3 5 4f"},
		{13.0, "13f"},
		{1 + 2i, "1i2"},
		{[]z{2i}, ",0i2"},
		{[]z{2i, 3.5 + 7i}, "0i2 3.5i7"},
		{l{1, 2, l{4, 5}}, "(1;2;(4;5))"},
		{d{l{5, 5.5}, iv{1, 2}}, "((5;5.5)!1 2)"},
	}
	for _, tc := range testCases {
		fmt.Printf("%v ?= %q\n", tc.x, tc.s)
		x := K(tc.x)
		y := fms(x)
		r := G(y).([]c)
		if reflect.DeepEqual(r, []byte(tc.s)) == false {
			t.Fatalf("expected: %q got %s (%q)\n", tc.s, string(r), string(r))
		}
		dec(y)
	}
}
func TestStr(t *testing.T) {
	ini()
	for _, x := range []s{"a", "b", "aa", "bb", "alpha", "betagammadelta"} {
		n := len(x)
		if n > 8 {
			n = 8
		}
		if r := G(K(x)); r != x[:n] {
			t.Fatalf("expected %s got %s\n", x, r)
		}
	}
	if u := sym(8 + K("abcdefgh")<<2); u != 0x6162636465666768 {
		t.Fatalf("%x\n", u)
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
	h := k(0)
	for i := k(0); i < k(len(m.k)); i += 4 {
		a, b, c, d := m.k[i+0], m.k[i+1], m.k[i+2], m.k[i+3]
		if a == 0 && b == 0 && c == 0 && d == 0 {
			continue
		}
		fmt.Printf("0x%04x %08x %08x %08x %08x", i, a, b, c, d)
		if i == h {
			tp := m.k[i] >> 28
			if tp == 0 {
				fmt.Printf("  %d", m.k[i])
				h += 1 << (m.k[i] - 2)
				nf := m.k[i+1]
				if nf > 0 && nf < 64 {
					fmt.Printf(" illegal fp")
				} else if nf > 0 && m.k[nf]>>28 != 0 {
					fmt.Printf(" fp is not free")
				}
			} else {
				atoms := "?cifzsgld?????"
				vects := "?CIFZSGLD?????"
				tp, n := typ(i)
				bt := bk(tp, n)
				if n == atom {
					fmt.Printf(" %c%d +%d", atoms[tp], bt, b)
				} else {
					fmt.Printf(" %c%d #%d +%d", vects[tp], bt, n, b)
				}
				h += 1 << (bt - 2)
			}
		}
		fmt.Println()
	}
}
func fpck(s s) { // check free pointers
	for i := 4; i < 32; i++ {
		nf := m.k[i]
		if nf > 0 && (nf < 64 || m.k[nf]>>28 != 0) {
			xxd()
			panic("fpck " + s + " bad pointer in free-list: @" + strconv.Itoa(int(i)))
		}
	}
	h := k(0)
	for i := k(0); i < k(len(m.k)); i += 4 {
		if i == h {
			tp := m.k[i] >> 28
			if tp == 0 {
				h += 1 << (m.k[i] - 2)
				nf := m.k[i+1]
				if nf > 0 && (nf < 64 || m.k[nf]>>28 != 0) {
					xxd()
					panic("fpck " + s + " illegal free-pointer")
				}
			} else {
				tp, n := typ(i)
				bt := bk(tp, n)
				h += 1 << (bt - 2)
			}
		}
	}
}
func pr(x k, a ...interface{}) {
	fmt.Printf(":%x ", x)
	r := fms(inc(x))
	_, n := typ(r)
	s := s(m.c[8+r<<2 : 8+n+r<<2])
	dec(r)
	fmt.Println(a, s)
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
				fmt.Printf("free block at %x with bt %d\n", a, t)
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
	kstr := func(dst k, s string) { // byte order independend
		u, n := uint64(0), len(s)
		if n > 8 {
			n = 8
		}
		for i := 0; i < n; i++ {
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
			return str(sym(8 + x<<2))
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
				r[i] = str(sym(8 + 8*k(i) + x<<2))
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
