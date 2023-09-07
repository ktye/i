package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	. "github.com/ktye/wg/module"
)

var sourcefiles []string
var filestart []int32
var filelines [][]int

type Q struct {
	p int32
	s string
}

var tracemap map[Q]bool

func trace(s string) {
	if tracemap == nil {
		tracemap = make(map[Q]bool)
	}
	tracemap[Q{srcp, s}] = true
}
func tracestart() {
	x := src()
	n := nn(x)
	b := Bytes[int32(x) : int32(x)+n]
	//fmt.Println("z.k!")
	//fmt.Println(string(b))
	sourcefiles = append(sourcefiles, "z.k")
	filelines = append(filelines, getlinesb(b))
	filestart = append(filestart, 0)
}
func tracefile(x K) {
	p := int32(x)
	n := nn(x)
	name := string(Bytes[p : p+n])

	if strings.HasSuffix(name, ".k") {
		sourcefiles = append(sourcefiles, name)
		filelines = append(filelines, getlines(name))
		filestart = append(filestart, nn(src()))
	}
}
func tracetest(x K) {
}
func traceend() {
	fmt.Fprintln(os.Stderr, "sourcefiles", sourcefiles)
	fmt.Fprintln(os.Stderr, "filestart", filestart)
	for k := range tracemap {
		fmt.Fprintln(os.Stderr, k.p, k.s, fileposition(-1+k.p))
	}
}

func fileposition(x int32) string {
	j := -1
	f := "z.k"
	b := int32(0)
	for i, a := range filestart {
		if x < a {
			x -= b
			return f + linecol(j, int(x))
		}
		f = sourcefiles[i]
		b = a
		j++
	}
	x -= b
	return f + linecol(j, int(x))
}
func linecol(j, x int) string {
	if j < 0 {
		return "+" + strconv.Itoa(x)
	}
	if j >= len(filelines) {
		return "+?"
	}
	for i, n := range filelines[j] {
		if x <= n {
			return ":" + strconv.Itoa(1+i) + ":" + strconv.Itoa(x)
		}
		x -= 1 + n
	}
	return "+$"
}
func getlines(name string) (r []int) {
	b, e := os.ReadFile(name)
	if e != nil {
		panic(e)
	}
	return getlinesb(b)
}
func getlinesb(b []byte) (r []int) {
	v := bytes.Split(b, []byte("\n"))
	for _, s := range v {
		r = append(r, len(s))
	}
	return r
}
