package main

import (
	"fmt"
	"reflect"
)

// kc ki kf kz ks  (k from byte/int/float64/complex128/string)
// kC kI kF kZ kS  (k from ..slices)
// ck ik fk zk sk  (.. from k without unref)
// cK iK fK zK sK  (..        with    unref)
// Ck Ik Fk Zk Sk  (.. slices from k without unref)
// CK IK FK ZK SK  (.. slices        with    unref)

func kc(x byte) (r uint32)       { return kC([]byte{x}) }
func ki(x int) (r uint32)        { return kI([]int{x}) }
func kf(x float64) (r uint32)    { return kF([]float64{x}) }
func kz(x complex128) (r uint32) { return kZ([]complex128{x}) }
func ks(x string) (r uint32)     { return sc(kC([]byte(x))) }

// func kC(x []byte) (r uint32) in k.go
func kI(x []int) (r uint32) {
	r = mk(2, uint32(len(x)))
	p := int(2 + r>>2)
	for i, u := range x {
		MI[p+i] = uint32(int32(u))
	}
	return r
}
func kF(x []float64) (r uint32) {
	r = mk(3, uint32(len(x)))
	p := int(1 + r>>3)
	for i, u := range x {
		MF[p+i] = u
	}
	return r
}
func kZ(x []complex128) (r uint32) {
	r = mk(4, uint32(len(x)))
	p := 1 + r>>3
	for _, u := range x {
		MF[p] = real(u)
		MF[p+1] = imag(u)
		p += 2
	}
	return r
}
func kS(x []string) (r uint32) {
	r = mk(5, uint32(len(x)))
	p := int(2 + r>>2)
	for i, u := range x {
		t := sc(kC([]byte(u)))
		MI[p+i] = I(t + 8)
		dx(t)
	}
	return r
}
func tc(t, x uint32) {
	if tp(x) != t {
		panic("type")
	}
}
func cK(x uint32) (r byte)         { r = ck(x); dx(x); return r }
func iK(x uint32) (r int)          { r = ik(x); dx(x); return r }
func fK(x uint32) (r float64)      { r = fk(x); dx(x); return r }
func zK(x uint32) (r complex128)   { r = zk(x); dx(x); return r }
func sK(x uint32) (r string)       { tc(5, x); x = cs(x); r = sk(x); dx(x); return r }
func ck(x uint32) (r byte)         { tc(1, x); return MC[8+x] }
func ik(x uint32) (r int)          { tc(2, x); return int(MI[2+x>>2]) }
func fk(x uint32) (r float64)      { tc(3, x); return MF[1+x>>3] }
func zk(x uint32) (r complex128)   { tc(4, x); return complex(MF[1+x>>3], MF[2+x>>3]) }
func sk(x uint32) (r string)       { tc(5, x); return string(Ck(I(I(kkey) + I(8+x)))) }
func CK(x uint32) (r []byte)       { r = Ck(x); dx(x); return r }
func IK(x uint32) (r []int)        { r = Ik(x); dx(x); return r }
func FK(x uint32) (r []float64)    { r = Fk(x); dx(x); return r }
func ZK(x uint32) (r []complex128) { r = Zk(x); dx(x); return r }
func SK(x uint32) (r []string)     { r = Sk(x); dx(x); return r }
func Ck(x uint32) (r []byte) {
	tc(1, x)
	n := nn(x)
	x += 8
	r = make([]byte, n)
	copy(r, MC[x:x+n])
	return r
}
func Ik(x uint32) (r []int) {
	tc(2, x)
	n := nn(x)
	r = make([]int, int(n))
	p := 2 + x>>2
	for i := uint32(0); i < n; i++ {
		r[i] = int(MI[p+i])
	}
	return r
}
func Fk(x uint32) (r []float64) {
	tc(3, x)
	n := nn(x)
	r = make([]float64, n)
	copy(r, MF[1+x>>3:1+n+x>>3])
	return r
}
func Zk(x uint32) (r []complex128) {
	tc(4, x)
	n := int(nn(x))
	r = make([]complex128, n)
	p := 1 + x>>3
	for i := range r {
		r[i] = complex(MF[p], MF[p+1])
		p += 2
	}
	return r
}
func Sk(x uint32) (r []string) {
	tc(5, x)
	n := int(nn(x))
	r = make([]string, n)
	p := int(2 + x>>2)
	for i := range r {
		r[i] = ski(MI[p+i])
	}
	return r
}
func ski(off uint32) string { return string(Ck(I(I(kkey) + off))) }

