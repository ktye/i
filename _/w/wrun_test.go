package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/validate"
	"github.com/go-interpreter/wagon/wasm"
)

const trace = false

func TestWagon(t *testing.T) {
	if broken {
		t.Skip()
	}
	b, e := ioutil.ReadFile("t")
	if e != nil {
		t.Fatal(e)
	}
	mb, e := ioutil.ReadFile("../../k.w")
	if e != nil {
		t.Fatal(e)
	}
	v := strings.Split(strings.TrimSpace(string(b)), "\n")
	for i := range v {
		if len(v[i]) == 0 || v[i][0] == '/' {
			continue
		}
		vv := strings.Split(v[i], " /")
		if len(vv) != 2 {
			panic("test file")
		}
		in := strings.TrimSpace(vv[0])
		exp := strings.TrimSpace(vv[1])
		m, tab, data := run(bytes.NewReader(mb))
		if e = runWagon(tab, m.wasm(tab, data), strings.Fields(in), exp); e != nil {
			t.Fatalf("%d: %s", i+1, e)
		}
	}
}

type K struct {
	m     *wasm.Module
	vm    *exec.VM
	stack []uint32
}

func runWagon(tab []segment, b []byte, args []string, exp string) error {
	fmt.Println(args, exp)
	m, e := wasm.ReadModule(bytes.NewReader(b), nil)
	if e != nil {
		return e
	}
	if e := validate.VerifyModule(m); e != nil {
		return e
	}
	vm, e := exec.NewVM(m)
	if e != nil {
		return e
	}
	if trace {
		fmt.Println("memory", len(vm.Memory()))
	}
	K := K{m: m, vm: vm, stack: []uint32{16}}
	K.call("ini")
	K.pop()
	for _, s := range args {
		if o, e := K.call(s); o == true {
			if e != nil {
				return e
			}
			continue
		}
		if strings.HasSuffix(s, "dump") {
			if n, e := strconv.Atoi(strings.TrimPrefix(s, "dump")); e == nil {
				K.dump(0, uint32(n))
			} else {
				K.dump(0, 100)
			}
			continue
		}
		K.push(K.parseVector(s))
	}
	if len(K.stack) != 1 {
		return fmt.Errorf("stack #" + strconv.Itoa(len(K.stack)))
	}
	// compare result
	r := K.pop()
	got := K.kst(r)
	if got != exp {
		return fmt.Errorf("expected/got:\n%s\n%s", exp, got)
	}
	// free result and check for memory leaks
	K.push(r)
	if _, e := K.call("dx"); e != nil {
		return e
	}
	if e := leak(K.vm.Memory()); e != nil {
		return e
	}
	return nil
}
func (K *K) push(x ...uint32) {
	K.stack = append(K.stack, x...)
}
func (K *K) pop() (r uint32) {
	r = K.stack[len(K.stack)-1]
	K.stack = K.stack[:len(K.stack)-1]
	return r
}
func (K *K) call(s string) (bool, error) {
	m, vm := K.m, K.vm
	x, ok := m.Export.Entries[s]
	if !ok {
		return false, fmt.Errorf("function does not exist: %s", s)
	}
	fidx := m.Function.Types[x.Index]
	ftyp := m.Types.Entries[fidx]
	n := len(ftyp.ParamTypes)
	var e error
	var res interface{}
	if n == 1 {
		res, e = vm.ExecCode(int64(x.Index), uint64(K.pop()))
	} else if n == 2 {
		y := K.pop()
		res, e = vm.ExecCode(int64(x.Index), uint64(K.pop()), uint64(y))
	} else {
		return true, fmt.Errorf("%s expects %d arguments", s, n)
	}
	if e != nil {
		return true, e
	}
	if res != nil {
		if u, o := res.(uint32); o {
			K.push(uint32(u))
		} else {
			return true, fmt.Errorf("%s: result type %T", s, res)
		}
		if trace {
			r := K.stack[len(K.stack)-1]
			fmt.Printf("%s: %v(%x)\n", s, r, r)
		}
	} else if trace {
		fmt.Printf("%s: nil\n", s)
	}
	return true, nil
}
func (K *K) mk(t, n uint32) uint32 {
	K.push(t, n)
	if _, e := K.call("mk"); e != nil {
		panic(e)
	}
	return K.pop()
}
func (K *K) parseVector(s string) uint32 {
	m := K.vm.Memory()
	if len(s) > 0 && s[0] == '"' {
		s = strings.Trim(s, `"`)
		b := []c(s)
		p := K.mk(1, uint32(len(b)))
		for i := 0; i < len(b); i++ {
			m[8+int(p)+i] = b[i]
		}
		return p
	}
	if len(s) > 0 && s[0] == '`' {
		v := strings.Split(s[1:], "`")
		sn := uint32(len(v))
		sv := K.mk(4, sn)
		for i := uint32(0); i < sn; i++ {
			b := v[i]
			rn := uint32(len(b))
			r := K.mk(1, rn)
			for k := uint32(0); k < rn; k++ {
				m[8+r+k] = b[k]
			}
			binary.LittleEndian.PutUint32(m[sv+8+4*i:], uint32(r)) // sI(sv+8+4*i, r)
		}
		return sv
	}
	if len(s) > 0 && s[0] == '(' {
		return K.parseList(s[1:])
	}
	f := strings.Index(s, ".")
	v := strings.Split(s, ",")
	n := uint32(len(v))
	iv := make([]int64, n)
	fv := make([]float64, n)
	var e error
	for i, s := range v {
		if f == -1 {
			iv[i], e = strconv.ParseInt(s, 10, 32)
		} else {
			fv[i], e = strconv.ParseFloat(s, 64)
		}
		if e != nil {
			panic(fmt.Errorf("parse: %s", s))
		}
	}
	if f == -1 {
		x := K.mk(2, n)
		for i := uint32(0); i < n; i++ {
			binary.LittleEndian.PutUint32(m[x+8+i*4:], uint32(iv[i]))
		}
		return x
	} else {
		x := K.mk(3, n)
		for i := uint32(0); i < n; i++ {
			binary.LittleEndian.PutUint64(m[x+8+i*8:], math.Float64bits(fv[i]))
		}
		return x
	}
}
func (K *K) parseList(s string) uint32 {
	if len(s) == 0 || s[len(s)-1] != ')' {
		panic("parse list")
	} else if len(s) == 1 {
		return K.mk(5, 0)
	}
	r := make([]uint32, 0)
	s = s[:len(s)-1]
	l, a := 0, 0
	for i, c := range s {
		if c == '(' {
			l++
		} else if c == ')' {
			l--
			if l < 0 {
				panic(")")
			}
		} else if l == 0 && c == ';' {
			r = append(r, K.parseVector(s[a:i]))
			a = i + 1
		}
	}
	r = append(r, K.parseVector(s[a:]))
	x := K.mk(5, uint32(len(r)))
	m := K.vm.Memory()
	for k := range r {
		binary.LittleEndian.PutUint32(m[8+x+4*uint32(k):], r[k])
	}
	return x
}
func u64(v interface{}) uint64 {
	switch x := v.(type) {
	case uint32:
		return uint64(x)
	case uint64:
		return x
	case float64:
		return math.Float64bits(x)
	default:
		panic(x)
	}
}
func (K *K) dump(a, n k) {
	m := K.vm.Memory()
	fmt.Printf("%.8x ", a)
	for i := k(0); i < n; i++ {
		p := a + 4*i
		x := get(m, p)
		fmt.Printf(" %.8x", x)
		if i > 0 && (i+1)%8 == 0 {
			fmt.Printf("\n%.8x ", p+4)
		} else if i > 0 && (i+1)%4 == 0 {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
}

type k = uint32

func (K *K) kst(a k) s {
	m := K.vm.Memory()
	x := get(m, a)
	t, n := x>>29, x&536870911
	var f func(i int) s
	var tof func(s) s = func(s s) s { return s }
	istr := func(i int) s {
		if n := int32(get(m, 8+4*k(i)+a)); n == -2147483648 {
			return "0N"
		} else {
			return strconv.Itoa(int(n))
		}
	}
	fstr := func(i int) s {
		if f := getf(m, a+8+8*k(i)); math.IsNaN(f) {
			return "0n"
		} else {
			return strconv.FormatFloat(f, 'g', -1, 64)
		}
	}
	sstr := func(i int) s {
		r := get(m, a+8+4*k(i))
		rn := get(m, r) & 536870911
		return string(m[r+8 : r+8+rn])
	}
	sep := " "
	switch t {
	case 1:
		return `"` + string(m[a+8:a+8+n]) + `"`
	case 2:
		f = istr
	case 3:
		f = fstr
		tof = func(s s) s {
			if strings.Index(s, ".") == -1 {
				return s + "f"
			}
			return s
		}
	case 4:
		f = sstr
		sep = "`"
		tof = func(s s) s { return "`" + s }
	case 5:
		f = func(i int) s { return K.kst(get(m, 8+4*uint32(i)+a)) }
		sep = ";"
		tof = func(s s) s { return "(" + s + ")" }
	default:
		panic(fmt.Sprintf("nyi: kst: t=%d", t))
	}
	r := make([]s, n)
	for i := range r {
		r[i] = f(i)
	}
	return tof(strings.Join(r, sep))
}
func get(m []byte, a k) k        { return binary.LittleEndian.Uint32(m[a:]) }
func getf(m []byte, a k) float64 { return math.Float64frombits(binary.LittleEndian.Uint64(m[a:])) }
func mark(m []byte) { // mark bucket type within free blocks
	for t := k(4); t < 32; t++ {
		p := get(m, 4*t) // free pointer of type t
		for p != 0 {
			m[4+p] = 0
			m[8+p] = c(t)
			p = get(m, p) // pointer to next free
		}
	}
}
func leak(m []byte) error {
	mark(m)
	p := k(256) // first data block
	for p < k(len(m)) {
		// a free block has refcount==0 at p+4 and bucket type at p+8 (after marking)
		if get(m, p+4) != 0 {
			return fmt.Errorf("non-free block at %d(%x)", p, p)
		}
		t := m[p+8]
		if t < 4 || t > 31 {
			return fmt.Errorf("illegal bucket type %d at %d(%x)", t, p, p)
		}
		dp := 1 << t
		p += k(dp)
	}
	return nil
}
