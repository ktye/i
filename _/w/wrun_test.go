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
	for i := range args {
		if u, e := strconv.ParseUint(args[i], 10, 64); e == nil {
			stack = append(stack, u)
		} else if args[i] == "dump" {
			dump(vm.Memory(), k(stack[len(stack)-2]), k(stack[len(stack)-1]))
			stack = stack[:len(stack)-2]
		} else {
			if e := call(args[i]); e != nil {
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
	if t != 2 {
		panic("nyi: ~I")
	}
	r := make([]s, n)
	for i := range r {
		r[i] = strconv.Itoa(int(get(m, 8+4*k(i)+a)))
	}
	return strings.Join(r, " ")
}
func get(m []byte, a k) k { return binary.LittleEndian.Uint32(m[a:]) }
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
		t := get(m, p+8)
		if t < 4 || t > 31 {
			return fmt.Errorf("illegal bucket type %d at %d(%x)", t, p, p)
		}
		dp := 1 << t
		p += k(dp)
	}
	return nil
}
