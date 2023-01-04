function w $g(w %x, w %y){
@start
 %.1 =w sub %x, %y
 ret %.1
}
function w $f(w %x){
@start
 %.1 =w call $g(w 1, w %x)
 ret %.1
}
