
sel :!0
path:(`Tree;1)
exec:{$[@f:lup x;$[#path::x;disp walk path;1_,tags::Tags];f@0]}
/exec:{$[#path::x;disp walk path;1_,tags::Tags]}
slct:{sel::x}
walk:{[p];$[1~#p;*.,,*p;.{(x;y)}/p]}
push:{exec path::path,,sel}
drop:{exec@-1_path}

disp:{t:@x;$[1~t;x;7~t;dis7 x;6~t;kstr'x; 'kx]}
dis7:{k:$!x;pad:rpad[;|/#'k]; ((pad'k),'"|"),' 'k'.x}
rpad:{y#x,(y-#x)#" "}
kstr:{$[1~@x;x; 'kx]}

lup: {*.,,x}

/

Tags:"List Dict Ints Floats Symbols Tree Back"    /application default (or 0-level)
tags:Tags                                         /tags for current path

List:("alpha";"beta";"gamma")
Dict:`alpha`b!(1 2;3 4 5;`symbol)
Ints:10-!8
Floats:129 'F!10
Symbols:`abc`d`efghi
Tree:(`alpha`beta`gamma!(1 2 3;`a`b`c!1 2 3;("first line";"second line")))

Back:drop

"this is a listbox-ui example application

application source: ktye.github.io/example.k
listbox source    : ktye.github.io/listbox.html"