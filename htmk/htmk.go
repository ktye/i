package main

import("os";"net/http";"html";"strings";_ "embed")

// get put post
//

// get["/table/t/{i:row}/{i:col}"]{x..}
// put[".."]
// post[".."]{[body]..}
// patch delete

func main(){
 kinit();r("get");r("put");r("post");r("delete");r("patch")
 for _, f:=range os.Args[1:]{b,e:=os.ReadFile(f);fatal(e);dx(Val(KC(b)))}
 http.HandleFunc("/dev.html",dev)
 P:=os.Getenv("P");if P==""{P=":3001"};http.ListenAndServe(P,nil)}
func r(m string){}
func KC(b []byte)uint64{r:=mk(Ct,int32(len(b)));copy(Bytes[int32(r):],b);return r}
func fatal(e error){if e!=nil{panic(e)}}


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
