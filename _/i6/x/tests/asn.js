const f=function(x){
 let r
 r=I(I(1)+x);
 return r;
}
const h=function(x){
 let r
 f(x);
 r=f(x);
 return r;
}
