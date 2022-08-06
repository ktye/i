
serve:{x:"\n"\:x@&~x="\r"; \q:*x;        /q:request(1st line)
 $[~h~(-#h:" HTTP/1.1")#q; r404          /http
   ~get q;                 r404          /GET requests
   1~#q:4_(-#h)_q;         html index    /root
   "/?"~2#q;               text `k@. 2_q /k-expr
   0~#f:<`$cwd,q;          r404          /file not found
   http[mime ext q;f]]}                  /send file
  
ext:{`$(x?".")_x}
prefix:{x~(#x)#y}
get:prefix"GET "

http:{"HTTP/1.1 200 OK\nContent-type: ",x,"\nContent-Length: ",($#y),"\n\n",y}
mime:``html`js!("text/plain";"text/html";"text/javascript")
html:http mime`html
text:http mime`
r404:"HTTP/1.1 404 Not Found\nContent-Length: 0\n\n"

index:"<!DOCTYPE html>
<head><meta charset=\"utf-8\"><title>kosmic</title></head>
<body>a b c</body></html>\n"
