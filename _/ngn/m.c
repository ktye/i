// ngn/k, (c) 2019-2020 ngn, GNU AGPLv3 - https://bitbucket.org/ngn/k/raw/master/LICENSE
#include"a.h"

//ktye
void  *mmap(void *, size_t, int, int, int, off_t);
int    munmap(void *, size_t);
int    write(int x, char *s, int n) { return 0; }
#define PROT_READ 0
#define PROT_WRITE 0
#define MAP_PRIVATE 0 
#define MAP_ANON 0


#define nxt(x) ((A*)dat(x))[-2]
L mu;S A mx[40];A glb,cn[tn],ci[3][5];dbg(S I ml;/*prevent allocations*/)S NI A oom()_(write(1,"oom\n",4);exit(1);0)
S A ma(UI b)_(dbg(asrt(!ml));asrt(4<b);P(b>=ZZ(mx),oom())mu+=1ll<<b;A x=mx[b];P(x,mx[b]=nxt(x);AB(b,x))I i=b+1;W(i<ZZ(mx)&&!mx[i],i++)
 Y(i<ZZ(mx),x=mx[i];mx[i]=nxt(x))
 E(i=max(b,24);V*p=mmap(0,1ll<<i,PROT_READ|PROT_WRITE,MAP_PRIVATE|MAP_ANON,-1,0);P(((L)p>>4)==-1,oom())x=(A)(p+ZA))
 W(i>b,i--;A y=mx[i]=(A)(x+(1ll<<i));nxt(y)=0)AB(b,x))
A1(mr0,dbg(ml++);asrt(x);P(pkd(x),0)asrt(Ar(x)>0);P(--Ar(x),0)C b=xb;P(!b,AT(tn,x))nxt(x)=mx[b];mx[b]=x&~(-1ull<<48);mu-=1ll<<b;Y(ref(x),mrn(xn,xa))x)
A1(mr1,dbg(ml--);P(!x,0)P(xt==tn,munmap(xc-ZP,xn+ZP);0)dbg(ms(xc,0xab,xn*ZT[xt]));0)
A1(mr,mr1(mr0(x)))V mrn(I n,O A*a){F(n,mr(a[i]))}A1(mRa,F(xn,mR(xai))x)A1(mut,P(pkd(x)||Ar(x)==1,x)A u=atnv(xt,xn,xc);xr;Y(ref(u),mRa(u))u)
V*memcpy(V*x,O V*y,L n)_(C*p=x;O C*q=y;F(n,p[i]=q[i])x)
V*memset(V*x,I c,L n)_(C*p=x;F(n,*p++=c);x)
I strcmp(O C*s,O C*t)_(W(*s&&*s==*t,s++;t++)*s-*t)
L strlen(O C*x)_(O C*p=x;W(*p,p++)p-x)
A atn(C t,L n)_(A x=ma(64-__builtin_clzll(ZA+7+n*ZT[t]));Ar(x)=1;AT(t,AN(n,x)))
A atnv(C t,L n,O V*v)_(A u=atn(t,n);mc(uc,v,n*ZT[t]);u)A1(ax,Ar(x)>1?Ar(x)--,atn(xt,xn):x)A0(aa0,am(enl(as(0)),enl(au0)))
A1(a1,atnv(tA,1,A(x)))A2(a2,atnv(tA,2,A(x,y)))A3(a3,atnv(tA,3,A(x,y,z)))A2(aM,atnv(tM,2,A(x,y)))A2(am,atnv(tm,2,A(x,y)))
A aCn(O C*s,L n)_(atnv(tC,n,s))A aCm(O C*s,O C*t)_(aCn(s,t-s))A aCz(O C*s)_(aCn(s,mn(s)))
A gkv(A*p)_(A x=*p;asrt(xn==2);P(Ar(x)>1,--Ar(x);*p=mR(xx);mR(xy))*p=xx;A u=xy;mr(AN(0,x));u) //get keys and values
A room(A x,L l)_(asrt(xtC);L n=xn;P(Ar(x)==1&&ZA+n+l<=1ll<<xb,x)A u=AN(n,atn(tC,n+l));mc(uc,xc,n);xr;u)
S C*sy1;S C sy0[1<<16]__attribute__((aligned(ZA)));C*syp(I i)_(sy0+i)
I syP(O C*s)_(C c=*s;P(c<128&&(!c||!s[1]),syC(c))O C*p=sy0+256;W(p<sy1,O C*q=p,*t=s;W(*q&&*q==*t,q++;t++)P(!*q&&!*t,p-sy0)p=q;W(*p,p++)p++)
 I r=sy1-sy0,n=mn(s)+1;P(sy1+n>sy0+Z sy0,write(1,"syms ",5);oom())mc(sy1,s,n);sy1+=n;r)
I line(C*p,C*q)_(A x=val(aCm(p,q));P(x,mr(out(x));1)epr(0))S C*skp(C*p)_(W(*p=='/'&&p[1]==10,p+=3;W(*p&&(p[-1]-10||p[-2]-'\\'||p[-3]-10),p++))p)
A1(ldf,x=N(u1c(x));ed(!xn||xc[xn-1]-10,x)xc[xn-1]=0;C*p=xc;Y(*p=='#'&&p[1]=='!',p+=2;W(*p&&*p-10,p++))
 W(*p,C*q=p=skp(p);W(*q&&(*q-10||q[1]==32||q[1]=='}'),q++)Nx(line(p,q));p=q+!!*q)xr;au0)
V init(){tilh((V*)sy0,128);sy1=sy0+256;glb=aa0();
 cn[tA]=a0();cn[tC]=cn[tc]=ac(32);cn[tL]=cn[tl]=al(_0Nl);cn[tI]=cn[ti]=ai(_0Ni);cn[tD]=cn[td]=ad(_0n);cn[tS]=cn[ts]=as(0);F(tn-to,cn[to+i]=au0)
 ci[0][0]=ai(0);ci[0][1]=ai(1);ci[0][2]=ai(_0Wi);ci[0][3]=ai(-_0Wi);ci[0][4]=ai(_0Ni);
 S L l[]={0,1,_0Wl,-_0Wl,_0Nl};F(5,ci[1][i]=al(l[i]))S D d[]={0,1,_0w,-_0w,_0n};F(5,ci[2][i]=ad(d[i]))}


//ktye (cut)
