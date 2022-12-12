
F:`!()                                         /function map: type to stringer


cc:{[o]opts::o
 F[`prg`mem`con`var`tab`fun`arg`sym`res`loc`ast`lod`sto]:(ccprg;ccmem;cccon;ccvar;cctab;ccfun;ccarg;ccsym;ccres;ccloc;ccast;cclod;ccsto)
 F[`add`and`bnd`ant`asn`cal`cli`cnd`dfr`div`drp`eql`cst]:(ccadd;ccand;ccbnd;ccant;ccasn;cccal;cccli;cccnd;ccdfr;ccdiv;ccdrp;cceql;cccst)
 F[`for`get`Get`gte`jmp`les`lit`lte`mod`mor`mul`neg    ]:(ccfor;ccget;ccGet;ccgte;ccjmp;ccles;cclit;cclte;ccmod;ccmor;ccmul;ccneg)
 F[`neq`nop`not`orr`bor`xor`ret`shl`shr`stm`sub`swc`typ]:(ccneq;ccnop;ccnot;ccorr;ccbor;ccxor;ccret;ccshl;ccshr;ccstm;ccsub;ccswc;cctyp)

 ccren'`main`ldexp`hypot                       /rename libm
 i:i@&{(0N~|/(T 1_c)?`swc`stm)&3=#c:&P=x}'i:&T=`cnd  /keep ternary
 rwcnd'|j@&~(j:&(~`=S)&(T=`cnd)|T=`swc)':i     /rewrite cnd/swc expressions 
 S[0]:$[0N~i:o?`lib;S 0;o 1+i]                 /emit`lib`name (instead of main + remove func main)
 n:(&#P){x[y]+:1}/P;n[0]-:1;                   /number of children per node
 s:{c:(0,(#x)-z)^x                             /cut args from stack
   (c 0),,,/(F T y)[y;|c 1]}/[();|!#P;|n]      /fold over the tree in reverse order
*s}


hx8:{[i]"0x",`x@|C i+!8}
em:{"(",x,")"}

INS:{[i;x]T::ins[i;T;*x];P::ins[i;@[P;&P>i-1;+;1];x 1];P[1+i]:i;I::ins[i;I;x 2];S::ins[i;S;x 3]}
ins:{[i;x;y](i#x),y,i_x}
del:{T::T n:(!#P)_x;P::@[P n;&(P n)>x-1;-;1];I::I n;S::S n}
rwcnd:{c:|1_&P=x;S[x]:`;P[x]:P x-1;c INS\:(T x;1+x;I x; S x:x-1);del x;restm'|i@&(T P i:&`stm=T)':`ret`asn}
restm:{c:*|&P=x;INS[c;(T x-1;x;I x-1;S x-1)];P[x]:P[x-1];del[x-1]} /fix ret stm

ccren:{S[&(S=x)&T':`fun`cal]:s:`$"_",$x;}

nl:{$[#x:"\n"/:x;x,"\n";x]}
nyi:{y;em x}

cctype:`i`u`j`k`f!("int32_t";"uint32_t";"int64_t";"uint64_t";"double")

ccprg:{
 f:(S i)!i:&T=`fun
 fpre:nl@ccsig'f@?S i@&i<f S i:&T=`cal
 fncs: ,/y@&`fun=T c:1_&P=0
 cons:nl y@&`con=T c
 vars:nl y@&`var=T c
 runt:$[`nort':opts;"";ccrnt]
 main:$[|/`nort`lib':opts;"";ccmain]
 data:$[#D;"memcpy(_M,\"",("\\x","\\x"/:(#D)^`x D),"\",",($#D),");";""]
 tabl:nl y@&`tab=T c
 tabl:$[#tabl;"_F=malloc(",($1+|/I@&`tab=T),"*sizeof(void*));\n";""],tabl
 mems:nl y@&`mem=T c
 init:$[|/(#mems;#data;#tabl);nl("void init(){\n";mems;data;tabl;"}\n");""]
 (runt;fpre;cons;vars),(,fncs),(ccfmt init),main}
 
ccfun:{
 $[(`main~S x)&`lib':opts; :""; 0]
 args:", "/:y@&`arg=T c:&x=P
 ret:$[#r:", "/:y@&`res=T c;r;"void"]
 locs:nl@{(cctype x)," ",y,";"}'[!g;.g:","/:'(y i)@=(S c)i:&`loc=T c]
 defr:     *y@&`dfr=T c
 setj:$[#defr;"if(!setjmp(_jb_)){";""]
 body: 2_-1_y@*&`ast=T c
 ccfmt@,/(ret;" ";$S x;"(";args;"){\n";setj;locs;body;defr)}

ccsig:{
 args:","/:cctype S c@&`arg=T c:&x=P
 ret:$[#r:cctype S c@&`res=T c;*r;"void"]
 ret," ",($S x),"(",args,");"}

ccmem:{y;("Memory";$[`b~S x;"2";""];"(";$I x;");")}
cctab:{y;("_F[";$I x;"]=";$S x;";\n")}

/floating point literals that don't roundtrip (written by mk flt)
ccflts:(0x25fab843dda06440;0x2cde25af6a56113f;0x182d4454fb21e93f;0x485555555555c5bf;0xcd3b7f669ea0e63f;0x935555555555e53f;0xc64bac7e4f7e92be;0x4b5555555555a53f;0x5d1f29a9e5e55abe;0xf544c819a001fa3e;0x9b1a86a049faa8bd;0xe69d3f334f500340;0xec62bdfb9e516840;0xf16bd2c541bdbbbe;0xaf788e1dc571cc3f;0x182d4454fb210940;0x93bdbe166cc166bf;0x7a806b5b542830c0;0xcd9cd11ffdd8e53d;0xea49b0133f567e40;0xfe822b654715f73f;0x000000000000a03c;0x0000000000001000;0x075c143326a6913c;0x8eec28fd693650c0;0x00000000000030be;0x04fa97999999d93f;0xff799f501344d33f;0xa1487d56e31dc73e;0x03dfbf19a0012abf;0xd0f710111111813f;0x9425f7a17f00ecbf;0x730288368cc052c0;0x3bbee2d2180e7b40;0xef39fafe422e8640;0x44523edf12f1c23f;0x000000002d44643e;0x7051cc989846e83c;0x25ba052dbfb85ec0;0x51302dd5104987c0;0x0000000000003043;0x399d52a246df913f;0x00000040fb21e93f;0x555555555555c53f;0xde03cb966446c73f;0xffffffffffffef7f;0x0000e0fe422ee63f;0xd0a4be726937663e;0x9fc678d0099ac33f;0x000000000000e043;0x053f4e7b9dee213e;0x182d4454fb21f93f;0x3c60145bc4db3840;0x0ee526157bcbdb3f;0x914fc1166cc156bf;0x000000000000303e;0x83c8c96d305ff43f;0x763c7935ef39ea3d;0x599322942449d23f;0x8dedb5a0f7c6b03e;0xf8c1631adca54c40)!("165.02700983169885";"6.613756321437934e-05";"0.7853981633974483";"-0.1666666666666663";"0.7071067811865476";"0.6666666666666735";"-2.755731417929674e-07";"0.041666666666666595";"-2.5050747762857807e-08";"2.4801587288851704e-05";"-1.1358536521387682e-11";"2.414213562373095";"194.5506571482614";"-1.6533902205465252e-06";"0.22222198432149784";"3.141592653589793";"-0.0027777777777015593";"-16.157537187333652";"1.5896230157654656e-10";"485.3903996359137";"1.4426950408889634";"1.1102230246251565e-16";"2.2250738585072014e-308";"6.123233995736766e-17";"-64.85021904942025";"-3.725290298461914e-09";"0.3999999999940942";"0.3010299956639812";"2.7557313621385722e-06";"-0.0001984126982958954";"0.008333333333322118";"-0.8750608600031904";"-75.00855792314705";"432.88106049129027";"709.782712893384";"0.14798198605116586";"3.774894707930798e-08";"2.6951514290790595e-15";"-122.88666844901361";"-745.1332191019411";"4.503599627370496e+15";"0.017453292519943295";"0.7853981256484985";"0.16666666666666666";"0.1818357216161805";"1.7976931348623157e+308";"0.6931471803691238";"4.1381367970572385e-08";"0.15313837699209373";"9.223372036854776e+18";"2.087570084197473e-09";"1.5707963267948966";"24.858464901423062";"0.4342944819032518";"-0.0013888888888873056";"3.725290298461914e-09";"1.2732395447351628";"1.9082149292705877e-10";"0.2857142874366239";"1e-06";"57.29577951308232")
ccflt:{x:C x+!8;$[#r:ccflts x;r;$*`f x]}
cclit:{hi:{*`i C 4+x+!4};lo:{*`i C x+!4}
 y;i:I x;$[`f~t:S x;ccflt i
           `i~t;$[0N~i;"-0x80000000";$i]
           `u~t;$[i<0;"0x",`x@|`c@,i;$i]
           `j~t;$[i~0;,"0";(-1~h:hi i)&0>l:lo i;$l;(h~0)&l>0;$l;hx8 i],"ll"
           `k~t;$[i~0;,"0";(0~h:hi i)&0<l:lo i;$l;hx8 i],"ull";"???"]}
 

ccvar:{" "/:(cctype S 1+x;  $S x;,"=";(*y),";")}
cccon:{"const ",(cctype S 1+x)," ",($S x),"=",(*y),";"}

cccst:{("(";cctype S x;")";,"(";y 1;,")")}

cctyp:{y;cctype S x} 
ccres:{y;cctype S x}
ccloc:{y}
ccsym:{y;$S x}
ccarg:{(cctype S x;" ";*y)}

ccstm:{ "{\n",(,/y,\:";\n"),"}\n"}
/ccstm:{"{\n",(";\n"/:y),"}\n"}
ccast:ccstm
ccget:{y;$S x}
ccGet:ccget
ccret:{"return ",(", "/:y),";\n"}

ccsiz:`b`i`j`f!`I8`I32`I64`F64
cclod:{($ccsiz S x;"(";y 0;")")}
ccsto:{("Set";$ccsiz S x;"(";y 0;",";y 1;")")}

cccal:{($S x;,"("; ","/:y ;,")")}
ccdrp:{y}

cccli:{n:I x;("((";$[#r:cctype S x;r;"void"];"(*)("; ","/:n#(1+n)_y; "))_F[";*y;"])(";  ","/:n#1_y; ")")}
ccasn:{r:(s:$S x),"=",(*|y),";\n";
 $[((n:3+2*#s)#r)~s,"=(",s,"+";$["+=1"~-3#r:s,"+=",(n_-3_r);s,"++";r],";\n";r]}
cccnd:{$[S x; (ccnd@*y;"?(";y 1;"):(";y 2;")") ;"if",(ccnd y 0),(ccemb ccsem y 1),$[3~#y;"else",(ccemb ccsem y 2);""]]}
ccnop:{y;""}
ccswc:{("switch(";*y;"){\n"),(1_(@[("case ",/:$-1+!#y); (-1+#y)*I x; "default"] ,\:":"),'(ccsem'y),\:"break;\n"),"}" }
 
ccfor:{$[I x;ccwhl[x;y];("for(;";ccnd y 0;";";$[";\n"~-2#r:y 1;-2_r;r];")";y 2)]}
ccwhl:{((,"do"),(,-12_y 2),(,y 1)),("}while";ccnd y 0;";\n")}
ccjmp:{y;$[I x;"break";"continue"]}
ccdfr:{("else{";*y;";}\n}")}

cc1:{x,*z}
cc2:{z;("(";z 0;x;z 1;")")}
ccneg:cc1"-"; ccnot:cc1"!"
ccadd:cc2"+"; ccsub:cc2"-"; ccmul:cc2"*"; ccdiv:cc2"/"; ccmod:cc2"%"
cceql:cc2"==";ccgte:cc2">=";ccles:cc2"<"; cclte:cc2"<=";ccmor:cc2">"; ccneq:cc2"!="
ccand:cc2"&"; ccbnd:cc2"&&";ccant:cc2"&~";ccorr:cc2"|"; ccbor:cc2"||";ccxor:cc2"^"; ccshl:cc2"<<";ccshr:cc2">>"

ccnd:{$["!=0)"~t:-4#x;(-4_x),")";"==0)"~t;"(!",(1_-4_x),")";x]}
ccsem:{$[(*|x)':";}\n";x;x,";"]}
ccemb:{$["{"=*x;x;"{",x,"}"]}

ccfmt:{
 x:@[x;1+find["}\nelse";x];" "]
 x:@[x;1+find["}\nwhil";x];" "]
 x:@[x;1+find["};\n"   ;x];" "]
 l:"\n"\:x;
 l@:&~l~\:"r=r;"
 l@:&~(""~/:l)|(,";")~/:l;s:_32+&'(-"}"=*'l)+-1_0,+\+/'("{"=l)-"}"=l
 ,/(s,'l),\:"\n"}

/ run time
ccrnt:"#include<stdio.h>
#include<stdlib.h>
#include<stdint.h>
#include<string.h>
#include<setjmp.h>
#include<math.h>
#define F64abs       fabs
#define F64sqrt      sqrt
#define F64floor     floor
#define F64min       fmin
#define F64max       fmax
#define F64copysign  copysign
#define Exit exit
char *_M, *_M2;
void **_F;
int32_t _memorysize, _memorysize2;
FILE *_fd_;
int    args_;
char **argv_;
static jmp_buf _jb_;
void Memory (int32_t x){_memorysize =x;_M =calloc(x,64*1024);}
void Memory2(int32_t x){_memorysize2=x;_M2=calloc(x,64*1024);}
int32_t Memorysize( ){return _memorysize; }
int32_t Memorysize2(){return _memorysize2;}
int32_t Memorygrow(int32_t delta){
 int32_t r=_memorysize;
 _memorysize+=delta;
 _M=(char *)realloc(_M, 64*1024*_memorysize);
 return r;
}
int32_t Memorygrow2(int32_t delta){
 int32_t r=_memorysize2;
 _memorysize2+=delta;
 _M2=(char *)realloc(_M2, 64*1024*_memorysize2);
 return r;
}
#define I8(x)       ((int8_t *)_M)[x]
#define I32(x)      ((int32_t*)_M)[(x)>>2]
#define I64(x)      ((int64_t*)_M)[(x)>>3]
#define F64(x)      ( (double*)_M)[(x)>>3]
#define SetI8(x,y)  ((int8_t *)_M)[x]=(y)
#define SetI32(x,y) ((int32_t*)_M)[(x)>>2]=(y)
#define SetI64(x,y) ((int64_t*)_M)[(x)>>3]=(y)
#define SetF64(x,y) (( double*)_M)[(x)>>3]=(y)
#define I32B(x)     (int32_t)(x)
void Memorycopy (int32_t dst, int32_t src, int32_t n){ memcpy(_M +dst, _M +src, n); }
void Memorycopy2(int32_t dst, int32_t src, int32_t n){ memcpy(_M2+dst, _M +src, n); }
void Memorycopy3(int32_t dst, int32_t src, int32_t n){ memcpy( _M+dst, _M2+src, n); }
void Memoryfill(int32_t p, int32_t v, int32_t n){ memset(_M+p, (int)v, (size_t)n); }
int32_t  I32clz(uint32_t x) { return (int32_t)__builtin_clz((unsigned int)x); }
int32_t  I64popcnt(uint64_t x){ return (int32_t)__builtin_popcountll(x); }
double   F64reinterpret_i64(uint64_t x){return *(double  *)&x;}
uint64_t I64reinterpret_f64(double   x){return *(uint64_t*)&x;}
int32_t Args(){ return args_; }
int32_t Arg(int32_t i, int32_t r){
 if(i>=args_) return 0;
 if(r ==   0) return strlen(argv_[i]);
 memcpy(_M+r,argv_[i],strlen(argv_[i]));
 return 0;
}
int32_t Read(int32_t file, int32_t nfile, int32_t dst){
 static char *filebuf = NULL;
 static size_t      n = 0;
 if(dst != 0){ memcpy(_M+dst,filebuf,n); return 0; }
 char name[512];
 if(nfile > 511) return -1;
 memcpy(name, _M+file, nfile);
 name[nfile] = (char)0;
 FILE *fp = fopen(name, \"rb\");
 if(fp==NULL){if(filebuf!=NULL)free(filebuf);n=0;return -1;}
 fseek(fp, 0, SEEK_END);
 n = ftell(fp);
 fseek(fp, 0, SEEK_SET);
 if(filebuf != NULL) free(filebuf);
 filebuf = malloc(n);
 if(n != fread(filebuf, 1, n, fp)){ fclose(fp); return -1; }
 fclose(fp);
 return n;
}
int32_t Write(int32_t file, int32_t nfile, int32_t src, int32_t n){
 if(nfile == 0){ fwrite(_M+src, 1, n, stdout); return 0; }
 char name[512];
 memcpy(name, _M+file, nfile);
 name[nfile] = (char)0;
 FILE *fp = fopen(name, \"wb\");
 if(fp == NULL){ return -1; }
 fwrite(_M+src, 1, n, fp);
 fclose(fp);
 return 0;
}
int32_t ReadIn(int32_t dst, int32_t n){
 char *r = fgets(_M+dst, n, stdin);
 if(r==NULL){ //todo eof
  return 0;
 }else return (int32_t)strnlen(_M+dst,n);
}
int64_t Native(int64_t x, int64_t y){
#ifdef NATIVE
 return cnative(x, y);
#else
 return 0*(x+y);
#endif
}
void panic(int32_t x) { longjmp(_jb_,1); }
/*
#define P(...) printf(__VA_ARGS__);fflush(stdout);
void Printi(int32_t x, int32_t y){ fprintf(stderr,\"%d %d\\n\",x,y);fflush(stderr); }
void Prinf(int32_t x, double y){ fprintf(stderr,\"%d %.6lf\\n\",x,y);fflush(stderr); }
void printk_(uint64_t y){
 y=y&~(uint64_t)1;
 if((y>>59)>16){
  int32_t n=I32(((int32_t)(y))-12);
  fprintf(stderr,\"(%d#%u) \",n,(uint32_t)(y>>59));
 }else fprintf(stderr,\"(%u) \",(uint32_t)(y>>59));
}
void Printk(int32_t x, uint64_t y){
 fprintf(stderr,\"%d \", x);printk_(y);fprintf(stderr,\"\\n\");fflush(stderr);
}
void Printu(int32_t x, uint64_t y){
 fprintf(stderr,\"%d %llu\\n\", x, y);fflush(stderr);
}
void Printl(int32_t x, uint64_t y){
 y=y&~(uint64_t)1;
 if(23!=y>>59){
  fprintf(stderr,\"%dL !%llu\\n\",x,y>>59);
  fflush(stderr);return;
 }
 int32_t yp=(int32_t)y;
 int32_t n=I32(yp-12);
 fprintf(stderr,\"%dL %d# \",x,n);fflush(stderr);
 if(n>5)n=5;
 int32_t e=yp+8*n;
 for(;yp<e;yp+=8)printk_(I64(yp));
 fprintf(stderr,\"\\n\");
 fflush(stderr);
}
*/
"
ccmain:"int main(int args, char **argv){
 args_=(int32_t)args;
 argv_=argv;
 init();
 _main();
 return 0;
}
"
