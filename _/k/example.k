path:(`Tree;1)
exec:{$[#p:path;disp walk p;1_,tags::Tags]}

walk:{[p];*.,,*p} /todo index

back:{exec path::-1_path}

disp:{t:@x;$[1~t;x;7~t;dis7 x;kstr'x]}
dis7:{k:$!x;pad:rpad[;|/#'k]; ((pad'k),'"|"),' 'k'.x}
rpad:{y#x,(y-#x)#" "}                                               /rpad["alpha";10]
kstr:{$[1~@x;x; 'kx]}

/

Tags:"List Dict Ints Floats Symbols Tree"    /application default (or 0-level)
tags:Tags                                    /tags for current path

List:("alpha";"beta";"gamma")
Dict:`alpha`b!(1 2;3 4 5;`symbol)
Ints:10-!8
Floats:129 'F!10
Symbols:`abc`d`efghi
Tree:(`alpha`beta`gamma!(1 2 3;`a`b`c!1 2 3;("first line";"second line")))


"this is a listbox-ui example application

application source: ktye.github.io/example.k
listbox source    : ktye.github.io/listbox.html"