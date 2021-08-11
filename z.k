/               64   72  80   88   96    104   112 120 128 136  144 152  160 168 176 184          216    224    232   240    248   256            312
``x`y`z`k`l`t`p`while`in`find`abs`imag`conj`angle`solve`dot`sin`cos`".("`".)"`exp`log`fill`a`b`c`".qr"`".slv"`".dot"`".rf"`".rf1"`".rz"`f`g`h`i`j`m`n
`".x"`".y"`".z"`".k"`".l" /320..352

`".x":{,/+"0123456789abcdef"@(x%16;16/x:0+x)} /`x@ (hex) 320(index.go Atx)
`".y":{x}    /328
`".z":{`Z@x} /336
`".k":{(.`".kxy")[x;1000000]} /344

`".l":{t:@x;pad:{(|/#'x)#'x};lst:.`".l";k:(.`".kxy")[;1]
 kt:{x:$[`T~@x;T x;pad ("";"-"),$x];(x,'"|"),'T y}
 d:{r:!x;x:.x;$[`T~@x;kt[r;x];,'[,'[pad(k'r);"|"];k'x]]}
 T:{(,*x),(,(#*x)#"-"),1_x:" "/:'+pad'$(!x),'.x}
 dd:("";,"..")20<#x
 x:$[x~*x;x;(20&#x)#x]
 $[`L~t;k'x;`D~t;d x;`T~t;T x;,k x],dd} /352

`".kxy":{t:@x;n:#x;k:(.`".kxy")[;y];m:y
 q:{c,("\\"/:(0,i)^@[x;i;(qs!"tnr\"\\")x i:&x in qs:"\t\n\r\"\\"]),c:_34}
 s:{$[|/x in "\t\n\r"__!31;"0x",`x@x;q x]}
 /s:{q x}
 a:{x:$x;$[`c~t;s x;`s~t;"`",x;x]}
 d:{r:"!",k@.x;n:#!x;x:k@!x;$[(1~n)|(@.x)in`D`T;"(",x,")";x],r}
 v:{m*:(.`".kstm")t;  dd:("";"..")m<#x;x:(m&#x)#x
  x:$[`L~t;k'x;`C~t;x;$x]
  x:$[`B~t;(*'x),"b";`C~t;s x;`S~t;c,(c:"`")/:x;`L~t;$[1~n;*x;"(",(";"/:x),")"];" "/:x]
  (("";",")(1~n)),x,dd}
 $[n~0;(.`".kst0")@t;`T~t;"+",d@+x;`D~t;d x;x~*x;a x;v x]} /344

 
`".kst0":`B`C`I`S`F`Z`L!("0#0b";c,c:_34;"!0";"0#`";"0#0.";"0#0a";"()")
`".kstm":`B`C`I`S`F`Z`L!100 100 30 30 20 10 20

`".dotmv":{{+/x*y}\:[x;y]}
`".dot":{f:.`".dotmv";f[x;y]}

`".slv":{q:$[`i~@*|x;x;(.`".qr")@x];s:.`".qslv";$[`L~@y;s/:[q;y];s[q;y]]}

`".qslv":{H:x 0;r:x 1;n:x 2;m:x 3;j:0;K:!m
 while[j<n;y[K]-:(+/(conj H[j;K])*y K)*H[j;K];K:1_K;j+:1]
 i:n-1;J:!n;y[i]%:r@i
 while[i;j:i_J;i-:1;y[i]:(y[i]-+/H[j;i]*y@j)%r@i]
 n#y}

`".qr":{K:!m:#*x;I:!n:#x;j:0;r:n#0a;turn:$[`Z~@*x;{(-x)angle angle y};{x*1. -1@y>0}]
 while[j<n;I:1_I
  r[j]:turn[s:0. abs/j_x j;xx:x[j;j]]
  x[j;j]-:r[j]
  x[j;K]%:%s*(s+abs xx)
  x[I;K]-:{+/x*y}/:[(conj x[j;K]);x[I;K]]*\:x[j;K]
  K:1_K;j+:1];(x;r;n;m)}

ej:{(y j),'x_z i j:&~0N=i:(z x)?y x} /sym t1 t2
avg:{(+/x)%0.+#x}
var:{(+/x*x:(x-avg x))%-1+#x}
std:{%var x}

`".rf": {.5+(x?0)%4294967295.}
`".rf1":{.5+(1.+x?0)%4294967295.}
`".rz": {(%-2*log(.`".rf1")x)angle 360*(.`".rf")x}
