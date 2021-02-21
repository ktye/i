package j

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"j/jgo"
	"testing"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/wasm"
)

func TestJ(t *testing.T) {
	ms := uint32(16)
	jjj := jj{}
	jgo := jgo.J{}
	jwa := newWagon(t)
	jjj.J(ms)
	jgo.J(ms)
	jwa.J(ms)
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
		fmt.Println(string(b))
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

		runtest(t, jjj, in, exp)
		runtest(t, jgo, in, exp)
		runtest(t, jwa, in, exp)
	}
}
func runtest(t *testing.T, j jer, b []byte, exp string) {
	B := M
	defer func() { M = B }()

	for _, c := range b {
		if j.J(uint32(c)) != 0 {
			t.Fatal("early value")
		}
	}
	r := j.J(10)
	if r == 0 {
		t.Fatal("zero")
	}
	M = j.M()
	s := X(I(4))
	s = "(" + s[1:len(s)-1] + ")"
	if s != exp {
		t.Fatalf("got %q\nexp %q\n", r, exp)
	}
	Leak()
	cls()
}

type jj struct{}

func (o jj) J(x uint32) uint32 { return J(x) }
func (o jj) M() []uint32       { return M }

type jer interface {
	J(x uint32) uint32
	M() []uint32
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
