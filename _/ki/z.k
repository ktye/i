/        64     72    80    88     96     104    112    120
`k`l`a`b`while`"rf."`"rz."`"uqs."`"uqf."`"gdt."`"lin."`"grp."`"enc."

`x:,/+"0123456789abcdef"@16 16\256!            /`x@ hex
`t:`39         /`t@ token
`p:`46         /`p@ parse

`enc:{$[#y;+(&'(|/c)-c:#'r),'r:{x\y}/[x;y];(#x)#,!0]}  /x\Y
`uqs:{x@&1,1_~x~'x@-1+!#x:^x}                  /?^x
`uqf:{x@&(!#x)=x?x}                            /?xL
`gdt:{[t;g]($[g;{x@>y x};{x@<y x}])/(,!#t),|.t}  /<t >t  grade table

`grp:{(x@*'g)!g:(&~a~'a@-1+!#a:x i)^i:<x}      /=x

abs:`32;find:`31

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
 $[`T~t;"+",d[+y;k];`D~t;d[y;k];0~#y;(`C`I`S`L!("\"\"";"!0";"0#`";"()"))t;y~*y;a y;v[y;k;x]]}
 
`k:`kxy 1000000
`d:{x-(*x),-1_x}  /deltas

