package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

var broken = true // ../../k.w

func TestB(t *testing.T) {
	testCases := []struct {
		sig string
		b   string
		e   string
	}{
		{"I:I", "x?[;x:4;x:6];x", "024002400240024020000e020001020b0c010b410421000c010b410621000c000b2000"},
		{"I:I", "x?[x:4;;x:6];x", "024002400240024020000e020001020b410421000c020b0c000b410621000c000b2000"},
		{"0:I", "x::5", "2000 4102 74 4105 360200"},
		{"0:I", "(x,1)::5", "2000 4102 74 4101 4102 6c6a 4105 360200"},
		{"I:I", "x?!;x", "20000440 00 0b 2000"},
		{"I:I", "-1+x", "417f 2000 6a"},
		{"I:I", "x-1", "200041016b"},
		{"I:II", "x:I 4+x;x", "4104 2000 6a 410274 28 0200 2100 2000"},
		{"I:I", "x?[x:4;x:5;x:6];x", "024002400240024020000e020001020b410421000c020b410521000c010b410621000c000b2000"},
		{"I:I", "I?255j&1130366807310592j>>J?8*x", "42ff0142808290c080828102410820006cad8883a7"},
		{"I:I", "(x<6)?/x+:1;x", "0240 0340 2000 4106 49 45 0d01 20004101 6a 2100 0c00 0b0b2000"},
		{"I:I", "(x<6)?/(x+:1;x+:1);x", "0240 0340 2000 4106 49 45 0d01 200041016a2100 200041016a2100 0c00 0b0b2000"},
		{"I:I", "1/(x+:1;?x>5);x", "0240 0340 2000 4101 6a 2100 2000 4105 4b  0d01 0c00 0b0b2000"},
		{"I:III", "$[x;y;z]", "2000 047f 2001 05 2002 0b"},
		{"I:I", "(x>3)?(:-x);x", "2000 4103 4b 0440 4100 2000 6b 0f 0b 2000"},
		{"I:I", "(x>3)?x+:1;x", "2000 4103 4b 0440 2000 4101 6a 2100 0b 2000"},
		{"I:II", "x::y;I x", "2000 410274 2001 360200 2000 410274 280200"},
		{"I:I", "x/r:r+i;r", "20000440410021020340200120026a2101200241016a22022000490d000b0b2001"},
		{"I:I", "x/r+:i;r", "20000440410021020340200120026a2101200241016a22022000490d000b0b2001"},
		{"I:II", "x+y", "20002001 6a"},
		{"I:II", "r:x;r+:y;r", "2000 2102 2002 2001 6a 2102 2002"},
		{"I:I", "x/r:i;r", "2000044041002102034020022101200241016a22022000490d000b0b2001"},
		{"I:II", "(3+x)*y", "4103 2000 6a 2001 6c"},
		{"I:I", "1+x", "410120006a"},
		{"F:FF", "(x*y)", "20002001 a2"},
		{"F:FF", "x-y", "20002001 a1"},
		{"F:FF", "3.*x+y", "44 0000000000000840 20002001 a0 a2"},
		{"I:I", "x:1+x;x*2", "4101 2000 6a 2100 2000 4102 6c"},
	}
	for n, tc := range testCases {
		f := newfn(tc.sig, tc.b)
		e := f.parse(nil, nil, nil)
		b := string(hex(e.bytes()))
		s := trim(tc.e)
		if b != s {
			t.Fatalf("%d: expected/got:\n%s\n%s", n+1, s, b)
		}
		fmt.Println(b)
		ctest(t, tc.sig, tc.b)
	}
}
func TestRun(t *testing.T) {
	m, data := run(strings.NewReader("add:I:II{x+y}/cnt\n/\n/sum:I:I{x/r+:i;r}\n/"))
	g := s(hex(m.wasm(data)))
	e := "0061736d0100000001070160027f7f017f0302010005030100010707010361646400000a0b010901027f200020016a0b"
	if e != g {
		t.Fatalf("expected/got\n%s\n%s\n", e, g)
	}
}
func ctest(t *testing.T, sig, b s) {
	b = jn("f:", sig, "{", b, "}")
	m, data := run(strings.NewReader(b))
	out := m.cout(data)
	if len(out) == 0 {
		t.Fatal("no output")
	}
	//fmt.Println(string(out))
}
func hex(a []c) []c {
	var r bytes.Buffer
	for _, b := range a {
		hi, lo := hxb(b)
		r.WriteByte(hi)
		r.WriteByte(lo)
	}
	return r.Bytes()
}
func newfn(sig string, body string) fn {
	var buf bytes.Buffer
	buf.WriteString(body)
	buf.WriteByte('}')
	v := strings.Split(sig, ":")
	if len(v) != 2 {
		panic("signature")
	}
	f := fn{src: [2]int{1, 0}, Buffer: buf}
	f.t = typs[v[0][0]]
	for _, c := range v[1] {
		f.locl = append(f.locl, typs[byte(c)])
	}
	f.args = len(v[1])
	return f
}
func trim(s string) string { return strings.Replace(s, " ", "", -1) }

