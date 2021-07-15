/32 40 48   56 64   72  80   88   96    104   112                       200   208    216    224
 `k`l`while`in`find`abs`imag`conj`angle`solve`dot`a`b`c`d`e`f`g`h`i`j`".kst"`".qr"`".slv"`".dot"
 
`".kst":{t:@x;n:#x;k:.(`".kst")
 q:{c,x,c:_34}
 a:{x:$x;$[`c~t;q x;`s~t;"`",x;x]}
 d:{r:"!",k@.x;x:k@!x;$[1~n;"(",x,")";x],r}
 v:{x:$[`L~t;k'x;t in `D`C;x;$x]
  x:$[`B~t;(*'x),"b";`C~t;q x;`S~t;c,(c:"`")/:x;`L~t;$[1~n;*x;"(",(";"/:x),")"];`D~t;d x;" "/:x]
  (("";",")(1~n)&~`D~t),x}
 $[n~0;(.`".kst0")@t;`T~t;"+",k@+x;x~*x;a x;v x]}
`".kst0":`B`C`I`S`F`Z`L!("0#0b";c,c:_34;"!0";"0#`";"0#0.";"0#0a";"()")

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
