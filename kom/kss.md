# kss

```
fI:{[xI;yI](!#xI)+xI*yI}
```
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
```go
/kom rce
func f(x, y K) (r K) {
 return Add(Til(Cnt(rx(x)),Mul(x,y))
}
```
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
  var t3 = (F64(xF+i)-t1) / (-1+xn))
  t2 += t3*t3
 }
 free(xF, bt(8*xn))
 return sqrt(t2)
}
```
