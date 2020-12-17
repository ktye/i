package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	var m string
	var start, size int
	var t, d bool
	flag.StringVar(&m, "m", "r5", "machine")
	flag.IntVar(&start, "start", 0, "start address")
	flag.IntVar(&size, "size", 32*1024, "memory(bytes)")
	flag.BoolVar(&t, "t", false, "text input (decimal uint32)")
	flag.BoolVar(&d, "d", false, "disassemble")
	flag.Parse()

	var b = make([]byte, size)
	var prog []byte
	var e error
	if flag.NArg() == 0 {
		prog, e = ioutil.ReadAll(os.Stdin)
	} else {
		prog, e = ioutil.ReadFile(flag.Arg(0))
	}
	if e == nil && t {
		var buf bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(prog))
		for s.Scan() {
			w, e := strconv.ParseInt(s.Text(), 10, 32)
			fatal(e)
			binary.Write(&buf, binary.LittleEndian, int32(w))
		}
		prog = buf.Bytes()
	}
	fatal(e)
	copy(b, prog)

	M := map[string]machine{
		"r5": &R5{},
	}
	if x := M[m]; x == nil {
		panic("unknown machine: " + m)
	} else {
		if d {
			x.Dump(b, uint32(start))
		} else {
			x.Start(b, uint32(start))
		}
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

type machine interface {
	Start([]byte, uint32)
	Dump([]byte, uint32)
}
