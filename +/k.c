#include<stdio.h>
#include"ktye.h"  // ktye/k implementation

// external libraries
void loadlib(const char*);


// c-api implementation
#define K uint64_t
void kinit(){ ktye_cinit();ktye_kinit(); }
K K2(char v, K x, K y) { 
 int32_t p=1+ktye_idx(v,288,253);
 return ((K(*)(K,K))_F[64+p])(x,y); 
}
void KR(const char *name, void *fp, int arity) {
 //todo
}
//todo..

int main(int args, char **argv){
 args_=(int32_t)args;
 argv_=argv;
 kinit();

 loadlib("lib");

 ktye_doargs();
 printf("ktye/k+\n");
 ktye_store();
 while(1){
  printf(" ");
  K x = ktye_read( );
  ktye_try(x);
 }
 return 0;
}

