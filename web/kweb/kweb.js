import { K } from '../k.js'
import {plot,draw} from './plot.js'

function ge(x){return document.getElementById(x)}
function ce(x){return document.createElement(x)}
function su(u){return (u.length)?new TextDecoder("utf-8").decode(u):""}
function us(s){return new TextEncoder("utf-8").encode(s)}
function pd(e){if(e){e.preventDefault();e.stopPropagation()}};
function rm(p){while(p.firstChild)p.removeChild(p.firstChild)}


//predefined ui types and defaults per type.
let UI={
 table:      uitable,
 listbox:    uilistbox,
 select:     uiselect,
 tree:       uitree,
 input:      uiinput,
 button:     uibutton,
 edit:       uiedit,
 h1:         uih1,
 tty:        uitty,
 text:       uitext,
 c:uitext,   i:uiinput,  s:uiinput, f:uiinput,  z:uiinput,
 C:uitext,   I:uilistbox,S:uiselect,F:uilistbox,Z:uilistbox,
 L:uilistbox,D:uitable,  T:uitable, m:uitext,p:uitext,l:uitext
}
UI["0"]=uitext

function register(s,f){UI[s]=f} //register user supplied elements: kweb.register(str,function(dst,x))



let nodes={}

function update(){
 for(let id in nodes){
  let dst=nodes[id]
  if(dst.offsetParent===null)continue     //skip invisible
  if(dst.s){
   let k=K.Kx(".",dst.s)
   if(k!=dst.k)gui(id,null,dst.s,null,K.ref(dst.d))    //global variable has changed
   else        K.unref(k)
  }
  if(dst.e){
   gui(id,null,null,K.ref(dst.e),K.ref(dst.d)) //evaluate expr
  }
}}


function show(d,x){
 //convert short form to dict.
 let id="",t=K.TK(d)
 if("s"==t){ id=K.sK(d); d=K.Kx("!",K.Ks(""),K.Ks("")) }
 if("S"==t){ id=K.sK(K.Kx("*",K.ref(d))); d=K.Kx("`type!*|",d) }
 if("D"==t){ let kid=K.Kx("@",K.ref(d),K.Ks("id"))
  if("s"==K.TK(kid))id=K.sK(kid)
  else K.unref(kid)
 }
 
 //hide: show[`id;`]
 if(x==K.Ks("")){ge(id).classList.remove("hidden");K.unref(d);return K.Ks(id)}
 
 switch(K.TK(x)){
 case "l": return gui(id,null,null,x,d)           //expr(lambda)
 case "s": return gui(id,null,x,null,d)           //symbol
 default:  return gui(id,x,null,null,d)           //value
}}


function gui(id,x,s,e,d){id=(id=="")?"uid"+String(Object.keys(nodes).length):id
 let dst=ge(id);
 if(dst==null){dst=ce("div");dst.id=id;document.body.appendChild(dst)}
 dst.classList.remove("hidden")
 dst.classList.add(...classes(d))
 if(('d' in dst)&&dst.d)K.unref(dst.d); dst.d=d  //dict
 if(('e' in dst)&&dst.e)K.unref(dst.e); dst.e=e  //expr
 if(('k' in dst)&&dst.k)K.unref(dst.k); dst.k=x  //k-value
 dst.s=s                                         //symbol
 rm(dst)
 nodes[id]=dst
 if(s){dst.k=K.Kx(".",s)}                        //symbol
 if(e){dst.k=K.Kx(".",K.ref(e),K.KL([]))}        //expr
 let t=K.TK(dst.k)
 
 let u=K.JK(K.Kx("@",K.ref(d),K.Ks("type")))
 let f=UI[(("string"!=typeof u)||u=="")?t:u]
 f(dst,K.ref(dst.k))
 
 let a=K.JK(K.Kx("`id`type_",K.ref(d)))
 let keys=Object.keys(a),cld=dst.firstChild;
 if(e)cld.readOnly=true
 for(let i=0;i<keys.length;i++){let ki=keys[i];let ai=a[ki]
  cld[ki]=(ki.startsWith("on"))?jsevent(ki,dst):ai
 }
 return K.Ks(id)
}
function classes(d){
 let c=K.Kx("@",K.ref(d),K.Ks("class"))
 switch(K.TK(c)){
 case "S": return K.SK(c)
 case "s": let r=K.sK(c);return(r=="")?[]:[K.sK(c)]
 default:  K.unref(c);return []
}}
function jsevent(s,dst){
 return function(e){
 let a=[];let f=K.Kx("@",K.ref(dst.d),K.Ks(s))
  switch(s){
  case "onchange":  a=[K.KJ(e.target.value)]; break
  case "onkeydown": a=[K.Ks(e.target.key)];   break
  default:
  }
  K.unref(K.Kx(".",f,K.KL(a)))
  update()
 }
}

