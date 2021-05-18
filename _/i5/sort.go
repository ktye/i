package k

import (
	srt "sort"
)

//^ ^6 3 2 3 /2 3 3 6
//^ x:6 3 2 3;^x /2 3 3 6
func sort(x T) T {
	switch v := x.(type) {
	case B:
		r := use(v).(B)
		srt.Slice(r.v, func(i, j int) bool { return blt(r.v[i], r.v[j]) })
		return r
	case C:
		r := use(v).(C)
		srt.Slice(r.v, func(i, j int) bool { return r.v[i] < r.v[j] })
		return r
	case I:
		r := use(v).(I)
		srt.Ints(r.v)
		return r
	case F:
		r := use(v).(F)
		srt.Float64s(r.v)
		return r
	case Z:
		r := use(v).(Z)
		srt.Slice(r.v, func(i, j int) bool { return zlt(r.v[i], r.v[j]) })
		return r
	case S:
		r := use(v).(S)
		srt.Strings(r.v)
		return r
	case L:
		r := use(v).(L)
		srt.Slice(r.v, func(i, j int) bool { return llt(r.v[i], r.v[j]) })
		return r
	default:
		panic("type")
	}
}
func gradeup(x T) T {
	xv, o := x.(vector)
	if o == false {
		panic("type")
	}
	r := seq(0, xv.ln())
	switch v := x.(type) {
	case B:
		srt.SliceStable(r.v, func(i, j int) bool { return blt(v.v[r.v[i]], v.v[r.v[j]]) })
	case C:
		srt.SliceStable(r.v, func(i, j int) bool { return v.v[r.v[i]] < v.v[r.v[j]] })
	case I:
		srt.SliceStable(r.v, func(i, j int) bool { return v.v[r.v[i]] < v.v[r.v[j]] })
	case F:
		srt.SliceStable(r.v, func(i, j int) bool { return v.v[r.v[i]] < v.v[r.v[j]] })
	case Z:
		srt.SliceStable(r.v, func(i, j int) bool { return zlt(v.v[r.v[i]], v.v[r.v[j]]) })
	case S:
		srt.SliceStable(r.v, func(i, j int) bool { return v.v[r.v[i]] < v.v[r.v[j]] })
	case L:
		srt.SliceStable(r.v, func(i, j int) bool { return llt(v.v[r.v[i]], v.v[r.v[j]]) })
	default:
		panic("type")
	}
	dx(x)
	return r
}
func gradedown(x T) T {
	xv, o := x.(vector)
	if o == false {
		panic("type")
	}
	r := seq(0, xv.ln())
	switch v := x.(type) {
	case B:
		srt.SliceStable(r.v, func(j, i int) bool { return blt(v.v[r.v[i]], v.v[r.v[j]]) })
	case C:
		srt.SliceStable(r.v, func(j, i int) bool { return v.v[r.v[i]] < v.v[r.v[j]] })
	case I:
		srt.SliceStable(r.v, func(j, i int) bool { return v.v[r.v[i]] < v.v[r.v[j]] })
	case F:
		srt.SliceStable(r.v, func(j, i int) bool { return v.v[r.v[i]] < v.v[r.v[j]] })
	case Z:
		srt.SliceStable(r.v, func(j, i int) bool { return zlt(v.v[r.v[i]], v.v[r.v[j]]) })
	case S:
		srt.SliceStable(r.v, func(j, i int) bool { return v.v[r.v[i]] < v.v[r.v[j]] })
	case L:
		srt.SliceStable(r.v, func(j, i int) bool { return llt(v.v[r.v[i]], v.v[r.v[j]]) })
	default:
		panic("type")
	}
	dx(x)
	return r
}
