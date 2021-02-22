package j

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"j/jgo"
	"os"
	"testing"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
)

func TestJ(t *testing.T) {
	ms := uint32(16)
	js := []jer{jj{}, jgo.J{}, newWagon(t)}
	all := func(js []jer, x uint32) {
		for _, j := range js {
			j.J(x)
		}
	}
	all(js, ms)
	b, err := ioutil.ReadFile("t")
	if err != nil {
		t.Fatal(err)
	}
	v := bytes.Split(b, []byte{10})
	for _, b := range v {
		if len(b) == 0 || b[0] == '(' { // (section)
			continue
		}
		if string(b) == `\` {
			break
		}
		i := bytes.IndexByte(b, '(')
		if i < 0 {
			t.Fatal("no (")
		}
		e := bytes.IndexByte(b, ')') // first) (other comments possible)
		if e < 0 {
			t.Fatal("no )")
		}
		in := b[:i]
		exp := string(b[i : 1+e])

		os.Stdout.Write(in)
		os.Stdout.WriteString(exp)
		for i, j := range js {
			runtest(t, j, in, exp)
			if i > 0 {
				memcompare(t, js[0], js[i])
			}
		}
		os.Stdout.WriteString("\n")
	}
}
func memcompare(t *testing.T, aj, bj jer) {
	a, b := aj.M(), bj.M()
	if len(a) != len(b) {
		t.Fatalf("memory#: %d != %d", len(a), len(b))
	}
	for i, u := range a {
		if u != b[i] {
			fmt.Println()
			Dump(aj, 100)
			Dump(bj, 100)
			t.Fatalf("mem differs at %d (%x): %x %x\n", i, i, u, b[i])
		}
	}
}
func runtest(t *testing.T, j jer, b []byte, exp string) {
	for _, c := range b {
		if j.J(uint32(c)) != 0 {
			t.Fatal("early value")
		}
	}
	r := j.J(10)
	if r == 0 {
		t.Fatal("zero")
	}
	m := j.M()
	s := SX(m, m[1])
	s = "(" + s[1:len(s)-1] + ")"
	if s == exp {
		os.Stdout.WriteString(" ok")
	} else {
		t.Fatalf("got %q\nexp %q\n", s, exp)
	}
	n := ln(m, m[1])
	for i := uint32(0); i < n; i++ {
		j.J('_')
	}
	j.J(10)
	Leak(j)
}

func newWagon(t *testing.T) wk {
	fatal := func(e error) {
		if e != nil {
			t.Fatal(e)
		}
	}
	b, e := ioutil.ReadFile("j.wasm")
	fatal(e)
	m, e := wasm.ReadModule(bytes.NewReader(b), nil)
	fatal(e)
	vm, e := exec.NewVM(m)
	fatal(e)
	return wk{m, vm}
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
