function w $f(w %x){
@start
 %.1 =w sgt %x, 0
 jnz %.1, @i1, @e1
 @i1
 ret 1
 jmp @f1
 @e1
 ret 2
 @f1
}
