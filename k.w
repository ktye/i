sin:F:F{}cos:F:F{}exp:F:F{}log:F:F{}atan2:F:FF{}hypot:F:FF{}draw:V:III{}grow:I:I{}printc:V:II{}
ini:I:I{0::134480132;4::67372048;12::1887966018;128::x;p:256;i:8;(i<x)?/((4*i)::p;p*:2;i+:1);kkey::enl 1 mk 0;kval::enl 0;xyz::((mks 120)cat mks 121)cat mks 122;x}kkey:{132}kval:{136}xyz:{148}
bk.I:II{r:32-*7+y*C x;(r<4)? :4;r}
mk:I:II{t:x bk y;i:4*t;m:4*I 128;(~I i)?/((i>='m)?(m:grow 1+i%4;128::m;i::1<<i>>2;m:i;i-:4);i+:4);(128~i)?!;a:I i;i::I a;j:i-4;(j>=4*t)?/(u:a+1<<j>>2;u::I j;j::u;j-:4);a::y|x<<29;(a+4)::1;a}
mki:I:I{r:2 mk 1;(r+8)::x;r}mkf:I:F{r:3 mk 1;(r+8)::x;r}
mkd:I:II{(~5~tp x)?(r:4 mk 1;(r+8)::F?0;(r+16)::F?1; :x add y mul r);r:nn x;(1~r)?(~1~nn y)?y:enl y;y:lx y;(~r~nn y)?!;r:x l2 y;r::2|7<<29;r}
mkc:I:I{r:1 mk 1;(r+8)::C?x;r}mks:I:I{sc mkc x}mkz:I:II{r:x mkd y;r::2|6<<29;r}l2.I:II{r:6 mk 2;(r+8)::x;(r+12)::y;r}l3.I:III{r:6 mk 3;(r+8)::x;(r+12)::y;(r+16)::z;r}
nn:I:I{(x<256)? :1;536870911&I x}tp:I:I{(x<256)? :0;(I x)>>29}v1:{xt:tp x;xn:nn x;xp:8+x}v2:{v1;yt:tp y;yn:nn y;yp:8+y}r8:{rp:r+8}
fr.V:I{v1;t:4*xt bk xn;x::I t;t::x}dx:V:I{(x>255)?(xr:I x+4;(x+4)::xr-1;(1~xr)?(v1;((~xt)+xt>5)?xn/(dx I xp+4*i);fr x))}dxr:{dx x;r}dxyr:{dx x;dx y;r}rx:V:I{x rxn 1}rxn.V:II{(x>255)?(x+:4;x::y+I x)}rl.V:I{v1;xn/(rx I xp;xp+:4)}rld.V:I{rl x;dx x}kvdx:{rld x;k:I x+8;v:I x+12}kvdy:{rld y;k:I y+8;v:I y+12}
lx.I:I{v1;(xt~6)? :x;((xt~7)+xn~1)? :enl x;(0~xt)?!;r:6 mk xn;r8;x rxn xn;xn/(rp::x atx mki i;rp+:4);dxr}
til:I:I{v1;(4~xt)? :zim x;(6~xt)? :x ech 161;(7~xt)?(r:I xp;rx r;dx x; :r);(~2~xt)?!;n:I xp;dx x;(n<'0)? :tir -n;seq(0;n;1)}seq.I:III{r:2 mk y;rp:8+r;y/(rp::z*i+x;rp+:4);r}tir.I:I{r:2 mk x;rp:4+r+4*x;x/(rp::i;rp-:4);r}
ext:{(~xn~yn)?((xn~1)?(x:x take yn;xn:yn;xp:x+8);(yn~1)?(y:y take xn;yn:xn;yp:y+8))}
upx.I:II{t:tp x;yt:tp y;(t~yt)? :x;((t~7)+(yt~7))?!;(yt~6)? :lx x;n:nn x;(t<yt)?/(x:up(x;t;n);t+:1);x}
up.I:III{r:(y+1) mk z;xp:x+8;r8;$[(1~y);z/(rp::C xp+i;rp+:4);(2~y);z/(rp::F?'I xp;rp+:8;xp+:4);(3~y);z/(rp::F xp;(rp+8)::0.;xp+:8;rp+:16);!];dxr}upxy:{x:x upx y;y:y upx x}
atx:I:II{v2;(~xt)?( :x cal enl y);(xt~7)? :atd(x;y;yt);(yt>5)? :ecr(x;y;64);(yt~3)?(xt<5)? :x phi y;(~yt~2)?!;r:xt mk yn;r8;w:C xt;yn/(yi:I yp;(xn<=yi)?!;mv(rp;xp+w*yi;w);rp+:w;yp+:4);(xt>5)?rl r;(yn~1)?(xt~6)?(rx I r+8;dx r;r:I r+8);dxyr}
atm.I:II{(~nn y)?(dx x; :y);rx y;f:fst y;t:y drop 1;(1<nn t)? :(x atx f)atm t;t:fst t;nf:nn f;((~f)+~nf~1)?(f?x:x atx f; :ecl(x;t;64));(x atx f)atx t}
atd.I:III{k:I x+8;v:I x+12;(z~5)?(rx k;y:k fnd y;z:2);rx v;dx x;v atx y}
cal:I:II{y:lx y;v2;xt? :x atm y;(yn~1)?(((sadv x)|sadv x-128)? :(I.x)(fst y);(x<128)?x+:128);(x<128)?((~yn~2)?!;rld y; :(I.x)(I yp;I yp+4));(x<256)?((~yn~1)?!; :(I.x)(fst y));(xn~2)?(rld x;a:I xp;$[1~yn; :(I.a+128)(fst y;I xp+4);2~yn;(rld y; :(I.a)(I yp;I yp+4;I xp+4));!]);(xn~3)?(rl x;(1~yn)?y:fst y;r:asi(I x+12;I x+16;y);v:I x+8;dx x; :v cal r);(xn~4)?(a:I x+20;(a>yn)?(a-:yn;a/y:y lcat 0; :prj(x;y;seq(yn;a;1))); :lcl(x;y;0));!;x}
lcl.I:III{fn:I x+20;(~fn)?(dx y;y:6 mk 0);(~fn~nn y)?!;yp:y+8;a:I x+16;ap:a+8;an:nn a;l:2 mk an;lp:l+8;sp:8+I kval;an/(d:sp+4*I ap;lp::I d;v:0;(i<fn)?(v:I yp;rx v;yp+:4);d::v;ap+:4;lp+:4);dx y;r:I x+12;rx r;r:evl r;sp:8+I kval;lp:l+8;ap:a+8;an/(d:sp+4*I ap;dx I d;d::I lp;ap+:4;lp+:4);dx l;dxr}
rev.I:I{(7~tp x)?(kvdx;:(rev k)mkd rev v);n:nn x;(~n)? :x;x atx tir n}fst:I:I{v1;(~xn)?(dx x;(xt~0)? :0;(xt~5)? :sc 1 mk 0;(xt>5)? :6 mk 0; :(mki xt)cst mkc 0);(~xt)? :x;(xt~7)?( :fst val x);x atx mki 0}lst.I:I{(7~tp x)? :lst val x;(~nn x)? :fst x;x atx mki (nn x)-1}
cut.I:II{v2;(yt~7)?((xt~2)?((~1~xn)?!;kvdy;rx x; :(x cut k)mkd x cut v);rx y; :((til y)exc x)tkd y);(~xt~2)?!;(xn~1)?(r:y drop I xp;dx x; :r);r:6 mk xn;r8;xn/(a:I xp;b:I xp+4;(i~xn-1)?b:yn;(b<a)?!;rx y;rp::y atx seq(a;b-a;1);xp+:4;rp+:4);dxyr}
rsh.I:II{v2;(yt~7)?((xt~2)?((~1~xn)?!;kvdy;rx x; :(x rsh k)mkd x rsh v); :x tkd y);(~xt~2)?!;n:prod(xp;xn);r:y take n;(xn~1)?((yt~6)?(n~1)?r:enl r;dx x; :r);xn-:1;xe:xp+4*xn;xn/(m:I xe;n:n%m;n:xp prod xn-i;r:(seq(0;n;m))cut r;(1~m)?(i>0)?$[6~tp r;r:r ech 172;r:enl r];xe-:4);(1~I xe)?r:enl r;dxr}prod.I:II{r:1;y/(r*:I x;x+:4);r}

take.I:II{(~nn x)?x:fst x;xn:nn x;o:0;(y<'0)?(o:xn+y;y:-y;(o<'0)? :x);r:seq(o;y;1);(xn<y)?(r8;y/(rp::i\xn;rp+:4));x atx r}drop:I:II{v1;a:y;(y<'0)?(y:-y;a:0);(y>xn)?(dx x; :xt mk 0);x:x atx seq(a;xn-y;1);(xt~6)?(1~xn-y)?x:enl x;x}
tkd.I:II{t:tp x;kvdy;(~t~5)?!;rx k;x:k fnd x;rx x;v:v atx x;(1~nn x)?v:enl v;(k atx x)mkd v}
phi.I:II{n:nn y;r:4 mk n;r8;yp:y+8;n/(p:0x399d52a246df913f*F yp;rp::cos p;(rp+8)::sin p;rp+:16;yp+:8);dx y;x mul r}
use.I:I{(1~I x+4)? :x;v1;r:xt mk xn;r8;mv(rp;xp;xn*C xt);dx x;r}mv:V:III{z/(x+i)::C?C y+i}
cat:I:II{v2;(~xt)?(x:enl(x);xt:6);(xt~yt)? :x ucat y;(xt~6)? :x ucat lx y;(yt~6)? :(lx x)ucat y;(lx x)ucat lx y}
ucat:I:II{v2;(xt>5)?(rl x;rl y);(xt~7)?(r:((I x+8)ucat I y+8)mkd(I x+12)ucat I y+12;dx x;dx y; :r);r:xt mk xn+yn;w:C xt;mv(r+8;xp;w*xn);mv(r+8+w*xn;yp;w*yn);dxyr}
lcat:I:II{x:use x;v1;((xt bk xn)<(xt bk xn+1))?(r:xt mk xn+1;rld x;mv(r+8;xp;4*xn);x:r;xp:x+8);(xp+4*xn)::y;x::(xn+1)|6<<29;x}
enl:I:I{r:6 mk 1;(r+8)::x;r}cnt.I:I{(7~tp x)?x:til x;r:mki nn x;dxr}typ.I:I{v1;r:2 mk 1;(8+r)::xt;dxr}not.I:I{t:tp x;(t>5)? :x ech 126;(~t)?((~x)? :mki 1;dx x; :mki 0);x eql mki 0}
wer.I:I{v1;(xt~1)? :prs x;(xt~4)? :zan(x;xn;xp);(xt~6)? :flp x;(~xt~2)?!;n:0;xn/(n+:I xp;xp+:4);xp:8+x;r:2 mk n;r8;j:0;(j<xn)?/((I xp)/(rp::j;rp+:4);xp+:4;j+:1);dxr}
mtc.I:II{r:2 mk 1;(r+8)::x match y;dxyr}match.I:II{(x~y)? :1;(~(I x)~I y)? :0;v1;yp:y+8;m:0;$[~xt; :1;1~xt;n:xn;2~xt;n:xn<<2;3~xt;n:xn<<3;4~xt;n:xn<<4;5~xt;n:xn<<2;(xn/((~((I xp) match I yp))? :0;xp+:4;yp+:4); :1)];n/(~(C xp+i)~C yp+i)? :0;1}
fnd.I:II{v2;(~xt~yt)?!;r:2 mk yn;r8;w:C yt;yn/(rp::x fnx yp;rp+:4;yp+:w);dxyr}fnx.I:II{v1;eq:8+xt;w:C xt;xn/(((I.eq)(xp;y))? :i;xp+:w);xn}
lop.I:III{t:tp y;(~t)?( :fxp(x;y;z));(6~t)?(rld y;f:I y+12;y:I y+8; :whl(y;x;f;z));dx z;0}
jon.I:II{r:lop(x;y;0);r? :r;v1;((~xt~6)+~xn)?(dx y; :x);rl x;r:I xp;y rxn xn-2;(xn-1)/(xp+:4;r:(r cat y)cat I xp);dxr}
spl.I:II{r:lop(x;y;enl 6 mk 0);r? :r;rx x;yn:nn y;r:x fds y;(~nn r)?(dx r;:enl x);r:((mki 0)cat r)cut x;rn:(nn r)-1;r8;rn/(rp+:4;rp::(I rp) drop yn);r}
fds.I:II{v2;((~xt~yt)+xt>5)?!;(xn<yn)?(dx x;dx y; :2 mk 0);(~yn)?(dx x;dx y; :(seq(0;xn;1))drop 1);r:2 mk 0;w:C xt;eq:8+xt;j:0;(j<xn)?/(a:0;yn/(k:w*i;a+:(I.eq)(xp+k;yp+k));(a~yn)?(r:r ucat mki j;j+:yn-1;xp+:w*yn-1);xp+:w;j+:1);dxyr}
exc.I:II{n:mki nn y;rx x;x atx wer n eql y fnd x}
srt.I:I{rx x;x atx grd x}gdn.I:I{rev grd x}grd.I:I{v1;r:seq(0;xn;1);y:seq(0;xn;1);r8;msrt(y+8;rp;0;xn;xp;xt);dxyr}
msrt.V:IIIIII{((x3-z)>=2)?(c:(x3+z)%2;msrt(y;x;z;c;x4;x5);msrt(y;x;c;x3;x4;x5);mrge(x;y;z;x3;c;x4;x5))}
mrge.V:IIIIIII{k:z;j:x4;w:C x6;i:z;(i<x3)?/(c:k>=x4;(~c)?$[j>=x3;c:0;c:(I.x6)(x5+w*I x+k<<2;x5+w*I x+j<<2)];$[c;(a:j;j+:1);(a:k;k+:1)];(y+i<<2)::I x+a<<2;i+:1)}
gtc.I:II{(C x)>C y}gti.I:II{(I x)>'I y}gtf.I:II{(F x)>F y}eqc.I:II{(C x)~C y}eqi.I:II{(I x)~ I y}eqf.I:II{((I x)~I y)*(I x+4)~I y+4}eqz.I:II{(x eqf y)*(x+8)eqf y+8}eqL.I:II{(I x)match I y}
gtl.I:II{x:I x;y:I y;v2;(~xt~yt)? :xt>yt;n:xn;(yn<xn)?n:yn;w:C xt;n/(a:xp+i*w;b:yp+i*w;((I.xt)(a;b))? :1;((I.xt)(b;a))? :0);xn>yn}
sc:I:I{k:I kkey;n:nn k;x:enl x;r:k fnx x+8;$[r<n;dx x;(kkey::k cat x;kval::(I kval)lcat 0)];r:mki r;r::1|5<<29;r}
cs:I:I{r:I x+8;r:I 8+(I kkey)+4*r;rx r;dxr}
eql.I:II{cmp(x;y;1)}mor.I:II{cmp(x;y;0)}les.I:II{cmp(y;x;0)}cmp.I:III{upxy;v2;ext;(xt~6)? :ecd(x;y;62-z);f:xt;z?f+:8;w:C xt;r:2 mk xn;r8;xn/(rp::(I.f)(xp;yp);xp+:w;yp+:w;rp+:4);dxyr}
min.I:II{mia(x;y;38)}max.I:II{mia(x;y;124)}mia.I:III{upxy;v2;ext;(xt~6)? :ecd(x;y;z);rx x;rx y;$[z~38;a:x les y;a:x mor y];a:wer a;rx a;asi(y;a;x atx a)}
nd.I:IIII{upxy;v2;ext;(xt~6)? :ecd(x;y;x3);w:C xt;f:z+xt;r:xt mk xn;r8;xn/((V.f)(xp;yp;rp);xp+:w;yp+:w;rp+:w);dxyr}
nm.I:III{v1;(xt>5)? :ech(x;z);r:use x;r8;w:C xt;y+:xt;xn/(((V.y)(xp;rp));xp+:w;rp+:w);(xt~4)?(y~19)? :zre r;r}
nmf.I:II{dx x;x:I x+8;y:use(mki 3)cst y;yp:y+8;(nn y)/(yp::(F.x)(F yp);yp+:8);y}
add.I:II{nd(x;y;143;43)}sub.I:II{nd(x;y;147;45)}mul.I:II{nd(x;y;151;42)}diw.I:II{nd(x;y;155;37)}mod.I:II{nd(x;y;23;7)}
adc.V:III{z::C?(C x)+C y}adi.V:III{z::(I x)+I y}adf.V:III{z::(F x)+F y}adz.V:III{adf(x;y;z);adf(x+8;y+8;z+8)}
suc.V:III{z::C?(C x)-C y}sui.V:III{z::(I x)-I y}suf.V:III{z::(F x)-F y}suz.V:III{suf(x;y;z);suf(x+8;y+8;z+8)}
muc.V:III{z::C?(C x)*C y}mui.V:III{z::(I x)*I y}muf.V:III{z::(F x)*F y}muz.V:III{z::((F x)*F y)-(F y+8)*F x+8;(z+8)::((F x)*F y+8)+(F x+8)*F y}
dic.V:III{z::C?(C x)%C y}dii.V:III{z::(I x)%'I y}dif.V:III{z::(F x)%F y}moi.V:III{x:I x;y:I y;z::(y+I?x\'y)\'y}
diz.V:III{a:F x;b:F x+8;c:F y;d:F y+8;$[(+c)>=(+d);(r:d%c;p:c+r*d;z::(a+b*r)%p;(z+8)::(b-a*r)%p);(r:c%d;p:d+r*c;z::(b+a*r)%p;(z+8)::((b*r)-a)%p)]}
abx.I:I{nm(x;15;171)}neg.I:I{nm(x;19;173)}sqr.I:I{nm(x;27;165)}
abc.V:II{c:C x;$[c is 1;y::C?c-32;y::C?c]}abi.V:II{i:I x;$[(i<'0);y::-i;y::i]}abf.V:II{y::+F x}abz.V:II{y::(F x)hypot F x+8}
nec.V:II{c:C x;$[c is 2;y::C?c+32;y::C?c]}nei.V:II{y::-I x}nef.V:II{y::-F x}nez.V:II{y::-F x;(y+8)::-F x+8}
sqc.V:II{!}sqi.V:II{!}sqf.V:II{y::%F x}sqz.V:II{y::F x;(y+8)::-F x+8}
lgf.I:I{v1;(~xt~3)?!;x:use x;xp:x+8;xn/(xp::log F xp;xp+:8);x}
zre.I:I{x zri 0}zim.I:I{x zri 8}zri.I:II{v1;r:3 mk xn;r8;xp+:y;xn/(rp::F xp;rp+:8;xp+:16);dxr}zan.I:III{r:3 mk y;r8;y/(rp::(F z)ang F z+8;z+:16;rp+:8);dxr}
drv.I:II{r:0 mk 2;(r+8)::x;(r+12)::y;r}ecv.I:I{40 drv x}epv.I:I{41 drv x}ovv.I:I{123 drv x}riv.I:I{125 drv x}scv.I:I{91 drv x}liv.I:I{93 drv x}
ech.I:II{(tp y)? :y bin x;(7~tp x)?(kvdx; :k mkd v ech y);x:lx x;v1;r:6 mk xn;r8;rl x;(y<120)?y+:128;xn/(rx y;rp::y atx I xp;xp+:4;rp+:4);dxyr}
ecp.I:II{rx x;p:fst x;epi(p;x;y)}epi.I:III{n:nn y;(~n)?(dx x;dx z; :y);y rxn n;z rxn n;r:6 mk n;r8;n/(yi:y atx mki i;rx yi;rp::z cal yi l2 x;x:yi;rp+:4);dx yi;dx y;dx z;r}
ovr.I:II{t:tp y;(2~t)? :x mod y;ovs(x;y;0;0)}scn.I:II{t:tp y;t?(t<5)? :x diw y;ovs(x;y;enl 6 mk 0;0)}ovi.I:III{ovs(y;z;0;x)}sci.I:III{ovs(y;z;enl 6 mk 0;x)}scl.V:II{x?(rx y;xp:x+8;xp::(I xp)lcat y)}
ovs.I:IIII{n:nn x;x rxn n;r:x3;o:1;(~r)?(r:fst x;o:0;n-:1;z scl r);y rxn n;n/(r:y cal r l2 x atx mki i+1-o;z scl r);dx x;dx y;(~z)? :r;dx r;fst(z)}
fxp.I:III{t:x;rx x;1?/(rx x;rx y;r:y atx x;((r match x)+r match t)?(dx x;dx y;dx t;z?r:(fst z)lcat r; :r);z scl x;dx x;x:r);x}
ecr.I:III{(7~tp y)?(kvdy; :k mkd ecr(x;v;z));n:nn y;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(rp::z cal x l2 y atx mki i;rp+:4);dx z;dxyr}
ecl.I:III{(7~tp x)?(kvdx; :k mkd ecl(v;y;z));n:nn x;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(rp::z cal (x atx mki i)l2 y;rp+:4);dx z;dxyr}
whl.I:IIII{t:tp x;t?(((~t~2)+~1~nn x)?!;dx x; :nlp(y;z;x3;I x+8));r:y;x3 scl r;n:mki 0;1?/(rx x;rx r;t:x atx r;(t match n)?(dx t;dx n;dx z;dx x;x3?(dx r;r:fst x3); :r);dx t;rx z;r:z atx r;x3 scl r);x}
bin.I:II{v2;(~xt~yt)?!;r:2 mk yn;r8;w:C xt;yn/(rp::ibin(xp;yp;xn;xt);rp+:4;yp+:w);dxyr}
ibin.I:IIII{k:0;j:z-1;w:C x3;1?/((k>'j)? :k-1;h:(k+j)>>1;$[(I.x3)(x+w*h;y);j:h-1;k:h+1]);x}
nlp.I:IIII{(x3<0)?!;r:x;y rxn x3;z scl x;x3/(r:y atx r;z scl r);dx y;z?(dx r;r:fst z);r}
ecd.I:III{v2;ext;n:nn x;r:6 mk n;r8;x rxn n;y rxn n;z rxn n;n/(c:mki i;rx c;rp::z cal(x atx c)l2 y atx c;rp+:4);dx z;dxyr}
val:I:I{v1;$[~xt;((x<256)? :x;rl x;r:6 mk xn;mv(r+8;x+8;4*xn);(xn~4)?(r+20)::mki I r+20;dx x);(1~xt);(r:prs x;n:(2<nn r)*58~I r+8;r:evl r;n?(dx r; r:0));(5~xt);r:lup x;(6~xt);r:evl x;(7~xt);(r:I x+12;rx r;dx x);!];r}
lup:I:I{r:I(I kval)+8+4*I x+8;rx(r);dxr}
asn:I:II{(~5~tp x)?!;p:(I kval)+8+4*I x+8;dx I p;p::y;rx y;dx x;y}
asd.I:I{rld x;v:I x+8;s:I x+12;a:I x+16;u:I x+20;(~v~58)?(rx s;r:lup s;a?(rx a;r:r atx a);u:v cal r l2 u);r:u;rx r;a?(rx s;u:asi(lup s;a;u));dx s asn u;r}
asi.I:III{v2;(xt~7)?(yt<6)?(kvdx;(yt~5)?(rx k;y:k fnd y); :k mkd asi(v;y;z));
(yt~6)?((xt~7)?(kvdx;rx y;f:fst y;$[f;(5~tp f)?(rx k;f:k fnd f);(f:seq(0;nn k;1))]; :k mkd asi(v;(enl f)cat y drop 1;z)); ((~xt~6)+~yt~6)?!; r:x take xn;(1~xn)?r:enl r;r8;rx y;a:fst(y);y:y drop 1;(1~nn y)?y:fst y;(~a)?a:seq(0;xn;1);(~2~tp a)?!;an:nn a;ap:a+8;(an~1)?(dx a;ri:rp+4*I ap;ri::asi(I ri;y;z); :r);(~yn~2)?!;(~6~tp z)?(z:(enl z)take an);(~an~nn z)?!;rxn(y;an-1);rl z;zp:z+8;an/(ri:rp+4*I ap;ri::asi(I ri;y;I zp);ap+:4;zp+:4);dx a;dx z; :r);
(~yt~2)?!;zt:tp z;zn:nn z;zp:8+z;(yn>1)?(zn~1)?(~zn~yn)?((~zn~1)?!;z:z take yn;zn:yn;zp:z+8);
(xt<6)?((~zt~xt)?!;r:use x;r8;w:C xt;yn/(k:I yp;mv(rp+w*k;zp;w);yp+:4;zp+:w);dx y;dx z; :r);
(xt~6)?(r:x take xn;(1~xn)?r:enl r;(1~yn)?(z:enl z;zn:1;zt:6);(~6~zt)?z:lx z;r8;(~yn~zn)?!;zp:z+8;rl z;yn/(k:I yp;(~k<xn)?!;t:rp+4*k;dx I t;t::I zp;yp+:4;zp+:4);dx y;dx z; :r);!;x}
swc.I:I{v1;i:1;(i<xn)?/(r:I xp+4*i;rx r;r:evl r;((~i\2)|(i~xn-1))?(dx x; :r);dx r;i+:1;(~I r+8)?i+:1);dx x;0}
ras.I:II{v:I x+8;(y~3)?(v<256)?((v~58)+v>128)?((v>128)?v-:128;r:I x+12;r rxn 2;s:fst r;a:r drop 1;$[nn a;(a:ltr a;an:nn a;(an~1)?a:fst a);(dx a;a:0)];u:I x+16;rx u;dx x; :(l3(v;s;a))lcat evl u);0}
ltr.I:I{v1;(~xt~6)? :x;rl x;r:6 mk xn;r8;xn/(rp::evl (I xp);rp+:4;xp+:4);dxr}rtl.I:I{v1;(~xt~6)? :x;rl x;r:6 mk xn;rp:r+8+4*xn;xp+:4*xn;xn/(rp-:4;xp-:4;rp::evl I xp);dxr}
evl.I:I{v1;(~xt~6)?((xt~5)?(xn~1)? :lup x; :x);(~xn)? :x;(xn~1)? :rtl fst x;v:I xp;(v~36)?(xn>3)?  :swc x;r:x ras xn;r? :asd r;(v~128)? :lst ltr x;x:rtl x;xn:nn x;xp:x+8;(v~64)?(xn~4)?(rl x;r:asi(I x+12;I x+16;I x+20);dx x; :r);(xn~2)?(rl x;r:(I xp)atx I xp+4;dx x; :r);a:(xp+4)fnl xn-1;a?(rx I x+8; :prj(I x+8;x drop 1;a));rx I xp;(I xp)cal x drop 1}
prj.I:III{r:0 mk 3;(r+8)::x;(r+12)::y;(r+16)::z;r}
fnl.I:II{r:0;y/((~I x)?((~r)?r:2 mk 0;r:r ucat mki i);x+:4);r}
uqg.I:II{v1;r:xt mk 0;n:0;w:C xt;xn/(m:r fnx xp;(m~n)?(rx x;r:r cat x atx mki i;y?y:y lcat 2 mk 0;n+:1);y?(yi:y+8+4*m;yi::(I yi)cat mki i);xp+:w);y?r:r l2 y;dxr}
unq.I:I{x uqg 0}grp.I:I{x uqg 6 mk 0}
flr.I:I{v1;(xt>5)? :x ech 223;(~xt)?(dx x; :mki x);(1~xt)?(dx x; :C xp);(2~xt)?(r:1 mk xn;r8;xn/((rp+i)::C?I xp;xp+:4);dx x; :r);(xt~3)?(r:2 mk xn;r8;xn/(rp::I?'F xp;xp+:8;rp+:4);dx x; :r);(xt~4)? :zre x;!;x}
ang.F:FF{p:57.29577951308232*y atan2 x;(p<0.)?(p+:360.);p}
cst.I:II{v2;(xt~5)?(yt~1)?(dx x; :sc y);((~xt~2)+~xn~1)?!;dx x;x:I x+8;(x<'0)?(x:-x;n:yn%C x;(~yn~n*C x)?!;r:use y;r::n|x<<29; :r);(0~yn)?(dx y;(7~x)? :(5 mk 0)mkd 6 mk 0; :x mk 0);((yt>x)+yt>4)?!;(8~x)?(n:yn*C yt;r:use y;r::n|1<<29; :r);(yt<'x)?/(y:up(y;yt;yn);yt+:1);y}
flp.I:I{n:nn I x+8;m:nn x;(x ovr 44)atx ecr((mki n)mul seq(0;m;1);seq(0;n;1);43)}
rnd.I:I{(~1073741825~I x)?!;dx x;x:I x+8;r:2 mk x;r8;x/(rp::rng 0;rp+:4);r}
rng.I:I{x:I 12;x^:x<<13;x^:x>>17;x^:x<<5;12::x;x}
xxx.I:I{!;x}
drw.I:II{v2;w:I xp;(yt~7)?draw(w;I xp+4;y);(yt~2)?draw(w;yn%w;y);dx x;dx y;0}
sadv.I:I{$[x~39;r:1;x~47;r:1;x~92;r:1;r:0];r}

out.I:I{(~x)? :x;rx x;r:x;(~1~tp r)?r:kst x;(r+8)printc nn r;dx r;x}
kst:I:I{t:tp x;(~nn x)?(t>1)?(t<6)?(dx x;r:((mkc 48)cc 35)cc 48;(t~3)?r:r cc 46;(t~4)?r:r cc 97;(t~5)?(r+10)::C?96; :r);(7~t)?(kvdx;k:(kst k)cc 33;(0~nn v)?(dx v;v:mki 0); :k ucat kst v);$[6~t;((1~nn x)? :(mkc 44)ucat kst fst x;x:x ech 235);x:str x];$[~t;r:x;1~t;r:((mkc 34)ucat x)cc 34;5~t;r:(mkc 96)ucat x jon mkc 96;6~t;r:((mkc 40)ucat x jon mkc 59)cc 41;r:x jon mkc 32];r}
str.I:I{v1;(xt~1)? :x;(~xt)? :x cg xn;((xt>5)+~xn~1)? :x ech 164;$[2~xt;r:ci I xp;3~xt;r:cf F xp;4~xt;r:(F xp)cz F xp+8;5~xt;(rx x;r:cs x);!];dxr}
cc.I:II{n:nn x;((1 bk n)<1 bk n+1)? :x ucat mkc y;(x+8+n)::C?y;x::1+I x;x}
ng.I:II{y?x:(mkc 45)ucat x;x}
cg.I:II{((~x)+x~128)? :1 mk 0;(x<127)? :mkc x;(x<256)? :(mkc x-128)cc 58;(y~2)?(rl x;r:(str I x+12)cat str I x+8);(y~3)?(rl x;dx I x+16;r:kst I x+12;(r+8)::C?91;(r+7+nn r)::C?93;r:(str I x+8)ucat r);(y~4)?(r:I x+8;rx r);dxr}
ci.I:I{(~x)? :mkc 48;m:0;(x<'0)?(x:-x;m:1);r:1 mk 0;x?/(c:x\10;r:r cc 48+c;x%:10);(~nn r)?r:r cc 48;(rev r)ng m}
cf.I:F{(~x~x)? :(mkc 48)cc 110;(x~0.)? :((mkc 48)cc 46)cc 48;m:0;(x<0.)?(m:1;x:-x);(x>0xffffffffffffef7f)? :((mkc 48)cc 119)ng m;e:0;(x>1000.)?/(e+:3;x%:1000.);d:7;(x<1.)?(d+:1;(x<0.1)?(d+:1;(x<0.01)?(d+:1;(x<0.001)?(d:7;(x<1.0)?/(e-:3;x*:1000.)))));n:I?'x;r:ci n;x-:F?n;d-:nn r;(d<'1)?d:1;r:r cc 46;t:0;d/(x*:10.;n:I?x;r:r cc 48+n;x-:F?n;t:(1+t)*~n+~i);r:r drop(-t);e?r:(r cc 101)ucat ci e;r ng m}
cz.I:FF{a:x hypot y;p:I?0.5+x ang y;((cf a)cc 97)ucat ci p}

prs:I:I{v1;(~xt~1)?!;xn+:xp;8::xp;xn?(47~C xp)?8::xp com xn;r:sq xn;$[1~nn r;r:fst r;r:128 cat r];dxr}
sq.I:I{r:6 mk 0;q:(pt x)ex x;q?r:r lcat q;1/(v:ws x;p:I 8;(~v)?(v:C p;v:~(v~59)+v~10);v?((p<x)?8::1+p; :r);8::1+p;(~nn r)?r:r lcat 0;r:r lcat(pt x)ex x);!;x}
ex.I:II{((~x)+(ws y))? :x;r:C I 8;((r is 32)+r~10)? :x;r:pt y;(isv r)?(~isv x)? :l3(r;x;(pt y)ex y);l2(x;r ex y)}
pt.I:I{r:tok x;(~r)?(p:I 8;(p~x)? :0;l:123~C p;(l+40~C p)?(8::1+p;$[l;(a:0;(91~C 1+p)?(8::2+p;a:sq x;(~nn a)?a:a lcat 5 mk 0;a:a ovr 44);r:sq x;r:lam(p;I 8;r;a));(r:sq x;n:nn r;(n~1)?r:fst r;(n>1)?r:enl r)]));1/(p:I 8;b:C p;((p~x)+32~C p-1)? :r;$[b is 16;r:(tok x)l2 r;b~91;(8::1+p;p:sq x;(~nn p)?p:p lcat 0;r:(enl r)cat p); :r]);!;r}
isv.I:I{v1;(~xt)? :1;(xt~6)?(xn~2)?(a:I xp;(a<256)?((a is 16)|((a-128)is 16))? :1);0}
lac.I:II{v1;(xt~6)?((1~xn)?(5~tp I x+8)? :y;xn/(y:(I xp)lac y;xp+:4));(xt~5)?(xn~1)?(p:I xp;(p>y)?(p<4)? :p);y}
loc.I:II{v1;(~xt~6)? :y;xn/(y:(I xp)loc y;xp+:4);xp:x+8;(xn~3)?(58~I xp)?(r:I xp+4;rx r;s:fst r;n:nn y;(n~y fnx s+8)?(rx s;y:y cat s);dx s);y}
lam.I:IIII{$[1~nn z;z:fst z;z:128 cat z];(~x3)?(r:I xyz;rx r;x3:r take z lac 0);v:nn x3;x3:z loc x3;n:y-x;t:1 mk n;mv(t+8;x;n);r:0 mk 4;(r+8)::t;(r+12)::z;(r+16)::x3;(r+20)::v;r}
ws.I:I{p:I 8;(47~C p)?(b:C p-1;((b~32)+b~10)?p:p com x);1/((p~x)?(8::p; :1);b:C p;((b~10)+(b is 64))?(8::p; :0);p+:1;(47~C p)?p:p com x);x}com.I:II{(x<y)?/((10~C x)? :x;x+:1);x}
tok.I:I{(ws x)? :0;p:I 8;b:C p;((b is 32)+b~10)? :0;5/(r:(I.i+136)(b;p;x);r? :r);0}
pui.I:III{(~x is 4)? :0;r:0;((x is 4)*y<z)?/(r*:10;r+:I?x-48;y+:1;x:C y);r?(120~x)? :0;8::y;r}
pin.I:III{u:pui(x;y;z);(~y~I 8)? :mki I?u;(x~45)?(y+:1;(y<z)?(x:C y;8::y;u:pui(x;y;z);(~y~I 8)? :mki -I?u;8::y-1;));0}
pfd.F:IIF{g:1.0;(z<0.)?g:-g;1/(b:C x;$[(x<y)*b is 4;(g*:0.1;z+:g*F?b-48);(8::x; :z)];x+:1);z}
pfl.I:III{m:0;(x~45)?(t:C y-1;((t~34)+(t~93)+(t~41)+t is 7)? :0;m:1);r:pin(x;y;z);y:I 8;((y~z)+~r)? :r;(46~C y)?(r:up(r;2;1);r8;rp::pfd(y+1;z;F rp));y:I 8;(y<z)?(101~C y)?(8::y+1;q:pin(C 1+y;1+y;z);(~q)?(8::y; :r);e:I q+8;dx q;f:F r+8;(e<'0)?/(f%:10.;e+:1);(e>0)?/(f*:10.;e-:1);(r+8)::f);m?(f:F r+8;(f>0.)? (r+8)::-f);r}
num.I:III{r:pfl(x;y;z);(~r)? :r;y:I 8;x:C y;(y<z)?((119~x)+(110~x)+(112~x)+97~x)?((2~tp r)?r:up(r;2;1);y+:1;8::y;(~97~x)?((112~x)?f:3.141592653589793*F r+8;(110~x)?f:0x000000000000f8ff;(119~x)?f:0x000000000000f07f;(r+8)::f; :r);r:up(r;3;1);a:pfl(C y;y;z);(~a)?a:mki 0;(2~tp a)?a:up(a;2;1);r:r atx a);r}
nms.I:III{r:num(x;y;z);(~r)? :r;1/(y:I 8;x:C y;((y+2)>z)? :r;(~x~32)? :r;y+:1;8::y;q:num(C y;y;z);(~q)?(8::y-1; :r);r:r upx q;q:q upx r;r:r cat q);r}
vrb.I:III{(~x is 24)? :0;(32~C y-1)?((x~92)?(8::1+y; :160);(x~39)?y+:1);r:C y;(z>1+x)?(58~C 1+y)?(y+:1;r+:128);8::1+y;r}
chr.I:III{(48~x)?(120~C 1+y)? :(2+y)phx z;(~x~34)? :0;a:1+y;1/(y+:1;(y~z)?!;(34~C y)?(n:y-a;r:1 mk n;mv(r+8;a;n);8::1+y; :r));r}
phx.I:II{r:1 mk 0;h:1;q:0;1/(c:C x;((y<=x)+~c is 5)?(8::x; :r);c-:(48*c<58)+87*c>96;h:~h;h?r:r cc c+q<<4;q:c;x+:1);x}
nam.I:III{(~x is 3)? :0;a:y;1/(y+:1;((y~z)+~(C y)is 7)?(n:y-a;r:1 mk n;mv(r+8;a;n);8::y; :sc r));x}
sym.I:III{((y~z)+~x~96)? :0;y+:1;x:C y;8::y;(y<z)?(r:nam(x;y;z);r? :r;r:chr(x;y;z);r? :sc r);r:5 mk 1;(r+8)::0;r}
sms.I:III{r:sym(x;y;z);(~r)? :r;1/(y:I 8;q:sym(C y;y;z);(~q)? :enl r;r:r cat q);r}
is.I:II{y&cla x}cla.I:I{(128<x-32)? :0;C 128+x}
160!{204840484848485040604848484848504444444444444444444448604848484848424242424242424242424242424242424242424242424242424240506048484041414141414141414141414141414141414141414141414141414048604800}

000:{xxx; gtc; gti; gtf; xxx; gti; gtl; mod; xxx; eqc; eqi; eqf; eqz; eqi; eqL; xxx; abc; abi; abf; abz; nec; nei; nef; nez; xxx; moi; xxx; xxx; sqc; sqi; sqf; sqz}
032:{xxx; mkd; xxx; rsh; cst; diw; min; ecv; ecd; epi; mul; add; cat; sub; cal; ovv; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; asn; xxx; les; eql; mor; fnd}
064:{atx; xxx; xxx; xxx; xxx; xxx; nmf; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; sci; scv; ecl; exc; cut}
096:{xxx; xxx; xxx; xxx; drw; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ovi; max; ecr; mtc; xxx}
128:{xxx; sin; cos; exp; log; xxx; xxx; xxx; chr; nms; vrb; nam; sms; xxx; xxx; xxx; adc; adi; adf; adz; suc; sui; suf; suz; muc; mui; muf; muz; dic; dii; dif; diz}
160:{out; til; xxx; cnt; str; sqr; wer; epv; ech; ecp; fst; abx; enl; neg; val; riv; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; lst; xxx; grd; grp; gdn; unq}
192:{typ; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; scn; liv; spl; srt; flr}
224:{xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; kst; lgf; xxx; xxx; xxx; prs; xxx; rnd; xxx; xxx; xxx; xxx; xxx; xxx; xxx; xxx; ovr; rev; jon; not; xxx}

\
01234567   xt:x>>29      xt~0(function)  x<256(basic) x<128(dyadic)
Fcifzsld   xn:x&^7<<29   xn~2(derived)   adv  verb
4148x444                 xn~3(proj)      verb argv empty-index             
                         xn~4(lambda)    str  tree args arity   
                                              
+  add abx               abs:+z          bytes
-  sub neg                                 0..  7  type sizes   0 1 4 8 16 4 4 0
*  mul fst                                 8.. 11  parse cur (pp)
%  div sqr               conj:%z          12.. 15  rng state
&  min wer   prs flp     ang:&z           16..127  free pointers (4*i) for 4..31
|  max rev                               128..131  memsize log2
<  les grd                               132..135  k-tree keys
>  mor gdn                               136..139  k-tree values
=  eql grp                               140..143  trp line
~  mtc not               ~f              144..147  trp col
!  mkd til   seq         z:re!im im:!z   148..151  `x`y`z
,  cat enl                               152..155  
^  exc asc                               156..159  
$  str cst   sc cs                       160..255  char map az|AZ|NM|VB|AD|TE
#  rsh cnt   take                        256.....  buckets/heap
_  drp flr   drop        re:_z f:__i i:_f
?  fnd unq   fnd fnx                     (:;`x;y)          assign      x:y
@  atx typ               z:abs@ang z@ang (+;(`x;a;b;c);y)  assign(m/i) s[a;b;c]+:y
.  cal val               im:. z          (;a;b;c)   (*128) sequence    a;b;c     
                                         ((/;+);1 2 3)     adverbs     +/1 2 3 
+'x  ech(168)  x+'y  ecd(40)             x'y  bin          ::x(last)   `a`b:1 2
+/x  ovr(251)  x+/y  ovi(123)  whl nlp   x/y  mod
+\x  scn(219)  x+\y  sci(91)  (c;f)/:x   x\y  y%x   \(help)  \\(exit)  \d(dump)
+':x ecp(169)  x+':y epi(41)  (c;f)\:x   x':y win?  dropfile(fs[name]) \lm(ld m.k)
+/:x fxp(253)  x+/:y ecl(125) (n;f)/:x   x/:y join  \wFILE(download)   \w(k.ws)
+\:x fxp(221)  x+\:y ecr(93)  (n;f)\:x   x\:y split \c(clear console)  \e var(edt)
