package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		scn := bufio.NewScanner(os.Stdin)
		for scn.Scan() {
			scan([]byte(scn.Text()))
		}
	} else {
		b, e := ioutil.ReadFile(os.Args[1])
		fatal(e)
		scan(b)
	}
}
func scan(b []byte) {
	a, sp, e := token(b)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(a)
	fmt.Println(sp)
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
