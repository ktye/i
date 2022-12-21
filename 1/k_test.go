package main

import (
	"bytes"
	. "github.com/ktye/wg/module"
	"os"
	"strconv"
	"strings"
	"testing"
)

// var tot, top int32

func TestVerbs(t *testing.T) {
	k1()
	j := til(w(3))
	_ = j

	T(fnd(j, w(2)), "2")                                    // (!3)?2
	T(fnd(j, til(w(5))), "(0 1 2 3 3)")                     // (!3)?!5
	T(cts(cat(w(1), w(2)), j), "(((0) (1 2)) ((0 1) (2)))") // 1 2^!3
	T(cut(w(1), j), "((0) (1 2))")                          // 1^!3 (cut/split at index)
	T(rot(w(5)), "(5)")                                     // %5  rotate
	T(rot(j), "(1 2 0)")                                    // %!3  rotate
	T(mtv(j, til(w(3))), "1")                               // (!3)~j
	T(mtc(3, 3), "1")                                       // 1~1
	T(mtc(3, j), "0")                                       // 1~j
	T(mtc(3, 5), "0")                                       // 1~2
	T(mtv(j, wer(w(3))), "0")                               // (!3)~&3
	T(wer(j), "(1 2 2)")                                    // &!3
	T(wer(w(3)), "(0 0 0)")                                 // &3
	T(ovr(77, til(w(3))), "0")                              // &/!3
	T(ovr(w('+'), til(w(0))), "0")                          // +/!3
	T(ovr(w('+'), j), "3")                                  // +/!3
	T(rev(j), "(2 1 0)")                                    // |x is x@xn-1+!xn:#x
	T(tak(w(-2), w(3)), "(3 3)")                            // -2#3
	T(drp(w(-3), j), "()")                                  // underdrop
	T(drp(w(-5), j), "()")                                  // underdrop
	T(drp(w(-1), j), "(0 1)")                               // taildrop
	T(drp(w(5), j), "()")                                   // overdrop
	T(drp(w(1), j), "(1 2)")                                // drop
	T(tak(w(4), w(3)), "(3 3 3 3)")                         // 4#3 scalar take
	T(tak(w(-5), j), "(1 2 0 1 2)")                         // undertake
	T(tak(w(-2), j), "(1 2)")                               // take tail
	T(tak(w(5), j), "(0 1 2 0 1)")                          // overtake
	T(cal(67, l2(til(w(4)), w(3))), "(0 1 2 0)")            // (!4)!3
	T(cal(w('!'), l2(til(w(4)), w(3))), "(0 1 2 0)")        // (!4)!3  mod
	T(cal(w('*'), l2(til(w(4)), til(w(2)))), "(0 1 0 3)")   // 0 1 2 3*0 1
	T(cal(w('*'), l2(j, j)), "(0 1 4)")                     // 0 1 2*0 1 2
	T(cal(w('*'), l2(j, w(-2))), "(0 -2 -4)")               // 0 1 2*-2
	T(cal(w('*'), l2(w(-2), j)), "(0 -2 -4)")               // -2*0 1 2
	T(cal(w('-'), til(w(2))), "-1")                         // 0-1
	T(cal(w('-'), enl(j)), "(0 -1 -2)")                     // -!3
	T(cal(w('-'), w(3)), "-3")                              // -3
	T(cal(w(46), l2(j, w(1))), "1")                         // .[j;0] => dyadic
	T(atx(j, til(w(0))), "()")
	T(atx(j, cat(enl(j), enl(j))), "((0 1 2) (0 1 2))")
	T(atx(j, j), "(0 1 2)")
	T(atx(j, w(2)), "2")
	T(atx(w(3), w(9)), "3")
	T(neg(w(5)), "-5")     // -5
	T(neg(j), "(0 -1 -2)") // -!3
	T(fst(enl(j)), "(0 1 2)")
	T(fst(j), "0")
	T(max(w(2), w(2)), "2")
	T(cat(w(1), w(2)), "(1 2)")
	T(el(w(1)), "(1)")
	T(el(enl(w(1))), "(1)")
	T(enl(w(1)), "(1)")
	T(cnt((w(3))), "1")
	T(cnt(til(w(5))), "5")
	T(til(w(0)), "()")
	T(til(w(-3)), "(-3 -2 -1)")
	T(til(w(5)), "(0 1 2 3 4)")
	T(w(-2), "-2")
	//readtests("readme")
}
func T(a int32, b string) {
	s := tostring(a)
	if s != b {
		panic("got:" + s + "\nnot:" + b)
	}
	println(s)
}
func tostring(x int32) string {
	if x&1 != 0 {
		return strconv.Itoa(int(x >> 1))
	} else {
		xn := n(x)
		if xn < 0 || xn > 20 {
			panic("xn")
		}
		u := make([]string, xn)
		for i := int32(0); i < xn; i++ {
			u[i] = tostring(I32(x + 4*i))
		}
		return "(" + strings.Join(u, " ") + ")"
	}
}
func readtests(file string) (r [][2]string) {
	b, e := os.ReadFile(file)
	fatal(e)
	v := bytes.Split(b, []byte("\n"))
	for i := range v {
		s := string(v[i])
		if len(s) > 0 && s[0] == ' ' {
			r = append(r, readcase(s[1:]))
		}
	}
	return r
}
func readcase(s string) (r [2]string) {
	var o bool
	r[0], r[1], o = strings.Cut(s, " /")
	if !o {
		panic("testcase: " + s)
	}
	return
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
