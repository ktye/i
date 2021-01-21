package main

import (
	"testing"
)

func eval(s string) uint32   { return val(kC([]byte(s))) }
func kstval(s string) string { return string(CK(kst(eval(s)))) }

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
func TestRand(t *testing.T) {
	kinit()
	for _, in := range []string{"#randf 3", "#randi 3", "#randz 3", "#randn 3", "#shuffle `beta`alpha`gamma"} {
		s := kstval(in)
		if s != "3" {
			t.Fatal()
		}
	}
	bleak()
}
func TestSolve(t *testing.T) {
	kinit()
	if s := kstval(`A:3 5#randf 15;,/(@;#)@\:solve[A;randf 5]`); s != "3 3" {
		t.Fatal(s)
	}
	if s := kstval(`A:5 6#randz 30;,/(@;#)@\:solve[A;randz 6]`); s != "4 5" {
		t.Fatal(s)
	}
	if s := kstval(`A:3 5#randf 15;,/(@;#)@\:solve[A;10 5#randf 50]`); s != "6 10" {
		t.Fatal(s)
	}
	if s := kstval(`A:5 6#randz 30;,/(@;#)@\:solve[A;11 6#randz 66]`); s != "6 11" {
		t.Fatal(s)
	}
	if s := kstval(`A:3 5#randf 15;x:randf 3;b:,/mul[&A;x];1.e-7>|/+x-solve[A;b]`); s != "1" {
		t.Fatal(s)
	}
	if s := kstval(`A:3 5#randz 15;x:randz 3;b:,/mul[&A;x];1.e-7>|/+x-solve[A;b]`); s != "1" {
		t.Fatal(s)
	}
	if s := kstval(`cond diag 1a10 3a30 2a40`); s != "3.0" {
		t.Fatal(s)
	}
	bleak()
}
func TestDiag(t *testing.T) {
	kinit()
	if s := kstval(`,/(@;#)@\:diag 5`); s != "6 5" {
		t.Fatal(s)
	}
	if s := kstval(`,/(@;#)@\:diag 4 4#randf 24`); s != "3 4" {
		t.Fatal(s)
	}
	if s := kstval(`,/(@;#)@\:diag 5 5#randz 25`); s != "4 5" {
		t.Fatal(s)
	}
	bleak()
}
func TestMul(t *testing.T) {
	kinit()
	if s := kstval(`t:{,/(@*x;#x;#*x)}`); s != "" {
		t.Fatal(s)
	}
	if s := kstval(`A:3 5#randf 15;B:4 5#randf 20;t mul[A;B]`); s != "3 4 3" {
		t.Fatal(s)
	}
	if s := kstval(`A:3 5#randz 15;B:4 5#randz 20;t mul[A;B]`); s != "4 4 3" {
		t.Fatal(s)
	}
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
