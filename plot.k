
aax:{am:{`am \x;(!x)!abs@.x};an:{(!x)!angle@.x}
 i:&{$[1~#x;0b;`Z~@$[L~@v:(.x)1;*v;v]]}'x;$[i~!0;x;[r:`r \&1_~':`i \i:^i,!#x;x:x i;`nx \#x;x[r]:am'x r;x[1+r]:an'x 1+r; \x]]}


plot:{n:1;
 /multi:{x[`pw]%:n;x[`px]:(!n)*x`pw;x}
 multi:{x}
 plt::multi@(`plot@)'aax$[(@x)':`L`T;*(x;n:#x);,x]
 Show[draw[`plts@`plt;(PW;PH)];.`"pltcl.";.`"pltzo."]}
 


xy:`alpha`beta!(0.+!10;(sin 0.2*!10;cos 0.2*!10))
po:`abc!,1a30+.1*?10a
aa:`x`y!(0.+!10;*po)
ba:`x`y!(1.+!#i;i:68 56 42 29 0 39 8 61)

/plot xy
/plot(xy;po;ba)
