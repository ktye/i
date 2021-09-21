/               64      72    80    88     96     104    112    120                          
``x`y`z`k`l`t`p`while`".rf"`".rz"`".uqs"`".uqf"`".gdt"`".lin"`".odo"
`".x"`".y"`".z"`".k"`".l"          /128.. (index.go:Atx)
`".x":{,/+"0123456789abcdef"@(x%16;16/x:0+x)} /`x@ (hex)
`".y":{x}                          /132
`".z":{`Z@x}                       /140
`".k":{(.`".kxy")[x;1000000]}      /148

`".l":{t:@x;pad:{(|/#'x)#'x};lst:.`".l";k:(.`".kxy")[;1]
 kt:{x:$[`T~@x;T x;pad ("";"-"),$x];(x,'"|"),'T y}
 d:{r:!x;x:.x;$[`T~@x;kt[r;x];,'[,'[pad(k'r);"|"];k'x]]}
 T:{$[`L':@'.x;,k x;(,*x),(,(#*x)#"-"),1_x:" "/:'+pad'$(!x),'.x]}
 dd:("";,"..")20<#x
 x:$[x~*x;x;(20&#x)#x]
 $[`L~t;k'x;`D~t;d x;`T~t;T x;,k x],dd} /352

`".kxy":{t:@x;n:#x;k:(.`".kxy")[;y];m:y
 q:{c,("\\"/:(0,i)^@[x;i;(qs!"tnr\"\\")x i:&x':qs:"\t\n\r\"\\"]),c:_34}
 s:{$[|/x':"\t\n\r"__!31;"0x",`x@x;q x]}
 /s:{q x}
 a:{x:$x;$[`c~t;s x;`s~t;"`",x;x]}
 d:{r:"!",k@.x;n:#!x;x:k@!x;$[(1~n)|(@.x)':`D`T;"(",x,")";x],r}
 v:{m*:(.`".kstm")t;  dd:("";"..")m<#x;x:(m&#x)#x
  x:$[`L~t;k'x;`C~t;x;$x]
  x:$[`B~t;(*'x),"b";`C~t;s x;`S~t;c,(c:"`")/:x;`L~t;$[1~n;*x;"(",(";"/:x),")"];" "/:x]
  (("";",")(1~n)),x,dd}
 $[n~0;(.`".kst0")@t;`T~t;"+",d@+x;`D~t;d x;x~*x;a x;v x]} /344

`".uqs":{x@&~0b~':x:^x} 
`".uqf":{x@&(!#x)=x?x}
`".gdt":{[t;g](!#t){x g y x}/|.t}

`".kst0":`B`C`I`S`F`Z`L!("0#0b";c,c:_34;"!0";"0#`";"0#0.";"0#0a";"()")
`".kstm":`B`C`I`S`F`Z`L!100 100 30 30 20 10 20

`".lin":{$[`L~@z;(.`".lin")[x;y]'z;[dx:0.+1_-':x;dy:0.+1_-':y;b:(-2+#x)&0|x'z;(y b)+(dy b)*(z-x b)%dx b]]}
`".odo":{{x/y}'[x;(!*/x)%/:{(*/x)%\ x}x]}

dot:{dotmv:{{+/x*y}\:[x;y]};dotmv[x;y]}

solve:{qslv:{H:x 0;r:x 1;n:x 2;m:x 3;j:0;K:!m
 while[j<n;y[K]-:(+/(conj H[j;K])*y K)*H[j;K];K:1_K;j+:1]
 i:n-1;J:!n;y[i]%:r@i
 while[i;j:i_J;i-:1;y[i]:(y[i]-+/H[j;i]*y@j)%r@i]
 n#y}
 q:$[`i~@*|x;x;qr x];$[`L~@y;qslv/:[q;y];qslv[q;y]]}

qr:{K:!m:#*x;I:!n:#x;j:0;r:n#0a;turn:$[`Z~@*x;{(-x)@angle y};{x*1. -1@y>0}]
 while[j<n;I:1_I
  r[j]:turn[s:abs/j_x j;xx:x[j;j]]
  x[j;j]-:r[j]
  x[j;K]%:%s*(s+abs xx)
  x[I;K]-:{+/x*y}/:[(conj x[j;K]);x[I;K]]*\:x[j;K]
  K:1_K;j+:1];(x;r;n;m)}

any:`30;abs:`32;sin:`44;cos:`39;find:`31;fill:`38;imag:`33;conj:`34;angle:`35;exp:`42;log:`43

ej:{(y j),'x_z i j:&~0N=i:(z x)?y x} /sym t1 t2
avg:{(+/x)%0.+#x}
var:{(+/x*x:(x-avg x))%-1+#x}
std:{%var x}

`".rf": {.5+(x?0)%4294967295.}
`".rf1":{.5+(1.+x?0)%4294967295.}
`".rz": {(%-2*log(.`".rf1")x)@360.*(.`".rf")x}
