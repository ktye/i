function w $f(w %x, w %y){
@start
 %.1 =l mul 8, %x
 %.2 =l add %_F, %.1
 %.3 =l loadl %.2
 %.4 =w call %.3(w %y)
 ret %.4
}
