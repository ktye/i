j:I:I{                                                          // single entry point
 (~I 0)?ini 16
 ...
}

ini.V:I:{0::x;p:128;i:7;(i<x)?/((4*i)::p;p*:2);8::mk 0;12:mk 0} // todo P:mk 0;T:P
bk.I:I{r:32-*7+4*x;(r<4)? :4;r}                                 // bucket type .. mk allocator
mk.I:I{t:bk x;i:4*t;m:4*I 0;(~I i)?/((i>='m)?!;i+:4);a:I i;i::I a;j:i-4;(j>=4*t)?/(u:a+1<<j>>2;u::I j;j::u;j-:4);a::1;(a+4)::x}

rx.I:I{(~x&7)?(x::1+I x);x}                                                     // ref
dx.I:I{(x&~x&7)?( x::(I x)-1; (~I x)?(n:I x+4;p:x+8;n/(dx I p;p+:4);fr x);x)}   // unref
fr.I:I{p:4*bk I 4+x;x::I p;p::x}                                                // free
nn.I:I{I 4+x}                                                                   // length

//33+  !    "    #    $    %    &    '    (    )    *    +    ,    -    .    / 
000:{stk; dup; cnt; nyi; mod; min; whl; xxx; xxx; mul; add; cat; sub; exe; div}

//33+  :    ;    <    =    >    ?    @
025:{asn; xxx; lti, eql, gti, ife, atx}

//33+  ^    _   // `
061:{max; pop}

//33+  |    }    ~
091:{rol; xxx; swp}