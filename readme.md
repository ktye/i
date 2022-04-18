# Ref

```
src:   github.com/ktye/i          6#k               (public domain)
jit:   wavm run --enable multi-memory   --enable exception-handling 
                --mount-root . k.wat [args]  (github.com/WAVM/WAVM)
build: go install ktye.github.com/wg/cmd/wg  (compile wasm from go)
       wg . > k.wat         (wat2wasm: github.com/WebAssembly/wabt)
       wat2wasm --enable-bulk-memory --enable-simd k.wat  -o k.wasm
online:ktye.github.io  (wasm needs multi-value/bulk-memory/simd128)
c/extend: +/k.c                           (header-library "ktye.h")
args:     data file.k file.t                  (assign/execute/test)
web:      drop file assigns                      `file<c (download)

repl:  \\(exit) \t1000 expr(bench) \m(totalmem) \c(clear+checkrefs)

+ flp add       '  ech  bin   x y' lin     b bool     1b    110b   
- neg sub       /  rdc  mod   f n/ ndo     c char     "x"   "ab"        
* fst mul       \  scn  div   x//y dec     i int      2     1 2         
% sqr div       ': ecp  in    y\\x enc     s symbol   `a    ``c`d       
! til key       /: ecr  split   f/:fix     f float    2.    1. 2.    
& wer min       \: ecl  join    f\:fix     z complex  1a    1a20 2a     
| rev max       while[c;a;b;..]            L list     (1;2 3)      
< asc les       $[a;b;...]      cond       D dict     `a`b!1 2
> dsc mor       @[x;i;+;y]      amend      T table    +`a`b!.. 
= grp eql       .[x;i;+;y]      dmend      v verb     +
~ not mtc       {a+b}.d         env        c comp     1+/*%  
, enl cat       k?t             group      d derived  +/
^ srt cut       k!t             key        l lambda   {x+y}
# cnt tak       t,d t,t t,'t(h) join       x native   c-extension
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
[k.wasm](https://github.com/ktye/i/releases/download/latest/k.wasm)|webassembly binary module||
[k-wavm.wat](https://github.com/ktye/i/releases/download/latest/k-wavm.wat)|webassembly text format (with try/catch)|`wavm run --enable multi-memory --enable exception-handling --mount-root . k-wavm.wat`|
[k.go](https://github.com/ktye/i/releases/download/latest/k.go)|bundled package k|`go build k.go`|
[kg.go](https://github.com/ktye/i/releases/download/latest/kg.go)|main program|`go build kg.go`|
[k.f](https://github.com/ktye/i/releases/download/latest/k.f)|fortran|`gfortran k.f`|
[k.c](https://github.com/ktye/i/releases/download/latest/k.c)|c|`gcc k.c -lm`|
[ktye.h](https://github.com/ktye/i/releases/download/latest/ktye.h)|single header library|see [k+](https://github.com/ktye/i/tree/master/%2B)|
[k+.tar.gz](https://github.com/ktye/i/releases/download/latest/k%2B.tar.gz)|k + extensions|linux: `sh mk.lin`|

  
