

dotv:{{+/x*y}\:[x;y]}
dot:{$[~`L~@y;dotv[x;y];dotv/:[x;y]]}

/see /_/i5/mat/mat.go
qr:{K:!m:#*x;I:!n:#x;j:0;r:n#0a
 while[j<n
  I: \1_I
  s:abs/j_x j
`s \s
  xx:x[j;j]
  r[j]:(-s)angle angle xx
`r \r
  f:%s*(s+abs xx)
`f \f
  x[j;j]-:r[j]
  x[j;K]%:f 
`x \x
`j \j
`K \K
`I \I
`xik \x[I;K]
  sum:{+/x*y}/:[(conj x[j;K]);x[I;K]]
`sum \sum
  x[I;K]:x[I;K]-sum*\:x[j;K]
  K:1_K;j+:1]
  \(x;r;n;m)}


 
A:+0a0+(1 -2a90 3;5a90 3 2;2 3 1;4 -1 1)
q:qr A
