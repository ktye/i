# ⍳ interpret

## Usage
The package exports only two functions and no types. The functions are *P* and *E*.
```go
    // Make a new interpreter with an empty k-tree
    a := make(map[interface{}]interface{}) 
    
    // Parse: func P(s string) interface{}
    l := P("k-expression")  // The return value is an ast
    
    // Eval: func E(a map[interface{}]interface{}, l interface{}) interface{}
    v := E(a, l)           // The return value is a list of values

```

## Types
The package does not define any special types. Most Go types can be used directly, even user defined ones.
- If a type is a slice, it is interperted as a list and all vector functions are supported. The default list, is `[]interface{}`.
- If a type is a map, it is used as a dictionary. The default dict is: `map[interface{}]interface{}`
- A struct is a dictionary as well. It's keys are strings.
- If it is a function, it can be used as one.
- Numbers are converted to float64 or complex128. The conversion is only done if math functions are called. Structural functions preserve the types.
- Numerical methods (Neg, Add, Div, ...) can be overloaded for custom types and used directly, by defining a method:
  - `type myInt string;  func (s myInt) Neg() interface{} { return "-" + s }`
- To create a value with a custom type, call a custom function that returns one.
- To create a custom function, add it to the k-tree: `a["myfunc"] = func(...interface{}) interface{} { return 42 }`

## Verbs

```
    a          l          a-a        l-a        a-l        l-l     triad  tetrad
+   idn        flp        [add]      [add]      [add]      [add]    -      -        ⍉
-   [neg]      [neg]      [sub]      [sub]      [sub]      [sub]    -      -      
*   fst        fst        [mul]      [mul]      [mul]      [mul]    -      -        ×
%   sqr        [sqr]      [div]      [div]      [div]      [div]    -      -        √÷
!   iot        odo        mod        -          mod>       mkd      -      -        ⍳
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
