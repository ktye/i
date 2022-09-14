package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

const max = 70

// version1: for each line, if too long, try to split+indent at lowest level (recurse)
func ikfmt1(s string) (r string) {
	f := func(v [][]byte) (r [][]byte, done bool) {
		for i := range v {
			if len(v[i]) > max && cansplit(v[i]) {
				l := lowestlevel(v[i])
				w := splitlevel(v[i], l)
				r = v[:i]
				r = append(r, w...)
				r = append(r, v[1+i:]...)
				return r, false
			}
		}
		return v, true
	}
	b := level_inline(s)
	v := [][]byte{b}
	for {
		var done bool
		v, done = f(v)
		if done {
			break
		}
	}
	for i := range v {
		v[i] = semi(v[i])
	}
	return string(bytes.Join(v, []byte{10}))
}
func semi(b []byte) []byte {
	for i, c := range b {
		if c < 32 {
			b[i] = ';'
		}
	}
	return b
}
func level_inline(s string) []byte {
	l := 1
	b := []byte(s)
	for i, c := range b {
		if c == '(' || c == '[' {
			l++
		} else if c == ')' || c == ']' {
			l--
		} else if c == ';' {
			b[i] = byte(l)
		}
	}
	return b
}
func cansplit(b []byte) bool {
	for _, c := range b {
		if c < 32 {
			return true
		}
	}
	return false
}
func lowestlevel(b []byte) (l byte) {
	l = 255
	for _, c := range b {
		if c < 32 && c < l {
			l = c
		}
	}
	return l
}
func splitlevel(b []byte, l byte) (r [][]byte) {
	space := strings.Repeat(" ", int(l))
	r = bytes.Split(b, []byte{l})
	for i := range r {
		if i > 0 {
			r[i] = append([]byte(space), r[i]...)
		}
	}
	return r
}

// version 0: always split and indent by level
func ikfmt0(s string) (r string) {
	l := level(";" + s)
	v := strings.Split(s, ";")
	for i := range l {
		v[i] = strings.Repeat(" ", l[i]) + v[i]
	}
	r = strings.Join(v, "\n")
	return r[1:len(r)]
}
func level(s string) (r []int) {
	l := 1
	for _, c := range s {
		if c == '(' || c == '[' {
			l++
		} else if c == ')' || c == ']' {
			l--
		} else if c == ';' {
			r = append(r, l)
		}
	}
	return r
}

func main() {
	var b []byte
	if len(os.Args) > 1 {
		var e error
		b, e = os.ReadFile(os.Args[1])
		fatal(e)
	} else {
		b = []byte(examples)
	}
	for _, s := range bytes.Split(b, []byte{10}) {
		i := bytes.IndexByte(s, '{')
		if i < 0 {
			fmt.Println(string(s))
		} else {
			fmt.Print(string(s[:i]))
			fmt.Println(ikfmt1(string(s[1+i:])))
		}
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

const examples = `repl::{[x]n:nn x;xp:i x;s:0;$[n>0;s:B xp;$[((B xp)=92)&&n>1;c4:B 1+xp;$[(B 1+xp)=92;Exit 0;$[c4=109;dx x;dx[Out[Ki[I 128]]]]];:()]];x:val x;$[x~k 0 0;$[s=32;dx[Out x];write[cat1[join[Kc[10];Lst x];Kc[10]]]]]}
doargs::{[]a:ndrop[1;getargv[]];an:nn a;ap:i a;ee:Ku[k 0 25901];i1:0;(i1<an;i1:i1+1)'[x2:x0[ap];$[(match[x2;ee])~0;$[i1<an-1;dx[x2];x2:x1[ap];dx[ee];dx a;repl[x2]];Exit 0];dofile[x2;readfile[rx[x2]]];ap:ap+8];dx[ee];dx a}
dofile::{[x,c]kk:Ku[k 0 27438];tt:Ku[k 0 29742];xe:ntake[-2;rx x];$[(match[xe;kk])~0;dx[val c];$[(match[xe;tt])~0;test c;dx[Asn[sc[rx x];c]]]];dx[xe];dx x;dx[tt];dx[kk]}
Out:k:{[x]write[cat1[Kst[rx x];Kc[10]]];x}
Otu:k:{[x,y]write[cat1[Kst x;Kc[58]]];Out y}
read:k:{[]r:mk[Ct;504];ntake[ReadIn[i r;504];r]}
write::{[x]Write[0;0;i x;nn x];dx x}
longempty:{this is a long line with an empty value;;^that was the empty value}`
