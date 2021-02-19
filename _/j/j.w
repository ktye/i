j:I:I{
 (~I 0)?ini 16;
 !;0}

ini.V:I{0::x;p:128;i:7;(i<x)?/((4*i)::p;p*:2);4::mk 0;12::mk 0;s:mk 5;(8+s)::mk 0;(12+s)::I s+8}

bk.I:I{r:32-*7+4*x;(r<4)? :4;r}
mk.I:I{t:bk x;i:4*t;m:4*I 0;(~I i)?/((i>='m)?!;i+:4);a:I i;i::I a;k:i-4;(k>=4*t)?/(u:a+1<<k>>2;u::I k;k::u;k-:4);a::1;(a+4)::x;a}

rx.I:I{(~x&7)?(x::1+I x);x}                                                     // ref
dx.V:I{x?(~x&7)?(r:(I x)-1;x::r;(~r)?(n:I x+4;p:x+8;n/(dx I p;p+:4);fr x))}     // unref
fr.V:I{p:4*bk I 4+x;x::I p;p::x}                                                // free
nn.I:I{I 4+x}                                                                   // length

lc.I:II{n:nn x;(~n)?(dx x;r:mk 1;(8+r)::y; :r);(1~I x)?((bk 1+n)~bk n)?((4+lp x)::y;(4+x)::1+I 4+x; :x);r:mk 1+n;xp:8+x;rp:8+r;n/(rp::rx I xp;rp+:4;xp+:4);rp::y;dx x;r}
pc.I:I{s:I 8;p:I 8+s;t:I 12+s;q:pa t;r:lc(t;x);(12+s)::r;(t~p)?((8+s)::r; :r);(lp q)::r;r}
pa.I:I{p:I 8+I 8;1?/((~nn p)? :p;l:I(lp p);((l~x)+(~l)+p~x)? :p;p:l);p}
lp.I:I{4+x+4*nn x}
fi.I:I{(0~nn x)?!;4+x+4*nn x}
use.I:I{(~nn x)?!;r:rx I x+8;dx x;r}
ipo.I:I{x:po(x);(~x&1)?!;x%'2}
lpo.I:I{x:po(x);(~x&7)?!;x}

pi:V:I{pu 1+2*x}
add.V:I{pi(ipo x)+ipo x}           /+
sub.V:I{pi(-ipo x)+ipo x}          /-
mul.V:I{pi(ipo x)*ipo x}           /*
div.V:I{swp x;pi(ipo x)%ipo x}     //
mod.V:I{swp x;pi(ipo x)\'ipo x}    /%
eql.V:I{pi(ipo x)~ipo x}           /=
gti.V:I{pi(ipo x)>'ipo x}          />
lti.V:I{pi(ipo x)<'ipo x}          />

stk.V:I{!}                         /!
dup.V:I{x:po x;pu x;pu x}          /"
cat.V:I{!} /,
cnt.V:I{!} /#
rol.V:I{x:po x;y:po x;z:po x;pu x;pu z;pu y}


swp.V:I{x:po x;y:po x;pu x;pu y} /~
drp.V:I{x:po x}                  /_
amd.V:I{v:po x;i:ipo x;a:use lpo x;n:nn a;$[i~n;a:a lc v;(i<'0)+i>'n;!;(ap:8+a+4*i;x:rx I ap;ap::v)];pu a}                     /$
atx.V:I{i:ipo x;l:lpo x;((i<0)+i>='nn l)?!;pu rx I 8+l+4*i;dx l}                                                               /@
asn.V:I{y:fi lpo x;(~2~y&3)?!;v:po x;(v&7)?(v:(mk 0)lc v);s:I 12;p:s fn y;(~p)?(s:s lc y;s:s lc 1;p:lp s);dx I p;p::v;12::s}   /:
fn:I:II{n:(nn x)>>1;p:x+8;n/((y~I p)? :4+p;p+:8);0}
po.I:I{s:I 4;n:nn s;(~n)?!;p:lp s;r:I(lp s);p::0;n-:1;((bk 1+n)~bk n)?((4+s)::n; :r);q:mk n;qp:q+8;sp:s+8;n/(qp:rx I sp;qp+:4;sp+:4);dx s;4::q;r}
pu.V:I{4::(I 4)lc x}

//33+  !    "    #    $    %    &    '    (    )    *    +    ,    -    .    / 
000:{stk; dup; cnt; amd; mod; min; whl; xxx; xxx; mul; add; cat; sub; exe; div}

//33+  :    ;    <    =    >    ?    @
025:{asn; xxx; lti, eql, gti, ife, atx}

//33+  ^    _   // `
061:{max; drp}

//33+  |    }    ~
091:{rol; xxx; swp}