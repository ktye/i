ini:I:I{0::1130366807310592j;8::x;p:512;i:9;(i<x)?/((4*i)::p;p::i;p*:2;i+:1);x} /x:16(64k)
bk:I:II{32-*7+y*C x}                                                            /bucket type for vector t,n sz:I:II{8+y*C x}  
mk:I:II{t:x bk y;i:4*t;(~I i)?/i+:4;i::I 4+a:I i;a::y|x<<29;(a+4)::1;a}         /todo grow, split
gi:I:I{I x}gj:J:I{J x}gf:F:I{F x}gc:I:I{C x}
/til:I:I{v1;xt?[e;;;;e;e;I xd];xt tox jota I xd}
/jota:I:I{r:2 mk x;x/(8+r+i)::i;r}

\
01234567   xt:x>>29       xn:x&536870911 (-1+1<<29)
Fcifzsld   xt=0(function) x<256(basic) 
0148x440   ft=xn&0xff00  (derived, proj, comp, lambda, native)     
	   fn=xn&0xff    (argn)

+ add abs
- sub neg
* mul fst
% div sqr
& min wer
| max rev
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
