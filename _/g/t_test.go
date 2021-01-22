package main

import (
	"fmt"
	"testing"
)

var tinit = "t:`a`b`c!(1+!10;2+!10;`q`w`e`q`w`e`q`w`e`q)"
var ttests = [][2]string{
	{"t{~c=`q}", "`a`b`c!(2 3 5 6 8 9;3 4 6 7 9 10;`w`e`w`e`w`e)"},
	{"+/'`a`b#t", "`a`b!(55;65)"},
	{"(*;+/)'`a`b#t", "``a`b!(`*`+{;(1;55);(2;65))"},
	{"(`first`sum!(*;+/))'`a`b#t", "``a`b!(`first`sum;(1;55);(2;65))"},
	{"`c=t", "(`q`w`e;`a`b;((1 4 7 10;2 5 8;3 6 9);(2 5 8 11;3 6 9;4 7 10)))"},
}

func TestT(t *testing.T) {
	kinit()
	dx(eval(tinit))
	for _, tc := range ttests {
		s := kstval(tc[0])
		fmt.Println(s)
		if s != tc[1] {
			t.Fatalf("%s\ngot:%s\nexp:%s\n", tc[0], s, tc[1])
		}
	}
	bleak()
}