function uitext(dst,x){ //text node
 let s=ce("span");s.textContent=(K.TK(x)=="C")?K.CK(x):K.CK(K.Kx("`k@",x));dst.appendChild(s)
}

function uiinput(dst,x){ //input element from Cisfz
 let t=K.TK(x)
 let s=K.CK(("C"==t)?x:K.Kx("$",x))
 let e=ce("input");e.type="text";e.value=s;e.defaultValue=s;e.readOnly=dst.classList.contains("readonly")
 e.onchange=function(evt){
  let x=K.Kx("$",K.Ks(t),K.KC(evt.target.value))
  if(x==0)  e.value=e.defaultValue
  else{     e.value=K.CK(("C"==t)?K.ref(x):K.Kx("$",K.ref(x)));e.defaultValue=e.value
   if(dst.s)K.KA(dst.s,x)
   else     K.unref(x)
   update()
 }}
 dst.appendChild(e)
}
function uibutton(dst,x){
 let b=ce("input");b.type="button";b.classList.add("kweb-button")
 b.value=K.CK(x)
 dst.appendChild(b)
}
function uicheckbox(dst,x){ //b
 uiinput(dst,x)
 dst.firstChild.type=checkbox
}
function uiedit(dst,x){
 let ta=ce("textarea");ta.classList.add("kweb-textarea");ta.readOnly=dst.classList.contains("readonly")
 ta.value=K.CK(x)
 ta.onchange=function(evt){if(dst.s)K.KA(dst.s,K.KC(ta.value));update()}
 dst.appendChild(ta)
}
function uih1(dst,x){
 let h=ce("h1");h.classList.add("kweb-h1");h.textContent=K.CK(x);dst.appendChild(h)
}


let O=function(x){console.log("out k>",x)}                           //default k output

function uitty(dst,x){
 let tty=ce("textarea");tty.value=K.CK(x)
 O=function(x){tty.value+=x}      //redirect k output to tty
 dst.evl=function(s){let c=s.startsWith("\n  ")?" ":"";s=s.trim();if(!s.length)return
  tty.value+="\n";krep(c+s);tty.value+=" ";tty.scrollTop=tty.scrollHeight
  update()
 }
 tty.onkeydown=function(e){
  if(("Enter"==e.key)&&(0<tty.value.length)){pd(e);
   let v=tty.value; let i=v.lastIndexOf("\n");
   dst.evl((i<0)?v:v.slice(i))
 }}
 dst.appendChild(tty)
}

function html(s){ return s.replace(/[\u00A0-\u9999<>\&]/g,((i)=>`&#${i.charCodeAt(0)};`)) }
function uiselect(dst,x){ //create select element from vectors
 let n=K.NK(x)
 let S=Array(n)
 if(K.TK(x)=="L"){
  for(let i=0;i<n;i++){let xi=K.Kx("@",K.ref(x),K.Ki(i))
   S[i]=K.CK(("C"==K.TK(xi))?xi:K.Kx("`k@",xi))
  }
  K.unref(x)
 }else S=K.LK(K.Kx("$",x)).map(K.CK)
 let s=ce("select")
 for(let i=0;i<n;i++){
  let o=ce("option")
  o.value=S[i]
  o.innerHTML=html(S[i]).replace(/%20/g,"&nbsp;")
  o.classList.add("kweb-option")
  s.appendChild(o)
 }
 s.selectedIndex=-1
 dst.appendChild(s)
 dst.classList.add("kweb-select")
}
function uilistbox(dst,x){ //listbox from vectors, or T D
 if(-1<"TD".indexOf(K.TK(x)))x=K.Kx("`l@",x)
 uiselect(dst,x)
 let lb=dst.firstChild;lb.multiple=true
}
function uidicttab(dst,x){ //D (only S!..)
 let [k,v]=K.LK(x)
 if(K.TK(k)!="S"){K.unref(k);K.unref(v);return}
 let S=K.SK(k), n=S.length
 let ta=ce("table")
 ta.classList.add("kweb-dict")
 for(let i=0;i<n;i++){
  let tr=ce("tr")
   let th=ce("th")
    th.textContent=S[i]
   tr.appendChild(th)
   let td=ce("td")
    td.textContent=K.CK(K.Kx("`k@",K.Kx("@",K.ref(v),K.Ki(i))))
   tr.appendChild(td)
  ta.appendChild(tr)
 }
 dst.appendChild(ta)
 K.unref(v)
 //todo: editable
}

