# modus interpreti
## current state
kvm instructions are stored in a flat k list (not in byte code).
the list contains k values (e.g. constants) and primitives.
primitives are `+-*..` monadic and dyadic and special vm operations (drop, jump, ..)

evaluation from chars of k-code to result is: `r=exec(parse(tok(x)))`.

parse returns a kvm instruction list and exec executes it from start to end.

## limits
the first limit is the stack size which is fixed and rather small.

another limit is that exec is called recursively when evaluating lambda functions.
their code is stored in the function value and calling a lambda function executes this code.
the call itself can occur anywhere. the vm cannot know it will call a lambda function.

## inband control flow
`while` and `cond` are represented with jumps in the byte code. they could also be external functions (like lambdas)

this could simplify both parser and execution.

e.g.
```
// a while object is a list with two instruction lists: l2(condition, body)
func whl(x K) K {
  c := x0(x)
  b := r1(x)
  r := K(0)
  for {
    p := exec(rx(c))
    dx(p)
    if int32(p) != 0 {
      break
    }
    dx(r)
    r = exec(rx(b))
  }
  dx(c)
  dx(b)
  return r
}
```

## tail calls/stateless execution
the current exec uses an accumulator and runs over the list.
a stateless executor would store all information in global variables (or heap), e.g.

### instruction pointer, return stack
instruction lists are still distributed in k variables.
the ip could contain the list value and the index, or point directly to the address of the current instruction.
in the second case it is int32 and all instruction lists must be terminated with a return instruction.

the return stack could be an int vector. the top position is the last value.
to prevent too many relocations when using k vectors as stacks, they could be preinitialized with dummy values to have a certain initial length.


```
// global
IP int32  //instruction pointer
RS K      //return stack (ints)
ST K      //value stack, k values

func step() int32 {
  for {
    x := K(I64(IP)) //fetch current instruction
    if x == RET {
      rest() // explained later
      IP = ipop(RS)
      if IP == END {
        return 0 //halt
      } else {
        break
      }
    }
  }
  t := tp(x)
  //kpush t if value
  //switch if verb, similar as in exec.go buth without the accumulator
  IP += 8
  return 1
}
```

## lambda call
calling a lambda function does 4 things (ktye/k uses dyanamic scope): `func lambda(x K) K {..}`
save shadowed variables(the local list including args), assign arguments, execute instructions, restore variables, unref the function.

it must be rewritten to return before execution. e.g.:

```
func lambda(x K) K {
  loc := x1(x) //local list (symbols)
  // save values of the symbols in loc in a k-list
  // kpush that list to the value stack
  // kpush loc
  // assign arguments
  code := x2(x) //get instructions
  dx(x) //unref lambda
  RS = rpush(IP)
  IP = int32(code)
}
```

when step sees a return instruction it calls `rest()` which restores the variable list to the state before calling the lambda function
```
func rest() {
  syms := kpop()
  vlist := kpop()
  // assign each value in vlist to the corresponding symbol in syms.
}
```

## exec
instead of `r=exec(parse(tok(x)))` we do
```
  r := parse(tok(x))
  rpush(IP)
  IP = int32(r)
  for step() {}
  r = kpop()
```

