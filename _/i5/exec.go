package k

func (k *K) fold(x []token, src []byte) (r []token) { // todo: *| sum ..
	if x == nil {
		return x
	}
	//fmt.Printf("fold>%#v\n", x)
	var p int
	k.ctx.push(src, &p)
	defer k.ctx.restore("fold")

	push := func(x ...token) { r = append(r, x...) }
	pop := func() token {
		n := len(r) - 1
		if n < 0 {
			panic("stack")
		}
		t := r[n]
		r = r[:n]
		return t
	}
	span := func(t T, x ...token) (r token) {
		r.t = t
		if len(x) > 0 {
			r.p = x[0].p
			r.e = x[0].e
			for _, t := range x {
				r.p = mini(r.p, t.p)
				r.e = maxi(r.e, t.e)
			}
		}
		return r
	}
s:
	for _, v := range x {
		//fmt.Printf("fold v=%T : %+v\n", v.t, v.t)
		p = v.p
		switch u := v.t.(type) {
		case adverb:
			f := pop()
			if con(f) {
				push(span(k.derive(u, f.t, sstr(f)+string(u)), v, f))
			} else {
				push(f, v)
			}
		case verb:
			push(span(k.Func[string(u)], v))
		case Call1:
			f := pop()
			x := pop()
			if con(f) && con(x) {
				push(span(f.t.(F1).call1(x.t), f, x, v))
			} else {
				push(x, f, v)
			}
		case Call2:
			f := pop()
			x := pop()
			y := pop()
			if con(f) && con(x) && con(y) {
				push(span(f.t.(F2).call2(x.t, y.t), f, x, y, v))
			} else {
				push(y, x, f, v)
			}
		case makelist:
			n := int(u)
			toks := r[len(r)-n:]
			for _, t := range toks {
				if con(t) == false {
					push(v)
					continue s
				}
			}
			var l L
			l.v = make([]T, n)
			for i := range l.v {
				l.v[i] = pop().t
			}
			l.init()
			push(span(l.uf(), toks...))
		default:
			push(v)
		}
	}
	return r
}
func con(x token) bool {
	switch x.t.(type) {
	case varname, makelist, Call1, Call2, copula, λ:
		return false
	default: // todo: rand-verb
		return true
	}
}

func (k *K) exec(x []token, src []byte) T {
	if x == nil {
		return x
	}
	var p int
	//fmt.Printf("exec>%#v\n", x)
	k.ctx.push(src, &p)
	defer k.ctx.restore("exec")

	for _, v := range x {
		p = v.p
		//fmt.Printf("exec v=%T [%d]: %+v\n", v.t, v.p, v.t)
		switch u := v.t.(type) {
		case verb:
			panic("no-verbs-here")
		case F12:
			k.push(u)
		case F1:
			k.push(u.call1(k.pop()))
		case F2:
			x := k.pop()
			k.push(u.call2(x, k.pop()))
		case Call1:
			f := k.pop()
			k.push(f.(F1).call1(k.pop()))
		case Call2:
			f := k.pop()
			x := k.pop()
			k.push(f.(F2).call2(x, k.pop()))
		case Call4:
			x := k.pop()
			y := k.pop()
			z := k.pop()
			if u {
				k.push(k.dmend(x, y, z, k.pop()))
			} else {
				k.push(k.amend(x, y, z, k.pop()))
			}
		case adverb:
			x := k.pop()
			k.push(k.derive(u, x, sstr(x)+string(u)))
		case makelist:
			var l L
			l.v = make([]T, int(u))
			for i := range l.v {
				l.v[i] = k.pop()
			}
			l.init()
			k.push(l.uf())
		case Link:
			x := k.pop()
			k.push(link(x, k.pop()))
		case copula:
			k.assign(u)
		case varname:
			k.push(k.Lookup(string(u)))
		case drp:
			dx(k.pop())
		default:
			k.push(u)
		}
	}
	return k.pop()
}

func (k *K) push(x T) { k.stack = append(k.stack, x) }

//; 1;2 /2
func (k *K) pop() T {
	n := len(k.stack) - 1
	if n < 0 {
		return nil
	}
	r := k.stack[n]
	k.stack = k.stack[:n]
	return r
}
