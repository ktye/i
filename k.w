ini:I:I{8::x;p:256;i:11;(i<16)?/(p*:2;(4*i)::p;(4*p)::i;i+:1);x}      /x:16(64k)
st:I:I {I?255j&1130366807310592j>>J?8*x}                              /element size for type
sz:I:II{8+y*st x}                                                     /storage size for vector of t,n
bk:I:II{*-1+x sz y}                                                   /bucket type for vector t,n
mk:I:II{t:x bk y;i:4*t;(~I i)?/i+:4;i::a:I 4+i;a::n|x<<29;(a+4)::1;a} /todo grow, split

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
