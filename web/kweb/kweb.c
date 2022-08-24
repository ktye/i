/* kweb local application server
 * 
 * it serves static files attached to the binary (kweb.exe)
 * other files are served from disk
 * post requests write files to disk
 *
 * see mk for building
 */

#include<windows.h>
#include<winsock2.h>
#include<stdio.h>
#include<stdlib.h>

   
#define E(c,s) if(c){printf("^%s\n",s);exit(1);}

#define OK "HTTP/1.1 200 OK\nContent-Type: "
#define HM "text/html\n"
#define JS "text/javascript\n"
#define CS "text/css\n"
#define TX "text/plain\n"
#define LN "Content-Length: "
#define GT "GET /"
#define PS "POST /"
#define HT " HTTP/1.1"
#define R4 "HTTP/1.1 404 Not Found\nContent-Length: 0\n\n"
#define NX closesocket(f);continue
#define writec(x,y)  send((x),(y),sizeof(y)-1,0)
#define write(x,y,z) send((x),(y),(z),0)
#define WSA WSADATA wsa;WSAStartup(MAKEWORD(2,0),&wsa)


// files are attached to the binary
// each file has a header line: \filename
int  NF;     //attached files
char*F[32];  //name
char*P[32];  //data
int  N[32];  //len

int flen(FILE *f){
 fseek(f,0,SEEK_END);
 int r=ftell(f);
 fseek(f,0,SEEK_SET);
 return r;
}

void fsys(char *a0){
 FILE *f=fopen(a0,"rb");
 E(!f,"^a0");
 int n=flen(f);
 char *a=malloc(n);
 E(n!=fread(a,1,n,f),"<a0");
 fclose(f);

 NF=0;
 char s[10];
 s[0]='\n';s[1]='\\';
 memcpy(2+s,"k.wasm\n",7);
 
 for(int i=0;i<n-1;i++){
  if((a[i]!='\n')||a[1+i]!='\\')continue;
  char *c=2+i+a;
  if(!NF){if(strncmp(c,"k.wasm\n",10))continue;  //k.wasm is first attachment.
   printf("kweb.exe:%d\n", i);
  }
  c=strchr(c,10);
  if(!c)continue;
  *c++='\0';
  F[NF]=2+i+a;
  P[NF]=c;
  N[NF]=n-(c-a);
  if(0<NF)N[NF-1]-=n-i;
  NF++;
 }
 for(int i=0;i<NF;i++){
  printf("%8s:%d\n",F[i],N[i]);
 }
}

void sendfile(SOCKET d,FILE *f,void *ignore,int n){
 char b[1024];
 for(;;){
  int n=fread(b,1,1024,f);
  if(!n)return;
  write(d,b,n);
}}

void hdr(SOCKET d, char *c, int n){
 writec(d,OK);
 if(strstr(c,".html"))    writec(d,HM);
 else if(strstr(c,".js")) writec(d,JS);
 else if(strstr(c,".css"))writec(d,CS);
 else                     writec(d,TX);
 writec(d,LN);
 char sz[32];
 itoa(n,sz,10);
 write(d,sz,strlen(sz));
 write(d,"\n\n",2);
}

int getfile(SOCKET d,char *c){
 printf("get: %s ",c);
 if(!c[0])memcpy(c,"a.html",7);
 for(int i=0;i<NF;i++){  //{from mem}
  if(!strcmp(F[i],c)){
   hdr(d,c,N[i]);
   write(d,P[i],N[i]);
   printf("{%d ok}\n",N[i]);
   return 1;
 }}
 FILE *f=fopen(c,"rb");  //[from disk]
 if(!f)return 0;
 int n=flen(f);
 hdr(d,c,n);
 sendfile(d,f,NULL,n);
 printf("[%d ok]\n",n);
 fclose(f);
 return 1;
}


int postfile(SOCKET s,char *name, char *c,int n,int m){int t=n; //n: still in c, m:content-length
 printf("post:%s [%d ", name, m);
 FILE *f=fopen(name,"wb");
 if(!f)return 0;
 fwrite(c,1,n,f);
 char b[1024];
 while(t<m){
  int n=recv(s,b,1023,0);
  t+=n;
  fwrite(b,1,n,f);
  if(n==0)break;
 }
 printf("ok]\n");
 fclose(f);return 1;
}

void serve(int port){
 SOCKET s=socket(AF_INET,SOCK_STREAM,6);E((s<0),"socket")
 int opt;E(setsockopt(s,SOL_SOCKET,SO_REUSEADDR,(char*)&opt, sizeof(opt))<0,"sockopt")
 //struct sockaddr_in in={AF_INET, htons(port),{htonl(0x7f000001)}}; //does not work on windows.
 struct sockaddr_in in;in.sin_family=AF_INET;in.sin_addr.s_addr=inet_addr("127.0.0.1");in.sin_port=htons(port);
 E(bind(s,(SOCKADDR*)&in,sizeof(in)),"bind")
 E(listen(s,64),"listen")
 
 printf("browse to http://127.0.0.1:%d\n",port);
 if(port==8088)system("explorer http://127.0.0.1:8088");
 for(;;){
  struct sockaddr_in c; uint32_t cn=sizeof(c); SOCKET f=accept(s,&c,&cn);E(f<0,"accept")
  char q[2040];int n=recv(f,q,sizeof(q)-1,0);E(n<0,"read");
  n=(n>sizeof(q)-1)?sizeof(q)-1:n;q[n]='\0';
  
  //printf("q>%s\n", q);
  
  if(n>sizeof GT){
   if(!strncmp(GT,q,-1+sizeof GT)){
    char *c=strstr(q,HT);
    if(c==NULL)goto NE;
    *c='\0';
    if(!getfile(f,q+5))goto NE;
    NX;
  }}
  if(n>sizeof PS){
   if(!strncmp(PS,q,-1+sizeof PS)){
    char *c=strstr(q,HT);if(c==NULL)goto NE;
    *c='\0';
    c=strstr(1+c,LN);if(c==NULL)goto NE;
    int m=atoi(c+sizeof(LN)-1);
    c=strstr(1+c,"\r\n\r\n");if(c==NULL)goto NE;
    if(!postfile(f,6+q,4+c,n-(c-q)-4,m))goto NE;
    NX;
  }}
  
NE:
  printf("[404]\n");
  writec(f,R4);
  closesocket(f);
 }
}


int main(int args, char **argv){ fsys(argv[0]);WSA;serve((2==args)?atoi(argv[1]):8088); }
