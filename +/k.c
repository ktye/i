#include<stdio.h>
#include"k.h"
#define NATIVE
int64_t native(int64_t x,int64_t y);
#include"ktye.h"  // ktye/k implementation

// external libraries
#ifdef MAT
void loadmat();
#else
void loadmat(){printf("not loading mat\n");}
#endif

#ifdef DRW
void loaddrw();
#else
void loaddrw(){}
#endif

#ifdef RAY
void loadray();
#else
void loadray(){}
#endif

#ifdef SQL
void loadsql();
#else
void loadsql(){}
#endif



int main(int args, char **argv){
 args_=(int32_t)args;
 argv_=argv;
 kinit();

 loadmat();
 loaddrw();
 loadray();
 loadsql();

 ktye_doargs();
 printf("ktye/k+\n");
 ktye_store();
 while(1){
  printf(" ");fflush(stdout);
  K x = ktye_readfile(ktye_mk(Ct,0));
  ktye_try(x);
 }
 return 0;
}
