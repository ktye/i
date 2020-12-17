package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

type R5 struct {
	pc uint32
	m  []byte
	x  [32]uint32
	y  [32]float64
	f  map[string]func(uint32, uint32, uint32)
}

func (r *R5) Start(b []byte, p uint32) {
	r.init(b, p)
	for {
		u := r.u(r.pc)
		if u == 0 {
			return
		}
		s, a, b, c := r.dec(u)
		f := r.f[s]
		f(a, b, c)
		r.pc += 4
	}
}
func (r *R5) Dump(b []byte, p uint32) {
	r.init(b, p)
	for {
		u := r.u(r.pc)
		if u == 0 {
			return
		}
		fmt.Printf("%8x %8x ", r.pc, u)
		s, a, b, c := r.dec(u)
		fmt.Printf("%5s,%d,%d,%d\n", s, int32(a), int32(b), int32(c))
		if r.f[s] == nil {
			panic(s)
		}
		r.pc += 4
	}
}
func (r *R5) init(b []byte, p uint32) {
	r.pc = p
	r.m = b
	r.f = map[string]func(uint32, uint32, uint32){
		"xxx": r.xxx, "sb": r.sb, "sw": r.sw, "lb": r.lb, "lw": r.lw, "lbu": r.lbu,
		"addi": r.addi, "slli": r.slli, "srli": r.srli, "xori": r.xori, "ori": r.ori,
		"andi": r.andi, "add": r.add, "sub": r.sub, "sll": r.sll,
		"slt": r.slt, "sltu": r.sltu, "xor": r.xor, "srl": r.srl,
		"sra": r.sra, "or": r.or, "and": r.and, "mul": r.mul, "div": r.div,
		"divu": r.divu, "rem": r.rem, "remu": r.remu, "beq": r.beq,
		"bne": r.bne, "bge": r.bge, "bgeu": r.bgeu, "blt": r.blt,
		"bltu": r.bltu, "fld": r.fld, "fsd": r.fsd, "fadd": r.fadd,
		"fsub": r.fsub, "fmul": r.fmul, "fdiv": r.fdiv, "fsqr": r.fsqr,
		"fle": r.fle, "flt": r.flt, "feq": r.feq,
		"fcvtwd": r.fcvtwd, "fcvtdw": r.fcvtdw,
	}
}

func (r *R5) dec(x uint32) (string, uint32, uint32, uint32) {
	switch x & 127 {
	case 3:
		return rvI(fL, x)
	case 7:
		return "fld", rd(x), rs1(x), immI(x)
	case 19:
		return rvI(fI, x)
	case 35:
		return rvS(fs, x)
	case 39:
		return "fsd", immS(x), rs1(x), rs2(x)
	case 51:
		if f7(x) == 0 {
			return rvR(fR, x)
		} else if f7(x) == 1 {
			return rvR(fM, x)
		}
		return rvR(fS, x)
	case 83:
		if 0x02000000&x == 0 && 0xf0000000&x != 0 {
			return rvD(fD[3&x>>27], x)
		} else if f7(x) == 81 {
			return rvD(fE[f3(x)], x)
		} else if f7(x) == 45 {
			return "fsqr", rd(x), rs1(x), 0
		} else if f7(x) == 97 && rs2(x) == 0 {
			return "fcvtwd", rd(x), rs1(x), 0
		} else if f7(x) == 105 && rs2(x) == 0 {
			return "fcvtdw", rd(x), rs1(x), 0
		}
	case 99:
		return rvB(fB, x)
	}
	panic(x)
}

func rvI(t [8]string, x uint32) (string, uint32, uint32, uint32) {
	return t[f3(x)], rd(x), rs1(x), immI(x)
}
func rvR(t [8]string, x uint32) (string, uint32, uint32, uint32) {
	return t[f3(x)], rd(x), rs1(x), rs2(x)
}
func rvB(t [8]string, x uint32) (string, uint32, uint32, uint32) {
	return t[f3(x)], rs1(x), rs2(x), immB(x)
}
func rvS(t [3]string, x uint32) (string, uint32, uint32, uint32) {
	return t[f3(x)], rs1(x), rs2(x), immS(x)
}
func rvD(t string, x uint32) (string, uint32, uint32, uint32) { return t, rd(x), rs1(x), rs2(x) }

