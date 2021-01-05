package main

import (
	"bufio"
	"fmt"
	"math"
	"math/cmplx"
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
	Mark       bool
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
func numOp2(i *Interpreter, fi func(x, y int) int, fr func(x, y float64) float64, fz func(x, y complex128) complex128) {
	y := i.Pop()
	x := i.Pop()
	xt, yt := numType(x), numType(y)
	if xt*yt == 0 {
		i.err("not-numeric")
	}
	for xt < yt {
		x, xt = uptype(x, xt)
	}
	for yt < xt {
		y, yt = uptype(y, yt)
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
func add(i *Interpreter) {
	numOp2(i, func(x, y int) int { return x + y }, func(x, y float64) float64 { return x + y }, func(x, y complex128) complex128 { return x + y })
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
	d := make(Dictionary)
	d[Name("stack")] = Operator(pstack) // we only have pstack
	d[Name("pstack")] = Operator(pstack)
	d[Name("pop")] = Operator(pop)
	d[Name("exch")] = Operator(exch)
	d[Name("=")] = Operator(_print)
	d[Name("==")] = Operator(_print)
	return d
}
func pstack(i *Interpreter) {
	for n := len(i.v.stack) - 1; n >= 0; n-- {
		fmt.Println(i.v.stack[n].String())
	}
	fmt.Println() // to separate multiple calls. gs does not do this.
}
func exch(i *Interpreter) {
	n := len(i.v.stack)
	i.v.stack[n-1], i.v.stack[n-2] = i.v.stack[n-2], i.v.stack[n-1]
}
func pop(i *Interpreter)    { _ = i.Pop() }
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
	i.prompt()
	for s.Scan() {
		i.Run(s.Text())
		i.prompt()
	}
}
