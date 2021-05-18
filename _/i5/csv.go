package k

import (
	"bytes"
	gcsv "encoding/csv"
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
	"strings"
)

//csv ""csv"ab|cd|ef\ngh|ij|kl\n" /(("ab";"gh");("cd";"ij");("ef";"kl"))  (auto-detect)
//csv ",if"csv"1,2\n3,4\n5,6\n" /(1 3 5;2 4 6.)
//csv ",2hiffs"csv"x\n\n1,2,0,abc\n2,3,90,gh" /(1 2;2 3.;0 90.;`abc`gh)
//csv ";izs"csv"1;2;0;abc\n2;3;90;gh" /(1 2;2a0 3a90;`abc`gh)
func csvread(x, y T) T {
	c, o := y.(C)
	if o == false {
		panic("type")
	}
	defer dx(y)

	var f csvformat
	switch v := x.(type) {
	case byte:
		f = newCsvFormat([]byte{v}, c.v)
	case C:
		f = newCsvFormat(v.v, c.v)
	default:
		panic("type")
	}
	return f.decode(c.v)
}

func csvwrite(x T) T {
	switch v := x.(type) {
	case D:
		if v.tab {
			defer dx(x)
			return csvl(v.k.(S).v, rx(v.v).(L))
		}
	case L:
		return csvl(nil, v)
	}
	panic("type")
}

type csvformat struct {
	skip    int
	sep     byte
	comment byte
	panic   bool
	typ     []byte
	col     []int
}

func newCsvFormat(b []byte, d []byte) csvformat {
	var f csvformat
	if len(b) == 0 {
		i := bytes.Index(d, []byte{10})
		if i < 0 {
			panic("empty csv")
		}
		h := d[:i]
		var sep = byte(',')
		n := 0
		for _, s := range ",;|\t" {
			if m := bytes.Count(h, []byte{byte(s)}); m > n {
				sep = byte(s)
				n = m
			}
		}
		b = []byte{sep}
	}
	f.sep = b[0]
	b = b[1:]

	i := 0
	next := func(b []byte) (int, byte, []byte) {
		i++
		v, m := tNum(b)
		if m > 0 {
			b = b[m:]
			if len(b) == 0 {
				panic("format")
			}
			n, o := v.(int)
			if !o || n <= 0 {
				panic("format")
			}
			i = n
		}
		n := i
		if b[0] == 'z' {
			i++
		} else if b[0] == 'h' {
			i = 0
		}
		return n, b[0], b[1:]
	}
	var n int
	var t byte
	for b != nil && len(b) > 0 {
		n, t, b = next(b)
		if bytes.IndexByte([]byte("hbcifzs*#!"), t) < 0 {
			panic("format")
		}
		switch t {
		case 'h':
			f.skip = n
		case '#':
			f.comment = '#'
		case '!':
			f.panic = true
		case '*':
		default:
			f.typ = append(f.typ, t)
			f.col = append(f.col, n)
		}
	}
	return f
}

