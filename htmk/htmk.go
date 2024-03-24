package main

import("os";"net/http";"html";"strings";_ "embed")

// get put post patch delete

// get["/table/t/{i:row}/{i:col}"]{x..}
// put[".."]
// post[".."]{[body]..}
// patch delete

func main(){
 kinit();r("get",0);r("put",1);r("post",2);r("delete",3);r("patch",4)
 for _, f:=range os.Args[1:]{b,e:=os.ReadFile(f);fatal(e);dx(Val(KC(b)))}
 http.HandleFunc("/dev.html",dev)
 P:=os.Getenv("P");if P==""{P=":3001"};http.ListenAndServe(P,nil)}
func r(s string,i uint64){Asn(Ks(s),ti(14,int32(l2(i,KC(s)))))}
func Ks(s string)uint64{return sc(KC(s))}
func KC(b []byte)uint64{r:=mk(Ct,int32(len(b)));copy(Bytes[int32(r):],b);return r}
func fatal(e error){if e!=nil{panic(e)}}

func Native(x,y int64)int64{
 S:=[]string{"GET","PUT","POST","DELETE","PATCH"}[x];s,m:=p(CK(x0(y)));f:=r1(y)
 http.HandleFunc(S+" "+s,func(w http.ResponseWriter,r*http.Request){
  for k:=range m{dx(Asn(sK(k),m[k](r.PathValue(k))))}
  b,_:=io.ReadAll(r.Body);w.Write(b(lambda(f,KC(b))))})}
func b(x uint64)[]byte{if tp(x)!=Ct{x=Kst(x)};return CK(x)}

func p(s string)(string,map[string]func(string)uint64){
 //todo: parse "/path/{i:row}/{f:num}"
}

func dev(w http.ResponseWriter, r *http.Request){
 var f []string
 d, _ := os.ReadDir(".");for _,x:=range d{if s:=x.Name();strings.HasSuffix(s,".html")||strings.HasSuffix(s,".css")||strings.HasSuffix(s,".js")||strings.HasSuffix(s,".k"){f=append(f,"<span class='link'>"+html.EscapeString(s)+"</span>")}}
 w.Write([]byte(strings.Replace(devhtml,"FILES",strings.Join(f," "),1)))}

const devhtml=`
<!DOCTYPE html>
<head><meta charset="utf-8"><title></title>
<style>
body{margin:0;overflow:hidden}
.link{}
.link:hover{pointer:cursor}
</style>
</head>
<body style="display:grid;grid-template-columns:auto 1fr">
<div style="height:100vh;display:flex;flex-flow:column;font-family:monospace">
<div>FILES<button>write</button></div>
<textarea style="flex-grow:1;resize:horizontal">
file content
</textarea>
</div>
<iframe style="width:100%;height:100%" src="/"></iframe></body></html>
`
