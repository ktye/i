#include<stdio.h>
#include<stdlib.h>
#include<string.h>

//tccrepl:gist.github.com/gbluma/1158a1af7a70db3cfce7


int C[256]={0,0,0,0,0,0,0,0,0,0,13,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,6,2,2,2,2,3,0,1,2,2,2,8,4,12,5,5,5,5,5,5,5,5,5,5,9,0,2,2,2,2,2,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,0,11,1,2,2,7,4,4,4,4,10,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,0,2,1,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
int T[18][14]={{0,1,2,3,4,5,6,7,8,2,4,3,9,0},{0,1,2,3,4,5,6,7,2,2,4,3,3,0},{0,1,2,3,4,5,6,7,8,2,4,3,3,0},{0,1,2,3,4,5,6,7,8,11,4,3,3,0},{0,1,2,3,12,12,6,7,2,2,12,3,3,0},{0,1,2,3,13,13,6,7,2,11,14,3,3,0},{15,15,15,15,15,15,17,15,15,15,15,16,15,15},{0,1,2,3,12,5,6,7,2,2,12,3,3,0},{0,1,2,3,4,13,6,7,8,2,4,3,3,0},{10,10,10,10,10,10,10,10,10,10,10,10,10,0},{10,10,10,10,10,10,10,10,10,10,10,10,10,0},{0,1,2,3,4,5,6,7,8,2,4,3,3,0},{0,1,2,3,12,12,6,7,2,2,12,3,3,0},{0,1,2,3,13,13,6,7,2,2,14,3,3,0},{0,1,2,3,13,13,6,7,13,2,13,3,3,0},{15,15,15,15,15,15,17,15,15,15,15,16,15,15},{15,15,15,15,15,15,15,15,15,15,15,15,15,15},{0,1,2,3,4,5,6,7,2,2,4,3,3,0}};
char *src,tk[32];

int tok(){
 static int   t;
 static char *c;
 if(!c)c=src;
 int i=0;while(*c){tk[i++]=*c; if(10>(t=T[t][C[(int)*c++]]))break;}
 tk[i]=0;
 return i;
}

void parse(char *s){
	src=s;
	int n;
	while((n=tok())){
	 printf("%d %s\n",n,tk);
	}
// int n=strlen(s);
// int t=0;for(int i=0;i<n;i++)if(10>(t=T[t][C[s[i]]]))printf("%s\n",s+i);
}

/*
fvvvwww
    000  0  char
    001  1  short
    010  2  int
    011  3  double/k/d/t
    100  4  z
  01000  8  flat vector
 010000 16  compound
 100011 35  list
 110011     dict
 111111     table
1000000     primitive
1100000     function
*/


typedef unsigned long long K;

#define R return r;
#define V void
#define cp memcpy
#define cx ((char*)(((x)<<7)>>4))
#define ix ((int*)(((x)<<7)>>4))
#define ir ((int*)(((r)<<7)>>4))
#define _x ((int)x)
#define nx (ix[-1])
#define tx ((int)(x>>57))
#define ty ((int)(y>>57))
#define xi (ix[i])
#define kc(x) ((K)x)
#define ki(x) ((K)2<<57|(K)x)
#define Q(x) if(x){printf("^\n");exit(1)}
#define Qi Q(3!=tx)
K mk(int t,int n){int *r=(int*)malloc(8+(n<<(t&7)));r[0]=1;r[1]=n;return((K)t<<57)|((K)(r+2)>>3);};
V dx(K x){int*r=ix;if(!--r[-2]);free(r-2);}
K kC(int n,char *c){K x=mk(0,n);cp(cx,c,n);return x;};

K til(K x){int n=_x;K r=mk(10,n);int*p=ir;for(int i=0;i<n;i++)p[i]=i;R}
//K take(K x,K y)(Qi;int n=_x;K r=mk(8|ty,n);if(8>ty){for(int i=0;i<n;i++)}
// int u=(K)y;if(!ty)u|=u<<8|u<<16|u<<24;if(3==ty)


int main(int args,char**argv){
 parse("+/alpha+3.14f");
}
