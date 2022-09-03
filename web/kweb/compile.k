
np:{n:(&#x){x[y]+:1}/x;n[0]-:1}  /number of children
dp:{+/0<(x@)\:!#x}               /depth
lp:{(~~)#@[&#x;x i;i:!#x]}       /last child


/visiblility from checked vector
/root is always visible (1,1_)
vc:{[c;p]&1,1_&/1_(c@P@)\:!#p}

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
 d:dp P                                      /depth
 m:1+|/d                                     /1+max depth
 c:&~np P                                    /child nodes
 l:lp P                                      /last child nodes
 i:!#P
 r:(m*#P)#" "
 r:@[r;(-1+  d l)+m*{x+!y-x}'[1+P l;l];"|"]  /step(3)
 r:@[r;(-1+1_d i)+m*1_i;"+"]                 /step(2)
 r:@[r;(-1+  d l)+m*l;"L"]                   /└ for last child
 r:@[r;d+m*!#P;("T";"-")@~np P]              /step (1)
 (#P)^r}


dt:drawtree P  /ascii tree
c:&#P          /checked state
v:,0
clicktree:{v::vc[c[v x]:~c[v x];P];}
