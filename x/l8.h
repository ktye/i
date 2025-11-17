//gcc -Os -s -nostdlib -static -Wl,--nmagic -Wl,--build-id=none -fno-asynchronous-unwind-tables k.c
//(objcopy --remove-section .comment a.out b.out)

typedef unsigned int  uint;
typedef unsigned long ulong;

asm(".globl _start;_start:mov %rsp,%rdi;call cmain");

__attribute((naked))long brk(long x){asm("mov $12,%rax;syscall;ret");}
__attribute((naked))int  open(char*file,int flags,uint mode){asm("mov $2,%rax;syscall;ret");} //flag(rd0 wr|cr|tr577) mode(0644 420)
__attribute((naked))void close(int fd){asm("mov $3,%rax;syscall;ret");}
__attribute((naked))uint lseek(int fd,ulong o,uint w){asm("mov $8,%rax;syscall;ret");}        //0(set) 2(end)
__attribute((naked))long sread( long fd,char*b,uint n){asm("xor %rax,%rax;syscall;ret");}
__attribute((naked))long swrite(long fd,char*b,uint n){asm("mov $1,%rax;syscall;ret");}
__attribute((naked))void Exit(int x){asm("mov $60,%rax;syscall");}

                    static double F64copysign(double x,double y){union{double f;ulong i;}ux={x},uy={y};ux.i&=-1ul/2;ux.i|=uy.i&1ul<<63;return ux.f;}
                    static double F64floor(double x){double f=(double)(long)x;return x==f?f:x<0?f-1.0:f;}
__attribute((naked))static double F64abs(double x){ asm("xor %eax,%eax;dec %rax;shr %rax;movq %rax,%xmm1;andpd %xmm1,%xmm0;ret;");}
__attribute((naked))static double F64sqrt(double x){ asm("sqrtsd %xmm0,%xmm0;ret");}
__attribute((naked))static double F64min(double x,double y){asm("minsd %xmm1,%xmm0;ret;\n");}
__attribute((naked))static double F64max(double x,double y){asm("maxsd %xmm1,%xmm0;ret;\n");}

/*
static void wl(ulong x){ //debug-only
 char s[32];s[31]=10;s[30]='0';
 if(!x){swrite(1,s+30,2);return;}
 int i=31;while(x){s[--i]='0'+(x%10);x/=10;}
 swrite(1,s+i,32-i);}
*/

static char*M_;
static ulong brk0;
#define I_ ((int*)M_)
#define U_ ((ulong*)M_)
static void*F_[];   //dispatch table
#define int32_t int //some literals
static ulong pages_=0,pages__=0;  //a page is 64k
static int Memorysize(void){return pages_;}
static int Memorygrow(int d){pages_+=d;if(pages_>16384)return -1;brk(brk0+(pages_<<16));return pages_-d;}
static void Memory(int x){brk0=brk(0);M_=(char*)brk0;Memorygrow(x);}
static void Memory2(int x){}

//try/catch
static void*jb_[5];static int jb__=0;
#define setjmp  __builtin_setjmp
#define longjmp __builtin_longjmp
static void panic(int x){if(!jb__)Exit(1);longjmp(jb_,1); }
static void store(void){ulong n=pages_<<13;if(pages__){/*copy backwards*/ulong o=((ulong)M_-brk0)>>3;for(int i=0;i<n;i++)U_[i-o]=U_[i];M_=(char*)brk0; }
 pages__=pages_;brk(brk0+(pages_<<17));for(int i=0;i<n;i++)U_[i+n]=U_[i];M_=(char*)(brk0+(pages_<<16));}
static void catch(void){pages_=pages__;pages__=0;M_=(char*)brk0;brk(brk0+(pages_<<16));}
/* k repl does:
   main(){ store(); for(;;){ write(Ku(32)); ulong x=readfile(mk(Ct,0)); try(x); }}
   try(ulong x){ jb__=1;if(!setjmp(jb_)){ repl(x); store(); }else{catch();      }}

 brk0<--pages_-->brk
  M_...
  +--------------+     without try/catch   M_ is always brk0, brk increases occasionally


 brk0            M_            brk
  +--------------+--------------+                  store() initial call: double memory, copy forward set new M_
         +----------------^ copy forward left-to-right is ok

 brk0   pages__            pages_           brk
  +--------------+---------------------------+     during execution (try) active k increases memory (pages_>pages__)
     old k mem           active k mem


 brk0   pages__            
  M_                      brk
  +------------------------+-----------------+     store() (no error) copy backwards, reset M_, shrink brk
          ^---------------------+  left-to-right

 brk0      pages__         M_         pages_          brk   also in store() prepare for next try:
  +------------------------+---------------------------+    increase brk, copy forward


catch():
 brk0     pages_          brk
  M_                              in case of errors:
  +------------------------+      reset pages_, M_ and brk,  pages__=0
 */


#define I8(x)             (int)(M_[x])
#define I32(x)                 (I_[(x)>>2])
#define I64(x)           (long)(U_[(x)>>3])
#define F64(x)       ((double*)U_)[(x)>>3]
static void SetI8( int x,int y){M_[x]=(char)(y);}
static void SetI32(int x,int y){I_[(x)>>2]=(y);}
static void SetI64(int x,long y){U_[(x)>>3]=(long)(y);}
static void SetF64(int x,double  y){((double*)U_)[(x)>>3]=(y);}

void memcpy(char*d,char*s,int n){for(;n;n--)*d++=*s++;}
int strlen(char*c){char *p;for(p=c;*p;++p);return p-c;}

#define I32B(x) (int)(x)
static void Memorycopy(int d,int s,int n){memcpy(M_ +d, M_+s,(ulong)n); }
static int I32clz(int x){ return(int)__builtin_clz((uint)x);}
static double F64reinterpret_i64(ulong  x){union{ulong i;double f;}u;u.i=x;return u.f;}
static ulong  I64reinterpret_f64(double x){union{ulong i;double f;}u;u.f=x;return u.i;}

char**argv;
static int Args(void){return((int*)argv)[0];};
static int Arg(int i,int r){if(!r)return strlen(argv[1+i]);memcpy(M_+r,argv[1+i],strlen(argv[1+i]));return 0;}
static int Read( int f,int nf,int d){static int fd=0;static uint sz=0;if(d!=0){sread(fd,M_+d,sz);close(fd);return 0;}
 char c=M_[f+nf];M_[f+nf]=0;fd=open(M_+f,0,0);M_[f+nf]=c;if(fd<0)return 0;sz=lseek(fd,0,2);lseek(fd,0,0);return(int)sz;}
static int Write(int f,int nf,int s,int n){if(!nf){swrite(1,M_+s,n);return 0;}
 char c=M_[f+nf];M_[f+nf]=0;int fd=open(M_+f,577,420);M_[f+nf]=c;if(fd<0)return -1;swrite(fd,M_+s,n);close(fd);return 0;}
static int ReadIn(int d,int n){return(int)sread(1,M_+d,n);}
static long Native(long x,long y){return(x+y)*0;}
