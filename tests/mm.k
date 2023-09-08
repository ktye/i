/halfkey example

A:`p"MM:{[a;w]Yk;yF;xF;f;rF;Rk  /matmul
 Y:+w                      /transpose w
 l:an;m:wn;n:Yn            /dimensions x:l*m y:m*n r:l*n
 R:l^0.+&Ki l*n            /allocate result
 N[l;x:a i;r:R i;Nj[n;y:Y j;f:0.;Nk[m;f+:(x k)*y k];r[j]:f]];
 x:rx x;y:rx y;r:rx r;R}   /fix refcounts"

codes:A 1+i:4*!(#(A:0,A))%4
names:A 2+i

 \simple[]