function uitable(dst,x){ //TD
 if(K.TK(x)=="D")return uidicttab(dst,x)
 let N=K.iK(K.Kx("#",K.ref(x)))
 let L=K.LK(x)                    //[keys,values]
 let S=K.SK(L[0])
 let ta=ce("table") 
 ta.classList.add("kweb-table")
 ta.cols=S                        //store column names/types
 ta.coltype=Array(S.length)
  let tr=ce("tr")
   for(let i=0;i<S.length;i++){let th=ce("th");th.textContent=S[i];tr.appendChild(th)}
  ta.appendChild(tr)
  let V=[]
  for(let i=0;i<S.length;i++){
   let v=K.Kx("@",K.ref(L[1]),K.Ki(i))
   ta.coltype[i]=K.TK(v)
   V[i]=SV(v)
  }
  K.unref(L[1])
  for(let i=0;i<N;i++){
   let tr=ce("tr")
   tr.i=i                         //store row index
   for(let j=0;j<S.length;j++){
    let td=ce("td")
     td.textContent=V[j][i]
    tr.appendChild(td)
   }
   ta.appendChild(tr)
  }
 dst.appendChild(ta)
 if(dst.s&&(!dst.classList.contains("readonly")))editTable(dst)
}
function editTable(dst){ //make table editable (link with k)
 dst.querySelectorAll("td").forEach(x=>{
  x.contentEditable=true
  x.onbeforeinput=function(e){let et=e.target
   if(!('old' in et))et.old=et.textContent
  }
  x.onkeydown=function(e){let et=e.target
   if(e.key=="Enter"){
    let i=et.parentElement.i,j=et.cellIndex
    let t=et.parentElement.parentElement
    let ct=t.coltype[j]
    let s=K.Ks(ct.toLowerCase())
    let v=K.Kx("$",s,K.KC((ct=="S")?'`"'+et.textContent+'"':et.textContent))
    if(v==0){x.classList.add("kweb-invalid");return false}
    et.textContent=K.CK(K.Kx("$",K.ref(v)))
    et.classList.remove("kweb-editing","kweb-invalid")
    et.blur()
    v=K.Kx(".",dst.k,K.Kx(",",K.Ks(t.cols[j]),K.Ki(i)),v)
    K.KA(dst.s,v)
    dst.k=K.ref(v)
    delete et.old
    update()
    return false
   }else if(e.key=="Escape"){
    if('old' in et)et.textContent=et.old
    delete et.old
    et.blur()
    et.classList.remove("kweb-editing","kweb-invalid")
    return false
   }
   et.classList.add("kweb-editing")
  }
})}
function SV(x){ //strings from vector 
 let l
 switch(K.TK(x)){
 case 'B':case 'I':case 'S':case 'F':case 'Z':l=K.Kx("$",x);break
 default:{console.error("nyi");return[]}
 }
 let r=[]
 let n=K.NK(x)
 for(let i=0;i<n;i++)r.push(K.CK(K.Kx("@",K.ref(l),K.Ki(i))))
 K.unref(l);return r
}

function uitree(dst,x){ //treeview for D
 console.log("nyi treeview")
}


