package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Interpreter struct {
	v, e stack // operand, dictionary, execution
	d    []Dictionary
	o    io.Writer
}
type stack struct {
	stack []Value
}

func New(o io.Writer) Interpreter {
	var i Interpreter
	i.d = []Dictionary{mkBuiltins(), make(Dictionary), make(Dictionary)}
	i.o = o
	return i
}
func (i *Interpreter) Push(v Value)        { i.v.Push(v) }
func (i *Interpreter) Pop() (r Value)      { return i.v.Pop() }
func (i *Interpreter) err(e string)        { panic(e) }
func (i *Interpreter) Exec(k *Interpreter) { i.Push(k) }
func (i *Interpreter) String() string      { return "save" }
func (i *Interpreter) Clone() Value { // the interpreter is a Value itself (for save/restore)
	k := New(i.o)
	k.v = i.v.clone()
	k.e = i.e.clone()
	return &k
}

func (s *stack) Push(v Value) { s.stack = append(s.stack, v) }
func (s *stack) Pop() (r Value) {
	r = s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return r
}
func (s stack) clone() (r stack) {
	r.stack = make([]Value, len(s.stack))
	for i, v := range s.stack {
		r.stack[i] = v.Clone()
	}
	return r
}

type Value interface {
	Exec(i *Interpreter)
	Clone() Value // deep copy
	String() string
}
type Quoter interface {
	Quote() string
}

type (
	Boolean    bool
	Integer    int
	Real       float64
	Complex    complex128
	Mark       string
	Name       string
	Null       bool
	Operator   func(*Interpreter)
	Array      []Value
	String     string
	Dictionary map[Value]Value
)

func (b Boolean) Exec(i *Interpreter) { i.Push(b) }
func (b Boolean) String() string      { return strconv.FormatBool(bool(b)) }
func (b Boolean) Clone() Value        { return b }
func (n Integer) Exec(i *Interpreter) { i.Push(n) }
func (n Integer) String() string      { return strconv.Itoa(int(n)) }
func (n Integer) Clone() Value        { return n }
func (r Real) Exec(i *Interpreter)    { i.Push(r) }
func (r Real) String() string         { return strconv.FormatFloat(float64(r), 'g', -1, 64) }
func (r Real) Clone() Value           { return r }
func (z Complex) Exec(i *Interpreter) { i.Push(z) }
func (z Complex) String() string {
	r, phi := cmplx.Polar(complex128(z))
	phi *= 180.0 / math.Pi
	if phi < 0 {
		phi += 360.0
	}
	if r == 0.0 {
		phi = 0.0 // We want predictable angles in this case.
	}
	if phi == -0.0 || phi == 360.0 {
		phi = 0.0
	}
	ang := fmt.Sprintf("%.1f", phi)
	if strings.HasSuffix(ang, ".0") {
		ang = ang[:len(ang)-2]
	}
	return fmt.Sprintf("%v@%s", r, ang)
}
func (z Complex) Clone() Value         { return z }
func (m Mark) Exec(i *Interpreter)     { i.Push(m) }
func (m Mark) String() string          { return string(m) }
func (m Mark) Clone() Value            { return m }
func (n Name) Exec(i *Interpreter)     { i.Push(n) } // todo lookup
func (n Name) String() string          { return string(n) }
func (n Name) Clone() Value            { return n }
func (o Operator) Exec(i *Interpreter) { o(i) }
func (o Operator) String() string      { return runtime.FuncForPC(reflect.ValueOf(o).Pointer()).Name() }
func (o Operator) Clone() Value        { return o }
func (a Array) Exec(i *Interpreter)    { i.Push(a) }
func (a Array) String() string         { return fmt.Sprintf("%v", []Value(a)) }
func (a Array) Clone() Value {
	r := make(Array, len(a))
	for i, v := range a {
		r[i] = v.Clone()
	}
	return r
}
func (s String) Exec(i *Interpreter)     { i.Push(s) }
func (s String) String() string          { return string(s) }
func (s String) Quote() string           { return "(" + string(quote(string(s))) + ")" }
func (s String) Clone() Value            { return s }
func (d Dictionary) Exec(i *Interpreter) { i.Push(d) }
func (d Dictionary) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "<<")
	for k, v := range d {
		fmt.Fprintf(&b, "%s %s ", k, v)
	}
	fmt.Fprintf(&b, ">>")
	return b.String()
}
func (d Dictionary) Clone() Value { panic("nyi"); return d }

