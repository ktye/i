/ki parses i source and writes the IR as a k table
/see ik.k for a description of i.

kiC:"";
kix:{y:@[y;1;+;#*x];y[1]:0|y 1;x,'y} /catentate subtrees and translate parents, 0N->0

kitab:{i:&"t"~/:x;s:x i-2; {(`tab;0;x;y)}'[,/(x i-1)+!'#'s;,/s]}
kivar:{v:x -2+i:&(`64)~/:x;s:x -1+i;{(`con`lit;0N 1;0N,kival y;(`$-1_c),`$*|c:$x)}'[s;v]}
kicon:{v:x -3+i:&(`74)~/:x;s:x -2+i;{(`var`lit;0N 1;0N,kival y;(`$-1_c),`$*|c:$x)}'[s;v]}
kival:{$[`i~t:@x;x;[kiC,:`c x:$[`f~t;`c@,x;x];:-8+#kiC]]}
kilam:{i:&"l"~/:x; kifun'[x i-1;x i-6;x i-8]}

kifun:{[s;r;f]f:.f;(,`fun;,0N;,0;,s)kix/(kiarg f;(+kires'$r);kiloc f;kiast[f;$r])}
kiarg:{a:(n:x 3)#x 1; j:0N 0 i:2/!2*n; (`arg`sym@i;j;j;,/+2^(kits'a),a)}
kiloc:{n:#l:(,`)_(x 3)_x 1;j:0N 0 i:2/!2*n; (`loc`sym@i;j;j;,/+2^(kits'l),l)}
kires:{(`res;0N;0N;`$x)}
kiast:{x:*x;t:$[1~#y;`$y;`]
 x:$[1<#y;-2_x;x] /multi-res
 r:()kiprs/x
 r:((-#y)_r),,(`ret;0N;0N;t)kix/|(-#y)#r
 r:(`ast;0N;0N;`)kix/r}


kiprs:{t:@y;
 f2:{kix[kix[(x;0N;2;kits(*|y)3);*|y];y@-2+#y]}
 f:(0;kiasn;f2`add;f2`sub;f2`mul;kiatx;0)@0|(0;`64;`66;`67;`68;`83)?y
 $[~f~0;(-2_x),,f x;y~`84;((-2+*|r)_x),,-1_r:kical x;`s~t;x,,(`sym;0N;0;y);`S~t;x,,(n#`sym;n#0N;&n:#y;kidrp y);`i~t;x,,(`lit;0N;y;`i);x]}
 

kiatx:{kix[(`cal;0N;0N;(*|x)3);x@-2+#x]}
kical:{c:(`cal;0N;0N;(*|x)3);$[()~p:x[-2+#x];c,0;(c kix/n#-2_x),n:-p 2]}
kiasn:{kix[kix[(`asn;0N;#(*|x)3;`);*|x];x@-2+#x]}
kidrp:{@[x;&`=x;`"__"]}

kits:{`$*|$x}
kiel:{i@&`v=@'x[-7+i:&(`64)~/:x]}    /lambda assignments -> `64->"l"
kiev:{i@&`v=@'x[-6+i:&(`64)~/:x]}    /lambda (void)      -> `64->"v"
kiet:{i@&`i=@'x[-1+i:&(`64)~/:x]}    /table assignments  -> `64->"t"


kiget:{f:&`fun=*x;s:x[3;i:1+&(`loc=*x)|`arg=*x];fi:=f'i
 a:&`sym=*x
 a:a@&~x[0;x[1;a]]':`arg`loc`asn
 fa:f'a
 g:0N=(s fi fa)?'x[3;a]
 x[0;a]:`get`Get g
 x[2;a]:0N}

ki:{kiC::"";x:`p@x;x:@[x;(kiel;kiev;kiet)@\:x;"lvt"]
 D:$[#i:&`C=@'x;,/x i;""]
 l:(`prg`mem;0 0N;0 1;(*x),`a)
 l:-1_'l /delete mem, todo rm
 l:{x kix/y}/[l;(kitab;kivar;kicon;kilam)@\:x]
 s:s!`$-1_'$s:?(l 3)@i:&`sym=*l;l[3;i]:s l[3;i] /remove types from sym
 l:kiget l                                      /sym->get|Get
 `C`D`T`P`I`S!(kiC;D),l}

kixx:{d:ki x;`<"\n"/:(`lxy 40 100)@+`T`P`I`S!-4#.d}

/x:kixx"`blank /program
/
/
//functions
//f[ii]:{[](,1;,2)}
/g[i]:{[]`xi`:f[];xi}
/"
