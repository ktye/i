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

/ `"l.":{t:@x;k:`kxy 1
/  kt:{x:$[`T~@x;T x;`pad("";"-"),$x];(x,'"|"),'T y}
/  d:{r:!x;x:.x;$[`T~@x;kt[r;x];,'[,'[`pad(k'r);"|"];k'x]]}
/  T:{$[`L':@'.x;,k x;(,*x),(,(#*x)#"-"),1_x:" "/:'+`pad@'$(!x),'.x]}
/  dd:("";,"..")20<#x:$[(@x)':`L`D`T;x;x~*x;x;[t:`L;,x]]
/  x:$[x~*x;x;(20&#x)#x]
/  $[`D~t;d x;`T~t;T x;x~*x;,k x;k'x],dd} /352

`"l.":{ /k
 kt:{[x;y;k;T]x:$[`T~@x;T[x;k];`pad("";"-"),$x];(x,'"|"),'T[y;k]}
 d:{[x;k;kt;T]r:!x;x:.x;$[`T~@x;kt[r;x;k;T];,'[,'[`pad(k'r);"|"];k'x]]}
 T:{[x;k]$[`L':@'.x;,k x;(,*x),(,(#*x)#"-"),1_x:" "/:'+`pad@'$(!x),'.x]}
 t:@x;k:`kxy 1
 dd:("";,"..")20<#x:$[(@x)':`L`D`T;x;x~*x;x;[t:`L;,x]]
 x:$[x~*x;x;(20&#x)#x]
 $[`D~t;d[x;k;kt;T];`T~t;T[x;k];x~*x;,k x;k'x],dd} /352

/ `"kxy.":{t:@y;n:#y;k:`kxy x;m:x
/  q:{c,("\\"/:(0,i)^@[x;i;(qs!"tnr\"\\")x i:&x':qs:"\t\n\r\"\\"]),c:_34}
/  s:{$[|/x':"\t\n\r"__!31;"0x",`x@x;q x]}
/  a:{x:$x;$[`c~t;s x;`s~t;"`",x;x]}
/  d:{r:"!",k@.x;n:#!x;x:k@!x;$[(1~n)|(@.x)':`D`T;"(",x,")";x],r}
/  v:{m*:(.`".kstm")t; dd:("";"..")m<#x;x:(m&#x)#x
/   x:$[`L~t;k'x;`C~t;x;$x]
/   x:$[`B~t;(*'x),"b";`C~t;s x;`S~t;c,(c:"`")/:x;`L~t;$[1~n;*x;"(",(";"/:x),")"];" "/:x]
/   (("";",")(1~n)),x,dd}
/  $[n~0;(.`".kst0")@t;`T~t;"+",d@+y;`D~t;d y;y~*y;a y;v y]} /344

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

`"split.":{$[`L~@x;`split@'x;" "\:$[" "=x@-1+#x:x@&~i&~':i:" "=x;-1_x;x]]} /multiple whitespace/leading/trailing
`"edit.":{et:{"+",(`k@!x),"!",el@.x}
 el:{t:(`S`B`C`I`F`Z!``b`c`i`f`z)@@'x; (`k@t),"$'+`split@'\"\\n\"\\:-1_1_\"\n",("\n"/:" "/:'+`pad@'$x),"\n\""}
 $[`T~t:@x;et x;`L~t;el x;`C~t;"-1_1_\"\n",x,"\n\"";"*",el@,x]}

/linear algebra
dot:{[xt;y]{+/x*y}\:[xt;y]} /col-major

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


/statistics
avg:{(+/x)%0.+#x}
var:{(+/x*x:(x-avg x))%-1+#x}
std:{%var x}

rem:{x/x+x/y} /like apl for negative arguments
ej:{(y j),'x_z i j:&~0N=i:(z x)?y x} /sym t1 t2


/serialize
`"pack.":{w:{(`c@,#x),x};($t),$[`s~t:@x;`pack@$x;x~*x;w `c@,x;`L~@x;(`c@,#x),,/`pack@'x;(@x)':`D`T;(`pack@.x),`pack@!x;`S~t;,/`pack@$x;w `c x]}
`"unpack.":{s:x;g:{[n]r:n#s;s::n_s;r};n:{*`i@g 4};u:{x;$[(t:*g 1)':"bcifz";*(`$t)g n[];t~"s";`$u 0;t~"S";`$u 0;t~"L";u'!n[];t~"D";(u 0)!u 0;t~"T";+(u 0)!u 0;(`$_t+32)g n[]]};u 0}

csv:{c:{s:`$'x@i:&x':"ifzs";n:`i$" "\:-1_@[x;i;" "];y[a]:(y[a],''"a"),''y[1+a:&s=`z];s$'y n};s:$[" "~(*x);`split@;(*x)\:];x:1_x;y:+s'$[`L~@y;y;"\n"\:y];$[#x;c[x;y];y]} / csv[";1z3i2f"; "input.."]

/plot
`"plot.":{[d]l:$!d;v:.d; t:$[2~#d;`xy;`polar];
 y:$[t~`xy; $[`L~@y:v 1;y;,y];          $[`L~@y:_*v;y;,y]]
 x:$[t~`xy; $[`L~@x:v 0;x;(,x)@(#y)#0]; $[`L~@x:imag@*v;x;,x]]
 xt:`tics(&/&/x;|/|/x);yt:`tics(&/&/y;|/|/y)
 a:$[t~`xy;(xt 0;*-1#xt;yt 0;*-1#yt);(-a;a;-a;a:*|`tics@0.,|/|/abs@*v)]
 c:c@(#c:11826975 950271 2924588 2631638 12412820 4937356 12744675 8355711 2276796 13614615)/!#x
 style:$[t~`polar;"..";`i~@**y;"||";"--"] / -.| line points bar
 size: $[t~`polar;2;style~"||";(--/((**x),-1#*x))%-1+#*x ;2]
 lines:{`style`size`color`x`y!(style;size;z;x;0.+y)}'[x;y;c]
 `L`T`t`l`a!(lines;"";t;l;a)}
 
plotwh:{[x;fs;W;H]; ; a:x`a;T:x`T;grey:13882323
 C:(W%2;H%2);R:(W%2)&(H%2)-fs
 dst:$[`xy~x`t;(fs;W-fs;H-fs;fs);((C-R),C+R)0 2 3 1];rdst:(fs;fs;W-2*fs;H-2*fs)
 xs:(a 0 1)(dst 0 1)' /transform axis to canvas
 ys:(a 2 3)(dst 2 3)'
 bars:{[l]$["|"':l`style;(`color;l`color),,/{(`Rect;((-dx%2)+xs x;ys y;dx:-/xs(l`size;0.);(ys a 2)-ys y))}'[l`x;l`y];()]}
 line:{[l]$["-"':l`style;(`linewidth;l`size;`color;l`color;`poly;(xs l`x;ys l`y));()]}
 dots:{[l]$["."':l`style;(`color;l`color),,/{(`Circle;(xs x;ys y;1.5*l`size))}'[l`x;l`y];()]}
 c:(`clip;(0;0;W;H);`font;("monospace";fs);`color;0;`text;((W%2;fs);1;T))
 xy:{[]c,:(`text;((fs;H);0;$a 0);`text;((W%2;H);1;(x`l)0);`text;((W-fs;H);2;$a 1))
       c,:(`Text;((fs;H-fs);0;$a 2);`Text;((fs;H%2);2;(x`l)1);`Text;((fs;fs);2;$a 3))
       c,:(`color;0;`linewidth;2;`rect;rdst)      /todo: clip rdst
       c,:(`linewidth;1;`color;grey)
       c,:(`clip;rdst)
       c,:,/{(`line;0.+(x;dst 2;x;dst 3))}'xs`tics x[`a;0 1]
       c,:,/{(`line;0.+(dst 0;x;dst 1;x))}'ys`tics x[`a;2 3]}
 po:{[]c,:(`text;((C 0;H);1;(x`l)0);`text;(C+.75*R;6;$(x`a)3))
       c,:(`font;("monospace";_fs*.8)),,/{(`text;(C+R*(_;imag)@'x;y;z))}'[1@270.+a;0 0 6 6 4 4 2 2;$a:30 60 120 150 210 240 300 330]
       c,:(`color;0),/{(`line;,/+C+(R-fs%2;R)*/:(_;imag)@'x)}'1@30.*!12
       /c,:(`clip;C,R) /bug in cairo?
       c,:(`color;grey;`linewidth;1;`line;((-R)+*C;C 1;R+*C;C 1);`line;(*C;(-R)+C 1;*C;R+C 1))
       c,:,/{(`circle;0.+C,x)}'r:(xs@`tics 0.,x[`a;3])-*C
       c,:(`color;0;`linewidth;2;`circle;C,R)}
 $[`xy~x`t;xy[];po[]]
 c,:,/bars'x`L
 c,:,/line'x`L
 c,:,/dots'x`L}

`"tics.":{[minmax]nice:{[x;r]f:x%0.+10^ex:_log[10;x];(1 2 5 10.@*&(~f>1 2 5 0w;f<1.5 3 7 0w)[r])*10^ex};e:nice[-/|minmax;0];s:nice[e%4.;1];n:_1.5+e%s;$[~(minmax 1)>*-2#r:s*(_(*minmax)%s)+!n;-1_r;r]}

`"ceg.":{(x i)!0-':1+i:&(1_~~':x),1b} /#'= (^x)
hist:{$[`i~@x;hist[(x;&/y;|/y);y];(Y;(`38)[0;(`ceg@^1+((d%2.0)+-1_Y:(x 1)+(d:(--/1_x)%-1.+n)*!n)'y)@!n:_0.+*x])]} /[xn;y] [xn xmin xmax;y] e.g. plot `x`y!hist[10 -4 4.;?-1000]