func rd(x uint32) uint32   { return 31 & (x >> 7) }
func rs1(x uint32) uint32  { return 31 & (x >> 15) }
func rs2(x uint32) uint32  { return 31 & (x >> 20) }
func f3(x uint32) uint32   { return 7 & (x >> 12) }
func f7(x uint32) uint32   { return 127 & (x >> 25) }
func immI(x uint32) uint32 { return uint32(int32(x) >> 20) }

/*
func B(x int32) uint32 {
	if x < 0 {
		x = 4096 + x
	}
	x = 128*(x%32) + 128*((x/2048)%2) + 33554432*(x/32) + 16777216*(x/4096)
	return uint32(x)
}
*/
func immB(x uint32) uint32 {
	return (30 & (x >> 7)) | ((x & 128) << 4) | ((0x7e000000 & x) >> 20) | uint32(int32(0x80000000&x)>>19)
}
func immS(x uint32) uint32 { return uint32(int32(0xfe000000&x)>>20) | (31 & (x >> 7)) }

var fL = [8]string{"lb", "xxx", "lw", "xxx", "lbu", "xxx", "xxx", "xxx"}
var fI = [8]string{"addi", "slli", "xxx", "xxx", "xori", "srli", "ori", "andi"}
var fR = [8]string{"add", "sll", "slt", "sltu", "xor", "srl", "or", "and"}
var fM = [8]string{"mul", "xxx", "xxx", "xxx", "div", "divu", "rem", "remu"}
var fS = [8]string{"sub", "xxx", "xxx", "xxx", "xxx", "sra", "xxx", "xxx"}
var fB = [8]string{"beq", "bne", "xxx", "xxx", "blt", "bge", "bltu", "bgeu"}
var fs = [3]string{"sb", "xxx", "sw"}
var fD = [4]string{"fadd", "fsub", "fmul", "fdiv"}
var fE = [4]string{"fle", "flt", "feq", "xxx"}

func (r *R5) xxx(c, a, b uint32)  { panic("illegal instruction") }
func (r *R5) sb(c, a, b uint32)   { r.m[a+b] = byte(c) }
func (r *R5) sw(c, a, b uint32)   { r.su(a+b, c) }
func (r *R5) lb(c, a, b uint32)   { r.x[c] = uint32(int8(r.m[r.x[a]+b])) }
func (r *R5) lw(c, a, b uint32)   { r.x[c] = r.u(r.x[a] + b) }
func (r *R5) lbu(c, a, b uint32)  { r.x[c] = uint32(r.m[r.x[a]+b]) }
func (r *R5) addi(c, a, b uint32) { r.x[c] = r.x[a] + b }
func (r *R5) slli(c, a, b uint32) { r.x[c] = r.x[a] << (b & 31) }
func (r *R5) srli(c, a, b uint32) { r.x[c] = r.x[a] >> b }
func (r *R5) xori(c, a, b uint32) { r.x[c] = r.x[a] ^ b }
func (r *R5) ori(c, a, b uint32)  { r.x[c] = r.x[a] | b }
func (r *R5) andi(c, a, b uint32) { r.x[c] = r.x[a] & b }

