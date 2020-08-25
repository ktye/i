path:()
exec:{$[x~`edit;edit[];0~#x;1_,tags::Tags;0~@f:.x;disp f walk path;disp walk $[5~@x;path::,x;x]]}
walk:{$[0~#x;();atx[.*x;1_x]]}
push:{exec path::path,,$[ist u:walk path;-1+x^0;x]}
back:{disp walk path::-1_path}
edit:{tako 1;    'k walk path}
ok:  {s:*path;asn[s;$[ist n:.s;tas[n;x];$[#p:1_path;@[n;$[1~#p;*p;p];x];x]]];}

tako:{tags::(tags^`ok),$[x;`ok;0#`]}
disp:{t:@x;tako t<6;$[ist x;dis8 x;6~t; 'k'x;7~t;dis7 x; 'kx]}
dis7:{k:$!x;pad:rpad[;|/#'k]; ((pad'k),'"|"),' 'k'.x} /display dict
dis8:{" "/:'&rpad''[l;|/'#''l:tst''(!x),/:'.x]}       /display table
rpad:{y#x,(y-#x)#" "}
kst: {$[1~@x; x; 'kx]}
tst: {$[5~@x;$x; 'kx]}
ist: {$[~7~@x;0;~#x;0;&/(n:#*.x)~'#'.x;n>1;0]}        /dict is table
atx: {$[#y;atx[$[ist x;tix[x;*y];x@*y];1_y];x]}
tix: {@[;y]'x}                                        /index table x by row y


data:`!()
drop:{data,:x!,y;"data<-",$x}
/asn: {.(:;x;$[6~@y;$[1~#x;,y;,,,y];y])}
asn: {.(:;x;$[6~@y;,y;y])}
tas: { \(`tas;x;y);}

/

Tags:`List`Dict`Ints`Floats`Symbols`Tree`Table`draw`edit
tags:Tags                                             /tags for current path
tag: {" "/:$?tags}                                    /tags may be updated dynamically based on path

List:("alpha";"beta";"gamma")
Dict:`alpha`b`c!(1 2;3 4 5;`symbol)
Ints:10-!8
Floats:129 'F!10
Symbols:`abc`d`efghi
Tree:(`alpha`beta`gamma!(1 2 3;`a`b`c!1 2 3;("first line";"second line")))
Table:`abc`def`g`h!(`x`y`zz;9+!3;"ABC";(`a`b!1 2;`c`d!5 6;`a`d!7 9))

draw:{30 'd 'r1500;("one";"two";"three")}

tags^`ok
tags^`ok
tags^`ok

"this is a listbox-ui example application

double-click on words in the tag bar; then select, navigate or edit.
ESC is Back

application source: ktye.github.io/example.k
listbox source    : ktye.github.io/listbox.html
k.wasm binary     : ktye.github.io/k.wasm
k source          : https://raw.githubusercontent.com/ktye/i/master/k.w"