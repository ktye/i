package main

import (
	"fmt"
	"j"
	"os"
	"strings"
)

func main() {
	b := []byte(strings.Join(os.Args[1:], " ") + "\n")
	for _, c := range b {
		r := j.Step(uint32(c))
		if r != 0 {
			fmt.Println(j.XX(r))
			j.Dx(r)
		}
	}
	j.Leak()
}
