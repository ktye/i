package i

import "reflect"

// monadic verbs

func idn(x interface{}) interface{} { return x }
func flp(x interface{}) interface{} { return e("TODO") }
func neg(x interface{}) interface{} { return nm(x, rneg, zneg, "Neg") }
func fst(v interface{}) interface{} {
	//function first (x) { return (x.t == 4) ? first(x.v) : (x.t != 3) ? x : len(x) ? x.v[0]:k(3,[]); }
	// TODO dict
	if n := ln(v); n < 0 {
		return v
	} else if n == 0 {
		return nil
	}
	return at(v, 0)
}
func sqr(x interface{}) interface{} { return nm(x, rsqr, zsqr, "Sqr") }
func iot(x interface{}) interface{} { return e("TODO") }
func odo(x interface{}) interface{} { return e("TODO") }
func wer(x interface{}) interface{} { return e("TODO") }
func rev(x interface{}) interface{} { return e("TODO") }
func asc(x interface{}) interface{} { return e("TODO") }
func dsc(x interface{}) interface{} { return e("TODO") }
func eye(x interface{}) interface{} { return e("TODO") }
func grp(x interface{}) interface{} { return e("TODO") }
func not(x interface{}) interface{} { return nm(x, rnot, znot, "Not") }
func enl(x interface{}) interface{} {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case reflect.Func, reflect.Slice, reflect.Map:
		return []interface{}{x}
	}
	l := reflect.MakeSlice(reflect.SliceOf(v.Type()), 1, 1)
	l.Index(0).Set(v)
	return l.Interface()
}
func is0(x interface{}) interface{} { return e("TODO") }
func cnt(x interface{}) interface{} { return e("TODO") }
func flr(x interface{}) interface{} { return e("TODO") }
func fmt(x interface{}) interface{} { return e("TODO") }
func fgn(x interface{}) interface{} { return e("TODO") }
func unq(x interface{}) interface{} { return e("TODO") }
func evl(x interface{}) interface{} { return e("TODO") }

// dyadic verbs

