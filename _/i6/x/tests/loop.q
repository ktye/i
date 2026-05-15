function  $f(){
@start
 %n =w copy 0
 %i =w copy 0
 @f11
 %.1 =w slt %i, 3
 jnz %.1, @f11c, @f11e
 @f11c
 %n =w add %n, 1
 %i =w add %i, 1
 jmp @f11
 @f11e
}
