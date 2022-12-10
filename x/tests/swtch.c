int32_t f(int32_t x){
 switch(x){
  case 0:return (1+x);
  break;
  default:return x;
  break;
 };
}
int32_t g(int32_t x){
 switch(x){
  case 0:return (1+x);
  break;
  case 1:return x;
  break;
 };
 return 0;
}
