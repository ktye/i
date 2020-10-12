// +build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
)

func main() {
	f := "c:/k/ktye.github.io/ngn.wasm"
	if len(os.Args) == 2 {
		f = os.Args[1]
	}
	b, e := ioutil.ReadFile(f)
	fatal(e)
	m, e := wasm.ReadModule(bytes.NewReader(b), hostFuncs)
	fatal(e)
	w, e := exec.NewVM(m)
	fatal(e)
	ini := int64(m.Export.Entries["init"].Index)
	rep := int64(m.Export.Entries["rep"].Index)
	hep, ok := w.GetGlobal("__heap_base")
	if !ok {
		fatal(fmt.Errorf("no global __heap_base"))
	}
	B = hep
	_, e = w.ExecCode(ini)
	fatal(e)
	repl(func() error { _, e := w.ExecCode(rep); return e })
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
func repl(rep func() error) {
	s := bufio.NewScanner(os.Stdin)
	fmt.Printf(" ")
	for s.Scan() {
		S = s.Text() + "\n"
		rep()
		fmt.Printf(" ")
	}
}

var S string
var B uint64

func hostFuncs(name string) (*wasm.Module, error) {
	read := func(proc *exec.Process, fd, p, n int32) int32 {
		b := []byte(S)
		_, e := proc.WriteAt(b, int64(p))
		fatal(e)
		return int32(len(b))
	}
	writ := func(proc *exec.Process, fd, p, n int32) int32 {
		b := make([]byte, n)
		_, e := proc.ReadAt(b, int64(p))
		fatal(e)
		os.Stdout.Write(b)
		return n
	}
	now := func(proc *exec.Process) int64 { return 0 }
	O := func(proc *exec.Process, x int64) int64 { return 0 }
	munmap := func(proc *exec.Process, x, y int32) int32 { return 0 }
	exit := func(proc *exec.Process, x int32) { os.Exit(0) }
	mmap := func(proc *exec.Process, a, b, c, d, e, f int32) int32 {
		if B == 0 {
			panic("oom")
		}
		r := int32(B)
		B = 0
		return r
		return 0
	}
	m := wasm.NewModule()
	m.Types = &wasm.SectionTypes{
		Entries: []wasm.FunctionSig{
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI64}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64}},                                                                                                // O
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},                                                          // read,write
			{Form: 0, ParamTypes: []wasm.ValueType{}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64}},                                                                                                                 // now
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI64, wasm.ValueTypeI64}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64}},                                                                             // v0c, v1c
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}},                                                                             // munmap
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{}},                                                                                                                 //exit
			{Form: 0, ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32}, ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32}}, //mmap
		},
	}
	m.FunctionIndexSpace = []wasm.Function{
		{Sig: &m.Types.Entries[1], Host: reflect.ValueOf(read), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[1], Host: reflect.ValueOf(writ), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[2], Host: reflect.ValueOf(now), Body: &wasm.FunctionBody{}},
		{Sig: &m.Types.Entries[0], Host: reflect.ValueOf(O), Body: &wasm.FunctionBody{}},      // u0c
		{Sig: &m.Types.Entries[0], Host: reflect.ValueOf(O), Body: &wasm.FunctionBody{}},      // u1c
		{Sig: &m.Types.Entries[0], Host: reflect.ValueOf(O), Body: &wasm.FunctionBody{}},      // cmd
		{Sig: &m.Types.Entries[3], Host: reflect.ValueOf(O), Body: &wasm.FunctionBody{}},      // v0c
		{Sig: &m.Types.Entries[3], Host: reflect.ValueOf(O), Body: &wasm.FunctionBody{}},      // v1c
		{Sig: &m.Types.Entries[0], Host: reflect.ValueOf(O), Body: &wasm.FunctionBody{}},      // frk
		{Sig: &m.Types.Entries[4], Host: reflect.ValueOf(munmap), Body: &wasm.FunctionBody{}}, // munmap
		{Sig: &m.Types.Entries[5], Host: reflect.ValueOf(exit), Body: &wasm.FunctionBody{}},   // exit
		{Sig: &m.Types.Entries[6], Host: reflect.ValueOf(mmap), Body: &wasm.FunctionBody{}},   // mmap
		{Sig: &m.Types.Entries[0], Host: reflect.ValueOf(O), Body: &wasm.FunctionBody{}},      // hop
		{Sig: &m.Types.Entries[0], Host: reflect.ValueOf(O), Body: &wasm.FunctionBody{}},      // hcl
	}
	m.Export = &wasm.SectionExports{
		Entries: map[string]wasm.ExportEntry{
			"read":   {FieldStr: "read", Kind: wasm.ExternalFunction, Index: 0},
			"write":  {FieldStr: "write", Kind: wasm.ExternalFunction, Index: 1},
			"now":    {FieldStr: "now", Kind: wasm.ExternalFunction, Index: 2},
			"u0c":    {FieldStr: "u0c", Kind: wasm.ExternalFunction, Index: 3},
			"u1c":    {FieldStr: "u1c", Kind: wasm.ExternalFunction, Index: 4},
			"cmd":    {FieldStr: "cmd", Kind: wasm.ExternalFunction, Index: 5},
			"v0c":    {FieldStr: "v0c", Kind: wasm.ExternalFunction, Index: 6},
			"v1c":    {FieldStr: "v1c", Kind: wasm.ExternalFunction, Index: 7},
			"frk":    {FieldStr: "frk", Kind: wasm.ExternalFunction, Index: 8},
			"munmap": {FieldStr: "munmap", Kind: wasm.ExternalFunction, Index: 9},
			"exit":   {FieldStr: "exit", Kind: wasm.ExternalFunction, Index: 10},
			"mmap":   {FieldStr: "mmap", Kind: wasm.ExternalFunction, Index: 11},
			"hop":    {FieldStr: "hop", Kind: wasm.ExternalFunction, Index: 12},
			"hcl":    {FieldStr: "hcl", Kind: wasm.ExternalFunction, Index: 13},
		},
	}
	return m, nil
}
