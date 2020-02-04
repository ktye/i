ini:I:I{0::1130366807310592j;128::x;p:512;i:9;(i<x)?/((4*i)::p;p::i;p*:2;i+:1);x} /x:16(64k)
mk:I:II{t:32-*7+y*C x;i:4*t;(~I i)?/i+:4;(128~i)?!;i::I 4+a:I i;j:i-4;(j>=4*t)?/(j-:4;u:a+1<<j-2;u::I j;j::u);a::y|x<<29;(a+4)::1;a}  /make vector(t,n) (todo grow)
gi:I:I{I x}gj:J:I{J x}gf:F:I{F x}gc:I:I{C x}
/til:I:I{v1;xt?[e;;;;e;e;I xd];xt tox jota I xd}
/jota:I:I{r:2 mk x;x/(8+r+i)::i;r}

\
01234567   xt:x>>29       xn:x&536870911 (-1+1<<29)
Fcifzsld   xt=0(function) x<256(basic) 
0148x440   ft=xn&0xff00  (derived, proj, comp, lambda, native)     
	   fn=xn&0xff    (argn)

+ add abs                 memory
- sub neg                   0..  7   type sizes
* mul fst                   8.. 11   k-tree/key   pointer
% div sqr                  12.. 15   k-tree/value pointer
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
