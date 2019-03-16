# ⍳ interpret

## Usage
The package exports two functions *P* and *E* and no types.
```go
    var  a map[interface{}]interface{} // k tree
    l := P("!3")  // Parse expr to ast (list []interface{})
    v := E(l, a)  // Evaluate to value (interface{})
```

## Types
The interpreter uses `float64` and `complex128` for numbers, slices of them for uniform vectors,
`[]interface{}` for a general list, `map[interface{}]interface{}` for dicts, `string` for symbols,
`[]rune` for character vectors and `func` for functions.

But any user supplied types can be used as well:
- If a type is a slice, it is interperted as a list and all vector functions are supported.
- If a type is a map or a struct, it is used as a dictionary with string keys in the latter case.
- Any func can be used as a function.
- Types can overload all numerical operators, by implementing the interface:
  - Monadic example: `func (t myType) Neg() interface{}`
  - Dyadic example: `func (t myType) Add(b interface{}, recvIsLeft bool) interface{}`
  - The overloaded type may be a slice, that receives the vector directly, or it is called for every atom.
- Convertible numeric types (to float or complex) can be used for all mathematical operators. The results are converted back to the original type (including booleans).
- To use native Go values or functions, add them to the k-tree directly: `a["myfunc"] = func(...interface{}) interface{} { return 42 }`

## Verbs

```
    a          l          a-a        l-a        a-l        l-l     triad  tetrad
+   flp        flp        [add]      [add]      [add]      [add]    -      -        ⍉
-   [neg]      [neg]      [sub]      [sub]      [sub]      [sub]    -      -      
*   fst        fst        [mul]      [mul]      [mul]      [mul]    -      -        ×
%   [sqr]      [sqr]      [div]      [div]      [div]      [div]    -      -        √÷
!   til        odo        mod        -          mod>       mkd      -      -        ⍳
&   wer        wer        [min]      [min]      [min]      [min]    -      -        ⍸⌊
|   rev        rev        [max]      [max]      [max]      [max]    -      -        ⌽⌈
<   asc        asc        [les]      [les]      [les]      [les]    -      -        ⍋
>   dsc        dsc        [mor]      [mor]      [mor]      [mor]    -      -        ⍒
=   eye        grp        [eql]      [eql]      [eql]      [eql]    -      -        ⌸
~   [not]      [not]      mch        mch        mch        mch      -      -        ≡
,   enl        enl        cat        cat        cat        cat      -      -        ∊
^   is0        [is0]      ept        ept        ept        ept      -      -        ∧
#   cnt        cnt        tak        rsh        tak        rsh      -      -        ⍴≢↑ 
_   [flr]      [flr]      drp        fil        drp        cut      -      -        ↓ (not ⌊)
$   fmt        [fmt]      cst        cst        cst        cst      -      -        ⍕∪
?   rng        unq        rnd        fnd        pik        fnd>     spl    -        ⍷
@   typ        typ        atx        atx        atx        atx      amd    amd      ⌶⌷
.   evl        evl        cal        cal        cal        cal      dmd    dmd      ⍎
'   -          -          -          bin        -          bin>     -      -  
/   -          -          -          -          pak        pak      -      -        
\   -          -          -          upk        spl        -        -      -   
':  -          -          -          -          win        -        -      -  
```

## Adverbs

```
       mv/nv      dv         l-mv        l-dv      3+v
':     -          ecp        -           ecp       -       ⍨
'      ech        ecd        ecd         ecd       eca     ¨
/:     -          -          ecr         ecr       -       ⌿
\:     -          -          ecl         ecl       -       ⍀
/      fxd        ovr        fxw         ovd       ova   
\      scf        scn        scw         scd       sca  
```

## Others
```
:  ←
x  ⍺
y  ⍵
o  ∇
```

## Ref
https://github.com/JohnEarnest/ok serves as a reference. It comes with documentation.
