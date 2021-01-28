package main

import (
	"fmt"
)

type Number interface {
	type int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, complex128
}

type Real interface {
	type int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64
}

func GradeUp[T Real](x []T) (r []int) {
	r = Til(len(x))
	sort.SliceStable(r, func(i, j int) bool { return x[r[i]] < x[r[j]] })
}

func GradeDown[T Real](x []T) (r []int) {
	r = Til(len(x))
	sort.SliceStable(r, func(i, j int) bool { return x[r[i]] > x[r[j]] })
}

func Index[T any](x []T, i []int) (r []T) {
	r = make([]T, len(i))
	for i, k := range i {
		r[i] = x[k]
	}
	return r
}

func Not[T Number](x []T) []bool { // or uint8?
	r := make([]bool, len(x))
	for i, u := range x {
		r[i] = x != 0
	}
	return r
}

func Add[T Number](x, y []T) []T {
	for i, u := range y { 
		x[i] += u
	}
	return x
}
func Add1[T Number](x []T, y T) []T {
	for i := range x {
		x[i] += y
	}
	return x
}
func Sub[T Number](x, y []T) []T {
	for i, u := range y {
		x[i] -= u
	}
	return x
}
func Sub1[T Number](x []T, y T) []T {
	for i := range x {
		x[i] -= y
	}
	return x
}
func Sub2[T Number](x T, y []T) []T {
	for i, u := range y {
		y[i] = x - u
	}
	return y
}

func Til(n int) (r []int) {
	r = make([]int, n)
	for i := range r {
		r[i] = i
	}
	return r
}

func Take[T any](n int, y []T) (r []T) {
	if n < 0 {
		return Rev(Take(-n, Rev(y)))
	}
	r = make([]T, n)
	for i := range r {
		r[i] = y[i%len(y)]
	}
	return r
}

func Rev[T any](x []T) []T {
	if len(x) < 2 {
		return x
	}
	k := len(x) - 1
	for i := 0; i < len(x)/2; i++ {
		x[i], x[k] = x[k], x[i]
		k--
	}
	return x
}

func main() {
	o := fmt.Println
	o(Til(5))
	o(Take(2, Til(5)))
	o(Take(7, Til(5)))
	o(Take(-2, Til(5)))
	o(Add([]int{1, 2, 3}, []int{4, 5, 6}))
	o(Add1([]int{1, 2, 3}, 4))
	o(Sub1([]float64{1, 2, 3}, 4))
}

