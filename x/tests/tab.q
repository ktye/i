function  $f(){
@start
 
}
function w $g(w %x){
@start
 ret %x
}
function w $h(w %x, w %y){
@start
 %.1 =w add %x, %y
 ret %.1
}
function init(){
@start



%_F =l call $malloc(24);
%fi =l add %_F, 0
storel $f, %fi
%fi =l add %_F, 8
storel $g, %fi
%fi =l add %_F, 16
storel $h, %fi

}
