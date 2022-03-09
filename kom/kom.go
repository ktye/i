package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	. "github.com/ktye/wg/module"
)

//go:embed k_.go
var gsrc []byte // source of the k implementation

var mo []string              // {"nul", "Idy", "Flp", ... } see ../k.go (func init)
var dy []string              // {"Asn", "Dex", "Add", ... }
var kfns [][]byte            // generated code for k functions
var tablen int               // number of predefined functions in indirect call table (k.go:init "Functions...")
var maxargs int              //
var lmbdas map[string]string // compile identical functions only once
var krep bool                // append repl to main

// kom args > a.go
// args:
//  *.k   compile k
//  *.go  add native go
//  *.t   add tests
//  repl  add repl loop
func main() {
	kinit()
	ksrc := zksrc() // k source code (z.k + program)
	//ksrc = ksrc[:76] // symbols 76
	//ksrc = ksrc[:519] // math 519
	//ksrc = ksrc[:1559] // print 1559

	//ksrc = append(ksrc, addtests(ksrc, 0, -1)...)
	//ksrc = append(ksrc, []byte(xxx)...)
	for _, a := range os.Args[1:] {
		if strings.HasSuffix(a, ".k") {
			ksrc = append(ksrc, fileread(a)...)
		} else if strings.HasSuffix(a, ".go") {
			gsrc = append(gsrc, fileread(a)...)
		} else if strings.HasSuffix(a, ".t") {
			ksrc = append(ksrc, addtests(a, 0, -1)...)
		} else if a == "repl" {
			krep = true
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
	konsts = make(map[string]int)
	lmbdas = make(map[string]string)

	var out bytes.Buffer

	x := Prs(KC(ksrc)) // byte-code as L
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

	// remove zk() (don't run interpreted initialization code)
	a = bytes.Index(gsrc, []byte("\nfunc zk()"))
	n := bytes.Index(gsrc[a:], []byte("\n}\n"))
	gsrc = append(gsrc[:a], gsrc[2+a+n:]...)
	gsrc = bytes.Replace(gsrc, []byte("zk()"), nil, 1)

	// patch function table in gsrc: (add generated k functions)
	a = bytes.Index(gsrc, []byte("kom:FTAB "))
	n = 1 + bytes.Index(gsrc[a:], []byte(")\n"))
	var buf bytes.Buffer
	buf.Write(gsrc[:a])
	if len(kfns) > 0 {
		fmt.Fprintf(&buf, "\n\tFunctions(%d, ", tablen)
		for i := 0; i < len(kfns); i++ {
			if i > 0 {
				fmt.Fprintf(&buf, ", ")
			}
			fmt.Fprintf(&buf, "f_%d", i)
		}
		buf.WriteString(")\n")
	}
	buf.Write(gsrc[a+n:])
	io.Copy(os.Stdout, &buf) // emit patched go interpreter source

	os.Stdout.Write([]byte(rtext())) // add runtime extensions "kal", "lmb"
	os.Stdout.Write([]byte(gcomment(string(ksrc))))

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
	if krep {
		os.Stdout.Write([]byte(replsrc))
	}
	os.Stdout.Write([]byte("}\n"))
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
func tests(file string) (r [][2]string) {
	b := fileread(file)
	if len(b) > 0 && b[len(b)-1] == 10 {
		b = b[:len(b)-1]
	}
	v := bytes.Split(b, []byte{10})
	for i := range v {
		ab := strings.Split(string(v[i]), " /")
		r = append(r, [2]string{ab[0], ab[1]})
	}
	return r
}
func addtests(file string, a, n int) []byte { // all: a,n=0,-1
	var b []byte
	skip := make(map[string]bool)
	for _, s := range skiptests() {
		skip[s] = true
	}
	quote := func(s string) string {
		s = strings.ReplaceAll(s, `\`, `\\`)
		s = strings.ReplaceAll(s, `"`, `\"`)
		s = `"` + s + `"`
		if len(s) == 3 {
			s = "," + s
		}
		return s
	}
	b = append(b, []byte("ktest:{`<$[x~`k y;\"ok   \";`fail+\"fail \"],($z),\"\\n\"}\n")...)
	all := tests(file)
	if n < 0 {
		n = len(all) - a
	}
	for i := a; i < a+n; i++ {
		tc := all[i]
		if skip[tc[0]] {
			tc[0] = "`skip"
			tc[1] = "`skip"
		}
		//s := "`<$[" + quote(tc[1]) + "~`k@([" + tc[0] + `]);"ok   ";"fail "],($` + strconv.Itoa(1+i) + `),"\n"` + "\n"
		b = append(b, []byte(fmt.Sprintf("ktest[%s;[%s];%d]\n", quote(tc[1]), tc[0], 1+i))...)
	}
	//println(string(b))
	return b
}
func fileread(f string) []byte {
	b, e := os.ReadFile(f)
	fatal(e)
	return b
}
func gcomment(s string) string {
	v := strings.Split(s, "\n")
	for i := range v {
		v[i] = "// " + v[i]
	}
	return strings.Join(v, "\n") + "\n\n"
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
func zksrc() []byte { // extract z.k from gsrc, do not include plot.
	s, e := strconv.Unquote(string(gsearch("Data(600, ")))
	fatal(e)
	a := strings.Index(s, "PW:800;PH:600;") // plot
	if a < 0 {
		panic("cannot find plot in z.k")
	}
	s = s[:a]
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
	//  replace the leading 0 with -1 and, jif and offset with -2 and jmp-offset with -3
	// a return is identity+jmp [idy Ki(1048576) 320]
	//  replace all with -4
	xn := nn(x)
	p := int32(x)
	e := p + 8*xn
	for p < e {
		u := K(I64(p))
		if u == 384 { // jif
			j := int32(I64(p - 8))
			if I64(p+j) == 320 { // jump
				o := int32(I64(p + j - 8))
				if o < 0 { // negative jump offset (end of while)
					//println("WHILE")
					SetI64(p+j+o, -1) // replace leading 0
					SetI64(p-8, -2)   // jif offset
					SetI64(p, -2)     // replace jif with 0xfffffffffffffffe
					SetI64(p+j-8, -3) // replace negative jump offset with fffffffffffffffd
				}
			}
		} else if u == 320 && int32(I64(p-8)) == 1048576 {
			//println("RETURN")
			SetI64(p-16, -4) // identity
			SetI64(p-8, -4)  // jump offset
			SetI64(p, -4)    // jump
		}
		p += 8
	}
}

var locs []string         // local variables of current lambda
var lo []string           // ssa variables
var as []string           // assignment stack
var lm map[string]bool    // k-name to ssa name
var ce []int32            // cond end positions
var symtab map[string]int // k symbol index table for runtime
var konsts map[string]int // k runtime constants

func kom(w io.Writer, x K, locals []string, args int) string { // see ../exec.go (func exec) how k executes byte code

	// kom is called recursively and (re)stores it's state in globals
	save_lo, save_locs, save_lm, save_ce := lo, locs, lm, ce
	lo, locs, lm, ce = locals, locals, make(map[string]bool), make([]int32, 0)
	for i := 0; i < args; i++ {
		lm[locals[i]] = true
	}
	defer func() { lo, locs, lm, ce = save_lo, save_locs, save_lm, save_ce }()

	xn := nn(x)
	p := int32(x)
	e := p + 8*xn
	markwhile(x)
	for p < e {
		u := K(I64(p))
		//println("u", u, int32(u), "tp", tp(u))
		//if p+8 < e {
		//	n := K(I64(p + 8))
		//	println(" n", n, int32(n), "tp", tp(n))
		//}
		if u == 0xffffffffffffffff { // while
			fmt.Fprintf(w, "%s := K(0)\nfor {\n", ssa())
		} else if u == 0xfffffffffffffffe { // jif within while
			fmt.Fprintf(w, "dx(%s)\n", v0())
			fmt.Fprintf(w, "if int32(%s) == 0 { break; }\n", v0())
			lpop()
			printstack(w)
			p += 16
		} else if u == 0xfffffffffffffffd { // jmp offset in while
			r := lo[len(lo)-2]
			fmt.Fprintf(w, "dx(%s); %s = %s\n}\n", r, r, v0())
			lpop()
			printstack(w)
			p += 8
		} else if u == 0xfffffffffffffffc { // return
			s := v0()
			lpop()
			for _, s := range lo {
				fmt.Fprintf(w, "dx(%s)\n", s)
			}
			fmt.Fprintf(w, "return %s\n", s)
			p += 16
		} else if p+8 < e && K(I64(p+8)) == 384 { // jif in cond
			fmt.Fprintf(w, "%s := K(0)\n", ssa())
			printstack(w)
			fmt.Fprintf(w, "dx(%s)\n", v1()) // cond's return var
			fmt.Fprintf(w, "if int32(%s) != 0 {\n", v1())
			ppop()
			printstack(w)
			p += 8
		} else if p+8 < e && K(I64(p+8)) == 320 { // jmp in cond
			ce = append(ce, p+16+int32(u))
			fmt.Fprintf(w, "%s = %s\n", v1(), v0())
			lpop()
			printstack(w)
			fmt.Fprintf(w, "} else {\n")
			p += 8
		} else if tp(u) != 0 { // noun
			if isAsn(u, e, p) {
			} else if isLup(u, e, p) {
				p += 8 // skip .
				s := sK(u)
				if isglobal(s + "_") {
					fmt.Fprintf(w, "%s := Val(Ks(%d)) // .%s\n", ssa(), intern(s), s)
				} else {
					s += "_"
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
				ppop()
				ppop()
				ppop()
				printstack(w)
			case 3: // 192..255  tetradic
				fmt.Fprintf(w, "// tetradic %v\n", int32(u))
				printstack(w)
				s := "Dmd"
				if 211 == int32(u) {
					s = "Amd"
				}
				fmt.Fprintf(w, "%s := %s(%s, %s, %s, %s)\n", ssa(), s, v1(), v2(), vn(3), vn(4))
				ppop()
				ppop()
				ppop()
				ppop()
				printstack(w)
			case 4: // 256       drop
				fmt.Fprintf(w, "dx(%s) //drop\n", v0())
				lpop()
				printstack(w)
			case 5: // 320       jump
				panic("!jump\n")
			case 6: // 384       jump if not
				panic("!jump if not\n")
			default: //448..     quoted verb
				v := int(unquote(K(int32(u))))
				c := byte(' ')
				t := "0:+-*%!&|<>=~,^#_$?@."
				if v > 0 && v < len(t) {
					c = t[v]
				} else if v-64 > 0 && v-64 < len(t) {
					c = t[v-64]
				}
				fmt.Fprintf(w, "%s := K(%d) // `%d (%c)\n", ssa(), v, v, c)
			}
		}
		p += 8

		for len(ce) > 0 && ce[len(ce)-1] == p { // cond-end
			ce = ce[:len(ce)-1]
			fmt.Fprintf(w, "%s = %s\n", lo[len(lo)-2], lo[len(lo)-1])
			lpop()
			printstack(w)
			fmt.Fprintf(w, "}\n")
		}
	}
	if len(lo) == 0 {
		return ""
	}
	return v0()
}
func printstack(w io.Writer) {} // fmt.Fprintf(w, "// lo> %v\n", lo) }
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
		return "konst_" + strconv.Itoa(g)
	} else {
		r := "konst_" + strconv.Itoa(len(konsts))
		konsts[s] = len(konsts)
		return r
	}
}
func isAsn(u K, e, p int32) bool {
	if tp(u) == 4 && e > p+8 {
		v := K(I64(p + 8))
		if tp(v) == 0 && int32(v) == 64 {
			as = append(as, sK(u)+"_")
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
	for _, v := range locs {
		if v == s {
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
		s = strings.TrimSuffix(s, "_")
		fmt.Fprintf(w, "%s := Asn(Ks(%d), %s) // %s:\n", ssa(), intern(s), v1(), s)
		ppop()
	} else {
		fmt.Fprintf(w, "dx(%s); %s = rx(%s)\n", s, s, y)
		lpop()
		lo = append(lo, s)
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
				m++
				for i := sz * n; i < 8*m; i++ { // blank trailing memory
					SetI8(int32(x)+int32(i), 0)
				}
			}
			for i := 0; i < m; i++ {
				fmt.Fprintf(&buf, "SetI64(k_p+%d, int64(%d))\n", 8*i, I64(int32(x)+int32(8*i)))
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
		return fmt.Sprintf("Ks(%d)", intern(s))
	case 5: // f
		return fmt.Sprintf("Kf(%v)", fstr(F64(int32(x))))
	case 6: // z
		return fmt.Sprintf("Kz(%v,%v)", fstr(F64(int32(x))), fstr(F64(int32(x)+8)))
	case 13: // lambda
		return kmLambda(x)
	case 17: // B
		return vec(1)
	case 18: // C
		return fmt.Sprintf("rx(%s)", konst(mkC(string(CK(rx(x))))))
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
		fmt.Fprintf(&buf, "k_p = int32(@)\n")
		if n > 0 {
			for i := int32(0); i < n; i++ {
				y := ati(rx(x), i)
				s := string(CK(Kst(rx(y))))
				fmt.Fprintf(&buf, "SetI64(k_p+%d, int64(%s)) // %s\n", 8*i, kmNoun(y), s)
			}
			//fmt.Fprintf(&buf, "%s = uf(%s)\n", s, s)
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
	if ary > maxargs {
		maxargs = ary
	}
	code := x0(int32(x))
	locs := x1(int32(x))
	kstr := x2(int32(x))
	ls := string(CK(rx(kstr)))
	if s, o := lmbdas[ls]; o {
		return s
	}
	v := SK(locs)
	for i := range v {
		v[i] += "_"
	}
	r := kom(&buf, code, v, ary)

	n := "f_" + strconv.Itoa(len(kfns))
	b := []byte("func " + n + "(" + strings.Join(v[:ary], ", "))
	if ary > 0 {
		b = append(b, ' ', 'K')
	}
	b = append(b, []byte(") K { // "+single(string(CK(rx(kstr))))+"\n")...)
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
	s := fmt.Sprintf("lmb(%d, %d)", tablen+fn, ary)
	lmbdas[ls] = s
	return s
}
func kmMo(w io.Writer, u K) {
	if u == 0 {
		fmt.Fprintf(w, "%s := K(0)\n", ssa())
		return
	}
	fmt.Fprintf(w, "%s := %s(%s)\n", ssa(), mo[int32(u)], v1())
	ppop()
	printstack(w)
}
func kmDy(w io.Writer, u K) {
	if int32(u) == 64 { // asn
		kmAsn(w)
		printstack(w)
	} else {
		fmt.Fprintf(w, "%s := %s(%s, %s)\n", ssa(), dy[int32(u)-64], v1(), v2())
		ppop()
		ppop()
		printstack(w)
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
		fmt.Fprintf(w, "SetI64(int32(%s)+%d, int64(%s))\n", v0(), 8*i, v1())
		ppop()
	}
	if n > 0 {
		fmt.Fprintf(w, "%s = uf(%s)\n", v0(), v0())
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
func lpop()           { lo = lo[:len(lo)-1] }                   // pop last ssa
func ppop()           { lo[len(lo)-2] = lo[len(lo)-1]; lpop() } // pop previous local
func vn(n int) string { return lo[len(lo)-n-1] }                // n'th var name from top
func v0() string      { return vn(0) }
func v1() string      { return vn(1) }
func v2() string      { return vn(2) }
func fstr(f float64) string {
	if f != f {
		return "na"
	} else if math.IsInf(f, 1) {
		return "inf"
	} else if math.IsInf(f, -1) {
		return "-inf"
	} else {
		return strconv.FormatFloat(f, 'g', -1, 64)
	}
}

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
	for i := 0; i < len(konsts); i++ {
		v = append(v, "konst_"+strconv.Itoa(i))
	}
	fmt.Fprintf(w, "var %s K\n", strings.Join(v, ", "))
}
func constnt(w io.Writer) { // initialized in main
	if len(konsts) == 0 {
		return
	}
	fmt.Fprintf(w, "// runtime constants\n")
	ord := make([]string, len(konsts))
	for s, i := range konsts {
		ord[i] = s
	}
	for k, c := range ord {
		s := "konst_" + strconv.Itoa(k)
		fmt.Fprintf(w, "%s = %s\n", s, strings.Replace(c, "@", s, -1))
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
	if len(b) > 0 {
		copy(Bytes[int32(r):], b)
	}
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
func CK(x K) []byte {
	if x == 0 {
		return []byte{}
	}
	r := Bytes[int32(x) : int32(x)+nn(x)]
	rx(x)
	return r
}

const head string = `func main() {
`

func rtext() string {
	s := `
func lmb(x, a int32) K { // store compiled lambda function as a k value (type xf)
	l := l2(Ki(x), mk(Ct, 0))
	SetI32(int32(l)-12, a)         // arity as length
	return K(int32(l)) | K(14)<<59 // native function type xf
}
func kal(f, x K) (r K) { // call kom-compiled lambda function
	n := nn(x)
	xp := int32(x)
	fp := int32(f)
	_, _ = xp, fp
	lfree(x)
	switch(n) {
`
	for i := 0; i <= maxargs; i++ {
		t := strings.Join(strings.Split(strings.Repeat("K", i), ""), ",")
		a := ""
		if i > 0 {
			a = "K(I64(xp))"
			for j := 1; j < i; j++ {
				a += ", K(I64(xp+" + strconv.Itoa(8*j) + "))"
			}
		}
		s += fmt.Sprintf("case %d:\nr=Func[fp].(func(%s) K)(%s)\n", i, t, a)
	}

	return s + `
	default:
		r = trap(Err)
	}
	dx(f)
	return r
}`
}

const replsrc = `
doargs()
write(Ku(2932601077199979))
store()
for {
	write(Ku(32))
	x := read()
	try(x)
}
`

func skiptests() []string {
	return []string{
		"{x+y*z}[;1 2;]",                     // cannot string compiled functions ($l)
		"f:{x+y};f 3",                        // $l
		"({x+y+z}3)4",                        // $l
		"({x+y+z}[;1;])[;2]",                 // $l
		"{[a]a}",                             // $l
		"(.{[]x+y})3",                        // $l
		"(.{a+b})3",                          // .l
		"(.{[]s:1+s:1})1",                    // .l
		"2#{x+y}",                            // $l
		"x:1",                                // prints 1 not null
		"{x+y}",                              // $l
		".{x+y}",                             // .l
		"(.{a:3*x;a:5;x:2})1",                // .l
		"(.{a:3*x})3",                        // .l
		"{a+b}.`a`b!1 2",                     // envcall is not supported (l.d)
		"(+`a`b!(1 2 3;4 5 6)){a>1}",         // envcall is not supported (t@d)
		"`unpack 0x690400000002000000",       // unpack uses dynamic scope
		"x~`unpack `pack x:1.2 3.4",          // unpack
		"x~`unpack `pack x:(1;2 3;`beta)",    // unpack
		"x~`unpack `pack x:`a`b!2 3",         // unpack
		"x~`unpack `pack x:+`a`b!(1 2;3 4.)", // unpack
	}
}

// todo: refcount optimization
//  k$ := rx(?)
//  dx(k$) //drop
// => delete both lines
//  k$ := Asn(?)
//  dx(k$)
// => dx(Asn(?))
