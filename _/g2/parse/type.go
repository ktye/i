package main

import (
	"strings"
	"time"
)

type A interface{}
type Byte []byte                         // "a" "abc" 0x12ef         ""      " "   `b 32
type Symbol []string                     // `a`b `"alpha"            0#`     `
type Number []float64                    // 1 1.0 1.1 -1.23e-002     !0      0n
type Complex []complex128                // 1.2a310 2a               0#0a    0na
type Time []time.Time                    // 2021.03.11T12:15:17      0#0T    0T
type Duration []time.Duration            // 2.5h 2ms 3us             0#0s    0s
type List []A                            // (1;2 3;4.5)              ()      ""
type Dict [2]A                           // `a`b!(1;2 3)
type Table [2]A                          // +`a`b!(1 2;3 4)
type Verb string                         // + / ':
type Func struct{}                       // {x+y}
type Derived struct{ Verb, Adverb Verb } // +/
type Projection []A                      // +[;3] 2-
type Composition []A                     // +-
type F1 func(x A) A                      // native functions
type F2 func(x, y A) A
type F3 func(x, y, z A) A

type SrcError struct {
	Src []byte
	Pos int
	Err error
}

func (s SrcError) Error() string {
	o := 0
	for i := 0; i < s.Pos; i++ {
		if s.Src[i] == 10 {
			o = i + 1
		}
	}
	s.Pos -= o
	s.Src = s.Src[o:]
	for i, c := range s.Src {
		if c == 10 && i > s.Pos {
			s.Src = s.Src[:i]
			break
		}
	}

	if s.Pos > 30 {
		n := s.Pos - 20
		s.Pos -= n
		s.Src = s.Src[20:]
	}
	if len(s.Src) > 70 {
		s.Src = append(s.Src[:70], '.', '.')
	}
	space := ""
	if s.Pos > 0 {
		space = strings.Repeat(" ", s.Pos)
	}
	return string(s.Src) + "\n" + space + "^" + s.Err.Error()
}
