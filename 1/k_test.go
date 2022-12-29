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

func TestK(t *testing.T) {
	k1()
	j := til(w(3))
	P := func(a, b string) { T(tok(ks(a)), b) }
	E := func(a, b string) { T(exe(ks(a+"\n"), 0), b) }
	O := func(x int32) { out(x); o(10) }
	_, _, _, _ = j, P, E, O

	//cond? 1+("4*x" "3*x" "x>4" "1+x")$x<0
	//      $[x<0;1+x;x>4;3*x;4*x]



	dotests("readme")
	E(`+\!4`, "(0 1 3 6)")
	E("1-(2)", "-1")
	E("3-(1+2)", "0")
	E("(1+2)-3", "0")
	E("-/2 3 4", "-5")
	E("+/2", "2")
	E("+/1 2 4", "7")
	E("--2", "2")
	E("-2", "-2")
	E("1+2", "3")
	E("1", "1")
	P(`abc"def`, "(97 98 99)")                  // skip open quotation
	P(`abc"def"`, "(97 98 99 ((100 101 102)))") //quotation
	P("123+4 5\n", "((123) 43 (4) 32 (5))")
	P("123\n", "((123))")
	T(ks("123"), "(49 50 51)")
	T(cal(w('\''), l2(l2(w('/'), w('+')), l2(j, j))), "(3 3)") // +/'(!3;!3)  ' ((/ +) (j j))
	T(amd(j, til(w(5)), rev(til(w(5)))), "(1 0 2)") // @[!3;|!5;|!5]   i:i mod xn
	T(amd(j, w(0), w(3)), "(3 1 2)")                // @[!3;0;3]
	T(val(w('y')), "0")                             // .y  not assigned
	T(asn(w('x'), j), "(0 1 2)")                    // x:!3
	T(val(w('x')), "(0 1 2)")                       // .x
	T(las(j), "2")                                  // _!3  last
	T(not(w(0)), "1")                               // ~!3
	T(not(j), "(1 0 0)")                            // ~!3
	T(flp(cat(cat(w(1), enl(ints(2, 3))), enl(j))), "((1 2 0) (1 3 1) (1 2 2))")
	T(flp(l2(j, j)), "((0 0) (1 1) (2 2))")                              // +(!3;!3)
	T(ech(w('#'), l2(j, j)), "(3 3)")                                    // #'(!3;!3)
	T(grp(ints(3, 2, 2, 1, 3, 2, 1)), "((3 2 1) ((0 4) (1 2 5) (3 6)))") // =3 2 3 1 3 2 1
	T(unq(cat(j, j)), "(0 1 2)")                                         // ?(!3),!3
	T(unq(j), "(0 1 2)")                                                 // ?!3
	T(gup(ints(1, 8, 1, 2, 5, 9)), "(0 2 3 4 1 5)")                      // <x
	T(gdn(ints(1, 8, 1, 2, 5, 9)), "(5 1 4 3 0 2)")                      // >x
	T(fnd(j, w('9')), "3")                                               // (!3)?9
	T(fnd(j, w(2)), "2")                                                 // (!3)?2
	T(fnd(j, til(w(5))), "(0 1 2 3 3)")                                  // (!3)?!5
	T(cts(cat(w(1), w(2)), j), "(((0) (1 2)) ((0 1) (2)))")              // 1 2^!3
	T(cut(w(1), j), "((0) (1 2))")                                       // 1^!3 (cut/split at index)
	T(rot(w(5)), "(5)")                                                  // %5  rotate
	T(rot(j), "(1 2 0)")                                                 // %!3  rotate
	T(mtv(j, til(w(3))), "1")                                            // (!3)~j
	T(mtc(3, 3), "1")                                                    // 1~1
	T(mtc(3, j), "0")                                                    // 1~j
	T(mtc(3, 5), "0")                                                    // 1~2
	T(mtv(j, wer(w(3))), "0")                                            // (!3)~&3
	T(wer(j), "(1 2 2)")                                                 // &!3
	T(wer(w(3)), "(0 0 0)")                                              // &3
	T(ovr(77, til(w(3))), "0")                                           // &/!3
	T(ovr(w('+'), til(w(0))), "0")                                       // +/!3
	T(ovr(w('+'), j), "3")                                               // +/!3
	T(rev(j), "(2 1 0)")                                                 // |x is x@xn-1+!xn:#x
	T(tak(w(-2), w(3)), "(3 3)")                                         // -2#3
	T(drp(w(-3), j), "()")                                               // underdrop
	T(drp(w(-5), j), "()")                                               // underdrop
	T(drp(w(-1), j), "(0 1)")                                            // taildrop
	T(drp(w(5), j), "()")                                                // overdrop
	T(drp(w(1), j), "(1 2)")                                             // drop
	T(tak(w(4), w(3)), "(3 3 3 3)")                                      // 4#3 scalar take
	T(tak(w(-5), j), "(1 2 0 1 2)")                                      // undertake
	T(tak(w(-2), j), "(1 2)")                                            // take tail
	T(tak(w(5), j), "(0 1 2 0 1)")                                       // overtake
	T(cal(67, l2(til(w(4)), w(3))), "(0 1 2 0)")                         // (!4)!3
	T(cal(w('!'), l2(til(w(4)), w(3))), "(0 1 2 0)")                     // (!4)!3  mod
	T(cal(w('*'), l2(til(w(4)), til(w(2)))), "(0 1 0 3)")                // 0 1 2 3*0 1
	T(cal(w('*'), l2(j, j)), "(0 1 4)")                                  // 0 1 2*0 1 2
	T(cal(w('*'), l2(j, w(-2))), "(0 -2 -4)")                            // 0 1 2*-2
	T(cal(w('*'), l2(w(-2), j)), "(0 -2 -4)")                            // -2*0 1 2
	T(cal(w('-'), til(w(2))), "-1")                                      // 0-1
	T(cal(w('-'), enl(j)), "(0 -1 -2)")                                  // -!3
	T(cal(w('-'), w(3)), "-3")                                           // -3
	T(cal(w(46), l2(j, w(1))), "1")                                      // .[j;0] => dyadic
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

	O(w(123))
	O(ks("abc"))
	O(j)

	println("top/tot", top/1024, tot/1024, "k")
}
func T(a int32, b string) {
	s := tostring(a)
	if s != b {
		panic("got:" + s + "\nnot:" + b)
	}
	println(s)
}
func ks(s string) int32 {
	r := rm(int32(len(s)))
	for _, c := range s {
		c1(r, w(int32(c)))
	}
	return r
}
func tostring(x int32) string {
	if x == 0 {
		return "null"
	} else if x&1 != 0 {
		return strconv.Itoa(int(x >> 1))
	} else {
		xn := n(x)
		if xn < 0 || xn > 30 {
			panic("xn")
		}
		u := make([]string, xn)
		for i := int32(0); i < xn; i++ {
			u[i] = tostring(I32(x + 4*i))
		}
		return "(" + strings.Join(u, " ") + ")"
	}
}
func ints(x ...int) int32 {
	n := len(x)
	r := rm(int32(n))
	for _, i := range x {
		c1(r, w(int32(i)))
	}
	return r
}
func dotests(file string) {
	ts := readtests(file)
	for _, t := range ts {
		print(t[0])
		x := ks(t[0])
		println(tostring(x))
		r := exe(x, 0)
		s := tostring(r)
		if t[1] != s {
			panic("got " + s)
		}
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
	r[0] += "\n"
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
