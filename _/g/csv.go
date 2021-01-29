package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// parse csv to table
//  csv x     char   (default format in global CSV)
//  csv`file
//  csv`file1`file2
//  csv[x;y]  format; char
// format dict (global CSV)
//  `a`b`c`d!(0;1.0;"2";3a)  "2" column 2 string as symbol, 3a reads abs deg from column 3 and 4
//  `comma`a`b..!(";";..)     default: blanks
//  `skip`comma`..!(2;",";..) header lines
// return value
//  table `a`b`c`d!(values..)
//  table `file`a`b`c`d!(values...) (for file input)
func csv1(x uint32) (r uint32) {
	d := val(ks("CSV"))
	if d == 0 || tp(d) != 7 {
		panic("csv: var CSV must be a format dict")
	}
	if tp(x) == 5 { // file
		xn := nn(x)
		if xn == 0 {
			panic("csv: no inputs")
		} else if xn == 1 {
			rx(x)
			return csv3(d, read1(x), x)
		}
		// multiple files: `file`col1`col2!(`file1`file1`file1..;..)
		for i := i(0); i < xn; i++ {
			rx(x)
			t := csv1(atx(x, mki(i)))
			if i == 0 {
				r = t
			} else {
				r = dcat(r, val(t))
			}
		}
		return r
	}
	return csv3(d, x, 0)
}
func csv2(x, y uint32) uint32 { return csv3(x, y, 0) }
func csv3(x, y, z uint32) uint32 {
	type parser func(s [][]byte, i int) uint32
	fparse := func(s [][]byte, i int) uint32 {
		f, e := strconv.ParseFloat(strings.Replace(string(s[i]), ",", ".", -1), 64)
		if e != nil {
			return kf(math.Float64frombits(18444492273895866368)) // 0n
		}
		return kf(f)
	}
	sskip := ks("skip")
	scoma := ks("comma")
	skip, comma := 0, []byte(nil)
	k, v := kvd(x)
	if MI[2+k>>2] == MI[2+sskip>>2] {
		rx(v)
		skip = iK(fst(v))
		k = drop(k, 1)
		v = drop(v, 1)
	}
	if MI[2+k>>2] == MI[2+scoma>>2] {
		rx(v)
		comma = CK(fst(v))
		k = drop(k, 1)
		v = drop(v, 1)
	}
	dx(sskip)
	dx(scoma)

	table := mk(6, 0)
	index := make([]int, nn(k))
	parse := make([]parser, nn(k))
	for i := uint32(0); i < nn(k); i++ {
		vi := MI[2+i+v>>2]
		if t := tp(vi); t == 1 {
			index[i] = int(ck(vi) - '0')
			parse[i] = func(s [][]byte, i int) uint32 { return ks(string(s[i])) }
			table = lcat(table, mk(5, 0))
		} else if t == 2 {
			index[i] = ik(vi)
			parse[i] = func(s [][]byte, i int) uint32 {
				n, e := strconv.Atoi(string(s[i]))
				if e != nil {
					return ki(0)
				}
				return ki(n)
			}
			table = lcat(table, mk(2, 0))
		} else if t == 3 {
			index[i] = int(fk(vi))
			parse[i] = fparse
			table = lcat(table, mk(3, 0))
		} else if t == 4 {
			index[i] = int(real(zk(vi)))
			parse[i] = func(s [][]byte, i int) uint32 { return atx(fparse(s, i), fparse(s, i+1)) }
			table = lcat(table, mk(4, 0))
		} else {
			panic("csv2: format")
		}
	}
	dx(v)

	split := func(b []byte) [][]byte {
		return bytes.Split(b, comma)
	}
	if len(comma) == 0 {
		split = bytes.Fields
	}
	lines := bytes.Split(CK(y), []byte{10})
	if skip > len(lines) {
		skip = len(lines) - 1
	}
	if len(lines) > 0 && len(lines[0]) > 3 {
		if b := lines[0]; b[0] == 0xef && b[1] == 0xbb && b[2] == 0xbf {
			lines[0] = b[3:] // bom
		}
	}
	lines = lines[skip:]
	m := uint32(0)
	for line, b := range lines {
		line++
		if len(b) > 0 && b[0] == 13 {
			b = b[1:]
		}
		if len(b) > 0 && b[len(b)-1] == 13 {
			b = b[:len(b)-1]
		}
		b = bytes.TrimSpace(b)
		if len(b) == 0 {
			continue
		}
		v := split(b)
		for i := range index {
			n := index[i]
			p := parse[i]
			if n < 0 || n >= len(v) {
				panic(fmt.Sprintf("csv: line %d: parse column %d (only %d)", 1+line, 1+n, len(v)))
			}
			tp := 2 + uint32(i) + table>>2
			MI[tp] = ucat(MI[tp], p(v, n))
		}
		m++
	}
	if z != 0 {
		k = ucat(ks("file"), k)
		table = ucat(enl(take(z, m)), table)
	}
	return mkd(k, table)
}