func (f csvformat) decode(b []byte) L {
	for i := 0; i < f.skip; i++ {
		n := bytes.IndexByte(b, 10)
		if n < 0 {
			panic("empty csv")
		}
		b = b[1+n:]
	}
	r := gcsv.NewReader(bytes.NewReader(b))
	r.Comma = rune(f.sep)
	r.Comment = rune(f.comment)
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	a, e := r.ReadAll()
	if e != nil {
		panic(e)
	}
	if len(a) == 0 {
		return KL([]T{})
	}
	if len(f.typ) == 0 {
		f.typ = []byte(strings.Repeat("c", len(a[0])))
		f.col = seq(1, 1+len(f.typ)).v
	}
	l := make([]T, len(f.typ))
	for i := range f.typ {
		l[i] = f.parseColumn(a, f.typ[i], f.col[i]-1)
	}
	return KL(l)
}
func (f csvformat) parseColumn(a [][]string, t byte, j int) T {
	col := f.readcol(a, j)
	switch t {
	case 'b':
		r := make([]bool, len(col))
		for i, s := range col {
			b, e := strconv.ParseBool(s)
			if e != nil && f.panic {
				panic(fmt.Errorf("line %d: col %d: %s", i, j, e))
			}
			r[i] = b
		}
		return KB(r)
	case 'c':
		r := make([]T, len(col))
		for i, s := range col {
			r[i] = KC([]byte(s))
		}
		return KL(r)
	case 'i':
		r := make([]int, len(col))
		for i, s := range col {
			b, e := strconv.Atoi(s)
			if e != nil && f.panic {
				panic(fmt.Errorf("line %d: col %d: %s", i, j, e))
			}
			r[i] = b
		}
		return KI(r)
	case 'f':
		return KF(f.parsefloats(col, j))
	case 'z':
		r := f.parsefloats(col, j)
		p := f.parsefloats(f.readcol(a, 1+j), 1+j)
		z := make([]complex128, len(r))
		for i := range r {
			z[i] = cmplx.Rect(r[i], p[i]*math.Pi/180.0)
		}
		return KZ(z)
	case 's':
		return KS(col)
	default:
		panic("type")
	}
}
func (f csvformat) parsefloats(a []string, j int) []float64 {
	r := make([]float64, len(a))
	for i, s := range a {
		b, e := strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
		if e != nil {
			if f.panic {
				panic(fmt.Errorf("line %d: col %d: %s", i, j, e))
			} else {
				b = math.NaN()
			}
		}
		r[i] = b
	}
	return r
}
func (f csvformat) readcol(a [][]string, j int) []string {
	col := make([]string, len(a))
	for i, cols := range a {
		if j >= len(cols) {
			if f.panic {
				panic(fmt.Errorf("line %d: column range: %d >= %d", i, j, len(cols)))
			} else {
				col[i] = ""
			}
		} else {
			col[i] = a[i][j]
		}
	}
	return col
}

func csvl(h []string, x L) C {
	h, x = zsplit(h, x)
	cols := len(x.v)
	x = flip(x).(L)
	rows := len(x.v)
	defer dx(x)

	var buf bytes.Buffer
	w := gcsv.NewWriter(&buf)
	if h != nil {
		if e := w.Write(h); e != nil {
			panic(e)
		}
	}
	a := make([]string, cols)
	for i := 0; i < rows; i++ {
		rw := x.v[i].(vector)
		for j := range a {
			rw.ref()
			a[j] = astring(atv(rw, j))
		}
		w.Write(a)
	}
	w.Flush()
	return KC(buf.Bytes())
}
func astring(x T) string {
	switch v := x.(type) {
	case bool:
		if v {
			return "true"
		}
		return "false"
	case C:
		return string(v.v)
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case string:
		return v
	default:
		panic("type")
	}
}
func zsplit(h []string, x L) ([]string, L) {
	zcols := make(map[int]bool)
	for i := range x.v {
		if _, o := x.v[i].(Z); o {
			zcols[i] = true
			break
		}
	}
	if len(zcols) == 0 {
		return h, x
	}
	if h != nil {
		hh := make([]string, len(h)+len(zcols))
		j := 0
		for i, s := range h {
			hh[j] = s
			j++
			if zcols[i] {
				hh[j] = "angle"
				j++
			}
		}
		h = hh
	}
	r := make([]T, len(x.v)+len(zcols))
	j := 0
	for i, v := range x.v {
		if zcols[i] {
			r[j], r[j+1] = raz(v.(Z))
			j += 2
		} else {
			r[j] = rx(v)
			j++
		}
	}
	dx(x)
	return h, KL(r)
}
func raz(z Z) (F, F) {
	r := make([]float64, z.ln())
	a := make([]float64, z.ln())
	for i := range z.v {
		r[i], a[i] = cmplx.Polar(z.v[i])
		a[i] *= 180.0 / math.Pi
	}
	return KF(r), KF(a)
}
