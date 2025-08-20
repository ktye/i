# ktye/k intro ([qr](#qr-decomposition-least-squares);[lu](#lu-decomposition);[svd](#singular-value-decomposition);[fft](#fft);[stats](#statistics);[trans](#transpose-permute-axes))

```
ktye/k   ktye.github.io/k.html
+ flp add  '  ech pri both bin
- neg sub  /  ovr fix echright
* fst mul  \  scn fix eachleft
% sqr div      / join   decode
! til key  mod \ split  encode
& wer min  $[a;b;...]     cond
| rev max  while[c;a;b;d;e;..]
< asc les  f:{x+y}   [bl;o;ck]
> dsc mor            "chars" c
= grp eql  01234567   1 2 3  i
~ not mtc   :+-*%&|  .4 5 6. f
, enl cat  <>=~!,^#   2a300  z
^ srt cut  _$?@.     (1;2 3) L
# cnt tak           `a`b!5 6 D
_ flr drp  t,d t,t t,'t   join
$ str cst           k!t    key
? unq fnd  in       k?t  group
@ typ atx  @[x;i;+;y]    amend
. val cal  .[x;i;+;y]    dmend
                              
abs sin cos exp log find angle
imag conj  types:cisfzLDTvcdlx
?n(uniform) ?-n(normal) ?z(bi)
n?n(with)   random   -n?n(w/o)
```


# get/compile
get any version from here:

