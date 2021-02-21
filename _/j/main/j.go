// +build ignore
package main

// go build j.go

import (
	"fmt"
	. "j"
	"os"
	"strings"
)

func main() {
	J(16)
	// Dump(100)
	b := []byte(strings.Join(os.Args[1:], " ") + "\n")
	for _, c := range b {
		r := J(uint32(c))
		if r != 0 {
			fmt.Println(X(M[1]))
		}
	}
	Leak()
}
