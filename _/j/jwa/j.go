package jwa

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
	"j/x"
	"reflect"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
)

//go:embed j.wasm
var jwasm []byte

type I = uint32

func New() interface {
	J(x I) I
	M() []I
} {
	m, e := wasm.ReadModule(bytes.NewReader(jwasm), hostFuncs)
	fatal(e)
	vm, e := exec.NewVM(m)
	fatal(e)
	return wk{m, vm}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
func hostFuncs(name string) (*wasm.Module, error) { // imported as module "ext"
	stk := func(proc *exec.Process, y uint32) {
		b := make([]byte, proc.MemSize())
		_, e := proc.ReadAt(b, 0)
		if e != nil {
			panic(e)
		}
		m := make([]uint32, len(b)/4)
		for i := range m {
			m[i] = binary.LittleEndian.Uint32(b[4*i:])
		}
		fmt.Println(x.X(m, m[1]))
	}
	draw := func(proc *exec.Process, x, y uint32) { fmt.Println("draw", x, y) }
	m := wasm.NewModule()
	m.Types = &wasm.SectionTypes{
		Entries: []wasm.FunctionSig{
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI32}, ReturnTypes: nil},
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: nil},
		},
	}
	m.FunctionIndexSpace = []wasm.Function{
		{Sig: &m.Types.Entries[0], Host: reflect.ValueOf(stk), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[1], Host: reflect.ValueOf(draw), Body: &wasm.FunctionBody{}},
	}
	m.Export = &wasm.SectionExports{
		Entries: map[string]wasm.ExportEntry{
			"stk":  {FieldStr: "stk", Kind: wasm.ExternalFunction, Index: 0},
			"draw": {FieldStr: "draw", Kind: wasm.ExternalFunction, Index: 1},
		},
	}
	return m, nil
}

type wk struct {
	m  *wasm.Module
	vm *exec.VM
}

func (k wk) J(a uint32) uint32 {
	x, ok := k.m.Export.Entries["j"]
	if !ok {
		panic("no j")
	}
	res, e := k.vm.ExecCode(int64(x.Index), uint64(a))
	if e != nil {
		panic(e)
	}
	return res.(uint32)
}
func (k wk) M() []uint32 {
	b := k.vm.Memory()
	r := make([]uint32, len(b)>>2)
	binary.Read(bytes.NewReader(b), binary.LittleEndian, &r)
	return r
}
