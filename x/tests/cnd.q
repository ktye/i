function w $f(w %x){
@start
 %.3 =w sgt %x, 1
 jnz %.3, @i2, @e2
 @i2
 ret %x
 @e2
 %.2 =w sgt %x, 3
 jnz %.2, @i1, @e1
 @i1
 ret %x
 jmp @f1
 @e1
 %.1 =w sub 0, %x
 ret %.1
 @f1
}
