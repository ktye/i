ini:I:I{0::289360691419414784j;128::x;p:256;i:8;(i<x)?/((4*i)::p;p*:2;i+:1);x} /x:16(64k)
bk:I:II{r:32-*7+y*I?C x;(r<4)? :4;r}mk:I:II{t:x bk y;i:4*t;(~I i)?/i+:4;(128~i)?!;a:I i;i::I a;j:i-4;(j>=4*t)?/(u:a+1<<j>>2;u::I j;j::u;j-:4);a::y|x<<29;(a+4)::1;a}
mki:I:I{r:2 mk 1;(r+8)::x;r}mkf:I:F{r:3 mk 1;(r+8)::x;r}mkd:I:II{v2;ext;r:7 mk 2;(r+8)::x;(r+12)::y;r}mkz:I:II{r:x mkd y;r::2|6<<29;r}
v1:{xt:(I x)>>29;xn:(I x)&536870911;xp:8+x}v2:{v1;yt:(I y)>>29;yn:(I y)&536870911;yp:8+y}
fr:0:I{v1;t:4*xt bk xn;x::I t;t::x}dx:0:I{(x>255)?(xr:I x+4;(x+4)::xr-1;(1~xr)?(v1;(xt>3)?xn/(dx I xp+4*i);fr x))}dxr:{dx x;r}dxyr:{dx x;dx y;r}rx:0:I{(x>255)?(x+:4;x::1+I x)}rl:0:I{v1;xn/(rx I xp;xp+:4)}
til:I:I{v1;(~2~xt)?!;n:I xp;dx x;(n<'0)? :tir -n;seq(0;n;1)}seq:I:III{r:2 mk y;rp:8+r;y/(rp::z*i+x;rp+:4);r}tir:I:I{r:2 mk x;rp:4+r+4*x;x/(rp::i;rp-:4);r}
ext:{(~xt~yt)?!;((xn~1)&yn>1)?(x:x take yn;xn:yn;xp:x+8);((yn~1)&xn>1)?(y:y take xn;yn:xn;yp:y+8);(~xn~yn)?!}
upx:{(xt>=5)?!;(yt>=5)?!;(xt<yt)?/(x:up(x;xt;xn);xt+:1);(yt<xt)?/(y:up(y;yt;yn);yt+:1);xp:x+8;yp:y+8}
up:I:III{r:(y+1) mk z;xp:x+8;rp:r+8;y?[;z/(rp::I?C xp+i;rp+:4);z/(rp::F?I xp;rp+:8;xp+:4);!];dxr}
atx:I:II{v2;(~yt~2)?!;r:xt mk yn;rp:r+8;xt?[!;atc;atI;atF;atL;atL;!];dxyr} atc:{(yn/((rp+i)::C?32;yi:I yp;(yi<xn)?(rp+i)::C xp+yi;yp+:4))}
atL:{(nas:1 mk 0;yn/(rp::nas;yi:I yp;(yi<xn)?rp::I xp+4*yi;rp+:4;yp+:4);rl r;dx nas)}atT:{(yn/(rp::naT;yi:I yp;(yi<xn)?rp::T xp+W*yi;rp+:W;yp+:4))}naI:{-2147483648}naF:{9221120237041090561f}
rev:I:I{v1;(~xn)? :x;x atx tir xn}fst:I:I{v1;(xt~7)?(rx 12+x;dx x; :fst 12+x);x atx mki 0}
cut:I:II{v2;(~xt~2)?!;(xn~1)?(r:y drop I xp;dx x; :r);r:5 mk xn;rp:r+8;xn/(a:I xp;b:I xp+4;(i~xn-1)?b:yn;(b<a)?!;rx y;rp::y atx seq(a;b-a;1);xp+:4;rp+:4);dxyr}
rsh:I:II{v2;(~xt~2)?!;n:prod(xp;xn);r:y take n;(xn~1)?(dx x; :r);xn-:1;xe:xp+4*xn;xn/(m:I xe;n:n%m;n:xp prod xn-i;r:(seq(0;n;m))cut r;xe-:4);dxr}prod:I:II{r:1;y/(r*:I x;x+:4);r}
take:I:II{v1;r:til mki y;rp:r+8;(xn<y)?(y/(rp+4*i)::i\xn);x atx r}drop:I:II{v1;(y>xn)?!;x atx seq(y;xn-y;1)}
use:{(~1~I v+4)?(r:xt mk xn;rp+r+8;mv(rp;xp;xn*I?C xt);dx x;x:r;xp:x+8)}mv:0:III{z/((x+i)::C y+i)}
cat:I:II{v2;((~xt)|~yt)?!;(xt~yt)? :x ucat y;(xt~5)? :x lcat y;!;x}
ucat:I:II{v2;(xt>4)?(rl x);(xt>5)?(r:(x+8)mkd x+12;dx x;dx y; :r);r:xt mk xn+yn;w:I?C xt;mv(r+8;xp;w*xn);mv(r+8+w*xn;yp;w*yn);dxyr}
lcat:I:II{v1;((xt bk xn)<(xt bk xn+1))?(r:xt mk xn+1;mv(r+8;xp;4*xn);dx x;x:r;xp:x+8);(xp+4*xn)::y;x::(xn+1)|5<<29;x}
enl:I:I{(5 mk 0) lcat x} cnt:I:I{v1;dx x;mki xn}tip:I:I{v1;r:2 mk 1;(8+r)::xt;dxr} not:I:I{x eql mki 0}
wer:I:I{v1;(~xt~2)?!;n:0;xn/(n+:I xp;xp+:4);xp:8+x;r:2 mk n;rp:r+8;xn/((I xp)/(rp::i;rp+:4);xp+:4);dxr}
mtc:I:II{r:2 mk 1;(r+8)::x match y;dxyr}match:I:II{(x~y)? :1;(~(I x)~I y)? :0;v1;yp:y+8;m:0;xt?[ :1;nn:xn;nn:xn<<2;nn:xn<<3;(xn/((~((I xp) match I yp))? :0;xp+:4;yp+:4); :1)];nn/(~(C xp+i)~C yp+i)? :0;1}
fnd:I:II{v2;(~xt~yt)?!;r:2 mk yn;rp:r+8;w:I?C yt;yn/(rp::x fnx yp;rp+:4;yp+:w);dxyr}fnx:I:II{v1;xt?[!; :fnc(xp;xn;I?C y); :fni(xp;xn;I y); :fnj(xp;xn;J y);;;!;!; :fnl(xp;xn;I y)];x}fnc:I:III{y/((C?z)~C x+i)? :i;y}fni:I:III{y/((z~I x)? :i;x+:4);y}fnj:I:IIJ{y/((z~J x)? :i;x+:8);y}fnl:I:III{y/((z match I x)? :i;x+:4);y}
exc:I:II{r:2 mk 1;(r+8)::(I y)&536870911;rx x;x atx wer r eql y fnd x}
srt:I:I{rx x;x atx grd x}gdn:I:I{rev grd x}
grd:I:I{v1;r:seq(0;xn;1);y:seq(0;xn;1);rp:r+8;msrt(y+8;rp;0;xn;xp;xt);dxyr}
msrt:0:IIIIII{((x3-z)>=2)?(c:(x3+z)%2;msrt(y;x;z;c;x4;x5);msrt(y;x;c;x3;x4;x5);mrge(x;y;z;x3;c;x4;x5))}
mrge:0:IIIIIII{k:z;j:x4;w:I?C x6;i:z;(i<x3)?/(c:k>=x4;(~c)?$[j>=x3;c:0;c:(I.x6)(x5+w*I x+k<<2;x5+w*I x+j<<2)];$[c;(a:j;j+:1);(a:k;k+:1)];(y+i<<2)::I x+a<<2;i+:1)}
gtc:I:II{(C x)>C y}gti:I:II{(I x)>'I y}gtf:I:II{(F x)>F y} eqc:I:II{(C x)~C y}eqi:I:II{(I x)~ I y}eqf:I:II{(F x)~F y}
gtl:I:II{x:I x;y:I y;v2;(~xt~yt)? :xt>yt;n:xn;(yn<xn)?n:yn;w:I?C xt;n/(a:xp+i*w;b:yp+i*w;((I.xt)(a;b))? :1;((I.xt)(b;a))? :0);xn>yn}
eql:I:II{cmp(x;y;1)}mor:I:II{cmp(x;y;0)}les:I:II{cmp(y;x;0)}cmp:I:III{v2;upx;ext;f:xt;z?f+:8;w:I?C xt;r:2 mk xn;rp:r+8;xn/(rp::(I.f)(xp;yp);xp+:w;yp+:w;rp+:4);dxyr}
1:{gtc;gti;gtf;gtl;gtl} 9:{eqc;eqi;eqf;match;match}

\
/not:I:I{v1;xt?[;(notC);;(notf);; :x lrc 126;(notZ); :x drc 126;(notI)];x}
not:I:I{v1;xt?[;!;;!;;!;!;!; :x];x}
notT:{r:2 mk xn;rp+r+8;xn/(rp::~T xp;rp+:4;xp+:W);dx x; :r}
/notf:{r:2 mk xn;rp+r+8;xn/(rp::0.~F xp;rp+:4;xp+:8);dx x; :r}
notZ:{rl x;dx x; :(not x+8)max not x+12}
/lrc:I:II{v1;rl x;r:5 mk xn;rp:r+8;xn/(rp::.y xp;rp+:4;xp+:4);dxr} // todo call indirect .y
/drc:I:II{v1;rl x;r:(x+8)mkd.x+12;dxr}
up:I:II{v1; (xt<y)?/xt?[;ic;fi; :x mkz mkf 0.;!];dxr} ic:{(r:2 mk xn;rp:8+r;xn/(rp::I?C xp+i;rp+:4);dx x;x:r;xt:2)} fi:{(r:3 mk xn;rp:8+r;xn/(rp::F?I xp;rp+:8;xp+:4);dx x;x:r;xt:3)}
/eql:I:II{v2;ext;xt?[;eqC;eqI;eqF;;!;!;!;eqI];dxyr} eqT:{(r:2 mk xn;rp:r+8;xn/(rp::(T xp)~T yp;rp+:4;xp+:W;yp+:W))}


\
cnt:I:I{v1;decr x;mki xn}
tip:I:I{v1;r:2 mk 2;(8+r)::xt;(12+r)::xn;dxr}
sumi:I:II{y/r+:I x+1;r}
wer:I:I{v1;(~2~xt)?!;xn/(0>'I xp+i)?!;r:2 mk xn sumi 8+xp;rp:8+r;xn/(I xp)/(rp::i;rp+:4);dxr}
enl:I:I{v1;r:1 mk 6;(8+r)::x;r}
neg:I:I{v1;xt?[!;!;;;;!; :45 lrc x; :45 drc x;r:xt mk xn];decr x;(2~xt)?(rp:r+8;xn/(rp::0-I xp;rp+:4;xp+:4); :r);(4~xt)?(xt:3;xn*:2);rp:8+r;xn/(rp::-F xp;rp+:4;xp+:4);r}
unq:I:I{v1;xt?[!;;;;;;;!;r:xt mk 0];rn:0;xn/(~rn~r fnd1 xi:x ati i)?(rn+:1;r:r cat xi);(xt~6)?lnc r;dxr}
ati:I:I{v1;r:xt mk 1;rp:8+r;xt?[!;rp::C xp;;rp::J xp;(rp::J xp;(8+rp):J xp+8);;!;rp::I xp];(6~xt)?(rp::inc rp);r}
fnd:I:II{v2;r:2 mk yn];rp:r+8;yn/(yi:y ati i;rp::x fnd1 yi);dxyr}fnd1:I:II{v2;decr y;r:xn; xn/(xi:x ati i)mc y}
mc:I:II{(x~y)? :0; xt?[ :1; (xn/((~(C xp)~(C yp) :0);xp+:1;yp+:1); :1); ;xn*:2;xn*:4; ; :x mtchl y; :(xp mtl yp)&&(4+xp)mtl 4+yp) ]; xn/((~(I xp)~(I yp) :0);xp+:4;yp+:4); :1}
mtc:I:II{v2;r:2 mk 1;rp::x mc y;dxyr}mtl:I:II{v2;(xn/(xi:x ati;yi y ati i;r:(xi mc yi);(~r)?dxyr);dxyr)}

\
odo5:{p#’x(&#)’_(p:*/x)%*\x}

r8:{rp:r+8}
add:I:II{v2; ..conform: r:xn mk xt; r8; xt?[adddC;adddI;adddF;adddZ!]; dxyr }
addd:T{ nd1+nd2 }subbb:T{ nd1-nd2 }mulll:T{ nd1*nd2 }
nd1:{vx:xn>1;vy:yn>1; xn/(rp,i)::(T xp,i*vx) }nd2:{ T yp,i*vy; dxyr }  /vx vy: 0(atom) 1(vector)

p,i → (addr,offset,width)
e.g. (p,i)::I 144,0

01234567   xt:x>>29       xn:x&536870911 (-1+1<<29)
Fcifslzd   xt~0(function) x<256(basic) 
01484444   ft~xn&0xff00  (derived, proj, lambda, native)  composition==lambda?
	   fn~xn&0xff    (argn)

+ add abs                 abs:+z                   memory
- sub neg                                          0..  7   type sizes   0 1 4 8 16 4 4 0
* mul fst                                          8.. 11   k-tree/key   pointer     todo
% div sqr                                         12.. 15   k-tree/value pointer     todo
& min wer                                         16..127   free pointers (4*i) for bt i, i:4..31
| max rev                                        128..131   memsize log2
< les grd                              
> mor gdn                                        function tables (call indirect)
= eql grp                                        T1:I:I (#128)      T2:I:I (#128)
~ mtc not   match                                                   0..7 gtT 
! key til   seq           z:re!im  im:!z
, cat enl                 
^ exc asc                 
$ str cst   
# rsh cnt   take
_ drp flr   drop          ang:_z
? fnd unq   fnd fnx fnc fni..              
@ atx typ                 z:abs@ang  z@ang
. cal val                 re:. z
