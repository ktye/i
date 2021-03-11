package main

import (
	"math/big"
	"time"
)

type A interface {
	Ref() A
	Unref()
	Zero() A
	Empty() A
	Index([]int) A
	Element(int) A
	Len() int
}

// C, Byte             "a" "abc" 0x12ef
// S, Symbol           `a`b `"alpha"
// I, Int              1 2 -123
// F, Float            1.0 1.1 -1.23e-002
// Z, Complex          1.2a310 2a
// B, Big              1234i
// Tv, Time            2021.03.11T12:15:17
// Dv, Duration        2.5h 2ms 3us
// L list              (1;2 3;4.5)
// D dict              `a`b!(1;2 3)
// T table             +`a`b!(1 2;3 4)
// Verb                +
// Adverb              /
// Lambda              {x+y}
// Native              func(x, y A) A
// Projection          +[;3]
// Composition         +-

type V byte
type L struct {
	Rc int
	A  []A
}
type I struct {
	Rc int
	A  []int
}
type Int int
type F struct {
	Rc int
	A  []float64
}
type Float float64
type Big struct {
	Rc int
	*big.Int
}
type Duration time.Duration

func (l *L) Ref() A { l.Rc++; return l }
func (l *L) Unref() {
	l.Rc--
	if l.Rc < 0 {
		for i := range l.A {
			l.A[i].Unref()
		}
	}
}
func (l *L) Zero() A  { return nil /*return &C{' '}*/ }
func (l *L) Empty() A { return &L{} }
func (l *L) Len() int { return len(l.A) }
func (l *L) Index(x []int) A {
	a := make([]A, len(x))
	for i, j := range x {
		a[i] = l.A[j].Ref()
	}
	return &L{A: a}
}
func (l *L) Element(i int) A { return l.A[i].Ref() }

func (i *I) Ref() A   { i.Rc++; return i }
func (i *I) Unref()   { i.Rc-- }
func (i *I) Zero() A  { return Int(0) }
func (i *I) Empty() A { return &I{} }
func (i *I) Len() int { return len(i.A) }
func (i *I) Index(x []int) A {
	a := make([]int, len(x))
	for k, j := range x {
		a[k] = i.A[j]
	}
	return &I{A: a}
}
func (i *I) Element(x int) A { return Int(i.A[x]) }

func (i Int) Ref() A   { return i }
func (i Int) Unref()   {}
func (i Int) Zero() A  { return Int(0) }
func (i Int) Empty() A { return &I{} }
func (i Int) Len() int { return 1 }
func (i Int) Index(x []int) A {
	a := make([]int, len(x))
	for j := range x {
		a[j] = int(i)
	}
	return &I{A: a}
}
func (i Int) Element(x int) A { return i }

func (f *F) Ref() A   { f.Rc++; return f }
func (f *F) Unref()   { f.Rc-- }
func (f *F) Zero() A  { return Float(0) }
func (f *F) Empty() A { return &F{} }
func (f *F) Len() int { return len(f.A) }
func (f *F) Index(x []int) A {
	a := make([]float64, len(x))
	for k, j := range x {
		a[k] = f.A[j]
	}
	return &F{A: a}
}
func (f *F) Element(x int) A { return Float(f.A[x]) }

func (f Float) Ref() A   { return f }
func (f Float) Unref()   {}
func (f Float) Zero() A  { return Float(0) }
func (f Float) Empty() A { return &F{} }
func (f Float) Len() int { return 1 }
func (f Float) Index(x []int) A {
	a := make([]float64, len(x))
	for j := range x {
		a[j] = float64(f)
	}
	return &F{A: a}
}
func (f Float) Element(x int) A { return f }

func (i *Big) Ref() A   { i.Rc++; return i }
func (i *Big) Unref()   { i.Rc-- }
func (i *Big) Zero() A  { var u big.Int; return &Big{Int: &u} }
func (i *Big) Empty() A { return &L{} }
func (i *Big) Len() int { return 1 }
func (i *Big) Index(x []int) A {
	a := make([]A, len(x)) // todo uni
	for j := range x {
		a[j] = i.Ref()
	}
	return &L{A: a}
}
func (i *Big) Element(x int) A { return i.Ref() }
func (i *Big) String() string  { return i.Int.String() + "i" }

func (d Duration) Ref() A   { return d }
func (d Duration) Unref()   {}
func (d Duration) Zero() A  { return Duration(0) }
func (d Duration) Empty() A { return &L{} }
func (d Duration) Index(x []int) A {
	a := make([]A, len(x)) // todo uni
	for j := range x {
		a[j] = d
	}
	return &L{A: a}
}
func (d Duration) Element(int) A { return d }
func (d Duration) Len() int      { return 1 }

func (v V) Ref() A        { return v }
func (v V) Unref()        {}
func (v V) Zero() A       { return V(0) }
func (v V) Empty() A      { return &L{} }
func (v V) Index([]int) A { return nil }
func (v V) Element(int) A { return nil }
func (v V) Len() int      { return 1 }
