# html/js version

- k.js is a js module that loads k.wasm and provides an interface similar to the c-api ../+/k.h
- k.wasm is not checked in. it's the core that is built by [wg](https://github.com/ktye/wg) from ../*.go
- k6.html is a minimal web page that demonstrates it's usage
- k.html (to be reworked) is the front page on ktye.github.io

## js interface
The js provides an api interface to k that does not e.g. include a filesystem or a repl.
k.wasm uses calls compatible to the webassembly system interface (wasi) for io/args/clock.
This allows it to be used out of the box by compatible runtimes such as [wavm](https://github.com/WAVM/WAVM).

### initialization
`import{K}from './k.js'`
As compiling a webassembly module is usually done asynchronously, the initialization takes an callback f that is invoked afterwards:
```
K.kinit(f, read, write) 
```
read and write are js functions that are called from k for file io.
E.g. an application could trigger a download whenever k writes a file.
`write("", output_string)` is called when k wants to write to stdout, e.g. in a repl.

### k values
k values in ktye/k are 64bit integers, which are BigInt in js.
`TK(x)` returns the type which is

- "c", "C" char(s)
- "i", "I" 32 bit int(s)
- "s", "S" symbol(s)
- "f", "F" 64 bit floats(s)
- "F" is also returned for complex in ktye/k
- "L" general list
- "D", "T" dict and table
- undefined otherwise, e.g. functions

`NK(x)` returns the number of elements for arrays and should only be called on those.

### create K values (K from js)
Kc(x) /char atom, x is numeric or string(first char)
Ks(x) /symbol from string
Ki(x) /integer from a number
Kf(x) /float64 from a number

KC(x) /chars, x is a string or utf8array
KS(x) /symbols from strings
KI(x) /ints from numeric array or Int32Array
KF(x) /floats from numeric array of Float64Array
KL(x) /k list from a js array of k values

cK(x) /returns a number -128..127
CK(x) /returns a string
iK(x) /returns a number int32 range
IK(x) /returns an Int32Array


### example
serve k6.html, k.js and k.wasm locally. open k6.html in a browser and type on the js console:

```
> x=K.KC("alpha")
< -8070450532247893328n
> K.CK(x)
< 'alpha'
```

### execute
```
Kx("+/", x, y) 
```
