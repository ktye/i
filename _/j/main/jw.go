// +build ignore
package main

// go build jw.go

import (
	"fmt"
	"j"
	"j/jgo"
	"os"
	"strings"
)

func main() {
	jg := jgo.J{}
	jg.J(16)
	j.M = jg.M()
	// j.Dump(100)
	b := []byte(strings.Join(os.Args[1:], " ") + "\n")
	for _, c := range b {
		r := jg.J(uint32(c))
		if r != 0 {
			j.M = jg.M()
			fmt.Println(j.X(j.M[1]))
		}
	}
	j.M = jg.M()
	j.Leak()
}
