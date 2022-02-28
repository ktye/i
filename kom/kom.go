package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	. "github.com/ktye/wg/module"
)

//go:embed k_.go
var gsrc []byte // source of the k implementation

var xxx = `f:{x+a:1+x};f 3` // `{*a:-a:x}` // `*1+2`

var mo []string   // {"nul", "Idy", "Flp", ... } see ../k.go (func init)
var dy []string   // {"Asn", "Dex", "Add", ... }
var kfns [][]byte // generated code for k functions
var tablen int    // number of predefined functions in indirect call table (k.go:init "Functions...")

func main() {
	ksrc := zksrc()    // k source code (init + program)
	ksrc = []byte(xxx) //delete z.k for now
	kinit()
	for _, a := range os.Args[1:] {
		if strings.HasSuffix(a, ".k") {
			ksrc = append(ksrc, fileread(a)...)
		} else if strings.HasSuffix(a, ".go") {
			gsrc = append(gsrc, fileread(a)...)
		} else {
			fatal(fmt.Errorf("unknown suffix: " + a))
		}
	}

	var e error
	mo = monadics()
	dy = dyadics()
	tablen, e = strconv.Atoi(string(gsearch("kom:FTAB ")))
	fatal(e)
	symtab = make(map[string]int)

	var out bytes.Buffer
	// out := bytes.NewBuffer(gsrc)

	x := Prs(KC(ksrc)) // byte-code as L
	kom(&out, x, []string{}, 0)

	// patch function table in gsrc: (add generated k functions)
	// todo: overrite "nyi" in function table at 98 with "kom"
	a := bytes.Index(gsrc, []byte("kom:FTAB "))
	n := 2 + bytes.Index(gsrc[a:], []byte(")\n"))
	var buf bytes.Buffer
	if len(kfns) > 0 {
		buf.Write(gsrc[:a])
		fmt.Fprintf(&buf, "\n\tFunctions(%d, ", tablen)
		for i := 0; i < len(kfns); i++ {
			if i > 0 {
				fmt.Fprintf(&buf, ", ")
			}
			fmt.Fprintf(&buf, "f%d", i)
		}
		buf.WriteString(")\n")
	}
	buf.Write(gsrc[a+n:])
	// io.Copy(os.Stdout, &buf) // emit patched go interpreter source

	os.Stdout.Write([]byte(rtext)) // runtime extensions

	// compiling src generates code for k functions that is emmitted first
	for _, b := range kfns {
		os.Stdout.Write(b)
	}
	os.Stdout.Write([]byte(head))
	io.Copy(os.Stdout, &out)
	os.Stdout.Write([]byte("}\n"))
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
func fileread(f string) []byte {
	b, e := os.ReadFile(f)
	fatal(e)
	return b
}
func gsearch(s string) []byte {
	a := bytes.Index(gsrc, []byte(s))
	if a < 0 {
		panic(fmt.Errorf("cannot find %q in gsrc", s))
	}
	a += len(s)
	n := bytes.Index(gsrc[a:], []byte(")\n"))
	b := gsrc[a : a+n]
	return b
}
func zksrc() []byte { // extract z.k from gsrc
	s, e := strconv.Unquote(string(gsearch("Data(600, ")))
	fatal(e)
	return []byte(s)
}
func monadics() []string { // extract monadics from gsrc function table
	b := gsearch("Functions(00,")
	v := strings.Split(string(b), ",")
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}
	return v
}
func dyadics() []string { // extract dyadics from gsrc function table
	b := gsearch("Functions(64,")
	v := strings.Split(string(b), ",")
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}
	return v
}

var as []string           // assignment stack
var lo []string           // ssa variables
var ar int                // arity
var lm map[string]bool    // k-name to ssa name
var symtab map[string]int // k symbol index table for runtime

