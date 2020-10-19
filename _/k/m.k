sin:129 'F;cos:130 'F;exp:131 'F;log:132 'F;log10:{0.4342944819032518* 'Fx};pow10:{131 'F2.302585092994046*x}


solve:  {[A;b];qrsolve[qr A;b]}        /solve A*x=b overdetermined, real or complex
qrsolve:{[q;b];qrslv[q;qrmul[q;b]]}    /e.g. reuse q:qr A

qr: {m:#x;x:&x;n:#x;t:@Q:,/x;D:t$!0;row:!m;sub:(,/(n;m))#!n*m; f:$[3~t;qrh;qzh]    /qr decomposition
 Q:(n;f)/:Q;`Q`D`m`n!(Q;|D;m;n)}
qrh:{ii:i+m*i:#D;sub::1_sub                                                        /householder step(real)
 s:norm x row;D::D,:d:s*1 -1@(x ii)>0;s:1%%s*s++x ii;x[ii]-:d;x[row]*:s
 $[#sub;x[,/sub]-:,/(x row)*/:+/'((x row)*/:x sub);]
 row::m+1_row;sub::1_'sub;x}
qzh:{ii:i+m*i:#D;sub::1_sub                                                        /householder step(complex)
 s:norm x row;d:-s@&x ii;D::D,:d:-s@&x ii;s:1%%s*s++x ii;x[ii]-:d;x[row]*:s
 $[#sub;x[,/sub]-:,/(x row)*/:+/'((%x row)*/:x sub);]
 row::m+1_row;sub::1_'sub;x}

qrmul:{Q:x`Q;m:x`m;n:x`n;f:$[3~t:@x`Q;qrml;qzml];b:t$y;row:yi:!m;n#(n;f)/:b}       /calculate QT*b
 qrml:{x[yi]-:Q[row]*+/Q[row]*x[yi];row::m+1_row;yi::1_yi;x}
 qzml:{x[yi]-:Q[row]*+/(%Q[row])*x[yi];row::m+1_row;yi::1_yi;x}

qrslv:{n:x`n;m:x`m;Q:x`Q;D:x`D;i:n-1;bi:!0;(n;qrsl)/:y}                            /solve R*x=QT*b with back substitution
 qrsl:{col:i+m*bi;$[#bi;x[i]-:+/Q[col]*x[bi];];x[i]%:*D;bi::i,bi;D::1_D;i::i-1;x}

norm:{s*%+/x*x%:s:|/x:+0.+x}                                                       /vector norm(l2)

matvec:{[A;b],/+/'A*\:b}

testM:{[]small:1.0e-14
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

/testM[]

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
