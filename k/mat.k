

dotv:{{+/x*y}\:[x;y]}
dot:{$[~`L~@y;dotv[x;y];dotv/:[x;y]]}

/see /_/i5/mat/mat.go
qr:{K:!m:#*x;I:!n:#x;j:0;r:n#0a
 while[j<n
  s:abs/j_x j
  xx:x[j;j]
  r[j]:(-s)angle angle xx
  f:%s*(s+abs xx)
  x[j;j]-:r[j]
  x[j;K]%:f 
  sum:{+/x*y}/:[(conj x[j;K]);x[I;K]]
`sum \sum
`sxkj \sum*x[j;K]
  x[I;K]:x[I;K]-\:sum*x[j;K]
  K:1_K;j+:1]
 (x;r;n;m)}


 
A:+0a0+(1 -2a90 3;5a90 3 2;2 3 1;4 -1 1)
qr A