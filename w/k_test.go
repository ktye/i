package kw

import "testing"

func TestRst(t *testing.T) {
	rst()
	e := map[int]k{8: 256, 9: 512, 10: 1024, 11: 2048, 12: 4096, 13: 8192, 14: 16384, 15: 32768}
	for i := 4; i < 30; i++ {
		if get(k(4*i)) != e[i] {
			t.Fatal()
		}
	}
	pfl()
	mk(1, 3)
	pfl()
	xxd()
}
func pfl() {
	for i := 0; i < 30; i++ {
		println(i, get(k(4*i)))
	}
}
func xxd() { // memory dump
	t := [16]b{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 97, 98, 99, 100, 101, 102}
	s2 := func(x b) (b, b) { return t[x>>4], t[x&0xF] }
	l := make([]b, 49)
	for i := 0; i < len(l); i++ {
		l[i] = 32
	}
	n := 0
	u := k(0)
	e := true
	h := 0
	s := make([]b, 32)
	s[3] = '#'
	tp, tn, rc := b(0), 0, k(0)
	for i := 0; i < len(m); i += 2 {
		if i == h {
			tp, tn = typ(k(i))
			c := [8]b{'x', 'i', 'f', 'z', 's', 'g', 'd', 'l'}
			s[0] = c[tp]
			if tn >= 0 {
				s[0] -= 32
			}
			rc = get(k(i + 4))
			bt := m[i] & 63
			s[1], s[2] = s2(bt)
			h += 1 << bt
		}
		if n == 0 {
			l[0], l[1] = s2(b(u >> 24))
			l[2], l[3] = s2(b(u >> 16))
			l[4], l[5] = s2(b(u >> 8))
			l[6], l[7] = s2(b(u))
			u += 16
			n = 8
		}
		l[n+1], l[n+2] = s2(m[i])
		l[n+3], l[n+4] = s2(m[i+1])
		if m[i] != 0 || m[i+1] != 0 {
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
					s[0], tn = 0, -1
				}
				println()
			}
			e = true
		}
	}
}
