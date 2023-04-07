
np:{n:(&#x){x[y]+:1}/x;n[0]-:1}  /number of children
dp:{+/0<(x@)\:!#x}               /depth
lp:{(~~)#@[&#x;x i;i:!#x]}       /last child


/visiblility from checked vector
/root is always visible (1,1_)
vc:{[c;p]&1,1_&/1_c@(p@)\:!#p}

/nodes that are last children of their parents


/ (+ 1 (- (* 2 3)) (% 4))
/          i  d  P        (1)              (2)              (3)
/ ┬     +  0  0  0         ┬     +          ┬     +          ┬     +
/ └┬    1  1  1  0          ┬    1          └┬    1          └┬    1
/  ├┬   -  2  2  1           ┬   -           ├┬   -           ├┬   -
/  │└┬  *  3  3  2            ┬  *            └┬  *           │└┬  *
/  │ ├─ 2  4  4  3             ─ 2             ├─ 2           │ ├─ 2
/  │ └─ 3  5  4  3             ─ 3             └─ 3           │ └─ 3
/  └┬   %  6  1  1           ┬   %           └┬   %           └┬   %
/   └─  4  7  2  6            ─  4            └─  4            └─  4
/  
/ (1) draw at depth d:                  ┬ for nodes ─ for leaves
/ (2) draw at depth d-1:                └ if last child, ├ otherwise
/ (3) connect with parent at depth d-1: │ for rows 1+P i..i-1

drawtree:{[P]
 d:dp P                         /depth
 r:m:1+|/d                      /1+max depth
 a:~np P                        /child nodes
 w: lp P                        /last child nodes
 t:!#P
 r:((d w)+r*{x+!y-x}'[1+P w;w]  /step(3)    "|"
    (1_d)+r*1_t                 /step(2)    "+"
    (d w)+r*w)-1                /last child "L"
 e:@[(m*#P)#" ";r;"|+L"]
 e:@[e;d+m*!#P;"T-"a];(#P)^e}   /step (1)

nodes:($#P)," nodes"
dt:drawtree P    /ascii tree
c:(1+&#P)-T=`fun /checked state (collapse function nodes)
v:vc[c;P]        /visibility
clicktree:{v::vc[c[v x]:~c[v x];P];}
