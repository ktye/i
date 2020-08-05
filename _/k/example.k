path:,`
exec:{$[#$*p:path;disp walk p;()]}

walk:{[p];*.,,*p} /todo index

disp:{$[1~@x;x;7~@x;dis7 x;kstr'x]}
dis7:{k:$!x;pad:rpad[;|/#'k]; ((pad'k),'"|"),' 'k'.x}
rpad:{y#x,(y-#x)#" "}                                               /rpad["alpha";10]
kstr:{$[@x~1;x; 'kx]}

/

tags:"List Dict"

List:("alpha";"beta";"gamma")
Dict:`a`b!(1 2;3 4 5;`symbol)


"this is a listbox-ui example application

application source: ktye.github.io/example.k
listbox source    : ktye.github.io/listbox.html"