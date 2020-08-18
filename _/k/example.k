path:()
exec:{$[0~#x;1_,tags::Tags;0~@f:.x;disp f walk path;disp walk $[5~@x;path::,x;x]]}
walk:{[p];$[0~#p;();1~#p;*.,,*p;.{(x;y)}/p]}
push:{x:$[7~@d:walk path;$[1~#s:(!d)x;,s;s];x];exec path:: \path,,x}
back:{tags^:`ok;exec path::-1_path}
edit:{tags,:`ok;  'k walk     path}
ok:  {tags^:`ok;      asn[path; \x];}

disp:{t:@x;$[1~t;x;ist x;dis8 x;6~t; 'k'x;7~t;dis7 x; 'kx]}
dis7:{k:$!x;pad:rpad[;|/#'k]; ((pad'k),'"|"),' 'k'.x} /display dict
dis8:{n:1+|/'#''t:&kst''(,$!*x),.'x;,/'&rpad''[t;n]}  /display table
rpad:{y#x,(y-#x)#" "}
kst: {$[1~@x;x; 'kx]}
ist: {$[~6~@x;0;~|/7~'@'x;0;&/~':!'x;1;0]}             /is table

data:`!()
drop:{data,:x!,y;"data<-",$x}
asn: {.(:;x;$[6~@y;$[1~#x;,y;,,,y];y])}

/

Tags:`List`Dict`Ints`Floats`Symbols`Tree`Table`draw`back`edit
tags:Tags                                             /tags for current path
tag: {" "/:$?tags}                                    /tags may be updated dynamically based on path

List:("alpha";"beta";"gamma")
Dict:`alpha`b!(1 2;3 4 5;`symbol)
Ints:10-!8
Floats:129 'F!10
Symbols:`abc`d`efghi
Tree:(`alpha`beta`gamma!(1 2 3;`a`b`c!1 2 3;("first line";"second line")))
Table:(`abc`def`g!("x";2;3.0);`abc`def`g!("y";3;2.1);`abc`def`g!("zz";1;2.3))

draw:{30 'd 'r1500;("one";"two";"three")}

"this is a listbox-ui example application

double-click on words in the tag bar; then select, navigate or edit.
ESC is Back

application source: ktye.github.io/example.k
listbox source    : ktye.github.io/listbox.html
k.wasm binary     : ktye.github.io/k.wasm
k source          : https://raw.githubusercontent.com/ktye/i/master/k.w"