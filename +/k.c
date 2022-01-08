#include<stdio.h>
#include"k.h"
#include"ktye.h"  // ktye/k implementation

// external libraries
void loadmat();
//void loadimg();
void loadray();


int main(int args, char **argv){
 args_=(int32_t)args;
 argv_=argv;
 kinit();

 loadmat();
 //loadimg();
 loadray();

 ktye_doargs();
 printf("ktye/k+\n");
 ktye_store();
 while(1){
  printf(" ");fflush(stdout);
  K x = ktye_read( );
  ktye_try(x);
 }
 return 0;
}
