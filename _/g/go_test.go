package main

import (
	"testing"
)

func eval(s string) uint32   { return val(kC([]byte(s))) }
func kstval(s string) string { return string(CK(kst(eval(s)))) }
func TestT(t *testing.T) {
	kinit()
	dx(eval("t:`a`b`c!(1+!10;2+!10;3+!10)"))
	if s := kstval("#*where[t;5<t`c]"); s != "7" {
		panic(s)
	}
	dx(eval("a:12"))
	if s := kstval("#*t 'w{a>4}"); s != "6" {
		panic(s + "?6")
	}
	if s := kstval("a"); s != "12" {
		panic(s)
	}
	bleak()
}
func TestG(t *testing.T) {
	kinit()
	r := kstval("a:{x+y};a[3;4]")
	if r != "7" {
		t.Fatal()
	}
	bleak()
}
func TestPlot(t *testing.T) {
	kinit()
	dx(eval("plot 1 2 3"))
	bleak()
}
func TestGo(t *testing.T) {
	testR := func(x bool) {
		if x == false {
			t.Fatal()
		}
	}
	k := func(s string) uint32 { return eval(s) }

	var r uint32
	kinit()

	var f float64
	r = k("1.2")
	testG(r, &f)
	testR(f == 1.2)
	dx(r)

	var I []int
	r = k("4 3 2")
	testG(r, &I)
	testR(len(I) == 3 && I[0] == 4)
	dx(r)

	var F []float64
	r = k("1 2 3.0")
	testG(r, &F)
	testR(len(F) == 3 && F[2] == 3.0)
	dx(r)

	var Z []complex128
	r = k("1a90 1a180")
	testG(r, &Z)
	testR(len(Z) == 2 && real(Z[1]) == -1)
	dx(r)

	var i int
	r = k("13")
	testG(r, &i)
	testR(i == 13)
	dx(r)

	var y string
	r = k("`alpha")
	testG(r, &y)
	testR(y == "alpha")
	dx(r)

	var p point
	r = k("`X`Y!(3;4.5)")
	testG(r, &p)
	testR(p == point{3, 4.5})
	dx(r)

	var P []point
	r = k("(`X`Y!(3;4.5);`X`Y!(1;2.0)")
	testG(r, &P)
	testR(len(P) == 2 && P[1].Y == 2.0)
	dx(r)
}
func testG(x uint32, r interface{}) {
	m := memstore()
	G(x, r)
	memcompare(m, "testG")
}

type point struct {
	X int
	Y float64
}