// stack operators
func pop(i *Interpreter) { _ = i.Pop() }
func exch(i *Interpreter) {
	n := len(i.v.stack)
	i.v.stack[n-1], i.v.stack[n-2] = i.v.stack[n-2], i.v.stack[n-1]
}
func dup(i *Interpreter) { x := i.Pop(); i.Push(x); i.Push(x) }
func _copy(i *Interpreter) {
	n := i.Pop().(Integer)
	if n < 0 {
		i.err("range")
	}
	o := len(i.v.stack) - int(n)
	for k := 0; k < int(n); k++ {
		i.Push(i.v.stack[o+k])
	}
}
func index(i *Interpreter) { cvi(i); i.Push(i.v.stack[len(i.v.stack)-int(i.Pop().(Integer))-1]) }
func roll(i *Interpreter) {
	k := int(i.Pop().(Integer))
	n := int(i.Pop().(Integer))
	if n < 0 {
		i.err("range")
	}
	k %= n
	v := make([]Value, n)
	copy(v, i.v.stack[len(i.v.stack)-n:])
	i.v.stack = i.v.stack[:len(i.v.stack)-n]
	for j := 0; j < n; j++ {
		i.Push(v[(n+j-k)%n])
	}
}
func clear(i *Interpreter) { i.v.stack = i.v.stack[:0] }
func count(i *Interpreter) { i.Push(Integer(len(i.v.stack))) }
func mark(i *Interpreter)  {}
func cleartomark(i *Interpreter) {
	counttomark(i)
	n := int(i.Pop().(Integer))
	i.v.stack = i.v.stack[:len(i.v.stack)-n-1]
}
func counttomark(i *Interpreter) {
	l := len(i.v.stack)
	for k := 0; k < l; k++ {
		if _, ok := i.v.stack[l-k-1].(Mark); ok {
			i.Push(Integer(k))
			return
		}
	}
	i.err("unmatchedmark")
}
func array(i *Interpreter) {
	n := int(i.Pop().(Integer))
	l := len(i.v.stack)
	if n < 0 {
		i.err("range")
	} else if n > l {
		i.err("stack")
	}
	a := make(Array, n)
	copy(a, i.v.stack[l-n:])
	i.v.stack = i.v.stack[:l-1-n]
	i.Push(a)
}

