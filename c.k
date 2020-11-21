funcs:{(,/!'x@i)!,/.'x@i:&,/7=@'x}
exports:{*|x}

/go a way
go:    {f:funcs x;o:gohead,/gofunc'[!f;.'.f];o}
gofunc:{""/:("func ";$x;"(";gosig y;")";goret y;"{";golocs y;gobody y;"}";_10)}
gosig: {a:1_x 2;n:-1+x 3;","/:(n#a){($x)," ",gotype y}'n#a}
goret: {$[`V~r:*x 2;"";"(",($r)," ",(gotype r),")"]}
gotype:{$["Z"<$x;"i";"f"]}
golocs:{l:(x 3)_x 2;$[#l;(";"/:{"var ",($x)," ",gotype x}'l),";";""]}
gobody:{".."}
gohead:"package main
type c=byte
type i=uint32
type f=float64
var C []c;var I []i;var F []f;var T[256]interface{}
func init(){C=make([]c 16);msl()
type sl struct {p uintptr;l,c int}
func msl() {
 cp := *(*sl)(unsafe.Pointer(&C));ip := *(*sl)(unsafe.Pointer(&I));jp := *(*sl)(unsafe.Pointer(&J));fp := *(*sl)(unsafe.Pointer(&F))
 ip.l, ip.c, ip.p = cp.l*4, cp.c*2, cp.p; fp.l, fp.c, fp.p = ip.l*2, ip.c*2, jp.p
 C = *(*[]c)(unsafe.Pointer(&cp));I = *(*[]i)(unsafe.Pointer(&ip));F = *(*[]f)(unsafe.Pointer(&fp))}
"

f:funcs k
g:{gofunc[x;.f x]}
