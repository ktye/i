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
		m, data := run(bytes.NewReader(mb))
		if e = runWagon(m.wasm(data), strings.Fields(in), exp); e != nil {
			t.Fatalf("%d: %s", i+1, e)
		}
	}
}

func runWagon(b []byte, args []string, exp string) error {
	fmt.Println(args, exp)
	h := []string{"16", "ini"}
	args = append(h, args...)
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

	var stack []uint64
	call := func(s string) error {
		x, ok := m.Export.Entries[s]
		if !ok {
			return fmt.Errorf("unknown func: " + s)
		}
		fidx := m.Function.Types[x.Index]
		ftyp := m.Types.Entries[fidx]
		n := len(ftyp.ParamTypes)
		pop := make([]uint64, n)
		copy(pop, stack[len(stack)-n:])
		stack = stack[:len(stack)-n]
		res, e := vm.ExecCode(int64(x.Index), pop...)
		if e != nil {
			return e
		}
		if res != nil {
			stack = append(stack, u64(res))
			if trace {
				fmt.Printf("%s %v: %v(%x)\n", s, pop, res, res)
			}
		} else if trace {
			fmt.Printf("%s %v: nil\n", s, pop)
		}
		return nil
	}
	pushVector := func(s string) bool {
		if idx := strings.Index(s, ","); idx == -1 {
			return false
		}
		s = strings.Trim(s, ",")
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
				return false
			}
		}
		m := vm.Memory()
		if f == -1 {
			stack = append(stack, 2, uint64(n))
			if e = call("mk"); e != nil {
				return false
			}
			x := uint32(stack[len(stack)-1])
			for i := uint32(0); i < n; i++ {
				binary.LittleEndian.PutUint32(m[x+8+i*4:], uint32(iv[i]))
			}
		} else {
			stack = append(stack, 3, uint64(n))
			if e = call("mk"); e != nil {
				return false
			}
			x := uint32(stack[len(stack)-1])
			for i := uint32(0); i < n; i++ {
				binary.LittleEndian.PutUint64(m[x+8+i*8:], math.Float64bits(fv[i]))
			}
		}
		return true
	}
	for _, s := range args {
		if pushVector(s) {
			continue
		}
		if u, e := strconv.ParseUint(s, 10, 64); e == nil {
			stack = append(stack, u)
		} else if s == "dump" {
			dump(vm.Memory(), k(stack[len(stack)-2]), k(stack[len(stack)-1]))
			stack = stack[:len(stack)-2]
		} else if strings.HasPrefix(s, `"`) {
			s = strings.Trim(s, `"`)
			b := []c(s)
			stack = append(stack, 1, uint64(len(b)))
			if e := call("mk"); e != nil {
				return e
			}
			m := vm.Memory()
			p := stack[len(stack)-1]
			for i := 0; i < len(b); i++ {
				m[8+int(p)+i] = b[i]
			}
		} else {
			if e := call(s); e != nil {
				return e
			}
		}
	}
	if len(stack) != 2 { // [16, result]
		return fmt.Errorf("stack size")
	}
	// compare result
	got := kst(k(stack[1]), vm.Memory())
	if got != exp {
		return fmt.Errorf("expected/got:\n%s\n%s", exp, got)
	}
	// free result and check for memory leaks
	if e := call("dx"); e != nil {
		return e
	}
	if e := leak(vm.Memory()); e != nil {
		return e
	}
	return nil
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
func dump(M []byte, a, n k) {
	fmt.Printf("%.8x ", a)
	for i := k(0); i < n; i++ {
		p := a + 4*i
		x := get(M, p)
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

func kst(a k, m []byte) s {
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
	default:
		panic("nyi: kst t~CI")
	}
	r := make([]s, n)
	for i := range r {
		r[i] = f(i)
	}
	return tof(strings.Join(r, " "))
}
func get(m []byte, a k) k        { return binary.LittleEndian.Uint32(m[a:]) }
func getf(m []byte, a k) float64 { return math.Float64frombits(binary.LittleEndian.Uint64(m[a:])) }
func mark(m []byte) { // mark bucket type within free blocks
	for t := k(4); t < 32; t++ {
		p := get(m, 4*t) // free pointer of type t
		for p != 0 {
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
