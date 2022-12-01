/ki parses i source and writes the IR as a k table
/see ik.k for a description of i.

kiC:"";
kix:{y:@[y;1;+;#*x];y[1]:0|y 1;x,'y} /catentate subtrees and translate parents, 0N->0
ki0:(0#`;!0;!0;0#`)

kitab:{i:&"t"~/:x;s:x i-2; {(`tab;0;x;y)}'[,/(x i-1)+!'#'s;,/s]}
kivar:{v:x -2+i:&(`64)~/:x;s:x -1+i;{(`var`lit;0N 1;0N,kival y;(`$-1_c),`$*|c:$x)}'[s;v]}
kicon:{v:x -3+i:&(`74)~/:x;s:x -2+i;{(`con`lit;0N 1;0N,kival y;(`$-1_c),`$*|c:$x)}'[s;v]}
kival:{$[`i~t:@x;x;[kiC,:`c x:$[`f~t;`c@,x;x];:-8+#kiC]]}
kilam:{i:&"l"~/:x; kifun'[x i-1;x i-6;x i-8]}
kilan:{i:&"n"~/:x; kifun'[x i-1;0#`;x i-7]}

kifun:{[s;r;f]f:.f;(,`fun;,0N;,0;,s)kix/(kiarg f;($[#r:kires'$r;+r;ki0]);kiloc f;kiast[f;$r])}
kiarg:{a:(n:x 3)#x 1;      j:0N 0 i:2/!2*n;(`arg`sym@i;,/+|2^(2*n)#2*!n;j;,/+2^(kits'a),a)}
kiloc:{n:#l:(,`)_(x 3)_x 1;j:0N 0 i:2/!2*n;(`loc`sym@i;,/+|2^(2*n)#2*!n;j;,/+2^(kits'l),l)}
kires:{(`res;0N;0N;`$x)}
kiast:{x:*x;t:$[1~#y;`$y;`]
 x:$[1<#y;-2_x;x] /multi-res
 r:()kiprs/x
 
 r:$[(0~#y)&(`cal~***|r);$[*|kis(*|r)3;(-1_r),,(`drp;0N;0N;`)kix/,*|r;r];r] /drop last call if niladic func

 
 r:$[#y;((-#y)_r),,(`ret;0N;0N;t)kix/|(-#y)#r;r]
 r:(`ast;0N;0N;`)kix/r}


kiprs:{t:@y;
 f2:{kix[kix[(x;0N;2;kits(*|y)3);*|y];y@-2+#y]}
 f:(0;kiasn;f2`add;f2`sub;f2`mul;0)@0|(0;`64;`66;`67;`68)?y
 $[~(y~())|f~0;(-2_x),,f x;y~`83;kiatx x;y~`84;kifnc x;`s~t;x,,(`sym;0N;0;y);`S~t;x,,(n#`sym;n#0N;&n:#y;kidrp y);`i~t;x,,(`lit;0N;y;`i);x]}
 

kiatx:{$[`sym~**|x;(-2_x),,kix[$[(s:(*|x)3)':`i`u`j`k`f;(`cst`typ;0N 0;0N 0N;s,kits a 3);(`cal;0N;0N;s)];a:x@-2+#x];x]}
kifnc:{$[`sym~**|x;kical x;kicli x]}
kical:{s:kis f:(*|x)3;a:*s;r:s 1;c:(`cal;0N;0N;f);$[a;((-2-a)_x),,(c kix/|a#-2_x);((-1-a)_x),,c]}
kicli:{n:#x;an:#a:$(x n-2)3;r:+(rn#`res;rn#0N;(rn:#r)#0N;r:`$'$(x n-3)3); e:(x n-4); x:-4_x; b:+(an#`arg;an#0N;an#0N;`$'a);((-an)_x),,((((`cli;0N;an;`)kix/,e)kix/(-an)#x)kix/b)kix/r}


kiasn:{kix[kix[(`asn;0N;#(*|x)3;`);*|x];x@-2+#x]}
kidrp:{@[x;&`=x;`"__"]}

kits:{`$*|$x}
kiel:{i@&`v=@'x[-7+i:&(`64)~/:x]}    /lambda assignments -> `64->"l"
kiev:{i@&`v=@'x[-6+i:&(`64)~/:x]}    /lambda (void)      -> `64->"n"
kiet:{i@&`i=@'x[-1+i:&(`64)~/:x]}    /table assignments  -> `64->"t"

kis:`!()                             /signature          -> `f`g!(2 0;0 2)  (#arg;#res)
kisig:{r:(*r)!+1_r:(x i-1;{(.x)3}'x i-8;#'$x -6+i:&"l"~/:x);r,(x i-1)!+({(.x)3}'x -7+i;(#i:&"n"~/:x)#0)}



kiget:{f:&`fun=*x;s:x[3;i:1+&(`loc=*x)|`arg=*x];fi:=f'i
 a:&`sym=*x
 a:a@&~(x[0;x[1;a]])':`arg`loc`asn
 fa:f'a
 g:0N=(s fi fa)?'x[3;a]
 x[0;a]:`get`Get g
 x[2;a]:0N}

ki:{kiC::"";x:`p@x;x:@[x;(kiel;kiev;kiet)@\:x;"lnt"]; kis::kisig x
 D:$[#i:&`C=@'x;,/x i;""]
 l:(`prg`mem;0 0N;0 1;(*x),`a)
 l:-1_'l /delete mem, todo rm
 l:{x kix/y}/[l;(kitab;kivar;kicon;kilam;kilan)@\:x]
 s:s!`$-1_'$s:?(l 3)@i:&`sym=*l;l[3;i]:s l[3;i] /remove types from sym
 l:kiget l                                      /sym->get|Get
 `C`D`T`P`I`S!(kiC;D),l}

kixx:{d:ki x;`<"\n"/:(`lxy 40 100)@+`T`P`I`S!-4#.d}

/x:kixx"`cast /program
/
//functions
/f[j]:{[xi]`j xi}
/"