// arithmetic operators
func add(i *Interpreter) {
	numOp2(i, 1, 0, func(x, y int) int { return x + y }, func(x, y float64) float64 { return x + y }, func(x, y complex128) complex128 { return x + y })
}
func sub(i *Interpreter) {
	numOp2(i, 1, 0, func(x, y int) int { return x - y }, func(x, y float64) float64 { return x - y }, func(x, y complex128) complex128 { return x - y })
}
func mul(i *Interpreter) {
	numOp2(i, 1, 0, func(x, y int) int { return x * y }, func(x, y float64) float64 { return x * y }, func(x, y complex128) complex128 { return x * y })
}
func div(i *Interpreter) {
	numOp2(i, 2, 0, func(x, y int) int { return x / y }, func(x, y float64) float64 { return x / y }, func(x, y complex128) complex128 { return x / y })
}
func idiv(i *Interpreter) {
	numOp2(i, 0, 1, func(x, y int) int { return x / y }, nil, nil)
}
func mod(i *Interpreter) {
	numOp2(i, 0, 1, func(x, y int) int { return x % y }, func(x, y float64) float64 { return math.Mod(x, y) }, nil)
}
func abs(i *Interpreter) {
	x := i.Pop()
	switch v := x.(type) {
	case Integer:
		if v < 0 {
			v = -v
		}
		i.Push(v)
	case Real:
		i.Push(Real(math.Abs(float64(v))))
	case Complex:
		i.Push(Real(cmplx.Abs(complex128(v))))
	default:
		i.err("type")
	}
}
func neg(i *Interpreter) {
	numOp1(i, 0, 0, func(x int) int { return -x }, func(x float64) float64 { return -x }, func(x complex128) complex128 { return -x })
}
func ceiling(i *Interpreter) {
	numOp1(i, 0, 2, func(x int) int { return x }, func(x float64) float64 { return math.Ceil(x) }, nil)
}
func floor(i *Interpreter) {
	numOp1(i, 0, 2, func(x int) int { return x }, func(x float64) float64 { return math.Floor(x) }, nil)
}
func round(i *Interpreter) {
	numOp1(i, 0, 2, func(x int) int { return x }, func(x float64) float64 { return math.Round(x) }, nil)
}
func truncate(i *Interpreter) {
	numOp1(i, 0, 2, func(x int) int { return x }, func(x float64) float64 { return math.Trunc(x) }, nil)
}
func sqrt(i *Interpreter) {
	numOp1(i, 2, 2, nil, func(x float64) float64 { return math.Sqrt(x) }, nil)
}
func atan(i *Interpreter) {
	numOp2(i, 2, 2, nil, func(x, y float64) float64 { return math.Atan2(x, y) }, nil)
}
func cos(i *Interpreter) {
	numOp1(i, 2, 2, nil, func(x float64) float64 { return math.Cos(x) }, nil)
}
func sin(i *Interpreter) {
	numOp1(i, 2, 2, nil, func(x float64) float64 { return math.Sin(x) }, nil)
}
func exp(i *Interpreter) { // pow
	numOp2(i, 2, 2, nil, func(x, y float64) float64 { return math.Pow(x, y) }, nil)
}
func ln(i *Interpreter) { // log base e
	numOp1(i, 2, 2, nil, func(x float64) float64 { return math.Log(x) }, nil)
}
func log(i *Interpreter) { // log base 10
	numOp1(i, 2, 2, nil, func(x float64) float64 { return math.Log10(x) }, nil)
}
func _rand(i *Interpreter) { i.Push(Integer(rand.Int())) }
func srand(i *Interpreter) {
	x := i.Pop()
	if numType(x) != 1 {
		i.err("type")
	}
	rand.Seed(int64(x.(Integer)))
}
func eq(i *Interpreter) { x, y, _ := numTp2(i, 0, 0); i.Push(Boolean(x == y)) }
func numType(v Value) int {
	switch v.(type) {
	case Integer:
		return 1
	case Real:
		return 2
	case Complex:
		return 3
	default:
		return 0
	}
}
func upNum(a, b Value) (Value, Value) {
	max := func(x, y int) int {
		if x > y {
			return x
		} else {
			return y
		}
	}
	if at, bt := numType(a), numType(b); at > 0 && bt > 0 {
		a, _ = uptype(a, max(at, bt))
		b, _ = uptype(b, max(at, bt))
		return a, b
	}
	return a, b
}
func uptype(v Value, t int) (Value, int) {
	if t == 1 {
		return Real(v.(Integer)), 2
	} else if t == 2 {
		return Complex(complex(v.(Real), 0)), 3
	}
	panic("unreachable")
}
func numOp1(i *Interpreter, minType, maxType int, fi func(x int) int, fr func(x float64) float64, fz func(x complex128) complex128) {
	x := i.Pop()
	xt := numType(x)
	if xt == 0 {
		i.err("type")
	}
	for xt < minType {
		x, xt = uptype(x, xt)
	}
	if maxType > 0 && xt > maxType {
		i.err("type")
	}
	switch xt {
	case 1:
		i.Push(Integer(fi(int(x.(Integer)))))
	case 2:
		i.Push(Real(fr(float64(x.(Real)))))
	case 3:
		i.Push(Complex(fz(complex128(x.(Complex)))))
	}
}
func numTp2(i *Interpreter, minType, maxType int) (x, y Value, t int) {
	y = i.Pop()
	x = i.Pop()
	xt, yt := numType(x), numType(y)
	if xt*yt == 0 {
		return x, y, 0
	}
	for xt < yt {
		x, xt = uptype(x, xt)
	}
	for yt < xt {
		y, yt = uptype(y, yt)
	}
	for xt < minType {
		x, xt = uptype(x, xt)
		y, yt = uptype(y, yt)
	}
	t = xt
	if maxType > 0 && xt > maxType {
		t = 0
	}
	return x, y, t
}
func numOp2(i *Interpreter, minType, maxType int, fi func(x, y int) int, fr func(x, y float64) float64, fz func(x, y complex128) complex128) {
	x, y, t := numTp2(i, minType, maxType)
	switch t {
	case 0:
		i.err("type")
	case 1:
		i.Push(Integer(fi(int(x.(Integer)), int(y.(Integer)))))
	case 2:
		i.Push(Real(fr(float64(x.(Real)), float64(y.(Real)))))
	case 3:
		i.Push(Complex(fz(complex128(x.(Complex)), complex128(y.(Complex)))))
	}
}

