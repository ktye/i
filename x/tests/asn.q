function w $f(w %x){
@start
 %r =w add 1, %x
 ret %r
}
function w $h(w %x){
@start
 call $f(w %x)
 %r =w call $f(w %x)
 ret %r
}
