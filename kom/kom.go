package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	. "github.com/ktye/wg/module"
)

//go:embed k_.go
var gsrc []byte // source of the k implementation

var xxx = "$[1;2;3]" // `{*a:-a:x}` // `*1+2`

var mo []string   // {"nul", "Idy", "Flp", ... } see ../k.go (func init)
var dy []string   // {"Asn", "Dex", "Add", ... }
var kfns [][]byte // generated code for k functions
var tablen int    // number of predefined functions in indirect call table (k.go:init "Functions...")

func main() {
	ksrc := zksrc() // k source code (init + program)
	//ksrc = []byte(xxx) //delete z.k for now
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
	konsts = make(map[string]string)

	var out bytes.Buffer

	x := Prs(KC(ksrc)) // byte-code as L
	markwhile(x)
	v := kom(&out, x, []string{}, 0)
	if v != "" { // fix last node (we could also print).
		out.Write([]byte("dx(" + v + ")\n"))
	}

	// replace function table@98 "nyi" with "kal" (call native function)
	a := bytes.Index(gsrc, []byte("Functions(64,"))
	if a < 0 {
		fatal(fmt.Errorf("cannot find 'Functions(64,' in gsrc"))
	}
	a += bytes.Index(gsrc[a:], []byte("nyi"))
	gsrc[0+a] = 'k'
	gsrc[1+a] = 'a'
	gsrc[2+a] = 'l'

	// patch function table in gsrc: (add generated k functions)
	a = bytes.Index(gsrc, []byte("kom:FTAB "))
	n := 1 + bytes.Index(gsrc[a:], []byte(")\n"))
	var buf bytes.Buffer
	buf.Write(gsrc[:a])
	if len(kfns) > 0 {
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
	io.Copy(os.Stdout, &buf) // emit patched go interpreter source

	os.Stdout.Write([]byte(rtext)) // add runtime extensions "kal", "LAMBDA"

	// emit code for compiled k functions
	for _, b := range kfns {
		os.Stdout.Write(b)
	}
	constdecl(os.Stdout)
	os.Stdout.Write([]byte("func main() {\nkinit()\n"))
	symbols(os.Stdout)
	if needkp() {
		os.Stdout.Write([]byte("var k_p int32\n"))
	}
	constnt(os.Stdout)
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
func markwhile(x K) { // mark while loops upfront
	// a while loop is [0 ... jif(to the end) ... jump(back)]
	// replace the leading 0 with -1 and the jif with -2
	xn := nn(x)
	p := int32(x)
	e := p + 8*xn
	for p < e {
		u := K(I64(p))
		j := int32(I64(p - 8))
		if u == 384 { // jif
			if I64(p+j) == 320 { // jump
				o := int32(I64(p + j - 8))
				if o < 0 { // negative jump offset (end of while)
					SetI64(p+j+o, -1) // replace leading 0
					SetI64(p, -2)     // replace jif with 0xfffffffffffffffe
					SetI64(p+j-8, -3) // replace negative jump offset with fffffffffffffffd
				}
			}
		}
		p += 8
	}
}

var as []string              // assignment stack
var lo []string              // ssa variables
var rs []string              // block return var stack
var ar int                   // arity
var lm map[string]bool       // k-name to ssa name
var ce []int32               // cond end positions
var symtab map[string]int    // k symbol index table for runtime
var konsts map[string]string // k runtime constants

func kom(w io.Writer, x K, locals []string, args int) string { // see ../exec.go (func exec) how k executes byte code

	// kom is called recursively and (re)stores it's state in globals
	save_lo, save_rs, save_ar, save_lm, save_ce := lo, rs, ar, lm, ce
	lo, rs, ar, lm, ce = locals, make([]string, 0), args, make(map[string]bool), make([]int32, 0)
	for i := 0; i < args; i++ {
		lm[locals[i]] = true
	}
	defer func() { lo, rs, ar, lm, ce = save_lo, save_rs, save_ar, save_lm, save_ce }()

	xn := nn(x)
	p := int32(x)
	e := p + 8*xn
	for p < e {
		u := K(I64(p))

		if u == 0xffffffffffffffff { // while
			fmt.Fprintf(w, "%s := K(0)\nfor {\n", ssa())
			rs = append(rs, v0())
		} else if u == 0xfffffffffffffffe { // jif within while
			fmt.Fprintf(w, "if %s != 0 { dx(%s); break; }\n", v0(), v0())
		} else if u == 0xfffffffffffffffd {
			r := rs[len(rs)-1]
			rs = rs[:len(rs)-1]
			fmt.Fprintf(w, "dx(%s); %s = %s\n}\n%s := %s\n", r, r, v0(), ssa(), r)
			p += 8
		} else if K(I64(p+8)) == 384 { // jif in cond
			fmt.Fprintf(w, "%s := K(0)\ndx(%s)\n", ssa(), v1()) // cond's return var
			fmt.Fprintf(w, "if int32(%s) != 0 {\n", v1())
			rs = append(rs, v0())
			p += 8
		} else if K(I64(p+8)) == 320 { // jmp in cond
			ce = append(ce, p+16+int32(u))
			fmt.Fprintf(w, "%s = %s\n} else {\n", rs[len(rs)-1], v0())
			p += 8
		} else if tp(u) != 0 { // noun
			if isAsn(u, e, p) {
			} else if isLup(u, e, p) {
				p += 8 // skip .
				s := sK(u)
				if isglobal(s) {
					fmt.Fprintf(w, "%s := Val(Ks(%d)) // .%s\n", ssa(), intern(s), s)
				} else {
					lo = append(lo, s)
					fmt.Fprintf(w, "rx(%s) // .%s\n", s, s)
				}
			} else if isLst(u, e, p) {
				mkLst(w, int(int32(u)))
				p += 8
			} else {
				fmt.Fprintf(w, "%s := %s\n", ssa(), kmNoun(u))
			}
		} else {
			switch int32(u) >> 6 {
			case 0: //   0..63   monadic
				kmMo(w, u)
			case 1: //  64..127  dyadic
				kmDy(w, u)
			case 2: // 128       dyadic indirect
				fmt.Fprintf(w, "%s := Cal(%s, l2(%s, %s))\n", ssa(), v1(), v2(), vn(3))
			case 3: // 192..255  tetradic
				fmt.Fprintf(w, "var %s K\nif 211 == int32(dx(%s)) {\n", ssa(), v1())
				s := fmt.Sprintf(" %s = Amd(%s, %s, %s, %s)\n", v0(), v1(), v2(), vn(3), vn(4))
				w.Write([]byte(s))
				fmt.Fprintf(w, "} else {\n")
				w.Write([]byte(strings.Replace(s, "Amd", "Dmd", 1)))
				fmt.Fprintf(w, "}\n")
			case 4: // 256       drop
				fmt.Fprintf(w, "dx(%s)\n", v0())
			case 5: // 320       jump
				fmt.Fprintf(w, "!jump\n")
			case 6: // 384       jump if not
				fmt.Fprintf(w, "!jump if not\n")
			default: //448..     quoted verb
				v := int(unquote(K(int32(u))))
				c := byte(' ')
				t := "0:+-*%!&|<>=~,^#_$?@."
				if v > 0 && v < len(t) {
					c = t[v]
				} else if v-64 > 0 && v-64 < len(t) {
					c = t[v-64]
				}
				fmt.Fprintf(w, "%s = K(%d) // `%d (%c)\n", ssa(), v, v, c)
			}
		}
		p += 8

		if len(ce) > 0 && ce[len(ce)-1] == p {
			rs = append(rs, v0())
			for len(ce) > 0 && ce[len(ce)-1] == p { // cond-end
				ce = ce[:len(ce)-1]
				fmt.Fprintf(w, "%s = %s\n}\n", rs[len(rs)-2], rs[len(rs)-1])
				rs = rs[:len(rs)-1]
			}
			lo = append(lo, rs[len(rs)-1])
			rs = rs[:len(rs)-1]
		}
	}
	if len(lo) == 0 {
		return ""
	}
	return v0()
}
func intern(s string) int32 {
	if i, o := symtab[s]; o {
		return int32(i)
	} else {
		i := 8 * len(symtab)
		symtab[s] = i
		return int32(i)
	}
}
func konst(s string) (r string) {
	if g, o := konsts[s]; o {
		return g
	} else {
		r := "konst_" + strconv.Itoa(len(konsts))
		konsts[s] = r
		return r
	}
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
		fmt.Fprintf(w, "%s := Asn(Ks(%d), %s) // %s:\n", ssa(), intern(s), v1(), s)
	} else {
		lo = append(lo, s)
		fmt.Fprintf(w, "dx(%s); %s = rx(%s)\n", s, s, y)
	}
}
func kmNoun(x K) string {
	t := tp(x)
	vec := func(sz int) string { // initialize all vector sizes using 64bit int data
		n := int(nn(x))
		s := string(CK(Kst(rx(x))))
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "mk(%d, %d) // %s\n", t, n, s)
		fmt.Fprintf(&buf, "k_p = int32(@)\n")
		if n > 0 {
			m := sz * n / 8
			if 8*m < sz*n {
				for i := 8 * m; i < sz*n; i++ { // blank trailing memory
					SetI8(int32(i), 0)
				}
				m++
			}
			for i := 0; i < m; i++ {
				fmt.Fprintf(&buf, "SetI64(k_p+%d, %d)\n", 8*i, I64(int32(x)+int32(8*i)))
			}
		}
		return fmt.Sprintf("rx(%s)", konst(string(buf.Bytes())))
	}
	switch t {
	case 1: // b
		return fmt.Sprintf("Kb(%d)", int32(x))
	case 2: // c
		return fmt.Sprintf("Kc(%d)", int32(x))
	case 3: // i
		return fmt.Sprintf("Ki(%d)", int32(x))
	case 4: // s
		s := sK(x)
		return fmt.Sprintf("Ks(%d) // %s", intern(s), s)
	case 5: // f
		return fmt.Sprintf("Kf(%v)", F64(int32(x)))
	case 6: // z
		return fmt.Sprintf("Kz(%v,%v)", F64(int32(x)), F64(int32(x)+8))
	case 13: // lambda
		return kmLambda(x)
	case 17: // B
		return vec(1)
	case 18: // C
		return fmt.Sprintf("rx(%s)", konst(mkC(string(CK(x)))))
	case 19: // I
		return vec(4)
	case 20: // S
		n := nn(x)
		s := string(CK(Kst(rx(x))))
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "mk(20, %d) // %s\n", n, s)
		fmt.Fprintf(&buf, "k_p = int32(@)\n")
		if n > 0 {
			for i := int32(0); i < n; i++ {
				y := ati(rx(x), i)
				S := sK(y)
				s := intern(S)
				dx(y)
				fmt.Fprintf(&buf, "SetI32(k_p+%d, %d) // %s\n", 4*i, s, S)
			}
		}
		return fmt.Sprintf("rx(%s)", konst(string(buf.Bytes())))
	case 21: // F
		return vec(8)
	case 22: // Z
		return vec(16)
	case 23: // L
		n := nn(x)
		s := string(CK(Kst(rx(x))))
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "mk(23, %d) // %s\n", n, s)
		fmt.Fprintf(&buf, "konst_p = int32(@)\n")
		if n > 0 {
			for i := int32(0); i < n; i++ {
				y := ati(rx(x), i)
				s := string(CK(Kst(rx(y))))
				fmt.Fprintf(&buf, "SetI64(k_p+%d, int64(%s)) // %s\n", 8*i, kmNoun(y), s)
			}
		}
		return fmt.Sprintf("rx(%s)", konst(string(buf.Bytes())))
	// functions cannot exist, dicts, tables should not.
	// todo: folded constants? compositions?  +/?  (+;-)?
	default:
		fatal(fmt.Errorf("noun type %d", t))
		return "?"
	}
}
func kmLambda(x K) string {
	single := func(s string) string { return strings.Replace(s, "\n", ";", -1) }
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
	b = append(b, []byte(") K { // "+single(string(CK(kstr)))+"\n")...)
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
	return fmt.Sprintf("LAMBDA(%d)", tablen+fn)
}
func kmMo(w io.Writer, u K) {
	if u == 0 {
		fmt.Fprintf(w, "%s := K(0)\n", ssa())
		return
	}
	fmt.Fprintf(w, "%s := %s(%s)\n", ssa(), mo[int32(u)], v1())
}
func kmDy(w io.Writer, u K) {
	if int32(u) == 64 { // asn
		kmAsn(w)
	} else {
		fmt.Fprintf(w, "%s := %s(%s, %s)\n", ssa(), dy[int32(u)-64], v1(), v2())
	}
}
func isLst(u K, e, p int32) bool {
	if tp(u) == 3 && p < e { // is u always known at compile time (int const)?
		v := K(I64(p + 8))
		if tp(v) == 0 && int32(v) == 27 { // next verb is 'lst' (create dyanmic list)
			return true // #u
		}
	}
	return false
}
func mkLst(w io.Writer, n int) { // build list at runtime
	fmt.Fprintf(w, "%s := mk(23, %d)\n", ssa(), n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(w, "SetI64(int32(%s)+%d, %s)\n", v0(), 8*i, vn(1+i))
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

// generate code creates symbols in the right order at runtime.
func symbols(w io.Writer) {
	if len(symtab) == 0 {
		return
	}
	fmt.Fprintf(w, "// symbol table\n")
	v := make([]string, len(symtab))
	i := 0
	for s := range symtab {
		v[i] = s
		i++
	}
	sort.Slice(v, func(i, j int) bool { return symtab[v[i]] < symtab[v[j]] })
	for i, s := range v {
		fmt.Fprintf(w, "sc(%s) // `%s %d\n", mkC(s), s, 8*i)
	}
}
func mkC(s string) string { // "Cat(Cat(Ku(123..), Ku(5678..)), Ku(9012..))"
	if s == "" {
		return "mk(Ct, 0)"
	}
	enc := func(x []byte) (r uint64) {
		var o uint64 = 1
		for _, b := range x {
			r += o * uint64(b)
			o <<= 8
		}
		return r
	}
	b := []byte(s)
	i := 0
	var r string
	for i < len(b) {
		s := fmt.Sprintf("Ku(%d)", enc(b[i:int(mini(int32(8+i), int32(len(b))))]))
		if i == 0 {
			r = s
		} else {
			r = "Cat(" + r + ", " + s + ")"
		}
		i += 8
	}
	return r
}

// generate code that initializes k runtime constants.
func constdecl(w io.Writer) { // declared in global context
	if len(konsts) == 0 {
		return
	}
	var v []string
	for _, k := range konsts {
		v = append(v, k)
	}
	fmt.Fprintf(w, "var %s K\n", strings.Join(v, ", "))
}
func constnt(w io.Writer) { // initialized in main
	if len(konsts) == 0 {
		return
	}
	fmt.Fprintf(w, "// runtime constants\n")
	for c, k := range konsts {
		fmt.Fprintf(w, "%s = %s\n", k, strings.Replace(c, "@", k, -1))
	}
}
func needkp() bool { // need variable k_p
	for c := range konsts {
		if strings.Index(c, "@") >= 0 {
			return true
		}
	}
	return false
}

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
func kal(f, x K) (r K) { // call kom-compiled lambda function
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
	        r = Func[p].(func(K, K, K) K)(K(I64(p)), K(I64(p+8)), K(I64(p+16)))
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
