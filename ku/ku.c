#include<stdint.h>
#include"fenster.h"

uint64_t cnative(uint64_t x,uint64_t y);
#define NATIVE
#include"k.h"

int32_t W,H;
uint32_t *O;
uint64_t cnative(uint64_t x,uint64_t y){ //P[color;indexlist]
 x=Fst(rx(y));y=Las(y);
 dx(x);dx(y);if((tp(x)!=it)||(tp(y)!=It))return 0;
 int32_t  n=nn(y),yp=(int32_t)y>>2;
 for(int32_t i=0;i<n;i++){
  int32_t j=((int32_t*)_M)[yp+i];
  if((j>=0)&&(j<W*H))O[j]=(uint32_t)x;
 }
}


const char *p0="W:H:100;B:P 255;G:P 255*256;R:P 255*256*256;wh:!W*H;K:{`key \\x;$[x~82;R wh;x~71;G wh;x~66;B wh;P[0;wh]];}";

uint64_t ks(const char *p){
 int32_t  n=(int32_t)strlen(p);
 uint64_t r=mk(Ct,n);
 memcpy(_M+(int32_t)r,p,n);
 return r;
}

int main(int args, char **argv){
 args_=args;
 argv_=argv;
 init();
 kinit();
 Asn(sc(Ku('P')),((l2(0,Ki(1)))&0xffffffff)|((uint64_t)xf<<59)); //P:native func
 dx(val(ks(p0)));
 
 doargs();
 W=(int32_t)Val(sc(Ku('W')));
 H=(int32_t)Val(sc(Ku('H')));
 printf("W=%d H=%d\n", W, H);
 uint64_t K=sc(Ku('K'));
 O=calloc(W*H,4);
 struct fenster f={
  .title="ku",
  .width=W,
  .height=H,
  .buf=O,
 };
 fenster_open(&f);
 int64_t now=fenster_time();
 while(!fenster_loop(&f)){
  for(int i=0;i<128;i++)if(f.keys[i])dx(Atx(Val(K),Ki(i)));
  if(f.keys[27])break;
  int64_t time=fenster_time();
  if(time-now<1000/60)fenster_sleep(time-now);
  now=time;
 }
 fenster_close(&f);
}
