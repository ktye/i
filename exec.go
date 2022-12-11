package main

import (
	. "github.com/ktye/wg/module"
)

type f1 = func(K) K
type f2 = func(K, K) K
type f3 = func(K, K, K) K
type f4 = func(K, K, K, K) K

func quoted(x K) int32 { return I32B(int32(x) >= 448 && tp(x) == 0) }
func quote(x K) K      { return x + 448 }
func unquote(x K) K    { return x - 448 }

func exec(x K) K {
	srcp = 0
	xn := nn(x)
	if xn == 0 {
		dx(x)
		return 0
	}
	a := K(0) // accumulator
	p := int32(x)
	e := p + 8*xn
	//kdb:var arg K
	for p < e {
		u := K(I64(p))
		if tp(u) != 0 {
			push(a)
			a = rx(u)
		} else {
			switch int32(u) >> 6 {
			case 0: //   0..63   monadic
				//kdb:arg=l1(rx(a))
				//kdb:fpush(u,arg)
				a = Func[marksrc(u)].(f1)(a)
				//kdb:dx(arg);fpop()
			case 1: //  64..127  dyadic
				//kdb:arg=l2(rx(a),rx(K(I64(sp-8))))
				//kdb:fpush(u,arg)
				a = Func[marksrc(u)].(f2)(a, pop())
				//kdb:dx(arg);fpop()
			case 2: // 128       dyadic indirect
				//kdb:arg=l2(rx(K(I64(sp-8))),rx(K(I64(sp-16))))
				marksrc(a)
				//kdb:b:=rx(a)
				//kdb:fpush(b,arg)
				a = Cal(a, l2(pop(), pop()))
				//kdb:dx(b);dx(arg);fpop()
			case 3: // 192..255  tetradic
				//kdb:arg=l2(l2(rx(a),rx(K(I64(sp-8)))),l2(rx(K(I64(sp-16))),rx(K(I64(sp-24)))))
				a = Func[marksrc(u)].(f4)(a, pop(), pop(), pop())
				//kdb:dx(arg);fpop()
			case 4: // 256       drop
				dx(a)
				a = pop()
			case 5: // 320       jump
				p += int32(a)
				a = pop()
			case 6: // 384       jump if not
				u = pop()
				p += int32(a) * I32B(int32(u) == 0)
				dx(u)
				a = pop()
			default: //448..     quoted verb
				push(a)
				a = rx(u - 448)
			}
		}
		//vcount: vcount(a)
		p += 8
		continue
	}
	dx(pop())
	dx(x)
	return a
}

//kdb:func fpush(f, x K){}
//kdb:func fpop(){}

func marksrc(x K) int32 {
	if p := 0xffffff & int32(x>>32); p != 0 {
		srcp = p
	}
	return int32(x)
}
func push(x K) {
	SetI64(sp, int64(x))
	sp += 8
	if sp == 512 {
		trap(Stack)
	}
}
func pop() K {
	sp -= 8
	if sp < 256 {
		trap(Stack)
	}
	//return K(I64(sp))
	r := K(I64(sp))
	return r
}
func lst(n K) K {
	rn := int32(n)
	r := mk(Lt, rn)
	rp := int32(r)
	for i := int32(0); i < rn; i++ {
		SetI64(rp, int64(pop()))
		rp += 8
	}
	return uf(r)
}
func nul(x K) K { push(x); return 0 }
func lup(x K) K {
	vp := I32(8) + int32(x)
	r := x0(K(vp))
	return r
}
func Asn(x, y K) K {
	if tp(x) != st {
		trap(Type)
	}
	vp := I32(8) + int32(x)
	dx(K(I64(vp)))
	SetI64(vp, int64(rx(y)))
	return y
}
func Amd(x, i, v, y K) K {
	xt := tp(x)
	if xt == st {
		return Asn(x, Amd(Val(x), i, v, y))
	}
	if xt < 16 {
		trap(Type)
	}
	if tp(i) == Lt { // @[;;v;]/[x;y;i]
		n := nn(i)
		for j := int32(0); j < n; j++ {
			x = Amd(x, ati(rx(i), j), rx(v), ati(rx(y), j))
		}
		dx(i)
		dx(v)
		dx(y)
		return x
	}

	if xt > Lt {
		r := x0(x)
		x = r1(x)
		if xt == Tt && tp(i)&15 == it { // table-assign-rows
			if tp(y) > Lt {
				y = Val(y)
			}
			return key(r, Dmd(x, l2(0, i), v, y), xt)
		}
		r = Unq(Cat(r, rx(i)))
		return key(r, Amd(ntake(nn(r), x), Fnd(rx(r), i), v, y), xt)
	}
	if i == 0 {
		if v == 1 {
			if tp(y) < 16 {
				y = ntake(nn(x), y)
			}
			dx(x)
			return y
		}
		return Cal(v, l2(x, y))
	}
	if tp(v) != 0 || v != 1 {
		y = cal(v, l2(Atx(rx(x), rx(i)), y))
	}
	ti, yt := tp(i), tp(y)
	if xt&15 != yt&15 {
		x, xt = explode(x), Lt
	}
	if ti == it {
		if xt != yt+16 {
			x = explode(x)
		}
		return sti(x, int32(i), y)
	}
	if yt < 16 {
		y = ntake(nn(i), y)
		yt = tp(y)
	}
	if xt == Lt {
		y = explode(y)
	}
	return stv(x, i, y)
}
func Dmd(x, i, v, y K) K {
	if tp(x) == st {
		return Asn(x, Dmd(Val(x), i, v, y))
	}
	i = explode(i)
	f := Fst(rx(i))
	if nn(i) == 1 {
		dx(i)
		return Amd(x, f, v, y)
	}

	if f == 0 {
		f = seq(nn(x))
	}
	i = ndrop(1, i)

	if tp(f) > 16 { // matrix-assign
		n := nn(f)
		if nn(i) != 1 {
			trap(Rank)
		}
		i = Fst(i)
		if tp(f) == It && tp(x) == Tt {
			t := rx(x0(x))
			return key(t, Dmd(r1(x), l2(Fnd(t, i), f), v, y), Tt)
		}
		if tp(f) != It || tp(x) != Lt {
			trap(Nyi) // Dt
		}
		x = use(x)
		for j := int32(0); j < n; j++ {
			rj := int32(x) + 8*I32(int32(f)+4*j)
			SetI64(rj, int64(Amd(K(I64(rj)), rx(i), rx(v), ati(rx(y), j))))
		}
		dx(f)
		dx(i)
		dx(v)
		dx(y)
		return x
	}
	x = rx(x)
	return Amd(x, f, 1, Dmd(Atx(x, f), i, v, y))
}

//vcount: func vcount(x K) {
//vcount: 	i := Cnt(rx(x))
//vcount: 	Printf("vcount %d\n", int32(i))
//vcount: }
