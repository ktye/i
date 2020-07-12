norm:{s*%+/x*x%:s:|/x:+0.+x}

qr:{m:#x;x:&x;n:#x;n qrh/(,/x;!0;m;n)}
qrh:{Q:x 0;i:#R:x 1;m:x 2;n:x 3
(Q;R;m;n)}

\

/k9
/8-local limit
/matrix assign

qr: {[A]n:#*A;m:#A;A:,/+A;(n;qrh)/:(A;!0;m;n)}       /qr decomposition A:mxn m>n
qrh:{[x]Q:x 0;i:#R:x 1;m:x 2;n:x 3;ii:i+i*m          /householder step Q:nxm
  s:nrm@Q(row:ii+!m-i);col:m*(i+1)_!n
  R,:d:s*1 -1@Q[ii]<0
  s:1%sqrt s*s+abs Q ii;Q[ii]-:d;Q[row]*:s
  Q[,/rc]-:,/Q[row]*/+/Q[row]*Q[rc:row+\col]
  (A;R;m;n)}
  
abs:{[x]x[i]:-x i:&x<0;x}
nrm:{[x]s*sqrt+/x*x%:s:|/x:abs@x}                   /l2 vector norm



idx:{[i;k;n],/k+/i*n}

/matlab version(does not transpose A)
function [A,d]=HouseholderQR(A);
 [m,n]=size(A);
 for i=1:n,
 s=norm(A(i:m,i));
 if s==0, error(’rank(A)<n’), end
 if A(i,i)>=0, d(i)=-s; else d(i)=s; end
 fak=sqrt(s*(s+abs(A(i,i))));
 A(i,i)=A(i,i)-d(i);
 A(i:m,i)=A(i:m,i)/fak;
 if i<n, % transformation of the rest of the matrix G:=G-u*(u’*G)
 A(i:m,i+1:n)=A(i:m,i+1:n)-A(i:m,i)*(A(i:m,i)’*A(i:m,i+1:n));
 end
end