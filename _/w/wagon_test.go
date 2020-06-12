package main

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
)

func TestWagon(t *testing.T) {
	// t.Skip()
	newk := func(b []byte) K {
		m, e := wasm.ReadModule(bytes.NewReader(b), hostFuncs)
		if e != nil {
			t.Fatal(e)
		}
		vm, e := exec.NewVM(m)
		if e != nil {
			t.Fatal(e)
		}
		return wk{m, vm}
	}
	ktest(newk, t)
}

type wk struct {
	m  *wasm.Module
	vm *exec.VM
}

func (k wk) Memory() []byte { return k.vm.Memory() }
func (k wk) Call(s string, argv ...uint32) uint32 {
	m, vm := k.m, k.vm
	x, ok := m.Export.Entries[s]
	if !ok {
		panic(fmt.Errorf("function does not exist: %s", s))
	}
	fidx := m.Function.Types[x.Index]
	ftyp := m.Types.Entries[fidx]
	n := len(ftyp.ParamTypes)
	var e error
	var res interface{}
	if n != len(argv) {
		panic(fmt.Errorf("%s expects %d arguments (got %d)", s, n, len(argv)))
	}
	if n == 1 {
		res, e = vm.ExecCode(int64(x.Index), uint64(argv[0]))
	} else if n == 2 {
		res, e = vm.ExecCode(int64(x.Index), uint64(argv[0]), uint64(argv[1]))
	} else if n == 3 {
		res, e = vm.ExecCode(int64(x.Index), uint64(argv[0]), uint64(argv[1]), uint64(argv[2]))
	} else {
		panic(fmt.Errorf("%s expects %d arguments", s, n))
	}
	if e != nil {
		panic(e)
	}
	switch v := res.(type) {
	case nil:
		return 0
	case uint32:
		return v
	default:
		panic(fmt.Errorf("%s returns %T", s, res))
	}
}
func hostFuncs(name string) (*wasm.Module, error) { // imported as module "ext"
	sin := func(proc *exec.Process, x float64) float64 { return math.Sin(x) }
	cos := func(proc *exec.Process, x float64) float64 { return math.Cos(x) }
	log := func(proc *exec.Process, x float64) float64 { return math.Log(x) }
	atan2 := func(proc *exec.Process, x, y float64) float64 { return math.Atan2(x, y) }
	hypot := func(proc *exec.Process, x, y float64) float64 { return math.Hypot(x, y) }
	draw := func(proc *exec.Process, x, y, z uint32) { panic("dummy-draw") }
	grow := func(proc *exec.Process, x uint32) uint32 { panic("dummy-grow") }
	printc := func(proc *exec.Process, x, y uint32) {
		b := make([]byte, int(y))
		proc.ReadAt(b, int64(x))
		fmt.Printf("%s\n", string(b))
	}

	m := wasm.NewModule()
	m.Types = &wasm.SectionTypes{
		Entries: []wasm.FunctionSig{
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeF64}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeF64}},
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeF64, wasm.ValueTypeF64}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeF64}},
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: nil},
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: nil},
		},
	}
	m.FunctionIndexSpace = []wasm.Function{
		{Sig: &m.Types.Entries[0], Host: reflect.ValueOf(sin), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[0], Host: reflect.ValueOf(cos), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[0], Host: reflect.ValueOf(log), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[1], Host: reflect.ValueOf(atan2), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[1], Host: reflect.ValueOf(hypot), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[2], Host: reflect.ValueOf(draw), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[3], Host: reflect.ValueOf(grow), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[4], Host: reflect.ValueOf(printc), Body: &wasm.FunctionBody{}},
	}
	m.Export = &wasm.SectionExports{
		Entries: map[string]wasm.ExportEntry{
			"sin":    {FieldStr: "sin", Kind: wasm.ExternalFunction, Index: 0},
			"cos":    {FieldStr: "cos", Kind: wasm.ExternalFunction, Index: 1},
			"log":    {FieldStr: "log", Kind: wasm.ExternalFunction, Index: 2},
			"atan2":  {FieldStr: "atan2", Kind: wasm.ExternalFunction, Index: 3},
			"hypot":  {FieldStr: "hypot", Kind: wasm.ExternalFunction, Index: 4},
			"draw":   {FieldStr: "draw", Kind: wasm.ExternalFunction, Index: 5},
			"grow":   {FieldStr: "grow", Kind: wasm.ExternalFunction, Index: 6},
			"printc": {FieldStr: "printc", Kind: wasm.ExternalFunction, Index: 7},
		},
	}
	return m, nil
}
