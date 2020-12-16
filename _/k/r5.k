
/ ( (add;0;1;2;`label); ..


asm:{l:,/{*|x}'x@a:&,/4=#'x;.[*x;x 1 2 3]}

R:{[op;f3;f7;rd;rs1;rs2]op+/128 4096 32768 1048576 33554432*(rd;f3;rs1;rs2;f7)}
I:{[op;f3;rd;rs1;imm]op+/128 4096 32768 1048576*(rd;f3;rs1;imm)}
S:{[op;f3;rs1;rs2;imm]op+/128 4096 32768 1048576 33554432*(32/imm;f3;rs1;rs2;32\imm)}
B:{[op;f3;rs1;rs2;imm]op+/128 128 4096 32768 1048576 33554432 33554432*(32/imm;2/2048\imm;f3;rs1;rs2;64/32\imm;4096/imm)}

add:  R[51;0;0]   /add  rd,rs1,rs2
sub:  R[51;0;32]  /sub  rd,rs1,rs2
addi: I[19;0]     /addi rd,rs1,imm
mul:  R[51;0;1]   /mul  rd,rs1,rs2
div:  R[51;4;1]   /div  rd,rs1,rs2
divu: R[51;5;1]   /divu rd,rs1,rs2

xor:  R[51;4;0]   /xor  rd,rs1,rs2
or:   R[51;6;0]   /or   rd,rs1,rs2
and:  R[51;7;0]   /and  rd,rs1,rs2

lb:   I[3;0]      /lb rd rs1 imm
lw:   I[3;2]      /lw rd rs1 imm
sb:   S[35;0]     /sb rs1 rs2 imm
sw:   S[35;2]     /sw rs1 rs2 imm
fld:  I[7;3]      /fld rd rs1 imm
fsd:  S[39;3]     /fsd rs1 rs2 imm

beq:  B[99;0]     /beq  rs1 rs2 imm
bne:  B[99;1]     /bne  rs1 rs2 imm
blt:  B[99;4]     /blt  rs1 rs2 imm
bge:  B[99;5]     /bge  rs1 rs2 imm
bltu: B[99;6]     /bltu rs1 rs2 imm
bgeu: B[99;7]     /bgeu rs1 rs2 imm
