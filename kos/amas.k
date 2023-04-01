/armv6 assembler

/O R X Y M L A  table: opcode,dst,src,arg,imm,label,addr
/s s s s i s s  label/addr are symbolic and resolved

/conds, shifts, ...?

op:2097152*`and`eor`sub`rsb`add`adc`sbc`src`tst`teq`cmp`cmn`orr`mov`bic`mvn?
rm:`r0`r1`r2`r3`r4`r5`r6`r7`r8`r9`r10`fp`ip`sp`lr`pc`cpsc?
rn:65536*rm@
rd: 4096*rm@

rot:{y@32/32-x-!32}
mmi:{2//rot[2*y;(32#2)\\x]}
imm:(,/mmi/:[!256;!16])?   /encode possible immediates

al:3758096384              /always                 0b111...
im:33554432+imm@
orrr:al+/(op;rd;rn;rm)@'   /orrr`add`r3`r2`r3  -> -528338941 e0823003
orri:al+/(op;rd;rn;im)@'   /orri`sub`r2`r3,256 -> -498914303 e2432c01


/encode: go.dev/play/p/41_1hQtVWVW
