import { K } from '../k.js'

let T=[
 //k-type    class      ui element
 ["T",       "",        uitable],      //table
 ["BISFZL",  "listbox", uiselect],     //
 ["BISFZL",  "",        uiselect],     //select(dropdown)
 ["D",       "tree",    uitree],       //tree view
 ["D",       "",        uidicttab],    //table on row per key
 ["Cisfz",   "input",   uiinput],      //input element, veryfied by k
 ["C",       "edit",    uitextarea],   //text area
 ["C",       "h1",      uih1],         //header
 ["C",       "tty",     uitty],        //k-console
 ["C",       "",        uispan],       //span element
]

function ge(x){return document.getElementById(x)}
function ce(x){return document.createElement(x)}
function su(u){return (u.length)?new TextDecoder("utf-8").decode(u):""}
function us(s){return new TextEncoder("utf-8").encode(s)}
function pd(e){if(e){e.preventDefault();e.stopPropagation()}};
function rm(p){while(p.firstChild)p.removeChild(p.firstChild)}

let nodes={}

function update(){
 for(let id in nodes){
  let dst=nodes[id]
  if(dst.offsetParent===null)continue     //skip invisible
  if(dst.s){ 
   let k=K.Kx(".",dst.s)
   if(k!=dst.k)gui(id,null,dst.s,null)    //global variable has changed
  }
  if(dst.e)gui(id,null,null,K.ref(dst.e)) //evaluate expr
}}


function show(kid,x){let id=K.sK(K.Kx("*",K.ref(kid)))
 let c=("S"==K.TK(kid))?K.SK(K.Kx("_",K.Ki(1),kid)):[]  //classes
 switch(K.TK(x)){
 case "l": return gui(id,null,null,x,c)           //expr(lambda)
 case "s": return gui(id,null,x,null,c)           //symbol
 default:  return gui(id,x,null,null,c)           //value
}}

function gui(id,x,s,e,c){id=(id=="")?"uid"+String(Object.keys(nodes).length):id
 let dst=ge(id);
 if(dst==null){dst=ce("div");dst.id=id;document.body.appendChild(dst)}
 if(c)dst.classList.add(...c)
 if(('e' in dst)&&dst.e)K.unref(e);     dst.e=e  //expr
 if(('k' in dst)&&dst.k)K.unref(dst.k); dst.k=x  //k-value
 dst.s=s                                         //symbol
 rm(dst)
 nodes[id]=dst
 if(s){dst.k=K.Kx(".",s)}                        //symbol
 if(e){dst.k=K.Kx(".",K.ref(e),K.KL([]))}        //expr
 let t=K.TK(dst.k)
 console.log("gui type ", t)
 let f=function(x,y){console.log("gui:type",t)}
 for(let i=0;i<T.length;i++){let t0=T[i][0],t1=T[i][1],t2=T[i][2]
  if(-1<t0.indexOf(t)){
   if((t1!="")&&dst.classList.contains(t1)){f=t2;break}
   //else{                                    f=t2;break}
 }}
 f(dst,K.ref(dst.k))
 return K.Ks(id)
}

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

function uispan(dst,x){  //text node for C
 let s=ce("span");s.textContent=K.CK(x);dst.appendChild(s)
}

function uiinput(dst,x){ //input element from Cisfz
 let t=K.TK(x)
 let s=("C"==t)?K.CK(x):K.Kx("$",x)
 let e=ce("input");e.type="text";e.value=s;e.defaultValue=s;e.readOnly=dst.classList.contains("readonly")
 e.onchange=function(evt){
  let x=K.Kx("$",K.Ks(t),K.KC(evt.target.value))
  if(x==0)  e.value=e.defaultValue
  else{     e.value=K.CK(K.ref(x));e.defaultValue=e.value
   if(dst.s)K.KA(dst.s,x)
   else     K.unref(x)
 }}
 dst.appendChild(e)
}
function uitextarea(dst,x){
 let ta=ce("textarea");ta.classList.append("kweb-textarea");ta.readOnly=dst.classList.contains("readonly")
 ta.value=K.CK(x)
 ta.onchange=function(evt){if(dst.s)K.KA(dst.s,K.KC(ta.value))}
 dst.appendChild(ta)
}
function uih1(dst,x){
 let h=ce("h1");h.classList.append("kweb-h1");h.textContent=K.CK(x);h.appendChild(h)
}

let O=function(x){console.log("k>",x)}                           //default k output

function uitty(dst,x){
 let tty=ce("textarea");tty.value=K.CK(x)
 O=function(x){tty.value+=x;tty.scrollTop=tty.scrollHeight}      //redirect k output to tty
 tty.onkeydown=function(e){
  if(("Enter"==e.key)&&(0<tty.value.length)){pd(e);
   let v=tty.value; let i=v.lastIndexOf("\n");
   let s=((i<0)?v:v.slice(i)).trim()
   if(!s.length)return
   evl(s); tty.scrollTop=tty.scrollHeight;
 }}
 dst.appendChild(dst)
}

function uiselect(dst,x){ //create select element from vectors
 let n=K.NK(x)
 let S=Array(n)
 if(K.TK(x)=="L"){
  for(let i=0;i<n;i++)S[i]=K.Kx("`k@",K.Kx("@",K.ref(x),K.Ki(i)))
  K.unref(x)
 }else S=K.LK(K.Kx("$",x)).map(K.CK)
 let s=ce("select")
 for(let i=0;i<n;i++){
  let o=ce("option")
   o.textContent=S[i]
  s.appendChild(o)
 }
 dst.appendChild(s)
}
function uidicttab(dst,x){ //create table from x(dict), only for symbol keys
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
    td.textContent=K.Kx("`k@",K.Kx("@",K.ref(v),K.Ki(i)))
   tr.appendChild(td)
  ta.appendChild(tr)
 }
 dst.appendChild(ta)
 K.unref(v)
 //todo: editable
}

function uitable(dst,x){ //create table from x(table)
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

function uitree(dst,x){ //treeview for D
 console.log("nyi treeview")
}


/*
function style(x){
 let s=ce("style");s.innerText=K.CK(x);
 document.head.appendChild(s);return BigInt(0)
}
*/

document.update=update //for custom updates from k

// drop files (execute .k)
window.ondragover=function(e){pd(e)}
window.ondrop=function(e){pd(e);if(e.dataTransfer.items){for(let i=0;i<e.dataTransfer.items.length;i++){if(e.dataTransfer.items[i].kind=='file'){let file=e.dataTransfer.items[i].getAsFile();addfile(file)}}}else for(let i=0;i<e.dataTransfer.files.length;i++)addfile(e.dataTransfer.files[i])}
function addfile(x){
 let r=new FileReader()
 r.onload=function(){
  let u=new Uint8Array(r.result)
  if(x.name.endsWith(".k")){ document.body.innerHTML=""; ktry(su(u)) }
 }
 r.readAsArrayBuffer(x)
}

function init(start,kwasm){ //start k
 kwasm=(kwasm!==undefined)?kwasm:"../k.wasm"
 let ext={                  //wasm import module
  init: start,
  read: function(file)     {return new Uint8Array(0)},
  write:function(file,data){if(file===""){O(su(data))}else{}},
  show: show,
  js:   K.JS,
 }
 K.kinit(ext,kwasm)
}
function ktry(s){
 try     {let x=K._.Val(K.KC(s));K.save();return x}
 catch(e){console.log(e);K.restore();return false}
}


let kweb = {init,ktry,show,update}
export { kweb }