|file|what|compile|created with|
|---|---|---|
[k.f](https://github.com/ktye/i/releases/download/latest/k.f)|fortran version|`gfortran k.f`|[f.go](https://github.com/ktye/wg/blob/master/f77/f.go)|
[k.c](https://github.com/ktye/i/releases/download/latest/k.c)|c version|`gcc k.c -lm`|[cc.k](https://github.com/ktye/i/blob/master/x/cc.k)|
[kv.c](https://github.com/ktye/i/releases/download/latest/kv.c)|c simd5 |`clang-18 -O3 -mavx2 kv.c -lm`|[c.go](https://github.com/ktye/wg/blob/master/c.go)|
[k.go](https://github.com/ktye/i/releases/download/latest/k.go)|go version|`go build k.go`|[go.k](https://github.com/ktye/i/blob/master/x/go.k)|
[k.wasm](https://ktye.github.io/k.wasm)|run [online](https://ktye.github.io.k.html)||[wat.go](https://github.com/ktye/wg/blob/master/wat.go)|

e.g.
```
$wget https://github.com/ktye/i/releases/download/latest/k.c
$gcc -ok k.c -lm          #clang -ok -fwrapv -ffast-math k.c
$k
ktye/k
 █
```

# command line arguments
- `k a.k b.k c.k` executes all files in order, as if they were catenated and drops into the repl
- `k a.k -e '+/x'` loads a.k, evaluates the expression and exits. *e* stands for both: *eval* and *end* (like in awk)
- `k a.k -e` prevents the repl

executing a file parses and executes everything in one go. it is not line oriented as opposed to other apl/k.
If you want output from a line, print or debug:

# print, debug
debug with a backslash: `x+ \y`; it is also dyadic: to include a label ``x+`Y \y:3`` prints `` `Y:3``.
that means you sometimes have to include an @ if you want to force the monadic form.

in the repl output is converted to (clipped) 2d form by default.
if you prepand a space, it uses k syntax in one long line.

both forms are also available as `` `l x`` and `` `k x`` and return chars.
they are defined in [z.k](z.k) which is built in.

for file i/o there are no numeric verbs and there is only 1 form:
- read: `` <`file`` returns chars, e.g. ``x:<`file``
- write: `` `file<"chars only\n" ``
- to stdout: `` `<"..."``


# special forms, adverbs and overloads
the only keyword is `while`. there is block `[x;y;z]` similar to `*(z;y;x)` for use in cond `$[a;b;c]`.

there are 3 adverbs: `' / \` (each over scan). they have verb overloads if the left arg is not a function:
- `x'y` is bin (binary search) or lin (linear interpolation) depending on x being an atom or a vector.
- `x/y` and `x\y` are decode/encode or join/split depending on x being numeric or characteristic.
- `x+/y` and `x+\y` is each-right/each-left.
there is no reduce/scan with initial values, but higher order forms exist: `f3/[x;y;z]`.
- `-'x` is each-prior/pairs. it is used only for `:+-*%&|<=>` and only for unnested vector arguments. otherwise each is called.
- `f/x` and `f\x` is fixpoint if f is strictly monadic.

# dots, names
a dot is not part of a symbol: `a . b` can be written as `a.b`. to index a list/dict use `` d`a``.

there are no namespaces, only flat global variables and locals.

there are also no undefined variables/errors. when k creates a new symbol, it also creates an associated spot for a global variable which is zero (null verb).

# examples

## qr decomposition, least squares
```
qr:{K:!m:#*x;I:!n:#x;j:0;r:n#0a;turn:$[`Z~@*x;{(-x)@angle y};{x*1. -1@y>0}]
 while[j<n;I:1_I
  r[j]:turn[s:abs@abs/j_x j;xx:x[j;j]]
  x[j;j]-:r[j]
  x[j;K]%:%s*(s+abs xx)
  x[I;K]-:{+/x*y}/[(conj x[j;K]);x[I;K]]*\x[j;K]
  K:1_K;j+:1];(x;r;n;m)}

solve:{qslv:{H:x 0;r:x 1;n:x 2;m:x 3;j:0;K:!m
 while[j<n;y[K]-:(+/(conj H[j;K])*y K)*H[j;K];K:1_K;j+:1]
 i:n-1;J:!n;y[i]%:r@i
 while[i;j:i_J;i-:1;y[i]:(y[i]-+/H[j;i]*y@j)%r@i]
 n#y}
 q:$[`i~@*|x;x;qr x];$[`L~@y;qslv/[q;y];qslv[q;y]]}

dot:(+/*)\

A:(3 5 -8 12.;-2 3 3 0.;7 -8 2 1.)
b:dot[+A;x:1 2 3];0.00001>|/abs x-solve[A;b]

A:+0a0+(1 -2a90 3;5a90 3 2;2 3 1;4 -1 1);
0.0001>|/abs r-solve[A;dot[+A;r:1a30 2a30 3a30]]
```
qr works for both, real and complex input. A is stored in column major order (list of columns).

## lu decomposition
with partial pivoting
```
lu:{[A]i:0;k:!#A;P:!#A
 while[1<#k
  j:i+*&a=m:|/a:abs A[k;i]
  P[(i;j)]:P[(j;i)]
  A[(i;j)]:A[(j;i)]
  A[k:1_k;i]%:A[i;i]
  A[k;k]-:A[k;i]*\A[i;k]
  i+:1]
 (A;P)}

lusolve:{[LUP;b];A:*LUP;P:*|LUP
 r:{[x;i;a]x[i]-:+/(a k)*x k:!i}/[b P;!n:#A;A]
   {[x;i;a]x[i]:(x[i]-+/(a k)*x k:(1+i)_!#x)%a[i]}/[r;|!n;|A]}

A:5^?25
x:?5
b:(+/*)\[A;x]
x-lusolve[lu A;b]
```

## singular value decomposition

svd computes `A=USV'` using one-sided jacobi iterations.
both U and V are unitary and S contains the singular values on the diagonal.
```
svd:{[A];n:#A       /jacobi-hestenes  A:list of columns  real or complex
 d:{+/(conj x)*y}
 J:{(c;s:(z%a)*t*c:1%%1+t*t:((q>0)-q<0)%(abs q)+%1+q*q:(_0a+d[y;y]-d[x;x])%2*a:abs z)}
 R:{[c;s;x;y]((x*c)-y*conj s;(x*s)+y*c)}
 V:(!n)=/!n:#A
 P:{[A;ik]$[1e-14<abs z:d.A ik;[r:J.(A ik),z;V[ik]:R.r,V ik;A[ik]:R.r,A ik];A]}
 I:,/i,''&'i>/i:!n
 F:{[A]P/(,A),I}
 A:F/A
 U:A%s:abs/'A
 (U g;s g;V g:>s)}
```

P diagonolizes two columns of A using jacobi rotations computed by J and applied with R.
I contains all pairs of columns which are diagnolized in order.
since old columns may be modified by calls to P,
the procedure over all pairs must be repeated until convergence,
done by the fixed point iteration F/A.

if A is thin (more rows than columns), the svd can be done using R from the the qr decomposition of A.
V and the singular values are the same, but U must be premultiplied with Q.
the multiplication is done using the householder transformation stored in the qr decomposition.
this is similar to the premultiplication with QT within qrsolve, but applies the transformation in reverse order.
```
svq:{q:qr x;h:q 0;r:q 1;n:q 2;m:q 3;o:0|(1+&n)@!m
 qmul:{x:o*m#x;K:,m-1;j:n-1;while[-1<j;x[K]-:(+/(conj h[j;K])*x K)*h[j;K];K:(-1+*K),K;j-:1];x} /Q*
 d:{x*(!#x)=/!#x}
 (,qmul'U:D 0),1_D:svd R:(d r)+(n#'h)*i</i:!n}
```

```
eye:{(!x)=/!x}
dia:{x*eye@#x}
uni:{|/|/abs mul[+conj x;x]-eye@#x}              /test if unitary: maxabs I-x'*x
mul:{(+/*)/[x;y]}                                /r,x,y: list of rows
dvs:{[U;s;V]|/|/abs A-mul[U]mul[dia s;+conj V]}  /test if A is U*S*V'

A:3^?12  /or 3^?12a (complex)
q:svd A
`u \uni U:q 0
`v \uni V:q 2
`t \dvs.q
```

## fft
```
fft: {[f;x]{[x;p;k;e]@[x;k,j;(x[k]+x[j]*e),x[k]-x[j:k+p]*e]}/[x f 0;f 1;f 2;f 3]}
fftn:{[n]l:*&|2\n;e:angle[1;-(i:!n)*360.%n]
 (2/|2\i;2^/!l;|&'~!l#2;e@|{h!2*x}\!h:n%2)}

n:2^15
f:fftn n
r:fft[f;`z@?2*n]
```
see [details](https://github.com/ktye/i/tree/master/%2B/mat/fft)

or recursive (odd/even):
```
fft:{$[-1+n:#x;(x+r),(x:fft x o)-r:(1@(!n)*-180.%n)*fft[x 1+o:2*!n%:2];x]}  /radix-2 recursive odd/even split
ifft:{(conj fft conj x)%#x}                                                 /inverse is fft under conjugation
rfft:{.5*(x+y;1a270*x-y:conj(x:fft imag[x;y])n!n-!n:#x)}                    /real241 two for the price of one
fft2:+fft'+fft'                                                             /two-dim separate and do it again
```

## interpolation
```
lin:{$[`L~@z;lin[x;y]'z;[dx:0.+1_-'x;dy:0.+1_-'y;b:(-2+#x)&0|x'z;(y b)+(dy b)*(z-x b)%dx b]]}
lin[0 4.;3 4.;0.+!5] /3. 3.25 3.5 3.75 4.
```

## statistics
```
avg:{(+/x)%0.+#x}
var:{(+/x*x:(x-avg x))%-1+#x}
std:{%var x}
med: **|2^^
z95:{1.97*(((std@_x)^b)+(std imag x)^b)^1%b:3.2}
mavg:{(y-(-#y)#(-x)_y:+\y)%x}
mmax:{while[x-:1;y:{x|(*x),-1_x}y];y}
mmin:{while[x-:1;y:{x&(*x),-1_x}y];y}
```

## transpose (permute axes)
```
trans:{z@<(y x)/(!y)x} /y:input shape, output shape:y x, apl uses <x
trans[1 2 0;4 3 2;!24] /0 6 12 18 1 7 13 19 2 8 14 20 3 9 15 21 4 10 16 22 5 11 17 23
                       /in this case same as: ,/+^[*y;z]
```

