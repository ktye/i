package main

import (
	"bytes"
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

/*
func hostFuncs(name string) (*wasm.Module, error) { // imported as module "ext"
	sin := func(proc *exec.Process, x float64) float64 { return math.Sin(x) }
	cos := func(proc *exec.Process, x float64) float64 { return math.Cos(x) }
	log := func(proc *exec.Process, x float64) float64 { return math.Log(x) }
	atan2 := func(proc *exec.Process, x, y float64) float64 { return math.Atan2(x, y) }
	hypot := func(proc *exec.Process, x, y float64) float64 { return math.Hypot(x, y) }
	draw := func(proc *exec.Process, x, y, z uint32) uint32 { return 0 }
	grow := func(proc *exec.Process, x uint32) uint32 { return x } // not implemented for wrun_test
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
*/
