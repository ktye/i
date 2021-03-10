package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
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
	initFuncTable(m, w)
	initMemory(m, w)
	w.Write([]byte("}\n"))

	ext := writeImports(m, w)

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
	if len(b) != 2 || b[1] != 0x0b {
		panic("offset has remaining bytes")
	}
	return i
}
func initFuncTable(m *wasm.Module, w io.Writer) {
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
}
func initMemory(m *wasm.Module, w io.Writer) {
	if m.Memory == nil {
		fmt.Fprintf(w, "M = make([]byte, 64 * 1024)\n")
		return
	}
	if n := len(m.Memory.Entries); n != 1 {
		panic(fmt.Errorf("number of memory entries: %d (not 1)", n))
	}
	n := m.Memory.Entries[0].Limits.Initial
	fmt.Fprintf(w, "M = make([]byte, 64 * 1024 * %d)\n", n)
}
func writeImports(m *wasm.Module, w io.Writer) int {
	ext := 0
	if m.Import != nil {
		for _, im := range m.Import.Entries {
			if im.Type.Kind() == wasm.ExternalFunction {
				extfn(ext, m, w)
				ext++
			}
		}
	}
	return ext
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
type i64const uint64
type f32const float32
type f64const float64
type iload struct {
	s                  stringer
	align, offset      int32
	retsize, storesize int
	sign               byte
}
type istore struct {
	v, s                 stringer
	align, offset        int32
	inputsize, storesize int
}
type fload struct {
	s             stringer
	align, offset int32
	width         int
}
type fstore struct {
	v, s          stringer
	align, offset int32
	width         int
}
type ii struct{ x stringer }
type i8 struct{ x stringer }
type i16 struct{ x stringer }
type i32 struct{ x stringer }
type u32 struct{ x stringer }
type i64 struct{ x stringer }
type u64 struct{ x stringer }
type f32 struct{ x stringer }
type f64 struct{ x stringer }
type i32eqz struct{ x stringer }
type i32op1 struct {
	x  stringer
	op string
}
type i32op2 struct {
	y, x stringer
	op   string
}
type i64op1 struct {
	x  stringer
	op string
}
type i64op2 struct {
	y, x stringer
	op   string
}
type cmp struct {
	y, x stringer
	op   string
}

type localget int
type call struct {
	n int
	a []stringer
}
type icall struct {
	n       stringer
	in, out []wasm.ValueType
	a       []stringer
}
type ret []stringer
type fltf1 struct {
	size int
	x    stringer
	op   string
}
type fltf2 struct {
	size int
	x, y stringer
	op   string
}
type reinterp struct {
	x  stringer
	op string
}

func (i i32const) String() string { return "uint32(" + strconv.FormatUint(uint64(i), 10) + ")" }
func (i i64const) String() string { return "uint64(" + strconv.FormatUint(uint64(i), 10) + ")" }
func (i f32const) String() string { return fmt.Sprintf("float32(%v)", float32(i)) }
func (i f64const) String() string { return fmt.Sprintf("float64(%v)", float64(i)) }
func (l iload) String() string {
	return fmt.Sprintf("li%d%c%d(%s, %d, %d)", l.retsize, l.sign, l.storesize, l.s, l.align, l.offset)
}
func (l istore) String() string {
	return fmt.Sprintf("si%du%d(%s, %s, %d, %d)", l.inputsize, l.storesize, l.v, l.s, l.align, l.offset)
}
func (l fload) String() string {
	return fmt.Sprintf("lf%d(%s, %d, %d)", l.width, l.s, l.align, l.offset)
}
func (l fstore) String() string {
	return fmt.Sprintf("sf%d(%s, %s, %d, %d)", l.width, l.v, l.s, l.align, l.offset)
}
func (i ii) String() string  { return "int(" + i.x.String() + ")" }
func (i i8) String() string  { return "int8(" + i.x.String() + ")" }
func (i i16) String() string { return "int16(" + i.x.String() + ")" }
func (i i32) String() string {
	s := i.x.String()
	if c, ok := i.x.(i32const); ok {
		var n int32 = int32(c)
		s = strconv.Itoa(int(n))
	}
	return "int32(" + s + ")"
}
func (i i64) String() string {
	s := i.x.String()
	if c, ok := i.x.(i64const); ok {
		var n int64 = int64(c)
		s = strconv.FormatInt(n, 10)
	}
	return "int64(" + s + ")"
}
func (i u32) String() string    { return "uint32(" + i.x.String() + ")" }
func (i u64) String() string    { return "uint64(" + i.x.String() + ")" }
func (i f32) String() string    { return "float32(" + i.x.String() + ")" }
func (i f64) String() string    { return "float64(" + i.x.String() + ")" }
func (i i32eqz) String() string { return "ub(0 == " + i.x.String() + ")" }
func (i i32eqz) Bool() string   { return "(!" + boolean(i.x) + ")" }
func (i i32op1) String() string { return "uint32(" + i.op + "(" + i.x.String() + "))" }
func (i i32op2) String() string {
	if len(i.op) > 2 {
		return fmt.Sprintf("uint32(%s(%s,%s))"+i.op, i.x.String(), i.y.String())
	} else {
		return "uint32(" + i.x.String() + i.op + i.y.String() + ")"
	}
}
func (i i64op1) String() string { return "uint64(" + i.op + "(" + i.x.String() + "))" }
func (i i64op2) String() string { return "uint64(" + i.x.String() + i.op + i.y.String() + ")" }
func (i cmp) String() string    { return "ub(" + i.x.String() + i.op + i.y.String() + ")" }
func (i cmp) Bool() string      { return "(" + i.x.String() + i.op + i.y.String() + ")" }

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
	return fmt.Sprintf("F[%s].(func(%s)(%s))(%s)", c.n.String(), typelist(c.in, false), typelist(c.out, false), strings.Join(v, ", "))
}
func (r ret) String() string {
	v := make([]string, len(r))
	for i := range r {
		v[i] = r[i].String()
	}
	return "return " + strings.Join(v, ", ") + "\n"
}
func (o fltf1) String() string {
	if len(o.op) == 1 {
		return o.op + o.x.String()
	}
	if o.size == 32 {
		return "float32(" + o.op + "(float64(" + o.x.String() + ")))"
	} else {
		return o.op + "(" + o.x.String() + ")"
	}
}
func (o fltf2) String() string {
	if len(o.op) == 1 {
		return "(" + o.x.String() + o.op + o.y.String() + ")"
	}
	if o.size == 32 {
		return "float32(" + o.op + "(float64(" + o.x.String() + "), float64(" + o.y.String() + ")))"
	} else {
		return o.op + "(" + o.x.String() + ", " + o.y.String() + ")"
	}
}
func (r reinterp) String() string { return r.op + "(" + r.x.String() + ")" }

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
		loadint := func(ret, store int, sign byte) {
			var align, offset int32
			r = immi32(r[1:], &align)
			r = immi32(r[1:], &offset)
			s.push(iload{s.pop(), align, offset, ret, store, sign})
		}
		storeint := func(in, store int) {
			var align, offset int32
			r = immi32(r[1:], &align)
			r = immi32(r[1:], &offset)
			w.Write([]byte(istore{s.pop(), s.pop(), align, offset, in, store}.String() + "\n"))
		}
		loadfloat := func(width int) {
			var align, offset int32
			r = immi32(r[1:], &align)
			r = immi32(r[1:], &offset)
			s.push(fload{s.pop(), align, offset, width})
		}
		storefloat := func(size int) {
			var align, offset int32
			r = immi32(r[1:], &align)
			r = immi32(r[1:], &offset)
			w.Write([]byte(fstore{s.pop(), s.pop(), align, offset, size}.String() + "\n"))
		}
		float1 := func(size int, op string) { s.push(fltf1{size: size, x: s.pop(), op: op}) }
		float2 := func(size int, op string) { s.push(fltf2{size: size, y: s.pop(), x: s.pop(), op: op}) }

		switch b { // https://webassembly.github.io/spec/core/appendix/index-instructions.html
		case 0x00: // unreachable
			w.Write([]byte("panic(`unreachable`)\n"))
		case 0x02: // block
			var t byte
			t = r[1]
			r = r[1:]
			if t != 0x40 {
				panic("todo: block-type")
			}
			fmt.Fprintf(w, "//%s\n", newlabel("B"))
		case 0x03: // loop
			var t byte
			t = r[1]
			r = r[1:]
			if t != 0x40 {
				panic("todo: loop-type")
			}
			fmt.Fprintf(w, "%s:\n", newlabel("L"))
		case 0x04: // if
			var t byte
			t = r[1]
			r = r[1:]
			if t != 0x40 {
				panic("todo: if-type")
			}
			newlabel("I")
			fmt.Fprintf(w, "if %s {\n", boolean(s.pop()))
		case 0x05: // else
			if len(s.p) > 0 {
				//fmt.Fprintf(w, "%s\n", s.pop())
			}
			fmt.Fprintf(w, "} else {\n")
		case 0x0b: //end
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
			if r[1] != 0 {
				panic("call-indirect: table-index must be 0")
			}
			r = r[1:]
			t := m.Types.Entries[uint32(i)]
			c := icall{s.pop(), t.ParamTypes, t.ReturnTypes, s.pops(len(t.ParamTypes))}
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
		case 0x28: // i32.load
			loadint(32, 32, 'u')
		case 0x29: // i64.load
			loadint(64, 64, 'u')
		case 0x2a: // f32.load
			loadfloat(32)
		case 0x2b: // f64.load
			loadfloat(64)
		case 0x2c: // i32.load8s
			loadint(32, 8, 's')
		case 0x2d: // i32.load8u
			loadint(32, 8, 'u')
		case 0x2e: // i32.load16s
			loadint(32, 16, 's')
		case 0x2f: // i32.load16u
			loadint(32, 16, 'u')
		case 0x30: // i64.load8s
			loadint(64, 8, 's')
		case 0x31: // i64.load8u
			loadint(64, 8, 'u')
		case 0x32: // i64.load16s
			loadint(64, 16, 's')
		case 0x33: // i64.load16u
			loadint(64, 16, 'u')
		case 0x34: // i64.load32s
			loadint(64, 32, 's')
		case 0x35: // i64.load32u
			loadint(64, 32, 'u')
		case 0x36: // i32.store
			storeint(32, 32)
		case 0x37: // i64.store
			storeint(64, 64)
		case 0x38: // f32.store
			storefloat(32)
		case 0x39: // f64.store
			storefloat(64)
		case 0x3a: // i32.store8
			storeint(32, 8)
		case 0x3b: // i32.store16
			storeint(32, 16)
		case 0x3d: // i64.store8
			storeint(64, 8)
		case 0x3e: // i64.store16
			storeint(64, 16)
		case 0x3c: // i64.store32
			storeint(64, 32)
		case 0x41: // i32.const
			var i int32
			r = immi32(r[1:], &i)
			s.push(i32const(i))
		case 0x42: // i64.const
			var i int64
			r = immi64(r[1:], &i)
			s.push(i64const(i))
		case 0x43: // f32const
			f := math.Float32frombits(binary.LittleEndian.Uint32(r[1:]))
			r = r[4:]
			s.push(f32const(f))
		case 0x44: // f64const
			f := math.Float64frombits(binary.LittleEndian.Uint64(r[1:]))
			r = r[8:]
			s.push(f64const(f))
		case 0x45: // i32.eqz
			s.push(i32eqz{s.pop()})
		case 0x46: // i32.eq
			s.push(cmp{s.pop(), s.pop(), "=="})
		case 0x47: // i32.ne
			s.push(cmp{s.pop(), s.pop(), "!="})
		case 0x48: // i32.lts
			s.push(cmp{i32{s.pop()}, i32{s.pop()}, "<"})
		case 0x49: // i32.ltu
			s.push(cmp{s.pop(), s.pop(), "<"})
		case 0x4a: // i32.gts
			s.push(cmp{i32{s.pop()}, i32{s.pop()}, ">"})
		case 0x4b: // i32.gtu
			s.push(cmp{s.pop(), s.pop(), ">"})
		case 0x4c: // i32.les
			s.push(cmp{i32{s.pop()}, i32{s.pop()}, "<="})
		case 0x4d: // i32.leu
			s.push(cmp{s.pop(), s.pop(), "<="})
		case 0x4e: // i32.ges
			s.push(cmp{i32{s.pop()}, i32{s.pop()}, ">="})
		case 0x4f: // i32.geu
			s.push(cmp{s.pop(), s.pop(), ">="})
		case 0x5b: // f32.eq
			s.push(cmp{s.pop(), s.pop(), "=="})
		case 0x5c: // f32.ne
			s.push(cmp{s.pop(), s.pop(), "!="})
		case 0x5d: // f32.lt
			s.push(cmp{s.pop(), s.pop(), "<"})
		case 0x5e: // f32.gt
			s.push(cmp{s.pop(), s.pop(), ">"})
		case 0x5f: // f32.le
			s.push(cmp{s.pop(), s.pop(), "<="})
		case 0x60: // f32.ge
			s.push(cmp{s.pop(), s.pop(), ">="})
		case 0x61: // f64.eq
			s.push(cmp{s.pop(), s.pop(), "=="})
		case 0x62: // f64.ne
			s.push(cmp{s.pop(), s.pop(), "!="})
		case 0x63: // f64.lt
			s.push(cmp{s.pop(), s.pop(), "<"})
		case 0x64: // f64.gt
			s.push(cmp{s.pop(), s.pop(), ">"})
		case 0x65: // f64.le
			s.push(cmp{s.pop(), s.pop(), "<="})
		case 0x66: // f64.ge
			s.push(cmp{s.pop(), s.pop(), ">="})
		case 0x6a: // i32.add
			s.push(i32op2{s.pop(), s.pop(), "+"})
		case 0x6b: // i32.sub
			s.push(i32op2{s.pop(), s.pop(), "-"})
		case 0x6c: // i32.mul
			s.push(i32op2{s.pop(), s.pop(), "*"})
		case 0x6d: // i32.divs
			s.push(i32op2{i32{s.pop()}, i32{s.pop()}, "/"})
		case 0x6e: // i32.divu
			s.push(i32op2{s.pop(), s.pop(), "/"})
		case 0x6f: // i32.rems
			s.push(i32op2{i32{s.pop()}, i32{s.pop()}, "%"})
		case 0x70: // i32.remu
			s.push(i32op2{s.pop(), s.pop(), "%"})
		case 0x67: //i32.clz
			s.push(i32op1{s.pop(), "bits.LeadingZeros32"})
		case 0x68: //i32.ctz
			s.push(i32op1{s.pop(), "bits.TrailingZeros32"})
		case 0x69: //i32.popcnt
			s.push(i32op1{s.pop(), "bits.OnesCount32"})
		case 0x71: // i32.and
			s.push(i32op2{s.pop(), s.pop(), "&"})
		case 0x72: // i32.or
			s.push(i32op2{s.pop(), s.pop(), "|"})
		case 0x73: // i32.xor
			s.push(i32op2{s.pop(), s.pop(), "^"})
		case 0x74: // i32.shl
			s.push(i32op2{s.pop(), s.pop(), "<<"})
		case 0x75: // i32.shr_s
			s.push(i32op2{i32{s.pop()}, i32{s.pop()}, ">>"})
		case 0x76: // i32.shr_u
			s.push(i32op2{s.pop(), s.pop(), ">>"})
		case 0x77: //i32.rotl
			s.push(i32op2{ii{s.pop()}, s.pop(), "bits.RotateLeft32"})
		case 0x78: //i32.rotr
			s.push(i32op2{ii{i32op1{s.pop(), "-"}}, s.pop(), "bits.RotateLeft32"})
		case 0x79: //i64.clz
			s.push(i64op1{s.pop(), "bits.LeadingZeros64"})
		case 0x7a: //i64.ctz
			s.push(i64op1{s.pop(), "bits.TrailingZeros64"})
		case 0x7b: //i64.popcnt
			s.push(i64op1{s.pop(), "bits.OnesCount64"})
		case 0x7c: // i64.add
			s.push(i64op2{s.pop(), s.pop(), "+"})
		case 0x7d: // i64.sub
			s.push(i64op2{s.pop(), s.pop(), "-"})
		case 0x7e: // i64.mul
			s.push(i64op2{s.pop(), s.pop(), "*"})
		case 0x7f: // i64.divs
			s.push(i64op2{i32{s.pop()}, i32{s.pop()}, "/"})
		case 0x80: // i64.divu
			s.push(i64op2{s.pop(), s.pop(), "/"})
		case 0x81: // i64.rems
			s.push(i64op2{i64{s.pop()}, i64{s.pop()}, "%"})
		case 0x82: // i32.remu
			s.push(i64op2{s.pop(), s.pop(), "%"})
		case 0x83: // i64.and
			s.push(i64op2{s.pop(), s.pop(), "&"})
		case 0x84: // i64.or
			s.push(i64op2{s.pop(), s.pop(), "|"})
		case 0x85: // i64.xor
			s.push(i64op2{s.pop(), s.pop(), "^"})
		case 0x86: // i64.shl
			s.push(i64op2{s.pop(), s.pop(), "<<"})
		case 0x87: // i64.shr_s
			s.push(i64op2{i64{s.pop()}, i64{s.pop()}, ">>"})
		case 0x88: // i64.shr_u
			s.push(i64op2{s.pop(), s.pop(), ">>"})
		case 0x89: //i64.rotl
			s.push(i64op2{ii{s.pop()}, s.pop(), "bits.RotateLeft64"})
		case 0x8a: //i64.rotr
			s.push(i64op2{ii{i64op1{s.pop(), "-"}}, s.pop(), "bits.RotateLeft64"})
		case 0x8b: // f32abs
			float1(32, "math.Abs")
		case 0x8c: // f32neg
			float1(32, "-")
		case 0x8d: // f32ceil
			float1(32, "math.Ceil")
		case 0x8e: // f32floor
			float1(32, "math.Floor")
		case 0x8f: // f32trunc
			float1(32, "math.Trunc")
		case 0x90: // f32nearest
			float1(32, "math.Round")
		case 0x91: // f32sqrt
			float1(32, "math.Sqrt")
		case 0x92: // f32add
			float2(32, "+")
		case 0x93: // f32sub
			float2(32, "-")
		case 0x94: // f32mul
			float2(32, "*")
		case 0x95: // f32div
			float2(32, "/")
		case 0x96: // f32min
			float2(32, "math.Min")
		case 0x97: // f32max
			float2(32, "math.Max")
		case 0x98: // f32copysign
			float2(32, "math.Copysign")
		case 0x99: // f64abs
			float1(64, "math.Abs")
		case 0x9a: // f64neg
			float1(64, "-")
		case 0x9b: // f64ceil
			float1(64, "math.Ceil")
		case 0x9c: // f64floor
			float1(64, "math.Floor")
		case 0x9d: // f64trunc
			float1(64, "math.Trunc")
		case 0x9e: // f64nearest
			float1(64, "math.Round")
		case 0x9f: // f64sqrt
			float1(64, "math.Sqrt")
		case 0xa0: // f64add
			float2(64, "+")
		case 0xa1: // f64sub
			float2(64, "-")
		case 0xa2: // f64mul
			float2(64, "*")
		case 0xa3: // f64div
			float2(64, "/")
		case 0xa4: // f64min
			float2(64, "math.Min")
		case 0xa5: // f64max
			float2(64, "math.Max")
		case 0xa6: // f64copysign
			float2(64, "math.Copysign")
		case 0xa7: // i32wrap_i64
			s.push(u32{s.pop()})
		case 0xa8: // i32.trunc_f32_s
			s.push(u32{i32{fltf1{size: 32, x: s.pop(), op: "math.Trunc"}}})
		case 0xa9: // i32.trunc_f32_u
			s.push(u32{fltf1{size: 32, x: s.pop(), op: "math.Trunc"}})
		case 0xaa: // i32.trunc_f64_s
			s.push(u32{i32{fltf1{size: 64, x: s.pop(), op: "math.Trunc"}}})
		case 0xab: // i32.trunc_f64_u
			s.push(u32{fltf1{size: 64, x: s.pop(), op: "math.Trunc"}})
		case 0xac: // i64.extend_i32_s
			s.push(u64{i64{i32{s.pop()}}})
		case 0xad: // i64.extend_i32_u
			s.push(u64{s.pop()})
		case 0xae: // i64.trunc_f32_s
			s.push(u64{i64{fltf1{size: 64, x: s.pop(), op: "math.Trunc"}}})
		case 0xaf: // i64.trunc_f32_u
			s.push(u64{fltf1{size: 64, x: s.pop(), op: "math.Trunc"}})
		case 0xb0: // i64.trunc_f64_s
			s.push(u64{i64{fltf1{size: 64, x: s.pop(), op: "math.Trunc"}}})
		case 0xb1: // i64.trunc_f64_u
			s.push(u64{fltf1{size: 64, x: s.pop(), op: "math.Trunc"}})
		case 0xb2: // f32.convert_i32s
			s.push(f32{i32{s.pop()}})
		case 0xb3: // f32.convert_i32u
			s.push(f32{s.pop()})
		case 0xb4: // f32.convert_i64s
			s.push(f32{i64{s.pop()}})
		case 0xb5: // f32.convert_i64u
			s.push(f32{s.pop()})
		case 0xb6: // f32.demote_f64
			s.push(f32{s.pop()})
		case 0xb7: // f64.convert_i32u
			s.push(f64{s.pop()})
		case 0xb8: // f64.convert_i64s
			s.push(f64{i64{s.pop()}})
		case 0xba: // f64.convert_i64u
			s.push(f64{s.pop()})
		case 0xbb: // f64.promote_f32
			s.push(f64{s.pop()})
		case 0xbc: // i32.reinterpret_f32
			s.push(reinterp{s.pop(), "math.Float32bits"})
		case 0xbd: // i64.reinterpret_f64
			s.push(reinterp{s.pop(), "math.Float64bits"})
		case 0xbe: // f32.reinterpret_i32
			s.push(reinterp{s.pop(), "math.Float32frombits"})
		case 0xbf: // f64.reinterpret_i64
			s.push(reinterp{s.pop(), "math.Float64frombits"})
		case 0xc0: // i32.extend8_s
			s.push(u32{i32{i8{s.pop()}}})
		case 0xc1: // i32.extend16_s
			s.push(u32{i32{i16{s.pop()}}})
		case 0xc2: // i64.extend8_s
			s.push(u64{i64{i8{s.pop()}}})
		case 0xc3: // i64.extend16_s
			s.push(u64{i64{i16{s.pop()}}})
		case 0xc4: // i64.extend32_s
			s.push(u64{i64{i32{s.pop()}}})
		default:
			panic(fmt.Sprintf("unknown %x\n", b))
		}
		r = r[1:]
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
	return r[num-1:]
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
	return r[num-1:]
}
func immi64(r []byte, i *int64) []byte {
	const (
		int64Mask  int64 = 1 << 7
		int64Mask2       = ^int64Mask
		int64Mask3       = 1 << 6
		int64Mask4       = ^0
	)
	var shift int
	var b int64
	num := 0
	for shift < 64 {
		b = int64(r[num])
		num++
		*i |= (b & int64Mask2) << shift
		shift += 7
		if b&int64Mask == 0 {
			break
		}
	}

	if shift < 64 && (b&int64Mask3) == int64Mask3 {
		*i |= int64Mask4 << shift
	}
	return r[num-1:]
}

