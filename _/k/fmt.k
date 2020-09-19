fmt:{kst@&x}
kst:{$[6~@x;lst x;(1~#x)&5~@x;$x; 'kx]}
lst:{f:*x;$[0~#x;"()";1~#x; 'k*x;2~#x;mon.x;(128~_f)&~@f;seq 1_x;3~#x;dya.x;brk[*x;1_x]]}
mon:{$[(1~#x)&0~@x;$[iad@*x;(kst y),$*x;,/kst'(x;y)];spc[kst x;kst y]]}
dya:{$[(1~#x)&0~@x;(emb y),($x),kst z;brk[x;(y;z)]]}
brk:{,/(kst x;"[";";"/:kst'y;"]")}
emb:{$[6>@x;kst x;"(",(kst x),")"]}
seq:{";"/:kst'x}
iad:{6>adv?`$$x}
adv:,/`$'"_"\:"'_/_\_':_/:_\:"
chr:(_"0"+!10),(_"A"+!26),_"a"+!26
spc:{x,$[(~iad@*|x)&62>chr?*y;" ",y;y]}

test:{x~fmt x}
t:test'("*x"
"1 2 3"
"{x+y}"
"*x"
"+/1 2 3"
"f x"
"f`a`b`c"
"3+4"
"{x+y}[1;2]"
"1;2;3")
