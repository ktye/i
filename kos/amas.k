/armv6 assembler


/calculate all valid immediates (8bit value, 2*4bit rotate)
/encode number: imm[x] (low 12bits)  0N(invalid)
rot:{y@32/32-x-!32}
mmi:{2//rot[2*y;(32#2)\\x]}
imm:(,/mmi/:\:[!256;!16])!,/(!256)+/:\:256*!16
