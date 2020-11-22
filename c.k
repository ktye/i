k:( `g!{[r;x]x+:3}; `ff!{[r;x;y]x/(x+:4;y:0);y} )  /x+:y ;y+2


funcs:{(,/!'x@i)!,/.'x@i:&,/7=@'x}
exports:{*|x}

/go a way
go:    {f:funcs x;o:gohead,/gofunc'[!f;.'.f];o}
gofunc:{""/:("func ";$x;"(";gosig y;")";goret y;"{";golocs y;gobody y;"}";_10)}
gosig: {a:1_x 2;n:-1+x 3;","/:(n#a){($x)," ",gotype y}'n#a}
goret: {$[`V~r:*x 2;"";"(",($r)," ",(gotype r),")"]}
gotype:{$["Z"<$x;"i";"f"]}
golocs:{l:(x 3)_x 2;$[#l;(";"/:{"var ",($x)," ",gotype x}'l),";";""]}
gobody:{b:x 1;";"/:return()gotree/b}
gotree:{p:-1+#x;m:16/y;n:16\y
$[~y;  ;
  ~m;  x,:,$%y
 1~m;  x[p]:en gmon[go1;_n-128;x p] 
 2~m;  x:(-2_x),,$[63~n;goif;64~n;goat;46~n;gocal;gody[go2;(_n)]][x p;x p-1]
 4~m;  x,:,$_n
 5~m;  x,:,$%y
 6~m;  x:((-n)_x),,(-n)#x
 7~m;  x:(-4_x),,goas -4#x
 9~m;  x[p]:"return ",x p
15~m;  ;
 $y];  x }

goas:{ (x 2),(x 3),"=",x 0}
goat:{|/"CIF"~/:x;gold[x;y];x,en y}
gocal:{$["?{"~x;gowhl y;x,en@ ","/:|y]}
gowhl:{"for ",x[1]," {",(";"/:|*x),"}"}
gold:{s:$["I"~x;">>2";"F"~x;">>3";""];x,"[",(en y),s,"]"}
goif:{"if ",x,"{",(";"/:|y),"}"}
go1:" "\:"+\ math.Abs( ~ n32( _\ math.Floor( * bits.LeadingZeros32( % math.Sqrt("
go2:" "\:": = </ <= </' <=' >/ >= >/' >=' ! % ^ << _ >>"
gmon:{s:$[(#x)~n:x?,y; y; x 1+n];$["("~*|s;s,z,")";s,z]}
gody:{[l;v;x;y]  s:$[(#l)~n:l?,v;v;l 1+n]; x,s,y}
en:{"(",x,")"}

return:{$[#x;;:x];p:-1+#x;$["return"~6#x p;x;x[p]:"return ",x[p]];x}

gohead:"package main
type c=byte
type i=uint32
type f=float64
var C []c;var I []i;var F []f;var T[256]interface{}
func init(){C=make([]c,16);msl()}
type sl struct {p uintptr;l,c int}
func msl() {
 cp := *(*sl)(unsafe.Pointer(&C));ip := *(*sl)(unsafe.Pointer(&I));jp := *(*sl)(unsafe.Pointer(&J));fp := *(*sl)(unsafe.Pointer(&F))
 ip.l, ip.c, ip.p = cp.l*4, cp.c*2, cp.p; fp.l, fp.c, fp.p = ip.l*2, ip.c*2, jp.p
 C = *(*[]c)(unsafe.Pointer(&cp));I = *(*[]i)(unsafe.Pointer(&ip));F = *(*[]f)(unsafe.Pointer(&fp))}
"

f:funcs k
g:{gofunc[x;.f x]}
