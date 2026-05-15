int32_t f(int32_t x){
 int32_t r;
 r=(1+x);
 return r;
}
int32_t h(int32_t x){
 int32_t r;
 f(x);
 r=f(x);
 return r;
}
