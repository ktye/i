#include<stdio.h>
#include<stdlib.h>
#include<string.h>

//tccrepl:gist.github.com/gbluma/1158a1af7a70db3cfce7

typedef struct{char*s;char *p;int n;char t;char v;char a;}k; //string/srcp/tokenlength;
int C[256]={0,0,0,0,0,0,0,0,0,0,13,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,6,2,2,2,2,3,0,1,2,2,2,8,4,12,5,5,5,5,5,5,5,5,5,5,9,0,2,2,2,2,2,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,0,11,1,2,2,7,4,4,4,4,10,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,0,2,1,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
int T[18][14]={{0,1,2,3,4,5,6,7,8,2,4,3,9,0},{0,1,2,3,4,5,6,7,2,2,4,3,3,0},{0,1,2,3,4,5,6,7,8,2,4,3,3,0},{0,1,2,3,4,5,6,7,8,11,4,3,3,0},{0,1,2,3,12,12,6,7,2,2,12,3,3,0},{0,1,2,3,13,13,6,7,2,11,14,3,3,0},{15,15,15,15,15,15,17,15,15,15,15,16,15,15},{0,1,2,3,12,5,6,7,2,2,12,3,3,0},{0,1,2,3,4,13,6,7,8,2,4,3,3,0},{10,10,10,10,10,10,10,10,10,10,10,10,10,0},{10,10,10,10,10,10,10,10,10,10,10,10,10,0},{0,1,2,3,4,5,6,7,8,2,4,3,3,0},{0,1,2,3,12,12,6,7,2,2,12,3,3,0},{0,1,2,3,13,13,6,7,2,2,14,3,3,0},{0,1,2,3,13,13,6,7,13,2,13,3,3,0},{15,15,15,15,15,15,17,15,15,15,15,16,15,15},{15,15,15,15,15,15,15,15,15,15,15,15,15,15},{0,1,2,3,4,5,6,7,2,2,4,3,3,0}};
char*sp;
char*nxt(){static int t;static char *c;if(!c)c=sp;while(*c)if(10>(t=T[t][C[(int)*c++]]))return c-1;return 0;}
int  tok(){static char*a,*b;a=a?b:nxt();while(a&&!C[a[0]])a=nxt();if(!a)return 0;b=nxt();sp=a;return b?b-a:strlen(a);}
char*glo,*loc;
char*cat(char*a,char*b){int na=strlen(a);int nb=strlen(b);char*r=malloc(na+nb+1);memcpy(r,a,na);memcpy(r+na,b,nb);r[na+nb]=0;return r;};char*c3(char*a,char*b,char*c){return cat(cat(a,b),c);}
char dya[]="i:i+i I:i+I";
char mon[]="i:+i i:-i i:@I I:|I";
void err(k x){printf("^ %s\n",x.s);exit(1);}

k dy(k x,k op,k y){char d[4];d[0]=x.t;d[1]=op.s[0];d[2]=y.t;d[3]=0;printf("dya %s\n",d);char *p=strstr(dya, d);if(!p)err(op);x.s=c3(x.s,op.s,y.s);if(0<x.a+y.a)x.s=c3("(",x.s,")");x.t=p[-2];x.a=0;return x;}
k mo(k p,k x){printf("mo %s %s\n",p.s,x.s);return p;}
k t(){int n;
 k r;r.s=0;n=tok();if(!n)return r;
 r.p=sp;r.s=malloc(1+n);memcpy(r.s,sp,1+n);r.s[n]=0;
 char c0=r.s[0];
 r.v=1==n&&C[c0]==2;
 r.t='i';if((C[c0]==5||c0=='-'&&n>1)&&strchr(r.s,'.'))r.t='f';
 r.a=1;
 return r;
}
k e(k x){
 if(!x.s)return x;
 k y=t();if(!y.s)return x;
 if(y.v&&!x.v){
  k r=e(t());return dy(x,y,r);
 }
 k r=e(y);
 return mo(x,r);
}




int main(int args,char**argv){
 glo=cat("",";");loc=cat("",";");
 if(args>1){
  sp=argv[1];
  k r=e(t());
  printf("v=%d t=%c a=%d s=%s\n", r.v, r.t, r.a, r.s);
 }
}
