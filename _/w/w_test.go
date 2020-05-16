package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

var broken = false // ../../k.w

func TestB(t *testing.T) {
	testCases := []struct {
		sig string
		b   string
		e   string
	}{
		{"I:I", "(I.x) 3", "4103 2000 11 00 00"},
		{"I:IF", "x::y;x", "2000 2001 390300 2000"},
		{"I:II", "y+I?C x", "2001 2000 2d0000 6a"},
		{"I:II", "y+I?C x", "2001 2000 2d0000 6a"},
		{"I:I", "x?[;x:4;x:6];x", "024002400240024020000e020001020b0c010b410421000c010b410621000c000b2000"},
		{"I:I", "x?[x:4;;x:6];x", "024002400240024020000e020001020b410421000c020b0c000b410621000c000b2000"},
		{"V:I", "x::5", "2000 4105 360200"},
		{"V:I", "x::C?C 5", "2000 4105 2d0000 3a0000"},
		{"V:I", "0::1130366807310592j", "4100 42 8082 90c0 8082 8102 370300"},
		{"I:I", "x?!;x", "20000440 00 0b 2000"},
		{"I:I", "-1+x", "417f 2000 6a"},
		{"I:I", "x-1", "200041016b"},
		{"I:II", "x:I 4+x;x", "4104 2000 6a 280200 2100 2000"},
		{"I:I", "x?[x:4;x:5;x:6];x", "024002400240024020000e020001020b410421000c020b410521000c010b410621000c000b2000"},
		{"I:I", "I?255j&1130366807310592j>>J?8*x", "42ff0142808290c080828102410820006cad8883a7"},
		{"I:I", "(x<6)?/x+:1;x", "0240 0340 2000 4106 49 45 0d01 20004101 6a 2100 0c00 0b0b2000"},
		{"I:I", "(x<6)?/(x+:1;x+:1);x", "0240 0340 2000 4106 49 45 0d01 200041016a2100 200041016a2100 0c00 0b0b2000"},
		{"I:I", "1/(x+:1;?x>5);x", "0240 0340 2000 4101 6a 2100 2000 4105 4b  0d01 0c00 0b0b2000"},
		{"I:III", "$[x;y;z]", "2000044020010520020b"},
		{"I:I", "(x>3)?(:-x);x", "2000 4103 4b 0440 4100 2000 6b 0f 0b 2000"},
		{"I:I", "(x>3)?x+:1;x", "2000 4103 4b 0440 2000 4101 6a 2100 0b 2000"},
		{"I:II", "x::y;I x", "2000 2001 360200 2000 280200"},
		{"I:I", "x/r:r+i;r", "20000440410021020340200120026a2101200241016a22022000490d000b0b2001"},
		{"I:I", "x/r+:i;r", "20000440410021020340200120026a2101200241016a22022000490d000b0b2001"},
		{"I:II", "x+y", "20002001 6a"},
		{"I:II", "x\\y", "20002001 70"},
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
		e := f.parse(nil, nil, nil, map[string]int{"I:I": 0})
		b := hex.EncodeToString(e.bytes())
		s := trim(tc.e)
		if b != s {
			t.Fatalf("#%d:%s\n expected/got:\n%s\n%s", n+1, tc.b, s, b)
		}
		// fmt.Println(b)
		ctest(t, tc.sig, tc.b)
	}
}
func TestRun(t *testing.T) {
	testCases := [][2]string{
		{"add:I:II{x+y}/comment\n/\n/sum:I:I{x/r+:i;r}\n/", "0061736d0100000001070160027f7f017f030201000503010001070d02036d656d02000361646400000a09010700200020016a0b"},
		{"add:I:II{x+y}10!abcd ", "0061736d0100000001070160027f7f017f030201000503010001070d02036d656d02000361646400000a09010700200020016a0b0b080100410a0b02abcd"},
	}
	for _, tc := range testCases {
		m, tab, data := run(strings.NewReader(tc[0]))
		g := hex.EncodeToString(m.wasm(tab, data))
		e := tc[1] // "0061736d0100000001070160027f7f017f030201000503010001070d02036d656d02000361646400000a09010700200020016a0b"
		if e != g {
			t.Fatalf("expected/got\n%s\n%s\n", e, g)
		}
	}

}
func ctest(t *testing.T, sig, b s) {
	b = jn("f:", sig, "{", b, "}")
	m, tab, data := run(strings.NewReader(b))
	out := m.cout(tab, data)
	if len(out) == 0 {
		t.Fatal("no output")
	}
	//fmt.Println(string(out))
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
func TestHtml(t *testing.T) { // write k.html from ../../k.w
	if broken {
		t.Skip()
	}
	m, tab, data, src, err := KWasmModule()
	if err != nil {
		t.Fatal(err)
	}
	help := []byte{}
	if idx := bytes.Index(src, []byte{'\n', '\\'}); idx != -1 {
		help = src[idx+3:]
		src = src[:idx+1]
	}
	tests, err := ioutil.ReadFile("t")
	if err != nil {
		t.Fatal(err)
	}
	wasm := m.wasm(tab, data)
	gui, err := ioutil.ReadFile("dat_gui")
	if err != nil {
		t.Fatal(err)
	}
	var txt bytes.Buffer
	fmt.Fprintf(&txt, "k.w(%d b) %s [\\\\src \\\\tests \\\\h]\\n ", len(wasm), time.Now().Format("2006.01.02"))
	var b bytes.Buffer
	s, e := ioutil.ReadFile("k_html")
	if e != nil {
		t.Fatal(e)
	}
	s = bytes.Replace(s, []byte(`{{wasm}}`), []byte(base64.StdEncoding.EncodeToString(wasm)), 1)
	s = bytes.Replace(s, []byte(`{{tests}}`), []byte(base64.StdEncoding.EncodeToString(tests)), 1)
	s = bytes.Replace(s, []byte(`{{src}}`), []byte(base64.StdEncoding.EncodeToString(src)), 1)
	s = bytes.Replace(s, []byte(`{{help}}`), []byte(base64.StdEncoding.EncodeToString(help)), 1)
	s = bytes.Replace(s, []byte(`{{gui}}`), gui, 1)
	s = bytes.Replace(s, []byte(`{{cons}}`), txt.Bytes(), 1)
	//s = bytes.Replace(s, []byte(`{{fncs}}`), fns.Bytes(), 1)
	b.Write(s)
	if e := ioutil.WriteFile("k.html", b.Bytes(), 0644); e != nil {
		t.Fatal(e)
	}
}
func KWasmModule() (module, []segment, []dataseg, []byte, error) {
	var src io.Reader
	var srcb []byte
	if k, e := ioutil.ReadFile("../../k.w"); e != nil {
		return nil, nil, nil, nil, e
	} else {
		src = bytes.NewReader(k)
		srcb = k
	}
	m, tab, data := run(src)
	return m, tab, data, srcb, nil
}
func TestCout(t *testing.T) { // write k_h h_h from ../../k.w
	if broken {
		t.Skip()
	}
	m, tab, data, src, err := KWasmModule()
	if err != nil {
		t.Fatal(err)
	}
	help := bytes.Index(src, []byte{'\n', '\\'})
	if help != -1 {
		f, e := os.Create("h_h")
		if e != nil {
			t.Fatal(e)
		}
		defer f.Close()
		version := fmt.Sprintf("k.w(c) %s", time.Now().Format("2006.01.02"))
		q := fmt.Sprintf("%q", string(src[help+3:]))
		q = strings.Replace(q, "%", `%%`, -1)
		fmt.Fprintf(f, "const char *version=%q;\nV help(){printf(%s);}\n", version, q)
	}
	var dst bytes.Buffer
	dst.Write(m.cout(tab, data))
	if e := ioutil.WriteFile("k_h", dst.Bytes(), 0744); e != nil {
		t.Fatal(e)
	}
}
func TestGout(t *testing.T) { // write kw.go from ../../k.w
	if broken {
		t.Skip()
	}
	m, tab, data, _, err := KWasmModule()
	if err != nil {
		t.Fatal(err)
	}
	var dst bytes.Buffer
	io.Copy(&dst, strings.NewReader(gh))
	dst.Write(m.gout(tab, data))
	if e := ioutil.WriteFile("kw.go", dst.Bytes(), 0744); e != nil {
		t.Fatal(e)
	}
}

const gh = `// +build ignore

package main
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"math/bits"
	"math/cmplx"
	"os"
	"strconv"
	"strings"
	"unsafe"
)
func init() {
	NAN = math.NaN()
}
type C=byte
type I=uint32
type J=uint64
type F=float64
type SI=int32
type slice struct {
	p uintptr
	l int
	c int
}
var MC []C
var MI []I
var MJ []J
var MF []F
var NAN F
func sin(x F) F { return math.Sin(x) }
func cos(x F) F { return math.Cos(x) }
func atan2(x, y F) F { return math.Atan2(x, y) }
func hypot(x, y F) F { return math.Hypot(x, y) }
func draw(x, y, z I) { fmt.Printf("draw %x %x %x\n", x, y, z) }
func msl() { // update slice headers after set/inc MJ
	cp := *(*slice)(unsafe.Pointer(&MC))
	ip := *(*slice)(unsafe.Pointer(&MI))
	jp := *(*slice)(unsafe.Pointer(&MJ))
	fp := *(*slice)(unsafe.Pointer(&MF))
	fp.l, fp.c, fp.p = jp.l, jp.c, jp.p
	ip.l, ip.c, ip.p = jp.l*2, jp.c*2, jp.p
	cp.l, cp.c, cp.p = ip.l*4, ip.c*4, ip.p
	MF = *(*[]F)(unsafe.Pointer(&fp))
	MI = *(*[]I)(unsafe.Pointer(&ip))
	MC = *(*[]byte)(unsafe.Pointer(&cp))
}
func grow(x I) I { panic("nyi grow"); return x }
func clz32(x I) I { return I(bits.LeadingZeros32(x)) }
func clz64(x J) I { return I(bits.LeadingZeros64(x)) }
func i32b(x bool) I { if x { return 1 } else { return 0 } }
func n32(x I) I { if x == 0 { return 1 } else { return 0 } }
func dump(a, n I) { // type: cifzsld -> 2468ace
	p := a >> 2
	fmt.Printf("%.8x ", a)
	for i := I(0); i < n; i++ {
		x := MI[p+i]
		fmt.Printf(" %.8x", x)
		if i > 0 && (i+1)%8 == 0 {
			fmt.Printf("\n%.8x ", a+4*i+4)
		} else if i > 0 && (i+1)%4 == 0 {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
}
func X(x I) string {
	type s = string
	type i = I
	Z := func(a i) complex128 { return complex(MF[a>>3], MF[1+a>>3]) }
	if x == 0 || x == 128 {
		return ""
	}
	var t, n i
	n = 1
	if x > 255 {
		u := MI[x>>2]
		t = u>>29
		n = u&536870911
	}
	var f func(i i) s
	var tof func(s) s = func(s s) s { return s }
	istr := func(i i) s { return strconv.Itoa(int(int32(MI[i+2+x>>2]))) }
	fstr := func(i i) s {
		if f := MF[i+1+x>>3]; math.IsNaN(f) {
			return "0n"
		} else {
			return strconv.FormatFloat(f, 'g', -1, 64)
		}
	}
	zstr := func(i i) s {
		if z := Z(x + 8 + 16*i); cmplx.IsNaN(z) {
			return "0ni0n"
		} else {
			return strconv.FormatFloat(real(z), 'g', -1, 64) + "i" + strconv.FormatFloat(imag(z), 'g', -1, 64)
		}
	}
	sstr := func(i i) s {
		r := MI[(x + 8 + 4*i)>>2]
		rn := nn(r)
		return string(MC[r+8 : r+8+rn])
	}
	sep := " "
	switch t {
	case 0:
		fc := []byte(":+-*%&|<>=!~,^#_$?@.'/\\")
		if x < 128 && bytes.Index(fc, []byte{byte(x)}) != -1 {
			return string(byte(x))
		} else if x < 256 && bytes.Index(fc, []byte{byte(x - 128)}) != -1 {
			return string(byte(x-128)) + ":"
		} else if n == 3 {
			r := X(MI[(x + 12)>>2])
			return X(MI[(x+8)>>2]) + "[" + r[1:len(r)-1] + "]"
		} else if n == 4 {
			return sstr(0)
		} else {
			return fmt.Sprintf(" '(%d)", x)
		}
	case 1:
		return "\"" + string(MC[x+8:x+8+n]) + "\""
	case 2:
		f = istr
	case 3:
		f = fstr
		tof = func(s s) s {
			if strings.Index(s, ".") == -1 {
				return s + ".0"
			}
			return s
		}
	case 4:
		f = zstr
	case 5:
		f = sstr
		sep = string(96)
		if n == 0 { return "0#"+sep }
		tof = func(s s) s { return sep + s }
	case 6:
		if n == 1 {
			return "," + X(MI[(8+x)>>2])
		}
		f = func(i i) s { return X(MI[2+i+x>>2]) }
		sep = ";"
		tof = func(s s) s { return "(" + s + ")" }
	case 7:
		return X(MI[(x+8)>>2]) + "!" + X(MI[(x+12)>>2])
	default:
		panic(fmt.Sprintf("nyi: kst: t=%d", t))
	}
	r := make([]s, n)
	for k := range r {
		r[k] = f(i(k))
	}
	return tof(strings.Join(r, sep))
}
func runtest() {
	b, e := ioutil.ReadFile("t")
	if e != nil {
		panic(e)
	}
	v := strings.Split(string(b), "\n")
	for i := range v {
		if len(v[i]) == 0 {
			fmt.Println("skip rest")
			os.Exit(0)
		}
		if v[i][0] == '/' {
			continue
		}
		vv := strings.Split(v[i], " /")
		if len(vv) != 2 {
			panic("test file")
		}
		in := strings.TrimRight(vv[0], " \t\r")
		exp := strings.TrimSpace(vv[1])
		got := run(in)
		fmt.Println(in, "/", got)
		if exp != got {
			fmt.Println("expected:", exp)
			os.Exit(1)
		}
	}
}
func mark() {
	for t := uint32(4); t < 32; t++ {
		p := MI[t]
		for p != 0 {
			MI[1+(p>>2)] = 0
			MI[2+(p>>2)] = t
			p = MI[p>>2]
		}
	}
}
func leak() {
	//dump(0, 200)
	mark()
	p := uint32(64)
	for p < uint32(len(MI)) {
		if MI[p+1] != 0 {
			panic(fmt.Errorf("non-free block at %d(%x)", p<<2, p<<2))
		}
		t := MI[p+2]
		if t < 4 || t > 31 {
			panic(fmt.Errorf("illegal bucket type %d at %d(%x)", t, p<<2, p<<2))
		}
		dp := uint32(1) << t
		p += dp >> 2
	}
}
func run(s string) string {
	m0 := 16
	MJ = make([]J, (1<<m0)>>3)
	msl()
	mt_init()
	ini(16)
	x := mk(1, I(len(s)))
	copy(MC[x+8:], s)
	x = kst(val(x))
	n := nn(x)
	s = string(MC[8+x:8+x+n])
	dx(x)
	dx(MI[132>>2]) //kkey
	dx(MI[136>>2]) //kval
	dx(MI[148>>2]) //xyz
	leak()
	return s
}
func main() {
	//m0 := 16
	//MJ = make([]J, (1<<m0)>>3)
	//msl()
	//mt_init()
	//ini(16)
	runtest()
}
`