const head = `//+build ignore

package main
import (
	"encoding/binary"
	"math"
	"math/bits"
	"fmt"
)
var M []byte
var F []interface{}
func li32u32(addr, align, offset uint32) uint32 { return binary.LittleEndian.Uint32(M[addr+offset:]) }
func li32s8(addr, align, offset uint32) uint32  { return uint32(int32(M[addr+offset])) }
func li32u8(addr, align, offset uint32) uint32  { return uint32(M[addr+offset]) }
func li32s16(addr, align, offset uint32) uint32 { return uint32(binary.LittleEndian.Uint16(M[addr+offset:])) }
func li32u16(addr, align, offset uint32) uint32 { return uint32(int32(binary.LittleEndian.Uint16(M[addr+offset:]))) }

func li64u64(addr, align, offset uint32) uint64 { return binary.LittleEndian.Uint64(M[addr+offset:]) }
func li64s8(addr, align, offset uint32) uint64  { return uint64(int64(M[addr+offset])) }
func li64u8(addr, align, offset uint32) uint64  { return uint64(M[addr+offset]) }
func li64s16(addr, align, offset uint32) uint64 { return uint64(binary.LittleEndian.Uint16(M[addr+offset:])) }
func li64u16(addr, align, offset uint32) uint64 { return uint64(int32(binary.LittleEndian.Uint16(M[addr+offset:]))) }
func li64s32(addr, align, offset uint32) uint64 { return uint64(binary.LittleEndian.Uint32(M[addr+offset:])) }
func li64u32(addr, align, offset uint32) uint64 { return uint64(int32(binary.LittleEndian.Uint32(M[addr+offset:]))) }

func si32u8(value uint32, addr, align, offset uint32)  { M[addr+offset] = uint8(value) }
func si32u16(value uint32, addr, align, offset uint32) { binary.LittleEndian.PutUint16(M[addr+offset:], uint16(value)) }
func si32u32(value uint32, addr, align, offset uint32) { binary.LittleEndian.PutUint32(M[addr+offset:], value) }
func si64u8(value uint32, addr, align, offset uint32)  { M[addr+offset] = uint8(value) }
func si64u16(value uint32, addr, align, offset uint32) { binary.LittleEndian.PutUint16(M[addr+offset:], uint16(value)) }
func si64u32(value uint32, addr, align, offset uint32) { binary.LittleEndian.PutUint32(M[addr+offset:], uint32(value)) }
func si64u64(value uint64, addr, align, offset uint32) { binary.LittleEndian.PutUint64(M[addr+offset:], value) }

func sf32(value float32, addr, align, offset uint32) { si32u32(math.Float32bits(value), addr, align, offset) }
func sf64(value float64, addr, align, offset uint32) { si64u64(math.Float64bits(value), addr, align, offset) }
func lf32(addr, align, offset uint32) float32 { return math.Float32frombits(li32u32(addr, align, offset)) }
func lf64(addr, align, offset uint32) float64 { return math.Float64frombits(li64u64(addr, align, offset)) }

func ub(b bool) uint32 { if b { return 1 } else { return 0 } }
`
