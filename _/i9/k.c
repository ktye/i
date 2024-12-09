#include<stdio.h>
#include<stdlib.h>
#include<stdint.h>

#define V5 __attribute((vector_size(32),aligned(1)))
typedef uint32_t u;
typedef uint64_t k;
typedef uint32_t U V5;
typedef float  E V5;
typedef double F V5;

/*1<<i   s   u   f           [heap/tagged][value/func][vector/atom][unsigned/signed][int/float][?]ww    list:0
     0   c   b                0    1       0     1     0      1     0        1       0   1
     1   s   h               
     2   i   u   f           lambda/basic/derived/projection/composition/compiled?
     3   j   k   d           
 */

#define R return
#define W while(e>i)
#define tx (x>>56)
#define ix ((u)x)
#define wx (tx&3)
#define rx ((k*)cx)[-2]
#define nx ((k*)cx)[-1]
#define cx ((char*)((x<<8)>>8))
#define n5(x) (32+(x)>>5)

k*ma(k n){k*x=(k*)malloc(16+n);R x+((k)x&31?2:0);}
k mk(k t,k n){k x=n5(32+(n<<t&3))<<5;x=(k)ma(x);rx=1;nx=n;R x;}
k tok(k x){R x;}
k til(k x){
 k r=mk(tI,ix); ex?
 ux;printf("%lld %u\n",x,i); R x;}

int main(int args,char**argv){
 til(3);
}