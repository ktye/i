/*
builtin         intrinsic _mm512_..
-----------------------------------
ldmxcsr         _mm_setcsr(x)
pdep_di         _pdep_u64(x,y)
pext_di         _pext_u64(x,y)
sqrtps          sqrt_round_ps(x,y)
permvarqi       permutexvar_epi8(y,x)
permvarsi       permutexvar_epi32(y,x)
vpermi2varps    permutex2var_ps(x,y,z)
vpermi2varqi    permutex2var_epi8(x,y,z)
selectb         mask_mov_epi8(z,x,y)
cvtb2mask       movepi8_mask(x)
compressqi      mask_compress_epi8(y,z,x)
*/

#define At(e...) __attribute((e))
#define V4 At(vector_size(1<<4),aligned(1))
#define V6 At(vector_size(1<<6),aligned(1))
typedef unsigned char Z V6;
typedef unsigned i2;
typedef unsigned long u;
typedef float e2;

#define bu(f) __builtin_##f
#define bo(f) bu(ia32_##f)
#define zo(f) bu(ia32_##f##512)
#define zm(f) bu(ia32_##f##512_mask)

#ifdef XXX
#undef bo
#undef zo
#undef zm
#define bo(f) _ia32_##f
#define zo(f) _ia32_##f##512
#define zm(f) _ia32_##f##512_mask

void _ia32_ldmxcsr(i2 i){}
u _ia32_pdep_di(u x, u y){u r=0; for(u i=1; y; i+=i){if(x&i)r|=y&-y;y&=y-1;};return r;}
u _ia32_pext_di(u x, u y){u r=0,m=0;for(u i=0;i<64;i++){if(1&(y>>i))r|=(1l&(x>>i))<<m++;};return r;}

Z _ia32_sqrtps512(     e2 V6 x, u y     ){e2 V6 r;                 for(u i=0;i<16;i++)             r[i]=__builtin_sqrtf(x[i]);    return r;}
Z _ia32_permvarqi512(      Z x, Z y     ){Z r;                     for(u i=0;i<64;i++)             r[i]=x[y[i]];                  return r;}
Z _ia32_permvarsi512(      Z x, Z y     ){i2 V6 a=x,b=y,r;         for(u i=0;i<16;i++)             r[i]=b[a[i]];                  return r;}
Z _ia32_vpermi2varps512(   Z x, Z y, Z z){i2 V6 a=x,i=y,b=z,r;i2 m;for(u j=0;j<16;j++){m=i[j]&0x0f;r[j]=((0xf0&i[j]))?b[m]:a[m];};return r;}
Z _ia32_vpermi2varqi512(   Z x, Z y, Z z){Z r;i2 m;                for(u i=0;i<64;i++){m=63&y[i];  r[i]=(64&y[i])?z[m]:x[m];};    return r;}
Z _ia32_selectb_512(       u k, Z a, Z w){Z r;                     for(u i=0;i<64;i++)             r[i]=(1&(k>>i))?a[i]:w[i];     return r;}
u _ia32_cvtb2mask512(      Z x          ){u r=0;                   for(u i=0;i<64;i++)             r  |=(x[i]>>7)?(1l<<i):0;      return r;}
Z _ia32_compressqi512_mask(Z a, Z s, u k){Z r=s;i2 m=0;            for(u i=0;i<64;i++){if(1&k>>i)r[m++]=a[i];                   };return r;}

#endif

#include<stdio.h>

void U(u x){printf("%ld\n",x);}
void P(Z x){for(u i=0;i<64;i++)printf("%02x ",x[i]);printf("\n");}
void Q(i2 V6 x){for(u i=0;i<16;i++)printf("%08x ",x[i]);printf("\n");}
void F(e2 V6 x){for(u i=0;i<16;i++)printf("%6.3f ",x[i]);printf("\n");}
Z iota(){Z r;for(u i=0;i<64;i++)r[i]=i;return r;}
Z jota(){i2 V6 r;for(u i=0;i<16;i++)r[i]=i;return r;}
Z reva(){Z r;for(u i=0;i<64;i++)r[i]=63-i;return r;}
Z rewa(){i2 V6 r;for(u i=0;i<16;i++)r[i]=15-i;return r;}
Z alta(){i2 V6 r;for(u i=0;i<16;i++)r[i]=(i%2)?0:-1;return r;}
u pdep_di(){
 u x=123456789234,y=2302402;
 return bo(pdep_di)(x,y);}
u pext_di(){
 u x=123456789234,y=2302402;
 return bo(pext_di)(x,y);}
Z sqrtps(){
 e2 V6 x;for(u i=0;i<16;i++)x[i]=(e2)i;
 return zo(sqrtps)(x,4);}
Z permvarqi(){
 //Z x=iota(),y=iota();
 Z x={49,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17
  ,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33
  ,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,10,51
  ,52,53,54,55,56,57,58,59,60,61,62,63};
 Z y={255,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19
  ,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39
  ,40,41,42,43,44,45,46,47,8,49,50,51,52,53,54,55,56,57,58,59,60,61,62};
 return zo(permvarqi)(x,y);}
Z permvarsi(){
 i2 V6 x=jota(),y=jota();
 y[1]=4;y[2]=3;
 return zo(permvarsi)(x,y);}

Z vpermi2varps(){
 e2 V6 x=jota(),z=rewa();
 i2 V6 y=jota();
 y[0]=5;y[1]=2;y[2]=9;y[3]=8;y[4]=1;
 y[2]|=0xf0;
 y[4]|=0x10;
 return zo(vpermi2varps)(x,y,z);}
Z vpermi2varqi(){
 Z x=iota(),y=iota(),z=reva();
 y[0]=5;y[1]=2;y[3]=8;y[4]=1;
 y[2]|=64;
 y[4]|=128;
 return zo(vpermi2varqi)(x,y,z);}
Z selectb_(){
 u k=123456789;
 Z x=iota(),y=reva();
 return zo(selectb_)(k,x,y);}
u cvtb2mask(){
 Z x={0,0,255,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,255,
    0,0,255,0,0,0,0,0,0,255,0,0,255,0,0,0,0,255,
    0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
 u r = (u)zo(cvtb2mask)(x);
 return r;}
Z compressqi(){
 u k=0xff0fff0000000000;
 Z x=iota(),y=reva();
 return zm(compressqi)(x,y,k);
}

int main(){
 U(pdep_di());
 U(pext_di());
 F(sqrtps());
 P(permvarqi());
 P(permvarsi());
 Q(vpermi2varps());
 P(vpermi2varqi());
 P(selectb_());
 U(cvtb2mask());
 P(reva());
 P(compressqi());
}
