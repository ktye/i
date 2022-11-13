
fft: {[f;x]{[x;p;k;e]@[x;k,j;(x[k]+x[j]*e),x[k]-x[j:k+p]*e]}/[x f 0;f 1;f 2;f 3]}
fftn:{[n]l:*&|2\\n;e:angle[1;-(i:!n)*360.%n]
 (2//|2\\i;2^/:!l;|&'~!l#2;e@|(h/2*)\:!h:n%2)}
/ bitrev   powers index    exponentials


n:2^15
f:fftn n

i:1001;r:while[i-:1;fft[f;`z@?2*n]]
 \4#`f r

