f:{[nrhs]
 rows:10000
 cols:300

 A:cols^?rows*cols
 B:nrhs^?rows*nrhs
 |/|/solve[A;B]}
