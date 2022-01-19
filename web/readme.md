# html/js version

- k.js is a js module that loads k.wasm and provides an interface similar to the [c-api k.h](../+/k.h)
- k.wasm is not checked in. it's the core that is built by [wg](https://github.com/ktye/wg) from `../*.go` or `wget ktye.github.io/k.wasm`
- k6.html is a minimal web page that demonstrates it's usage
- k.html (to be reworked) is the front page on ktye.github.io

## js interface
`k.js` provides an api interface to k that does not e.g. include a filesystem or a repl.  
k.wasm uses calls compatible to the webassembly system interface (wasi) for io/args/clock.
This allows it to be used out of the box by compatible runtimes such as [wavm](https://github.com/WAVM/WAVM) as a k repl.

### initialization
In js: `import {K} from './k.js'`  
As compiling a webassembly module is usually done asynchronously, the initialization takes an callback `f` that is invoked afterwards:
```
K.kinit(f, read, write) 
```
read and write are js functions that are called from k for file io.
E.g. an application could trigger a download whenever k writes a file.
`write("", output_string)` is called when k wants to write to stdout, e.g. in a repl.

### k values
k values in ktye/k are 64bit integers, which are BigInt in js.
`TK(x)` returns the type which is

| `TK(x)` returns | atom/vector type |
| --- | --- |
| `"c"`, `"C"` | char(s) |
| `"i"`, `"I"` | 32 bit int(s) |
| `"s"`, `"S"` | symbol(s) |
| `"f"`, `"F"` | 64 bit floats(s) |
| `"F"`        | is also returned for complex in ktye/k |
| `"L"`        | general list |
| `"D"`, `"T"` | dict and table |
| not defined  | otherwise, e.g. functions |



`NK(x)` returns the number of elements for arrays and should only be called on those.

### k from js and js from k
| call | description |
| ---- | ----------- |
|      | **create k atoms** |
|Kc(x) |c(char) from number or first char of a string|
|Ks(x) |s(symbol) from string|
|Ki(x) |i(integer) from number|
|Kf(x) |f(float64) from number|
|      | **create k vectors** |
|KC(x) |C(chars) from string or Uint8Array|
|KS(x) |S(symbols) from array of strings|
|KI(x) |I(ints) from numeric array or Int32Array|
|KF(x) |F(floats) from numeric array or Float64Array|
|KL(x) |L(general list) from a js array of k values|
|      | **js from k atoms** |
|cK(x) |number(-128..127) from c|
|iK(x) |number(int32 range) from i|
|fK(x) |number from f |
|      | **js from k vectors** |
|CK(x) |string from C|
|IK(x) |Int32Array from I|
|FK(x) |Float64Array from F or from complex atom or vector|
|LK(x) |array of k values from L, 2-array from D(dict) and T(table)|
|      | **data pointer**|
|dK(x) |byte index (32bit range) into K.M|

`K.M` is the Webassembly.Memory object, an ArrayBuffer.

### calls
`Kx(f,...args)` takes a string argument `f` and a variable number of k-value arguments.
`f` is evaluated and should return a function value which is called with the supplied arguments. eg:

```
K.Kx("+", K.Ki(1), K.KF([1,2,3]))  /equivalent to 1+1 2 3.0
```

Without arguments `f` is evaluated as an expression, but not called.

### external function (call js from k)
todo..

### refcounting
Functions consume their arguments.

- `x=ref(x)` increase refcount, return x
- `unref(x)` decrease refcount and free memory if possible

For chain calls use, e.g.: `f(ref(x), x)`

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
