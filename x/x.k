/ (1+2)-3
/               -
/            +     3
/          1   2


/ 0 1 2 3 4  /node id
T:1 1 0 0 0  /type 1(dyadic) 0(leaf)
P:0 0 1 1 0  /parent
D:0 1 2 2 1  /depth
V: "-+123"   /node value

dp:{+/0<(x@)\:!#x} /depth from parent
 \D~dp P

/evaluate(convert to string), one function per node type
F:({y;        , V x          }   /leaf
   { "(",(y 0),(V x),y 1,")" })  /dyadic: x(index into orig table) y(already evaluated child strings)
   
   
t:+`t`p`d`v!(T;P;D;V)



i:&(t`d)=|/t`d     /vector of indexes where depth is max



f:{[t]t:t@&(t`d)=|/t`d
 (F t`t).'+(t`t;t`v) }
 
 /...fixed-point-scan

