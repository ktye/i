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
kom file.k > a.go
wg -c a.go > k.c
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

## limits

- no dynamic scope
- cannot print compiled lambda     `${x+y}`
- cannot decompose compiled lambda `.{x+y}`
- encalls are not supported         `{...}.``d   t{..}` (table-where)
- unpack in `z.k` uses dynamic scope

## benchmark
