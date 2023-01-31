#include<stdint.h>
#include"fenster.h"

uint64_t cnative(uint64_t x,uint64_t y);
#define NATIVE
#undef min
#undef max
#include"k.h"

int32_t W,H,SC;     //dimensions,scale(1|2)
uint32_t *O;        //pixels
uint32_t *kui=NULL; //event log
int nui=0,aui=0;

// todo conform xI yI
uint64_t cnative(uint64_t x,uint64_t y){ //P[color;pixelindex] /ci;pI
 x=Fst(rx(y));y=Las(y);
 dx(x);dx(y);if((tp(x)!=it)||(tp(y)!=It))return 0;
 int32_t  n=nn(y),yp=(int32_t)y>>2;
 if(SC==2){
  for(int32_t i=0;i<n;i++){
   int32_t j=((int32_t*)_M)[yp+i];
   if((j>=0)&&(j<W*H)){int jj= 2*((j%W) + 2*W*(j/W));
    O[jj      ]=(uint32_t)x;
    O[jj+1    ]=(uint32_t)x;
    O[jj+2*W  ]=(uint32_t)x;
    O[jj+2*W+1]=(uint32_t)x;
 }}}else{
  for(int32_t i=0;i<n;i++){
   int32_t j=((int32_t*)_M)[yp+i];
   if((j>=0)&&(j<W*H))O[j]=(uint32_t)x;
}}}


const char *p0="W:H:100;B:P 255;G:P 255*256;R:P 255*256*256;wh:!W*H;M:{`mouse \\(x)};K:{`key \\x;$[x~82;R wh;x~71;G wh;x~66;B wh;P[0;wh]];}";
const char *sl="`-=[];'\\,./";
const char *Sl="~_+{}:\"|<>?";

uint64_t ks(const char *p){
 int32_t  n=(int32_t)strlen(p);
 uint64_t r=mk(Ct,n);
 memcpy(_M+(int32_t)r,p,n);
 return r;
}

int shift(int k,int s){
 if((k>=65)&&(k<=90)&&!s)return k+32;
 if((k>=48)&&(k<=57)&&s) return ")!@#$%^&*("[k-48];
 if(s)for(int i=0;i<11;i++)if(sl[i]==k)return Sl[i];
 return k;
}

void logspc(void)   {if(nui==aui)kui       =(uint32_t*)realloc(kui,aui=4*nui*2);}
int logkey(int x)      {logspc();kui[nui++]=((uint32_t)x)&0x0000ffff;return x;}
int logmas(int x,int y){logspc();kui[nui++]=((x&0xffff)<<16)|(y&0xffff)|0x80000000;}
int writekui(){
 FILE *f=fopen("ku.in.k", "w");
 fprintf(f,"P:{y}\n");
 for(int i=0;i<nui;i++){uint32_t x=kui[i];
  if(x&0x80000000)fprintf(f,"M[%u,%u]\n",(x&0x7fff0000)>>16,x&0xffff);
  else            fprintf(f,"K %u\n",x);
 }
 fclose(f);
 printf("w ku.in.k\n");
 return 1;
}

int fensterl(){ //verb: fens·terln
 uint64_t M=sc(Ku('M'));
 uint64_t K=sc(Ku('K'));
 O=calloc(SC*SC*W*H,4);
 struct fenster f={
  .title="ku",
  .width=SC*W,
  .height=SC*H,
  .buf=O,
 };
 fenster_open(&f);
 int64_t now=fenster_time();
 while(!fenster_loop(&f)){
  for(int i=0;i<128;i++)if(f.keys[i]){dx(Atx(Val(K),Ki(logkey(shift(i,f.mod&2)))));f.keys[i]=0;}
  if(f.keys[27])break;
  if(f.mouse){f.x/=SC;f.y/=SC;logmas(f.x,f.y);dx(Atx(Val(M),Cat(Ki((int32_t)f.x),Ki((int32_t)f.y))));f.mouse=0;}
  int64_t time=fenster_time();
  if(time-now<1000/60)fenster_sleep(time-now);
  now=time;
 }
 fenster_close(&f);
 return 0;
}
int main(int args, char **argv){
 if((args>1)&&(argv[1][0]=='2')){args--;argv++;SC=2;}else SC=1;
 args_=args;
 argv_=argv;
 init();
 kinit();
 Asn(sc(Ku('P')),((l2(0,Ki(1)))&0xffffffff)|((uint64_t)xf<<59)); //P:native func
 dx(val(ks(p0)));
 doargs();
 W=(int32_t)Val(sc(Ku('W')));
 H=(int32_t)Val(sc(Ku('H')));
 printf("W=%d H=%d %d\n", W, H, SC*W*H);
 return(!setjmp(_jb_))?fensterl():writekui();
}

