package main

import (
	"fmt"
)

type A interface {
	Ref()
	Unref()
	Zero() A
	Empty() A
}

type numeric interface {
	type int, float64, complex128
}

func NewFixed[T numeric]() A {
	var a Fixed[T]
	a.a = make([]T, 0)
	return &a
}

type Fixed[T numeric] struct {
	a  []T
	rc int
}
type Atom[T numeric] T

func (a Atom[T]) Ref()   {}
func (a Atom[T]) Unref() {}
func (a Atom[T]) Zero() A {
	a = 0
	return a
}
func (A Atom[T]) Empty() A {
	return NewFixed[T]()
}

func (a *Fixed[T]) Ref()   { a.rc++ }
func (a *Fixed[T]) Unref() { a.rc-- }
func (a *Fixed[T]) Zero() A {
	var z Atom[T]
	return z
}
func (a *Fixed[T]) Empty() A {
	return NewFixed[T]()
}

func main() {
	var a Fixed[int]
	fmt.Println(a)
}

