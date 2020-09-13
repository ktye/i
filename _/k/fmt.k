fmt:{kst@&x}
kst:{$[6~@x;lst x;5~@x;$x; 'kx]}
lst:{$[2~#x;,/kst'x;k x;3~#x;dya.x;brk[*x;1_x]]}
dya:{$[0~@x;(emb y),($x),kst z;brk[x;(y;z)]]}
brk:{,/(kst x;"[";";"/kst'y;"]")}
emb:{"(",(kst x),")"}

/fmt"3+4"
fmt"*x"
