// -nostdlib -ffreestanding

//gcc -nostdlib -static

#include<stdint.h>
asm(".global _start;_start:;xorl %ebp,%ebp;movq 0(%rsp),%rdi;lea 8(%rsp),%rsi;call main;");
asm("Exit:;movq %rax,%rdi;movq $60,%rax;syscall;\n");
asm("F64sqrt:sqrtsd %xmm0,%xmm0;ret;\n");

void Exit(int x);
double F64sqrt(double);

/*
asm("write:;movq $1,%rax;syscall;ret;\n");
asm("F64abs:xor %eax,%eax;dec %rax;shr %rax;movq %rax,%xmm1;andpd %xmm1,%xmm0;ret;");
asm("F64min:minsd %xmm0,%xmm0;ret;\n");
asm("F64max:maxsd %xmm0,%xmm0;ret;\n");

double F64abs(double);
double F64min(double);
double F64max(double);
double F64copysign(double x,double y){union{double f;uint64_t i;}ux={x},uy={y};ux.i&=-1ull/2;ux.i|=uy.i&1ull<<63;return ux.f;}
double F64floor(double x){int i=(int)x;return(double)i;}

int write(int fd, const void *buf, unsigned count);
int ns(char*c){char *p;for(p=c;*p;++p);return p-c;}
*/

void main(int args,char*argv[]){
	/*
 for(int i=0;i<args;++i){
  int n=ns(argv[i]);
  write(1,argv[i],n);
  write(1,"\n",1);
 }
 */
 Exit((int64_t)F64sqrt((double)args));
}