func cvi(i *Interpreter) {
	x := i.Pop()
	switch v := x.(type) {
	case Integer:
		i.Push(x)
	case Real:
		i.Push(Integer(v))
	default:
		i.err("type")
	}
}

func (i *Interpreter) Run(s string) {
	if strings.HasPrefix(s, "<<") {
		s = strings.Replace(s, "<<", "«", 1)
	}
	s = strings.Replace(s, " <<", " «", -1)
	s = strings.Replace(s, ">> ", "» ", -1)
	token, b := []rune{}, []rune(s)
	for {
		token, b = i.Token(b)
		if len(token) > 0 {
			if v := i.parse(string(token)); v != nil {
				v.Exec(i)
			}
		}
		if len(b) == 0 {
			return
		}
	}
}
func (i *Interpreter) Token(b []rune) (token, tail []rune) {
	isSpace := func(r rune) bool {
		if r <= '\u00FF' {
			switch r {
			case ' ', '\t', '\n', '\v', '\f', '\r', '\u0085', '\u00A0':
				return true
			}
			return false
		}
		if '\u2000' <= r && r <= '\u200a' {
			return true
		}
		switch r {
		case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
			return true
		}
		return false
	}
	adv := func(v []rune) []rune {
		for len(v) > 0 && isSpace(v[0]) {
			v = v[1:]
		}
		return v
	}
	str := func(s []rune) (b, t []rune) {
		if len(s) == 0 {
			i.err("parse")
		}
		q := false
		for i, r := range s {
			if r == '\\' {
				q = !q
			} else if !q && r == ')' {
				return b, s[i+1:]
			}
			b = append(b, r)
		}
		i.err("parse-string")
		return b, t
	}
	if len(b) == 0 {
		return nil, nil
	}
	b = adv(b)
	if len(b) > 0 && b[0] == '[' {
		return b[:1], b[1:]
	}
	for i, r := range b {
		if i == 0 {
			if r == '(' {
				b, tail = str(b)
				break
			} else if r == '[' || r == '{' || r == '«' {
				return b[:1], adv(b[1:])
			}
		}
		if isSpace(r) {
			tail = b[i:]
			b = b[:i]
			break
		}
	}
	tail = adv(tail)
	if len(b) > 1 {
		if t := b[len(b)-1]; t == ']' || t == '}' || t == '»' {
			b = b[:len(b)-1]
			tail = append([]rune{t, ' '}, tail...)
		}
	}
	if len(b) > 0 && b[0] == '%' {
		return nil, nil
	}
	return b, tail
}
func quote(s string) (r []rune) {
	for _, v := range s {
		switch v {
		case '\n':
			v = 'n'
		case '\r':
			v = 'r'
		case '\t':
			v = 't'
		case '\\':
		case ')':
		default:
			r = append(r, v)
			continue
		}
		r = append(r, '\\')
		r = append(r, v)
	}
	return r
}
func unquote(s string) string {
	var r []rune
	q := false
	for _, v := range s {
		if v == '\\' {
			q = !q
			if q {
				continue
			}
		}
		if q {
			switch v {
			case 'n':
				v = '\n'
			case 'r':
				v = '\r'
			case 't':
				v = '\t'
			}
		}
		r = append(r, v)
	}
	return string(r)
}
func (i *Interpreter) parse(s string) Value {
	if s == "true" {
		return Boolean(true)
	} else if s == "false" {
		return Boolean(false)
	}
	if i, e := strconv.Atoi(s); e == nil {
		return Integer(i)
	}
	if f, e := strconv.ParseFloat(s, 64); e == nil {
		return Real(f)
	}
	if strings.HasPrefix(s, "(") {
		return String(unquote(s[1 : len(s)-2]))
	}
	if i := strings.Index(s, "a"); i > 0 {
		if abs, e := strconv.ParseFloat(s[:i], 64); e == nil {
			if ang, e := strconv.ParseFloat(s[:i], 64); e == nil {
				return Complex(cmplx.Rect(abs, math.Pi*ang/180.0))
			}
		}
	}
	name := Name(s)
	d := i.where(name)
	if d != nil {
		return d[name]
	}
	i.err("/undefined in " + s)
	return nil
}
func (i *Interpreter) where(v Value) Dictionary {
	for n := len(i.d) - 1; n >= 0; n-- {
		d := i.d[n]
		if _, ok := d[v]; ok {
			return d
		}
	}
	return nil
}
func mkBuiltins() Dictionary {
	return Dictionary{
		Name("pop"):         Operator(pop),
		Name("exch"):        Operator(exch),
		Name("dup"):         Operator(dup),
		Name("copy"):        Operator(_copy),
		Name("index"):       Operator(index),
		Name("roll"):        Operator(roll),
		Name("clear"):       Operator(clear),
		Name("count"):       Operator(count),
		Name("mark"):        Operator(func(i *Interpreter) { i.Push(Mark("mark")) }),
		Name("["):           Operator(func(i *Interpreter) { i.Push(Mark("[")) }),
		Name("]"):           Operator(func(i *Interpreter) { counttomark(i); array(i) }),
		Name("{"):           Operator(func(i *Interpreter) { i.Push(Mark("{")) }),
		Name("cleartomark"): Operator(cleartomark),
		Name("counttomark"): Operator(counttomark),

		Name("add"):      Operator(add),
		Name("div"):      Operator(div),
		Name("idiv"):     Operator(idiv),
		Name("mod"):      Operator(mod),
		Name("mul"):      Operator(mul),
		Name("sub"):      Operator(sub),
		Name("abs"):      Operator(abs),
		Name("neg"):      Operator(neg),
		Name("ceiling"):  Operator(ceiling),
		Name("floor"):    Operator(floor),
		Name("round"):    Operator(round),
		Name("truncate"): Operator(truncate),
		Name("sqrt"):     Operator(sqrt),
		Name("atan"):     Operator(atan),
		Name("cos"):      Operator(cos),
		Name("sin"):      Operator(sin),
		Name("exp"):      Operator(exp),
		Name("ln"):       Operator(ln),
		Name("log"):      Operator(log),
		Name("rand"):     Operator(_rand),
		Name("srand"):    Operator(srand), // no rrand

		Name("eq"): Operator(eq),

		Name("stack"):  Operator(pstack), // we only have pstack
		Name("pstack"): Operator(pstack),
		Name("="):      Operator(_print),
		Name("=="):     Operator(__print),
	}
}
func pstack(i *Interpreter) {
	for n := len(i.v.stack) - 1; n >= 0; n-- {
		fmt.Fprintln(i.o, i.v.stack[n].String())
	}
}

func _print(i *Interpreter) { v := i.Pop(); fmt.Fprintf(i.o, "%s\n", v) } // =, stack
func __print(i *Interpreter) { // ==, pstack
	v := i.Pop()
	if q, o := v.(Quoter); o {
		v = String(q.Quote())
		fmt.Fprintf(i.o, "%s\n", v)
	}
}
func (i *Interpreter) prompt() {
	if n := len(i.v.stack); n > 0 {
		fmt.Printf("PS<%d>", n)
	} else {
		fmt.Printf("PS>")
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	i := New(os.Stdout)
	if len(os.Args) > 1 {
		i.Run(strings.Join(os.Args[1:], " "))
		return
	}
	i.prompt()
	for s.Scan() {
		i.Run(s.Text())
		i.prompt()
	}
}
