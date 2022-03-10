# kom compiles k to static code

## build

```
sh mk
go install
```

## use

```
kom [args]
args:
 *.k   compile file.k
 *.go  append native code
 *.t   compile tests
 repl  add a repl loop
```

the output is a single file of a Go subset, that can be compiled with go or [wg](https://github.com/ktye/wg), e.g. to c:

```
kom file.k   [repl] > a.go
wg -c -prefix ktye_   a.go > k.c
gcc k.c -lm
```

## how

Kom parses the k source to byte code, than transforms byte code into native code. 
It bundles the interpreter source with the generated code and emits everything as one file.

Every lambda function is compiled into a native function and loaded to k's type system as a special function type (compiled lambda).
Local variables are compiled to native locals and globals remain k variables.

`while`, cond `$[..]` and return `:x` are detected from jumps in byte code and transformed to native control flow.

Constant arrays are stored in native globals and are computed once at startup.
Both constants and lambdas are are generated only once if they appear multiple times.

The resulting binary is still a full k-interpreter.
New k code at runtime is interpreted the normal way.
It is not a jit-compiler.

## examples
`f:{{x+y}/!x}` is compiled to

```
...
func f_53(x_, y_ K) K {    // {x+y}  f_53 is stored in function table at 434
        rx(y_)             // ref
        rx(x_)
        k4 := Add(x_, y_)  // call k primitive directly
        dx(x_)             // unref
        dx(y_) 
        return k4
}                          // => there is some unnecessary refcounting
func f_54(x_ K) K {        // {{x+y}/!x}
        rx(x_)
        k2 := Til(x_)      // !x
        k3 := lmb(434, 2)  // assign compiled lambda (f_53) as a k value
        k4 := rdc(k3)      // create reduction
        k5 := Atx(k4, k2)  // apply reduction
        dx(x_)
        return k5
}
```

```
func f_53(x_, y_ K) K {         // {while[x;y]}
        k2 := K(0)              // return value for the following loop       
        for {
                rx(x_)
                dx(x_)
                if int32(x_) == 0 {
                        break
                }
                rx(y_)
                dx(k2)
                k2 = y_
        }
        dx(x_)
        dx(y_)
        return k2
}
```

```
func f_53(x_, y_, z_ K) K {  // {$[x;y;z;x;y]}
        rx(x_)
        k4 := K(0)           // return value for cond
        dx(x_)
        if int32(x_) != 0 {
                rx(y_)
                k4 = y_
        } else {
                rx(z_)
                k5 := K(0)   // $[a;b;c;d;e] is $[a;b;$[c;d;e]]
                dx(z_)
                if int32(z_) != 0 { 
                        rx(x_)
                        k5 = x_
                } else {
                        rx(y_)
                        k5 = y_
                }
                k4 = k5
        }
        dx(x_)
        dx(y_)
        dx(z_)
        return k4
}
```

## limits

- no dynamic scope
- cannot print compiled lambda     `${x+y}`
- cannot decompose compiled lambda `.{x+y}`
- envcalls are not supported        `{...}.d` or `t{..}` (table-where)
- unpack in `z.k` uses dynamic scope
- error positions are not tracked

## benchmark

bench/ contains three benchmarks which compare compiled vs interpreted code using the c backend for both (gcc -O3).

```
           |inc      over     qr
-----------|--------------------------
compiled   |2.869999 5.009999 6.589999
interpreted|6.48     10.57    6.639999
ratio      |0.442901 0.473982 0.992469

values are time(user) best of 3.
```

program | what
--- | ---
inc | `f:{i:0;while[i<x;i+:1]}`
over | `f:{{x+y}/!x}`
qr  | solve a 10000 x 300 random matrix with 10000 x 100 rhs

solve is built-in `z.k`

```
solve:{qslv:{H:x 0;r:x 1;n:x 2;m:x 3;j:0;K:!m
 while[j<n;y[K]-:(+/(conj H[j;K])*y K)*H[j;K];K:1_K;j+:1]
 i:n-1;J:!n;y[i]%:r@i
 while[i;j:i_J;i-:1;y[i]:(y[i]-+/H[j;i]*y@j)%r@i]
 n#y}
 q:$[`i~@*|x;x;qr x];$[`L~@y;qslv/:[q;y];qslv[q;y]]}

qr:{K:!m:#*x;I:!n:#x;j:0;r:n#0a;turn:$[`Z~@*x;{(-x)@angle y};{x*1. -1@y>0}]
 while[j<n;I:1_I
  r[j]:turn[s:abs@abs/j_x j;xx:x[j;j]]
  x[j;j]-:r[j]
  x[j;K]%:%s*(s+abs xx)
  x[I;K]-:{+/x*y}/:[(conj x[j;K]);x[I;K]]*\:x[j;K]
  K:1_K;j+:1];(x;r;n;m)}
```

### how array-like is the program?
assumption: scalar code profits from compilation while large vectors don't care.

[kvc](../go/mk) is a special build of k that prints vector sizes `#x` of the accumulator after each vm instruction.  

statistics of `k qr.k -e 'f 1'` using rhs:1 instead of 100:  

- 500 000 vm instructions with
- 64 % of instructions are scalar  1~#x
- 34 % are large 100<#x
- with a maximum length of 603


### k vs lapack
`qr.k` runs at the same speed compiled or interpreted. How does that compare to pure lapack?  

[lapack](./bench/lapack.c) standard build with refblas is 5 times faster.

## todo
refcount elimination 
