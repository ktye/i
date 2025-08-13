/*
#include"fftw3.h"
#include<stdint.h>

#define N 32768


int32_t rand_=1592653589;
double rnd(){ //uniform double, same version as k.
 int32_t r = rand_;
 r^=(r<<13);r^=(r>>17);r^=(r<<5);rand_ = r;
 return 0.5+((double)r)/4294967295.0;
}


void P(double *x){printf("%lf %lf %lf %lf..\n",x[0],x[1],x[2],x[3]);}
void R(double *x){for(int i=0;i<2*N;i++)x[i]=rnd();}


int main(int args, char **argv){
 fftw_complex *x;
 fftw_plan p;
 

 x=(fftw_complex*)fftw_malloc(sizeof(fftw_complex)*N);
 p=fftw_plan_dft_1d(N,x,x,FFTW_FORWARD,FFTW_ESTIMATE);

 for(int i=0;i<1000;i++){
  R((double*)x);
  fftw_execute(p);
 }
 P((double*)x);
 return 0;
}

*/
