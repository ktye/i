trim:{(+/(_10 13 32)=*x)_x}
mf:{r:"{"\:x;`hdr`body!(trim/r 0;r 1)}
compile:{mf'-1_"}"\:*(_10 92)\:x} 
/wasm:{[x]0x0061736d01000000}
wasm:{x}
