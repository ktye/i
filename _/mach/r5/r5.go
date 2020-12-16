package r5

import (
	"encoding/binary"
	"fmt"
	"math"
)

var pc uint32
var m []byte
var x [32]uint32
var y [32]float64

func Start(b []byte, p uint32) {
	m = b
	pc = p
	for {
		f, a, b, c := dec(u(pc))
		f(a, b, c)
		pc += 4
	}
}

func dec(x uint32) (f, uint32, uint32, uint32) {
	switch x & 127 {
	case 3:
		return decI(fL, x)
	case 7:
		return fld, rd(x), rs1(x), immI(x)
	case 19:
		return decI(fI, x)
	case 35:
		return decS(fs, x)
	case 39:
		return fsd, immS(x), rs1(x), rs2(x)
	case 51:
		if f3(x) == 0 {
			return decR(fS, x)
		}
		return decR(fR, x)
	case 83:
		if 0x02000000&x == 0 && 0xf0000000&x != 0 {
			return decD(fD[3&x>>27], x)
		} else if f7(x) == 81 {
			return decD(fE[f3(x)], x)
		} else if f7(x) == 45 {
			return fsqr, rd(x), rs1(x), 0
		} else if f7(x) == 97 && rs2(x) == 0 {
			return fcvtwd, rd(x), rs1(x), 0
		} else if f7(x) == 105 && rs2(x) == 0 {
			return fcvtdw, rd(x), rs1(x), 0
		} else {
			return decD(fD, x)
		}
	case 99:
		return decB(fB, x)
	default:
		fpanic("decode %x\n", x)
	}
	panic(1)
}

func decI(t [8]f, x uint32) (f, uint32, uint32, uint32) { return t[f3(x)], rd(x), rs1(x), immI(x) }
func decR(t [8]f, x uint32) (f, uint32, uint32, uint32) { return t[f3(x)], rd(x), rs1(x), rs2(x) }
func decB(t [8]f, x uint32) (f, uint32, uint32, uint32) { return t[f3(x)], immB(x), rs1(x), rs2(x) }
func decS(t [3]f, x uint32) (f, uint32, uint32, uint32) { return t[f3(x)], immS(x), rs1(x), rs2(x) }
func decD(t f, x uint32) (f, uint32, uint32, uint32)    { return f, rd(x), rs1(x), rs2(x) }

func rd(x uint32) uint32   { return 31 & x >> 7 }
func rs1(x uint32) uint32  { return 31 & x >> 15 }
func rs2(x uint32) uint32  { return 31 & x >> 20 }
func f3(x uint32) uint32   { return 7 & x >> 12 }
func f7(x uint32) uint32   { return 127 & x >> 25 }
func immI(x uint32) uint32 { return uint32(int32(x) >> 20) }
func immB(x uint32) uint32 { return (30 & x >> 7) | (x&128)<<4 | uint32(int32(0xc0000000&x)>>18) }
func immS(x uint32) uint32 { return uint32(int32(0xfe000000&x)>>20) | (15 & x >> 7) }

var fL = [8]f{lb, xxx, lw, xxx, lbu, xxx, xxx, xxx}
var fI = [8]f{addi, slli, xxx, xxx, xori, srri, ori, andi}
var fR = [8]f{add, sll, slt, sltu, xor, srl, or, and}
var fS = [8]f{sub, xxx, xxx, xxx, xxx, sra, xxx, xxx}
var fB = [8]f{beq, bne, xxx, xxx, blt, bge, bltu, bgeu}
var fs = [3]f{sb, xxx, sw}
var fD = [4]f{fadd, fsub, fmul, fdiv}
var fE = [4]f{fle, flt, feq, xxx}

type f func(uint32, uint32, uint32)

func xxx(c, a, b uint32)  { panic("illegal instruction") }
func sb(c, a, b uint32)   { m[c+a] = byte(b) }
func sw(c, a, b uint32)   { su(c+a, b) }
func lb(c, a, b uint32)   { x[c] = uint32(int8(m[x[a]+b])) }
func lw(c, a, b uint32)   { x[c] = u(x[a] + b) }
func lbu(c, a, b uint32)  { x[c] = uint32(m[x[a]+b]) }
func addi(c, a, b uint32) { x[c] = x[a] + b }
func slli(c, a, b uint32) { x[c] = x[a] << (b & 31) }
func srri(c, a, b uint32) { q := bl(b > 32); x[c] = q*(uint32(int32(x[a])>>(b&31))) + (1-q)*x[a]>>b }
func xori(c, a, b uint32) { x[c] = x[a] ^ b }
func ori(c, a, b uint32)  { x[c] = x[a] | b }
func andi(c, a, b uint32) { x[c] = x[a] & b }

func add(c, a, b uint32)  { x[c] = x[a] + x[b] }
func sub(c, a, b uint32)  { x[c] = x[a] - x[b] }
func sll(c, a, b uint32)  { x[c] = x[a] << x[b] }
func slt(c, a, b uint32)  { x[c] = bl(int32(x[a]) < int32(x[b])) }
func sltu(c, a, b uint32) { x[c] = bl(x[a] < x[b]) }
func xor(c, a, b uint32)  { x[c] = x[a] ^ x[b] }
func srl(c, a, b uint32)  { x[c] = x[a] >> x[b] }
func sra(c, a, b uint32)  { x[c] = uint32(int32(x[a]) >> x[b]) }
func or(c, a, b uint32)   { x[c] = x[a] | x[b] }
func and(c, a, b uint32)  { x[c] = x[a] & x[b] }

func beq(c, a, b uint32)  { pc += bl(a == b) * c }
func bne(c, a, b uint32)  { pc += bl(a != b) * c }
func bge(c, a, b uint32)  { pc += bl(int32(a) >= int32(b)) * c }
func bgeu(c, a, b uint32) { pc += bl(a >= b) * c }
func blt(c, a, b uint32)  { pc += bl(int32(a) < int32(b)) * c }
func bltu(c, a, b uint32) { pc += bl(a < b) * c }

func fld(c, a, b uint32)    { y[c] = d(x[a] + b) }
func fsd(c, a, b uint32)    { sd(x[a]+c, y[b]) }
func fadd(c, a, b uint32)   { y[c] = y[a] + y[b] }
func fsub(c, a, b uint32)   { y[c] = y[a] - y[b] }
func fmul(c, a, b uint32)   { y[c] = y[a] * y[b] }
func fdiv(c, a, b uint32)   { y[c] = y[a] / y[b] }
func fsqr(c, a, b uint32)   { y[c] = math.Sqrt(y[a]) }
func fle(c, a, b uint32)    { x[c] = bl(y[a] <= y[b]) }
func flt(c, a, b uint32)    { x[c] = bl(y[a] < y[b]) }
func feq(c, a, b uint32)    { x[c] = bl(y[a] == y[b]) }
func fcvtwd(c, a, b uint32) { x[c] = uint32(int32(y[a])) }
func fcvtdw(c, a, b uint32) { y[c] = float64(int32(x[a])) }

func u(i uint32) uint32      { return binary.LittleEndian.Uint32(m[i:]) }
func d(i uint32) float64     { return math.Float64frombits(binary.LittleEndian.Uint64(m[i:])) }
func su(i, x uint32)         { binary.LittleEndian.PutUint32(m[i:], x) }
func sd(i uint32, x float64) { binary.LittleEndian.PutUint64(m[i:], math.Float64bits(x)) }
func bl(x bool) uint32 {
	if x {
		return 1
	}
	return 0
}
func fpanic(s string, a ...interface{}) { panic(fmt.Sprintf(s, a)) }
