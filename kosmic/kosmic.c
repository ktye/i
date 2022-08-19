
#define E(c,s) if(c){printf("^%s\n",s);exit(1);}

#define OK "HTTP/1.1 200 OK\nContent-Type: "
#define HM "text/html\n"
#define JS "text/javascript\n"
#define CS "text/css\n"
#define TX "text/plain\n"
#define LN "Content-Length: "
#define GT "GET /"
#define HT " HTTP/1.1"
#define R4 "HTTP/1.1 404 Not Found\nContent-Length: 0\n\n"
#define NX close(f);continue
#define writec(x,y) write((x),(y),sizeof(y)-1)


// files are attached to the binary
// each file has a header line: \filename
int  NF;     //attached files
char*F[32];  //name
char*P[32];  //data
int  N[32];  //len

int flen(int f){
 struct stat st;
 if(fstat(f,&st))return -1;
 return st.st_size;
}

void fsys(char *a0){
 int f=open(a0,O_RDONLY);
 E(!f,"^a0");
 int n=flen(f);
 E(0>n,"#a0");
 
 char *a=malloc(n);
 E(n!=read(f,a,n),"<a0");
 close(f);

 NF=0;
 char s[10];
 s[0]='\n';s[1]='\\';
 memcpy(2+s,"k.wasm\n",7);
 
 for(int i=0;i<n-1;i++){
  if((a[i]!='\n')||a[1+i]!='\\')continue;
  char *c=2+i+a;
  if(!NF)if(strncmp(c,"k.wasm\n",10))continue;  //k.wasm is first attachment.
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
  printf("%s:%d\n",F[i],N[i]);
 }
}

void hdr(int d, char *c, int n){
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

int getfile(int d,char *c){
 printf("get: %s ",c);
 if(!c[0])memcpy(c,"a.html",7);
 for(int i=0;i<NF;i++){  //from mem
  if(!strcmp(F[i],c)){
   hdr(d,c,N[i]);
   write(d,P[i],N[i]);
   printf("{%d ok}\n",N[i]);
   return 1;
  }
 }                       //from disk
 int f=open(c,O_RDONLY);
 if(!f)return 0;
 int n=flen(f);
 if(0>n){close(f);return 0;}
 hdr(d,c,n);
 sendfile(d,f,NULL,n);
 printf("[%d ok]\n",n);
 return 1;
}

void serve(int port){
 int s=socket(AF_INET,SOCK_STREAM,6);E((s<0),"socket")
 int opt;E(setsockopt(s,SOL_SOCKET,SO_REUSEADDR,&opt,4)<0,"sockopt")
 struct sockaddr_in in={AF_INET, htons(port),{htonl(0x7f000001)}};E(bind(s,&in,sizeof(in)),"bind")
 E(listen(s,64),"listen")
 
 printf("browse to http://127.0.0.1:%d\n",port);
 //if((port==8088)&&IsWindows())system("explorer http://127.0.0.1:8088");
 for(;;){
  struct sockaddr_in c; uint32_t cn=sizeof(c); int f=accept(s,&c,&cn);E(f<0,"accept")
  char q[2040];int n=read(f,q,sizeof(q)-1);E(n<0,"read");
  n=(n>sizeof(q)-1)?sizeof(q)-1:n;q[n]='\0';
  
  if(n>sizeof GT){
   if(strncmp(GT,q,sizeof GT)){
    char *c=strstr(q,HT);
    if(c==NULL)goto NE;
    *c='\0';
    if(!getfile(f,q+5))goto NE;
    NX;
   }
  }
  
NE:
  printf("[404]\n");
  writec(f,R4);
  close(f);
 }
}

int main(int args, char **argv){ fsys(argv[0]); serve((2==args)?atoi(argv[1]):8088); }
