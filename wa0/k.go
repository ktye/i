package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"

	_ "embed"
)

//go:embed k.wasm
var K []byte

var Stdout io.Writer
var Stdin io.Reader
var bak []byte
var filebuf []byte

func main() {
	t0 := time.Now()
	var t1, t2 time.Duration

	ctx := context.Background()
	Stdout = os.Stdout
	Stdin = os.Stdin

	var r wazero.Runtime
	if len(os.Args) > 0 && os.Args[1] == "-i" { // interpreter
		copy(os.Args[1:], os.Args[2:])
		os.Args = os.Args[:len(os.Args)-1]
		r = wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfigInterpreter().WithWasmCore2())
	} else { // jit
		r = wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfig().WithWasmCore2())
	}
	defer r.Close(ctx)

	// ktye/wg system interface (https://github.com/ktye/wg/blob/master/module/system.go): Exit,Args,Arg,Read,Write,ReadIn,Native
	var m api.Module
	get := func(addr, n int32) (r []byte) { r, _ = m.Memory().Read(ctx, uint32(addr), uint32(n)); return r }
	set := func(addr int32, b []byte) { m.Memory().Write(ctx, uint32(addr), b) }
	arg := func(i, r int32) int32 {
		if i >= int32(len(os.Args)) {
			panic("arg")
		}
		if r == 0 {
			return int32(len(os.Args[i]))
		}
		set(r, []byte(os.Args[i]))
		return 0
	}
	read := func(file, nfile, dst int32) int32 {
		if dst != 0 {
			set(dst, filebuf)
			return 0
		}
		name := string(get(file, nfile))
		b, e := os.ReadFile(name)
		if e != nil {
			return -1
		}
		filebuf = b
		return int32(len(filebuf))
	}
	write := func(file, nfile, src, n int32) int32 {
		b := get(src, n)
		if nfile == 0 {
			Stdout.Write(b)
			return 0
		}
		name := string(get(file, nfile))
		e := os.WriteFile(name, b, 0644)
		if e != nil {
			return -1
		}
		return 0
	}
	readin := func(dst, n int32) int32 {
		b := make([]byte, n)
		nr, e := Stdin.Read(b)
		if e != nil {
			return 0
		}
		if nr > 0 && b[nr-1] == 10 {
			nr -= 1
		}
		set(dst, b)
		return int32(nr)
	}
	exit := func(x int32) {
		t2 = time.Since(t0)
		fmt.Fprintf(os.Stderr, "compile %v, run %v, total %v\n", t1, t2-t1, t2)
		os.Exit(int(x))
	}

	_, e := r.NewModuleBuilder("env").
		ExportFunction("Exit", exit).
		ExportFunction("Args", func() int32 { return int32(len(os.Args)) }).
		ExportFunction("Arg", arg).
		ExportFunction("Read", read).
		ExportFunction("Write", write).
		ExportFunction("ReadIn", readin).
		ExportFunction("Native", func(x, y int64) int64 { panic("nyi"); return x }).
		Instantiate(ctx, r)
	fatal(e)

	m, e = r.InstantiateModuleFromBinary(ctx, K)
	fatal(e)

	t1 = time.Since(t0)

	call := func(s string, args ...uint64) []uint64 {
		r, e := m.ExportedFunction(s).Call(ctx, args...)
		fatal(e)
		return r
	}
	ecall := func(s string, args ...uint64) error { _, e := m.ExportedFunction(s).Call(ctx, args...); return e }
	store := func() {
		b, _ := m.Memory().Read(ctx, 0, m.Memory().Size(ctx))
		if len(b) > len(bak) {
			bak = make([]byte, len(b))
		}
		copy(bak, b)
	}
	restore := func() { m.Memory().Write(ctx, 0, bak) }

	// same as _start (../system.go: func main), but without internal try/catch
	var x []uint64
	call("kinit")
	call("doargs")
	x = call("Ku", 2932601077199979)
	call("write", x[0])

	store()
	for {
		x = call("Ku", 32)
		call("write", x[0]) // space
		x = call("read")
		e = ecall("repl", x[0])
		if e != nil {
			fmt.Println("restore")
			restore()
		} else {
			store()
		}
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
