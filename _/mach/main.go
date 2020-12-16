package main

import (
	"flag"
	"io/ioutil"
	"os"
)

func main() {
	var m string
	var start, size int
	flag.StringVar(&m, "m", "r5", "machine")
	flag.IntVar(&start, "start", 0, "start address")
	flag.IntVar(&size, "size", 32*1024, "memory(bytes)")
	flag.Parse()

	var b = make([]byte, size)
	var prog []byte
	var e error
	if len(flag.Args) == 0 {
		prog, e = ioutil.ReadAll(os.Stdin)
	} else {
		prog, e = ioutil.ReadFile(flag.Args[0])
	}
	fatal(e)
	copy(b, prog)
	switch m {
	case "r5":
		r5.Start(b, uint32(start))
	default:
		panic("unknown machine: " + m)
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
