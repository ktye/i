package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"

	_ "embed"
)

//go:embed k.wasm
var K []byte

func main() {

	ctx := context.Background()

	//r := wazero.NewRuntime()
	//see: https://github.com/tetratelabs/wazero/blob/main/config.go
	r := wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfig().WithFeatureMultiValue(true).WithFeatureBulkMemoryOperations(true).WithFeatureSIMD(true))
	defer r.Close(ctx)

	_, e := r.NewModuleBuilder("env").
		ExportFunction("Exit", func(x int32) { os.Exit(int(x)) }).
		ExportFunction("Args", func() int32 { fmt.Printf("Args: %d\n", len(os.Args)); return int32(len(os.Args)) }).
		ExportFunction("Arg", func(i, r int32) int32 { panic("nyi"); return 0 /*nyi*/ }).
		ExportFunction("Read", func(file, nfile, dst int32) int32 { panic("nyi") }).
		ExportFunction("Write", func(file, nfile, src, n int32) int32 { panic("nyi") }).
		ExportFunction("ReadIn", func(dst, n int32) int32 { panic("nyi") }).
		ExportFunction("Native", func(x, y int64) int64 { panic("nyi"); return x }).
		Instantiate(ctx, r)
	fatal(e)

	//module, err := r.InstantiateModuleFromBinary(ctx, K)
	_, e = r.InstantiateModuleFromBinary(ctx, K)
	fatal(e)

	/*

		// fmt.Println(module.ExportedFunction("fac").Call(ctx, 7))
		call := func(s string, args ...uint64) ([]uint64, error) {
			return module.ExportedFunction(s).Call(ctx, args...)
		}
		store := func(addr int32, data []byte) { module.Memory().Write(ctx, uint32(addr), data) }
		Kc := func(b []byte) uint64 {
			r, e := call("mk", uint64(18), uint64(len(b))) // Ct
			fatal(e)
			store(int32(r[0]), b)
			return uint64(r[0])
		}


		call("kinit")
		call("doargs")

		// call("store")

		// repl
		s := bufio.NewScanner(os.Stdin)
		fmt.Printf(" ")
		for s.Scan() {
			t := strings.TrimSpace(s.Text())
			if t == "" {
				continue
			}

			call("try", Kc([]byte(t)))
			fmt.Printf(" ")

		}
	*/
}

func fatal(e error) {
	if e != nil {
		panic(e)
		os.Exit(1)
	}
}
