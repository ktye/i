let f=function(x){
 let r
 r=I((1|0)+x);
 return r;
}
let h=function(x){
 let r
 f(x);
 r=f(x);
 return r;
}
