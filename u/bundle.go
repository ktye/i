// +build ignore

// attach a default ui application to the binary
//  go run bundle.go u t.k
package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	a := os.Args
	if len(a) != 3 {
		panic("args")
	}
	tk, err := ioutil.ReadFile(a[2])
	p(err)
	f, err := os.OpenFile(a[1], os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		p(err)
	}
	defer f.Close()
	_, err = f.Write(tk)
	p(err)
	_, err = f.Write([]byte(fmt.Sprintf("%15dk\\ui", uint64(len(tk)))))
	p(err)
}
func p(err error) {
	if err != nil {
		panic(err)
	}
}
