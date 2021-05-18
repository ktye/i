package k

import (
	"encoding/json"
	"reflect"
)

func New() *K {
	k := K{}
	k.Var = make(map[string]T)
	k.Func = make(map[string]Verb)
	f := func(s string, f1 func(T) T, f2 func(T, T) T, fi func(int, int) int) {
		k.Func[s] = Verb{stringer{s}, f1, f2}
	}
	f("+", flip, add, addi)          //r
	f("-", neg, sub, nil)            //r
	f("*", first, mul, muli)         //r
	f("%", sqrt, div, divi)          //r
	f("&", where, min, mini)         //r
	f("|", reverse, max, maxi)       //r
	f("!", k.til, dict, nil)         //r
	f("~", k.not, match, nil)        //r
	f(",", enlist, cat, nil)         //r
	f("^", sort, k.cut, nil)         //r
	f("=", group, equal, nil)        //r
	f("<", gradeup, less, nil)       //r
	f(">", gradedown, more, nil)     //r
	f("#", count, take, nil)         //r
	f("_", k.floor, k.drop, nil)     //r
	f("@", typ, k.atx, nil)          //r
	f("?", uniq, fnd, nil)           //r
	f("$", k.str, k.cast, nil)       //r
	f(".", k.val, k.call, nil)       //r
	f("abs", abs, nil, nil)          //r
	f("angle", angle, rotate, nil)   //r
	f("real", zreal, nil, nil)       //r
	f("imag", zimag, complx, nil)    //r
	f("conj", conj, nil, nil)        //r
	f("read", read, readdir, nil)    //r
	f("csv", csvwrite, csvread, nil) //r
	f("solve", qr, solve, nil)       //r
	f("\\", debug, debug2, nil)
	f("/", regex, nil, nil)
	f(":", identity, dex, nil)
	// "'", each, each2,             //r
	// "/", over, over2,             //r
	// "\\", scan, scan2,            //r
	// "':", pairs, pairs2,          //r
	// "/:", fix, eachright,         //r
	// "\\:", scanfix, eachleft,     //r
	return &k
}

type K struct {
	Var   map[string]T
	Func  map[string]Verb
	stack []T
	Trap  bool
	ctx
}

type T interface{}

type F1 interface {
	call1(x T) T
}
type F2 interface {
	call2(x, t T) T
}
type F12 interface {
	F1
	F2
	String() string
}
type f1 func(T) T
type f2 func(T, T) T
type Call1 struct{}
type Call2 struct{}
type Call4 bool
type Link struct{}

func (f f1) call1(x T) T    { return f(x) }
func (f f2) call2(x, y T) T { return f(x, y) }

type Verb struct {
	stringer
	f1
	f2
}
type derived struct {
	stringer
	f1
	f2
}
type train []F12

type vector interface {
	refcounter
	zero() T
	ln() int
	at(int) T
	atv([]int) vector
	set(int, T)
	setv([]int, vector)
}
type refcounter interface {
	ref()
	unref() int
	refs() int
}

type B struct {
	refcount
	v []bool
}
type C struct {
	refcount
	v []byte
}
type I struct {
	refcount
	v []int
}
type F struct {
	refcount
	v []float64
}
type Z struct {
	refcount
	v []complex128
}
type S struct {
	refcount
	v []string
}
type L struct {
	refcount
	v []T
}
type D struct {
	refcount
	tab  bool
	k, v vector
}

type refcount struct{ rc *int }

