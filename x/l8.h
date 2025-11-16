
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

//uint csr=0xbfc0;      // ldmxcsr(&csr); (must be a memarg)
//__attribute((naked))void ldmxcsr(uint*x){asm("ldmxcsr (%rdi);ret");}
                    static double F64copysign(double x,double y){union{double f;ulong i;}ux={x},uy={y};ux.i&=-1ul/2;ux.i|=uy.i&1ul<<63;return ux.f;}
                    static double F64floor(double x){long i=(long)(x<0?x-1.:x);return(long)i;}
__attribute((naked))static double F64abs(double x){ asm("xor %eax,%eax;dec %rax;shr %rax;movq %rax,%xmm1;andpd %xmm1,%xmm0;ret;");}
__attribute((naked))static double F64sqrt(double x){ asm("sqrtsd %xmm0,%xmm0;ret");}
__attribute((naked))static double F64min(double x,double y){asm("minsd %xmm1,%xmm0;ret;\n");}
__attribute((naked))static double F64max(double x,double y){asm("maxsd %xmm1,%xmm0;ret;\n");}

static void wl(ulong x){ //debug-only
 char s[32];s[31]=10;s[30]='0';
 if(!x){swrite(1,s+30,2);return;}
 int i=31;while(x){s[--i]='0'+(x%10);x/=10;}
 swrite(1,s+i,32-i);}

static char*M_;
#define I_ ((int*)M_)
#define U_ ((ulong*)M_)
static void*F_[];   //dispatch table
#define int32_t int //some literals
static ulong pages_=0;
static int Memorysize(void){return pages_;}
static int Memorygrow(int d){pages_+=d;if(pages_>16384)return -1;brk((pages_<<16)+(ulong)M_);return pages_-d;}
static void Memory(int x){M_=(char*)brk(0);Memorygrow(x);}
static void Memory2(int x){}

//todo try/catch
static int jb_,jb__;
static int setjmp(int x){return 0;}
static int Memorycopy2(int x,int y,int z){return 0;}
static int Memorycopy3(int x,int y,int z){return 0;}
static int Memorysize2(void){return 2;}
static int Memorygrow2(int x){return 2;}


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
static void panic(int x){Exit(1);}

