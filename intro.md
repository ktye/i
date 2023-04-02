a brief introduction about the pecularities of ktye/k  [[qr](#qr-decomposition-least-squares);[fft](#fft);[stats](#statistics)]

# get/compile
get any version from the [table](https://github.com/ktye/i#build) e.g.
```
$wget https://github.com/ktye/i/releases/download/latest/k.c
$gcc -ok k.c -lm
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

adverbs have verb overloads if the left arg is not a verb: `' / \ ': /: \:` are `bin mod div in join split`.
- split and join allow vectors (e.g. words) on the left and are not restricted to chars
- div: `x\y` is the same as `y%x`. division always remains within their domain, e.g. ints. to get floats uptype must be explicit: `x%2.`
- encode decode need a double slash: `x//y` and `x\\y` (i always implement them late when everything is filled)


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
  x[I;K]-:{+/x*y}/:[(conj x[j;K]);x[I;K]]*\:x[j;K]
  K:1_K;j+:1];(x;r;n;m)}

solve:{qslv:{H:x 0;r:x 1;n:x 2;m:x 3;j:0;K:!m
 while[j<n;y[K]-:(+/(conj H[j;K])*y K)*H[j;K];K:1_K;j+:1]
 i:n-1;J:!n;y[i]%:r@i
 while[i;j:i_J;i-:1;y[i]:(y[i]-+/H[j;i]*y@j)%r@i]
 n#y}
 q:$[`i~@*|x;x;qr x];$[`L~@y;qslv/:[q;y];qslv[q;y]]}

dot:(+/*)\:

A:(3 5 -8 12.;-2 3 3 0.;7 -8 2 1.)
b:dot[+A;x:1 2 3];0.00001>|/abs x-solve[A;b]

A:+0a0+(1 -2a90 3;5a90 3 2;2 3 1;4 -1 1);
0.0001>|/abs r-solve[A;dot[+A;r:1a30 2a30 3a30]]
```
qr works for both, real and complex input. A is stored in column major order (list of columns).

## fft
```
fft: {[f;x]{[x;p;k;e]@[x;k,j;(x[k]+x[j]*e),x[k]-x[j:k+p]*e]}/[x f 0;f 1;f 2;f 3]}
fftn:{[n]l:*&|2\\n;e:angle[1;-(i:!n)*360.%n]
 (2//|2\\i;2^/:!l;|&'~!l#2;e@|(h/2*)\:!h:n%2)}

n:2^15
f:fftn n
fft[f;`z@?2*n]
```
see [details](https://github.com/ktye/i/tree/master/%2B/mat/fft)

## statistics
```
avg:{(+/x)%0.+#x}
var:{(+/x*x:(x-avg x))%-1+#x}
std:{%var x}
med: **|2^^
```
