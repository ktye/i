randi:{y/ 'rx}                             /x ints [0 y)
randf:{( 'rx)%4.294967295e9}               /x floats [-0.5 0.5)
randn:{_randz x}                           /x floats normal distributed
randz:{(%-2.* 'l0.5+randf x)@360.*randf x} /x binormal distributed complex 
shuffle:{x@< 'r#x}                         /x vector
