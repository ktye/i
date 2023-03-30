/softfloat.k replaces all float-nodes with function calls in IR (modifies T P I S)

del:{T::T n:(!#P)_x;P::@[P n;&(P n)>x-1;-;1];I::I n;S::S n;}         /delete a node
ins:{[i;x;y](i#x),y,i_x}                                             /insert (todo vectorize)
INS:{[i;x]T::ins[i;T;*x];P::ins[i;@[P;&P>i-1;+;1];x 1];P[1+i]:i;I::ins[i;I;x 2];S::ins[i;S;x 3]}


/replace float-cast with calls
S[i]:`$($S i:i@&`f=S 1+i:&T=`cst),\:"fcst";T[i]:`cal;del'|1+i        /cast ji<-f [ij]fcst
S[i]:`$"f",/:($S 1+i:&(`f=S)&T=`cst),\:"cst";T[i]:`cal;del'|1+i      /cast f<-ji f[ij]cst


/insert casts when for float load/store
cast:{INS[z;(`cst;P z;0N;x)];INS[1+z;(`typ;z;0N;y)];P[2+z]-:1}       /insert cast node: uint64(F64(x))
cast[`k;`f]'|i:&(S=`f)&(T=`lod)                                      /load-f64: cast to k
cast[`j;`f]@'{(&P=x)1}'|i:&(S=`f)&(T=`sto)                           /store: cast 2nd child


/replace float types with u64 for many node types
S[&(S=`f)&T':`lit`arg`loc`res`ret`cli`typ`cnd`swc]:`k               /`f->`k
S[&(S=`f)&T':`lod`sto                            ]:`j


/replace ops with calls to softfloat implementations
f:f!`$"f",/:$f:`add`sub`mul`div`les`lte`eql`neq`mor`gte              /binary ops/cmp
S[i]:f T i:i@&`f=S i:&T':!f;T[i]:`cal                                /softfloat calls
S[i:&(T=`neg)&S=`f]:`fneg;  T[i]:`cal                                /negate


/replace float intrinsics (pseudo-functions) with softfloat implementations
{S[&S=y]:x}'[`fabs`fsqr`flor`fmin`fmax`fcps;`$"F64",/:$`abs`sqrt`floor`min`max`copysign]


/remove reinterpretations (now non-ops) and lift child nodes
j:1+i:&S':`"F64reinterpret_i64"`"I64reinterpret_f64"
T[i]:T j;I[i]:I j;S[i]:S j;del'|j
