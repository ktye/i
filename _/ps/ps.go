package main

import (
	"bufio"
	"fmt"
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
}
type stack struct {
	stack []Value
}

func New() Interpreter {
	var i Interpreter
	i.d = []Dictionary{mkBuiltins(), make(Dictionary), make(Dictionary)}
	return i
}
func (i *Interpreter) Push(v Value)   { i.v.Push(v) }
func (i *Interpreter) Pop() (r Value) { return i.v.Pop() }
func (i *Interpreter) err(e string)   { panic(e) }

func (s *stack) Push(v Value) { s.stack = append(s.stack, v) }
func (s *stack) Pop() (r Value) {
	r = s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return r
}

type Value interface {
	Exec(i *Interpreter)
	String() string
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
	Dictionary map[Value]Value
)

func (b Boolean) Exec(i *Interpreter) { i.Push(b) }
func (b Boolean) String() string      { return strconv.FormatBool(bool(b)) }
func (n Integer) Exec(i *Interpreter) { i.Push(n) }
func (n Integer) String() string      { return strconv.Itoa(int(n)) }
func (r Real) Exec(i *Interpreter)    { i.Push(r) }
func (r Real) String() string         { return strconv.FormatFloat(float64(r), 'g', -1, 64) }
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
func (m Mark) Exec(i *Interpreter)       { i.Push(m) }
func (m Mark) String() string            { return string(m) }
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
func (n Name) Exec(i *Interpreter)     { i.Push(n) } // todo lookup
func (n Name) String() string          { return string(n) }
func (o Operator) Exec(i *Interpreter) { o(i) }
func (o Operator) String() string      { return runtime.FuncForPC(reflect.ValueOf(o).Pointer()).Name() }

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
func mark(i *Interpreter)  { i.Push(Mark("mark")) }
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
func numOp2(i *Interpreter, minType, maxType int, fi func(x, y int) int, fr func(x, y float64) float64, fz func(x, y complex128) complex128) {
	y := i.Pop()
	x := i.Pop()
	xt, yt := numType(x), numType(y)
	if xt*yt == 0 {
		i.err("type")
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
	if maxType > 0 && xt > maxType {
		i.err("type")
	}
	switch xt {
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
	if len(b) == 0 {
		return nil, nil
	}
	b = adv(b)
	for i, r := range b {
		if r == '(' {
			panic("todo token string")
		}
		if isSpace(r) {
			tail = b[i:]
			b = b[:i]
			break
		}
	}
	tail = adv(tail)
	return b, tail
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
		Name("mark"):        Operator(mark),
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
		Name("srand"):    Operator(srand),

		Name("stack"):  Operator(pstack), // we only have pstack
		Name("pstack"): Operator(pstack),
		Name("="):      Operator(_print),
		Name("=="):     Operator(_print),
	}
}
func pstack(i *Interpreter) {
	for n := len(i.v.stack) - 1; n >= 0; n-- {
		fmt.Println(i.v.stack[n].String())
	}
	fmt.Println() // to separate multiple calls. gs does not do this.
}

func _print(i *Interpreter) { v := i.Pop(); fmt.Printf("%s\n", v) }
func (i *Interpreter) prompt() {
	if n := len(i.v.stack); n > 0 {
		fmt.Printf("PS<%d>", n)
	} else {
		fmt.Printf("PS>")
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	i := New()
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
