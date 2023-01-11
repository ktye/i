function  $f(){
@start
 %i =w copy 0
 @f7
 %.1 =w slt %i, 3
 jnz %.1, @f7c, @f7e
 @f7c
 %i =w mul %i, 2
 %i =w add %i, 1
 jmp @f7
 @f7e
}
