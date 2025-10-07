/draw tree from T P I S (unicode)
tree:{[]
 np:{n:{x[y]+:1}/(,&#x),x;n[0]-:1} /number of children
 dp:{+/0<(x@)\!#x}                 /depth
 lp:{(~~)#@[&#x;x i;i:!#x]}        /last child
 dt:{[P]
  d:dp P                           /depth
  r:m:1+|/d                        /1+max depth
  a:~np P                          /child nodes
  w: lp P                          /last child nodes
  t:!#P
  r:((d w)+r*{x+!y-x}'[1+P w;w]    /step(3)    "|"
     (1_d)+r*1_t                   /step(2)    "+"
     (d w)+r*w)-1                  /last child "L"
  e:@[(m*#P)#" ";r;"|+L"]
  e:@[e;d+m*!#P;"R-"a];(#P)^e}     /step (1)
 r:{y/x\z};
 `<r["+";"├"]r["R";"┬"]r["L";"└"]r["|";"│"]r["-";"─"]("\n"/(.`"lxy.")[70 300;+``i`P`T`I`S!(`$(dt P);!#P;P;T;I;S)]),"\n"
 `0}
 