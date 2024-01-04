/*
builtin         intrinsic _mm512_..
-----------------------------------
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

#define bu(f) __builtin_##f
#define bo(f) bu(ia32_##f)
#define zo(f) bu(ia32_##f##512)
#define zm(f) bu(ia32_##f##512_mask)

#ifdef XXX
#undef zo
#undef zm
#define zo(f) _ia32_##f##512
#define zm(f) _ia32_##f##512_mask

Z _ia32_sqrtps512(float V6 x, int y){float V6 r;for(int i=0;i<16;i++)r[i]=__builtin_sqrtf(x[i]);return r;}
Z _ia32_permvarqi512(Z x, Z y){Z r;for(int i=0;i<64;i++)r[i]=y[x[i]];return r;}
Z _ia32_permvarsi512(Z x, Z y){unsigned int V6 a=x,b=y,r; for(int i=0;i<16;i++)r[i]=b[a[i]]; return r;}
Z _ia32_vpermi2varps512(Z x, Z y, Z z){unsigned int V6 a=x,i=y,b=z,r;int m;for(int j=0;j<16;j++){m=i[j]&0x0f;r[j]=((0xf0&i[j]))?b[m]:a[m];};return r; }
Z _ia32_vpermi2varqi512(Z x, Z y, Z z){Z r;int m;for(int i=0;i<64;i++){m=63&y[i];r[i]=(64&y[i])?z[m]:x[m];};return r;}
Z _ia32_selectb_512(unsigned long k, Z a, Z w){Z r;for(int i=0;i<64;i++)r[i]=(1&(k>>i))?a[i]:w[i]; return r;}
unsigned long _ia32_cvtb2mask512(Z x){unsigned long r=0;for(int i=0;i<64;i++)r|=(x[i]>>7)?(1<<i):0; return r;}
Z _ia32_compressqi512_mask(Z a, Z s, unsigned long k){Z r=s;int m=0;for(int i=0;i<64;i++){if(1&k>>i){r[m++]=a[i];}};return r;}

#endif

#include<stdio.h>

void P(Z x){for(int i=0;i<64;i++)printf("%02x ",x[i]);printf("\n");}
void Q(unsigned int V6 x){for(int i=0;i<16;i++)printf("%08x ",x[i]);printf("\n");}
void F(float V6 x){for(int i=0;i<16;i++)printf("%6.3f ",x[i]);printf("\n");}
Z iota(){Z r;for(int i=0;i<64;i++)r[i]=i;return r;}
Z jota(){unsigned int V6 r;for(int i=0;i<16;i++)r[i]=i;return r;}
Z reva(){Z r;for(int i=0;i<64;i++)r[i]=63-i;return r;}
Z rewa(){unsigned int V6 r;for(int i=0;i<16;i++)r[i]=15-i;return r;}
Z alta(){unsigned int V6 r;for(int i=0;i<16;i++)r[i]=(i%2)?0:-1;return r;}
Z sqrtps(){
 float V6 x;for(int i=0;i<16;i++)x[i]=(float)i;
 return zo(sqrtps)(x,4);}
Z permvarqi(){
 Z x=iota(),y=iota();
 y[1]=4;y[2]=3;
 return zo(permvarqi)(x,y);}
Z permvarsi(){
 unsigned int V6 x=jota(),y=jota();
 y[1]=4;y[2]=3;
 return zo(permvarsi)(x,y);}

Z vpermi2varps(){
 float V6 x=jota(),z=rewa();
 unsigned int V6 y=jota();
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
 unsigned long k=123456789;
 Z x=iota(),y=reva();
 return zo(selectb_)(k,x,y);}
unsigned long cvtb2mask(){
 Z x=iota();
 x[2]|=128;
 x[4]|=64;
 return zo(cvtb2mask)(x);}
Z compressqi(){
 unsigned long k=0xff0fff0000000000;
 Z x=iota(),y=reva();
 return zm(compressqi)(x,y,k);
}

int main(){
 F(sqrtps());
 P(permvarqi());
 P(permvarsi());
 Q(vpermi2varps());
 P(vpermi2varqi());
 P(selectb_());
 printf("%lu\n",cvtb2mask());
 P(reva());
 P(compressqi());
}
