/sin:F:F{}cos:F:F{}atan2:F:FF{}
trap:V:II{140::x;144::y}
ini:I:I{0::289360742959022340j;128::x;p:256;i:8;(i<x)?/((4*i)::p;p*:2;i+:1);kkey::5 mk 0;kval::6 mk 0;xyz::((mks 120)cat mks 121)cat mks 122;x}kkey:{132}kval:{136}xyz:{148} /x:16(64k)
bk:I:II{r:32-*7+y*C x;(r<4)? :4;r}mk:I:II{t:x bk y;i:4*t;(~I i)?/i+:4;(128~i)?!;a:I i;i::I a;j:i-4;(j>=4*t)?/(u:a+1<<j>>2;u::I j;j::u;j-:4);a::y|x<<29;(a+4)::1;a}
mki:I:I{r:2 mk 1;(r+8)::x;r}mkf:I:F{r:3 mk 1;(r+8)::x;r}mkd:I:II{v2;(~xn~yn)?!;(~xt~5)?!;(~yt~6)?y:lx y;r:x l2 y;r::2|7<<29;r}mkc:I:I{r:1 mk 1;(r+8)::C?x;r}mks:I:I{sc mkc x}mkz:I:II{r:x mkd y;r::2|6<<29;r}l2:I:II{r:6 mk 2;(r+8)::x;(r+12)::y;r}l3:I:III{r:6 mk 3;(r+8)::x;(r+12)::y;(r+16)::z;r}
nn:I:I{(x<256)? :1;536870911&I x}tp:I:I{(x<256)? :0;(I x)>>29}v1:{xt:tp x;xn:nn x;xp:8+x}v2:{v1;yt:tp y;yn:nn y;yp:8+y}r8:{rp:r+8}
ary:I:I{(x<128)? :2;(x<256)? :1;n:nn x;(n~2)? :1;(~n~4)?!;nn I x+16}
fr:V:I{v1;t:4*xt bk xn;x::I t;t::x}dx:V:I{(x>255)?(xr:I x+4;(x+4)::xr-1;(1~xr)?(v1;((~xt)+xt>4)?xn/(dx I xp+4*i);fr x))}dex:I:II{dx x;y}dxr:{dx x;r}dxyr:{dx x;dx y;r}rx:V:I{x rxn 1}rxn:V:II{(x>255)?(x+:4;x::y+I x)}rl:V:I{v1;xn/(rx I xp;xp+:4)}rld:V:I{rl x;dx x}
lx:I:I{v1;(xt~6)? :x;r:6 mk xn;r8;x rxn xn;xn/(rp::x atx mki i;rp+:4);dxr}lrc:I:II{v1;rl x;r:6 mk xn;r8;xn/(rp::((I.y)(xp));rp+:4;xp+:4);dxr}drc:I:II{rl x;r:(x+8)mkd (I.y)(x+12);dxr}drx:I:II{rx(12+x);dx x;(I.y)(12+x)}
til:I:I{v1;(~2~xt)?!;n:I xp;dx x;(n<'0)? :tir -n;seq(0;n;1)}seq.I:III{r:2 mk y;rp:8+r;y/(rp::z*i+x;rp+:4);r}tir.I:I{r:2 mk x;rp:4+r+4*x;x/(rp::i;rp-:4);r}
ext:{(~xn~yn)?(((xn~1)&yn>1)?(x:x take yn;xn:yn;xp:x+8);((yn~1)&xn>1)?(y:y take xn;yn:xn;yp:y+8))}

upx:I:II{t:tp x;yt:tp y;(t~yt)? :x;((t~7)+(yt~7))?!;(yt~6)? :lx x;n:nn x;(t<yt)?/(x:up(x;t;n);t+:1);x}
up:I:III{r:(y+1) mk z;xp:x+8;r8;y?[;z/(rp::C xp+i;rp+:4);z/(rp::F?'I xp;rp+:8;xp+:4);!];dxr}upxy:{x:x upx y;y:y upx x}

atx:I:II{v2;(~xt)?( :x cal enl y);(xt~7)? :atd(x;y;yt);(yt>5)? :ecr(x;y;64);(~yt~2)?!;r:xt mk yn;r8;w:C xt;yn/(yi:I yp;(xn<=yi)?!;mv(rp;xp+w*yi;w);rp+:w;yp+:4);(xt>4)?rl r;(yn~1)?(xt~6)?(rx I r+8;dx r;r:I r+8);dxyr}
atm:I:II{(~nn y)?(dx x; :y);rx y;f:fst y;t:y drop 1;(6~tp t)?(1~nn t)?t:fst t;nf:nn f;((~f)+~nf~1)?(f?x:x atx f; :ecl(x;t;64));(x atx f)atx t}
atd:I:III{k:I x+8;v:I x+12;(z~5)?(rx k;y:k fnd y;z:2);rx v;dx x;v atx y}

cal:I:II{v2;xt? :x atm y;(~yt~6)?!;(yn~1)?(((sadv x)|sadv x-128)? :((I.x)(fst y));(x<128)?x+:128);(x<128)?((~yn~2)?!;rld y; :(I.x)(I yp;I yp+4));(x<256)?((~yn~1)?!; :(I.x)(fst y));(xn~2)?(rld x;a:I xp;yn?[; :((I.a+128)(fst y;I xp+4));(rld y; :((I.a)(I yp;I yp+4;I xp+4)));!]);(xn~4)? :x lcl y;!;x}
lcl:I:II{fn:I x+20;(~fn~nn y)?!;a:I x+16;rx a;t:I x+12;rx t;an:nn I x+16;(fn<an)?y:y lcat(enl 0)take an-fn;d:a mkd y;r:lst t lev d;dx x;dx d;r}

rev:I:I{n:nn x;(~n)? :x;x atx tir n}fst:I:I{v1;(xt~7)?( :x drx 170);x atx mki 0}lst:I:I{v1;(xt~7)?( :x drx 186);x atx mki xn-1}
cut:I:II{v2;(~xt~2)?!;(xn~1)?(r:y drop I xp;dx x; :r);r:6 mk xn;r8;xn/(a:I xp;b:I xp+4;(i~xn-1)?b:yn;(b<a)?!;rx y;rp::y atx seq(a;b-a;1);xp+:4;rp+:4);dxyr}
rsh:I:II{v2;(~xt~2)?!;n:prod(xp;xn);r:y take n;(xn~1)?(dx x; :r);xn-:1;xe:xp+4*xn;xn/(m:I xe;n:n%m;n:xp prod xn-i;r:(seq(0;n;m))cut r;xe-:4);dxr}prod:I:II{r:1;y/(r*:I x;x+:4);r}
take:I:II{xn:nn x;r:seq(0;y;1);(xn<y)?(r8;y/(rp::i\xn;rp+:4));x atx r}drop:I:II{v1;(y>xn)?!;(xt~6)?(1~xn-y)?( :enl lst x);x atx seq(y;xn-y;1)}
use:I:I{(1~I x+4)? :x;v1;r:xt mk xn;r8;mv(rp;xp;xn*C xt);dx x;r}mv:V:III{z/(x+i)::C?C y+i}
cat:I:II{v2;(~xt)?(x:enl(x);xt:6);(xt~yt)? :x ucat y;(xt~6)? :x ucat lx y;(yt~6)? :(lx x)ucat y;!;x}
ucat:I:II{v2;(xt>4)?(rl x;rl y);(xt~7)?(r:((x+8)ucat y+8)mkd(x+12)ucat y+12;dx x;dx y; :r);r:xt mk xn+yn;w:C xt;mv(r+8;xp;w*xn);mv(r+8+w*xn;yp;w*yn);dxyr}
lcat:I:II{x:use x;v1;((xt bk xn)<(xt bk xn+1))?(r:xt mk xn+1;rld x;mv(r+8;xp;4*xn);x:r;xp:x+8);(xp+4*xn)::y;x::(xn+1)|6<<29;x}
enl:I:I{r:6 mk 1;(r+8)::x;r} cnt:I:I{r:mki nn x;dxr}typ:I:I{v1;r:2 mk 1;(8+r)::xt;dxr} not:I:I{x eql mki 0}
wer:I:I{v1;(xt~1)? :prs x;(~xt~2)?!;n:0;xn/(n+:I xp;xp+:4);xp:8+x;r:2 mk n;r8;xn/((I xp)/(rp::i;rp+:4);xp+:4);dxr}
mtc:I:II{r:2 mk 1;(r+8)::x match y;dxyr}match:I:II{(x~y)? :1;(~(I x)~I y)? :0;v1;yp:y+8;m:0;xt?[ :1;n:xn;n:xn<<2;n:xn<<3;n:xn<<4;(xn/((~((I xp) match I yp))? :0;xp+:4;yp+:4); :1)];n/(~(C xp+i)~C yp+i)? :0;1}
fnd:I:II{v2;(~xt~yt)?!;r:2 mk yn;r8;w:C yt;yn/(rp::x fnx yp;rp+:4;yp+:w);dxyr}fnx:I:II{v1;eq:8+xt;w:C xt;xn/(((I.eq)(xp;y))? :i;xp+:w);xn}
exc:I:II{rx x;x atx wer (mki nn y)eql y fnd x}
srt:I:I{rx x;x atx grd x}gdn:I:I{rev grd x}grd:I:I{v1;r:seq(0;xn;1);y:seq(0;xn;1);r8;msrt(y+8;rp;0;xn;xp;xt);dxyr}
msrt:V:IIIIII{((x3-z)>=2)?(c:(x3+z)%2;msrt(y;x;z;c;x4;x5);msrt(y;x;c;x3;x4;x5);mrge(x;y;z;x3;c;x4;x5))}
mrge:V:IIIIIII{k:z;j:x4;w:C x6;i:z;(i<x3)?/(c:k>=x4;(~c)?$[j>=x3;c:0;c:(I.x6)(x5+w*I x+k<<2;x5+w*I x+j<<2)];$[c;(a:j;j+:1);(a:k;k+:1)];(y+i<<2)::I x+a<<2;i+:1)}
gtc:I:II{(C x)>C y}gti:I:II{(I x)>'I y}gtf:I:II{(F x)>F y} eqc:I:II{(C x)~C y}eqi:I:II{(I x)~ I y}eqf:I:II{(F x)~F y}eqz:I:II{(x eqf y)? :(x+8)eqf y+8;0}eqL:I:II{(I x)match I y}
gtl:I:II{x:I x;y:I y;v2;(~xt~yt)? :xt>yt;n:xn;(yn<xn)?n:yn;w:C xt;n/(a:xp+i*w;b:yp+i*w;((I.xt)(a;b))? :1;((I.xt)(b;a))? :0);xn>yn}
sc:I:I{r:enl x;r::1|5<<29;r}cs:I:I{r:x+8;rx r;dxr}
eql:I:II{cmp(x;y;1)}mor:I:II{cmp(x;y;0)}les:I:II{cmp(y;x;0)}cmp:I:III{upxy;v2;ext;(xt~6)? :ecd(x;y;62-z);f:xt;z?f+:8;w:C xt;r:2 mk xn;r8;xn/(rp::(I.f)(xp;yp);xp+:w;yp+:w;rp+:4);dxyr}
min:I:II{mia(x;y;38)}max:I:II{mia(x;y;124)}mia:I:III{upxy;v2;ext;(xt~6)? :ecd(x;y;z);rx x;rx y;$[z~38;a:x les y;a:x mor y];a:wer a;rx a;asi(y;a;x atx a)}
nd:I:IIII{upxy;v2;ext;(xt~6)? :ecd(x;y;x3);w:C xt;f:z+xt;r:xt mk xn;r8;xn/((V.f)(xp;yp;rp);xp+:w;yp+:w;rp+:w);dxyr}
nm:I:III{v1;(xt>5)? :ech(x;z);r:use x;r8;w:C xt;y+:xt;xn/(((V.y)(xp;rp));xp+:w;rp+:w);(xt~4)?(y~19)? :zre r;r}
add:I:II{nd(x;y;143;43)}sub:I:II{nd(x;y;147;45)}mul:I:II{nd(x;y;151;42)}diw:I:II{nd(x;y;155;37)}
adc:V:III{z::(C x)+C y}adi:V:III{z::(I x)+I y}adf:V:III{z::(F x)+F y}adz:V:III{adf(x;y;z);adf(x+8;y+8;z+8)}
suc:V:III{z::(C x)-C y}sui:V:III{z::(I x)-I y}suf:V:III{z::(F x)-F y}suz:V:III{suf(x;y;z);suf(x+8;y+8;z+8)}
muc:V:III{z::(C x)*C y}mui:V:III{z::(I x)*I y}muf:V:III{z::(F x)*F y}muz:V:III{z::((F x)*F x+8)-(F y)*F y+8;(z+8)::((F x)*F y+8)+(F x+8)*F y}
dic:V:III{z::(C x)%C y}dii:V:III{z::(I x)%I y}dif:V:III{z::(F x)%F y}
diz:V:III{a:F x;b:F x+8;c:F y;d:F y+8;$[(+c)>=(+d);(r:d%c;p:c+r*d;z::(a+b*r)%p;(z+8)::(b-a*r)%p);(r:c%d;p:d+r*c;z::(a*r+b)%p;(z+8)::(b*r-a)%p)]}
abx:I:I{nm(x;15;171)}neg:I:I{nm(x;19;173)}sqr:I:I{nm(x;27;165)}
abc:V:II{c:C x;$[craz c;y::C?c-32;y::C?c]}abi:V:II{i:I x;$[(i<'0);y::0-i;y::i]}abf:V:II{y::+F x}abz:V:II{a:F x;b:F x+8;y::%a*a+b*b}
nec:V:II{c:C x;$[crAZ c;y::C?c+32;y::C?c]}nei:V:II{y::0-I x}nef:V:II{y::-F x}nez:V:II{y::-F x;(y+8)::-F x+8}
sqc:V:II{!}sqi:V:II{!}sqf:V:II{y::%F x}sqz:V:II{y::F x;(y+8)::-F x+8}
zre:I:I{x zri 0}zim:I:I{x zri 8}zri:I:II{v1;(xt~4)?!;r:3 mk xn;r8;xp+:y;xn/(rp::F xp;rp+:8;xp+:16);dxr}
crAZ:I:I{(x>64)?(x<91)? :1;0}craz:I:I{(x>96)?(x<123)? :1;0}
drv:I:II{r:0 mk 2;(r+8)::x;(r+12)::y;r}ecv:I:I{40 drv x}epv:I:I{41 drv x}ovv:I:I{123 drv x}scv:I:I{91 drv x}
ech:I:II{(tp y)? :y bin x;(7~tp x)?(rld x;k:I x+8;v:I x+12; :k mkd v ech y);x:lx x;v1;r:6 mk xn;r8;rl x;(y<120)?y+:128;xn/(rx y;rp::y atx I xp;xp+:4;rp+:4);dxyr}
ecp:I:II{rx x;p:fst x;epi(p;x;y)}epi:I:III{n:nn y;(~n)?(dx x;dx z; :y);y rxn n;z rxn n;r:6 mk n;r8;n/(yi:y atx mki i;rx yi;rp::z cal yi l2 x;x:yi;rp+:4);dx yi;dx y;dx z;r}
ovr:I:II{ovs(x;y;0)}scn:I:II{ovs(x;y;enl 6 mk 0)}scl:V:II{x?(rx y;xp:x+8;xp::(I xp)lcat y)}
ovs:I:III{(1~ary y)? :fxp(x;y;z);n:nn x;x rxn n;r:fst x;y rxn n-1;z scl r;(n-1)/(r:y cal r l2 x atx mki i+1;z scl r);dx x;dx y;(~z)? :r;dx r;fst(z)}
fxp:I:III{t:x;rx x;1?/(rx x;rx y;r:y atx x;((r match x)+r match t)?(dx x;dx y;dx t;z?r:(fst z)lcat r; :r);z scl x;dx x;x:r);x}
ecr:I:III{(1~ary z)? :whl(x;y;z;0);(7~tp y)?(rld y;k:I y+8;v:I y+12; :k mkd ecr(x;v;z));n:nn y;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(rp::z cal x l2 y atx mki i;rp+:4);dx z;dxyr}
ecl:I:III{(1~ary z)? :whl(x;y;z;enl 6 mk 0);(7~tp x)?(rld x;k:I x+8;v:I x+12; :k mkd ecl(v;y;z));n:nn x;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(rp::z cal (x atx mki i)l2 y;rp+:4);dx z;dxyr}
whl:I:IIII{t:tp x;t?(((~t~2)+~1~nn x)?!;dx x; :nlp(y;z;x3;I x+8));r:y;x3 scl r;n:mki 0;1?/(rx x;rx z;r:z atx r;x3 scl r;rx r;t:x atx r;(t match n)?(dx t;dx n;dx z;dx x;x3?(dx r;r:fst x3); :r);dx t);x}
bin:I:II{v2;(~xt~yt)?!;r:2 mk yn;r8;w:C xt;yn/(rp::ibin(xp;yp;xn;xt);rp+:4;yp+:w);dxyr}
ibin:I:IIII{k:0;j:z-1;w:C x3;1?/((k>'j)? :k-1;h:(k+j)>>1;$[((I.x3)(x+w*h;y));j:h-1;k:h+1]);x}

nlp:I:IIII{(x3<0)?!;r:x;y rxn x3;z scl x;x3/(r:y atx r;z scl r);dx y;z?(dx r;r:fst z);r}
ecd:I:III{n:nn x;(~n~nn y)?!;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(c:mki i;rx c;rp::z cal(x atx c)l2 y atx c;rp+:4);dx z;dxyr}
val:I:I{v1;xt?[;r:(prs x) evl 0;;;;r:x evl 0;(rx x+12;r:x+12;dx x);!];r}
lup:I:II{kv;(m~n)?(y? :x lup 0;dx x; :0);r:I v+8+4*m;rx r;dxr}kv:{k:I kkey;v:I kval;y?(k:I y+8;v:I y+12);n:536870911&I k;m:k fnx 8+x}

asn:I:III{v1;(~xt~5)?!;rx z;kv;(n~m)?(y?!;kkey::k cat x;kval::v lcat z; :z);vp:8+v+4*m;dx I vp;vp::z;dx x;z}
asd:I:II{rld x;v:I x+8;s:I x+12;a:I x+16;u:I x+20;(~v~58)?(rx s;r:s lup y;a?(rx a;r:r atx a);u:v cal r l2 u);a?(rx s;u:asi(s lup y;a;u));rx s;r:asn(s;y;u);dx s;r}

asi:I:III{v2;(xt~7)?(yt~5)?(rld x;k:I xp;v:I xp+4;rx k;y:k fnd y; :k mkd asi(v;y;z));
(yt~6)?((xt~7)?(rld x;k:I xp;v:I xp+4;rx y;f:fst y;$[f;(rx k;f:k fnd f);(f:seq(0;nn k;1))]; :k mkd asi(v;(enl f)cat y drop 1;z)); ((~xt~6)+~yt~6)?!; r:x take xn;r8;rx y;a:fst(y);y:y drop 1;(1~nn y)?y:fst y;(~a)?a:seq(0;xn;1);(~2~tp a)?!;an:nn a;ap:a+8;(an~1)?(dx a;ri:rp+4*I ap;ri::asi(I ri;y;z); :r);(~yn~2)?!;(~6~tp z)?(z:(enl z)take an);(~an~nn z)?!;rxn(y;an-1);rl z;zp:z+8;an/(ri:rp+4*I ap;ri::asi(I ri;y;I zp);ap+:4;zp+:4);dx a;dx z; :r);
(~yt~2)?!;zt:tp z;zn:nn z;zp:8+z;(yn>1)?(zn~1)?(~zn~yn)?((~zn~1)?!;z:z take yn;zn:yn;zp:z+8); (xt<5)?((~zt~xt)?!;r:xt mk xn;r8;w:C xt;mv(rp;xp;w*xn);yn/(k:I yp;mv(rp+w*k;zp;w);yp+:4;zp+:w);dx x;dx y;dx z; :r); (xt~6)?(r:x take xn;r8;z:lx z;zp:z+8;rl z;yn/(k:I yp;t:rp+4*k;dx I t;t::I zp;yp+:4;zp+:4);dx y;dx z; :r);!;x}

swc:I:II{v1;i:1;(i<xn)?/(r:I xp+4*i;rx r;r:r evl y;((~i\2)|(i~xn-1))?(dx x; :r);dx r;i+:1;(~I r+8)?i+:1);dx x;0}

/ras:I:III{v:I x+8;(y~3)?(v<256)?((v~58)+v>128)?(rld x;r:I x+12;rx r;u:I x+16;s:fst r;z?(v~58)?(l:I z;n:nn l;(n~l fnx s+8)?(rx s;z::l cat s));(v>128)?v-:128;a:r drop 1;$[~nn a;(dx a;a:0);a:enl(a)];x:(l3(v;enl s;a))lcat u);x}
ras:I:III{v:I x+8;(y~3)?(v<256)?((v~58)+v>128)?((v>128)?v-:128;r:I x+12;r rxn 2;s:fst r;a:r drop 1;$[nn a;(a:a lev z;an:nn a;(an~1)?a:fst a);(dx a;a:0)];u:I x+16;rx u;dx x; :(l3(v;s;a))lcat u evl z);0}

lev:I:II{v1;(~xt~6)? :x;rl x;r:6 mk xn;r8;xn/(rp::(I xp)evl y;rp+:4;xp+:4);dxr}
evl:I:II{v1;(~xt~6)?((xt~5)?(xn~1)?(r:x lup y;(~r)?!; :r); :x);(~xn)? :x;(xn~1)? :(fst x)lev y;v:I xp;(v~36)?(xn>3)?  :x swc y;r:ras(x;xn;y);r? :asd(r;y);x:x lev y;xn:nn x;xp:x+8;(v~128)?(xn>2)? :lst x;(xn~2)?(rl x;r:(I xp)atx I xp+4;dx x; :r);rx I xp;(I xp)cal x drop 1}
xxx:I:I{!;x}str:I:I{!;x}grp:I:I{!;x}unq:I:I{!;x}flr:I:I{!;x}cst:I:II{!;x}
sadv:I:I{$[x~39;r:1;x~47;r:1;x~92;r:1;r:0];r}

/parse
prs:I:I{v1;(~xt~1)?!;8::xp;r:sq xp+xn;$[1~nn r;r:fst r;r:128 cat r];dxr}
sq:I:I{r:enl(pt x)ex x;1/(v:ws x;p:I 8;(~v)?v:~(C p)~59;v?((p<x)?8::1+p; :r);8::1+p;r:r lcat(pt x)ex x);!;x}
ex:I:II{((~x)+(ws y)+((C I 8)is 32))? :x;r:pt y;(isv r)?(~isv x)? :l3(r;x;(pt y)ex y);l2(x;r ex y)}
pt:I:I{r:tok x;(~r)?(p:I 8;(p~x)? :0;l:123~C p;(l+40~C p)?(8::1+p;r:sq x;$[l;r:lam(p;I 8;r);(n:nn r;(n~1)?r:fst r;(n>1)?r:enl r)]));1/(p:I 8;b:C p;(p~x)? :r;$[b is 16;r:(tok x)l2 r;b~91;(8::1+p;r:(enl r)cat sq x); :r]);!;r}
isv:I:I{v1;(~xt)? :1;(xt~6)?(xn~2)?(a:I xp;(a<256)?((a is 16)|((a-128)is 16))? :1);0}

lac:I:II{v1;(xt~6)?xn/(y:(I xp)lac y;xp+:4);(xt~5)?(xn~1)?(p:I xp;(1~nn p)?(r:(C 8+p)-119;(r>y)?(r<4)? :r));y}

loc:I:II{v1;(~xt~6)? :y;xn/(y:(I xp)loc y;xp+:4);xp:x+8;(xn~3)?(58~I xp)?(r:I xp+4;rx r;s:fst r;n:nn y;(n~y fnx s+8)?(rx s;y:y cat s);dx s);y}
lam:I:III{$[91~C 1+x;(rx z;a:((fst z)drop 1)ovr 44;z:z drop 1);(r:I xyz;rx r;a:r take z lac 0)];v:nn a;a:z loc a;n:y-x;t:1 mk n;mv(t+8;x;n);r:0 mk 4;(r+8)::t;(r+12)::z;(r+16)::a;(r+20)::v;r}

ws:I:I{p:I 8;1/((p~x)?(8::p; :1);b:C p;((b~10)+(b is 64))?(8::p; :0);p+:1);x}
tok:I:I{(ws x)? :0;p:I 8;b:C p;(b is 32)? :0;n:0<((C p-1)is 4);(5-n)/(r:((I.i+n+136)(b;p;x));r? :r);0}
pun:I:III{(~x is 4)? :0;((x is 4)*y<z)?/(r*:10;r+:x-48;y+:1;x:C y);8::y;mki r}
pin:I:III{r:pun(x;y;z);r? :r;(x~45)?(y+:1;(y<z)?(x:C y;8::y;r:pun(x;y;z);(~r)?(8::y-1);r8;rp::-I rp; :r));0}
pfl:I:III{r:pin(x;y;z);(~r)? :r;y:I 8;(46~C y)?(r:up(r;2;1);y+:1;8::y;(y<z)?(x:C y;q:pun(x;y;z);q?(q:up(q;2;1);f:1.0;((I 8)-y)/f*:10.0;r8;rp::(F rp)+(F 8+q)%f;dx q)));r}
num:I:III{pfl(x;y;z)} /todo z
nms:I:III{r:num(x;y;z);(~r)? :r;1/(y:I 8;x:C y;((y+2)>z)? :r;(~x~32)? :r;y+:1;8::y;q:num(C y;y;z);(~q)? :r;r:r upx q;q:q upx r;r:r cat q);r}
vrb:I:III{(~x is 24)? :0;(x~39)?((C y-1)~32)?y+:1;r:C y;(z>1+x)?(58~C 1+y)?(y+:1;r+:128);8::1+y;r}
chr:I:III{(~x~34)? :0;a:1+y;1/(y+:1;(y~z)?!;(34~C y)?(n:y-a;r:1 mk n;mv(r+8;a;n);8::1+y; :r));r}
nam:I:III{(~x is 3)? :0;a:y;1/(y+:1;((y~z)+~(C y)is 7)?(n:y-a;r:1 mk n;mv(r+8;a;n);8::y; :sc r));x}
sym:I:III{(~x~96)? :0;y+:1;x:C y;8::y;(y<z)?(r:nam(x;y;z);r? :r;r:chr(x;y;z);r? :sc r);r:5 mk 1;(r+8)::1 mk 0;r}
sms:I:III{r:sym(x;y;z);(~r)? :r;1/(y:I 8;q:sym(C y;y;z);(~q)? :enl r;r:r cat q);r}
is:I:II{y&cla x}cla:I:I{(128<x-32)? :0;C 128+x}
160!204840484848485040604848484848504444444444444444444448604848484848424242424242424242424242424242424242424242424242424240506048484041414141414141414141414141414141414141414141414141414048604800

000:{xxx; gtc; gti; gtf; gtl; gtl; xxx; xxx; xxx; eqc; eqi; eqf; eqz; eqL; eqL; xxx; abc; abi; abf; abz; nec; nei; nef; nez; xxx; xxx; xxx; xxx; sqc; sqi; sqf; sqz}
032:{xxx; mkd; xxx; rsh; cst; diw; min; ecv; ecd; epi; mul; add; cat; sub; cal; ovv; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; dex; xxx; les; eql; mor; fnd}
064:{atx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ecl; scv; xxx; exc; cut}
096:{xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ecr; max; xxx; mtc; xxx}
128:{xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; nms; vrb; chr; nam; sms; xxx; xxx; xxx; adc; adi; adf; adz; suc; sui; suf; suz; muc; mui; muf; muz; dic; dii; dif; diz}
160:{xxx; til; xxx; cnt; str; sqr; wer; epv; ech; ecp; fst; abx; enl; neg; val; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; lst; xxx; grd; eql; gdn; unq}
192:{typ; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; scn; xxx; xxx; srt; flr}
224:{xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; prs; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ovr; rev; xxx; not; xxx}

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
&  min wer   prs flp                               16..127   free pointers (4*i) for bt i, i:4..31
|  max rev                                        128..131   memsize log2
<  les grd                                        132..135   k-tree keys
>  mor gdn                                        136..139   k-tree values
=  eql grp                                        140..143   trp line
~  mtc not   match                                144..147   trp col
!  mkd til   seq           z:re!im  re:!z         148..151   `x`y`z
,  cat enl                                        152..155   
^  exc asc                                        156..159   
$  str cst   sc cs                                160..255   char map az|AZ|NM|VB|AD|TE
#  rsh cnt   take                                 256.. ∞    buckets/heap
_  drp flr   drop          ang:_z                 
?  fnd unq   fnd fnx fnc fni..                     
@  atx typ                 z:abs@ang  z@ang        
.  cal val                 im:. z

+'x  ech(168)      x+'y  ecd(40)        x'y  bin?            parse tree (left to right)
+/x  ovr(251),fxp  x+/y  ecr(123),n/whl x/y  mod,mmul(L)?    (:;`x;y)           assign       x:y
+\x  scn(219),fxp  x+\y  ecl(91),n/whl  x\y  y%x,solve(L)?   (+;(`x;a;b;c);y)   assign(m/i)  s[a;b;c]+:y
+':x ecp(169)      x+':y epi(41)        x':y win?            (;a;b;c)   (*128)  sequence     a;b;c  ::x(last) :[x;y](dex)
+/:x ?(253)        x+/:y ?(125)         x/:y join?           ((/;+);1 2 3)      adverbs      +/1 2 3                      
+\:x ?(221)        x+\:y ?(93)          x\:y split?       
