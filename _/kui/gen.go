// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// generate h.go (k.w source and help)
// append date to k.go
func main() {
	b, e := ioutil.ReadFile("../../k.w")
	fatal(e)

	i := bytes.Index(b, []byte("\n\\\n"))
	if i == -1 {
		panic("k.w: cannot find \\")
	}

	src := b[:i]
	hlp := b[i+3:]

	var u bytes.Buffer
	fmt.Fprintf(&u, "package main\n\nconst source = %q\nconst help = %q\n", src, hlp)
	fatal(ioutil.WriteFile("h.go", u.Bytes(), 644))

	s := time.Now().Format("2006.01.02`\n")
	f, e := os.OpenFile("k.go", os.O_APPEND|os.O_WRONLY, 0644)
	fatal(e)
	defer f.Close()
	f.WriteString(s)
}

func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
