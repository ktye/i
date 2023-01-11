function  $f(){
@start
 %i =w copy 0
 @L
 %i =w add %i, 1
 %.1 =w sgt %i, 3
 jnz %.1, @i1, @e1
 @i1
 jmp @Le
 @e1
 jmp @L
 @Le
}
function  $g(){
@start
 %i =w copy 0
 @f26
 %i =w add %i, 1
 %.1 =w sgt %i, 3
 jnz %.1, @i1, @e1
 @i1
 jmp @f26e
 @e1
 jmp @f26
 @f26e
}
