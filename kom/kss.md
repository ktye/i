# kss (k static subset)
Kom compiled code is faster, because it eliminates the virtual machine.
It is still sub-optimal because it does a lot of checks and branches within the k implementation:
- every primitive checks for type and executes different code paths
- refcount is checked and updated for every vector type

If we defined an annotated subset of k, it can be compiled to more efficient code:

kss are special lambdas:
- arguments have type information in the name (`xI`: int-vector, `xi` int-atom, `xf` float...)
- code is still valid k code and evaluates to the same result when interpreted
- refcount of input vectors is 1 and space may be overwritten
- loops over the same range are combined if possible by the compiler

goal:
- faster vector code
- much faster scalar code, allows to implement primitives in k
- loop fusion not only prevents additional counting overhead, but also saves temporary vector allocations

how:
- k combines type and pointer in a value and stores length, refcount and data in heap memory
- kss compiled code only stores pointer and length in variables, e.g. registers.
- it has static information about type, no refcount and semi-static knowledge about length:
  - with `xI+yI` it knows that `xn~yn`
  - scalar functions e.g. `3*xI+yI` keep their length
  - structural functions update length, e.g. `xi#yI`
- loops are combined and simplified, e.g. within the loop `!#xn` is the loop counter itself, reductions can also be done during the same loop

## example
```
fI:{[xI;yI](!#xI)+xI*yI}
```
compiles with kom to:
```go
/kom ssa
func f(x, y K) (r K) {
 t1 := Mul(rx(x), rx(y))
 t2 := Cnt(rx(x)
 t3 := Til(t2)
 r = Add(t3, t1)
 dx(x)
 dx(y)
 return r
}
```
with lifetime analysis of k variables, some refcounts and intermediate values could be eliminated:
```go
/kom rce
func f(x, y K) (r K) {
 return Add(Til(Cnt(rx(x)),Mul(x,y))
}
```
under the hood these are still 3 loops and needs one more temporary vector allocation.
kss compiled, that could be a single loop with no additional storage:
```go
/kss
func fI(xI, xn, yI, yn int32) (int32, int32) {
 for i:=int32(0); i<yn<<2; i+=4 {
  SetI32(xI+i, (i>>2)+(I32(xI+i)*I32(yI+i)))
 }
 free(xI-8, bk(4*xn))
 return yI, yn
}
```

kss has to identify runs over the same length and can integrate the scan into the loop:
```
fI:{[xi;yF]xi#+\yF*yF}
```
```go
/kss
func fI(xi, yF, yn int32) (int32, int32) {
 f1 := 0.0
 for i:=int32(0);i<yn<<3;i+=8 {
  f1 += F64(yF+i)*F64(yF+8)
  SetF64(yF+i, f1)
 }
 return itakeF(xi, yF, yn)
}
```
`itakeF` would be a kss intrinsic implementation for `xi#xF`, e.g:
```go
func itakeF(n, xF, xn int32) (int32, int32) {
 if n<0 { return idropF(n+xn, xF); }
 if bk(8*n) == bk(8*xn) {
  return xF, n
 } else {
  var r int32 = 8+alloc(8*n)
  Memorycopy(r, xF, 8*n)
  free(xF-8)
  return r, n
 }
}
```
this algorithm for standard deviation would be compiled to two loops with no additional storage:
```
std:{[xF]%(+/xF*xF:(xF-(+/xF)%#xF))%-1+#xF}
```
```go
func std(xF, xn int32) (float64) {
 var t1 float64
 for i:=int32(0); i<xn<<3; i+=8 {
  t1 += F64(xF+i)
 }
 t1 /= xn
 var t2
 for i:=int32(0); i<xn<<3; i+=8 {
  var t3 = (F64(xF+i)-t1) / float64(-1+xn))
  t2 += t3*t3
 }
 free(xF, bt(8*xn))
 return sqrt(t2)
}
```

binary search is an example for a scalar algorithm that could be expressed with k code and compiled to an efficient native function:
```
bin:{[xI;yI]
 i:0;while[i<#yI
  k:0;j:-1+#xI
  while[~k<j;h:(k+j)%2;$[xI[h]>yi;j:h-1;k:h+1]
  yI[i]:k-1;i+:1];yI}
```
```go
func bin(xI, xn, yI, yn int32) (int32, int32) {
 var i int32 = 0
 for i<yn {
  var k int32 = 0
  var j int32 = -1 + xn
  for k > j {
   var h int32 = (k + j ) >> 1
   if I32(xI+4*h) > yi {
    j = h - 1
   } else {
    k = h + 1
  }
  yI[i>>2] = k - 1
  i++
 }
 free(xI-8, bk(4*xn)
 return yI, yn
}
```
