/softfloat.k replaces all float-nodes with function calls in IR (modifies T P I S)

del:{T::T n:(!#P)_x;P::@[P n;&(P n)>x-1;-;1];I::I n;S::S n;}         /delete a node


S[&(S=`f)&T':`lit`arg`loc`res`ret`cli`typ`cnd`swc]:`k                /`f->`k
S[&(S=`f)&T':`lod`sto                            ]:`j


S[i]:`$"f",/:($S i:i@&`f=S 1+i:&T=`cst),\:"cst";T[i]:`cal;del'|1+i   /cast ji<-f
S[i]:`$($S 1+i:&(`f=S)&T=`cst),\:"fcst";        T[i]:`cal;del'|1+i   /cast f->ji


f:f!`$"f",/:$f:`add`sub`mul`div`les`lte`eql`neq`mor`gte              /binary ops/cmp
S[i]:f T i:i@&`f=S i:&T':!f;T[i]=`cal                                / -> softfloat calls
S[i:&(T=`neg)&S=`f]:`fneg;  T[i]=`cal                                /negate

{S[&S=y]:x}'[`fabs`fsrt`flor`fmin`fmax`fcps;`$"F64",/:$`abs`sqrt`floor`min`max`copysign]

j:1+i:&S':`"F64reinterpret_i64"`"I64reinterpret_f64"                 /non-ops, lift child
T[i]:T j;I[i]:I j;S[i]:S j;del'|j

/todo: SetF64(x,f)   f now j  must be  SetI64(x, int64(uf))
/         F64(x)->f  f now j                     ^^^^^ insert cast node