package main

import (
	"fmt"
	"os"
	"strings"
)

// C:";*aq/:"
// A:-127#C?"*q*****;;****a*aaaaaaaaaa:;*****aaaaaaaaaaaaaaaaaaaaaaaaaa;/;**aaaaaaaaaaaaaaaaaaaaaaaaaaa;*;*"
// M:(C,"br")?m:1_'(";;*aq**";"*;*aq*:";"a;*bq**";"qrrr:/r";"/rrrrrr";":;*aq**";"b;*bq**";"rrrr:/r")
// t:&5>M\A@
// x:"ab+:4"
// t x
// 0 2 4

func main() {
	A, M := AM()
	s := "ab+:4"
	if len(os.Args) > 1 {
		s = os.Args[1]
	}
	fmt.Printf("%s %q\n", s, tok(A, M, s))
}
func tok(A []byte, M [][]byte, s string) []string {
	x := make([]byte, len(s))
	x[0] = A[s[0]]
	for i := 1; i < len(s); i++ {
		x[i] = M[x[i-1]][A[s[i]]]
	}
	r := make([]int, 0, len(s))
	for i := range x {
		if 5 > x[i] {
			r = append(r, i)
		}
	}
	r = append(r, len(s))
	v := make([]string, len(r)-1)
	for i := 0; i < len(r)-1; i++ {
		v[i] = s[r[i]:r[i+1]]
	}
	return v
}
func AM() ([]byte, [][]byte) {
	//    !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_`abcdefghijklmnopqrstuvwxyz{|}~
	a := "*q*****;;****a*aaaaaaaaaa:;*****aaaaaaaaaaaaaaaaaaaaaaaaaa;/;**aaaaaaaaaaaaaaaaaaaaaaaaaaa;*;*"
	a = strings.Repeat(";", 127-len(a)) + a
	A := make([]byte, 127)
	C := ";*aq/:"
	for i := range a {
		A[i] = find(C, a[i])
	}
	Cbr := C + "br"
	//              ;*aq/:     ;*aq/:     ;*aq/:     ;*aq/:     ;*aq/:     ;*aq/:     ;*aq/:     ;*aq/:
	m := []string{";;*aq**", "*;*aq*:", "a;*bq**", "qrrr:/r", "/rrrrrr", ":;*aq**", "b;*bq**", "rrrr:/r"}
	M := make([][]byte, len(m))
	for i := range m {
		M[i] = finds(Cbr, m[i][1:])
	}
	return A, M
}
func finds(s string, b string) []byte {
	r := make([]byte, len(b))
	for i := range b {
		r[i] = find(s, b[i])
	}
	return r
}
func find(s string, b byte) byte {
	i := strings.IndexByte(s, b)
	if i < 0 {
		return byte(len(s))
	}
	return byte(i)
}
