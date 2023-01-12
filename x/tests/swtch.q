function w $f(w %x){
@start
 %.2 =w eq %x, 0
 jnz %.2, @i1, @e1
 @i1
 %.1 =w add 1, %x
 ret %.1
 jmp @f1
 @e1
 ret %x
 @f1
}
function w $g(w %x){
@start
 %.3 =w eq %x, 0
 jnz %.3, @i2, @e2
 @i2
 %.1 =w add 1, %x
 ret %.1
 jmp @f2
 @e2
 %.2 =w eq %x, 1
 jnz %.2, @i1, @e1
 @i1
 ret %x
 jmp @f1
 @e1
 @f1
 @f2
 ret 0
}
