# Ref

```
src:   github.com/ktye/i          6#k               (public domain)
build: go install ktye.github.com/wg/cmd/wg  (compile wasm from go)
       wg .   > k.wat       (wat2wasm: github.com/WebAssembly/wabt)
       wat2wasm k.wat                                     -o k.wasm
online:ktye.github.io          (wasm needs multi-value/bulk-memory)
targets: go c fortran wasm web                   (x/ . f77/ . web/)
tests:                                                        (k.t)
c/extend: +/k.c                           (header-library "ktye.h")
kompile:  ahead of time                                      (kom/)
args:     data file.k file.t                  (assign/execute/test)
web:      drop file assigns, plot                `file<c (download)

+ flp add       '  ech  bin   x y' lin     c char     "x"   "ab"   
- neg sub       /  rdc  mod   f n/ ndo     i int      2     1 2    
* fst mul       \  scn  div   x//y dec     s symbol   `a    ``c`d  
% sqr div       ': ecp  in    y\\x enc     f float    2.    1. 2.  
! til key       /: ecr  split   f/:fix     z complex  1a    1a20 2a
& wer min       \: ecl  join    f\:fix     L list     (1;2 3)      
| rev max       while[c;a;b;..]            D dict     `a`b!1 2     
< asc les       $[a;b;...]      cond       T table    +`a`b!..     
> dsc mor       @[x;i;+;y]      amend      v verb     +            
= grp eql       .[x;i;+;y]      dmend      c comp     1+/*%        
~ not mtc       {a+b}.d         env        d derived  +/           
, enl cat       k?t             group      l lambda   {x+y}        
^ srt cut       k!t             key        x native   c-extension  
# cnt tak       t,d t,t t,'t(h) join                               
_ flr drp       t{a>5}          where     exec: t~`v: push         
$ str cst       c:<`file(read)             v:  0..63   monadic     
? unq fnd       `file<c(write)                64..127  dyadic      
@ typ atx       `@i(verb) (+)~`2             128       pop + dyadic
. val cal       .(1;2;`64+(+))  exec         129..255  tetradic    
                                             256       drop        
abs sin cos exp log any find fill            320/384   jmp, jmp-ifz
imag conj angle qr ej avg var std            448..     quoted verb 
solve dot plot hist (unpack) csv                                   
rand: ?n(uniform) ?-n(normal) ?z(binormal) n?n(with) -n?n(w/o) n?L 
```

# Build

|file|what|compile|
|---|---|---|
[k.f](https://github.com/ktye/i/releases/download/latest/k.f)|fortran version|`gfortran k.f`|
[k.c](https://github.com/ktye/i/releases/download/latest/k.c)|c version|`gcc k.c -lm`|
[k.js](https://github.com/ktye/i/releases/download/latest/k.js)|js version|`qjs --std k.js` (quickjs)|
[k.wasm](https://github.com/ktye/i/releases/download/latest/k.wasm)|webassembly binary||
[k.wasi](https://github.com/ktye/i/releases/download/latest/k.wasi)|standalone wasi binary||
[k.go](https://github.com/ktye/i/releases/download/latest/k.go)|go version|`go build k.go`|
[kdb.go](https://github.com/ktye/i/releases/download/latest/kdb.go)|[debugger](d)|`go build kdb.go`|

