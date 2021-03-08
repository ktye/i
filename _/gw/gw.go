package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/mathetake/gasm/wasm"
)

func main() {
	b, e := ioutil.ReadFile(os.Args[1])
	fatal(e)
	m, e := wasm.DecodeModule(bytes.NewReader(b))
	fatal(e)

	ftyps := m.SecFunctions // SecFunctions does not contain typs for imports
	for _, i := range m.SecImports {
		if i.Desc.Kind == 0 {
			ftyps = append([]uint32{*i.Desc.TypeIndexPtr}, ftyps...)
		}
	}
	fn(m, 0, ftyps, os.Stdout)
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

type stringer interface{ String() string }
type opener interface{ Open() }
type stack struct {
	p []stringer
}

func (s *stack) push(v stringer) {
	s.p = append(s.p, v)
}
func (s *stack) pop() stringer {
	i := len(s.p) - 1
	v := s.p[i]
	s.p = s.p[:i]
	return v
}
func (s *stack) pops(n int) []stringer {
	r := make([]stringer, n)
	for j := n - 1; j >= 0; j-- {
		r[j] = s.pop()
	}
	return r
}

type u32const uint32
type i32const uint32
type i32load struct {
	align, offset int32
	s             stringer
}
type u32 struct{ stringer }
type i32eqz struct{ stringer }
type i32add struct{ y, x stringer }
type i32sub struct{ y, x stringer }
type i32mul struct{ y, x stringer }
type i32or struct{ y, x stringer }
type i32shl struct{ y, x stringer }
type i32eq struct{ y, x stringer }
type i32ltu struct{ y, x stringer }
type i32gtu struct{ y, x stringer }
type localget int
type call struct {
	n int
	a []stringer
}
type ret []stringer

func (u u32const) String() string { return strconv.FormatUint(uint64(u), 10) }
func (i i32const) String() string { return strconv.FormatInt(int64(i), 10) }
func (l i32load) String() string  { return fmt.Sprintf("i32load(%d, %d, %s)", l.align, l.offset, l.s) }
func (u u32) String() string      { return "uint32(" + u.stringer.String() + ")" }
func (i i32eqz) String() string   { return fmt.Sprintf("(0==%s)", i.stringer.String()) }
func (i i32add) String() string   { return "(" + i.x.String() + "+" + i.y.String() + ")" }
func (i i32sub) String() string   { return "(" + i.x.String() + "-" + i.y.String() + ")" }
func (i i32mul) String() string   { return "(" + i.x.String() + "*" + i.y.String() + ")" }
func (i i32or) String() string    { return "(" + i.x.String() + "|" + i.y.String() + ")" }
func (i i32shl) String() string   { return "(" + i.x.String() + "<<" + i.y.String() + ")" }
func (i i32eq) String() string    { return "(" + i.x.String() + "==" + i.y.String() + ")" }
func (i i32ltu) String() string   { return "(" + i.x.String() + "<" + i.y.String() + ")" }
func (i i32gtu) String() string   { return "(" + i.x.String() + ">" + i.y.String() + ")" }
func (l localget) String() string { return "x" + strconv.Itoa(int(l)) }
func (c call) String() string {
	v := make([]string, len(c.a))
	for i, s := range c.a {
		v[i] = s.String()
	}
	return "f" + strconv.Itoa(c.n) + "(" + strings.Join(v, ", ") + ")"
}
func (r ret) String() string {
	v := make([]string, len(r))
	for i := range r {
		v[i] = r[i].String()
	}
	return "return " + strings.Join(v, ", ") + "\n"
}

func fn(m *wasm.Module, idx int, ftyps []uint32, w io.Writer) {
	c := m.SecCodes[idx]
	t := m.SecTypes[ftyps[idx]]

	for i, u := range m.SecTypes {
		fmt.Printf("SecTypes[%d] = %v\n", i, u)
	}
	for i, u := range ftyps {
		fmt.Printf("ftyps[%d] = %d\n", i, u)
	}

	nret := len(t.ReturnTypes)

	typelist := func(v []wasm.ValueType, in bool) string {
		tp := map[byte]string{0x7f: "uint32", 0x7e: "uint64", 0x7d: "float32", 0x7c: "float64"}
		var s []string
		for i, t := range v {
			u := tp[byte(t)]
			if in {
				u = "x" + strconv.Itoa(i) + " " + u
			}
			s = append(s, u)
		}
		return strings.Join(s, ", ")
	}
	fmt.Fprintf(w, "function f%d(%s) (%s) {\n", idx, typelist(t.InputTypes, true), typelist(t.ReturnTypes, false))
	s := stack{}
	r := c.Body[0:]
	locs := make(map[int32]bool)
	for len(r) > 0 {
		b := r[0]

		h := r[:len(r)]
		if len(h) > 8 {
			h = h[:8]
		}
		//fmt.Printf("> %x..\n", h)

		switch b { // https://webassembly.github.io/spec/core/appendix/index-instructions.html
		case 0x04: // if
			var t byte
			t = r[1]
			r = r[2:]
			if t != 0x40 {
				panic("todo: if-type")
			}
			fmt.Fprintf(w, "if %s {\n", s.pop())
		case 0x05: // else
			r = r[1:]
			fmt.Fprintf(w, "%s\n", s.pop())
			fmt.Fprintf(w, "} else {\n")
		case 0x0b: //end
			r = r[1:]
			fmt.Fprintf(w, "%s\n", s.pop())
			fmt.Fprintf(w, "}\n")
		case 0x0f: // return
			r = r[1:]
			s.push(ret(s.pops(nret)))
		case 0x10: // call
			var i int32
			r = immi32(r[1:], &i)
			//t := ftyps[i]
			//fmt.Println("idx", i, "t", t, m.SecTypes[t])
			//fmt.Println("func: ", typelist(m.SecTypes[ftyps[i]].InputTypes, true))
			s.push(call{int(i), s.pops(len(m.SecTypes[ftyps[i]].InputTypes))}) // todo multiple return values (not in one expression)
		case 0x20: // local.get
			var i int32
			r = immi32(r[1:], &i)
			s.push(localget(i))
		case 0x21: // local.set
			var i int32
			r = immi32(r[1:], &i)
			if locs[i] == false {
				fmt.Fprintf(w, "x%d := %s\n", i, s.pop().String())
				locs[i] = true
			} else {
				fmt.Fprintf(w, "x%d = %s\n", i, s.pop().String())
			}
		case 0x28: // i32.load memarg
			var align, offset int32
			r = immi32(r[1:], &align)
			r = immi32(r, &offset)
			s.push(i32load{align, offset, s.pop()})
		case 0x41: // i32.const
			var i int32
			r = immi32(r[1:], &i)
			//fmt.Println("const ", i)
			s.push(i32const(i))
		case 0x45: // i32.eqz
			r = r[1:]
			s.push(i32eqz{s.pop()})
		case 0x46: // i32.eq
			r = r[1:]
			s.push(i32eq{s.pop(), s.pop()})
		case 0x49: // i32.ltu
			r = r[1:]
			s.push(i32ltu{s.pop(), s.pop()})
		case 0x4b: // i32.gtu
			r = r[1:]
			s.push(i32gtu{s.pop(), s.pop()})
		case 0x6a: // i32.add
			r = r[1:]
			s.push(i32add{s.pop(), s.pop()})
		case 0x6b: // i32.sub
			r = r[1:]
			s.push(i32sub{s.pop(), s.pop()})
		case 0x6c: // i32.mul
			r = r[1:]
			s.push(i32mul{s.pop(), s.pop()})
		case 0x72: // i32.or
			r = r[1:]
			s.push(i32or{s.pop(), s.pop()})
		case 0x74: // i32.shl
			r = r[1:]
			s.push(i32shl{s.pop(), s.pop()})
		default:
			panic(fmt.Sprintf("unknown %x\n", b))
		}
	}
	if nret > 0 {
		s.push(ret(s.pops(nret)))
		w.Write([]byte(s.pop().String()))
	}
	fmt.Fprintf(w, "}\n")
}

func immu32(r []byte, u *uint32) []byte {
	const (
		p uint32 = 1 << 7
		q        = ^p
	)
	num := 0
	for shift := 0; shift < 35; shift += 7 {
		b := uint32(r[num])
		num++
		*u |= (b & q) << shift
		if b&p == 0 {
			break
		}
	}
	return r[num:]
}
func immi32(r []byte, i *int32) []byte {
	const (
		int32Mask  int32 = 1 << 7
		int32Mask2       = ^int32Mask
		int32Mask3       = 1 << 6
		int32Mask4       = ^0
	)
	var shift int
	var b int32
	num := 0
	for shift < 35 {
		b := int32(r[num])
		num++
		*i |= (b & int32Mask2) << shift
		shift += 7
		if b&int32Mask == 0 {
			break
		}
	}
	if shift < 32 && (b&int32Mask3) == int32Mask3 {
		*i |= int32Mask4 << shift
	}
	return r[num:]
}

/*
SecCodes     []*CodeSegment
type CodeSegment struct {
	NumLocals uint32
	Body      []byte
}
*/