//initialzation:
// - get & init k.wasm
// - get & exec a.k
// - exec user supplied start(js or k script)
// - initialize divs (connect to k variables)
// - run post function(js) if present
function initKweb(start,post,ak){
 return function(){
  start=("string"==typeof start)?[start]:start
  let dostart=function(){
   if("object"==typeof start)Promise.all(start.map(x=>fetch(x))).then(r=>Promise.all(r.map(r=>r.text()))).then(x=>{for(let i=0;i<x.length;i++)ktry(x[i]);initDivs(post)})
   else{                            start();                                  initDivs(post)}
  } 
  if(ak!==false){
   fetch("a.k").then(r=>r.text()).then(r=>{
    K._.Val(K.KC(r))
           dostart()
  })} else dostart()
}}
function initDivs(post){
 document.querySelectorAll("div").forEach(x=>{
  let d=function(){let k="",v="";if('kType'in x.dataset){k="type",v=x.dataset.kType}
   return K.Kx("!",K.Ks(k),K.Ks(v))
  }
       if('kVar' in x.dataset){cl(x);gui(x.id,null,K.Ks(x.dataset.kVar),null,d())}
  else if('kExpr'in x.dataset){cl(x);gui(x.id,null,null,K.Kx(x.dataset.kExpr),d())}
  else if('kVal' in x.dataset){cl(x);gui(x.id,K.Kx(x.dataset.kVal),null,null,d())}
 })
 let a=ce("a");a.id="_dl";a.style.display="none" //hidden download element
 document.body.appendChild(a)
 if("function"==typeof post)post()
}

function cl(x){if('d'in x)delete x.d;if('e'in x)delete x.e;if('k'in x)delete x.k} //for reinit

document.update=update //for custom updates from k

// drop files (execute .k)
let _dropfile=dropfile
window.ondragover=function(e){pd(e)}
window.ondrop=function(e){pd(e);if(e.dataTransfer.items){for(let i=0;i<e.dataTransfer.items.length;i++){if(e.dataTransfer.items[i].kind=='file'){let file=e.dataTransfer.items[i].getAsFile();addfile(file)}}}else for(let i=0;i<e.dataTransfer.files.length;i++)addfile(e.dataTransfer.files[i])}
function addfile(x){
 let r=new FileReader()
 r.onload=function(){
  let u=new Uint8Array(r.result)
  _dropfile(x.name,u)
 }
 r.readAsArrayBuffer(x)
}
function filedrop(f){_dropfile=f}
function dropfile(name,u){ //default drop handler can be overwritten with kweb.filedrop(f)
 if(name.endsWith(".k")){ document.body.innerHTML=""; ktry(su(u)) }
}
function writefile(name,u){let l=String(window.location)
 if(l.startsWith("http://localhost:")||l.startsWith("http://127.0.0.1:")) upload(name,u)
 else                                                                   download(name,u)
}
function download(name,u){
 let dl=ge("_dl");let b=new Blob([u],{type:"application/octet-stream"})
 dl.href=URL.createObjectURL(b);dl.download=name;dl.click()}
function upload(name,u){let l=window.location
 let h="http://"+l.hostname+":"+l.port+"/"+name
 fetch(h, {method:"POST",body:u/*new Blob(u)*/})
}
let files={}
function fsadd(name,u){files[name]=u}
function readfile(name){ // k read, e.g. x:<`name
 if(name in files)return files[name]
 
 let h = new XMLHttpRequest()
 h.open("GET", name, false)
 //h.responseType="arraybuffer" //cannot set type for synchronous request. text-only.
 h.send(null)
 return us((200===h.status)?h.responseText:"")
}


function init(o){
 nodes={}
 if(o.plot!==false) register('plot',plot)
 if(o.draw!==false) register('draw',draw)
 let start = ("start"in o)?o.start : function(){}
 let post  = ("post"in o) ?o.post:   function(){}
 let ext={                  //wasm import module
  init: initKweb(start,post,("ak"in o)?o.ak:true),
  read: readfile,
  write:function(file,data){if(file===""){O(su(data))}else{writefile(file,data)}},
  show: show,
  hide: function(id)       {ge(K.sK(K.ref(id))).classList.add("hidden");return id},
  js:   K.JS,
 }
 if("ext"in o)Object.assign(ext,o.ext)
 K.kinit(ext,("wasm"in o)?o.wasm:"../k.wasm")
}

function ktry(s){
 try     {let x=K._.Val(K.KC(s));K.save();return x}
 catch(e){console.log(e);K.restore();return false}
}
function krep(s){
 if((s=="\\")||s=="\\h"){help();return}
 try     {let x=K._.repl(K.KC(s));K.save();return x}
 catch(e){console.log(e);K.restore();return false}
}
function help(){fetch('readme').then(r=>r.text()).then(r=>O("\n"+r+" "))}


let kweb = {init,ktry,show,update,register,filedrop,fsadd}
export { kweb }
