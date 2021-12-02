#include<stdio.h>

#define KTYE
#include"k.h"

// external libraries
void loadlib(const char*);

// ktye internal api
int args_;
char **argv_;
void **_F;
void cinit();
void kinit();
void doargs();
void store();
void try();
K read();
int32_t idx(int32_t, int32_t, int32_t);

// c-api implementation
K K2(char v, K x, K y) { 
 int32_t p=1+idx(v,288,253);
 return ((K(*)(K,K))_F[64+p])(x,y); 
}
void KR(const char *name, void *fp, int arity) {
 //todo
}
//todo..

int main(int args, char **argv){
 args_=(int32_t)args;
 argv_=argv;
 cinit();
 kinit();

 loadlib("lib");

 doargs();
 printf("ktye/k+\n");
 store();
 while(1){
  printf(" ");
  K x = read( );
  try(x);
 }
 return 0;
}

