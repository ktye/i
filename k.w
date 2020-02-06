ini:I:I{0::1130366807310592j;128::x;p:512;i:9;(i<x)?/((4*i)::p;p::i;p*:2;i+:1);x} /x:16(64k)
mk:I:II{t:x bk y;i:4*t;(~I i)?/i+:4;(128~i)?!;i::I 4+a:I i;j:i-4;(j>=4*t)?/(j-:4;u:a+1<<j%4;u::I j;j::u);a::y|x<<29;(a+4)::1;a}
fr:0:I{v1;t:4*xt bk xn;x::I t;t::x}bk:I:II{32-*7+y*C x}
mki:I:I{r:2 mk 1;(r+4)::1;(r+8)::x;r}
decr:0:I{(x>255)?(xr:I x+4;(x+4)::xr-1;(1~xr)?fr x)}dxr:{decr x;r}
v1:{xt:(I x)>>29;xn:(I x)&536870911;xp:8+x;(4~xt)?xp:16+x}
til:I:I{v1;(~2~xt)?!;r:xt mk n:I xp;rp:8+r;n/(rp::i;rp+:4);dxr}
fst:I:I{v1;(~xt)?:x;(7~xt)?:x;r:xt mk 1;xt?[!;(r+8)::C xp;;(r+8)::J xp;((r+8)::J xp;(r+16)::J xp+8);;;!;(r+8)::I xp];dxr}

\
01234567   xt:x>>29       xn:x&536870911 (-1+1<<29)
Fcifzsld   xt~0(function) x<256(basic) 
0148x440   ft~xn&0xff00  (derived, proj, lambda, native)  composition==lambda?
	   fn~xn&0xff    (argn)

+ add abs                 memory
- sub neg                   0..  7   type sizes   0 1 4 8 16 4 4 0
* mul fst                   8.. 11   k-tree/key   pointer     todo
% div sqr                  12.. 15   k-tree/value pointer     todo
& min wer                  16..127   free pointers (4*i) for bt i, i:4..31
| max rev                 128..131   memsize log2
< les gup                 
> mor gdn                 
= eql grp                 
~ mtc not                 
! key til                 
, cat enl                 
^ exc asc                 
$ str cnt                 
_ drp flr                 
? fnd unq                 
@ atx typ                 
. cal val                 
