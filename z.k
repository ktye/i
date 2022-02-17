/               64     72    80    88     96     104    112    120    128                 
``x`y`z`k`l`a`b`while`"rf."`"rz."`"uqs."`"uqf."`"gdt."`"lin."`"odo."`"grp."

`"x.":{,/+"0123456789abcdef"@(x%16;16/x:256/256+x)} /`x@ (hex)
`"t.":`45         /`t@ token
`"p.":`46         /`p@ parse
`"b.":(`46)[`b;]  /`b@ reinterpret
`"c.":(`46)[`c;]
`"i.":(`46)[`i;]
`"s.":(`46)[`s;]
`"f.":(`46)[`f;]
`"z.":(`46)[`z;]

`"l.":{t:@x;k:`kxy 1
 kt:{x:$[`T~@x;T x;`pad("";"-"),$x];(x,'"|"),'T y}
 d:{r:!x;x:.x;$[`T~@x;kt[r;x];,'[,'[`pad(k'r);"|"];k'x]]}
 T:{$[`L':@'.x;,k x;(,*x),(,(#*x)#"-"),1_x:" "/:'+`pad@'$(!x),'.x]}
 dd:("";,"..")20<#x:$[(@x)':`L`D`T;x;x~*x;x;[t:`L;,x]]
 x:$[x~*x;x;(20&#x)#x]
 $[`D~t;d x;`T~t;T x;x~*x;,k x;k'x],dd} /352

`"pad.":{(|/#'x)#'x}
`"split.":{$[`L~@x;`split@'x;" "\:$[" "=x@-1+#x:x@&~i&~':i:" "=x;-1_x;x]]} /multiple whitespace/leading/trailing
`"edit.":{et:{"+",(`k@!x),"!",el@.x}
 el:{t:(`S`B`C`I`F`Z!``b`c`i`f`z)@@'x; (`k@t),"$'+`split@'\"\\n\"\\:-1_1_\"\n",("\n"/:" "/:'+`pad@'$x),"\n\""}
 $[`T~t:@x;et x;`L~t;el x;`C~t;"-1_1_\"\n",x,"\n\"";"*",el@,x]}

`"kxy.":{t:@y;n:#y;k:`kxy x;m:x
 q:{c,("\\"/:(0,i)^@[x;i;(qs!"tnr\"\\")x i:&x':qs:"\t\n\r\"\\"]),c:_34}
 s:{$[|/x':"\t\n\r"__!31;"0x",`x@x;q x]}
 a:{x:$x;$[`c~t;s x;`s~t;"`",x;x]}
 d:{r:"!",k@.x;n:#!x;x:k@!x;$[(1~n)|(@.x)':`D`T;"(",x,")";x],r}
 v:{m*:(.`".kstm")t; dd:("";"..")m<#x;x:(m&#x)#x
  x:$[`L~t;k'x;`C~t;x;$x]
  x:$[`B~t;(*'x),"b";`C~t;s x;`S~t;c,(c:"`")/:x;`L~t;$[1~n;*x;"(",(";"/:x),")"];" "/:x]
  (("";",")(1~n)),x,dd}
 $[n~0;(.`".kst0")@t;`T~t;"+",d@+y;`D~t;d y;y~*y;a y;v y]} /344
 
`"k.":`kxy 1000000

`"uqs.":{x@&~0b~':x:^x} 
`"uqf.":{x@&(!#x)=x?x}
`"gdt.":{[t;g](!#t){x g y x}/|.t}

`".kst0":`B`C`I`S`F`Z`L!("0#0b";c,c:_34;"!0";"0#`";"0#0.";"0#0a";"()")
`".kstm":`B`C`I`S`F`Z`L!100 100 30 30 20 10 20

`"lin.":{$[`L~@z;(.`"lin.")[x;y]'z;[dx:0.+1_-':x;dy:0.+1_-':y;b:(-2+#x)&0|x'z;(y b)+(dy b)*(z-x b)%dx b]]}
`"odo.":{{x/y}'[x;(!*/x)%/:{(*/x)%\ x}x]}
`"grp.":{(x@*'g)!g:(&~x~':x i)^i:<x}

dot:{dotmv:{{+/x*y}\:[x;y]};dotmv[x;y]}

solve:{qslv:{H:x 0;r:x 1;n:x 2;m:x 3;j:0;K:!m
 while[j<n;y[K]-:(+/(conj H[j;K])*y K)*H[j;K];K:1_K;j+:1]
 i:n-1;J:!n;y[i]%:r@i
 while[i;j:i_J;i-:1;y[i]:(y[i]-+/H[j;i]*y@j)%r@i]
 n#y}
 q:$[`i~@*|x;x;qr x];$[`L~@y;qslv/:[q;y];qslv[q;y]]}

qr:{K:!m:#*x;I:!n:#x;j:0;r:n#0a;turn:$[`Z~@*x;{(-x)@angle y};{x*1. -1@y>0}]
 while[j<n;I:1_I
  r[j]:turn[s:abs@abs/j_x j;xx:x[j;j]]
  x[j;j]-:r[j]
  x[j;K]%:%s*(s+abs xx)
  x[I;K]-:{+/x*y}/:[(conj x[j;K]);x[I;K]]*\:x[j;K]
  K:1_K;j+:1];(x;r;n;m)}

any:`30;abs:`32;sin:`44;cos:`39;find:`31;fill:`38;imag:`33;conj:`34;angle:`35;exp:`42;log:`43

ej:{(y j),'x_z i j:&~0N=i:(z x)?y x} /sym t1 t2
avg:{(+/x)%0.+#x}
var:{(+/x*x:(x-avg x))%-1+#x}
std:{%var x}
rem:{x/x+x/y} /like apl for negative arguments

`"rf.": {.5+(x?0)%4294967295.}
`"rf1.":{.5+(1.+x?0)%4294967295.}
`"rz.": {(%-2*log `rf1 x)@360.*`rf x}
`".html":{$[`L~@x;"<div style='display:flex;flex-direction:column'>",(,/(.`".html")'x),"</div>";`S~@x;"<div style='display:flex;flex-direction:row'>",(,/(.`".html")'x),"</div>";~`s~@x;"";{{"<",y," id='",x,"'></",y,">"}[x;$[y':`i`c`f`z`s`I`B`C`F`Z`S;"input";`b~y;"input type='checkbox'";`T~y;"table";`l~t;"button";"pre"]]}[$x;@.x]]}

`"pack.":{w:{(`c@,#x),x};($t),$[`s~t:@x;`pack@$x;x~*x;w `c@,x;`L~@x;(`c@,#x),,/`pack@'x;(@x)':`D`T;(`pack@.x),`pack@!x;`S~t;,/`pack@$x;w `c x]}
`"unpack.":{s:x;g:{[n]r:n#s;s::n_s;r};n:{*`i@g 4};u:{x;$[(t:*g 1)':"bcifz";*(`$t)g n[];t~"s";`$u 0;t~"S";`$u 0;t~"L";u'!n[];t~"D";(u 0)!u 0;t~"T";+(u 0)!u 0;(`$_t+32)g n[]]};u 0}

csv:{c:{s:`$'x@i:&x':"ifzs";n:`i$" "\:-1_@[x;i;" "];y[a]:(y[a],''"a"),''y[1+a:&s=`z];s$'y n};s:$[" "~(*x);`split@;(*x)\:];x:1_x;y:+s'$[`L~@y;y;"\n"\:y];$[#x;c[x;y];y]} / csv[";1z3i2f"; "input.."]

`"tics.":{[minmax]nice:{[x;r]f:x%0.+10^ex:_log[10;x];(1 2 5 10.@*&(~f>1 2 5 0w;f<1.5 3 7 0w)[r])*10^ex};e:nice[-/|minmax;0];s:nice[e%4.;1];n:_1.5+e%s;$[~(minmax 1)>*-2#r:s*(_(*minmax)%s)+!n;-1_r;r]}

`"ceg.":{(x i)!0-':1+i:&(1_~~':x),1b} /#'= (^x)
hist:{$[`i~@x;hist[(x;&/y;|/y);y];(Y;(`38)[0;(`ceg@^1+((d%2.0)+-1_Y:(x 1)+(d:(--/1_x)%-1.+n)*!n)'y)@!n:_0.+*x])]} /[xn;y] [xn xmin xmax;y] e.g. plot `x`y!hist[10 -4 4.;?-1000]
