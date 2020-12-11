package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func csv2(x, y uint32) uint32 {
	if tp(x) != 1 {
		panic("csv x must be char(format)")
	}
	if tp(y) != 1 {
		panic("csv y must be char(data)")
	}
	b := C[y+8 : y+8+nn(y)]
	defer dx(y)
	return parseCsv(sk(x), b)
}

// i: int
// f: float
// z+: complex from two columns abs angle
// a: complex from 1 column 1a30
// c: char array
// s: symbol
// *: repeat last
func parseCsv(format string, b []byte) uint32 {
	if len(b) > 2 && b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
		b = b[3:] // bom
	}
	skip := 0
	if len(format) > 0 && format[0] >= '0' && format[0] <= '9' {
		skip = int(format[0])
		format = format[1:]
	}
	for i := 0; i < skip; i++ {
		if idx := bytes.Index(b, []byte{'\n'}); idx >= 0 {
			b = b[idx+1:]
		}
	}
	r := csv.NewReader(bytes.NewReader(b))
	if len(format) > 0 && strings.IndexByte("-fzacs", format[0]) == -1 {
		r.Comma = rune(format[0])
		format = format[1:]
	}
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	cols := 0
	if len(records) > 0 {
		cols = len(records[0])
	}
	p := csvParsers(format, cols)
	for _, row := range records {
		for k := range p {
			p[k].parse(row)
		}
	}
	l := mk(6, uint32(len(p)))
	for i := uint32(0); i < uint32(len(p)); i++ {
		I[2+i+l>>2] = p[i].k()
	}
	return l
}
func csvParsers(format string, cols int) (r []csvparser) {
	if len(format) > 1 && format[len(format)-1] == '*' {
		format = format[:len(format)-1]
		r := string(format[len(format)-1])
		if r == "+" {
			r = "r+"
		}
		format += strings.Repeat(r, cols-len(format))
	}
	for i, c := range format {
		switch c {
		case '-':
		case 'c':
			r = append(r, &csvC{n: i})
		case 'i':
			r = append(r, &csvI{n: i})
		case 'f':
			r = append(r, &csvF{n: i})
		case 'z':
			if i+1 == len(format) || format[i+1] != '+' {
				panic("csv:format z must be followed by +")
			}
			r = append(r, &csvZ{n: i})
		case 'a':
			r = append(r, &csvA{n: i})
		case 's':
			r = append(r, &csvS{n: i})
		case '+':
			if i == 0 || format[i-1] != 'z' {
				panic("csv:format + must follow z")
			}
		default:
			panic(fmt.Sprintf("csv unknown format: %c", c))
		}
	}
	return r
}

type csvparser interface {
	parse([]string)
	k() uint32
}
type csvI struct {
	i []uint32
	n int
}
type csvF struct {
	f []float64
	n int
}
type csvZ struct {
	z []float64
	n int
}
type csvA csvZ
type csvC struct {
	s []string
	n int
}
type csvS csvC

func (c *csvI) parse(v []string) {
	i, e := strconv.ParseInt(v[c.n], 10, 32)
	if e != nil {
		panic(e)
	}
	c.i = append(c.i, uint32(int32(i)))
}
func (c *csvF) parse(v []string) {
	f, e := strconv.ParseFloat(strings.Replace(v[c.n], ",", ".", -1), 64)
	if e != nil {
		panic(e)
	}
	c.f = append(c.f, f)
}
func (c *csvZ) parse(v []string) { r, i := zstrings(v[c.n], v[1+c.n]); c.z = append(c.z, r, i) }
func (c *csvA) parse(v []string) {
	idx := strings.Index(v[c.n], "a")
	if idx == -1 {
		idx = strings.Index(v[c.n], "@")
	}
	if idx == -1 {
		panic(fmt.Sprintf("cannot parse %q as complex", v[c.n]))
	}
	r, i := zstrings(v[c.n][:idx], v[c.n][idx+1:])
	c.z = append(c.z, r, i)
}
func (c *csvC) parse(v []string) { c.s = append(c.s, v[c.n]) }
func (c *csvS) parse(v []string) { c.s = append(c.s, v[c.n]) }
func zstrings(a, b string) (float64, float64) {
	abs, e := strconv.ParseFloat(strings.Replace(a, ",", ".", -1), 64)
	if e != nil {
		panic(e)
	}
	ang, e := strconv.ParseFloat(strings.Replace(b, ",", ".", -1), 64)
	if e != nil {
		panic(e)
	}
	re, im := math.Sincos(ang / 180.0 * math.Pi)
	return abs * re, abs * im
}

func (c *csvI) k() (r uint32) {
	n := uint32(len(c.i))
	r = mk(2, n)
	copy(I[2+r>>2:], c.i)
	return r
}
func (c *csvF) k() (r uint32) {
	r = mk(3, uint32(len(c.f)))
	copy(F[1+r>>2:], c.f)
	return r
}
func (c *csvZ) k() (r uint32) {
	r = mk(4, uint32(len(c.z)))
	copy(F[1+r>>3:], c.z)
	return r
}
func (c *csvA) k() (r uint32) { return (*csvZ)(c).k() }
func (c *csvC) k() (r uint32) {
	n := uint32(len(c.s))
	r = mk(6, n)
	for i := range c.s {
		I[2+r>>2] = mkcs([]byte(c.s[i]))
	}
	return r
}
func (c *csvS) k() (r uint32) {
	r = mk(5, 0)
	for i := range c.s {
		r = ucat(r, sc(mkcs([]byte(c.s[i]))))
	}
	return r
}
