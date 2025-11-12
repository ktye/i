
//gcc -Os -s -nostdlib -static -Wl,--nmagic -Wl,--build-id=none -fno-asynchronous-unwind-tables k.c
//(objcopy --remove-section .comment a.out b.out)

typedef unsigned int  uint;
typedef unsigned long ulong;

asm(".globl _start;_start:mov %rsp,%rdi;call cmain");


                    static double F64copysign(double x,double y){union{double f;ulong i;}ux={x},uy={y};ux.i&=-1ul/2;ux.i|=uy.i&1ul<<63;return ux.f;}
                    static double F64floor(double x){long i=(long)x;return(long)i;}
__attribute((naked))static double F64abs(double x){ asm("xor %eax,%eax;dec %rax;shr %rax;movq %rax,%xmm1;andpd %xmm1,%xmm0;ret;");}
__attribute((naked))static double F64sqrt(double x){ asm("sqrtsd %xmm0,%xmm0;ret");}
__attribute((naked))static double F64min(double x,double y){asm("minsd %xmm0,%xmm0;ret;\n");}
__attribute((naked))static double F64max(double x,double y){asm("maxsd %xmm0,%xmm0;ret;\n");}

__attribute((naked))long sysread( long fd,char*b,int n){asm("xor %rax,%rax;syscall;ret");}
__attribute((naked))long syswrite(long fd,char*b,int n){asm("mov $1,%rax;syscall;ret");}
__attribute((naked))void Exit(int x){asm("mov $60,%rax;syscall");}

extern char _end[];
#define M_ _end
#define I_ ((int*)M_)
#define U_ ((ulong*)M_)
static void*F_[];   //dispatch table
#define int32_t int //some literals
static int Memorysize(void){ return 2; } //todo brk
static int Memorygrow(int delta){ return 2; } //todo
static void Memory(int x){}
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
static void Memorycopy(int dst,int src,int n){memcpy(M_ +dst, M_+src,(ulong)n); }
static int I32clz(int x){ return(int)__builtin_clz((uint)x);}
static double F64reinterpret_i64(ulong  x){union{ulong i;double f;}u;u.i=x;return u.f;}
static ulong  I64reinterpret_f64(double x){union{ulong i;double f;}u;u.f=x;return u.i;}

char**argv;
static int Args(void){return((int*)argv)[0];};
static int Arg(int i,int r){if(!r)return strlen(argv[1+i]);memcpy(M_+r,argv[1+i],strlen(argv[i]));return 0;}
static int Read( int f,int nf,int d){ return 0; } //todo
static int Write(int f,int nf,int s,int n){
 if(!nf){syswrite(1,M_+s,n);return 0;}
 return 0;} //todo file
static int ReadIn(int dst,int n){return(int)sysread(1,M_+dst,n);}
static long Native(long x,long y){return(x+y)*0;}
static void panic(int x){Exit(1);}

//void cmain(char**a){argv=a;main_();}
