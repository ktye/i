package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"testing"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/validate"
	"github.com/go-interpreter/wagon/wasm"
)

func TestWagon(t *testing.T) {
	if broken {
		t.Skip()
	}
	args := []string{
		"0", "200", "dump",
		"5", "mki",
		"0", "500", "dump",
	}
	b, e := ioutil.ReadFile("../../k.w")
	if e != nil {
		t.Fatal(e)
	}
	m, data := run(bytes.NewReader(b))
	runWagon(m.wasm(data), args)
}

func runWagon(b []byte, args []string) {
	h := []string{"16", "ini"}
	args = append(h, args...)
	m, e := wasm.ReadModule(bytes.NewReader(b), nil)
	if e != nil {
		panic(e)
	}
	if e := validate.VerifyModule(m); e != nil {
		panic(e)
	}
	vm, e := exec.NewVM(m)
	if e != nil {
		panic(e)
	}
	fmt.Println("memory", len(vm.Memory()))

	var stack []uint64
	for i := range args {
		if u, e := strconv.ParseUint(args[i], 10, 64); e == nil {
			stack = append(stack, u)
		} else if args[i] == "dump" {
			dump(vm.Memory(), stack[len(stack)-2], stack[len(stack)-1])
			stack = stack[:len(stack)-2]
		} else {
			x, ok := m.Export.Entries[args[i]]
			if !ok {
				panic("unknown func: " + args[i])
			}
			fidx := m.Function.Types[x.Index]
			ftyp := m.Types.Entries[fidx]
			n := len(ftyp.ParamTypes)
			pop := make([]uint64, n)
			copy(pop, stack[len(stack)-n:])
			stack = stack[:len(stack)-n]
			res, e := vm.ExecCode(int64(x.Index), pop...)
			if e != nil {
				panic(e)
			}
			stack = append(stack, u64(res))
			fmt.Printf("%s %v: %v(%x)\n", args[i], pop, res, res)
		}
	}
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
func dump(M []byte, a, n uint64) {
	fmt.Printf("%.8x  ", 0)
	for i, b := range M[a : a+n] {
		hi, lo := hxb(b)
		fmt.Printf("%c%c", hi, lo)
		if i > 0 && (i+1)%32 == 0 {
			fmt.Printf("\n%.8x  ", i+1)
		} else if i > 0 && (i+1)%16 == 0 {
			fmt.Printf("  ")
		} else if i > 0 && (i+1)%4 == 0 {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
}

/*

	for name, e := range m.Export.Entries {
		i := int64(e.Index)
		fidx := m.Function.Types[int(i)]
		ftype := m.Types.Entries[int(fidx)]
		switch len(ftype.ReturnTypes) {
		case 1:
			fmt.Fprintf(w, "%s() %s => ", name, ftype.ReturnTypes[0])
		case 0:
			fmt.Fprintf(w, "%s() => ", name)
		default:
			log.Printf("running exported functions with more than one return value is not supported")
			continue
		}



	for name, e := range m.Export.Entries {
		if name == "ini" {
			if o, e := vm.ExecCode(int64(e.Index), 16); e != nil {
				panic(e)
			} else {
				fmt.Println("ini 16:", o)
			}
		}
	}
}
*/
