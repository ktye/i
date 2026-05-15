function  $f(){
@start
 %i =w copy 0
 @f9
 %i =w add %i, 1
 %.1 =w slt %i, 2
 jnz %.1, @i1, @e1
 @i1
 jmp @f9 
 @e1
 %i =w mul %i, 2
 jmp @f9
 @f9e
}
