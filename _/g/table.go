package main

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
		k, v := kvd(x)
		rxn(k, 2)
		env := ech(k, 46)   // save    env: .'!t
		dx(ecd(k, v, 58))   // (!t):'(.t)
		y = atx(y, 0)       // y:y[]
		dx(ecd(k, env, 58)) // restore (!t):'env
	}
	if tp(y) != 2 {
		panic("where: y not ints|λ->ints")
	}
	return ecl(x, wer(y), 64) // x@\:&y
}

// where:{x@\:&$[@y;y;((!x):env;y[];(!x):'(.x);env:.'!x)1]}
// t:`a`b`c!(1+!10;2+!10;3+!10);
// t 'w{a>3} / or where[t;{a>3}]
// where[t;{a>3}]
