package main

import (
	"fmt"
	"testing"
)

// B-type immediates are even values in the range -4096..4094
func TestR5ImmB(t *testing.T) {
	for i := int32(-4096); i < 4095; i += 2 {
		r := int32(immB(B(i)))
		if r != i {
			O(uint32(i))
			O(uint32(B(i)))
			O(uint32(r))
			t.Fatal()
		}
	}
}

// S-type immediates are values in the range -2048..2047
func TestR5ImmS(t *testing.T) {
	for i := int32(-2048); i < 2048; i++ {
		r := int32(immS(S(i)))
		if r != i {
			O(uint32(i))
			O(uint32(S(i)))
			O(uint32(r))
			t.Fatal()
		}
	}
}

// I-type immediates are values in the range -2048..2047
func TestR5ImmI(t *testing.T) {
	for i := int32(-2048); i < 2048; i++ {
		r := int32(immI(I(i)))
		if r != i {
			O(uint32(i))
			O(uint32(I(i)))
			O(uint32(r))
			t.Fatal()
		}
	}
}

// J-type immediates are even values in the range −1048576..1048575
func TestR5ImmJ(t *testing.T) {
	for i := int32(-1048576); i < 1048575; i += 2 {
		r := int32(immJ(J(i)))
		if r != i {
			O(uint32(i))
			O(uint32(J(i)))
			O(uint32(r))
			t.Fatal()
		}
	}
}

func O(u uint32) { fmt.Printf("%08x %032b %d\n", u, u, int32(u)) }

// B:{[op;f3;rs1;rs2;imm]o:-2147483648*0>imm;imm:n12 imm;o+op+/128 4096 32768 1048576 33554432*((32/imm)+2/2048\imm;f3;rs1;rs2;(4096\imm)+64/32\imm)}
func B(x int32) uint32 {
	var r int32
	var o int32
	if x < 0 {
		x = 4096 + x
		o = -2147483648
	}
	r += 128 * ((x % 32) + ((x / 2048) % 2))
	r += 33554432 * (((x / 32) % 64) + x/4096)
	r += o
	return uint32(r)
}

// S:{[op;f3;rs1;rs2;imm]imm:n12 imm;op+/128 4096 32768 1048576 33554432*(32/imm;f3;rs1;rs2;32\imm)}
func S(x int32) uint32 {
	if x < 0 {
		x += 4096
	}
	return uint32(128*(x%32) + 33554432*(x/32))
}

// I:{[op;f3;rd;rs1;imm]op+/128 4096 32768 1048576*(rd;f3;rs1;imm)}
func I(x int32) uint32 {
	return uint32(1048576 * x)
}

// J:{[op;rd;r1;imm]o:-2147483648*0>imm;imm:n21 imm;o+op+/128 4096 524288*(rd;256/4096\imm;4096/imm)}
func J(x int32) uint32 { // -1048576..1048575
	var r int32
	var o int32
	if x < 0 {
		x = 2097152 + x
		o = -2147483648
	}
	r += 524288 * (x % 4096)       // (lower 12 bits) << 19
	r += 4096 * ((x / 4096) % 256) // 8 center bits remain
	r += o                         // sign bit
	return uint32(r)
}
