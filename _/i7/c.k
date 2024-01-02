/c.k the compiler.

print:{`<(`l@x),\"\n";x}

/8(args/scratch/rax:return) 7(saved by callee)
Rx:`rax`rdi`rsi`rdx`rcx`r8`r9`r11`rbp`r10`r11`r12`r13`r14`r15
all:,`main
compile:{(c1'x),\"\n"}

c1:{t: print lift print tf(.x)
 ($x),":\n ;",($.x)," todo.."}

tr:{y[`p]:0|x+y`p} /translate parent vector
pt:+`p`t`i`s!(!0;0#`;!0;0#`)
t1:+`p`t`i`s!(,0N),,'
tf:{               /table from func
 n:`n \#x          /arity
 v:`v \(.x)2       /vars args,locals
 b:split(*.x)      /body (instruction list)
 {x,tr[#x]y}/(,pt),te'b} /catenate subtables
te:{               /table from expression
 mo:{(-1_x),,(t1`op1,1,(`$$y)),tr[   #x;   *|x]}
 dy:{(-2_x),,(t1`op2,2,(`$$y)),tr[-1+#x;,/-2#x]}
 no:{$[`i~t:@x;t1(`lit;x;`);`s~t;t1(`var;0;x);0+`noun]}
 f:{$[~`v~@y;x,,no y;64>c:0+y;mo[x;y];dy[x;y]]}
*f/(,pt),x}

lift:{x[`p]-:+/p>/(p:x`p)i:&(`var=1_(x`t),`)&`"."=x`s;x i_!#x} / . var -> var

split:{1_'(&(`256)~/x)^x:(`256),x}

O:`<       /nathanotterness.com/2021/10/tiny_elf_modernized.html
O"[bits 64]\na0: equ 4096 * 40   ;load address
db 0x7f,'E','L','F' ;11
db 2\ndb 1\ndb 1\ndb 0\ndq 0\ndw 2\ndw 0x3e\ndd 1
dq main + a0\ndq phead\ndq 0\ndd 0\ndw 64\ndw 0x38
dw 1\ndw 0x40\ndw 0\ndw 0\n\nphead:
dd 1\ndd 5\ndq 0\ndq a0\ndq a0\ndq end\ndq end\ndq 0x200000\n\n"

O compile all
O"end:\n"

/rax rcx rdx rbx rsp rbp rsi rdi r8..r15

/main:
/ mov rax,97
/ call w
/
/q:
/ mov rax, 60
/ xor rdi, rdi
/ syscall
/
/w:
/ push rax
/ xor rax,rax
/ inc rax
/ mov edi,eax
/ mov rsi,rsp
/ mov edx,eax
/ syscall
/ pop rsi
/ ret

