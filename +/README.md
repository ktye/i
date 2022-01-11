# k+

k+ is a c build of k + wrappers around external c libraries.
the interface is the c-api:

# c-api (k.h)

the idea is to have a common k to c interface for multiple implementations of k.
currently ngn/k and ktye/k implement the interface.

those who write bindings to c libraries only need to know how the interface works without having to learn all implementation details of the interpreter.

the c-api is defined in [k.h](k.h)

# ktye.h

ktye.h is the c version of ktye/k.
It is generated from the source files in the toplevel directory using a custom compiler [ktye/wg](https://github.com/ktye/wg).
The generated c code is pretty large, the expressions are converted to many single assignments but the c compiler should be able to tranform this into acceptable binaries.

It should be portable and uses stdio, stdlib(malloc), setjmp(error recovery) and gcc portable vector instructions for simd128.

k.c is the main application that includes all extensions.

# building

the file mk builds the k+ binary for windows. it is tested with mingw from [web64devkit](https://nullprogram.com/blog/2020/09/25/).
other systems need different link flags.

# libraries

|dir|what|source|
|---|---|---|
|[mat](mat/mat.c)|solve,qr,eig,svd real+complex|lapack|
|[img](img/img.c)|2d vector drawing,rw png,r svg,r ttf|[nothings](https://github.com/nothings/stb),[nanosvg](https://github.com/memononen/nanosvg)|
|[ray](ray/ray.c)|show image in a window with interaction|[raylib](https://www.raylib.com/)|
