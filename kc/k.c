#include<stdio.h>
#include<stdlib.h> // 1+x:|x     _=use(x);n_=nn_(x);N(x){_[i]=1+x[n_+i-1]}
#include<string.h>

//tccrepl:gist.github.com/gbluma/1158a1af7a70db3cfce7

typedef struct{char*s;char *h;char *o;char *p;int n;char t;char v;char a;char f;char e;char*N;}k; //string/head/srcp/tokenlength/type/verb/atomic/func/side-effect
#include"p.h"
#define ax (x.t>'Z')
#define ay (y.t>'Z')
#define c3(a,b,c) c2(a,c2(b,c))
#define c4(a,b,c,d) c2(c2(a,b),c2(c,d))
#define c5(a,b,c,d,e) c2(c2(a,b),c3(c,d,e))
#define c6(a,b,c,d,e,f) c2(c3(a,b,c),c3(d,e,f))
int  cl(char c){return-97+(int)C[(int)c];}
char*sp;
char*nxt(){static int t;static char*c;if(!c)c=sp;while(*c)if(11>(t=T[t][cl(*c++)]-97))return c-1;return 0;}
int  tok(){static char*a,*b;a=a?b:nxt();while(a&&!cl(a[0]))a=nxt();if(!a)return 0;b=nxt();sp=a;return b?b-a:strlen(a);}
char*glo,*loc;
int  in(char*a,char c){return strchr(a,c)!=0;}
char*c2(char*a,char*b){if(!a)return b;if(!b)return a;int na=strlen(a);int nb=strlen(b);char*r=malloc(na+nb+1);memcpy(r,a,na);memcpy(r+na,b,nb);r[na+nb]=0;return r;};
char*br(char*a,char*b){int n=strlen(a);char*r=malloc(n+3);memcpy(1+r,a,n);r[0]=b[0];r[1+n]=b[1];r[2+n]=0;return r;}char*em(char*a){return br(a,"()");}
int  ssn=0;char*ssa(){char*r=malloc(3);r[0]='_';r[1]=ssn<10?'0'+ssn:ssn<36?'a'+ssn-10:'A'+ssn-36;r[2]=0;ssn++;return r;}
char sca[]="+-*%&|<>=";
char dya[]="i:i+i i:i*i I:i+I I:I+i I:I+I";
char mon[]="i:+i i:-i i:@I I:|I I:!i I:&i F:?i";
#define E(x) {printf("^%s %s\n",__func__,x.s);exit(1);}

k e(k);
k L(k x,char*n){if(!x.N){if(!x.a)x.o=c2(x.o,c2(x.h,x.s));}else if((!n)||!strcmp(x.N,n)){x.o=c2(x.o,c4("N(",x.N,"){",c4(x.h,"_[i_]=",x.s,";};")));x.h=0;x.s="_";};x.N=n;return x;}
k f1(k x,k y){char*s=c2(x.s,br(y.s,"()")); x.s=s;x.h=y.h;x.e|=y.e;return x;}
k at(k x,k y){if(x.s)return f1(x,y);x.s=c2(x.s,br(y.s,"[]"));return x;}
k ge(k x,k y){x.h=c2(x.h,y.h);char o=x.s[0];x.s=o=='!'?"i_":o=='&'?"0":"rf_()";x.N=y.s;x.t='I';return x;}
k dy(k x,k op,k y){if(op.s[0]=='@')return at(x,y);char d[4];d[0]=x.t;d[1]=op.s[0];d[2]=y.t;d[3]=0;char*p=strstr(dya,d);if(!p)E(op);

 if(!(ax+ay)){y.o=c2(c3(y.o, "xy_(", x.N),c3( ",", y.N, ");" ));  };y=ay?L(y,x.N):y;x=ax?L(x,y.N):x;
 x.s=c3(x.a?x.s:em(x.s),op.s,y.a?y.s:em(y.s));x.h=c2(y.h,x.h);x.o=c2(y.o,x.o);x.N=y.N;x.t=p[-2];x.a=0;return x;}
 
k mo(k op,k x){char o=op.s[0];if(x.t=='i'&&in("!?&",o))return ge(op,x);char d[3];d[0]=o;d[1]=x.t;d[2]=0;char*p=strstr(mon,d);if(!p)E(op); op.s=c2(op.s,x.a?x.s:em(x.s));op.t=p[-1];op.a=0;op.h=x.h;return op;}

k t(){int n;
 k r;r.s=0;r.h=0;r.o=0;r.e=0;r.N=0;n=tok();if(!n)return r;if(2==cl(sp[0]))return r;
 r.p=sp;r.s=malloc(1+n);memcpy(r.s,sp,1+n);r.s[n]=0;
 char c0=r.s[0];
 //printf("c0=%c\n",c0);
 if(c0=='('){r=e(t());r.v=0;return r;}
 r.v=1==n&&(3==cl(c0)||c0=='-');
 r.t='i';if((6==cl(c0)||(c0=='-'&&1<n))&&strchr(r.s,'.'))r.t='f';
 r.a=1;r.f=(c0=='f');
 //printf("t %c cl=%d v=%d a=%d n=%d\n",c0,cl(c0),r.v,r.a,n);
 return r;
}
k e(k x){
 if(!x.s)return x;
 k y=t();if(!y.s)return x;
 printf("yv=%d xv=%d\n",y.v,x.v);
 if(y.v&&!x.v){
  k r=e(t());return dy(x,y,r);
 }
 k r=e(y);
 return x.v?mo(x,r):at(x,r);
}




int main(int args,char**argv){
 glo=c2("",";");loc=c2("",";");
 if(args>1){
  //char*c;sp=argv[1];while((c=nxt())){printf("%s\n",c);}exit(0);
  //int n;sp=argv[1];while((n=tok())){printf("%d %s\n",n,sp);}exit(0);
  sp=argv[1];
  k r=e(t());
  printf("v:%d t:%c a:%d N:%s, h:%s s:%s\n", r.v, r.t, r.a, r.N, r.h?r.h:"", r.s);
  r=L(r,0);printf("%s %s %s\n",r.o,r.h,r.s);
 }
}
