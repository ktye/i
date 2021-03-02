package main

import (
	"bufio"
	"flag"
	"j"
	"j/jgo"
	"j/jwa"
	"os"
	_ "embed"
)

//go:embed j_.j
var lib []byte

func main() {
	var jj jer = jj{}
	var impl string
	var sz int
	flag.StringVar(&impl, "j", "go", "implementation: go|jgo|jwa")
	flag.IntVar(&sz, "sz", 16, "log2 mem size")
	flag.Parse()

	switch impl {
	case "go":
	case "jgo":
		jj = jgo.New()
	case "jwa":
		jj = jwa.New()
	default:
		panic("-j flag")
	}

	jj.J(uint32(sz))
	for _, b := range lib {
		jj.J(uint32(b))
	}
	jj.J(10)

	p()
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		b := []byte(s.Text())
		for _, u := range b {
			jj.J(uint32(u))
		}
		jj.J(10)
		p()
	}
}

func p() { os.Stdout.WriteString("j) ") }

type jer interface {
	J(x uint32) uint32
	M() []uint32
}

type jj struct{}

func (o jj) J(x uint32) uint32 { return j.J(x) }
func (o jj) M() []uint32       { return j.M }
