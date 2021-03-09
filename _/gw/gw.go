package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/go-interpreter/wagon/wasm"
)

func main() {
	b, e := ioutil.ReadFile(os.Args[1])
	fatal(e)
	Go(bytes.NewReader(b), os.Stdout)
}

func Go(r io.Reader, w io.Writer) {
	m, e := wasm.DecodeModule(r)
	fatal(e)
	w.Write([]byte(head))

	w.Write([]byte("func init() {"))
	var buf bytes.Buffer
	max := 0
	for _, e := range m.Elements.Entries {
		off := decodeOffset(e.Offset)
		for k, u := range e.Elems {
			idx := k + int(off)
			if idx > max {
				max = idx
			}
			fmt.Fprintf(&buf, "F[%d]=f%d;", idx, u)
		}
	}
	fmt.Fprintf(w, "F = make([]interface{}, %d);", 1+max)
	io.Copy(w, &buf)
	w.Write([]byte("}\n"))

	ext := 0
	for _, im := range m.Import.Entries {
		if im.Type.Kind() == wasm.ExternalFunction {
			extfn(ext, m, w)
			ext++
		}
	}

	for i, body := range m.Code.Bodies {
		fn(i+ext, m, body, w)
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
func decodeOffset(b []byte) (i int32) {
	if b[0] != 0x41 {
		panic("offset is not a constant")
	}
	b = immi32(b[1:], &i)
	if len(b) != 1 || b[0] != 0x0b {
		panic("offset has remaining bytes")
	}
	return i
}

type stringer interface{ String() string }
type booler interface{ Bool() string }
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

func boolean(s stringer) string {
	if b, ok := s.(booler); ok {
		return b.Bool()
	} else {
		return "(0 != " + s.String() + ")"
	}
}

type i32const uint32
type i32load struct {
	s             stringer
	align, offset int32
}
type i32store struct {
	s, v          stringer
	align, offset int32
}
type i32 struct{ stringer }
type i32eqz struct{ x stringer }
type i32clz struct{ x stringer }
type i32op2 struct {
	y, x stringer
	op   string
}
type i32cmp struct {
	y, x stringer
	op   string
}

type localget int
type call struct {
	n int
	a []stringer
}
type icall struct {
	n       int
	in, out []wasm.ValueType
	a       []stringer
}
type ret []stringer

func (i i32const) String() string { return "uint32(" + strconv.FormatUint(uint64(i), 10) + ")" }
func (l i32load) String() string  { return fmt.Sprintf("i32load(%s, %d, %d)", l.s, l.align, l.offset) }
func (l i32store) String() string {
	return fmt.Sprintf("i32store(%s, %s, %d, %d)", l.s, l.v, l.align, l.offset)
}
func (i i32) String() string {
	s := i.stringer.String()
	if c, ok := i.stringer.(i32const); ok {
		var n int32 = int32(c)
		s = strconv.Itoa(int(n))
	}
	return "int32(" + s + ")"
}
func (i i32eqz) String() string { return "ub(0 == " + i.x.String() + ")" }
func (i i32eqz) Bool() string   { return "(!" + boolean(i.x) + ")" }
func (i i32clz) String() string { return "uint32(bits.LeadingZeros32(" + i.x.String() + "))" }
func (i i32op2) String() string { return "uint32(" + i.x.String() + i.op + i.y.String() + ")" }
func (i i32cmp) String() string { return "ub(" + i.x.String() + i.op + i.y.String() + ")" }
func (i i32cmp) Bool() string   { return "(" + i.x.String() + i.op + i.y.String() + ")" }

func (l localget) String() string { return "x" + strconv.Itoa(int(l)) }
func (c call) String() string {
	v := make([]string, len(c.a))
	for i, s := range c.a {
		v[i] = s.String()
	}
	return "f" + strconv.Itoa(c.n) + "(" + strings.Join(v, ", ") + ")"
}
func (c icall) String() string {
	v := make([]string, len(c.a))
	for i, s := range c.a {
		v[i] = s.String()
	}
	return fmt.Sprintf("F[%d].(func(%s)(%s))(%s)", c.n, typelist(c.in, false), typelist(c.out, false), strings.Join(v, ", "))
}
func (r ret) String() string {
	v := make([]string, len(r))
	for i := range r {
		v[i] = r[i].String()
	}
	return "return " + strings.Join(v, ", ") + "\n"
}
func typestr(b byte) string {
	tp := map[byte]string{0x7f: "uint32", 0x7e: "uint64", 0x7d: "float32", 0x7c: "float64"}
	return tp[b]
}
func typelist(v []wasm.ValueType, in bool) string {
	var s []string
	for i, t := range v {
		u := typestr(byte(t))
		if in {
			u = "x" + strconv.Itoa(i) + " " + u
		}
		s = append(s, u)
	}
	return strings.Join(s, ", ")
}
func extfn(idx int, m *wasm.Module, w io.Writer) {
	sig, e := m.GetFunctionSig(uint32(idx))
	fatal(e)
	fmt.Fprintf(w, "func f%d(%s) (%s) {fmt.Println(\"f%d\"); return}\n", idx, typelist(sig.ParamTypes, true), typelist(sig.ReturnTypes, false), idx)
}
func fn(idx int, m *wasm.Module, body wasm.FunctionBody, w io.Writer) {
	sig, e := m.GetFunctionSig(uint32(idx))
	fatal(e)
	nret := len(sig.ReturnTypes)

	fmt.Fprintf(w, "func f%d(%s) (%s) {\n", idx, typelist(sig.ParamTypes, true), typelist(sig.ReturnTypes, false))
	lidx := int32(len(sig.ParamTypes))
	unused := make(map[int32]bool)
	for _, l := range body.Locals {
		s := ""
		for k := uint32(0); k < l.Count; k++ {
			if k > 0 {
				s += ", "
			}
			s += "x" + strconv.Itoa(int(lidx))
			unused[lidx] = true
			lidx++
		}
		fmt.Fprintf(w, "var %s %s\n", s, typestr(byte(l.Type)))
	}
	s := stack{}
	r := body.Code[0:]
	blocks := 0
	var labels []string
	for len(r) > 0 {
		b := r[0]

		h := r[:len(r)]
		if len(h) > 8 {
			h = h[:8]
		}
		//fmt.Printf("> %x..\n", h)

		br := func() {
			var i int32
			r = immi32(r[1:], &i)
			l := labels[len(labels)-1-int(i)]
			fmt.Fprintf(w, "goto %s\n", l)
		}
		newlabel := func(s string) string {
			l := s + strconv.Itoa(blocks)
			labels = append(labels, l)
			blocks++
			return l
		}
		localset := func() (i int32) {
			r = immi32(r[1:], &i)
			fmt.Fprintf(w, "x%d = %s\n", i, s.pop().String())
			return i
		}

		switch b { // https://webassembly.github.io/spec/core/appendix/index-instructions.html
		case 0x00: // unreachable
			r = r[1:]
			w.Write([]byte("panic(`unreachable`)\n"))
		case 0x02: // block
			var t byte
			t = r[1]
			r = r[2:]
			if t != 0x40 {
				panic("todo: block-type")
			}
			fmt.Fprintf(w, "//%s\n", newlabel("B"))
		case 0x03: // loop
			var t byte
			t = r[1]
			r = r[2:]
			if t != 0x40 {
				panic("todo: loop-type")
			}
			fmt.Fprintf(w, "%s:\n", newlabel("L"))
		case 0x04: // if
			var t byte
			t = r[1]
			r = r[2:]
			if t != 0x40 {
				panic("todo: if-type")
			}
			newlabel("I")
			fmt.Fprintf(w, "if %s {\n", boolean(s.pop()))
		case 0x05: // else
			r = r[1:]
			if len(s.p) > 0 {
				//fmt.Fprintf(w, "%s\n", s.pop())
			}
			fmt.Fprintf(w, "} else {\n")
		case 0x0b: //end
			r = r[1:]
			if len(s.p) > 0 {
				//fmt.Fprintf(w, "%s\n", s.pop())
			}
			l := labels[len(labels)-1]
			labels = labels[:len(labels)-1]
			if l[0] == 'B' {
				fmt.Fprintln(w, l+": //end")
			} else if l[0] == 'I' {
				fmt.Fprintf(w, "} //%s\n", l)
			} else {
				fmt.Fprintln(w, "//end "+l)
			}
		case 0x0c:
			br() // br
		case 0x0d: // brif
			fmt.Fprintf(w, "if %s {\n", boolean(s.pop()))
			br()
			fmt.Fprintf(w, "}\n")
		case 0x0f: // return
			r = r[1:]
			w.Write([]byte(ret(s.pops(nret)).String()))
		case 0x10: // call
			var i int32
			r = immi32(r[1:], &i)
			t, e := m.GetFunctionSig(uint32(i))
			fatal(e)

			c := call{int(i), s.pops(len(t.ParamTypes))}
			if nr := len(t.ReturnTypes); nr == 0 {
				w.Write([]byte(c.String() + ";\n"))
			} else if nr == 1 {
				s.push(c)
			} else {
				panic("todo: call: multiple return values")
			}
		case 0x11: // call indirect
			var i int32
			r = immi32(r[1:], &i)
			t, e := m.GetFunctionSig(uint32(i))
			fatal(e)
			c := icall{int(i), t.ParamTypes, t.ReturnTypes, s.pops(len(t.ParamTypes))}
			if nr := len(t.ReturnTypes); nr == 0 {
				w.Write([]byte(c.String() + ";\n"))
			} else if nr == 1 {
				s.push(c)
			} else {
				panic("todo: indirect-call: multiple return values")
			}
		case 0x20: // local.get
			var i int32
			r = immi32(r[1:], &i)
			if _, ok := unused[i]; ok {
				unused[i] = false
			}
			s.push(localget(i))
		case 0x21: // local.set
			localset()
		case 0x22: // local.tee
			i := localset()
			s.push(localget(i))
		case 0x28: // i32.load memarg
			var align, offset int32
			r = immi32(r[1:], &align)
			r = immi32(r, &offset)
			s.push(i32load{s.pop(), align, offset})
		case 0x36: // i32.store memarg
			var align, offset int32
			r = immi32(r[1:], &align)
			r = immi32(r, &offset)
			w.Write([]byte(i32store{s.pop(), s.pop(), align, offset}.String() + "\n"))
		case 0x41: // i32.const
			var i int32
			r = immi32(r[1:], &i)
			s.push(i32const(i))
		case 0x45: // i32.eqz
			r = r[1:]
			s.push(i32eqz{s.pop()})
		case 0x46: // i32.eq
			r = r[1:]
			s.push(i32cmp{s.pop(), s.pop(), "=="})
		case 0x47: // i32.ne
			r = r[1:]
			s.push(i32cmp{s.pop(), s.pop(), "!="})
		case 0x48: // i32.lts
			r = r[1:]
			s.push(i32cmp{i32{s.pop()}, i32{s.pop()}, "<"})
		case 0x49: // i32.ltu
			r = r[1:]
			s.push(i32cmp{s.pop(), s.pop(), "<"})
		case 0x4a: // i32.gts
			r = r[1:]
			s.push(i32cmp{i32{s.pop()}, i32{s.pop()}, ">"})
		case 0x4b: // i32.gtu
			r = r[1:]
			s.push(i32cmp{s.pop(), s.pop(), ">"})
		case 0x4c: // i32.les
			r = r[1:]
			s.push(i32cmp{i32{s.pop()}, i32{s.pop()}, "<="})
		case 0x4d: // i32.leu
			r = r[1:]
			s.push(i32cmp{s.pop(), s.pop(), "<="})
		case 0x4e: // i32.ges
			r = r[1:]
			s.push(i32cmp{i32{s.pop()}, i32{s.pop()}, ">="})
		case 0x4f: // i32.geu
			r = r[1:]
			s.push(i32cmp{s.pop(), s.pop(), ">="})
		case 0x6a: // i32.add
			r = r[1:]
			s.push(i32op2{s.pop(), s.pop(), "+"})
		case 0x6b: // i32.sub
			r = r[1:]
			s.push(i32op2{s.pop(), s.pop(), "-"})
		case 0x6c: // i32.mul
			r = r[1:]
			s.push(i32op2{s.pop(), s.pop(), "*"})
		case 0x6d: // i32.divs
			r = r[1:]
			s.push(i32op2{i32{s.pop()}, i32{s.pop()}, "/"})
		case 0x6e: // i32.divu
			r = r[1:]
			s.push(i32op2{s.pop(), s.pop(), "/"})
		case 0x6f: // i32.rems
			r = r[1:]
			s.push(i32op2{i32{s.pop()}, i32{s.pop()}, "%"})
		case 0x70: // i32.remu
			r = r[1:]
			s.push(i32op2{s.pop(), s.pop(), "%"})
		case 0x67: //i32.clz
			r = r[1:]
			s.push(i32clz{s.pop()})
		case 0x71: // i32.and
			r = r[1:]
			s.push(i32op2{s.pop(), s.pop(), "&"})
		case 0x72: // i32.or
			r = r[1:]
			s.push(i32op2{s.pop(), s.pop(), "|"})
		case 0x74: // i32.shl
			r = r[1:]
			s.push(i32op2{s.pop(), s.pop(), "<<"})
		case 0x75: // i32.shr_s
			r = r[1:]
			s.push(i32op2{i32{s.pop()}, i32{s.pop()}, ">>"})
		case 0x76: // i32.shr_u
			r = r[1:]
			s.push(i32op2{s.pop(), s.pop(), ">>"})
		default:
			panic(fmt.Sprintf("unknown %x\n", b))
		}
	}
	for i, b := range unused {
		if b {
			fmt.Fprintf(w, "_ = x%d\n", i)
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

const head = `//+build ignore

package main

import (
	"encoding/binary"
	"math/bits"
	"fmt"
)

var M []byte
var F []interface{}
func main() {}
func i32load(addr, align, offset uint32) uint32 { return binary.LittleEndian.Uint32(M[addr+offset:]) }
func i32store(addr, value, align, offset uint32) { binary.LittleEndian.PutUint32(M[addr+offset:], value) }
func ub(b bool) uint32 { if b { return 1 } else { return 0 } }
`
