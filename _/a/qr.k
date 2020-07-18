/k9(2020.07.14)
solve:  {[A;b]qrsolve[qr A;b]}       /solve A*x=b  (#A)>#*A  (overdetermined)
qrsolve:{[q;b]qrslv[q;qrmul[q;b]]}   /reuse qr

qr:{[A]m:#A;A:+A;n:#A;r:(n;qrh)/:(,/A;!0;m;n;!m;!(n;m))            /qr decomposition
    `Q`D`m`n!(r 0;|r 1;r 2;r 3)}

qrh:{[x]Q:x 0;ii:(1+x 2)*#D:x 1;row:x 4;sub:1_x 5                  /householder step
 s:norm Q row;D,:d:s*1 -1@Q[ii]>0;s:1%sqrt s*s+absa Q ii;Q[ii]-:d;Q[row]*:s
 $[#sub;Q[,/sub]-:,/Q[row]*/+/'(Q[row]*/Q[sub]);]
 (Q;D;x 2;x 3;(x 2)+1_row;1_'sub)}
 
qrmul:{[q;b]n:q`n;m:q`m;r:(n;qrml)/:(q`Q;b;!m;!m;m);n#r 1}         /calculate QT*b
qrml: {[x]Q:x 0;b:x 1;row:x 2;yi:x 3;m:x 4;b[yi]-:Q[row]*+/Q[row]*b[yi];(Q;b;m+1_row;1_yi;m)}

qrslv:{[q;b]n:q`n;r:(n;qrsl)/:(q`Q;b;q`D;n-1;q`m;n;!0);r 1}        /solve R*x=QT*b (back substitution)
qrsl: {[x]Q:x 0;b:x 1;D:x 2;i:x 3;m:x 4;n:x 5;bi:x 6;col:i+m*bi
 $[#bi;b[i]-:+/Q[col]*b[bi];];b[i]%:*D
 (Q;b;1_D;i-1;m;n;i,bi)}
 
norm::s*sqrt+/x*x%:s:|/x:abs x
abs:{[x]x;x[i]:-x[i:&x<0];x}
absa::$[x<0;-x;x]


/e.g. A*x=b (6x4)
matvec::+/'x*\y
A:4^?24
x:?4
b:matvec[A;x]
(`err;|/abs x-solve[A;b])
