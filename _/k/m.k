sin:129 'F;cos:130 'F;exp:131 'F;log:132 'F;log10:{0.4342944819032518* 'Fx};pow10:{131 'F2.302585092994046*x}

/linear algebra

/solve[A;b]          (solve A*x=b overdetermined, real or complex)
/qrsolve[q:qr A;b]   (reuse decomposition for multiple RHS)

solve:  {[A;b];qrsolve[qr A;b]}
qrsolve:{[q;b];qrslv[q;qrmul[q;b]]}
 
qr: {m:#x;x:&x;n:#x;x:,/x;f:$[3~t:@x;qrh;qzh]                                              /qr decomposition
 r:(n;f)/:(x;t$!0;m;n;!m;(,/(n;m))#!n*m);`Q`D`m`n!(r 0;|r 1;r 2;r 3)}
qrh:{Q:x 0;i:#D:x 1;m:x 2;n:x 3;row:x 4;sub:1_x 5;ii:i+i*m                                 /householder step(real)
 s:norm Q row;D,:d:s*1 -1@Q[ii]>0;s:1%%s*s++Q ii;Q[ii]-:d;Q[row]*:s
 $[#sub;Q[,/sub]-:,/Q[row]*/:+/'(Q[row]*/:Q[sub]);]
 (Q;D;m;n;m+1_row;1_'sub)}
qzh:{Q:x 0;i:#D:x 1;m:x 2;n:x 3;row:x 4;sub:1_x 5;ii:i+i*m                                 /householder step(complex)
 s:norm Q row;d:-s@&Q[ii];D,:d:-s@&Q[ii];s:1%%s*s++Q ii;Q[ii]-:d;Q[row]*:s
 $[#sub;Q[,/sub]-:,/Q[row]*/:+/'((%Q[row])*/:Q[sub]);]
 (Q;D;m;n;m+1_row;1_'sub)}

qrmul:{m:x`m;n:x`n;f:$[3~t:@x`Q;qrml;qzml];r:(n;f)/:(x`Q;t$y;!x`m;!m;x`m);(x`n)#r 1}       /calculate QT*b
qrml: {Q:x 0;b:x 1;row:x 2;yi:x 3;m:x 4;b[yi]-:Q[row]*+/Q[row]*b[yi];(Q;b;m+1_row;1_yi;m)}
qzml: {Q:x 0;b:x 1;row:x 2;yi:x 3;m:x 4;b[yi]-:Q[row]*+/(%Q[row])*b[yi];(Q;b;m+1_row;1_yi;m)}

qrslv:{n:x`n;r:(n;qrsl)/:(x`Q;y;x`D;n-1;x`m;x`n;!0);r 1}                                   /solve R*x=QT*b with back substitution
qrsl: {Q:x 0;b:x 1;D:x 2;i:x 3;m:x 4;n:x 5;bi:x 6;col:i+m*bi
 $[#bi;b[i]-:+/Q[col]*b[bi];];b[i]%:*D
 (Q;b;1_D;i-1;m;n;i,bi)}
 

norm:{s*%+/x*x%:s:|/x:+0.+x}                                                               /vector norm(l2)

matvec:{[A;b],/+/'A*\:b}

test:{[]small:1.0e-14
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

\
/matlab version(does not transpose A)
function [A,d]=qrh(A);
 [m,n]=size(A);
 for i=1:n,
 s=norm(A(i:m,i));
 if s==0, error(’rank(A)<n’), end
 if A(i,i)>=0, d(i)=-s; else d(i)=s; end
 f=sqrt(s*(s+abs(A(i,i))));
 A(i,i)=A(i,i)-d(i);
 A(i:m,i)=A(i:m,i)/f;
 if i<n
  A(i:m,i+1:n)=A(i:m,i+1:n)-A(i:m,i)*(A(i:m,i)’*A(i:m,i+1:n));
 end
end
