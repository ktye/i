//no std includes(cosmopolitan.h/a provides everything)
#include"k.h"     //c-api
#include"ktye.h"  //k-implementation

#define E(c,s) if(c){printf("^%s\n",s);exit(1);}
#define P(s)   printf("%s\n",s);fflush(stdout);


void serve(int port){
 int s=socket(AF_INET,SOCK_STREAM,6);E((s<0),"socket")
 int opt;E(setsockopt(s,SOL_SOCKET,SO_REUSEADDR,&opt,4)<0,"sockopt")
 struct sockaddr_in in={AF_INET, htons(port),{htonl(0x7f000001)}};E(bind(s,&in,sizeof(in)),"bind")
 E(listen(s,64),"listen")
 
 printf("browse to http://127.0.0.1:%d\n",port);
 if((port==8088)&&IsWindows())system("explorer http://127.0.0.1:8088");
 for(;;){
 
  //todo select s|stdin
 
  struct sockaddr_in c; uint32_t cn=sizeof(c); int f=accept(s,&c,&cn);E(f<0,"accept")
  char q[2040];int n=read(f,q,sizeof(q)-1);E(n<0,"read");
  n=(n>sizeof(q)-1)?sizeof(q)-1:n;q[n]='\0';

  K x=Kx(".",Ks("serve"));E(!x,"serve");
  K r=Kx("@",x,KC(q,n));  E('C'!=TK(r),"type");
  write(f,_M+(int32_t)r,NK(r));
  close(f);
 }
}

void dofile(const char *f){K a=ref(KC(f,strlen(f)));ktye_dofile(a,ktye_readfile(a));}
void ak(){          int f=open("a.k",O_RDONLY);if(0>f)return;close(f);dofile("a.k");}

int main(int args, char **argv){
 int p=8088;
 kinit();
 args_=args;argv_=argv;
 if(1==args)ak();
 if((1<args)){char *a=argv[1];
       if(!strcmp(a,"0")){p=0;             args_--;argv_++; }
  else if(0!=atoi(a)){    p=atoi(argv[1]); args_--;argv_++; }
  else p=0;
 }

 //todo: define/register extension
 char *cwd=getcwd(0,0);KA(Ks("cwd"),KC(cwd,strlen(cwd)));
 ktye_doargs();

 if(p)serve(p);

 ktye_store();
 P("kosmic")
 for(;;){printf(" ");fflush(stdout);ktye_try(ktye_read());}
}
