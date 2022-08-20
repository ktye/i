/               64     72    80    88     96     104    112    120    128                 
``x`y`z`k`l`a`b`while`"rf."`"rz."`"uqs."`"uqf."`"gdt."`"lin."`"odo."`"grp." /

`"x.":{,/+"0123456789abcdef"@(x%16;16/x:256/256+x)} /`x@ (hex)
`"t.":`45         /`t@ token
`"p.":`46         /`p@ parse
`"b.":(`46)[`b;]  /`b@ reinterpret
`"c.":(`46)[`c;]
`"i.":(`46)[`i;]
`"s.":(`46)[`s;]
`"f.":(`46)[`f;]
`"z.":(`46)[`z;]

`"uqs.":{x@&~0b~':x:^x}                           /?^x
`"uqf.":{x@&(!#x)=x?x}                            /?xL
`"gdt.":{[t;g](!#t)($[g;{x@>y x};{x@<y x}])/|.t}  /<t >t  grade table

`"odo.":{{y@(#y)/!x}/:[*/x;&'x#'|*\-1_1,|x]} /!I (odometer)
`"grp.":{(x@*'g)!g:(&~x~':x i)^i:<x}         /=x

any:`30;abs:`32;sin:`44;cos:`39;find:`31;fill:`38;imag:`33;conj:`34;angle:`35;exp:`42;log:`43


/ pretty print  `l@x  `k@x
`"pad.":{(|/#'x)#'x}

`"l.":{ /k
 kt:{[x;y;k;T]x:$[`T~@x;T[x;k];`pad("";"-"),$x];(x,'"|"),'T[y;k]}
 d:{[x;k;kt;T]r:!x;x:.x;$[`T~@x;kt[r;x;k;T];,'[,'[`pad(k'r);"|"];k'x]]}
 T:{[x;k]$[`L':@'.x;,k x;(,*x),(,(#*x)#"-"),1_x:" "/:'+`pad@'$(!x),'.x]}
 t:@x;k:`kxy 1
 dd:("";,"..")20<#x:$[(@x)':`L`D`T;x;x~*x;x;[t:`L;,x]]
 x:$[x~*x;x;(20&#x)#x]
 $[`D~t;d[x;k;kt;T];`T~t;T[x;k];x~*x;,k x;k'x],dd} /352

`"str.":{q:{c,("\\"/:(0,i)^@[x;i;(qs!"tnr\"\\")x i:&x':qs:"\t\n\r\"\\"]),c:_34}
 $[|/x':"\t\n\r"__!31;"0x",`x@x;q x]}

`"kxy.":{ /k m t n
 a:{t:@x;x:$x;$[`c~t;`str x;`s~t;"`",x;x]}
 d:{[x;k]r:"!",k@.x;n:#!x;x:k@!x;$[(1~n)|(@.x)':`D`T;"(",x,")";x],r}
 v:{[x;k;m;n]m*:(.`".kstm")t:@x; dd:("";"..")m<#x;x:(m&#x)#x
  x:$[`L~t;k'x;`C~t;x;$x]
  x:$[`B~t;(*'x),"b";`C~t;`str x;`S~t;c,(c:"`")/:x;`L~t;$[1~n;*x;"(",(";"/:x),")"];" "/:x]
  (("";",")(1~n)),x,dd}
 t:@y;n:#y;k:`kxy x;m:x
 $[n~0;(.`".kst0")@t;`T~t;"+",d[+y;k];`D~t;d[y;k];y~*y;a y;v[y;k;m;n]]} /344

`".kst0":`B`C`I`S`F`Z`L!("0#0b";c,c:_34;"!0";"0#`";"0#0.";"0#0a";"()")
`".kstm":`B`C`I`S`F`Z`L!100 100 30 30 20 10 20
 
`"k.":`kxy 1000000

/random
`"rf.": {.5+(x?0)%4294967295.}           /?n (uniform)  e.g. ?100
`"rf1.":{.5+(1.+x?0)%4294967295.}        
`"rz.": {(%-2*log `rf1 x)@360.*`rf x}    /?-n (normal)  ?z (binormal)   ?-100   ?100a

/interpolate: x y'z
`"lin.":{$[`L~@z;(.`"lin.")[x;y]'z;[dx:0.+1_-':x;dy:0.+1_-':y;b:(-2+#x)&0|x'z;(y b)+(dy b)*(z-x b)%dx b]]}
