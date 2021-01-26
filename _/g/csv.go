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
// return value
//  table `a`b`c`d!(values..)
//  table `file`a`b`c`d!(values...) (for file input)
func csv1(x uint32) uint32 {
	d := val(ks("CSV"))
	if d == 0 || tp(d) != 7 {
		panic("csv: var CSV must be a format dict")
	}
	if tp(x) == 5 { // file
		if nn(x) == 1 {
			return csv1(read1(x))
		}
		// multiple files: `file`col1`col2!(`file1`file1`file1..;..)
		rx(x)
		k, v := kvd(d)
		dx(v)
		k = ucat(ks("file"), k)
		v = ech(ech(x, 'C'), 46) // ,/' & .' C'x
		rx(v)
		n := ech(v, 35)      // #'
		v = ech(flp(v), 'u') // ,/'&v
		x = ecd(n, x, 35)    // n#'x
		v = ucat(enl(x), v)
		return mkd(k, v)
	}
	return csv2(d, x)
}
func csv2(x, y uint32) uint32 {
	type parser func(s [][]byte, i int) uint32
	fparse := func(s [][]byte, i int) uint32 {
		f, e := strconv.ParseFloat(strings.Replace(string(s[i]), ",", ".", -1), 64)
		if e != nil {
			return kf(math.Float64frombits(18444492273895866368)) // 0n
		}
		return kf(f)
	}
	k, v := kvd(x)
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

	comma := lupString("COMMA")
	split := func(b []byte) [][]byte {
		return bytes.Split(b, []byte(comma))
	}
	if comma == " " || comma == "" {
		split = bytes.Fields
	}
	b := CK(y)
	lines := bytes.Split(b, []byte{10})
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
	}
	dx(y)
	return mkd(k, table)
}
