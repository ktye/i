
dotv:{{+/x*y}\:[x;y]}
dot:{$[~`L~@y;dotv[x;y];dotv/:[x;y]]}

solve:{H:x 0;r:x 1;n:x 2;m:x 3;j:0;K:!m
 while[j<n;y[K]-:(+/(conj H[j;K])*y K)*H[j;K];K:1_K;j+:1]
 i:n-1;J:!n;y[i]%:r@i
 while[i;j:i_J;i-:1;y[i]:(y[i]-+/H[j;i]*y@j)%r@i]
 n#y}
 
qz:{K:!m:#*x;I:!n:#x;j:0;r:n#0a
 while[j<n;I:1_I
  r[j]:(-s:abs/j_x j)angle angle xx:x[j;j]
  x[j;j]-:r[j]
  x[j;K]%:%s*(s+abs xx)
  x[I;K]-:{+/x*y}/:[(conj x[j;K]);x[I;K]]*\:x[j;K]
  K:1_K;j+:1];(x;r;n;m)}



testmat:{A:+0a0+(1 -2a90 3;5a90 3 2;2 3 1;4 -1 1)
 Y:10.770329614269006a8.19859 13a52.61986 11a30 5a30
 `fail`ok@0.00001>|/abs solve[qz A;Y]-1a30 2a30 3a30}
