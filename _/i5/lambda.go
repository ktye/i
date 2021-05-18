package k

//λ {1+2} /{1+2}
//λ {1+2}[] /3
//λ {x+2}[2] /4
//λ {x+y}[2;3] /5
//λ {[a;b]3*a+b}[3;4] /21
//λ {[]3}[] /3
//λ {(a;y*a:x)}[2;3] /2 6
//λ {a+y*a:x}[2;3] /8
//λ a+{2*a:x}[2]+a:1 /6  (local assign)
//λ a+{2*a::x}[2]+a:1 /7  (global assign)
type λ struct {
	refcount
	src  []byte
	code []token
	save []T
	loc  []string
	ary  int
}

func (l λ) call(k *K) T {
	for i, s := range l.loc {
		l.save[i] = k.Var[s]
		k.Var[s] = nil
	}
	for i := l.ary - 1; i >= 0; i-- {
		k.Var[l.loc[i]] = k.pop()
	}

	for i := range l.code {
		rx(l.code[i])
	}
	r := k.exec(l.code, l.src)

	for i, s := range l.loc {
		if i < l.ary {
			dx(k.Var[s])
		}
		k.Var[s] = l.save[i]
	}
	l.unref()
	return r
}
func (l λ) dict(k *K) T {
	for i, s := range l.loc {
		l.save[i] = k.Var[s]
		k.Var[s] = nil
	}

	for i := range l.code {
		rx(l.code[i])
	}
	dx(k.exec(l.code, l.src))

	key := KS(l.loc)
	val := KL(make([]T, len(l.loc)))

	for i, s := range l.loc {
		val.v[i] = k.Var[s]
		k.Var[s] = l.save[i]
	}
	l.unref()
	return dict(key, val.uf())
}

func (l λ) unref() {
	if l.refcount.unref() == 0 {
		for _, c := range l.code {
			switch v := c.t.(type) {
			case refcounter:
				v.unref()
			}
		}
	}
}

func (l λ) String() string { return string(l.src) }
