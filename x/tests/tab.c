void f(){
}
int32_t g(int32_t x){
 return x;
}
int32_t h(int32_t x, int32_t y){
 return (x+y);
}
void init(){
 _F=malloc(3*sizeof(void*));
 _F[0]=f;
 _F[1]=g;
 _F[2]=h;
}
