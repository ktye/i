sin:F:F{}cos:F:F{}atan2:F:FF{}
ini:I:I{0::289360742959022336j;128::x;p:256;i:8;(i<x)?/((4*i)::p;p*:2;i+:1);x} /x:16(64k)
bk:I:II{r:32-*7+y*I?C x;(r<4)? :4;r}mk:I:II{t:x bk y;i:4*t;(~I i)?/i+:4;(128~i)?!;a:I i;i::I a;j:i-4;(j>=4*t)?/(u:a+1<<j>>2;u::I j;j::u;j-:4);a::y|x<<29;(a+4)::1;a}
mki:I:I{r:2 mk 1;(r+8)::x;r}mkf:I:F{r:3 mk 1;(r+8)::x;r}mkd:I:II{v2;ext;r:7 mk 2;(r+8)::x;(r+12)::y;r}mkz:I:II{r:x mkd y;r::2|6<<29;r}l2:I:II{r:6 mk 2;(r+8)::x;(r+12)::y;r}
v1:{xt:(I x)>>29;xn:(I x)&536870911;xp:8+x}v2:{v1;yt:(I y)>>29;yn:(I y)&536870911;yp:8+y}r8:{rp:r+8}
fr:V:I{v1;t:4*xt bk xn;x::I t;t::x}dx:V:I{(x>255)?(xr:I x+4;(x+4)::xr-1;(1~xr)?(v1;(xt>4)?xn/(dx I xp+4*i);fr x))}dxr:{dx x;r}dxyr:{dx x;dx y;r}rx:V:I{x rxn 1}rxn:V:II{(x>255)?(x+:4;x::y+I x)}rl:V:I{v1;xn/(rx I xp;xp+:4)}
lx:I:I{v1;(xt~6)? :x;r:6 mk xn;r8;xn/(rx x;rp::x atx mki i;rp+:4);dxr}
til:I:I{v1;(~2~xt)?!;n:I xp;dx x;(n<'0)? :tir -n;seq(0;n;1)}seq.I:III{r:2 mk y;rp:8+r;y/(rp::z*i+x;rp+:4);r}tir.I:I{r:2 mk x;rp:4+r+4*x;x/(rp::i;rp-:4);r}
ext:{(~xt~yt)?!;((xn~1)&yn>1)?(x:x take yn;xn:yn;xp:x+8);((yn~1)&xn>1)?(y:y take xn;yn:xn;yp:y+8);(~xn~yn)?!}
upx:{(xt>=5)?!;(yt>=5)?!;(xt<yt)?/(x:up(x;xt;xn);xt+:1);(yt<xt)?/(y:up(y;yt;yn);yt+:1);xp:x+8;yp:y+8}
up:I:III{r:(y+1) mk z;xp:x+8;r8;y?[;z/(rp::I?C xp+i;rp+:4);z/(rp::F?I xp;rp+:8;xp+:4);!];dxr}
atx:I:II{v2;(~xt)?( :x cal enl y);(~yt~2)?!;r:xt mk yn;r8;w:I?C xt;f:xt+128;yn/(yi:I yp;$[yi<xn;mv(rp;xp+w*yi;w);(V.f)rp];rp+:w;yp+:4);(xt>4)?rl r;(yn~1)?(xt~6)?(rx I r+8;dx r;r:I r+8);dxyr}
cal:I:II{v2;(~yt~6)?!;(yn~1)?(((sadv x)|sadv x-128)? :((I.x)(fst y)));(x<128)?((~yn~1)?!; :(I.x)(fst y));(x<256)?((~yn~2)?!;rl y;dx y; :(I.x)(I yp;I yp+4));(xn~2)?(rl x;dx x;a:I xp;yn?[; :((I.a)(fst y;I xp+4));(a+:128;rl y;dx y; :((I.a)(I yp;I yp+4;I xp+4)));!]);!;x}
rev:I:I{v1;(~xn)? :x;x atx tir xn}fst:I:I{v1;(xt~7)?(rx 12+x;dx x; :fst 12+x);x atx mki 0}
cut:I:II{v2;(~xt~2)?!;(xn~1)?(r:y drop I xp;dx x; :r);r:6 mk xn;r8;xn/(a:I xp;b:I xp+4;(i~xn-1)?b:yn;(b<a)?!;rx y;rp::y atx seq(a;b-a;1);xp+:4;rp+:4);dxyr}
rsh:I:II{v2;(~xt~2)?!;n:prod(xp;xn);r:y take n;(xn~1)?(dx x; :r);xn-:1;xe:xp+4*xn;xn/(m:I xe;n:n%m;n:xp prod xn-i;r:(seq(0;n;m))cut r;xe-:4);dxr}prod:I:II{r:1;y/(r*:I x;x+:4);r}
take:I:II{v1;r:til mki y;r8;(xn<y)?(y/(rp+4*i)::i\xn);x atx r}drop:I:II{v1;(y>xn)?!;x atx seq(y;xn-y;1)}
use:I:I{(1~I x+4)? :x;v1;r:xt mk xn;r8;mv(rp;xp;xn*I?C xt);dx x;r}mv:V:III{z/(x+i)::C y+i}
cat:I:II{v2;((~xt)|~yt)?!;(xt~yt)? :x ucat y;(xt~6)? :x lcat y;!;x}
ucat:I:II{v2;(xt>4)?(rl x);(xt>5)?(r:(x+8)mkd x+12;dx x;dx y; :r);r:xt mk xn+yn;w:I?C xt;mv(r+8;xp;w*xn);mv(r+8+w*xn;yp;w*yn);dxyr}
lcat:I:II{v1;((xt bk xn)<(xt bk xn+1))?(r:xt mk xn+1;rl x;dx x;mv(r+8;xp;4*xn);x:r;xp:x+8);(xp+4*xn)::y;x::(xn+1)|6<<29;x}
enl:I:I{(6 mk 0) lcat x} cnt:I:I{v1;dx x;mki xn}typ:I:I{v1;r:2 mk 1;(8+r)::xt;dxr} not:I:I{x eql mki 0}
wer:I:I{v1;(~xt~2)?!;n:0;xn/(n+:I xp;xp+:4);xp:8+x;r:2 mk n;r8;xn/((I xp)/(rp::i;rp+:4);xp+:4);dxr}
mtc:I:II{r:2 mk 1;(r+8)::x match y;dxyr}match:I:II{(x~y)? :1;(~(I x)~I y)? :0;v1;yp:y+8;m:0;xt?[ :1;nn:xn;nn:xn<<2;nn:xn<<3;nn:xn<<4;(xn/((~((I xp) match I yp))? :0;xp+:4;yp+:4); :1)];nn/(~(C xp+i)~C yp+i)? :0;1}
fnd:I:II{v2;(~xt~yt)?!;r:2 mk yn;r8;w:I?C yt;yn/(rp::x fnx yp;rp+:4;yp+:w);dxyr}fnx:I:II{v1;eq:8+xt;w:I?C xt;xn/(((I.eq)(xp;y))? :i;xp+:w);xn}
exc:I:II{rx x;x atx wer (mki (I y)&536870911)eql y fnd x}
srt:I:I{rx x;x atx grd x}gdn:I:I{rev grd x}grd:I:I{v1;r:seq(0;xn;1);y:seq(0;xn;1);r8;msrt(y+8;rp;0;xn;xp;xt);dxyr}
msrt:V:IIIIII{((x3-z)>=2)?(c:(x3+z)%2;msrt(y;x;z;c;x4;x5);msrt(y;x;c;x3;x4;x5);mrge(x;y;z;x3;c;x4;x5))}
mrge:V:IIIIIII{k:z;j:x4;w:I?C x6;i:z;(i<x3)?/(c:k>=x4;(~c)?$[j>=x3;c:0;c:(I.x6)(x5+w*I x+k<<2;x5+w*I x+j<<2)];$[c;(a:j;j+:1);(a:k;k+:1)];(y+i<<2)::I x+a<<2;i+:1)}
gtc:I:II{(C x)>C y}gti:I:II{(I x)>'I y}gtf:I:II{(F x)>F y} eqc:I:II{(C x)~C y}eqi:I:II{(I x)~ I y}eqf:I:II{(F x)~F y}eqz:I:II{(x eqf y)? :(x+8)eqf y+8;0}eqL:I:II{(I x)match I y}
gtl:I:II{x:I x;y:I y;v2;(~xt~yt)? :xt>yt;n:xn;(yn<xn)?n:yn;w:I?C xt;n/(a:xp+i*w;b:yp+i*w;((I.xt)(a;b))? :1;((I.xt)(b;a))? :0);xn>yn}
eql:I:II{cmp(x;y;1)}mor:I:II{cmp(x;y;0)}les:I:II{cmp(y;x;0)}cmp:I:III{v2;upx;ext;f:xt;z?f+:8;w:I?C xt;r:2 mk xn;r8;xn/(rp::(I.f)(xp;yp);xp+:w;yp+:w;rp+:4);dxyr}
nd:I:III{v2;upx;ext;w:I?C xt;f:z+xt;r:xt mk xn;r8;xn/((V.f)(xp;yp;rp);xp+:w;yp+:w;rp+:w);dxyr}
nm:I:II{v1;r:use x;r8;w:I?C xt;y+:xt;xn/(((V.y)(xp;rp));xp+:w;rp+:w);(xt~4)?(y~19)? :zre r;r}
add:I:II{nd(x;y;143)}sub:I:II{nd(x;y;147)}mul:I:II{nd(x;y;151)}diw:I:II{nd(x;y;155)}
adc:V:III{z::(C x)+C y}adi:V:III{z::(I x)+I y}adf:V:III{z::(F x)+F y}adz:V:III{adf(x;y;z);adf(x+8;y+8;z+8)}
suc:V:III{z::(C x)-C y}sui:V:III{z::(I x)-I y}suf:V:III{z::(F x)-F y}suz:V:III{suf(x;y;z);suf(x+8;y+8;z+8)}
muc:V:III{z::(C x)*C y}mui:V:III{z::(I x)*I y}muf:V:III{z::(F x)*F y}muz:V:III{z::((F x)*F x+8)-(F y)*F y+8;(z+8)::((F x)*F y+8)+(F x+8)*F y}
dic:V:III{z::(C x)%C y}dii:V:III{z::(I x)%I y}dif:V:III{z::(F x)%F y}
diz:V:III{a:F x;b:F x+8;c:F y;d:F y+8;$[(+c)>=(+d);(r:d%c;p:c+r*d;z::(a+b*r)%p;(z+8)::(b-a*r)%p);(r:c%d;p:d+r*c;z::(a*r+b)%p;(z+8)::(b*r-a)%p)]}
abs:I:I{x nm 15}neg:I:I{x nm 19}sqr:I:I{x nm 27}
abc:V:II{c:I?C x;$[craz c;y::C?c-32;y::C?c]}abi:V:II{i:I x;$[(i<'0);y::0-i;y::i]}abf:V:II{y::+F x}abz:V:II{a:F x;b:F x+8;y::%a*a+b*b}
nec:V:II{c:I?C x;$[crAZ c;y::C?c+32;y::C?c]}nei:V:II{y::0-I x}nef:V:II{y::-F x}nez:V:II{y::-F x;(y+8)::-F x+8}
sqc:V:II{!}sqi:V:II{!}sqf:V:II{y::%F x}sqz:V:II{y::F x;(y+8)::-F x+8}
zre:I:I{x zri 0}zim:I:I{x zri 8}zri:I:II{v1;(xt~4)?!;r:3 mk xn;r8;xp+:y;xn/(rp::F xp;rp+:8;xp+:16);dxr}
crAZ:I:I{(x>64)?(x<91)? :1;0}craz:I:I{(x>96)?(x<123)? :1;0}
nag:V:I{x::0}nac:V:I{x::C?32}nai:V:I{x::-2147483648}naf:V:I{x::9221120237041090561f}naz:V:I{naf x;naf x+8}nas:V:I{x::1 mk 0;(4+I x)::0}nal:V:I{x::0}
drv:I:II{r:0 mk 2;(r+8)::x;(r+12)::y;r}ecv:I:I{40 drv x}epv:I:I{41 drv x}ovv:I:I{123 drv x}scv:I:I{91 drv x}
ech:I:II{x:lx x;v1;r:6 mk xn;r8;rl x;(y<256)?(y>127)?y-:128;xn/(rx y;rp::y atx I xp;xp+:4;rp+:4);dxyr}
ecp:I:II{v1;(~xn)?(dx y; :fst x);x rxn -1+2*xn;y rxn xn-1;r:fst x;(xn-1)/r:r cat y cal (x atx mki 1+i)l2 x atx mki i;dxyr}
ovr:I:II{v1;(~xn)?(dx y; :fst x);x rxn xn;y rxn xn-1;r:fst x;(xn-1)/r:y cal r l2 x atx mki i+1;dxyr}
scn:I:II{v1;(~xn)?(dx y; :fst y);x rxn xn;y rxn xn-1;t:fst x;rx t;r:enl t;(xn-1)/(t:y cal t l2 x atx mki 1+i;rx t;r:r lcat t);dx t;dxyr}
ecr:I:III{v2;r:6 mk yn;r8;x rxn yn;y rxn yn;z rxn yn;yn/(rp::z cal x l2 y atx mki i;rp+:4);dx z;dxyr}
ecl:I:III{v2;r:6 mk xn;r8;x rxn xn;y rxn yn;z rxn yn;xn/(rp::z cal (x atx mki i)l2 y;rp+:4);dx z;dxyr}
ecd:I:III{!;x}epi:I:III{!;x}
val:I:I{v1;xt?[;;;;;r:evl x;(rx x+12;r:x+12;dx x);!];r}
evl:I:I{v1;(~xt~6)? :x;(xn~0)? :x;(xn~1)? :fst x;rl x;r:6 mk xn;r8;xn/(rp::evl I xp;rp+:4;xp+:4);r8;dx x;(2~xn)?(rl r;dx r; :(I rp)atx I rp+4);(~3~xn)?!;rx I rp;(I rp)cal r drop 1}
xx1:I:I{!;x}xx2:I:II{!;x}str:I:I{!;x}grp:I:I{!;x}unq:I:I{!;x}flr:I:I{!;x}cst:I:II{!;x}min:I:II{!;x}max:I:II{!;x}
sadv:I:I{$[x~39;1;x~47;1;x~92;1;0]}
1:{gtc;gti;gtf;gtl;gtl}9:{eqc;eqi;eqf;eqz;eqL;eqL}
16:{abc;abi;abf;abz;nec;nei;nef;nez}28:{sqc;sqi;sqf;sqz}
33:{til;xx1;cnt;str;sqr;wer;epv;ech;ecp;fst;abs;enl;neg;val}
60:{grd;grp;gdn;unq;typ}91:{scn;xx1;xx1;srt;flr}123:{ovr;rev;xx1;not}
128:{nag;nac;nai;naf;naz;nas;nal} // todo naz/nas
144:{adc;adi;adf;adz;suc;sui;suf;suz;muc;mui;muf;muz;dic;dii;dif;diz}
161:{mkd;xx2;rsh;cst;diw;min;ecv;ecd;epi;mul;add;cat;sub;cal;ovv}
188:{les;eql;mor;fnd;atx}219:{ecl;scv;xx1;exc;cut}251:{ecr;max;xx2;mtc}


\
odo5:{p#’x(&#)’_(p:*/x)%*\x}
flip:,'/,''

01234567   xt:x>>29       xn:x&536870911 (-1+1<<29)
Fcifzsld   xt~0(function) x<256(basic) 
0148x444   xn~2(derived)  adv  verb
	   xn~3(proj)     verb argv empty-index
	   xn~4(lambda)   str  tree args locals
	   
+  add abs                 abs:+z                   memory
-  sub neg                                          0..  7   type sizes   0 1 4 8 16 4 4 0
*  mul fst                                          8.. 11   k-tree/key   pointer     todo
%  div sqr                                         12.. 15   k-tree/value pointer     todo
&  min wer                                         16..127   free pointers (4*i) for bt i, i:4..31
|  max rev                                        128..131   memsize log2
<  les grd                                        132..135   k-tree root (dict)
>  mor gdn                                        function tables (call indirect)
=  eql grp                                        T1:I:I (#128)      T2:I:I (#128)
~  mtc not   match                                                   0..7 gtT 
!  key til   seq           z:re!im  im:!z
,  cat enl                 
^  exc asc                                         parse tree (left to right)
$  str cst                                         (.;`s;a;+;y)      assign       s[a]+:y
#  rsh cnt   take                                  (.;`s;(a;b;c);;y) depth assign s[a;b;c]:y
_  drp flr   drop          ang:_z                  (:;a;b;c)         sequence     a;b;c
?  fnd unq   fnd fnx fnc fni..                     ((/;+);1 2 3)     adverbs      +/1 2 3
@  atx typ                 z:abs@ang  z@ang
.  cal val                 re:. z

+'x  ech(40)  x+'y  ?(168)    x'y  ?
+/x  ovr(123) x+/y  ecr(251)  x/y  ?
+\x  scn(91)  x+\y  ecl(219)  x\y  ?
+':x ecp(41)  x':y  ?(169)    x':y ?
+/:x ?(125)   x+/:y ?(253)    x/:y sv?
+\:x ?(93)    x+\:y ?         x\:y vs?
