/ki parses i source and writes the IR as a k table
/see ik.k for a description of i.

kiC:"";
kix:{y:@[y;1;+;#*x];y[1]:0|y 1;x,'y} /catentate subtrees and translate parents, 0N->0

kitab:{i:&"t"~/:x;s:x i-2; {(`tab;0;x;y)}'[,/(x i-1)+!'#'s;,/s]}
kivar:{v:x -2+i:&(`64)~/:x;s:x -1+i;{(`con`lit;0N 1;0N,kival y;(`$-1_c),`$*|c:$x)}'[s;v]}
kicon:{v:x -3+i:&(`74)~/:x;s:x -2+i;{(`var`lit;0N 1;0N,kival y;(`$-1_c),`$*|c:$x)}'[s;v]}
kival:{$[`i~t:@x;x;[kiC,:`c x:$[`f~t;`c@,x;x];:-8+#kiC]]}
kilam:{i:&"l"~/:x; kifun'[x i-1;x i-6;x i-8]}

kifun:{[s;r;f]f:.f;(,`fun;,0N;,0;,s)kix/(kiarg f;(+kires'$r);kiloc f;kiast f)}
kiarg:{a:(n:x 3)#x 1; i:2/!2*n; (`arg`sym@i;0N 0 i;,/+2^(!n),&n;,/+2^(kits'a),a)}
kiloc:{n:#l:(x 3)_x 1;i:2/!2*n; (`loc`sym@i;0N 0 i;,/+2^(!n),&n;,/+2^(kits'l),l)}
kires:{(`res;0N;0;`$x)}
kiast:{x:*x
 r:()kiprs/x
 l:-1+#r
 r[l]:kix[(`ret;0N;0N;kits@*(r l)3);r l]
 r:(`ast;0N;0N;`)kix/r}            /todo `ret


kiprs:{t:@y;
 f2:{kix[kix[(x;0N;2;kits(*|y)3);*|y];y@-2+#y]}
 f:(0;kiasn;f2`add;f2`sub;f2`mul;0)@0|(0;`64;`66;`67;`68)?y
 $[~f~0;(-2_x),,f x;`s~t;x,,(`sym;0N;0;y);`i~t;x,,(`lit;0N;y;`i);x]}
 

kiasn:{kix[kix[(`asn;0N;1;`);*|x];x@-2+#x]}        /todo multi-assign(i-value)

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
 x[2;a]:g}

ki:{kiC::"";x:`p@x;x:@[x;(kiel;kiev;kiet)@\:x;"lvt"]
 D:$[#i:&`C=@'x;,/x i;""]
 l:(`prg`mem;0 0N;0 1;(*x),`a)
 l:-1_'l /delete mem, todo rm
 l:{x kix/y}/[l;(kitab;kivar;kicon;kilam)@\:x]
 s:s!`$-1_'$s:?(l 3)@i:&`sym=*l;l[3;i]:s l[3;i] /remove types from sym
 l:kiget l                                      /sym->get|Get
 `C`D`T`P`I`S!(kiC;D),l}

kixx:{d:ki x;`<"\n"/:(`lxy 40 100)@+`T`P`I`S!-4#.d}

x:kixx"`asn /program


/functions
f[i]:{[xi]ri:1+xi;ri}
g[ii]:{[xi](1+xi;1-xi)}
h[i]:{[xi]`xi`ri:g[xi];ri}"
