package jgo

import "math/bits"

type I = uint32
type SI = int32

var MI []I

func init() { MI = make([]I, 1<<16>>2); mt_init() }

type J struct{}

func (o J) J(x I) I { return j(x) }
func (o J) M() []I  { return MI }

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
func xxx(x I)     { panic("xxx") }
