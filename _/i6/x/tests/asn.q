function w $f(w %x){
@start
 %r =w add 1, %x
 ret %r
}
function w $h(w %x){
@start
 %.1 =w call $f(w %x)
 %r =w call $f(w %x)
 ret %r
}
