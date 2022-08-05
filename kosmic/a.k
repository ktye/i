port:8088

serve:{x:"\n"\:x@&~x="\r"; \q:*x;  
 $[~h~(-#h:" HTTP/1.1")#q; r404
   ~get q;                 r404
   1~#q:4_(-#h)_q;         html index
   "/?"~2#q;               r404          /todo k
   r404]}
  
/todo: readfile without trap
/todo: readfile + mime types

prefix:{x~(#x)#y}
get:prefix"GET "

http:{"HTTP/1.1 200 OK\nContent-type: ",x,"\nContent-Length: ",($#y),"\n\n",y}
html:http"text/html; charset=UTF-8"
r404:"HTTP/1.1 404 Not Found\nContent-Length: 0\n\n"

index:"<!DOCTYPE html>
<head><meta charset=\"utf-8\"><title>kosmic</title></head>
<body>a b c</body></html>\n"