func kom(w io.Writer, x K, locals []string, args int) string { // see ../exec.go (func exec) how k executes byte code

	// kom is called recursively and (re)stores it's state in globals
	save_lo, save_ar, save_lm := lo, ar, lm
	lo, ar, lm = locals, args, make(map[string]bool)
	for i := 0; i < args; i++ {
		lm[locals[i]] = true
	}
	defer func() { lo, ar, lm = save_lo, save_ar, save_lm }()

	xn := nn(x)
	p := int32(x)
	e := p + 8*xn
	for p < e {
		u := K(I64(p))
		if tp(u) != 0 { // noun
			if isAsn(u, e, p) {
			} else if isLup(u, e, p) {
				p += 8 // skip .
				s := sK(u)
				if isglobal(s) {
					fmt.Fprintf(w, "%s := Val(Ks(%d)) // .%s\n", ssa(), symtab[s], s)
				} else {
					lo = append(lo, s)
					fmt.Fprintf(w, "rx(%s) // .%s\n", s, s)
				}
			} else {
				kmNoun(w, u)
			}

		} else {
			switch int32(u) >> 6 {
			case 0: //   0..63   monadic
				kmMo(w, u)
			case 1: //  64..127  dyadic
				kmDy(w, u)
			case 2: // 128       dyadic indirect
				fmt.Fprintf(w, "dyadic indirect\n")
			case 3: // 192..255  tetradic
				fmt.Fprintf(w, "tetradic\n")
			case 4: // 256       drop
				fmt.Fprintf(w, "dx(%s)\n", v0())
			case 5: // 320       jump
				fmt.Fprintf(w, "jump\n")
			case 6: // 384       jump if not
				fmt.Fprintf(w, "jump if not\n")
			default: //448..     quoted verb
				fmt.Fprintf(w, "quoted\n")
			}
		}
		p += 8
	}
	if len(lo) == 0 {
		return ""
	}
	return v0()
}
func isAsn(u K, e, p int32) bool {
	if tp(u) == 4 && e > p+8 {
		v := K(I64(p + 8))
		if tp(v) == 0 && int32(v) == 64 {
			as = append(as, sK(u))
			return true
		}
	}
	return false
}
func isLup(u K, e, p int32) bool {
	if tp(u) == 4 && e > p+8 {
		v := K(I64(p + 8))
		if tp(v) == 0 && int32(v) == 20 {
			return true
		}
	}
	return false
}
func isglobal(s string) bool {
	for i := 0; i < ar; i++ {
		if lo[i] == s {
			return false
		}
	}
	return true
}
func kmAsn(w io.Writer) {
	y := v0()
	s := as[len(as)-1]
	as = as[:len(as)-1]
	if isglobal(s) {
		fmt.Fprintf(w, "%s := Asn(Ks(%d), %s) // %s:\n", ssa(), symtab[s], v1(), s)
	} else {
		lo = append(lo, s)
		fmt.Fprintf(w, "dx(%s); %s = rx(%s)\n", s, s, y)
	}
}
func kmNoun(w io.Writer, x K) {
	t := tp(x)
	switch t {
	case 3: // i const
		fmt.Fprintf(w, "%s := Ki(%d)\n", ssa(), int32(x))
	case 4: //
		fmt.Fprintf(w, "noun case4\n")
	case 13: // lambda
		kmLambda(w, x)
	default:
		fmt.Fprintf(w, "noun type %d\n", t)
	}
}
func kmLambda(w io.Writer, x K) {
	var buf bytes.Buffer
	ary := int(nn(x))
	code := x0(int32(x))
	locs := x1(int32(x))
	kstr := x2(int32(x))
	v := SK(locs)
	r := kom(&buf, code, v, ary)

	n := "f" + strconv.Itoa(len(kfns))
	b := []byte("func " + n + "(" + strings.Join(v[:ary], ", "))
	if ary > 0 {
		b = append(b, ' ', 'K')
	}
	b = append(b, []byte(") K { // "+string(CK(kstr))+"\n")...)
	if len(v) > ary {
		zeros := strings.Join(strings.Split(strings.Repeat("0", len(v)-ary), ""), ", ")
		b = append(b, []byte("var "+strings.Join(v[ary:], ", ")+" K = "+zeros+"\n")...)
	}
	b = append(b, buf.Bytes()...)
	for i := 0; i < len(v); i++ {
		b = append(b, []byte("dx("+v[i]+")\n")...)
	}
	b = append(b, []byte(fmt.Sprintf("return %s\n", r))...)
	b = append(b, '}', '\n')
	kfns = append(kfns, b)

	fn := len(kfns) - 1
	fmt.Fprintf(w, "%s := LAMBDA(%d) // f%d \n", ssa(), tablen+fn, fn)
}
func kmMo(w io.Writer, u K) {
	fmt.Fprintf(w, "%s := %s(%s)\n", ssa(), mo[int32(u)], v1())
}
func kmDy(w io.Writer, u K) {
	if int32(u) == 64 { // asn
		kmAsn(w)
	} else {
		fmt.Fprintf(w, "%s := %s(%s, %s)\n", ssa(), dy[int32(u)-64], v1(), v2())
	}
}
func ssa() string {
	n := len(lo)
	s := "k" + strconv.Itoa(n)
	for lm[s] { // create unique name
		n++
		s = "k" + strconv.Itoa(n)
	}
	lm[s] = true
	lo = append(lo, s)
	return s
}
func vn(n int) string { return lo[len(lo)-n-1] } // n'th var name from top
func v0() string      { return vn(0) }
func v1() string      { return vn(1) }
func v2() string      { return vn(2) }

func KC(b []byte) (r K) {
	r = mk(Ct, int32(len(b)))
	copy(Bytes[int32(r):], b)
	return r
}
func sK(x K) string { x = cs(x); return string(Bytes[int32(x) : int32(x)+nn(x)]) }
func SK(x K) []string {
	r := make([]string, int(nn(x)))
	for i := range r {
		r[i] = sK(ati(rx(x), int32(i)))
	}
	return r
}
func CK(x K) []byte { return Bytes[int32(x) : int32(x)+nn(x)] }

const head string = `func main() {
`
const rtext string = `
func kom(f, x K) (r K) { // call compiled lambda function
	n := nn(x)
	p := int32(f)
	lfree(x)
	switch(n) {
	case 0:
		r = Func[p].(func() K)()
	case 1:
	        r = Func[p].(func(K) K)(K(I64(p)))
	case 2:
	        r = Func[p].(func(K, K) K)(K(I64(p)), K(I64(p+8)))
	case 3:
	        r = Func[p].(func(K, K, K) K)(K(I64(p)), K(I64(p+8), K(I64(p+16)))
	default: // todo: track needed maxargs during compilation
		r = trap(Nyi)
	}
	dx(f)
	return r
}
func LAMBDA(x int32) K {
	l := l2(Ki(x), mk(Ct, 0))
	return K(int32(l)) | K(14)<<59 // native function type xf
}
`
