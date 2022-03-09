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

The output is a single file (Go subset) that can be compiled with Go or [wg](https://github.com/ktye/wg):

```
kom file.k   [repl] > a.go
wg -c -prefix ktye_   a.go > k.c
gcc k.c -lm
```

## how

Kom parses the k source to byte code, than transforms byte code into native code. 
It bundles the interpreter source with the generated code and emits everything as a single source file.

Every lambda function is compiled into a single native function and loaded to k's type system as a special function type (compiled lambda).
Local variables are compiled to native locals and globals remain k variables.

`while`, cond `$[..]` and return `:x` are detected from jumps in byte code and transformed to native control flow.

Constants are stored in native globals and are computed once at startup.

The resulting binary is still a full k-interpreter.
New k code at runtime is interpreted the normal way.
It is not a jit-compiler.

## example
`f:{{x+y}/!x}` is compiled to

```
...
func f_53(x_, y_ K) K {   // {x+y}  f_53 is stored in function table at 434
        rx(y_) // ref
        rx(x_)
        k4 := Add(x_, y_) // call k primitive directly
        dx(x_) // unref   (it still contains unnecessary refcounting)
        dx(y_) 
        return k4
}
func f_54(x_ K) K { // {{x+y}/!x}
        rx(x_)
        k2 := Til(x_)      // !x
        k3 := lmb(434, 2)  // assign compiled lambda (f_53) as a k value
        k4 := rdc(k3)      // create reduction
        k5 := Atx(k4, k2)  // apply reduction
        dx(x_)
        return k5
}
...
```

## limits

- no dynamic scope
- cannot print compiled lambda     `${x+y}`
- cannot decompose compiled lambda `.{x+y}`
- envcalls are not supported        `{...}.d` or `t{..}` (table-where)
- unpack in `z.k` uses dynamic scope

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
