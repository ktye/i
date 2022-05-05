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
/D~dp P

/non-recursive evaluate(convert to string), one function per node type
F:({y;         , V x            }   /leaf
   {,/"(",(y 0),(V x),(y 1),")" })  /dyadic: x(index into orig table) y(already evaluated child strings)
   
   
t:+`t`p`d`v!(T;P;D;V)


/non-compacting
flats1:{t:x`t;p:x`p            
 s:&#t                        /evaluated strings
 c:{x[y],:z}/[(#p)^!0;p;!#p]  /each node's child vector
 c[0]:1_c 0                   /remove root child
 l:|/d:x`d                    /depth level
 while[-2<l-:1;f:F t i:&d=1+l;*s[i]:f.'+(i;s c i)]}
 
 
flats:{t:x`t            
 s:&#t                  /evaluated strings
 c:(&=)/:[p;!#p:x`p]    /each node's child vector
 l:|/d:x`d              /depth level
 while[-2<l-:1;f:F t i:&d=1+l;*s[i]:f.'+(i;s c i)]}
 
 \flats1 t