func add(x, y interface{}) interface{} { return nd(x, y, radd, zadd, "Add") }
func sub(x, y interface{}) interface{} { return nd(x, y, rsub, zsub, "Sub") }
func mul(x, y interface{}) interface{} { return nd(x, y, rmul, zmul, "Mul") }
func div(x, y interface{}) interface{} { return nd(x, y, rdiv, zdiv, "Div") }
func mod(x, y interface{}) interface{} { return e("TODO") }
func mkd(x, y interface{}) interface{} { return e("TODO") }
func min(x, y interface{}) interface{} { return nd(x, y, rmin, zmin, "Min") }
func max(x, y interface{}) interface{} { return nd(x, y, rmax, zmax, "Max") } // cast to bool?
func les(x, y interface{}) interface{} { return nd(x, y, rles, zles, "Les") } // ?
func mor(x, y interface{}) interface{} { return nd(x, y, rmor, zmor, "Mor") } // ?
func eql(x, y interface{}) interface{} { return nd(x, y, reql, zeql, "Eql") } // ?
func mch(x, y interface{}) interface{} { return e("TODO") }
func cat(x, y interface{}) interface{} {
	// if (x.t==4&&y.t==4) { x=c(x); kmap(y.k, function(v) { dset(x,v,dget(y,v)); }); return x; };
	// return k(3, (x.t==3?x.v:[x]).concat(y.t==3?y.v:[y]));

	// TODO dict
	nx := ln(x)
	if nx < 0 {
		x = enl(x)
		nx = 1
	}
	ny := ln(y)
	if ny < 0 {
		y = enl(y)
		ny = 1
	}
	if t := reflect.TypeOf(x); t == reflect.TypeOf(y) {
		var l reflect.Value
		l = reflect.MakeSlice(t, nx+ny, nx+ny)
		for i := 0; i < nx; i++ {
			l.Index(i).Set(reflect.ValueOf(at(x, i)))
		}
		for i := 0; i < ny; i++ {
			l.Index(nx + i).Set(reflect.ValueOf(at(y, i)))
		}
		return l.Interface()
	}
	l := make([]interface{}, nx+ny)
	for i := 0; i < nx; i++ {
		l[i] = at(x, i)
	}
	for i := 0; i < ny; i++ {
		l[i+nx] = at(y, i)
	}
	return l
}
func tak(x, y interface{}) interface{} { return e("TODO") }
func rsh(x, y interface{}) interface{} {
	// if (y.t == 4) { return md(x, atx(y, x)); }
	// if (y.t != 3) { y = enlist(y); }
	// var a = first(x); var b = x.v[len(x)-1]; var c = 0;
	// function rshr(x, y, i) {
	// 	return krange(x.v[i].v, function(z) {
	// 		return i==len(x)-1 ? y.v[kmod(c++, len(y))] : rshr(x, y, i+1);
	// 	});
	// }
	// return na(a) ? (!len(y) ? y : cut(krange(len(y)/b.v, function(z) { return k(0, z*b.v); }), y)) :
	//        na(b) ? cut(krange(a.v, function(z) { return k(0, Math.floor(z*len(y)/a.v)); }), y) :
	//        rshr(l(x), len(y) ? y : enlist(y), 0);
	/*
		nx := ln(x)
		if nx < 0 {
			e("type")
		} else if nx == 0 {
			return nil
		}

		ny := ln(y)
		if ny < 0 {
			y = enl(y)
			ny = 1
		}

		a := fst(x)
		b := at(x, n-1)
		c := 0
		rshr := func(x, y interface{}, i int) {
			return krange(idx(at(x, i)), func(_ int) {
				if i == ln(x)-1 {
					c++
					return at(y, c%ln(y))
				}
				return rshr(x, y, i+1)
			})
		}
	*/

	return e("TODO")
	//return na(a) ? (!len(y) ? y : cut(krange(len(y)/b.v, function(z) { return k(0, z*b.v); }), y)) :
	//      na(b) ? cut(krange(a.v, function(z) { return k(0, Math.floor(z*len(y)/a.v)); }), y) :
	//     rshr(l(x), len(y) ? y : enlist(y), 0);
	/*
		if na(a) {
			// (!len(y) ? y : cut(krange(len(y)/b.v, function(z) { return k(0, z*b.v); }), y))
			if ny == 0 {
				return y
			}

		} else if na(b) {
				// return cut(...
			}
		}
	*/
}
func fil(x, y interface{}) interface{} { return e("TODO") }
func drp(x, y interface{}) interface{} { return e("TODO") }
func cut(x, y interface{}) interface{} {
	// return kzip(x, cat(drop(k1,x),count(y)), function(a, b) { // {x{x@y+!z-y}[y]'1_x,#y} ?
	// 	var r=[]; for(var z=p(a);z<p(b);z++) { r.push(lget(y,z)); } return k(3,r);
	// });
	return e("TODO")
}
func cst(x, y interface{}) interface{} { return e("TODO") }
func rnd(x, y interface{}) interface{} { return e("TODO") }
func fnd(x, y interface{}) interface{} { return e("TODO") }
func pik(x, y interface{}) interface{} { return e("TODO") }
func rfd(x, y interface{}) interface{} { return e("TODO") }
func atx(x, y interface{}) interface{} { return e("TODO") }
func cal(x, y interface{}) interface{} { return e("TODO") }
func bin(x, y interface{}) interface{} { return e("TODO") }
func rbn(x, y interface{}) interface{} { return e("TODO") }
func pak(x, y interface{}) interface{} { return e("TODO") }
func upk(x, y interface{}) interface{} { return e("TODO") }
func spl(x, y interface{}) interface{} { return e("TODO") }
func win(x, y interface{}) interface{} { return e("TODO") }
