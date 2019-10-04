// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// generate file t.go
func main() {
	t, err := os.Create("t.go")
	p(err)
	defer t.Close()
	tk, err := ioutil.ReadFile("t.k")
	p(err)
	fmt.Fprintf(t, "package main\n\nconst tk = %q\n", string(tk))
}
func p(err error) {
	if err != nil {
		panic(err)
	}
}
