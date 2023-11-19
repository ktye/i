package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	. "github.com/ktye/wg/module"
)

var save []byte

func newtest() {
	Stdout = os.Stdout
	rand_ = 1592653589
	if save == nil {
		Bytes = make([]byte, 64*1024)
		loc, xyz, pp, pe, sp, srcp, rand_ = 0, 0, 0, 0, 0, 0, 0
		Data(132, "\x00\x01@\x01\x01\x01\x01\t\x10`\x01\x01\x01\x01\x01\t\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\x01 \x01\x01\x01\x01\x01BBBBBBBBBBBBBBBBBBBBBBBBBB\x10\t`\x01\x01\x00\xc2\xc2\xc2\xc2\xc2\xc2BBBBBBBBBBBBBBBBBBBB\x10\x01`\x01") // k_test.go: TestClass
		Data(227, ":+-*%&|<>=~!,^#_$?@.':/:\\:vbcisfzldtmdplx00BCISFZLDT0")
		kinit()
		save = make([]byte, len(Bytes))
		copy(save, Bytes)
	} else {
		Bytes = make([]byte, len(save))
		copy(Bytes, save)
		pp, pe, sp = 0, 0, 2048
	}
}

func TestInit(t *testing.T) {
	newtest()
	reset()
}
func TestMk(t *testing.T) {
	newtest()
	x := Ki(3)
	y := Til(Ki(4))
	z := Add(x, y)
	if tp(z) != It || nn(z) != 4 {
		t.Fatal()
	}
	dx(z)
	reset()
}
func TestKst(t *testing.T) {
	newtest()
	s := sk(Ki(3))
	if s != "3" {
		t.Fatal(s)
	}
	reset()
}
func TestKT(t *testing.T) {
	newtest()
	b, e := os.ReadFile("k.t")
	if e != nil {
		t.Fatal(e)
	}
	v := bytes.Split(b, []byte{10})
	for i := range v {
		if len(v[i]) > 1 {
			test(mkchars(append(v[i], 10)))
			reset()
		}
	}
	reset()
}

