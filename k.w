sin:F:F{}cos:F:F{}log:F:F{}atan2:F:FF{}hypot:F:FF{}draw:V:III{}grow:I:I{}printc:V:II{}
trap:V:II{140::x;144::y}
ini:I:I{0::289360742959022340j;12::1887966018;128::x;p:256;i:8;(i<x)?/((4*i)::p;p*:2;i+:1);kkey::5 mk 0;kval::6 mk 0;xyz::((mks 120)cat mks 121)cat mks 122;x}kkey:{132}kval:{136}xyz:{148}
bk:I:II{r:32-*7+y*C x;(r<4)? :4;r}
mk:I:II{t:x bk y;i:4*t;m:4*I 128;(~I i)?/((i>='m)?(m:grow 1+i%4;128::m;i::1<<i>>2;m:i;i-:4);i+:4);(128~i)?!;a:I i;i::I a;j:i-4;(j>=4*t)?/(u:a+1<<j>>2;u::I j;j::u;j-:4);a::y|x<<29;(a+4)::1;a}
mki:I:I{r:2 mk 1;(r+8)::x;r}mkf:I:F{r:3 mk 1;(r+8)::x;r}
mkd:I:II{v2;ext;(xt~3)?(yt~3)?(r:4 mk xn;r8;xn/(rp::F xp;(rp+8)::F yp;rp+:16;xp+:8;yp+:8);dx x;dx y; :r);(~xt~5)?!;(~yt~6)?y:lx y;r:x l2 y;r::2|7<<29;r}
mkc:I:I{r:1 mk 1;(r+8)::C?x;r}mks:I:I{sc mkc x}mkz:I:II{r:x mkd y;r::2|6<<29;r}l2:I:II{r:6 mk 2;(r+8)::x;(r+12)::y;r}l3:I:III{r:6 mk 3;(r+8)::x;(r+12)::y;(r+16)::z;r}
nn:I:I{(x<256)? :1;536870911&I x}tp:I:I{(x<256)? :0;(I x)>>29}v1:{xt:tp x;xn:nn x;xp:8+x}v2:{v1;yt:tp y;yn:nn y;yp:8+y}r8:{rp:r+8}
ary:I:I{(x<128)? :2;(x<256)? :1;n:nn x;(n~2)? :1;I x+20}
fr:V:I{v1;t:4*xt bk xn;x::I t;t::x}dx:V:I{(x>255)?(xr:I x+4;(x+4)::xr-1;(1~xr)?(v1;((~xt)+xt>4)?xn/(dx I xp+4*i);fr x))}dex:I:II{dx x;y}dxr:{dx x;r}dxyr:{dx x;dx y;r}rx:V:I{x rxn 1}rxn:V:II{(x>255)?(x+:4;x::y+I x)}rl:V:I{v1;xn/(rx I xp;xp+:4)}rld:V:I{rl x;dx x}
lx:I:I{v1;(xt~6)? :x;((xt~7)+xn~1)? :enl x;r:6 mk xn;r8;x rxn xn;xn/(rp::x atx mki i;rp+:4);dxr}
til:I:I{v1;(~xt)?(xn~4)? :lcl(x;0;1);(4~xt)? :zim x;(7~xt)?(r:I xp;rx r;dx x; :r);(~2~xt)?!;n:I xp;dx x;(n<'0)? :tir -n;seq(0;n;1)}seq.I:III{r:2 mk y;rp:8+r;y/(rp::z*i+x;rp+:4);r}tir.I:I{r:2 mk x;rp:4+r+4*x;x/(rp::i;rp-:4);r}
ext:{(~xn~yn)?((xn~1)?(x:x take yn;xn:yn;xp:x+8);(yn~1)?(y:y take xn;yn:xn;yp:y+8))}
upx:I:II{t:tp x;yt:tp y;(t~yt)? :x;((t~7)+(yt~7))?!;(yt~6)? :lx x;n:nn x;(t<yt)?/(x:up(x;t;n);t+:1);x}
up:I:III{r:(y+1) mk z;xp:x+8;r8;y?[;z/(rp::C xp+i;rp+:4);z/(rp::F?'I xp;rp+:8;xp+:4);z/(rp::F xp;(rp+8)::0.0;xp+:8;rp+:16);!];dxr}upxy:{x:x upx y;y:y upx x}
atx:I:II{v2;(~xt)?( :x cal enl y);(xt~7)? :atd(x;y;yt);(yt>5)? :ecr(x;y;64);(yt~3)?(xt<5)? :x phi y;(~yt~2)?!;r:xt mk yn;r8;w:C xt;yn/(yi:I yp;(xn<=yi)?!;mv(rp;xp+w*yi;w);rp+:w;yp+:4);(xt>4)?rl r;(yn~1)?(xt~6)?(rx I r+8;dx r;r:I r+8);dxyr}
atm:I:II{(~nn y)?(dx x; :y);rx y;f:fst y;t:y drop 1;(1<nn t)? :(x atx f)atm t;t:fst t;nf:nn f;((~f)+~nf~1)?(f?x:x atx f; :ecl(x;t;64));(x atx f)atx t}
atd:I:III{k:I x+8;v:I x+12;(z~5)?(rx k;y:k fnd y;z:2);rx v;dx x;v atx y}
cal:I:II{v2;xt? :x atm y;(~yt~6)?!;(yn~1)?(((sadv x)|sadv x-128)? :((I.x)(fst y));(x<128)?x+:128);(x<128)?((~yn~2)?!;rld y; :(I.x)(I yp;I yp+4));(x<256)?((~yn~1)?!; :(I.x)(fst y));(xn~2)?(rld x;a:I xp;yn?[; :((I.a+128)(fst y;I xp+4));(rld y; :((I.a)(I yp;I yp+4;I xp+4)));!]);(xn~3)?(rl x;r:asi(I x+12;I x+16;y);v:I x+8;dx x; :v cal r);(xn~4)?(a:I x+20;(a>yn)?(a-:yn;a/y:y lcat 0; :prj(x;y;seq(yn;a;1))); :lcl(x;y;0));!;x}
lcl:I:III{fn:I x+20;(z+~fn)?(dx y;y:6 mk fn);(~fn~nn y)?!;y:y take fn;(1~fn)?y:enl y;a:I x+16;rx a;t:I x+12;rx t;an:nn I x+16;(fn<an)?(an-fn)/y:y lcat 0;d:a mkd y;r:lst t ltr d;dx x;z?(dx r; :d);dx d;r}
rev:I:I{(7~tp x)?(rld x;:(rev I x+8)mkd rev I x+12);n:nn x;(~n)? :x;x atx tir n}fst:I:I{v1;(~xn)?(dx x;(xt~0)? :0;(xt~5)? :sc 1 mk 0;(xt>5)? :6 mk 0; :(mki xt)cst mkc 0);(~xt)? :x;(xt~7)?( :fst val x);x atx mki 0}lst:I:I{(7~tp x)? :lst val x;x atx mki (nn x)-1}
cut:I:II{v2;(yt~7)?((xt~2)?(k:I yp;rx k;x:x cut k; :x tkd y);rx y; :((til y)exc x)tkd y);(~xt~2)?!;(xn~1)?(r:y drop I xp;dx x; :r);r:6 mk xn;r8;xn/(a:I xp;b:I xp+4;(i~xn-1)?b:yn;(b<a)?!;rx y;rp::y atx seq(a;b-a;1);xp+:4;rp+:4);dxyr}
rsh:I:II{v2;(yt~7)?((xt~2)?(k:I yp;rx k;x:x rsh k); :x tkd y);(~xt~2)?!;n:prod(xp;xn);r:y take n;(xn~1)?(dx x; :r);xn-:1;xe:xp+4*xn;xn/(m:I xe;n:n%m;n:xp prod xn-i;r:(seq(0;n;m))cut r;xe-:4);dxr}prod:I:II{r:1;y/(r*:I x;x+:4);r}
take:I:II{xn:nn x;o:0;(y<'0)?(o:xn+y;y:-y;(o<'0)? :x);r:seq(o;y;1);(xn<y)?(r8;y/(rp::i\xn;rp+:4));x atx r}drop:I:II{v1;a:y;(y<'0)?(y:0-y;a:0);(y>xn)?(dx x; :xt mk 0);x:x atx seq(a;xn-y;1);(xt~6)?(1~xn-y)?x:enl x;x}
tkd:I:II{t:tp x;k:I y+8;v:I y+12;rld y;(~t~5)?!;rx k;x:k fnd x;rx x;v:v atx x;(1~nn x)?v:enl v;(k atx x)mkd v}
phi:I:II{n:nn y;r:4 mk n;r8;yp:y+8;n/(p:0.017453292519943295*F yp;rp::cos p;(rp+8)::sin p;rp+:16;yp+:8);dx y;x mul r}
use:I:I{(1~I x+4)? :x;v1;r:xt mk xn;r8;mv(rp;xp;xn*C xt);dx x;r}mv:V:III{z/(x+i)::C?C y+i}
cat:I:II{v2;(~xt)?(x:enl(x);xt:6);(xt~yt)? :x ucat y;(xt~6)? :x ucat lx y;(yt~6)? :(lx x)ucat y;!;x}
ucat:I:II{v2;(xt>4)?(rl x;rl y);(xt~7)?(r:((I x+8)ucat I y+8)mkd(I x+12)ucat I y+12;dx x;dx y; :r);r:xt mk xn+yn;w:C xt;mv(r+8;xp;w*xn);mv(r+8+w*xn;yp;w*yn);dxyr}
lcat:I:II{x:use x;v1;((xt bk xn)<(xt bk xn+1))?(r:xt mk xn+1;rld x;mv(r+8;xp;4*xn);x:r;xp:x+8);(xp+4*xn)::y;x::(xn+1)|6<<29;x}
enl:I:I{r:6 mk 1;(r+8)::x;r} cnt:I:I{(7~tp x)?x:til x;r:mki nn x;dxr}typ:I:I{v1;r:2 mk 1;(8+r)::xt;dxr} not:I:I{x eql mki 0}
wer:I:I{v1;(xt~1)? :prs x;(xt~4)? :zan(x;xn;xp);(xt~6)? :flp x;(~xt~2)?!;n:0;xn/(n+:I xp;xp+:4);xp:8+x;r:2 mk n;r8;xn/((I xp)/(rp::i;rp+:4);xp+:4);dxr}
mtc:I:II{r:2 mk 1;(r+8)::x match y;dxyr}match:I:II{(x~y)? :1;(~(I x)~I y)? :0;v1;yp:y+8;m:0;xt?[ :1;n:xn;n:xn<<2;n:xn<<3;n:xn<<4;(xn/((~((I xp) match I yp))? :0;xp+:4;yp+:4); :1)];n/(~(C xp+i)~C yp+i)? :0;1}
fnd:I:II{v2;(~xt~yt)?!;r:2 mk yn;r8;w:C yt;yn/(rp::x fnx yp;rp+:4;yp+:w);dxyr}fnx:I:II{v1;eq:8+xt;w:C xt;xn/(((I.eq)(xp;y))? :i;xp+:w);xn}
jon:I:II{v1;((~xt~6)+~xn)?(dx y; :x);rl x;r:I xp;y rxn xn-2;(xn-1)/(xp+:4;r:(r cat y)cat I xp);dxr}
spl:I:II{rx x;yn:nn y;r:((mki 0)cat x fds y)cut x;rn:(nn r)-1;r8;rn/(rp+:4;rp::(I rp) drop yn);r}
fds:I:II{v2;((~xt~yt)+xt>5)?!;(xn<yn)?(dx x;dx y; :2 mk 0);(~yn)?(dx x;dx y; :(seq(0;xn;1))drop 1);r:2 mk 0;w:C xt;eq:8+xt;i:0;xn/(a:0;yn/(k:w*j;a+:((I.eq)(xp+k;yp+k)));(a~yn)?(r:r ucat mki i;i+:yn-1;xp+:w*yn-1);xp+:w);dxyr}
exc:I:II{rx x;x atx wer (mki nn y)eql y fnd x}
srt:I:I{rx x;x atx grd x}gdn:I:I{rev grd x}grd:I:I{v1;r:seq(0;xn;1);y:seq(0;xn;1);r8;msrt(y+8;rp;0;xn;xp;xt);dxyr}
msrt:V:IIIIII{((x3-z)>=2)?(c:(x3+z)%2;msrt(y;x;z;c;x4;x5);msrt(y;x;c;x3;x4;x5);mrge(x;y;z;x3;c;x4;x5))}
mrge:V:IIIIIII{k:z;j:x4;w:C x6;i:z;(i<x3)?/(c:k>=x4;(~c)?$[j>=x3;c:0;c:(I.x6)(x5+w*I x+k<<2;x5+w*I x+j<<2)];$[c;(a:j;j+:1);(a:k;k+:1)];(y+i<<2)::I x+a<<2;i+:1)}
gtc:I:II{(C x)>C y}gti:I:II{(I x)>'I y}gtf:I:II{(F x)>F y} eqc:I:II{(C x)~C y}eqi:I:II{(I x)~ I y}eqf:I:II{(J x)~J y}eqz:I:II{(x eqf y)? :(x+8)eqf y+8;0}eqL:I:II{(I x)match I y}
gtl:I:II{x:I x;y:I y;v2;(~xt~yt)? :xt>yt;n:xn;(yn<xn)?n:yn;w:C xt;n/(a:xp+i*w;b:yp+i*w;((I.xt)(a;b))? :1;((I.xt)(b;a))? :0);xn>yn}
sc:I:I{r:enl x;r::1|5<<29;r}cs:I:I{r:I x+8;rx r;dxr}
eql:I:II{cmp(x;y;1)}mor:I:II{cmp(x;y;0)}les:I:II{cmp(y;x;0)}cmp:I:III{upxy;v2;ext;(xt~6)? :ecd(x;y;62-z);f:xt;z?f+:8;w:C xt;r:2 mk xn;r8;xn/(rp::(I.f)(xp;yp);xp+:w;yp+:w;rp+:4);dxyr}
min:I:II{mia(x;y;38)}max:I:II{mia(x;y;124)}mia:I:III{upxy;v2;ext;(xt~6)? :ecd(x;y;z);rx x;rx y;$[z~38;a:x les y;a:x mor y];a:wer a;rx a;asi(y;a;x atx a)}
nd:I:IIII{upxy;v2;ext;(xt~6)? :ecd(x;y;x3);w:C xt;f:z+xt;r:xt mk xn;r8;xn/((V.f)(xp;yp;rp);xp+:w;yp+:w;rp+:w);dxyr}
nm:I:III{v1;(xt>5)? :ech(x;z);r:use x;r8;w:C xt;y+:xt;xn/(((V.y)(xp;rp));xp+:w;rp+:w);(xt~4)?(y~19)? :zre r;r}
add:I:II{nd(x;y;143;43)}sub:I:II{nd(x;y;147;45)}mul:I:II{nd(x;y;151;42)}diw:I:II{nd(x;y;155;37)}mod:I:II{nd(x;y;23;7)}
adc:V:III{z::C?(C x)+C y}adi:V:III{z::(I x)+I y}adf:V:III{z::(F x)+F y}adz:V:III{adf(x;y;z);adf(x+8;y+8;z+8)}
suc:V:III{z::C?(C x)-C y}sui:V:III{z::(I x)-I y}suf:V:III{z::(F x)-F y}suz:V:III{suf(x;y;z);suf(x+8;y+8;z+8)}
muc:V:III{z::C?(C x)*C y}mui:V:III{z::(I x)*I y}muf:V:III{z::(F x)*F y}muz:V:III{z::((F x)*F y)-(F y+8)*F x+8;(z+8)::((F x)*F y+8)+(F x+8)*F y}
dic:V:III{z::C?(C x)%C y}dii:V:III{z::(I x)%I y}dif:V:III{z::(F x)%F y}moi:V:III{x:I x;y:I y;z::(y+I?x\'y)\'y}
diz:V:III{a:F x;b:F x+8;c:F y;d:F y+8;$[(+c)>=(+d);(r:d%c;p:c+r*d;z::(a+b*r)%p;(z+8)::(b-a*r)%p);(r:c%d;p:d+r*c;z::(a*r+b)%p;(z+8)::(b*r-a)%p)]}
abx:I:I{nm(x;15;171)}neg:I:I{nm(x;19;173)}sqr:I:I{nm(x;27;165)}
abc:V:II{c:C x;$[craz c;y::C?c-32;y::C?c]}abi:V:II{i:I x;$[(i<'0);y::0-i;y::i]}abf:V:II{y::+F x}abz:V:II{y::(F x)hypot F x+8}
nec:V:II{c:C x;$[crAZ c;y::C?c+32;y::C?c]}nei:V:II{y::0-I x}nef:V:II{y::-F x}nez:V:II{y::-F x;(y+8)::-F x+8}
sqc:V:II{!}sqi:V:II{!}sqf:V:II{y::%F x}sqz:V:II{y::F x;(y+8)::-F x+8}
lgf:I:I{v1;(~xt~3)?!;x:use x;xp:x+8;xn/(xp::log F xp;xp+:8);x}
zre:I:I{x zri 0}zim:I:I{x zri 8}zri:I:II{v1;r:3 mk xn;r8;xp+:y;xn/(rp::F xp;rp+:8;xp+:16);dxr}zan:I:III{r:3 mk y;r8;y/(rp::(F z)ang F z+8;z+:16;rp+:8);dxr}
crAZ:I:I{(x>64)?(x<91)? :1;0}craz:I:I{(x>96)?(x<123)? :1;0}
drv:I:II{r:0 mk 2;(r+8)::x;(r+12)::y;r}ecv:I:I{40 drv x}epv:I:I{41 drv x}ovv:I:I{123 drv x}riv:I:I{125 drv x}scv:I:I{91 drv x}liv:I:I{93 drv x}
ech:I:II{(tp y)? :y bin x;(7~tp x)?(rld x;k:I x+8;v:I x+12; :k mkd v ech y);x:lx x;v1;r:6 mk xn;r8;rl x;(y<120)?y+:128;xn/(rx y;rp::y atx I xp;xp+:4;rp+:4);dxyr}
ecp:I:II{rx x;p:fst x;epi(p;x;y)}epi:I:III{n:nn y;(~n)?(dx x;dx z; :y);y rxn n;z rxn n;r:6 mk n;r8;n/(yi:y atx mki i;rx yi;rp::z cal yi l2 x;x:yi;rp+:4);dx yi;dx y;dx z;r}
ovr:I:II{t:tp y;(2~t)? :x mod y;(6~t)? :(ecl(y;x;42))ech ovv 43;ovs(x;y;0;0)}scn:I:II{t:tp y;t?(t<5)? :x diw y;(t>5)? :(I.81)(y;x);ovs(x;y;enl 6 mk 0;0)}ovi:I:III{ovs(y;z;0;x)}sci:I:III{ovs(y;z;enl 6 mk 0;x)}scl:V:II{x?(rx y;xp:x+8;xp::(I xp)lcat y)}
ovs:I:IIII{(1~ary y)? :fxp(x;y;z);n:nn x;x rxn n;r:x3;o:1;(~r)?(r:fst x;o:0;n-:1;z scl r);y rxn n;n/(r:y cal r l2 x atx mki i+1-o;z scl r);dx x;dx y;(~z)? :r;dx r;fst(z)}
fxp:I:III{t:x;rx x;1?/(rx x;rx y;r:y atx x;((r match x)+r match t)?(dx x;dx y;dx t;z?r:(fst z)lcat r; :r);z scl x;dx x;x:r);x}
ecr:I:III{(1~ary z)? :whl(x;y;z;0);(7~tp y)?(rld y;k:I y+8;v:I y+12; :k mkd ecr(x;v;z));n:nn y;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(rp::z cal x l2 y atx mki i;rp+:4);dx z;dxyr}
ecl:I:III{(1~ary z)? :whl(x;y;z;enl 6 mk 0);(7~tp x)?(rld x;k:I x+8;v:I x+12; :k mkd ecl(v;y;z));n:nn x;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(rp::z cal (x atx mki i)l2 y;rp+:4);dx z;dxyr}
whl:I:IIII{t:tp x;t?(((~t~2)+~1~nn x)?!;dx x; :nlp(y;z;x3;I x+8));r:y;x3 scl r;n:mki 0;1?/(rx x;rx z;r:z atx r;x3 scl r;rx r;t:x atx r;(t match n)?(dx t;dx n;dx z;dx x;x3?(dx r;r:fst x3); :r);dx t);x}
bin:I:II{v2;(~xt~yt)?!;r:2 mk yn;r8;w:C xt;yn/(rp::ibin(xp;yp;xn;xt);rp+:4;yp+:w);dxyr}
ibin:I:IIII{k:0;j:z-1;w:C x3;1?/((k>'j)? :k-1;h:(k+j)>>1;$[((I.x3)(x+w*h;y));j:h-1;k:h+1]);x}
nlp:I:IIII{(x3<0)?!;r:x;y rxn x3;z scl x;x3/(r:y atx r;z scl r);dx y;z?(dx r;r:fst z);r}
ecd:I:III{v2;ext;n:nn x;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(c:mki i;rx c;rp::z cal(x atx c)l2 y atx c;rp+:4);dx z;dxyr}
val:I:I{v1;xt?[((x<256)? :x;n:nn x;rl x;r:6 mk n;mv(r+8;x+8;4*n);(n~4)?(r+20)::mki I r+20;dx x; :r);(r:prs x;n:58~I r+8;r:r evl 0;n?(dx r; r:0));;;r:zim x;;r:x evl 0;(r:I x+12;rx r;dx x);!];r}
lup:I:II{kv;(m~n)?(y? :x lup 0;dx x; :0);r:I v+8+4*m;rx r;dxr}kv:{k:I kkey;v:I kval;y?(k:I y+8;v:I y+12);n:536870911&I k;m:k fnx 8+x}
asn:I:III{v1;(~xt~5)?!;rx z;kv;(n~m)?y?(y:0;kv);(n~m)?(y?!;kkey::k cat x;kval::v lcat z; :z);vp:8+v+4*m;dx I vp;vp::z;dx x;z}
asd:I:II{rld x;v:I x+8;s:I x+12;a:I x+16;u:I x+20;(~v~58)?(rx s;r:s lup y;a?(rx a;r:r atx a);u:v cal r l2 u);a?(rx s;u:asi(s lup y;a;u));rx s;r:asn(s;y;u);dx s;r}
asi:I:III{v2;(xt~7)?(yt~5)?(rld x;k:I xp;v:I xp+4;rx k;y:k fnd y; :k mkd asi(v;y;z));
(yt~6)?((xt~7)?(rld x;k:I xp;v:I xp+4;rx y;f:fst y;$[f;(rx k;f:k fnd f);(f:seq(0;nn k;1))]; :k mkd asi(v;(enl f)cat y drop 1;z)); ((~xt~6)+~yt~6)?!; r:x take xn;r8;rx y;a:fst(y);y:y drop 1;(1~nn y)?y:fst y;(~a)?a:seq(0;xn;1);(~2~tp a)?!;an:nn a;ap:a+8;(an~1)?(dx a;ri:rp+4*I ap;ri::asi(I ri;y;z); :r);(~yn~2)?!;(~6~tp z)?(z:(enl z)take an);(~an~nn z)?!;rxn(y;an-1);rl z;zp:z+8;an/(ri:rp+4*I ap;ri::asi(I ri;y;I zp);ap+:4;zp+:4);dx a;dx z; :r);
(~yt~2)?!;zt:tp z;zn:nn z;zp:8+z;(yn>1)?(zn~1)?(~zn~yn)?((~zn~1)?!;z:z take yn;zn:yn;zp:z+8); (xt<5)?((~zt~xt)?!;r:xt mk xn;r8;w:C xt;mv(rp;xp;w*xn);yn/(k:I yp;mv(rp+w*k;zp;w);yp+:4;zp+:w);dx x;dx y;dx z; :r); (xt~6)?(r:x take xn;r8;(~zt~6)?$[yn~1;(z:enl z;zn:1);z:lx z];(~yn~zn)?!;zp:z+8;rl z;yn/(k:I yp;t:rp+4*k;dx I t;t::I zp;yp+:4;zp+:4);dx y;dx z; :r);!;x}
swc:I:II{v1;i:1;(i<xn)?/(r:I xp+4*i;rx r;r:r evl y;((~i\2)|(i~xn-1))?(dx x; :r);dx r;i+:1;(~I r+8)?i+:1);dx x;0}
ras:I:III{v:I x+8;(y~3)?(v<256)?((v~58)+v>128)?((v>128)?v-:128;r:I x+12;r rxn 2;s:fst r;a:r drop 1;$[nn a;(a:a ltr z;an:nn a;(an~1)?a:fst a);(dx a;a:0)];u:I x+16;rx u;dx x; :(l3(v;s;a))lcat u evl z);0}
ltr:I:II{v1;(~xt~6)? :x;rl x;r:6 mk xn;r8;xn/(rp::(I xp)evl y;rp+:4;xp+:4);dxr}rtl:I:II{v1;(~xt~6)? :x;rl x;r:6 mk xn;rp:r+8+4*xn;xp+:4*xn;xn/(rp-:4;xp-:4;rp::(I xp)evl y);dxr}
evl:I:II{v1;(~xt~6)?((xt~5)?(xn~1)?(r:x lup y;(~r)?!; :r); :x);(~xn)? :x;(xn~1)? :(fst x)rtl y;v:I xp;(v~36)?(xn>3)?  :x swc y;r:ras(x;xn;y);r? :asd(r;y);(v~128)? :lst x ltr y;x:x rtl y;xn:nn x;xp:x+8;(xn~2)?(rl x;r:(I xp)atx I xp+4;dx x; :r);a:(xp+4)fnl xn-1;a?(rx I x+8; :prj(I x+8;x drop 1;a));rx I xp;(I xp)cal x drop 1}
prj:I:III{r:0 mk 3;(r+8)::x;(r+12)::y;(r+16)::z;r}
fnl:I:II{r:0;y/((~I x)?((~r)?r:2 mk 0;r:r ucat mki i);x+:4);r}
uqg:I:II{v1;r:xt mk 0;n:0;w:C xt;xn/(m:r fnx xp;(m~n)?(rx x;r:r cat x atx mki i;y?y:y lcat 2 mk 0;n+:1);y?(yi:y+8+4*m;yi::(I yi)cat mki i);xp+:w);y?r:r l2 y;dxr}
unq:I:I{x uqg 0}grp:I:I{x uqg 6 mk 0}
flr:I:I{v1;(xt>5)? :x ech 223;(xt~2)?(r:1 mk xn;r8;xn/((rp+i)::C?I xp;xp+:4);dx x; :r);(xt~3)?(r:2 mk xn;r8;xn/(rp::I?'F xp;xp+:8;rp+:4);dx x; :r);(xt~4)? :zre x;x}
ang:F:FF{p:57.29577951308232*y atan2 x;(p<0.0)?(p+:360.0);p}
cst:I:II{v2;(xt~5)?(yt~1)?(dx x; :sc y);((~xt~2)+~xn~1)?!;dx x;x:I x+8;(x<'0)?(x:-x;n:yn%C x;(~yn~n*C x)?!;r:use y;r::n|x<<29; :r);(0~yn)?(dx y;(7~x)? :(5 mk 0)mkd 6 mk 0; :x mk 0);((yt~2)?(~I y+8)?(dx y; :fst x mk 0));((yt>x)+yt>4)?!;(8~x)?(n:yn*C yt;r:use y;r::n|1<<29; :r);(yt<'x)?/(y:up(y;yt;yn);yt+:1);y}
flp:I:I{n:nn I x+8;m:nn x;(x ovr 44)atx ecr((mki n)mul seq(0;m;1);seq(0;n;1);43)}
rnd:I:I{(~1073741825~I x)?!;dx x;x:I x+8;r:2 mk x;r8;x/(rp::rng 0;rp+:4);r}
rng:I:I{x:I 12;x^:x<<13;x^:x>>17;x^:x<<5;12::x;x}
xxx:I:I{!;x}
drw:I:II{v2;w:I xp;(yt~7)?draw(w;I xp+4;y);(yt~2)?draw(w;yn%w;y);dx x;dx y;0}
sadv:I:I{$[x~39;r:1;x~47;r:1;x~92;r:1;r:0];r}

out:I:I{rx x;r:kst x;(r+8)printc nn r;dx r;x}
kst:I:I{t:tp x;(~nn x)?(t>1)?(t<6)?(dx x;r:((mkc 48)cc 35)cc 48;(t~3)?r:r cc 46;(t~4)?r:r cc 97;(t~5)?(r+10)::C?96; :r);(7~t)?(r:I x+8;t:I x+12;rld x;r:(kst r)cc 33;(0~nn t)?(dx t;t:mki 0); :r ucat kst t);$[6~t;((1~nn x)? :(mkc 44)ucat kst fst x;x:x ech 235);x:str x];t?[r:x;r:((mkc 34)ucat x)cc 34;;;;r:(mkc 96)ucat x jon mkc 96;r:((mkc 40)ucat x jon mkc 59)cc 41;r:x jon mkc 32];r}
str:I:I{v1;(xt~1)? :x;(~xt)? :x cg xn;((xt>5)+~xn~1)? :x ech 164;xt?[;;r:ci I xp;r:cf F xp;r:(F xp)cz F xp+8;(r:I xp;rx r);!];dxr}
cc:I:II{n:nn x;((1 bk n)<1 bk n+1)? :x ucat mkc y;(x+8+n)::C?y;x::1+I x;x}
ng:I:II{y?x:(mkc 45)ucat x;x}
cg:I:II{((~x)+x~128)? :1 mk 0;(x<127)? :mkc x;(x<256)? :(mkc x-128)cc 58;(y~3)?(rl x;dx I x+16;r:kst I x+12;(r+8)::C?91;(r+7+nn r)::C?93;r:(str I x+8)ucat r);(y~4)?(r:I x+8;rx r);dxr}
ci:I:I{(~x)? :mkc 48;m:0;(x<'0)?(x:0-x;m:1);r:1 mk 0;x?/(c:x\10;r:r cc 48+c;x%:10);(~nn r)?r:r cc 48;(rev r)ng m}
cf:I:F{(~x~x)? :(mkc 48)cc 110;(x~0.0)? :((mkc 48)cc 46)cc 48;m:0;(x<0.0)?(m:1;x:-x);(x>9218868437227405311f)? :((mkc 48)cc 119)ng m;e:0;(x>1000.0)?/(e+:3;x%:1000.0);d:7;(x<1.0)?(d+:1;(x<0.1)?(d+:1;(x<0.01)?(d+:1;(x<0.001)?(d:7;(x<1.0)?/(e-:3;x*:1000.0)))));n:I?'x;r:ci n;x-:F?n;d-:nn r;(d<'1)?d:1;r:r cc 46;t:0;d/(x*:10.0;n:I?x;r:r cc 48+n;x-:F?n;t:(1+t)*~n+~i);r:r drop(-t);e?r:(r cc 101)ucat ci e;r ng m}
cz:I:FF{a:x hypot y;p:I?0.5+x ang y;((cf a)cc 97)ucat ci p}

prs:I:I{v1;(~xt~1)?!;xn+:xp;8::xp;xn?(47~C xp)?8::xp com xn;r:sq xn;$[1~nn r;r:fst r;r:128 cat r];dxr}
sq:I:I{r:6 mk 0;q:(pt x)ex x;q?r:r lcat q;1/(v:ws x;p:I 8;(~v)?(v:C p;v:~(v~59)+v~10);v?((p<x)?8::1+p; :r);8::1+p;(~nn r)?r:r lcat 0;r:r lcat(pt x)ex x);!;x}
ex:I:II{((~x)+(ws y))? :x;r:C I 8;((r is 32)+r~10)? :x;r:pt y;(isv r)?(~isv x)? :l3(r;x;(pt y)ex y);l2(x;r ex y)}
pt:I:I{r:tok x;(~r)?(p:I 8;(p~x)? :0;l:123~C p;(l+40~C p)?(8::1+p;r:sq x;$[l;r:lam(p;I 8;r);(n:nn r;(n~1)?r:fst r;(n>1)?r:enl r)]));1/(p:I 8;b:C p;((p~x)+32~C p-1)? :r;$[b is 16;r:(tok x)l2 r;b~91;(8::1+p;p:sq x;(~nn p)?p:p lcat 0;r:(enl r)cat p); :r]);!;r}
isv:I:I{v1;(~xt)? :1;(xt~6)?(xn~2)?(a:I xp;(a<256)?((a is 16)|((a-128)is 16))? :1);0}
lac:I:II{v1;(xt~6)?xn/(y:(I xp)lac y;xp+:4);(xt~5)?(xn~1)?(p:I xp;(1~nn p)?(r:(C 8+p)-119;(r>y)?(r<4)? :r));y}
loc:I:II{v1;(~xt~6)? :y;xn/(y:(I xp)loc y;xp+:4);xp:x+8;(xn~3)?(58~I xp)?(r:I xp+4;rx r;s:fst r;n:nn y;(n~y fnx s+8)?(rx s;y:y cat s);dx s);y}
lam:I:III{$[91~C 1+x;(z:fst z;rx z;a:((fst z)drop 1)ovr 44;(~a)?a:5 mk 0;z:z drop 1);(r:I xyz;rx r;a:r take z lac 0)];v:nn a;a:z loc a;n:y-x;t:1 mk n;mv(t+8;x;n);r:0 mk 4;(r+8)::t;(r+12)::z;(r+16)::a;(r+20)::v;r}
ws:I:I{p:I 8;(47~C p)?(b:C p-1;((b~32)+b~10)?p:p com x);1/((p~x)?(8::p; :1);b:C p;((b~10)+(b is 64))?(8::p; :0);p+:1;(47~C p)?p:p com x);x}com:I:II{(x<y)?/((10~C x)? :x;x+:1);x}
tok:I:I{(ws x)? :0;p:I 8;b:C p;((b is 32)+b~10)? :0;5/(r:((I.i+136)(b;p;x));r? :r);0}
pun:I:III{(~x is 4)? :0;((x is 4)*y<z)?/(r*:10;r+:x-48;y+:1;x:C y);(~r)?(120~x)? :0;8::y;mki r}
pin:I:III{r:pun(x;y;z);r? :r;(x~45)?(y+:1;(y<z)?(x:C y;8::y;r:pun(x;y;z);(~r)?(8::y-1);r8;rp::-I rp; :r));0}
pfl:I:III{m:0;(x~45)?(t:C y-1;((t~93)+(t~41)+t is 7)? :0;m:1);r:pin(x;y;z);(~r)? :r;y:I 8;(46~C y)?(r:up(r;2;1);y+:1;8::y;(y<z)?(x:C y;q:pun(x;y;z);q?(q:up(q;2;1);f:1.0;((I 8)-y)/f*:10.0;r8;((F rp)<0.0)?(8+q)::-F 8+q;rp::(F rp)+(F 8+q)%f;dx q)));p:I 8;(p<z)?(101~C p)?(8::p+1;q:pin(C 1+p;1+p;z);(~q)?(8::p; :r);e:I q+8;dx q;f:F r+8;(e<'0)?/(f%:10.0;e+:1);(e>0)?/(f*:10.0;e-:1);(r+8)::f);m?(f:F r+8;(f>0.0)? (r+8)::-f);r}
num:I:III{r:pfl(x;y;z);(~r)? :r;y:I 8;x:C y;((119~x)+(110~x)+(112~x)+97~x)?((2~tp r)?r:up(r;2;1);y+:1;8::y;(~97~x)?((112~x)?f:3.141592653589793*F r+8;(110~x)?f:18444492273895866368f;(119~x)?f:9218868437227405312f;(r+8)::f; :r);r:up(r;3;1);a:pfl(C y;y;z);(~a)?a:mki 0;(2~tp a)?a:up(a;2;1);r:r atx a);r}
nms:I:III{r:num(x;y;z);(~r)? :r;1/(y:I 8;x:C y;((y+2)>z)? :r;(~x~32)? :r;y+:1;8::y;q:num(C y;y;z);(~q)?(8::y-1; :r);r:r upx q;q:q upx r;r:r cat q);r}
vrb:I:III{(~x is 24)? :0;(32~C y-1)?((x~92)?(8::1+y; :160);(x~39)?y+:1);r:C y;(z>1+x)?(58~C 1+y)?(y+:1;r+:128);8::1+y;r}
chr:I:III{(48~x)?(120~C 1+y)? :(2+y)phx z;(~x~34)? :0;a:1+y;1/(y+:1;(y~z)?!;(34~C y)?(n:y-a;r:1 mk n;mv(r+8;a;n);8::1+y; :r));r}
phx:I:II{r:1 mk 0;h:1;1/(c:C x;((y<=x)+~c is 5)?(8::x; :r);c-:(48*c<58)+87*c>96;h:~h;h?r:r cc c+q<<4;q:c;x+:1);x}
nam:I:III{(~x is 3)? :0;a:y;1/(y+:1;((y~z)+~(C y)is 7)?(n:y-a;r:1 mk n;mv(r+8;a;n);8::y; :sc r));x}
sym:I:III{(~x~96)? :0;y+:1;x:C y;8::y;(y<z)?(r:nam(x;y;z);r? :r;r:chr(x;y;z);r? :sc r);r:5 mk 1;(r+8)::1 mk 0;r}
sms:I:III{r:sym(x;y;z);(~r)? :r;1/(y:I 8;q:sym(C y;y;z);(~q)? :enl r;r:r cat q);r}
is:I:II{y&cla x}cla:I:I{(128<x-32)? :0;C 128+x}
160!{204840484848485040604848484848504444444444444444444448604848484848424242424242424242424242424242424242424242424242424240506048484041414141414141414141414141414141414141414141414141414048604800}

000:{xxx; gtc; gti; gtf; gtl; gtl; xxx; mod; xxx; eqc; eqi; eqf; eqz; eqL; eqL; xxx; abc; abi; abf; abz; nec; nei; nef; nez; xxx; moi; xxx; xxx; sqc; sqi; sqf; sqz}
032:{xxx; mkd; xxx; rsh; cst; diw; min; ecv; ecd; epi; mul; add; cat; sub; cal; ovv; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; dex; xxx; les; eql; mor; fnd}
064:{atx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; QQQ; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ecl; scv; sci; exc; cut}
096:{xxx; xxx; xxx; xxx; drw; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ecr; max; ovi; mtc; xxx}
128:{xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; chr; nms; vrb; nam; sms; xxx; xxx; xxx; adc; adi; adf; adz; suc; sui; suf; suz; muc; mui; muf; muz; dic; dii; dif; diz}
160:{out; til; xxx; cnt; str; sqr; wer; epv; ech; ecp; fst; abx; enl; neg; val; riv; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; lst; xxx; grd; grp; gdn; unq}
192:{typ; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; scn; liv; spl; srt; flr}
224:{xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; kst; lgf; xxx; xxx; xxx; prs; xxx; rnd; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ovr; rev; jon; not; xxx}

QQQ.I:II{n:nn y;r:3 mk 1;(r+8)::(y+8)norm n;dxyr}
norm.F:II{s:0.0;y/(v:F x;(~v~0.0)?$[s<v;(t:s%v;r:1.0+r*t*t;s:v);(t:v%s;r+:t*t)];x+:8);s*%r}

\
01234567   xt:x>>29       xn:x&536870911 (-1+1<<29)
Fcifzsld   xt~0(function) x<256(basic) x<128(dyadic)
4148x444   xn~2(derived)  adv  verb
	   xn~3(proj)     verb argv empty-index
	   xn~4(lambda)   str  tree args arity

+  add abx                 abs:+z              memory
-  sub neg                                     0..  7   type sizes   0 1 4 8 16 4 4 0
*  mul fst                                     8.. 11   parse cur (pp)
%  div sqr                 conj:%z            12.. 15   rng state
&  min wer   prs flp       ang:&z             16..127   free pointers (4*i) for bt i, i:4..31
|  max rev                                   128..131   memsize log2
<  les grd                                   132..135   k-tree keys
>  mor gdn                                   136..139   k-tree values
=  eql grp                                   140..143   trp line
~  mtc not   match                           144..147   trp col
!  mkd til   seq           z:re!im  im:!z    148..151   `x`y`z
,  cat enl                                   152..155   
^  exc asc                                   156..159   
$  str cst   sc cs                           160..255   char map az|AZ|NM|VB|AD|TE
#  rsh cnt   take                            256.....   buckets/heap
_  drp flr   drop          re:_z            
?  fnd unq   fnd fnx                         (:;`x;y)          assign      x:y
@  atx typ                 z:abs@ang  z@ang  (+;(`x;a;b;c);y)  assign(m/i) s[a;b;c]+:y
.  cal val                 im:. z            (;a;b;c)   (*128) sequence    a;b;c      ::x(last) 
                                             ((/;+);1 2 3)     adverbs     +/1 2 3  :[x;y](dex)                 
+'x  ech(168)      x+'y  ecd(40)        x'y  bin            b:A/x  x:A\b  qr:A\0   x:qr\b  (fz)    
+/x  ovr(251),fxp  x+/y  ecr(123),n/whl x/y  mod,mmul(L)?     
+\x  scn(219),fxp  x+\y  ecl(91),n/whl  x\y  y%x,solve(L)?  \(help)  \\(exit) \d(dump) \w(k.ws)
+':x ecp(169)      x+':y epi(41)        x':y win?           \c(clear console)       
+/:x ?(253)        x+/:y ovi(125)       x/:y join           \L100 F  (loop F[ui] with delay ms)
+\:x ?(221)        x+\:y sci(93)        x\:y split          \e(edit)  \eFILE  \e`VAR (ESC quit) 
