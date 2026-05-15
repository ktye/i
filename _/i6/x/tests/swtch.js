const f=function(x){
 switch(x){
  case 0:return I(I(1)+x);
  break;
  default:return x;
  break;
 } 
}
const g=function(x){
 switch(x){
  case 0:return I(I(1)+x);
  break;
  case 1:return x;
  break;
 } 
 return I(0);
}
