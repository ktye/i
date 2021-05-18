package k

import (
	"bytes"
	"fmt"
	"math"
	"math/cmplx"
	"regexp"
	"strconv"
	"strings"
)

type stringer struct {
	s string
}

func (s stringer) String() string { return s.s }

//$ $1.23 /"1.23"
//$ $1 2 3 /(,"1";,"2";,"3")
func (k *K) str(x T) T {
	if _, o := x.(vector); o {
		return each(f1(k.str), x)
	} else {
		return KC([]byte(sstr(x)))
	}
}

func debug(x T) T     { fmt.Println(sstr(x)); return x }
func debug2(x, y T) T { fmt.Println(sstr(x)+":", sstr(y)); dx(x); return y }

//$ `i$"3" /3
//$ `f$("3";"4.5") /3 4.5
//$ `z$"3.1" /3.1a0
//$ `z$"3a20" /3a20
//$ `$"alpha" /`alpha
//$ `$("a";"bc") /`a`bc
//$ 0$" ab c  " /"ab c"  (trim)
//$ 3$("a";"beta") /("a  ";"bet")  (pad)
func (k *K) cast(x, y T) T {
	if _, o := y.(L); o {
		return eachright(f2(k.cast), x, y)
	}
	if c, o := y.(byte); o {
		y = KC([]byte{c})
	}
	if s, o := x.(string); o {
		if c, o := y.(C); o {
			defer dx(y)
			ys := string(c.v)
			switch s {
			case "":
				return ys
			case "i":
				i, e := strconv.Atoi(ys)
				if e != nil {
					panic("parse")
				}
				return i
			case "f":
				switch ys {
				case "NaN", "0n":
					return math.NaN()
				case "0w", "+Inf":
					return math.Inf(1)
				case "-0w", "-Inf":
					return math.Inf(-1)
				}
				f, e := strconv.ParseFloat(ys, 64)
				if e != nil {
					panic("parse")
				}
				return f
			case "z":
				i := strings.Index(ys, "a")
				if i < 0 {
					f, e := strconv.ParseFloat(ys, 64)
					if e != nil {
						panic("parse")
					}
					return complex(f, 0)
				}
				f, e1 := strconv.ParseFloat(ys[:i], 64)
				g, e2 := strconv.ParseFloat(ys[1+i:], 64)
				if e1 == nil && e2 == nil {
					return rrot(f, g)
				}
				panic("parse")
			}
		}
	}
	if n, o := x.(int); o {
		if c, o := y.(C); o {
			defer dx(y)
			if n == 0 {
				c = use(c).(C)
				c.v = bytes.TrimSpace(c.v)
				return c
			} else if n > 0 {
				r := mk(byte(32), n).(C)
				copy(r.v, c.v)
				return r
			} else if n < 0 {
				n = -n
				r := mk(byte(32), n).(C)
				n = maxi(0, len(r.v)-len(c.v))
				copy(r.v, c.v[n:])
				return r
			}
		}
	}
	panic("type")
}

func kst(x T) T { defer dx(x); return KC([]byte(sstr(x))) }
func sstr(x T) string {
	switch v := x.(type) {
	case nil:
		return ""
	case stringer:
		return v.s
	case bool:
		if v {
			return "1b"
		}
		return "0b"
	case byte:
		return `"` + string(v) + `"`
	case string:
		return "`" + v
	case float64:
		return fdot(ftoa(v))
	case complex128:
		return absang(v)
	case projection:
		if v.s != "" {
			return v.s
		} else {
			return sstr(v.f) + v.x.brackets()
		}
	case *regexp.Regexp:
		return "(/" + strconv.Quote(v.String()) + ")"
	case interface{ String() string }:
		return v.String()
	default:
		return fmt.Sprint(x)
	}
}

func (b B) String() string {
	if len(b.v) == 0 {
		return "0#0b"
	}
	c := make([]byte, 1+len(b.v))
	for i, u := range b.v {
		if u {
			c[i] = '1'
		} else {
			c[i] = '0'
		}
	}
	c[len(c)-1] = 'b'
	if len(b.v) == 1 {
		return "," + string(c)
	}
	return string(c)
}
func (c C) String() string { return comma(len(c.v) == 1, strconv.Quote(string(c.v))) }
func (i I) String() string {
	if len(i.v) == 0 {
		return "!0"
	}
	return comma(len(i.v) == 1, unbrack(fmt.Sprint(i.v)))
}
func (f F) String() string {
	if len(f.v) == 0 {
		return "0#0."
	}
	v := make([]string, len(f.v))
	for i, u := range f.v {
		v[i] = ftoa(u)
	}
	s := comma(len(f.v) == 1, strings.Join(v, " "))
	if strings.Index(s, ".") < 0 && strings.Index(s, "e") < 0 && strings.Index(s, "n") < 0 && strings.Index(s, "w") < 0 {
		return s + "."
	}
	return s
}
func (z Z) String() string {
	if len(z.v) == 0 {
		return "0#0a"
	}
	return comma(len(z.v) == 1, absangs(z.v))
}
func (l L) String() string {
	v := make([]string, len(l.v))
	for i, s := range l.v {
		v[i] = sstr(s)
	}
	if len(l.v) == 1 {
		return "," + v[0]
	}
	return "(" + strings.Join(v, ";") + ")"
}
func (l L) brackets() string {
	s := l.String()
	return "[" + s[1:len(s)-1] + "]"
}
func (s S) String() string {
	if len(s.v) == 0 {
		return "0#`"
	}
	return "`" + strings.Join(s.v, "`") // todo `"a b"
}
func (d D) String() string {
	s := sstr(d.k)
	if s[0] == ',' {
		s = "(" + s + ")"
	}
	return s + "!" + sstr(d.v)
}
func (d derived) String() string { return d.s }
func (t train) String() (s string) {
	for _, f := range t {
		s += f.String()
	}
	return s
}
func comma(x bool, s string) string {
	if x {
		return "," + s
	}
	return s
}
func unbrack(s string) string { return s[1 : len(s)-1] }
func absangs(v []complex128) string {
	s := make([]string, len(v))
	for i, u := range v {
		s[i] = absang(u)
	}
	return strings.Join(s, " ")
}
func absang(z complex128) string {
	if cmplx.IsNaN(z) {
		return "0na"
	}
	r, phi := cmplx.Polar(z)
	phi *= 180 / math.Pi
	if phi < 0 {
		phi += 360
	}
	if r == 0 {
		phi = 0.0
	}
	if phi == -0.0 || phi == 360.0 {
		phi = 0.0
	}
	return fmt.Sprintf("%va%v", r, phi)
}
func ftoa(x float64) string {
	if math.IsNaN(x) {
		return "0n"
	} else if math.IsInf(x, 1) {
		return "0w"
	} else if math.IsInf(x, -1) {
		return "-0w"
	}
	return strconv.FormatFloat(x, 'g', -1, 64)
}
func fdot(s string) string {
	if strings.IndexAny(s, ".enw") < 0 {
		return s + "."
	}
	return s
}
