/risc-v assembler
asm:{l:,/{*|x}'x@a:&,/5=#'x;as'x,'!#x}
as: {.[*x;(x 1 2),$[5~@u:x 3;4*(a l?u)-*|x;u]]}

R:{[op;f3;f7;rd;rs1;rs2]op+/128 4096 32768 1048576 33554432*(rd;f3;rs1;rs2;f7)}
I:{[op;f3;rd;rs1;imm]op+/128 4096 32768 1048576*(rd;f3;rs1;imm)}
J:{[op;rd;r1;imm]o:-2147483648*0>imm;imm:n21 imm;o+op+/128 4096 524288*(rd;256/4096\imm;4096/imm)}
S:{[op;f3;rs1;rs2;imm]imm:n12 imm;op+/128 4096 32768 1048576 33554432*(32/imm;f3;rs1;rs2;32\imm)}
B:{[op;f3;rs1;rs2;imm]o:-2147483648*0>imm;imm:n12 imm;o+op+/128 4096 32768 1048576 33554432*((32/imm)+2/2048\imm;f3;rs1;rs2;(4096\imm)+64/32\imm)}
n12:{$[x<0;4096+x;x]};n21:{$[x<0;2097152+x;x]}

 jal:J[111];    jalr:I[103;0]
 add:R[51;0;0];  sub:R[51;0;32]; sll:R[51;1;0]; slt:R[51;2;0]
sltu:R[51;3;0];  xor:R[51;4;0];  srl:R[51;5;0]; sra:R[51;5;32]
  or:R[51;6;0];  and:R[51;7;0];  mul:R[51;0;1]; div:R[51;4;1] 
divu:R[51;5;1];  rem:R[51;6;1]; remu:R[51;7;1]
 addi:I[19;0];  slli:R[19;1;0]; srli:R[19;5;0]; xori:I[19;4]
  ori:I[19;6];  andi:I[19;7];    clz:R[19;1;48]
   lb:I[3;0];    lbu:I[3;4];      lw:I[3;2];      sb:S[35;0] 
   sw:S[35;2];   fld:I[7;3];     fsd:S[39;3] 
  beq:B[99;0];   bne:B[99;1];    blt:B[99;4];    bge:B[99;5]
 bltu:B[99;6];  bgeu:B[99;7]
rm:7                                                 /rounding-mode
  fadd:R[83;rm;1]; fsub:R[83;rm;5]; fmul:R[83;rm;9]; fdiv:R[83;rm;13]
  fsqr:R[83;rm;45]; fle:R[83;0;81];  flt:R[83;1;81];  feq:R[83;2;81]
fcvtwd:R[83;rm;97];               fcvtdw:R[83;rm;105]


/example
list:{ \(_10)/:$x}
list asm(   sb,1,21,2047,`start
            sw,1,21,2044
	    lb,1,21,2047
	    lw,1,21,2044
	   lbu,2,13,1234
	  addi,1,21,2005
	  slli,2,3,31
	  srli,2,3,29
	  xori,4,5,6
	   ori,7,8,-3
	  andi,9,0,-5
	   clz,1,2,0
           add,1,21,3
           sub,1,21,4 
	   sll,5,6,7
	   slt,29,30,31
	  sltu,1,2,3
	   xor,4,5,6
	   srl,7,8,9
	   sra,10,11,12
	    or,13,14,15
	   and,16,17,18
	   mul,19,20,21
	   div,22,23,24
	  divu,25,26,27
	   rem,28,29,30,`one  /label
	  remu,31,0,1
	   beq,1,2,`one       /resolve as relative address (-8)
	   bne,3,4,1026
	   blt,5,6,-2
	   bge,7,8,-4096
	  bltu,7,8,4094
	   fld,1,2,-2034
	   fsd,1,2,-2034
	  fadd,1,2,3
	  fsub,4,5,6
 	  fmul,7,8,9
	  fdiv,10,11,12
	  fsqr,13,14,0
	   fle,15,16,17
	   flt,18,19,20
	   feq,21,22,23
	fcvtwd,23,25,0
	fcvtdw,27,28,0
	   jal,1,0,`start
	  jalr,1,2,567
	   add,0,0,0) 
     
 