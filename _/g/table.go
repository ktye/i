package main

// sql
// t:`a`b`c!(1+!10;2+!10;3+!10)  /dict-table
// `a`c#t                        /select a,c from t
// {`a`sum!(a;a+b)}#t            /select a,sum:a+b from t
//    t 'w{a>3}                  /select from t where a>3
// `a#t 'w{a>3}                  /select a from t where a>3

// t 'w{a>3} / or where[t;{a>3}]
func where2(x, y uint32) uint32 { // t 'w bools | t 'w {a>3}
	_, cols := istab(x)
	if cols < 0 {
		panic("where: not a table")
	}
	if tp(y) == 0 && nn(y) == 5 {
		if MI[2+3+y>>2] != 0 {
			panic("where: only niladic lambdas allowed")
		}
		rx(x)
		y = expr(y, x)
	}
	if tp(y) != 2 {
		panic("where: y not ints|λ->ints")
	}
	return ecl(x, wer(y), 64) // x@\:&y
}
