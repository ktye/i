package main

import (
	"bufio"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

var line int
var file string

func main() {
	ini()
	table[21] = out
	rd := os.Stdin
	if len(os.Args) < 2 {
	} else {
		file = os.Args[1]
		if f, err := os.Open(file); err != nil {
			fatal(err.Error())
		} else {
			defer f.Close()
			rd = f
		}
	}
	sp()
	r := bufio.NewScanner(rd)
	for r.Scan() {
		line++
		do(r.Bytes())
		sp()
	}
}
func sp() {
	if file == "" {
		os.Stdout.Write([]c{' '})
	}
}
func do(s []byte) {
	defer stk()
	s = cmd(s)
	ns := k(len(s))
	c := mk(C, ns)
	cc := 8 + c<<2
	copy(m.c[cc:cc+ns], s)
	p := prs(c)
	r := evl(inc(p))
	if !isasn(p) {
		if m.k[r]>>28 != N {
			nl := mk(C, atom)
			m.c[8+nl<<2] = '\n'
			dec(out(cat(kst(evl(p)), nl)))
		} else {
			dec(r)
		}
	}
}
func out(x k) k {
	if t, n := typ(x); t != C {
		panic("type")
	} else {
		xp := ptr(x, t)
		print(string(m.c[xp : xp+n]))
	}
	return x
}
func isasn(x k) bool {
	if t, n := typ(x); t == L && n > 1 && m.k[m.k[2+x]]>>28 == N+2 && m.k[2+m.k[2+x]] == dyad {
		return true
	}
	return false
}
func cmd(b []byte) []byte {
	if len(b) < 2 || b[0] != '\\' {
		return b
	}
	switch b[1] {
	case 'v':
		return []c("lsv 0\n")
	case 'c':
		return []c("clv 0\n")
	case 'h':
		return []c("help 0\n")
	case '\\':
		panic("bye")
	default:
		return b
	}
}
func stk() {
	if c := recover(); c != nil {
		a, b := stack(c)
		println(a)
		println(file+":"+strconv.Itoa(line)+":", b)
		if file != "" {
			os.Exit(1)
		}
	}
}
func stack(c interface{}) (stk, err string) {
	for _, s := range strings.Split(string(debug.Stack()), "\n") {
		if strings.HasPrefix(s, "\t") {
			stk += "\n" + s[1:]
		}
	}
	err = "?"
	if s, o := c.(string); o {
		err = s
	} else if e, o := c.(error); o {
		err = e.Error()
	}
	return stk, err
}
func fatal(s string) { println(s); os.Exit(1) }
