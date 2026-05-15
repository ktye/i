function w $f(w %x){
@start
 %.2 =w sgt %x, 3
 jnz %.2, @i1, @e1
 @i1
 ret %x
 jmp @f1
 @e1
 %x =w mul 2, %x
 %.1 =w sub 0, %x
 ret %.1
 @f1
}
