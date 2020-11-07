package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"testing"
)

// Test wasm module generated from ../../k.w.

// K is implemented by a wasm interpreter (wagon/gasm)
type K interface {
	Memory() []byte
	Call(string, ...uint32) uint32
}

func ktest(newk func([]byte) K, t *testing.T) {
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
	m, tab, data := run(bytes.NewReader(mb))
	v := strings.Split(string(b), "\n")
	for i := range v {
		if len(v[i]) == 0 {
			fmt.Println("skip")
			return
		} else if v[i][0] == '/' {
			continue
		}
		vv := strings.Split(v[i], " /")
		if len(vv) != 2 {
			panic("test file")
		}
		in := strings.TrimRight(vv[0], " \t\r")
		exp := strings.TrimSpace(vv[1])
		fmt.Printf("%s /%s\n", in, exp)

		k := newk(m.wasm(tab, data))

		k.Call("ini", 16)
		s := k.Call("mk", uint32(1), uint32(len(in)))
		mem := k.Memory()
		copy(mem[s+8:], []byte(in))
		r := k.Call("kst", k.Call("val", s))
		n := k.Call("nn", r)
		mem = k.Memory()
		got := string(mem[r+8 : r+8+n])
		if exp != got {
			t.Fatalf("expected/got:\n%s\n%s\n", exp, got)
		}

		mem = k.Memory()
		k.Call("dx", r)
		k.Call("dx", get(mem, 132)) // kkey
		k.Call("dx", get(mem, 136)) // kval
		k.Call("dx", get(mem, 148)) // xyz
		if e := leak(mem); e != nil {
			panic(e)
		}
	}
}
func get(m []byte, a uint32) uint32   { return binary.LittleEndian.Uint32(m[a:]) }
func getf(m []byte, a uint32) float64 { return math.Float64frombits(binary.LittleEndian.Uint64(m[a:])) }
func mark(m []byte) { // mark bucket type within free blocks
	set := func(x, y uint32) { binary.LittleEndian.PutUint32(m[x:], y) }
	for t := uint32(4); t < 32; t++ {
		p := get(m, 4*t) // free pointer of type t
		for p != 0 {
			set(4+p, 0)
			set(8+p, t)
			p = get(m, p) // pointer to next free
		}
	}
}
func leak(m []byte) error {
	//dump(m, 0, 200)
	mark(m)
	p := uint32(256) // first data block
	for p < uint32(len(m)) {
		// a free block has refcount==0 at p+4 and bucket type at p+8 (after marking)
		if get(m, p+4) != 0 {
			//dump(m, 0, 200)
			return fmt.Errorf("non-free block at %d(%x)", p, p)
		}
		t := m[p+8]
		if t < 4 || t > 31 {
			return fmt.Errorf("illegal bucket type %d at %d(%x)", t, p, p)
		}
		dp := 1 << t
		p += uint32(dp)
	}
	return nil
}
