package k

import (
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
)

type Symbol string
type Quotation list
type value = interface{}
type list = []interface{}
type Interpreter struct {
	Var map[Symbol]value
	EX  []list
	PC  []int
	ST  list
}

//go:embed j.j
var j []byte

func New() *Interpreter {
	var i Interpreter
	v := map[Symbol]value{
		":":      i.asn,
		"if":     i.iff,
		"pop":    i.pop,
		"dup":    i.dup,
		"dip":    i.dip,
		"swap":   i.swap,
		"roll":   i.roll,
		"exec":   i.exec,
		"return": i.ret,
		"stack":  i.stack,
		"trap":   i.trap,
		"+":      i.add,
	}
	i.Var = v
	i.Run(Parse(j))
	return &i
}

func Parse(b []byte) (r list) {
	v := bytes.Fields(b)
	for _, si := range v {
		s := string(si)
		if n, e := strconv.Atoi(s); e == nil {
			r = append(r, n)
			continue
		}
		if n, e := strconv.ParseFloat(s, 64); e == nil {
			r = append(r, n)
			continue
		}
		if s[0] == '`' {
			r = append(r, s[1:])
			continue
		}
		r = append(r, Symbol(s))
	}
	i := 0
	for i < len(r) {
		if s, o := r[i].(Symbol); o && s == "]" {
			m := counttomark(r, i-1)
			q := make(Quotation, len(r[m+1:i]))
			copy(q, r[m+1:i])
			r[m] = q
			r = append(r[:1+m], r[i+1:]...)
			i = m + 1
		} else {
			i++
		}
	}
	return r
}
func counttomark(r list, n int) int {
	for i := n; i >= 0; i-- {
		if s, o := r[i].(Symbol); o && s == "[" {
			return i
		}
	}
	panic("unbalanced ]")
}
func (i *Interpreter) Run(x list) value {
	defer func() {
		if r := recover(); r != nil {
			i.trace()
		}
	}()
	i.EX = append(i.EX, x)
	i.PC = append(i.PC, 0)
	for {
		// return all values, if the program stack is empty
		if len(i.EX) == 0 {
			if len(i.ST) == 0 {
				return nil
			}
			r := i.ST
			i.ST = nil
			return r // or only top? error if #>1 ?
		}

		// pop program stack, if the top program is at the end
		ex := i.EX[len(i.EX)-1]
		c := &i.PC[len(i.PC)-1]
		if len(ex) == *c {
			i.EX = i.EX[:len(i.EX)-1]
			i.PC = i.PC[:len(i.PC)-1]
			continue
		}

		// execute one program step and advance program counter
		p := ex[*c]
		*c++

		if s, o := p.(Symbol); o {
			if v, o := i.Var[s]; o == false {
				panic(s + " undefined")
			} else {
				p = v
			}
			// quotations stored in a variable are executed on lookup
			if q, o := p.(Quotation); o {
				i.ST = append(i.ST, q)
				i.exec()
				continue
			}
		}

		switch v := p.(type) {
		case func():
			v()
		default:
			i.ST = append(i.ST, v)
		}
	}
}
func boolean(v value) bool {
	switch x := v.(type) {
	case bool:
		return x
	case int:
		return x != 0
	case float64:
		return x != 0
	default:
		panic(fmt.Errorf("not a boolean: %T", v))
	}
}
func (i *Interpreter) need(n int) (r int) {
	r = len(i.ST)
	if r < n {
		panic("underflow")
	}
	return r - 1
}
func (i *Interpreter) asn() {
	n := i.need(2)
	if s, o := i.ST[n-1].(string); o == false {
		panic(fmt.Errorf("assign: name is not a string: %T", i.ST[n-1]))
	} else {
		i.Var[Symbol(s)] = i.ST[n]
	}
	i.ST = i.ST[:n-1]
}
func (i *Interpreter) get() (r value) {
	n := i.need(1)
	r = i.ST[n]
	i.ST = i.ST[:n]
	return r
}
func (i *Interpreter) pop() {
	n := i.need(1)
	i.ST = i.ST[:n]
}
func (i *Interpreter) swap() {
	n := i.need(2)
	i.ST[n-1], i.ST[n] = i.ST[n], i.ST[n-1]
}
func (i *Interpreter) dup() {
	n := i.need(1)
	i.ST = append(i.ST, i.ST[n])
}
func (i *Interpreter) dip() {
	i.swap()
	v := i.get()
	i.EX = append(i.EX, Quotation{v})
	i.PC = append(i.PC, 0)
	i.exec()
}
func (i *Interpreter) roll() {
	n := i.need(3)
	i.ST[n-2], i.ST[n-1], i.ST[n] = i.ST[n], i.ST[n-2], i.ST[n-1]

}
func (i *Interpreter) iff() {
	n := i.need(2)
	b := boolean(i.ST[n-1])
	i.swap()
	i.pop()
	if b {
		i.exec()
	} else {
		i.pop()
	}
}
func (i *Interpreter) ret() {
	n := len(i.EX)
	i.EX = i.EX[:n-2]
	i.PC = i.PC[:n-2]
}
func (i *Interpreter) exec() {
	x := i.get()
	n := len(i.EX) - 1
	if i.PC[n] == len(i.EX[n]) { // save frame
		i.EX[n] = list(x.(Quotation))
		i.PC[n] = 0
	} else {
		i.EX = append(i.EX, list(x.(Quotation)))
		i.PC = append(i.PC, 0)
	}
}
func (i *Interpreter) stack() { fmt.Println(i.ST) }
func (i *Interpreter) trap()  { panic("trap") }
func (i *Interpreter) trace() {
	fmt.Println("EX")
	for j, ex := range i.EX {
		n := i.PC[j] - 1
		for k := range ex {
			if n == k {
				fmt.Printf(" >%v<", ex[k])
			} else {
				fmt.Printf(" %v", ex[k])
			}
		}
		fmt.Println()
	}
	fmt.Println("ST")
	for j := len(i.ST) - 1; j >= 0; j-- {
		if j == len(i.ST)-5 {
			fmt.Println("..")
			break
		}
		fmt.Printf("%2d %v\n", j, i.ST[j])
	}
}
func (i *Interpreter) add() {
	n := i.need(2)
	i.ST[n-1] = i.ST[n-1].(int) + i.ST[n].(int)
	i.pop()
}
