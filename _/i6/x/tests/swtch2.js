const f=function(x){
 return x;
}
const g=function(x){
 let r
 switch(x){
  case 0:r=f(x);
  break;
  default:{
   r=(x>I(5))?(I(x-I(3))):(I(x-I(2)));
  }
  break;
 } 
 return r;
}
