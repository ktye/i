#include<winsock2.h>
#include<stdio.h>
int main(int args, char **argv){

 /*win*/WSADATA wsa;int w=WSAStartup(18, &wsa);if(w){printf("^wsa %d\n",w);exit(1);}

 int f=socket(AF_INET,SOCK_STREAM,6);
 if(f<0){printf("^socket\n");exit(1);}

 //int opt;if(setsockopt(f,SOL_SOCKET,SO_REUSEADDR,&opt,4)<0){printf("^sockopt\n");exit(1);}
 /*win*/if(setsockopt(f,SOL_SOCKET,SO_REUSEADDR,"abcd",4)<0){printf("^sockopt\n");exit(1);}

 struct sockaddr_in in = {AF_INET, htons(8088), {htonl(0x7f000001)}};
 //if(bind(f,&in,sizeof(in))){printf("^bind\n");exit(1);}
 /*win*/if(bind(f,(SOCKADDR*)&in,sizeof(in))){printf("^bind: %s\n",WSAGetLastError());exit(1);}
}
