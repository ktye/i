/k:( `mk!{[r;x;y]j:i-4;(j>/4*t)?/(u:a+1^j_2;I[u]:I j;I[j]:u;j-:4);I[a]:y|x^29;I[a+4]:1;a}; 
/k:( `match!{[r;x;y]r}
/`f!{[r;x]r})

/k:60#k

funcs:{(,/!'x@i)!,/.'x@i:&,/7=@'x}
exports:{*|x}
obj:{ 'ox}

/go a way
go:    {f:funcs x;  o:gohead,/gofunc'[!f;.'.f];o}
gofunc:{""/:("func ";$x;"(";gosig y;")";goret y;"{";golocs y;gobody y;"}";_10)}
gosig: {a:1_x 2;n:-1+x 3;","/:(n#a){($x)," ",gotype y}'n#a}
goret: {$[`V~r:*x 2;"";"(",($r)," ",(gotype r),")"]}
gotype:{$["Z"<$x;"i";"f"]}
golocs:{l:(x 3)_x 2;$[#l;(";"/:{"var ",($x)," ",gotype x}'l),";";""]}
gobody:{b:x 1;";"/:return[`V~*x 2;()gotree/b]}
gotree:{p:-1+#x;m:16/y;n:16\y
$[~y;
 0~m;  x,:,$%y
 1~m;  x[p]:$[164~n;gocond[|x p];128>n;goforc x p;en gmon[go1;_n-128;x p]]
 2~m;  x:(-2_x),,$[63~n;goif;64~n;goat;46~n;gocal;gody[go2;(_n)]][x p;x p-1]
 4~m;  x,:,$_n
 5~m;  x,:,$%y
 6~m;  x:((-n)_x),,(-n)#x
 7~m;  x:(-4_x),,goas -4#x
 9~m;  x[p]:"return ",x p
15~m;
 $y];  x }

goas:{$[x 1;gost x;(x 2),(x 3),"=",x 0]}
gocond:{(5_,/{"else if ",(x 0),"{",(";"/:x 1),"}"}'((2\#x),2)#x),"else{",";"/:(*|x),"}"  }
goforc:{"for i:=i(0);i<",x,";i++"}
goforb:{x,"{",(";"/:y),"}"}
goat:{$[|/"CIF"~\:x;gold[x;y];"for"~3#x;goforb[x;y];x,en y]}
gocal:{$["?{"~x;gowhl y;|/go2~\:x; gody[go2;x;y 1;y 0]; "T"~x;goicl y; x,en@ ","/:rev y]}
goicl:{a:","/:-x 0;r:-x 1;"T[",(x 2),"].(func(",a,")",r,")"}
gowhl:{"for ",x[1]," {",(";"/:rev@*x),"}"}
gold:{x,"[",(en y),(shift x),"]"}
gost:{(x 2),"[",(en x 1),(shift x 2),"]=",x 0}
goif:{"if ",x,"{",(";"/:rev y),"}"}
go1:" "\:"+[ math.Abs( ~ n32( _\ math.Floor( * bits.LeadingZeros32( % math.Sqrt( ! trap("
go2:" "\:"~ == <( <' <{ <= <{( <=' >( >' >{ >= >{( >=' >{[ >= >[ > > >  % / ! % !( %' ^ << _ >> $( $' +[ + + + -[ - - - *[ * * * %( /' %[ /"
gmon:{s:$[(#x)~n:x?,y; y; x 1+n];$["("~*|s;s,z,")";s,z]}
gody:{[l;v;x;y]$[v~":"; :x,"=","y";v~"$"; :gocvt[x;y;""];v~"$("; :gocvt[x;y;"'"]; 0];s:$[(#l)~n:l?,v;v;l 1+n];"(",sign[s;x],sign[s;s],sign[s;y],")"}
gocvt:{(-x),"(",sign[z;y],")"}
en:{"(",x,")"}
ne:{$["("~*x;1_-1_x;x]}
shift:{$["I"~x;">>2";"F"~x;">>3";""]}
sign:{o:$["'"~*|x;-1_x;x];$[x~y;o;$[o~x;y;"s(",y,")"]]}
rev:{$[6~@x;|x;x]}

return:{$[x;:y;#y;;:y];p:-1+#y;$["return"~6#y p;y;y[p]:"return ",y[p]];y}

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