func (r *R5) add(c, a, b uint32)  { r.x[c] = r.x[a] + r.x[b] }
func (r *R5) sub(c, a, b uint32)  { r.x[c] = r.x[a] - r.x[b] }
func (r *R5) sll(c, a, b uint32)  { r.x[c] = r.x[a] << r.x[b] }
func (r *R5) slt(c, a, b uint32)  { r.x[c] = bl(int32(r.x[a]) < int32(r.x[b])) }
func (r *R5) sltu(c, a, b uint32) { r.x[c] = bl(r.x[a] < r.x[b]) }
func (r *R5) xor(c, a, b uint32)  { r.x[c] = r.x[a] ^ r.x[b] }
func (r *R5) srl(c, a, b uint32)  { r.x[c] = r.x[a] >> r.x[b] }
func (r *R5) sra(c, a, b uint32)  { r.x[c] = uint32(int32(r.x[a]) >> r.x[b]) }
func (r *R5) or(c, a, b uint32)   { r.x[c] = r.x[a] | r.x[b] }
func (r *R5) and(c, a, b uint32)  { r.x[c] = r.x[a] & r.x[b] }
func (r *R5) mul(c, a, b uint32)  { r.x[c] = r.x[a] * r.x[b] }
func (r *R5) div(c, a, b uint32)  { r.x[c] = uint32(int32(r.x[a]) / int32(r.x[b])) }
func (r *R5) divu(c, a, b uint32) { r.x[c] = r.x[a] / r.x[b] }
func (r *R5) rem(c, a, b uint32)  { r.x[c] = uint32(int32(r.x[a]) % int32(r.x[b])) }
func (r *R5) remu(c, a, b uint32) { r.x[c] = r.x[a] % r.x[b] }

func (r *R5) beq(c, a, b uint32)  { r.pc += bl(a == b) * c }
func (r *R5) bne(c, a, b uint32)  { r.pc += bl(a != b) * c }
func (r *R5) bge(c, a, b uint32)  { r.pc += bl(int32(a) >= int32(b)) * c }
func (r *R5) bgeu(c, a, b uint32) { r.pc += bl(a >= b) * c }
func (r *R5) blt(c, a, b uint32)  { r.pc += bl(int32(a) < int32(b)) * c }
func (r *R5) bltu(c, a, b uint32) { r.pc += bl(a < b) * c }

func (r *R5) fld(c, a, b uint32)    { r.y[c] = r.d(r.x[a] + b) }
func (r *R5) fsd(c, a, b uint32)    { r.sd(r.x[a]+c, r.y[b]) }
func (r *R5) fadd(c, a, b uint32)   { r.y[c] = r.y[a] + r.y[b] }
func (r *R5) fsub(c, a, b uint32)   { r.y[c] = r.y[a] - r.y[b] }
func (r *R5) fmul(c, a, b uint32)   { r.y[c] = r.y[a] * r.y[b] }
func (r *R5) fdiv(c, a, b uint32)   { r.y[c] = r.y[a] / r.y[b] }
func (r *R5) fsqr(c, a, b uint32)   { r.y[c] = math.Sqrt(r.y[a]) }
func (r *R5) fle(c, a, b uint32)    { r.x[c] = bl(r.y[a] <= r.y[b]) }
func (r *R5) flt(c, a, b uint32)    { r.x[c] = bl(r.y[a] < r.y[b]) }
func (r *R5) feq(c, a, b uint32)    { r.x[c] = bl(r.y[a] == r.y[b]) }
func (r *R5) fcvtwd(c, a, b uint32) { r.x[c] = uint32(int32(r.y[a])) }
func (r *R5) fcvtdw(c, a, b uint32) { r.y[c] = float64(int32(r.x[a])) }

func (r *R5) u(i uint32) uint32      { return binary.LittleEndian.Uint32(r.m[i:]) }
func (r *R5) d(i uint32) float64     { return math.Float64frombits(binary.LittleEndian.Uint64(r.m[i:])) }
func (r *R5) su(i, x uint32)         { binary.LittleEndian.PutUint32(r.m[i:], x) }
func (r *R5) sd(i uint32, x float64) { binary.LittleEndian.PutUint64(r.m[i:], math.Float64bits(x)) }
func bl(x bool) uint32 {
	if x {
		return 1
	}
	return 0
}
