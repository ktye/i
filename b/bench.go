package main

import (
	. "github.com/ktye/wg/module"
	"github.com/ktye/wg/wasi_unstable"
)

func init() {
	Memory(32768) // 2g
	Functions(0, addi, addf, subi, subf)
}

func main() {
	for i := int32(0); i < 268435456; i++ {
		SetI32(i, i)
	}
	writeln('!')
	isub(5)
	iadd(4)
	isub(4)
	vaddi(4)
	ndi(0, 4)
	ndi(2, 4)
	fadd(3.1)
	fsub(3.2)
	vaddf(3.3)
	ndf(1, 3.1)
	ndf(3, 3.2)
}
func vaddi(x int32) {
	write('V')
	write('+')
	a := clock()
	yp := int32(0)
	xv := I32x4splat(x)
	for i := int32(0); i < 67108864; i++ {
		I32x4store(yp, I32x4add(xv, I32x4load(yp)))
		yp += 16
	}
	bench(a)
}
func vaddf(x float64) {
	write('W')
	write('+')
	a := clock()
	yp := int32(0)
	xv := F64x2splat(x)
	for i := int32(0); i < 134217728; i++ {
		F64x2store(yp, F64x2add(xv, F64x2load(yp)))
		yp += 16
	}
	bench(a)
}
func iadd(x int32) {
	write('I')
	write('+')
	a := clock()
	yp := int32(0)
	for i := int32(0); i < 268435456; i++ {
		SetI32(yp, x+I32(yp))
		yp += 4
	}
	bench(a)
}
func isub(x int32) {
	write('I')
	write('-')
	a := clock()
	yp := int32(0)
	for i := int32(0); i < 268435456; i++ {
		SetI32(yp, x-I32(yp))
		yp += 4
	}
	bench(a)
}
func fadd(x float64) {
	write('F')
	write('-')
	a := clock()
	yp := int32(0)
	for i := int32(0); i < 268435456; i++ {
		SetF64(yp, x+F64(yp))
		yp += 8
	}
	bench(a)
}
func fsub(x float64) {
	write('F')
	write('-')
	a := clock()
	yp := int32(0)
	for i := int32(0); i < 268435456; i++ {
		SetF64(yp, x-F64(yp))
		yp += 8
	}
	bench(a)
}
func ndi(f int32, x int32) {
	write('i')
	write('0' + f)
	a := clock()
	yp := int32(0)
	for i := int32(0); i < 268435456; i++ {
		SetI32(yp, Func[f].(f2i)(x, I32(yp)))
		yp += 4
	}
	bench(a)
}
func ndf(f int32, x float64) {
	write('f')
	write('0' + f)
	a := clock()
	yp := int32(0)
	for i := int32(0); i < 268435456; i++ {
		SetF64(yp, Func[f].(f2f)(x, F64(yp)))
		yp += 8
	}
	bench(a)
}

type f2i = func(int32, int32) int32
type f2f = func(float64, float64) float64

func addi(x, y int32) int32     { return x + y }
func subi(x, y int32) int32     { return x - y }
func addf(x, y float64) float64 { return x + y }
func subf(x, y float64) float64 { return x - y }

func bench(a int64) { writei(clock() - a) }
func write(x int32) {
	SetI32(0, int32(8))
	SetI32(4, int32(1))
	SetI32(8, x)
	wasi_unstable.Fd_write(1, 0, 1, 0)
}
func writeln(x int32) {
	write(x)
	write('\n')
}
func writei(i int64) { // reverse
	write('|')
	if i == 0 {
		write('0')
	}
	for i != 0 {
		write('0' + int32(i%10))
		i /= 10
	}
	write('\n')
}
func clock() int64 {
	wasi_unstable.Clock_time_get(1, wasi_unstable.Timestamp(0), 0)
	return I64(0)
}
