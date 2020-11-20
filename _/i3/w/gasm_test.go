package main

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/mathetake/gasm/hostfunc"
	"github.com/mathetake/gasm/wasm"
)

func TestGasm(t *testing.T) {
	t.Skip() // wip
	newk := func(b []byte) K {
		m, e := wasm.DecodeModule(bytes.NewReader(b))
		if e != nil {
			panic(e)
		}
		vm, e := wasm.NewVM(m, gasmImport())
		if e != nil {
			panic(e)
		}
		return gk{m, vm}
	}
	ktest(newk, t)
}

type gk struct {
	m  *wasm.Module
	vm *wasm.VirtualMachine
}

func (k gk) Memory() []byte { return k.vm.Memory }
func (k gk) Call(s string, argv ...uint32) uint32 {
	args := make([]uint64, len(argv))
	for i := range args {
		args[i] = uint64(argv[i])
	}
	r, t, e := k.vm.ExecExportedFunction(s, args...)
	if e != nil {
		panic(e)
	}
	if len(t) != 1 || t[0] != wasm.ValueTypeI32 {
		panic("only a single i32 return value is supported")
	}
	return uint32(r[0])
}
func gasmImport() map[string]*wasm.Module {
	m := hostfunc.NewModuleBuilder()
	f := func(m *hostfunc.ModuleBuilder, name string, f reflect.Value) {
		m.MustSetFunction("ext", name, func(vm *wasm.VirtualMachine) reflect.Value { return f })
	}
	f(m, "sin", reflect.ValueOf(func(x float64) float64 { return math.Sin(x) }))
	f(m, "cos", reflect.ValueOf(func(x float64) float64 { return math.Cos(x) }))
	f(m, "log", reflect.ValueOf(func(x float64) float64 { return math.Log(x) }))
	f(m, "atan2", reflect.ValueOf(func(x, y float64) float64 { return math.Atan2(x, y) }))
	f(m, "hypot", reflect.ValueOf(func(x, y float64) float64 { return math.Hypot(x, y) }))
	f(m, "draw", reflect.ValueOf(func(x, y, z uint32) { panic("dummy-draw") }))
	f(m, "grow", reflect.ValueOf(func(x uint32) uint32 { panic("dummy-grow"); return x }))
	f(m, "printc", reflect.ValueOf(func(x, y uint32) { panic("dummy-printc") }))
	return m.Done()
}

func dump(m []byte, n uint32) {
	fmt.Printf("%.8x ", 0)
	for i := uint32(0); i < n; i++ {
		p := 4 * i
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