func mkchars(b []byte) (r K) {
	r = mk(Ct, int32(len(b)))
	copy(Bytes[int32(r):], b)
	return r
}
func sk(x K) string {
	x = Kst(x)
	n := nn(x)
	r := string(Bytes[int32(x) : int32(x)+n])
	dx(x)
	return r
}
func sK(x K) string {
	xp := int32(x)
	switch tp(x) {
	case 0:
		if x == 0 {
			return ""
		}
		s := []byte("0:+-*%!&|<>=~,^#_$?@.'/\\")
		var r string
		itoa := func(x int32) string { return strconv.Itoa(int(x)) }
		switch {
		case xp < 64:
			if xp < 23 {
				r = string(s[xp])
			} else {
				r = "`" + itoa(xp)
			}
			return r
		case xp < 128:
			if xp-64 < 23 {
				r = string(s[xp-64])
			} else {
				r = "`" + itoa(xp)
			}
			return r
		case xp == 211:
			return "@"
		case xp == 212:
			return "."
		case xp >= 448 && xp-448 < 23:
			return string(s[xp-448])
		default:
			return "`" + itoa(xp)
		}
	case ct:
		return strconv.Quote(string([]byte{byte(xp)}))
	case it:
		return strconv.Itoa(int(xp))
	case st:
		n := nn(K(I64(0)))
		if 8*n <= xp {
			panic("illegal symbol")
		}
		x = cs(x)
		dx(x)
		xp = int32(x)
		if nn(x) == 0 {
			return "`"
		}
		return "`" + string(Bytes[xp:xp+nn(x)])
	case 5:
		panic("float")
		return "?"
	case 6:
		panic("float")
		return "?"
	case cf:
		xn := nn(x)
		xp = int32(x) + 8*xn
		s := ""
		for i := int32(0); i < xn; i++ {
			xp -= 8
			s += sK(K(I64(xp)))
		}
		return s
	case df:
		a := []string{"'", "':", "/", "/:", "\\", "\\:"}
		r := sK(K(I64(xp)))
		p := I64(xp + 8)
		return r + a[int(p)]
	case pf:
		f := K(I64(xp))
		l := K(I64(xp + 8))
		i := K(I64(xp + 16))
		// if tp(f) == 0 && nn(i) == 1 && I32(int32(i)) == 1 {
		if nn(i) == 1 && I32(int32(i)) == 1 {
			return sK(K(I64(int32(l)))) + sK(f) // 1+
		}
		return "<prj>"
	case lf:
		x = K(I64(xp + 16))
		xp = int32(x)
		return string(Bytes[xp : xp+nn(x)])
		/*
			case Bt:
				r := bytes.Repeat([]byte{'0'}, int(nn(x)))
				for i := range r {
					if I8(xp+int32(i)) != 0 {
						r[i] = '1'
					}
				}
				return comma(1 == nn(x)) + string(r) + "b"
		*/
	case Ct:
		return comma(1 == nn(x)) + strconv.Quote(string(Bytes[xp:xp+nn(x)]))
	case It:
		if nn(x) == 0 {
			return "!0"
		}
		r := make([]string, nn(x))
		for i := range r {
			r[i] = strconv.Itoa(int(I32(xp + 4*int32(i))))
		}
		return comma(1 == nn(x)) + strings.Join(r, " ")
	case St:
		r := make([]string, nn(x))
		for i := range r {
			r[i] = sK(K(I32(xp)) | K(st)<<59)
			xp += 4
		}
		if nn(x) == 0 {
			return "0#`"
		}
		return comma(1 == nn(x)) + strings.Join(r, "")
	case 21:
		panic("float")
		return "?"
	case 22:
		panic("float")
		return "?"
	case Lt:
		r := make([]string, nn(x))
		for i := range r {
			r[i] = sK(K(I64(xp)))
			xp += 8
		}
		if len(r) == 1 {
			return "," + r[0]
		} else {
			return "(" + strings.Join(r, ";") + ")"
		}
	case Dt:
		return sK(K(I64(xp))) + "!" + sK(K(I64(xp+8)))
	case Tt:
		return "+" + sK(K(I64(xp))) + "!" + sK(K(I64(xp+8)))
	default:
		fmt.Println("type ", tp(x))
		panic("type")
	}
}
func comma(x bool) string {
	if x {
		return ","
	} else {
		return ""
	}
}

func reset() {
	if sp != 2048 {
		println(sp)
		panic("sp")
	}
	dx(src())
	dx(xyz)
	dx(K(I64(0)))
	dx(K(I64(8)))
	//additional 1024 for cpu stack (see kinit)
	if (uint32(1)<<uint32(I32(128)))-(1024+4096+mcount()) != 0 {
		panic("memcount")
	}
	for i := int32(5); i < 31; i++ {
		SetI32(4*i, 0)
	}
	kinit()
}
func mcount() uint32 {
	r := uint32(0)
	for i := int32(5); i < 31; i++ {
		n := fcount(4 * i)
		r += uint32(n) * (1 << uint32(i))
	}
	return r
}
func fcount(x int32) int32 {
	r := int32(0)
	for {
		if I32(x) == 0 {
			break
		}
		r++
		x = I32(x)
	}
	return r
}

func test(x K) {
	if tp(x) != Ct {
		trap() //type
	}
	l := ndrop(-1, split(Kc(10), rx(x)))
	n := nn(l)
	dx(l)
	for i := int32(0); i < n; i++ {
		testi(rx(x), i)
	}
	dx(x)
}
func testi(l K, i int32) {
	x := split(Ku(12064), ati(split(Kc(10), l), i))
	if nn(x) != 2 {
		trap() //length
	}
	y := x1(x)
	x = r0(x)
	dx(Out(ucat(ucat(rx(x), Ku(12064)), rx(y))))
	x = Kst(val(x))
	if match(x, y) == 0 {
		x = Out(x)
		trap() //test fails
	}
	dxy(x, y)
}
