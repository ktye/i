j:I:I{                                                          // single entry point
 (~I 0)?ini 16
 ...
}

ini.V:I{0::x;p:128;i:7;(i<x)?/((4*i)::p;p*:2);4::mk 0;12::mk 0;s:mk 5;(8+s)::mk 0;(12+s)::I s+8}

bk.I:I{r:32-*7+4*x;(r<4)? :4;r}                                 // bucket type .. mk allocator
mk.I:I{t:bk x;i:4*t;m:4*I 0;(~I i)?/((i>='m)?!;i+:4);a:I i;i::I a;j:i-4;(j>=4*t)?/(u:a+1<<j>>2;u::I j;j::u;j-:4);a::1;(a+4)::x}

rx.I:I{(~x&7)?(x::1+I x);x}                                                     // ref
dx.I:I{(x&~x&7)?( x::(I x)-1; (~I x)?(n:I x+4;p:x+8;n/(dx I p;p+:4);fr x);x)}   // unref
fr.I:I{p:4*bk I 4+x;x::I p;p::x}                                                // free
nn.I:I{I 4+x}                                                                   // length

pc.I:I{s:I 8;p:I 8+s;t:I 12+s;q:pa t;r:lc(t;x);(12+s)::r;(t~p)?((8+s::r); :r);(lp q)::r;r}
pa.I:I{p:I 8+I 8;1?/((~nn p)? :p;l:last p;((l~x)+(~l)+p~x)? :p;p:l);p}
lp.I:I{4+x+4*nn x}

//33+  !    "    #    $    %    &    '    (    )    *    +    ,    -    .    / 
000:{stk; dup; cnt; nyi; mod; min; whl; xxx; xxx; mul; add; cat; sub; exe; div}

//33+  :    ;    <    =    >    ?    @
025:{asn; xxx; lti, eql, gti, ife, atx}

//33+  ^    _   // `
061:{max; pop}

//33+  |    }    ~
091:{rol; xxx; swp}