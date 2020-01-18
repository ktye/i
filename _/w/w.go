package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

type c = byte

const (
	sFname = iota
	sRety
	sArgs
	sLocl
	sBody
	sNewl
)

type T c

const (
	I = T(0x07)
	F = T(0x7c)
)

// name:I:IIF::body..
type fn struct {
	name string
	rety T
	args []T
	locl []T
	sign int
}

var typs = map[b]T{'I': 0x07, 'F': 0x7c}

type module []f

func main() {
	var html bool
	flag.BoolVar(&html, "html", false, "html output")
	flag.Parse()
	rd := bufio.NewReader(os.Stdin)
	state := sFname
	line, char, hi := 1, 0, true
	err := func(s string) { fatal(fmt.Errorf("%d:%d: %s", line, char, s)) }
	var m module
	var f fn
	var p c
	for {
		b, e := rd.ReadByte()
		if e == io.EOF {
			if state != sFname {
				err("parse trailing")
			}
			wasm := m.emit()
			if html {
				page(wasm)
			} else {
				os.Stdout.Write(wasm)
			}
			return
		}
		char++
		if b == '\n' {
			line++
		}
		fatal(e)
		switch state {
		case sFname:
			if craZ(b) {
				f.name += string(b)
			} else if b == ':' {
				state++
			} else {
				err("parse function name")
			}
		case sRety:
			if f.rety == 0 {
				f.rety = typs[b]
				if f.rety == 0 {
					panic("parse return type")
				}
			} else if b == ':' {
				state++
			} else {
				err("parse return type")
			}
		case sArgs:
			if t := typs[b]; t == 0 && f.args == nil {
				err("parse args")
			} else if t != 0 {
				f.args = append(f.args, t)
			} else if b == ':' {
				state++
			} else {
				err("parse args")
			}
		case sLocl:
			if t := typs[b]; t != 0 {
				f.args = append(f.locl, t)
			} else if b == ':' {
				state++
			} else {
				err("parse locals")
			}
		case sBody:
			if b == '/' && hi {
				state == sCmnt
			} else if (b == ' ' || b == '\t') && hi {
				continue
			} else if (b == '\n') && hi {
				state = sNewl
			} else if crHx(b) {
				if hi {
					p, hi = xtoc(b)<<4, false
				} else {
					p, hi = p|xtoc(b), true
				}
			} else {
				err("parse body")
			}
		case sCmnt:
			if b == '\n' {
				state = sNewl
			}
		case sNewl:
			if b == ' ' || b == '\t' {
				state--
			} else {
				m = append(m, f)
				f = fn{}
				state = fname
			}
		default:
			err("internal parse state")
		}
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
func hxb(x c) (c, c) { h := "0123456789abcdef"; return h[x>>4], h[x&0x0F] }

func cr09(c c) bool { return c >= '0' && c <= '9' }
func craZ(c c) bool { return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') }
func cr0Z(c c) bool { return cr09(c) || craZ(c) }
func crHx(c c) bool { return cr09(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') }
func xtoc(x c) c {
	switch {
	case x < ':':
		return x - '0'
	case x < 'G':
		return 10 + x - 'A'
	default:
		return 10 + x - 'a'
	}
}

func (m module) emit() []c {
	var b bytes.Buffer
	b.Write([]c{0, 0x64, 0x73, 0x6d, 1, 0, 0, 0}) // header

	ts := section{t: 1} // type section
	sigs, sigv := make(map[string]int), make([]string, 0)
	for i, f := range m {
		s := string(f.sig())
		if n, o := sigs[s]; o == false {
			n = len(sigs[s])
			sigs[s] = n
			sigv = append(sigv, s)
			m[i].sign = n
		}
	}
	ts.Buffer.Write(lebu(len(sigv)))
	for i, s := range sigv {
		ts.Buffer.WriteString(s)
	}
	ts.add(b)

	//todo...
}

type section struct {
	bytes.Buffer
}

func (f fn) sig() (r []c) {
	r = append(r, 0x60)
	r = append(r, lebu(f.args))
	for _, t := range f.args {
		t = append(r, t)
	}
	r = append(r, 1)
	r = append(r, f.rety)
	return r
}

func (s section) add(w io.Writer) {
	w.WriteByte(s.t)
	w.Write(lebu(s.Buffer.Size()))
	w.Write(s.Buffer.Bytes())
}

func lebu(b []b, v uint64) []byte { // encode leb128
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

func page(w []c) {
	panic("todo")
}

/*
func lebs(b []b, v int64) []b { // encode signed leb128
	for {
		c := uint8(v & 0x7f)
		s := uint8(v & 0x40)
		v >>= 7
		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}
*/
