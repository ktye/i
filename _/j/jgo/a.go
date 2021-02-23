package jgo

import (
	"fmt"
	"j/x"
	"math/bits"
)

type I = uint32
type SI = int32

var MI []I

func init() { MI = make([]I, 1<<16>>2); mt_init() }

func New() interface {
	J(x I) I
	M() []I
} {
	return jJ{}
}

type jJ struct{}

func (o jJ) J(x I) I { return j(x) }
func (o jJ) M() []I  { return MI }

func n32(x I) I {
	if x == 0 {
		return 1
	} else {
		return 0
	}
}
func i32b(x bool) I {
	if x {
		return 1
	} else {
		return 0
	}
}
func clz32(x I) I { return I(bits.LeadingZeros32(x)) }
func stk(y I)     { fmt.Println(x.X(MI, MI[1])) }
func xxx(x I)     { panic("xxx") }
