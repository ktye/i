//go:build ignore

package main

import "encoding/binary"

// armv6 simulator (rpi0+1)

// doc: azeria-labs.com/writing-arm-assembly-part-1
// github.com/lseelenbinder/armsim/blob/master/armsim/instructions.go

func main() {
}

type cpu struct {
	r [32]uint32 // reg
	m []byte     // mem
}

func New(m int) *cpu {
	c := cpu{m: make([]byte, m)}
	return &c
}
func (c *cpu) load(p []byte, off uint32) {
	copy(c.m[off:], p)
	c.r[pc] = off
}
func (c *cpu) step() {
	u := binary.LittleEndian.Uint32(c.m[c.r[pc]:])
	c.r[pc] += 4
	dec(u)
}

// modes:
// svc(supervisor) hyp(hypervisor) irq(interrupt) fiq(fast-interrupt) abt(memory-abort) und(undefined-instruction)

const (
	r0 uint32 = iota
	r1
	r2
	r3
	r4
	r5
	r6
	r7
	r8
	r9
	r10
	fp
	ip
	sp   //stack pointer
	lr   //link register
	pc   //program counter
	cpsr //current program status
)

// cpsr bits
//
//	NZCVQ,GE,E: user writable
//	AIF,M: privilid. writable
//	JT: execution state priv writable
const (
	M0  = iota // M mode (usr,svc,..)
	M1         //
	M2         //
	M3         //
	T          // thumb 1(thumb)
	F          //
	I          //
	A          //
	E          // endian 0(little)
	R10        // 10-15 reserved
	R12        //
	R13        //
	R14        //
	R15        //
	GE0        // GE
	GE1        //
	GE2        //
	GE3        //
	R20        // reserved
	R21        //
	R22        //
	R23        //
	J          // jazelle
	R25        // reserved
	R26        //
	Q          //
	V          // overflow
	C          // carry
	Z          // zero
	N          // negative
)

const (
	EQ = iota // equal
	NE        // not equal
	CS        // carry set
	CC        // carry clear
	MI        // minus
	PL        // plus
	VS        // overflow
	VC        // no overflow
	HI        // unsigned higher
	LS        // unsigned lower or same
	GE        // signed greater than or eq
	LT        // signed less than
	GT        // signed greater than
	LE        // signed less than or eq
	AL        // always
)

func (c *cpu) dec(u uint32) {
	//              cccctttoooosnnnnddddaaaaass0mmmm
	//u := uint32(0b00000001111000000000000000000000)
	//fmt.Printf("%x %032b\n", u, u)
	cnd := (u & 0xf0000000) >> 28 // cond (4)
	typ := (u & 0x0e000000) >> 25 // type (3)
	opc := (u & 0x01e00000) >> 21 // opcode (4)
	rgn := (u & 0x000f0000) >> 16 // rn (4)
	rgd := (u & 0x0000f000) >> 12 // rd (4)
	sha := (u & 0x00000f80) >> 7  // shift amount (5)
	shs := (u & 0x00000060) >> 5  // shift (2)
	rgm := (u & 0x0000000f) >> 0  // rm (4)
	switch typ {
	case 0, 1:
		var y uint32
		if typ == 1 { // alisdair.mcdiarmid.org/arm-immediate-value-encoding/
			y = ror(u&0xff, 2*((u&0x1e00f00)>>8))
		} else {
			y = c.r[rgm]
		}
		c.r[rgd] = c.op(opc, 1&(u>>20) != 0, c.r[rgn], y)
	default:
		panic("todo")
	}
}
func ror(x, s, r uint32) uint32 { return (x >> r) | (x<<(32-r))&0xffffffff }
func (c *cpu) op(c uint32, s bool, x, y uint32) (r uint32) {
	carry := c.C()
	switch c {
	case 0:
		r = x & y // and
	case 1:
		r = x ^ y // xor
	case 2:
		r = x - y // sub
	case 3:
		r = y - x // rsb
	case 4:
		r = x + y // add
		u := uint64(x) + uint64(y)
		r = uint32(u)
		carry = u > 0xffffffff
	case 5: // adc
		u := uint64(x) + uint64(y) + uint64(carry)
		r = uint32(u)
		carry = u > 0xffffffff
	/*
		case 6:
			r = uint64(x) - uint64(y) // sbc
		case 7:
			r = y - x // rsc
			//todo carry
		case 8:
			if x != y {
				r = 1 // tst
			}
			//todo cond
		case 9:
			if x != y {
				r = 1 // teq
			}
			//todo cond
		case 10:
			x - y // dont update r
			// cmp
		case 11:
			x + y // dont update r
			// cmn
	*/
	case 12:
		r = x | y // orr
	case 13:
		r = y // mov
	case 14:
		r = x &^ y // bic
	case 15:
		r = ^y // mvn
	default:
		panic("op")
	}
	if s { // update carry if requested
		c := 1 << C
		if carry {
			c.r[cpsr] |= c
		} else {
			c.r[cpsr] &^= c
		}
	}
	return r
}
func (c *cpu) C() uint32 { return 1 & (c.r[cpsr] >> C) }
