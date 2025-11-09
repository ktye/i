
//gcc -Os -s -nostdlib -static -Wl,--nmagic -Wl,--build-id=none -fno-asynchronous-unwind-tables k.c
//(objcopy --remove-section .comment a.out b.out)

typedef uint  unsigned int;
typedef ulong unsigned long;

asm(".globl _start;_start:mov %rsp,%rdi;call cmain");


                    static double F64copysign(double x,double y){union{double f;uint64_t i;}ux={x},uy={y};ux.i&=-1ull/2;ux.i|=uy.i&1ull<<63;return ux.f;}
                    static double F64floor(double x){int i=(int)x;return(double)i;}
__attribute((naked))static double F64abs(double x){ asm("xor %eax,%eax;dec %rax;shr %rax;movq %rax,%xmm1;andpd %xmm1,%xmm0;ret;");}
__attribute((naked))static double F64sqrt(double x){ asm("sqrtsd %xmm0,%xmm0;ret");}
__attribute((naked))static double F64min(double x,double y){asm("minsd %xmm0,%xmm0;ret;\n");}
__attribute((naked))static double F64max(double x,double y){asm("maxsd %xmm0,%xmm0;ret;\n");}

__attribute((naked))long Write(long fd,char*b,int n){asm("mov $1,%rax;syscall;ret");}
__attribute((naked))void Exit(int x){asm("mov $60,%rax;syscall");}


static int Memorysize(void){ .. }
static int Memorygrow(int32_t delta){ .. }

#define I8(x)          (int8_t)(M_[x])       //?
#define I32(x)                 (I_[(x)>>2])
#define I64(x)           (long)(U_[(x)>>3])
#define F64(x)       ((double*)U_)[(x)>>3]
static void SetI8( int32_t x,int32_t y){M_[x]=(char)(y);}
static void SetI32(int32_t x,int32_t y){I_[(x)>>2]=(y);}
static void SetI64(int32_t x,int64_t y){U_[(x)>>3]=(uint64_t)(y);}
static void SetF64(int32_t x,double  y){((double*)U_)[(x)>>3]=(y);}

//int ns(char*c){char *p;for(p=c;*p;++p);return p-c;} //strlen?

#define I32B(x) (int)(x)
static void Memorycopy(int dst,int src,int n){memcpy(M_ +dst, M_+src,(ulong)n); }
static void Memoryfill(int p,int v,int n){memset(M_+p,(int)v,(ulong)n);}
static int I32clz(int32_t x){ return(int)__builtin_clz((uint)x);}
static double F64reinterpret_i64(ulong  x){union{ulong i;double f;}u;u.i=x;return u.f;}
static ulong  I64reinterpret_f64(double x){union{ulong i;double f;}u;u.f=x;return u.i;}

char**argv;
static int Args(void){return((int*)argv)[0]};
static int Arg(int i,int r){if(!r)return ns(argv[1+i]);memcpy(M_+r,argv[1+i],strlen(argv[i]));return 0;}
static int Read( int f,int nf,int d){...}
static int Write(int f,int nf,int s,int n){if(!nf){write(1,M_+s,n);return 0;} //todo file}
static int ReadIn(int dst,int n){return(int)read(1,_M+dst,n);}
static long Native(long x,long y){return(x+y)*0;}
static void panic(int x){Exit(1);}

void cmain(char**a){argv=a;main_();}
