/        64     72    80    88     96     104
`k`l`a`b`while`"rf."`"rz."`"gdt."`"grp."`"enc."

`x:,/+"0123456789abcdef"@16 16\256!            /`x@ hex
`t:`39         /`t@ token
`p:`46         /`p@ parse

`enc:{$[#y;+(&','(|/c)-c:#'r),'r:{x\y}/[x;y];(#x)#,!0]}  /x\Y
`gdt:{[t;g]($[g;{x@>y x};{x@<y x}])/(,!#t),|.t}  /<t >t  grade table

`grp:{(x@*'g)!g:(&~a~'a@-1+!#a:x i)^i:<x}      /=x

abs:`32;sin:`44;cos:`45;find:`31;imag:`33;conj:`34;angle:`35;exp:`42;log:`43

/ pretty print  `l@x  `k@x
`pad:{x@\!|/#'x}

`lxy:{ /k
 kt:{[x;y;k;T]x:$[`T~@x;T[x;k];`pad("";"-"),$x];(x,'"|"),'T[y;k]}
 d:{[x;k;kt;T]r:!x;x:.x;$[`T~@x;kt[r;x;k;T];,'[,'[`pad(k'r);"|"];k'x]]}
 T:{[x;k]$[`L?@'.x;,k x;(,*x),(,(#*x)#"-"),1_x:" "/'+`pad@'$(!x),'.x]}
 t:@y;k:`kxy@*x;h:*|x
 dd:("";,"..")h<#y:$[(@y)?`L`D`T;y;y~*y;y;[t:`L;,y]]
 y:$[y~*y;y;(h&#y)#y]
 $[`D~t;d[y;k;kt;T];`T~t;T[y;k];y~*y;,k y;k'y],dd}

`l:`lxy 70 20

`str:{q:{c,("\\"/(0,i)^@[x;i;(qs!"tnr\"\\")x i:&x?\qs:"\t\n\r\"\\"]),c:_34}
 $[|/x?\"\t\n\r"__!31;"0x",`x@x;q x]}

`kxy:{ /k m t n
 a:{t:@x;x:$x;$[`c~t;`str x;`s~t;"`",x;x]}
 d:{[x;k]r:"!",k@.x;n:#!x;x:k@!x;$[(n<2)|(@.x)?`D`T;"(",x,")";x],r}
 v:{[x;k;m]t:@x;x:(m&n:#x)#x
  x:$[`L~t;k'x;`C~t;x;$x]
  x:$[`C~t;`str x;`S~t;c,(c:"`")/x;`L~t;$[1~n;*x;"(",(";"/x),")"];" "/x]
  $[m<#x:(("";",")(1~n)),x;((m-2)#x),"..";x]}
 t:@y;k:`kxy x
 $[`T~t;"+",d[+y;k];`D~t;d[y;k];0~#y;(`C`I`S`F`Z`L!("\"\"";"!0";"0#`";"0#0.";"0@0a";"()"))t;y~*y;a y;v[y;k;x]]} /344
 
`k:`kxy 1000000

/random
`rf: {.5+(x?0)%4294967295.}           /?n (uniform)  e.g. ?100
`rf1:{.5+(1.+x?0)%4294967295.}        
`rz: {(%-2*log `rf1 x)@360.*`rf x}    /?-n (normal)  ?z (binormal)   ?-100   ?100a
