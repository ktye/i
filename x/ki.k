/ki parses i source and writes the IR as a k table
/see ik.k for a description of i.

/program IR
C:0x    /constants
D:0x    /data
T:0#`   /type
P:!0    /parent
I:!0    /ivalue
S:0#`   /svalue

/k appends a list of subtrees to IR (and translate parents)
/e.g. k((`con`lit;0N 1;0N 1;`x`i);(..));
k:{x:@[,/'+x;1;+;(#T)+,/n#'0,-1_+\n:#'*'x];T,:*x;P,:0|x 1;I,:x 2;S,:x 3}


/prg mem con var tab
/kiprg:{(`prg;x;0;x)}
/kimem:{(`mem;x;1;`a)}

/kiprg:{k(`prg;0;0*x)}
/kimem:{x;k(`mem;0;1;`a)}
/kicon:{... (i:*&(.)~/:x)~*&`const~/:x }

kival:{x;333} /todo: int32 value or intern in C
kivar:{v:x -2+i:&(`64)~/:x;s:x -1+i;{(`con`lit;0N 1;0N,kival y;(`$-1_c),`$*|c:$x)}'[s;v]}
kicon:{v:x -3+i:&(`74)~/:x;s:x -2+i;{(`var`lit;0N 1;0N,kival y;(`$-1_c),`$*|c:$x)}'[s;v]}

l:`p@<`"x.i"
 \k kivar l
