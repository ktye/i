int32_t f(int32_t x){
 return x;
}
int32_t g(int32_t x){
 int32_t r;
 switch(x){
  case 0:
  r=f(x);
  break;
  default:
  {
   if(x>5)r=(x-3);
   else r=(x-2);
   r=r;
  }
  break;
 };
 return r;
}
