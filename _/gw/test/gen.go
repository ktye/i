//+build ignore

package main

// generate out_test.go from t.w

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	w := os.Stdout
	w.Write([]byte(head))
	b, e := ioutil.ReadFile("t.w")
	fatal(e)
	scn := bufio.NewScanner(bytes.NewReader(b))
	for scn.Scan() {
		s := scn.Text()
		ci := strings.Index(s, ":")
		if len(s) > 0 && s[0] == 'f' && ci > 0 {
			f := s[:ci]
			v := strings.Split(scn.Text(), "\t")
			if len(v) > 2 {
				arg := v[len(v)-2][1:]
				exp := v[len(v)-1]
				a := strings.Split(arg, " ")
				writeTest(f, a, exp, w)
			}
		}
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
func writeTest(f string, a []string, exp string, w io.Writer) {
	fmt.Fprintf(w, "func TestF%s(t *testing.T) {\n\tr := %s(%s)\n", f, f, strings.Join(a, ", "))
	fmt.Fprintf(w, "\t"+`if s := fmt.Sprintf("%%v", r); s != "%s" { t.Fatalf("%s: %s != %%s", s)`, exp, f, exp)
	fmt.Fprintf(w, "\n\t} else { fmt.Println(%q, r) }\n}\n", f)
}

const head = `package main
import (
	"testing"
	"fmt"
)

// imported functions
func f0(x uint32) uint32 { i := int32(x); return uint32(-i) }
func f1(x, y uint32) uint32 { return x+y }

`
