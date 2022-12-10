int32_t f(int32_t x){
 return x;
}
int32_t g(int32_t x){
 int32_t r;
 switch(x){
  case 0:r=f(x);
  break;
  default:{
   r=(x>5)?((x-3)):((x-2));
   r=r;
  }
  break;
 } 
 return r;
}
