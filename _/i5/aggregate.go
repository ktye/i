package k

/*
func sum(x T) T {
	switch v := x.(type) {
	case bool, byte, int, float64, complex128:
		return x
	case I:
		dx(x)
		return sumi(v.v)
	case F:
		dx(x)
		return sumf(v.v)
	case Z:
		dx(x)
		return sumz(v.v)
	default:
		return over(F2(add), x)
	}
}

func sumi(x []int) (r int) {
	for _, i := range x {
		r += i
	}
	return r
}
func sumf(x []float64) (r float64) {
	if n := len(x); n < 128 {
		for _, i := range x {
			r += i
		}
	} else {
		n /= 2
		r = sumf(x[:n]) + sumf(x[n:])
	}
	return r
}
func sumz(x []complex128) (r complex128) {
	if n := len(x); n < 128 {
		for _, i := range x {
			r += i
		}
	} else {
		n /= 2
		r = sumz(x[:n]) + sumz(x[n:])
	}
	return r
}
*/

/*
func maxover(x T) T {
	if v, o := x.(vector); o == false {
		return x
	} else {
		if n := v.ln(); n == 0 {
			dx(x)
			return v.zero()
		} else if n == 1 {
			return atv(v, 0)
		}
	}
	switch v := x.(type) {
	case B:
		return maxb(v.v)
	case C:
		return maxc(v.v)
	case I:
		return maxi(v.v)
	case F:
		return maxf(v.v)
	default:
		return over(F2(max), x)
	}
}
func minover(x T) T {
	if v, o := x.(vector); o == false {
		return x
	} else {
		if n := v.ln(); n == 0 {
			dx(x)
			return v.zero()
		} else if n == 1 {
			return atv(v, 0)
		}
	}
	switch v := x.(type) {
	case B:
		return minb(v.v)
	case C:
		return minc(v.v)
	case I:
		return mini(v.v)
	case F:
		return minf(v.v)
	default:
		return over(F2(min), x)
	}
}

func maxb(v []bool) (r bool) {
	for _, u := range v {
		if u {
			return true
		}
	}
	return false
}
func maxc(v []byte) (r byte) {
	r = v[0]
	for _, u := range v {
		if u > r {
			r = u
		}
	}
	return r
}
func maxi(v []int) (r int) {
	r = v[0]
	for _, u := range v {
		if u > r {
			r = u
		}
	}
	return r
}
func maxf(v []float64) (r float64) {
	r = v[0]
	for _, u := range v {
		r = math.Max(r, u)
	}
	return r
}
func minb(v []bool) (r bool) {
	for _, u := range v {
		if !u {
			return false
		}
	}
	return true
}
func minc(v []byte) (r byte) {
	r = v[0]
	for _, u := range v {
		if u < r {
			r = u
		}
	}
	return r
}
func mini(v []int) (r int) {
	r = v[0]
	for _, u := range v {
		if u < r {
			r = u
		}
	}
	return r
}
func minf(v []float64) (r float64) {
	r = v[0]
	for _, u := range v {
		r = math.Min(r, u)
	}
	return r
}
*/