func TestCout(t *testing.T) {
	if broken {
		t.Skip()
	}
	var dst bytes.Buffer
	var src io.Reader
	if k, e := ioutil.ReadFile("../../k.w"); e != nil {
		t.Fatal(e)
	} else {
		src = bytes.NewReader(k)
	}
	io.Copy(&dst, strings.NewReader(kh))
	m, data := run(src)
	dst.Write(m.cout(data))
	io.Copy(&dst, strings.NewReader(kt1))
	for _, f := range m {
		if f.args == 1 && f.t == I && f.locl[0] == I && f.name != "ini" && f.name != "mki" {
			s := "\t\t} else if (Match(\"" + f.name + "\", a)) { n = f1(" + f.name + ", stack, n);\n"
			dst.WriteString(s)
		} else if f.args == 2 && f.t == I && f.locl[0] == I && f.locl[1] == I {
			s := "\t\t} else if (Match(\"" + f.name + "\", a)) { n = f2(" + f.name + ", stack, n);\n"
			dst.WriteString(s)
		}
	}
	io.Copy(&dst, strings.NewReader(kt2))
	if e := ioutil.WriteFile("k_c", dst.Bytes(), 0744); e != nil {
		t.Fatal(e)
	}
}

const kh = `#include<stdlib.h>
#include<stdio.h>
#include<stddef.h>
#include<malloc.h>
#include<math.h>
#define R return
typedef int32_t I;typedef int64_t J;typedef double F;typedef uint32_t uI;typedef uint64_t uJ;
I __builtin_clz(I x){I r;__asm__("bsr %1, %0" : "=r" (r) : "rm" (x) : "cc");R r^31;}
void trap() { exit(1); }
char *MC;I* MI;J* MJ;F *MF;`

const kt1 = `// Postfix test interface: e.g. 5 mki til rev fst 0 500 dump
I pop1(I *s, I n, I *x) {
	*x = s[n-1];
	return n-1;
}
I pop2(I *s, I n, I *x, I *y) {
	*x = s[n-2];
	*y = s[n-1];
	return n-2;
}
I push(I *s, I n, I x) {
	s[n] = x;
	return n+1;
}
I f1(I (*f)(I), I *s, I n) {
	printf("%d: ", s[n-1]);
	s[n-1] = f(s[n-1]);
	printf("%d\n", s[n-1]);
	return n;
}
I f2(I (*f)(I,I), I *s, I n) {
	printf("%d %d: ", s[n-2], s[n-1]);
	s[n-2] = f(s[n-2], s[n-1]);
	printf("%d\n", s[n-2]);
	return n-1;
}
I Number(char *s) {
	R strtol(s, (char **)NULL, 10);
}
I Match(char *a, char *b) {
	for (I i=0; ;i++) {
		if (a[i] != b[i]) return 0;
		if (a[i] == 0)    return 1;
	}
}
I Dump(I *s, I n) {
	I x = s[n-2];
	I y = s[n-1];
	I p = 0;
	printf("\n%08x  ", x);
	for (I i=x; i<x+y; i++) {
		printf("%02x", (uint8_t)MC[i]);
		p++;
		if ((i > x) && (p%32 == 0)) {
			printf("\n%08x  ", i+1);
		} else if ((i > x) && (p%4 == 0)) {
			printf(" ");
		}
	}
	return n-2;
}
I main(int args, char **argv){
	MC=malloc(64*1024);MI=(I*)MC;MJ=(J*)MC;MF=(F*)MC;
	I stack[32];
	I i, n = 0;
	I x, y, r;
	char *a;
	ini(16);
	for (i=1; i<args; i++) {
		a = argv[i];
		if (a[0] >= '0' && a[0] <= '9') {
			n = push(stack, n, Number(a));
			continue;
		}
		printf("%s ", argv[i]);
		if (Match("mki", a)) {
			n = f1(mki, stack, n);
`
const kt2 = `		} else if (Match("dump", a)) {
			n = Dump(stack, n);
		} else {
			printf("arg!");
			trap();
		}
	}
}
`
