\listbox.k

Tags:`List`Dict`Ints`Floats`Symbols`Tree`Table`draw`edit
tags:Tags                                             /tags for current path
tag: {" "/:$?tags}                                    /tags may be updated dynamically based on path

List:("alpha";"beta";"gamma")
Dict:`alpha`b`c!(1 2;3 4 5;`symbol)
Ints:10-!8
Floats:129 'F!10
Symbols:`abc`d`efghi
Tree:(`alpha`beta`gamma!(1 2 3;`a`b`c!1 2 3;("first line";"second line")))
Table:`abc`def`g`h`s!(`x`y`zz;9+!3;"ABC";(`a`b!1 2;`c`d!5 6;`a`d!7 9);("abc";"def";"ghijk"))
T:`a`b!(1 2 3;1.1 2.2 3.3)

draw:{30 'd 'r1500;("one";"two";"three")}


"this is a listbox-ui example application

double-click on words in the tag bar; then select, navigate or edit.
ESC is Back

application source: ktye.github.io/example.k
listbox source    : ktye.github.io/listbox.html
k.wasm binary     : ktye.github.io/k.wasm
k source          : https://raw.githubusercontent.com/ktye/i/master/k.w"
