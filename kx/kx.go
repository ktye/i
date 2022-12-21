package main

import (
	"bufio"
	_ "embed"
	"math"
	"math/cmplx"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//go:embed readme
var help []byte

var out *bufio.Writer

func main() {
	a := os.Args[1:]
	if len(a) == 0 || a[0] == "-h" {
		os.Stderr.Write(help)
		os.Exit(1)
	}

	kinit()
	if strings.HasSuffix(a[0], ".k") {
		x := KC(a[0])
		dofile(x, readfile(rx(x)))
		a = a[1:]
	}

	var e []string
	a, e = esplit(a)
	a, e = replace(a), replace(e)
	//fmt.Println("a:", a)
	//fmt.Println("e:", e)

	A := make([]uint64, len(a))
	for i, s := range a {
		A[i] = Prs(KC(s))
	}
	R := make([]string, len(a))

	I := sc(Ku(105)) // i
	X := sc(Ku(120)) // x
	i := int32(-1)
	s := bufio.NewScanner(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	for s.Scan() {
		i++
		dx(Asn(I, Ki(i)))                        // i
		dx(Asn(X, KL(strings.Fields(s.Text())))) // x

		n := 0
		for i := range A {
			R[i] = eval(A[i])
			n += len(R[i])
		}
		if n == 0 {
			continue
		}
		writes(R)
	}
	if len(e) > 0 {
		dx(Asn(I, Ki(1+i))) // i:number of rows
		R = make([]string, len(e))
		for i, s := range e {
			R[i] = eval(Prs(KC(s)))
		}
		writes(R)
	}
	out.Flush()
}
func esplit(a []string) (r, e []string) {
	for i := range a {
		if a[i] == "-e" {
			return a[:i], a[1+i:]
		}
	}
	return a, nil
}
func replace(v []string) []string {
	t := [][2]string{
		{`([0-9]+)([if])`, "(`$2$$x $1)"},    // 123i 123f
		{`([0-9]+)s`, "(`$$x $1)"},           // 123s
		{`([0-9]+)c`, "(x $1)"},              // 3c
		{`x([if])`, "(`$1$$x)"},              // xi xf
		{`xs`, "(`$$x)"},                     // xs
		{`^/([^/]*)/(.*)`, "$$[$1;$2;+]"},    // /cond/prog
		{`^\\([^\\]*)\\(.*)`, "$$[$1;$2;+]"}, // \cond\prog
		{`°`, "'"},                           // each
		// todo cond
	}
	for i, s := range v {
		for _, x := range t {
			re := regexp.MustCompile(x[0])
			s = re.ReplaceAllString(s, x[1])
		}
		v[i] = s
	}
	return v
}
func KC(s string) uint64 {
	r := mk(18, int32(len(s)))
	copy(Bytes[int32(r):], []byte(s))
	return r
}
func CK(x uint64) string { dx(x); return string(Bytes[int32(x) : int32(x)+nn(x)]) }
func KL(v []string) uint64 {
	r := mk(23, int32(len(v)))
	p := int32(r)
	for i := range v {
		SetI64(p+8*int32(i), int64(KC(v[i])))
	}
	return r
}
func eval(x uint64) string {
	return strs(exec(rx(x)))
}
func writes(x []string) { out.WriteString(strings.Join(x, " ") + "\n") }
func strs(x uint64) (r string) {
	p := int32(x)
	t := tp(x)
	switch t {
	case 2:
		r = string(p)
	case 3:
		r = strconv.Itoa(int(p))
	case 4:
		r = CK(cs(rx(x)))
	case 5:
		r = ftoa(F64(p))
	case 6:
		r = absang(complex(F64(p), F64(p+8)))
	case 18:
		r = CK(rx(x)) // C
	case 19, 20, 21, 22, 23:
		r = each(x) // IFSZL
	}
	dx(x)
	return r
}
func each(x uint64) string {
	n := nn(x)
	r := make([]string, n)
	for i := int32(0); i < n; i++ {
		r[i] = strs(ati(rx(x), i))
	}
	return strings.Join(r, " ")
}
func ftoa(f float64) string { return strconv.FormatFloat(f, 'g', 6, 64) }
func absang(z complex128) string {
	ang := 180.0 / math.Pi * cmplx.Phase(z)
	if ang < 0 {
		ang += 360.0
	}
	return ftoa(cmplx.Abs(z)) + " " + ftoa(ang)
}
