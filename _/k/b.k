/listbox.k
path:()
exec:{$[x~`edit;edit[];0~#x;1_,tags::Tags;0~@f:.*x;disp f walk path;disp walk $[5~@x;path::,x;x]]}
walk:{$[0~#x;();atx[v;1_x;ist v:.*x]]}
push:{exec path::path,,$[(1~#path)&ist u:walk path;-1+x^0;x]}
back:{disp walk path::-1_path}
edit:{tako 1;$[(1<#*|path)&(2~#path)&ist t:walk path;dis7 t; 'k t]}
teds:{(!x)!*'.x}
ok:  {s:*path;asn[s;$[ist n:.s;tas[n;1_path;x];$[#p:1_path;@[n;$[1~#p;*p;p];x];x]]];}

tako:{tags::(tags^`ok),$[x;`ok;0#`]}
disp:{t:@x;tako t<6;$[ist x;dis8 x;6~t; 'k'x;7~t;dis7 x; 'kx]}
dis7:{k:$!x;pad:rpad[;|/#'k]; ((pad'k),'"|"),' 'k'.x} /display dict
dis8:{" "/:'&rpad''[l;|/'#''l:tst''(!x),'.x]}         /display table
rpad:{y#x,(y-#x)#" "}
kst: {$[1~@x; x; 'kx]}
tst: {$[5~@x;$x; 'kx]}
ist: {$[~7~@x;0;~#x;0;&/(n:#*.x)~'#'.x;n>1;0]}        /dict is table
atx: {$[#y;atx[$[z;tix[x;*y];x@*y];1_y;0];x]}
tix: {@[;y]'x}                                        /index table x by row y

data:`!()
drop:{data,:x!,y;,/($x;": ";$#y;" bytes")}
asn: {.(:;x;$[6~@y;,y;y])}
tas: {$[0~#y;z;1~#y;@[x;(,!x),y;$[7~@z;.z;z]];@[x;@[y;0 1;y 1 0];z]]}
tag: {" "/:$?tags}
/m.k
sin:129 'F;cos:130 'F;exp:131 'F;log:132 'F;log10:{0.4342944819032518* 'Fx};pow10:{131 'F2.302585092994046*x}

/linear algebra

/solve[A;b]          (solve A*x=b overdetermined, real or complex)
/qrsolve[q:qr A;b]   (reuse decomposition for multiple RHS)

solve:  {[A;b];qrsolve[qr A;b]}
qrsolve:{[q;b];qrslv[q;qrmul[q;b]]}

qr: {m:#x;x:&x;n:#x;x:,/x;f:$[3~t:@x;qrh;qzh]                                              /qr decomposition
 r:n f/(x;t$!0;m;n;!m;(,/(n;m))#!n*m);`Q`D`m`n!(r 0;|r 1;r 2;r 3)}
qrh:{Q:x 0;i:#D:x 1;m:x 2;n:x 3;row:x 4;sub:1_x 5;ii:i+i*m                                 /householder step(real)
 s:norm Q row;D,:d:s*1 -1@Q[ii]>0;s:1%%s*s++Q ii;Q[ii]-:d;Q[row]*:s
 $[#sub;Q[,/sub]-:,/Q[row]*/+/'(Q[row]*/Q[sub]);]
 (Q;D;m;n;m+1_row;1_'sub)}
qzh:{Q:x 0;i:#D:x 1;m:x 2;n:x 3;row:x 4;sub:1_x 5;ii:i+i*m                                 /householder step(complex)
 s:norm Q row;d:-s@&Q[ii];D,:d:-s@&Q[ii];s:1%%s*s++Q ii;Q[ii]-:d;Q[row]*:s
 $[#sub;Q[,/sub]-:,/Q[row]*/+/'((%Q[row])*/Q[sub]);]
 (Q;D;m;n;m+1_row;1_'sub)}

qrmul:{m:x`m;n:x`n;f:$[3~t:@x`Q;qrml;qzml];r:n f/(x`Q;t$y;!x`m;!m;x`m);(x`n)#r 1}          /calculate QT*b
qrml: {Q:x 0;b:x 1;row:x 2;yi:x 3;m:x 4;b[yi]-:Q[row]*+/Q[row]*b[yi];(Q;b;m+1_row;1_yi;m)}
qzml: {Q:x 0;b:x 1;row:x 2;yi:x 3;m:x 4;b[yi]-:Q[row]*+/(%Q[row])*b[yi];(Q;b;m+1_row;1_yi;m)}

qrslv:{n:x`n;r:n qrsl/(x`Q;y;x`D;n-1;x`m;x`n;!0);r 1}                                      /solve R*x=QT*b with back substitution
qrsl: {Q:x 0;b:x 1;D:x 2;i:x 3;m:x 4;n:x 5;bi:x 6;col:i+m*bi
 $[#bi;b[i]-:+/Q[col]*b[bi];];b[i]%:*D
 (Q;b;1_D;i-1;m;n;i,bi)}
 

norm:{s*%+/x*x%:s:|/x:+0.+x}                                                               /vector norm(l2)

matvec:{[A;b],/+/'A*\b}



test:{[]
 small:1.0e-14
 A: 0.+(1 -2 3;5 3 2;2 3 1;4 -1 1)
 r: 0.+1 2 3
 b: matvec[A;r]
 x: solve[A;b]
 $[small>e:|/+x-r; \(`ok;e); \(`fail;e)]
 A: 0a0+(1 -2a90 3;5a90 3 2;2 3 1;4 -1 1)
 r: 1 2 3a30
 b: matvec[A;r]
 x: solve[A;b]
 $[small>e:|/+x-r; \(`ok;e); \(`fail;e)]
}
 



ucal:{[s;t;r]n:#r;i:(u:?r)?r;q:qr[t,'0a+i=(n#,!#u)];X:q qrsolve/&s;K:(#t 0)#'X}
/ s0:1a0 2a30
/ s:(s0;s0+2a0 0a0;s0+0a0 2a0)
/ t:(0a0 0a0;1a0 0a0;0a0 1a0)
/ r:`a`a`a

uidx:{[u;a]solve[1@0.+a;u]}
/ u:4.8a48 4.9a47 1.6a277 1.7a279
/ a:(0 0;0 0;0 180;0 180)



Tags:`List`idx
List:`u`a!(4.8a48 4.9a47 1.6a277 1.7a279;(0 0;0 0;0 180;0 180))
idx:{l:walk path;uidx[l`u;l`a]}


Tags:`List`unb
List:`s`t`r`c!((3a20 4a50;12.8a5 4.1a51;3.1a19 12.9a14;3a300 2a10);(0a0 0a0;10a0 0a0;0a0 10a0;0a0 0a0);`a`a`a`a;1 1 1 0)
unb:{w:&List`c;s:List[`s;w];t:List[`t;w];r:List[`r;w];K:qr[ucal[s;t;r]];l:walk path;K qrsolve/l`s}

tags: Tags