func (r refcount) ref() { *r.rc++ }
func (r refcount) unref() int {
	*r.rc--
	if *r.rc < 0 {
		panic("unref")
	}
	return *r.rc
}
func (r refcount) refs() int { return *r.rc }
func (r *refcount) init()    { i := 1; r.rc = &i }
func (b B) zero() T          { return false }
func (c C) zero() T          { return byte(0) }
func (i I) zero() T          { return 0 }
func (f F) zero() T          { return 0.0 }
func (z Z) zero() T          { return complex(0, 0) }
func (s S) zero() T          { return "" }
func (l L) zero() T          { return nil }
func (b B) ln() int          { return len(b.v) }
func (c C) ln() int          { return len(c.v) }
func (i I) ln() int          { return len(i.v) }
func (f F) ln() int          { return len(f.v) }
func (z Z) ln() int          { return len(z.v) }
func (s S) ln() int          { return len(s.v) }
func (l L) ln() int          { return len(l.v) }
func (b B) at(i int) T       { return b.v[i] }
func (c C) at(i int) T       { return c.v[i] }
func (i I) at(j int) T       { return i.v[j] }
func (f F) at(i int) T       { return f.v[i] }
func (z Z) at(i int) T       { return z.v[i] }
func (s S) at(i int) T       { return s.v[i] }
func (l L) at(i int) T       { return rx(l.v[i]) }
func (b B) atv(v []int) vector {
	r := make([]bool, len(v))
	for i, k := range v {
		r[i] = b.v[k]
	}
	return KB(r)
}
func (c C) atv(v []int) vector {
	r := make([]byte, len(v))
	for i, k := range v {
		r[i] = c.v[k]
	}
	return KC(r)
}
func (j I) atv(v []int) vector {
	r := make([]int, len(v))
	for i, k := range v {
		r[i] = j.v[k]
	}
	return KI(r)
}
func (f F) atv(v []int) vector {
	r := make([]float64, len(v))
	for i, k := range v {
		r[i] = f.v[k]
	}
	return KF(r)
}
func (z Z) atv(v []int) vector {
	r := make([]complex128, len(v))
	for i, k := range v {
		r[i] = z.v[k]
	}
	return KZ(r)
}
func (s S) atv(v []int) vector {
	r := make([]string, len(v))
	for i, k := range v {
		r[i] = s.v[k]
	}
	return KS(r)
}
func (l L) atv(v []int) vector {
	r := make([]T, len(v))
	for i, k := range v {
		r[i] = rx(l.v[k])
	}
	return KL(r)
}
func (b B) set(i int, u T) { b.v[i] = u.(bool) }
func (c C) set(i int, u T) { c.v[i] = u.(byte) }
func (j I) set(i int, u T) { j.v[i] = u.(int) }
func (f F) set(i int, u T) { f.v[i] = u.(float64) }
func (z Z) set(i int, u T) { z.v[i] = u.(complex128) }
func (s S) set(i int, u T) { s.v[i] = u.(string) }
func (l L) set(i int, u T) { rx(l.v[i]); l.v[i] = u }
func (b B) setv(i []int, u vector) {
	v := u.(B).v
	for k, j := range i {
		b.v[j] = v[k]
	}
}
func (c C) setv(i []int, u vector) {
	v := u.(C).v
	for k, j := range i {
		c.v[j] = v[k]
	}
}
func (n I) setv(i []int, u vector) {
	v := u.(I).v
	for k, j := range i {
		n.v[j] = v[k]
	}
}
func (f F) setv(i []int, u vector) {
	v := u.(F).v
	for k, j := range i {
		f.v[j] = v[k]
	}
}
func (z Z) setv(i []int, u vector) {
	v := u.(Z).v
	for k, j := range i {
		z.v[j] = v[k]
	}
}
func (s S) setv(i []int, u vector) {
	v := u.(S).v
	for k, j := range i {
		s.v[j] = v[k]
	}
}
func (l L) setv(i []int, u vector) {
	for k, j := range i {
		dx(l.v[j])
		l.v[j] = u.at(k)
	}
}

func (b B) MarshalJSON() ([]byte, error) { return json.Marshal(b.v) }
func (c C) MarshalJSON() ([]byte, error) { return json.Marshal(string(c.v)) }
func (i I) MarshalJSON() ([]byte, error) { return json.Marshal(i.v) }
func (f F) MarshalJSON() ([]byte, error) { return json.Marshal(f.v) }
func (s S) MarshalJSON() ([]byte, error) { return json.Marshal(s.v) }
func (l L) MarshalJSON() ([]byte, error) {
	r := []byte{'['}
	for i := range l.v {
		b, e := json.Marshal(l.v[i])
		if e != nil {
			return nil, e
		}
		r = append(r, b...)
		if i < len(l.v)-1 {
			r = append(r, ',')
		}
	}
	return append(r, ']'), nil
}

