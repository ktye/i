sin:F:F{}cos:F:F{}atan2:F:FF{}
ini:I:I{0::289360742959022336j;128::x;p:256;i:8;(i<x)?/((4*i)::p;p*:2;i+:1);kkey::5 mk 0;kval::6 mk 0;x}kkey:{132}kval:{136} /x:16(64k)
bk:I:II{r:32-*7+y*C x;(r<4)? :4;r}mk:I:II{t:x bk y;i:4*t;(~I i)?/i+:4;(128~i)?!;a:I i;i::I a;j:i-4;(j>=4*t)?/(u:a+1<<j>>2;u::I j;j::u;j-:4);a::y|x<<29;(a+4)::1;a}
mki:I:I{r:2 mk 1;(r+8)::x;r}mkf:I:F{r:3 mk 1;(r+8)::x;r}mkd:I:II{v2;ext;r:7 mk 2;(r+8)::x;(r+12)::y;r}mkz:I:II{r:x mkd y;r::2|6<<29;r}l2:I:II{r:6 mk 2;(r+8)::x;(r+12)::y;r}l3:I:III{r:6 mk 3;(r+8)::x;(r+12)::y;(r+16)::z;r}
nn:I:I{536870911&I x}v1:{(x>255)?(xt:(I x)>>29;xn:536870911&I x;xp:8+x)}v2:{v1;yt:(I y)>>29;yn:(I y)&536870911;yp:8+y}r8:{rp:r+8}
fr:V:I{v1;t:4*xt bk xn;x::I t;t::x}dx:V:I{(x>255)?(xr:I x+4;(x+4)::xr-1;(1~xr)?(v1;(xt>4)?xn/(dx I xp+4*i);fr x))}dex:I:II{dx x;y}dxr:{dx x;r}dxyr:{dx x;dx y;r}rx:V:I{x rxn 1}rxn:V:II{(x>255)?(x+:4;x::y+I x)}rl:V:I{v1;xn/(rx I xp;xp+:4)}
lx:I:I{v1;(xt~6)? :x;r:6 mk xn;r8;xn/(rx x;rp::x atx mki i;rp+:4);dxr}lrc:I:II{v1;rl x;r:6 mk xn;r8;xn/(rp::((I.y)(xp));rp+:4;xp+:4);dxr}drc:I:II{rl x;r:(x+8)mkd (I.y)(x+12);dxr}drx:I:II{rx(12+x);dx x;(I.y)(12+x)}
til:I:I{v1;(~2~xt)?!;n:I xp;dx x;(n<'0)? :tir -n;seq(0;n;1)}seq.I:III{r:2 mk y;rp:8+r;y/(rp::z*i+x;rp+:4);r}tir.I:I{r:2 mk x;rp:4+r+4*x;x/(rp::i;rp-:4);r}
ext:{(~xt~yt)?!;((xn~1)&yn>1)?(x:x take yn;xn:yn;xp:x+8);((yn~1)&xn>1)?(y:y take xn;yn:xn;yp:y+8);(~xn~yn)?!}
upx:{(xt>=5)?!;(yt>=5)?!;(xt<yt)?/(x:up(x;xt;xn);xt+:1);(yt<xt)?/(y:up(y;yt;yn);yt+:1);xp:x+8;yp:y+8}
up:I:III{r:(y+1) mk z;xp:x+8;r8;y?[;z/(rp::C xp+i;rp+:4);z/(rp::F?I xp;rp+:8;xp+:4);!];dxr}
atx:I:II{v2;(~xt)?( :x cal enl y);(~yt~2)?!;r:xt mk yn;r8;w:C xt;f:xt+128;yn/(yi:I yp;$[yi<xn;mv(rp;xp+w*yi;w);(V.f)rp];rp+:w;yp+:4);(xt>4)?rl r;(yn~1)?(xt~6)?(rx I r+8;dx r;r:I r+8);dxyr}
cal:I:II{v2;(~yt~6)?!;(yn~1)?(((sadv x)|sadv x-128)? :((I.x)(fst y));(x<128)?x+:128);(x<128)?((~yn~2)?!;rl y;dx y; :(I.x)(I yp;I yp+4));(x<256)?((~yn~1)?!; :(I.x)(fst y));(xn~2)?(rl x;dx x;a:I xp;yn?[; :((I.a+128)(fst y;I xp+4));(rl y;dx y; :((I.a)(I yp;I yp+4;I xp+4)));!]);!;x}
rev:I:I{n:nn x;(~n)? :x;x atx tir n}fst:I:I{v1;(xt~7)?( :x drx 170);x atx mki 0}lst:I:I{v1;(xt~7)?( :x drx 186);x atx mki xn-1}
cut:I:II{v2;(~xt~2)?!;(xn~1)?(r:y drop I xp;dx x; :r);r:6 mk xn;r8;xn/(a:I xp;b:I xp+4;(i~xn-1)?b:yn;(b<a)?!;rx y;rp::y atx seq(a;b-a;1);xp+:4;rp+:4);dxyr}
rsh:I:II{v2;(~xt~2)?!;n:prod(xp;xn);r:y take n;(xn~1)?(dx x; :r);xn-:1;xe:xp+4*xn;xn/(m:I xe;n:n%m;n:xp prod xn-i;r:(seq(0;n;m))cut r;xe-:4);dxr}prod:I:II{r:1;y/(r*:I x;x+:4);r}
take:I:II{xn:nn x;r:til mki y;r8;(xn<y)?(y/(rp+4*i)::i\xn);x atx r}drop:I:II{xn:nn x;(y>xn)?!;x atx seq(y;xn-y;1)}
use:I:I{(1~I x+4)? :x;v1;r:xt mk xn;r8;mv(rp;xp;xn*C xt);dx x;r}mv:V:III{z/(x+i)::C?C y+i}
cat:I:II{v2;(~xt)?(x:enl(x);xt:6);(xt~yt)? :x ucat y;(xt~6)? :x lcat y;!;x}
ucat:I:II{v2;(xt>4)?(rl x;rl y);(xt~7)?(r:((x+8)ucat y+8)mkd(x+12)ucat y+12;dx x;dx y; :r);r:xt mk xn+yn;w:C xt;mv(r+8;xp;w*xn);mv(r+8+w*xn;yp;w*yn);dxyr}
lcat:I:II{v1;((xt bk xn)<(xt bk xn+1))?(r:xt mk xn+1;rl x;dx x;mv(r+8;xp;4*xn);x:r;xp:x+8);(xp+4*xn)::y;x::(xn+1)|6<<29;x}
enl:I:I{(6 mk 0) lcat x} cnt:I:I{dx x;mki nn x}typ:I:I{v1;r:2 mk 1;(8+r)::xt;dxr} not:I:I{x eql mki 0}
wer:I:I{v1;(~xt~2)?!;n:0;xn/(n+:I xp;xp+:4);xp:8+x;r:2 mk n;r8;xn/((I xp)/(rp::i;rp+:4);xp+:4);dxr}
mtc:I:II{r:2 mk 1;(r+8)::x match y;dxyr}match:I:II{(x~y)? :1;(~(I x)~I y)? :0;v1;yp:y+8;m:0;xt?[ :1;n:xn;n:xn<<2;n:xn<<3;n:xn<<4;(xn/((~((I xp) match I yp))? :0;xp+:4;yp+:4); :1)];n/(~(C xp+i)~C yp+i)? :0;1}
fnd:I:II{v2;(~xt~yt)?!;r:2 mk yn;r8;w:C yt;yn/(rp::x fnx yp;rp+:4;yp+:w);dxyr}fnx:I:II{v1;eq:8+xt;w:C xt;xn/(((I.eq)(xp;y))? :i;xp+:w);xn}
exc:I:II{rx x;x atx wer (mki nn y)eql y fnd x}
srt:I:I{rx x;x atx grd x}gdn:I:I{rev grd x}grd:I:I{v1;r:seq(0;xn;1);y:seq(0;xn;1);r8;msrt(y+8;rp;0;xn;xp;xt);dxyr}
msrt:V:IIIIII{((x3-z)>=2)?(c:(x3+z)%2;msrt(y;x;z;c;x4;x5);msrt(y;x;c;x3;x4;x5);mrge(x;y;z;x3;c;x4;x5))}
mrge:V:IIIIIII{k:z;j:x4;w:C x6;i:z;(i<x3)?/(c:k>=x4;(~c)?$[j>=x3;c:0;c:(I.x6)(x5+w*I x+k<<2;x5+w*I x+j<<2)];$[c;(a:j;j+:1);(a:k;k+:1)];(y+i<<2)::I x+a<<2;i+:1)}
gtc:I:II{(C x)>C y}gti:I:II{(I x)>'I y}gtf:I:II{(F x)>F y} eqc:I:II{(C x)~C y}eqi:I:II{(I x)~ I y}eqf:I:II{(F x)~F y}eqz:I:II{(x eqf y)? :(x+8)eqf y+8;0}eqL:I:II{(I x)match I y}
gtl:I:II{x:I x;y:I y;v2;(~xt~yt)? :xt>yt;n:xn;(yn<xn)?n:yn;w:C xt;n/(a:xp+i*w;b:yp+i*w;((I.xt)(a;b))? :1;((I.xt)(b;a))? :0);xn>yn}
eql:I:II{cmp(x;y;1)}mor:I:II{cmp(x;y;0)}les:I:II{cmp(y;x;0)}cmp:I:III{v2;upx;ext;f:xt;z?f+:8;w:C xt;r:2 mk xn;r8;xn/(rp::(I.f)(xp;yp);xp+:w;yp+:w;rp+:4);dxyr}
nd:I:III{v2;upx;ext;w:C xt;f:z+xt;r:xt mk xn;r8;xn/((V.f)(xp;yp;rp);xp+:w;yp+:w;rp+:w);dxyr}
nm:I:II{v1;r:use x;r8;w:C xt;y+:xt;xn/(((V.y)(xp;rp));xp+:w;rp+:w);(xt~4)?(y~19)? :zre r;r}
add:I:II{nd(x;y;143)}sub:I:II{nd(x;y;147)}mul:I:II{nd(x;y;151)}diw:I:II{nd(x;y;155)}
adc:V:III{z::(C x)+C y}adi:V:III{z::(I x)+I y}adf:V:III{z::(F x)+F y}adz:V:III{adf(x;y;z);adf(x+8;y+8;z+8)}
suc:V:III{z::(C x)-C y}sui:V:III{z::(I x)-I y}suf:V:III{z::(F x)-F y}suz:V:III{suf(x;y;z);suf(x+8;y+8;z+8)}
muc:V:III{z::(C x)*C y}mui:V:III{z::(I x)*I y}muf:V:III{z::(F x)*F y}muz:V:III{z::((F x)*F x+8)-(F y)*F y+8;(z+8)::((F x)*F y+8)+(F x+8)*F y}
dic:V:III{z::(C x)%C y}dii:V:III{z::(I x)%I y}dif:V:III{z::(F x)%F y}
diz:V:III{a:F x;b:F x+8;c:F y;d:F y+8;$[(+c)>=(+d);(r:d%c;p:c+r*d;z::(a+b*r)%p;(z+8)::(b-a*r)%p);(r:c%d;p:d+r*c;z::(a*r+b)%p;(z+8)::(b*r-a)%p)]}
abx:I:I{x nm 15}neg:I:I{x nm 19}sqr:I:I{x nm 27}
abc:V:II{c:C x;$[craz c;y::C?c-32;y::C?c]}abi:V:II{i:I x;$[(i<'0);y::0-i;y::i]}abf:V:II{y::+F x}abz:V:II{a:F x;b:F x+8;y::%a*a+b*b}
nec:V:II{c:C x;$[crAZ c;y::C?c+32;y::C?c]}nei:V:II{y::0-I x}nef:V:II{y::-F x}nez:V:II{y::-F x;(y+8)::-F x+8}
sqc:V:II{!}sqi:V:II{!}sqf:V:II{y::%F x}sqz:V:II{y::F x;(y+8)::-F x+8}
zre:I:I{x zri 0}zim:I:I{x zri 8}zri:I:II{v1;(xt~4)?!;r:3 mk xn;r8;xp+:y;xn/(rp::F xp;rp+:8;xp+:16);dxr}
crAZ:I:I{(x>64)?(x<91)? :1;0}craz:I:I{(x>96)?(x<123)? :1;0}
nag:V:I{x::0}nac:V:I{x::C?32}nai:V:I{x::2147483648}naf:V:I{x::9221120237041090561f}naz:V:I{naf x;naf x+8}nas:V:I{x::1 mk 0;(4+I x)::0}nal:V:I{x::0}
drv:I:II{r:0 mk 2;(r+8)::x;(r+12)::y;r}ecv:I:I{40 drv x}epv:I:I{41 drv x}ovv:I:I{123 drv x}scv:I:I{91 drv x}
ech:I:II{x:lx x;v1;r:6 mk xn;r8;rl x;(y<120)?y+:128;xn/(rx y;rp::y atx I xp;xp+:4;rp+:4);dxyr}
ecp:I:II{n:nn x;(~n)?(dx y; :fst x);x rxn -1+2*n;y rxn n-1;r:fst x;(n-1)/r:r cat y cal (x atx mki 1+i)l2 x atx mki i;dxyr}
ovr:I:II{n:nn x;(~n)?(dx y; :fst x);x rxn n;y rxn n-1;r:fst x;(n-1)/r:y cal r l2 x atx mki i+1;dxyr}
scn:I:II{n:nn x;(~n)?(dx y; :fst y);x rxn n;y rxn n-1;t:fst x;rx t;r:enl t;(n-1)/(t:y cal t l2 x atx mki 1+i;rx t;r:r lcat t);dx t;dxyr}
ecr:I:III{n:nn y;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(rp::z cal x l2 y atx mki i;rp+:4);dx z;dxyr}
ecl:I:III{n:nn x;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(rp::z cal (x atx mki i)l2 y;rp+:4);dx z;dxyr}
ecd:I:III{!;x}epi:I:III{!;x}
val:I:I{v1;xt?[;r:(prs x) evl 0;;;;r:x evl 0;(rx x+12;r:x+12;dx x);!];r}
lup:I:II{kv;(m~n)?(y? :x lup 0;dx x; :0);r:I v+8+4*m;rx r;dxr}kv:{k:I kkey;v:I kval;y?(k:I y+8;v:I y+12);n:536870911&I k;m:k fnx 8+x}
asn:I:III{v1;(~xt~5)?!;rx z;kv;(n~m)?(y?!;kkey::k cat x;kval::v cat z; :z);vp:8+v+4*m;dx I vp;vp::z;dx x;z}
swc:I:II{v1;i:1;(i<xn)?/(r:I xp+4*i;rx r;r:r evl y;((~i\2)|(i~xn-1))?(dx x; :r);dx r;i+:1;(~I r+8)?i+:1);dx x;0}
evl:I:II{v1;(~xt~6)?((xt~5)?(xn~1)?(r:x lup y;(~r)?!; :r); :x);(~xn)? :x;(xn~1)? :fst x;v:I xp;(v~36)?(xn>3)?  :x swc y;rl x;r:6 mk xn;r8;xn/(rp::(I xp)evl y;rp+:4;xp+:4);r8;dx x;(xn>2)?((v~186)? :lst r;(v~174)?(v-:128;y:0);(v~46)?(xn~5)?(rl r;dx r;s:I rp+4;a:I rp+8;f:I rp+12;u:I rp+16;(~a)?(~f)? :asn(s;y;u);!));(2~xn)?(rl r;dx r; :(I rp)atx I rp+4);(~3~xn)?!;rx I rp;(I rp)cal r drop 1}
xxx:I:I{!;x}str:I:I{!;x}grp:I:I{!;x}unq:I:I{!;x}flr:I:I{!;x}cst:I:II{!;x}min:I:II{!;x}max:I:II{!;x}
sadv:I:I{$[x~39;r:1;x~47;r:1;x~92;r:1;r:0];r}

/parse
prs:I:I{v1;(~xt~1)?!;8::xp;r:sq xp+xn;$[1~nn r;r:fst r;r:186 cat r];dxr}
sq:I:I{r:enl(pt x)ex x;1/(v:ws x;p:I 8;(~v)?v:C p!59;v?((p<x)?8::1+p; :r);8::1+p;r:r lcat(pt x)ex x);!;x}
ex:I:II{((~x)+(ws y)+((C I 8)is 32))? :x;r:pt y;(isv r)?(~isv x)? :l3(r;x;(pt y)ex y);l2(x;r ex y)}
pt:I:I{r:tok x;(~r)?(p:I 8;(p~x)? :0;l:123~C p;(l|40~C p)?(8::1+p;r:sq x;n:nn r;(n~1)?r:fst r;(n>1)?r:enl r;l?r:lam(p;I 8;r)));1/(p:I 8;b:C p;(p~x)? :r;$[b is 16;r:(tok x)l2 r;b~91;(8::1+p;r:(enl r)cat sq x); :r]);!;r}
isv:I:I{v1;(~xt)? :1;(xt~6)?(xn~2)?(a:I xp;(a<256)?((a is 16)|((a-128)is 16))? :1);0}
lam:I:III{!;x}                                                                       /todo
ws:I:I{(x~I 8)? :1;0}                                                                /todo space
tok:I:I{(ws x)? :0;p:I 8;b:C p;(b is 32)? :0; 2/(r:((I.i+136)(b;p;x));r? :r);0}      /todo 2/->5/
num:I:III{(~x is 4)? :0;((C 1+y) is 4)?!;8::1+y;mki x-48}
vrb:I:III{(~x is 24)? :0;r:C y;(z>1+x)?(58~C 1+y)?(y+:1;r+:128);8::1+y;r}            /todo nam sym chr
is:I:II{y&cla x}cla:I:I{(128<x-32)? :0;C 128+x}
160!204840484848485040604848484848504444444444444444444448604848484848424242424242424242424242424242424242424242424242424240506048484041414141414141414141414141414141414141414141414141414048604800

000:{xxx; gtc; gti; gtf; gtl; gtl; xxx; xxx; xxx; eqc; eqi; eqf; eqz; eqL; eqL; xxx; abc; abi; abf; abz; nec; nei; nef; nez; xxx; xxx; xxx; xxx; sqc; sqi; sqf; sqz}
032:{xxx; mkd; xxx; rsh; cst; diw; min; ecv; ecd; epi; mul; add; cat; sub; cal; ovv; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; dex; xxx; les; eql; mor; fnd}
064:{atx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ecl; scv; xxx; exc; cut}
096:{xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ecr; max; xxx; mtc; xxx}
128:{nag; nac; nai; naf; naz; nas; nal; xxx; num; vrb; xxx; xxx; xxx; xxx; xxx; xxx; adc; adi; adf; adz; suc; sui; suf; suz; muc; mui; muf; muz; dic; dii; dif; diz}
160:{xxx; til; xxx; cnt; str; sqr; wer; epv; ech; ecp; fst; abx; enl; neg; val; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; lst; xxx; grd; eql; gdn; unq}
192:{typ; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; scn; xxx; xxx; srt; flr}
224:{xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ovr; rev; xxx; not; xxx}

\
odo5:{p#’x(&#)’_(p:*/x)%*\x}
fill:{@[y;&y=*0#y;:;x]} /aaron
flip:,'/,''

01234567   xt:x>>29       xn:x&536870911 (-1+1<<29)
Fcifzsld   xt~0(function) x<256(basic) x<128(dyadic)
0148x444   xn~2(derived)  adv  verb
	   xn~3(proj)     verb argv empty-index
	   xn~4(lambda)   str  tree args locals
	   
+  add abx                 abs:+z                   memory
-  sub neg                                          0..  7   type sizes   0 1 4 8 16 4 4 0
*  mul fst                                          8.. 11   parse cur (pp)
%  div sqr                                         12.. 15   
&  min wer                                         16..127   free pointers (4*i) for bt i, i:4..31
|  max rev                                        128..131   memsize log2
<  les grd                                        132..135   k-tree keys
>  mor gdn                                        136..139   k-tree values
=  eql grp                                        140..143   
~  mtc not   match                                144..147
!  key til   seq           z:re!im  im:!z         148..151
,  cat enl                                        152..155   
^  exc asc                                        156..159   
$  str cst                                        160..255   char map az|AZ|NM|VB|AD|TE
#  rsh cnt   take                                 256.. ∞    buckets/heap
_  drp flr   drop          ang:_z                 
?  fnd unq   fnd fnx fnc fni..                     
@  atx typ                 z:abs@ang  z@ang        
.  cal val                 re:. z

+'x  ech(40)  x+'y  ?(168)    x'y  ?              parse tree (left to right) xxx
+/x  ovr(123) x+/y  ecr(251)  x/y  ?              (.;`s;a;+;y)      assign       s[a]+:y
+\x  scn(91)  x+\y  ecl(219)  x\y  ?              (.;`s;(a;b;c);;y) depth assign s[a;b;c]:y
+':x ecp(41)  x':y  ?(169)    x':y ?              (::;a;b;c)        sequence     a;b;c  ::x(last) :[x;y](dex)
+/:x ?(125)   x+/:y ?(253)    x/:y sv?            ((/;+);1 2 3)     adverbs      +/1 2 3                      
+\:x ?(93)    x+\:y ?         x\:y vs?