func K(x interface{}) (r uint32) { return kgo(reflect.ValueOf(x)) }
func G(x uint32, r interface{})  { gok(x, reflect.ValueOf(r).Elem()) } // does not unref
func kgo(v reflect.Value) (r uint32) {
	if v.IsValid() == false {
		return 0
	}
	switch v.Kind() {
	case reflect.Bool:
		r = mk(2, 1)
		setBool(r, 0, v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r = mk(2, 1)
		setInt(r, 0, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r = mk(2, 1)
		setUint(r, 0, v)
	case reflect.Float64:
		r = mk(3, 1)
		setFloat(r, 0, v)
	case reflect.Complex128:
		r = mk(4, 1)
		setComplex(r, 0, v)
	case reflect.String:
		s := v.String()
		r = mk(1, uint32(len(s)))
		copy(MC[r+8:], s)
	case reflect.Slice:
		r = kslice(v)
	case reflect.Struct:
		r = kstruct(v)
	case reflect.Map:
		r = kmap(v)
	case reflect.Ptr:
		r = kgo(v.Elem())
	case reflect.Interface:
		r = kgo(v.Elem())
	default:
		panic("cannot convert go type: " + v.Kind().String())
	}
	return r
}
func kslice(v reflect.Value) (r uint32) {
	n := v.Len()
	t := v.Type().Elem()
	switch t.Kind() {
	case reflect.Bool:
		r = mk(2, uint32(n))
		for i := 0; i < n; i++ {
			setBool(r, uint32(i), v.Index(i))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r = mk(2, uint32(n))
		for i := 0; i < n; i++ {
			setInt(r, uint32(i), v.Index(i))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r = mk(2, uint32(n))
		for i := 0; i < n; i++ {
			setUint(r, uint32(i), v.Index(i))
		}
	case reflect.Float64:
		r = mk(3, uint32(n))
		for i := 0; i < n; i++ {
			setFloat(r, uint32(i), v.Index(i))
		}
	case reflect.Complex128:
		r = mk(4, uint32(n))
		for i := 0; i < n; i++ {
			setComplex(r, uint32(i), v.Index(i))
		}
	default:
		r = mk(6, uint32(n))
		for i := 0; i < n; i++ {
			MI[2+uint32(i)+r>>2] = kgo(v.Index(i))
		}
	}
	return r
}
func kstruct(x reflect.Value) (r uint32) {
	k := mk(5, 0)
	v := mk(6, 0)
	t := x.Type()
	for i := 0; i < x.NumField(); i++ {
		u := x.Field(i)
		n := t.Field(i).Name
		if isexported(n) {
			k = ucat(k, ks(n))
			v = lcat(v, kgo(u))
		}
	}
	return mkd(k, v)
}
func kmap(m reflect.Value) (r uint32) {
	k := mk(5, 0)
	v := mk(6, 0)
	iter := m.MapRange()
	for iter.Next() {
		k = ucat(k, ks(iter.Key().String()))
		v = lcat(v, kgo(iter.Value()))
	}
	return mkd(k, v)
}

func setBool(r, i uint32, v reflect.Value) {
	if v.Bool() {
		MI[2+i+r>>2] = 1
	} else {
		MI[2+i+r>>2] = 0
	}
}
func setInt(r, i uint32, v reflect.Value)   { MI[2+i+r>>2] = uint32(v.Int()) }
func setUint(r, i uint32, v reflect.Value)  { MI[2+i+r>>2] = uint32(v.Uint()) }
func setFloat(r, i uint32, v reflect.Value) { MF[1+i+r>>3] = v.Float() }
func setComplex(r, i uint32, v reflect.Value) {
	z := v.Complex()
	MF[1+2*i+r>>3] = real(z)
	MF[2+2*i+r>>3] = imag(z)
}

func gok(x uint32, r reflect.Value) { // does not unref
	switch r.Kind() {
	case reflect.Bool:
		if tp(x) != 2 || nn(x) != 1 {
			panic("expected bool")
		}
		r.SetBool(MI[2+x>>2] != 0)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if tp(x) != 2 || nn(x) != 1 {
			panic("expected int")
		}
		r.SetInt(int64(MI[2+x>>2]))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if tp(x) != 2 || nn(x) != 1 {
			panic("expected int")
		}
		r.SetUint(uint64(MI[2+x>>2]))
	case reflect.Float64:
		if tp(x) != 3 || nn(x) != 1 {
			panic("expected float")
		}
		r.SetFloat(MF[1+x>>3])
	case reflect.Complex128:
		if tp(x) != 4 || nn(x) != 1 {
			panic("expected complex")
		}
		r.SetComplex(complex(MF[1+x>>3], MF[2+x>>3]))
	case reflect.String:
		if t := tp(x); t == 1 {
			r.SetString(string(Ck(x)))
		} else if t != 5 || nn(x) != 1 {
			panic(fmt.Sprintf("expected string, got %d/%d", t, nn(x)))
		} else {
			r.SetString(sk(x))
		}
	case reflect.Slice:
		n := nn(x)
		t := r.Type()
		v := reflect.MakeSlice(t, int(n), int(n))
		switch t.Elem().Kind() {
		case reflect.Int:
			u := Ik(x)
			for i := range u {
				v.Index(i).SetInt(int64(u[i]))
			}
		case reflect.Float64:
			u := Fk(x)
			for i := range u {
				v.Index(i).SetFloat(float64(u[i]))
			}
		case reflect.Complex128:
			u := Zk(x)
			for i := range u {
				v.Index(i).SetComplex(u[i])
			}
		case reflect.String:
			if xt := tp(x); xt == 5 {
				u := Sk(x)
				for i := range u {
					v.Index(i).SetString(u[i])
				}
			} else if xt == 6 {
				for i := uint32(0); i < n; i++ {
					v.Index(int(i)).SetString(string(Ck(MI[2+i+x>>2])))
				}
			} else {
				panic(fmt.Sprintf("expected type for []string: %d", xt))
			}
		default:
			if tp(x) != 6 {
				panic(fmt.Sprintf("unkown slice type: %v, expect general list", t))
			}
			for i := uint32(0); i < n; i++ {
				gok(MI[2+i+x>>2], v.Index(int(i)))
			}
		}
		r.Set(v)
	case reflect.Struct:
		if xt := tp(x); xt != 7 {
			panic(fmt.Errorf("expected dict: xt=%d %v", xt, r.Type()))
		}
		keys := Sk(MI[2+x>>2])
		v := MI[3+x>>2]
		n := nn(v)
		if n != uint32(exportedFields(r)) {
			panic("number of dict/struct fields mismatches")
		}
		t := r.Type()
		j := uint32(0)
		for i := 0; i < r.NumField(); i++ {
			name := t.Field(int(i)).Name
			if isexported(name) {
				s := keys[j]
				if s != name {
					panic("expected dict field: " + s)
				}
				gok(MI[2+j+v>>2], r.Field(int(i)))
				j++
			}
		}
	case reflect.Map:
		if tp(x) != 7 {
			panic("expected dict")
		}
		t := r.Type()
		m := reflect.MakeMap(t)
		keys := Sk(MI[2+x>>2])
		v := MI[3+x>>2]
		n := nn(v)
		kt := t.Key()
		vt := t.Elem()
		for i := uint32(0); i < n; i++ {
			kk := reflect.New(kt).Elem()
			vv := reflect.New(vt).Elem()
			kk.SetString(keys[i])
			gok(MI[2+i+v>>2], vv)
			m.SetMapIndex(kk, vv)
		}
		r.Set(m)
	case reflect.Ptr:
		if x != 0 {
			p := reflect.New(r.Type().Elem())
			gok(x, p.Elem())
			r.Set(p)
			return
		}
	case reflect.Interface: // ignore
		if x != 0 {
			v := autoInterface(x)
			r.Set(v)
			return
		}
	default:
		panic("cannot convert to go type: " + r.Kind().String())
	}
}
func autoInterface(x uint32) (r reflect.Value) {
	switch tp(x) {
	case 1:
		var v string
		G(x, &v)
		r = reflect.ValueOf(v)
	case 2:
		var v []int
		G(x, &v)
		r = reflect.ValueOf(v)
	case 3:
		var v []float64
		G(x, &v)
		r = reflect.ValueOf(v)
	case 4:
		var v []complex128
		G(x, &v)
		r = reflect.ValueOf(v)
	case 5, 6:
		var v []string
		G(x, &v)
		r = reflect.ValueOf(v)
	default:
		panic(fmt.Errorf("autoVector nyi %d/%d", tp(x), nn(x)))
	}
	return r
}
func isexported(s string) bool {
	h := s[0]
	return h >= 'A' && h <= 'Z'
}
func exportedFields(v reflect.Value) (n int) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if isexported(t.Field(i).Name) {
			n++
		}
	}
	return n
}
