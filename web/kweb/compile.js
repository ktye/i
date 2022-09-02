import {kweb} from './kweb.js'
import {K} from '../k.js'

function ce(x) { return document.createElement(x) }
function ge(x) { return document.getElementById(x) }


function treeclick(e){
 console.log(e)
}

kweb.register("treetable",function(dst,x){
 let l=K.LK(x)
 let D=K.LK(l[0]),T=K.SK(l[1]),P=K.IK(l[2]),I=K.IK(l[3]),S=K.SK(l[4])
 let a=D.map(x=>K.CK(x).replaceAll("|","│").replace("+","├").replace("L","└").replace("-","─").replace("T","┬"))
 let p=ce("pre")
 for(let i=0;i<a.length;i++){
  p.textContent += a[i]+" "+T[i].padEnd(8)+" "+String(P[i]).padEnd(8)+(I[i]==-2147483648?"0N":String(I[i])).padEnd(8)+String(S[i])+"\n"
 }
 p.ondblclick=treeclick
 dst.appendChild(p)
})

function setinput(s){ 
 kweb.init({start:[s+".k","go.k","compile.k"]})
}

let sel=ge("sel")
let examples="asn`blank`cal`cast`cli`cnd`const`cont`drp`fun`heap`label`lit`loop`loop2`mem`swtch`tab`k".split("`")
for(let i=0;i<examples.length;i++){
 let o=ce("option");o.textContent=examples[i]
 sel.appendChild(o)
}
sel.selectedIndex=0
sel.onchange=function(e){setinput(examples[e.target.selectedIndex])}
setinput(examples[sel.selectedIndex])

ge("compile").onclick=function(){
 ge("out").textContent=K.CK(kweb.ktry("emit go``"))
}
fetch("compile.help").then(r=>r.text()).then(r=>ge("out").textContent=r)
