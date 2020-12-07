// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic("args")
	}
	name := os.Args[1]
	b, e := ioutil.ReadAll(os.Stdin)
	if e != nil {
		panic(e)
	}
	fmt.Printf("const %s = %q\n", name, string(b))
}
