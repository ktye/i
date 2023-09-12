/ ../intro.md  /: \: ': // \\ / \

qr:{K:!m:#*x;I:!n:#x;j:0;r:n#0a;turn:$[`Z~@*x;{(-x)@angle y};{x*1. -1@y>0}]
 while[j<n;I:1_I
  r[j]:turn[s:abs@abs/j_x j;xx:x[j;j]]
  x[j;j]-:r[j]
  x[j;K]%:%s*(s+abs xx)
  x[I;K]-:{+/x*y}/[(conj x[j;K]);x[I;K]]*\x[j;K]
  K:1_K;j+:1];(x;r;n;m)}

solve:{qslv:{H:x 0;r:x 1;n:x 2;m:x 3;j:0;K:!m
 while[j<n;y[K]-:(+/(conj H[j;K])*y K)*H[j;K];K:1_K;j+:1]
 i:n-1;J:!n;y[i]%:r@i
 while[i;j:i_J;i-:1;y[i]:(y[i]-+/H[j;i]*y@j)%r@i]
 n#y}
 q:$[`i~@*|x;x;qr x];$[`L~@y;qslv/[q;y];qslv[q;y]]}

dot:(+/*)\

A: \(3 5 -8 12.;-2 3 3 0.;7 -8 2 1.)
b: \dot[+A;x:1 2 3];0.00001>|/abs x-solve[A;b]

A: \+0a0+(1 -2a90 3;5a90 3 2;2 3 1;4 -1 1);
 \0.0001>|/abs r-solve[A;dot[+A;r:1a30 2a30 3a30]]

lu:{[A]i:0;k:!#A;P:!#A
 while[1<#k
  j:i+*&a=m:|/a:abs A[k;i]
  P[(i;j)]:P[(j;i)]
  A[(i;j)]:A[(j;i)]
  A[k:1_k;i]%:A[i;i]
  A[k;k]-:A[k;i]*\A[i;k]
  i+:1]
 (A;P)}

lusolve:{[LUP;b];A:*LUP;P:*|LUP
 r:{[x;i;a]x[i]-:+/(a k)*x k:!i}/[b P;!n:#A;A]
   {[x;i;a]x[i]:(x[i]-+/(a k)*x k:(1+i)_!#x)%a[i]}/[r;|!n;|A]}

A: \5^?25
x: \?5
b: \(+/*)\[A;x]
 \x-lusolve[lu A;b]

fft: {[f;x]{[x;p;k;e]@[x;k,j;(x[k]+x[j]*e),x[k]-x[j:k+p]*e]}/[x f 0;f 1;f 2;f 3]}
fftn:{[n]l:*&|2\n;e:angle[1;-(i:!n)*360.%n]
 (2/|2\i;2^/!l;|&'~!l#2;e@|{h!2*x}\!h:n%2)}

n:2^15
f:fftn n
r:fft[f;`z@?2*n]
 \`fft
 \#r
 \+/r

avg:{(+/x)%0.+#x}
var:{(+/x*x:(x-avg x))%-1+#x}
std:{%var x}
med: **|2^^


