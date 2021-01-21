package main

import "testing"

var tinit = "t:`a`b`c!(1+!10;2+!10;3+!10)"
var ttests = [][2]string{
	{"t{5<c}", "`a`b`c!(4 5 6 7 8 9 10;5 6 7 8 9 10 11;6 7 8 9 10 11 12)"},
}

func TestT(t *testing.T) {
	kinit()
	dx(eval(tinit))
	for _, tc := range ttests {
		s := kstval(tc[0])
		if s != tc[1] {
			t.Fatalf("%s\ngot:%s\nexp:%s\n", tc[0], s, tc[1])
		}
	}
	bleak()
}