func dx(x T) T {
	if v, o := x.(refcounter); o {
		v.unref()
	}
	return x
}
func rx(x T) T {
	if v, o := x.(refcounter); o {
		v.ref()
	}
	return x
}

func KB(x []bool) B {
	r := B{v: x}
	r.init()
	return r
}
func KC(x []byte) C {
	r := C{v: x}
	r.init()
	return r
}
func KI(x []int) I {
	r := I{v: x}
	r.init()
	return r
}
func KF(x []float64) F {
	r := F{v: x}
	r.init()
	return r
}
func KZ(x []complex128) Z {
	r := Z{v: x}
	r.init()
	return r
}
func KS(x []string) S {
	r := S{v: x}
	r.init()
	return r
}
func KL(x []T) L {
	r := L{v: x}
	r.init()
	return r
}

func mk(x T, n int) T {
	switch v := x.(type) {
	case nil:
		return KL(make([]T, n))
	case bool:
		r := KB(make([]bool, n))
		if v {
			for i := range r.v {
				r.v[i] = v
			}
		}
		return r
	case byte:
		r := KC(make([]byte, n))
		if v != 0 {
			for i := range r.v {
				r.v[i] = v
			}
		}
		return r
	case int:
		r := KI(make([]int, n))
		if v != 0 {
			for i := range r.v {
				r.v[i] = v
			}
		}
		return r
	case float64:
		r := KF(make([]float64, n))
		if v != 0 {
			for i := range r.v {
				r.v[i] = v
			}
		}
		return r
	case complex128:
		r := KZ(make([]complex128, n))
		if v != 0 {
			for i := range r.v {
				r.v[i] = v
			}
		}
		return r
	case string:
		r := KS(make([]string, n))
		if v != "" {
			for i := range r.v {
				r.v[i] = v
			}
		}
		return r
	default:
		panic("type")
	}
}

//@ @1 /`i
//@ @1 2 /`I
func typ(x T) T {
	defer dx(x)
	switch v := x.(type) {
	case bool:
		return "b"
	case byte:
		return "c"
	case int:
		return "i"
	case float64:
		return "f"
	case complex128:
		return "z"
	case string:
		return "s"
	case L:
		return "L"
	case D:
		if v.tab {
			return "T"
		}
		return "D"
	case B:
		return "B"
	case C:
		return "C"
	case I:
		return "I"
	case F:
		return "F"
	case Z:
		return "Z"
	case S:
		return "S"
	default:
		return reflect.TypeOf(x).String()
	}
}

func use2(x, y vector) vector {
	if x.refs() == 1 {
		y.unref()
		return x
	} else if y.refs() == 1 {
		x.unref()
		return y
	}
	return use(x)
}
func use(x vector) vector {
	if x.refs() == 1 {
		return x
	}
	r := mk(x.zero(), x.ln()).(vector)
	switch v := r.(type) {
	case B:
		copy(v.v, x.(B).v)
	case C:
		copy(v.v, x.(C).v)
	case I:
		copy(v.v, x.(I).v)
	case F:
		copy(v.v, x.(F).v)
	case Z:
		copy(v.v, x.(Z).v)
	case S:
		copy(v.v, x.(S).v)
	case L:
		copy(v.v, x.(L).v) //shallow
	default:
		panic("use-type")
	}
	x.unref()
	return r
}

func (k *K) Vars() (name []string, typs []byte, n []int, rc []int) {
	name = make([]string, len(k.Var))
	typs = make([]byte, len(k.Var))
	n = make([]int, len(k.Var))
	rc = make([]int, len(k.Var))
	i := 0
	for s, v := range k.Var {
		name[i] = s
		if r, o := v.(refcounter); o {
			rc[i] = r.refs()
		}
		if vec, o := v.(vector); o {
			n[i] = vec.ln()
		}
		rx(v)
		t := typ(v).(string)
		typs[i] = byte(t[0])
		i++
	}
	return
}
