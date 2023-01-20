function w $f(w %x){
@start
 ret %x
}
function w $g(w %x){
@start
 %.2 =w eq %x, 0
 jnz %.2, @i2, @e2
 @i2
 %r =w call $f(w %x)
 jmp @f2
 @e2
 %.1 =w sgt %x, 5
 jnz %.1, @i1, @e1
 @i1
 %r =w sub %x, 3
 jmp @f1
 @e1
 %r =w sub %x, 2
 @f1
 %r =w copy %r
 @f2
 ret %r
}
