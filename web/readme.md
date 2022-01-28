# html/js version

- k.js is a js module that loads k.wasm and provides an interface similar to the [c-api k.h](../+/k.h)
- k.wasm is not checked in. it's the core that is built by [wg](https://github.com/ktye/wg) from `../*.go` or `wget ktye.github.io/k.wasm`
- k6.html is a minimal web page that demonstrates it's usage
- k.html (to be reworked) is the front page on ktye.github.io
- l.wasm (see l.wat) is the link module, that imports a js function and provides it as a wasm function

## js interface
`k.wasm` is the k core, which does not include a filesystem or a repl.
It k.wasm uses calls compatible to the webassembly system interface (wasi) for io/args/clock.
This allows it to be used out of the box by compatible runtimes such as [wavm](https://github.com/WAVM/WAVM) as a k repl.  
`k.js` (a javascript module) provides an api to k, which allows custom read/write and extension functions.

### initialization
In js: `import {K} from './k.js'`  
As compiling a webassembly module is usually done asynchronously, the initialization takes an callback `f` that is invoked afterwards:
```
K.kinit(ext)
```
`ext` is a js object that contains js extension functions to be available in k.
See [external functions](#external-function-call-js-from-k).


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

Without arguments, `f` is evaluated as an expression, but not called.

### Assignment / lookup
`KA(sym,val)` assigns. Both `sym` and `val` are k objects.
To lookup a variable use the equivalent to `. sym`:

```
K.KA(K.Ks("name"), K.Ki(314))         // KA consumes arguments, no return value.
let i = K.iK(K.Kx(".", K.Ks("name"))) // lookup. In general we have to query the type.
```

### external function (call js from k)

The object passed to `K.kinit(ext)` contains js functions that are callable from k.  
There are 3 special functions:
- `init() ` is called after initializing k, which happens asynchronously
- `read(x_string)` is called when k reads a from a file: ``` <`name ```. It should return a Uint8Array.
- `write(x_string,y_uint8array)` is called when k wants to write a file: ``` `name<chars ```, e.g. to trigger a download.  
   `x_string` is empty when k wants to write to _stdout_, as in `1+ \2`. There is no _stdin_ equivalent.

All other functions functions in `ext` are called by k.wasm with the signature `i64:i64`.  
Both are k values. The argument is a general list with all arguments.
The list's length always matches the number of arguments of the provided function signature, e.g. 3 for ```ext.f=function(x,y,z){return 0}```.  
The extension functions (except for init, read and write) have to use the k-api (via global `K`), consume their arguments and return a k value.  
They are stored as k variables under their name.

_Implementation:_ k.wasm only imports wasi functions, no user defined js functions.
However it exports it's function table.
K calls the indirect function at index 98 to do a native/extern function call.
We cannot set the index with a js function directly. Only a wasm function can be put there by the js loader.
That's why there is another module _l.wasm_ whose sole purpose is to provide a wasm function that calls an imported js function.  


### reference counting
Functions consume their arguments.

- `x=ref(x)` increase refcount, return x
- `unref(x)` decrease refcount or free memory

For chain calls use, e.g.: `f(ref(x), x)`

### error save restore
All api calls are unprotected.
k errors trigger wasm traps and memory is in undefined state.
ktye/k does not have error values and does not rewind the state on errors.
Instead you can take snapshots and recover.
```
try      { ... ; K.save() }     // save copies all k wasm memory to a back buffer
catch(e) { ... ; K.restore() }  // restore k memory from back buffer
```

### example
Serve `k6.html`, `k.js`, `k.wasm` and `l.wasm` locally. Open k6.html in a browser and type in the js console:

```
> x=K.KC("alpha")
< -8070450532247893328n  // a k value returned as a BigInt
> K.CK(x)
< 'alpha'
```

see [k6.html](k6.html) for an example providing custom read/write functions and an extension (pitimes).

### libraries
(draw.js)[draw.js] extends `k.js` with additional functions.
See (k.html)[k.html] how it is used:
```
import { K } from './k.js'
import { D } from './draw.js'
...
var ext = {
 init: function()         {...},
 read: function(file)     {...},
 write:function(file,data){...}},
}
Object.assign(ext, D)              //this adds draw.js functions to the import object
K.kinit(ext)
```

`draw` adds functions `draw`, `show`, `showev` which are compatible to the [c-versions](../+/draw/readme.md)

