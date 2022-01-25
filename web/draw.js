import { K } from './k.js'

let D = {} // show showev draw


// draw L                 // draw[400 300;(`color;255*256;`rect;10;10;300;200)]
// `color;rgb
// `font;"20px monospace"
// `linewidth;w
// `rect;x;y;w;h
// `fillrect;x;y;w;h
// `circle;x;y;r  
// `fillcircle;x;y;r
// `line;x0;y0;x1;y1
// `poly;X;Y
// `fillpoly;X;Y
// `text;x;y;"text"
// `rtext;x;y;"text"      // rotated
// numbers may be i or f, poly/fillpoly are F
D.draw = function(x,y){
 if(K.TK(x) != 'I') return KE("draw: x type")
 if(K.TK(y) != 'L') return KE("draw: y type")
 let l=K.LK(y)
 let n=K.NK(y)
 let wh=K.IK(x)
 
 let cnv = ge("_cnv"); cnv.width=wh[0]; cnv.height=wh[1]
 let ctx = cnv.getContext("2d")
 
 let ck = function(s,c){ if(c==false) K.KE(s+" arg") }
 
 let setcolor = function(i){ 
  let s="rgb(" + i&255 + ", " + (i>>>8)&255 + ", " + (i>>>16)&255 + ")"
  ctx.fillStyle=s; ctx.strokeStyle=s
 }
 
 let nums = function(x){
  for(let i=0;i<x.length;i++){
   let xi=x[i], t = K.TK(xi)
   x[i] = (t=="i") ? K.iK(xi) : (t=="f") ? K.fK(xi) : K.KE("number-type") 
  }
  return x
 }
 let vec = function(x){ let t=K.TK(x); return (t=='F') ? K.FK(x) : (t=='I') ? K.IK(x) : KE("vec type") }
 
 let cmd=function(s,a){
  let n
  switch(s){
   case "color":
    ck("color", (a.length==1) && (K.TK(a[0]) == 'i'))
    setcolor(K.iK(i))
    break
   case "font":
    ck("font", (a.length==1) && (K.TK(a[0]) == 'C'))
    ctx.setfont( K.CK(e) )
    break
   case "linewidth":
    n = nums(a); ck("linewidth", a.length==1)
    ctx.lineWidth = n[0]
    break
   case "rect":
   case "fillrect":
    n = nums(a); ck("rect", n.length==4)
    if(s=="rect") ctx.strokeRect(...n)
    else          ctx.fillRect(...n)
    break
   case "circle":
   case "fillcircle":
    n = nums(a); ck("circ", a.length==3)
    ctx.beginPath()
    ctx.arc(n[0], n[1], n[2], 0, 2 * Math.PI)
    if(s=="circle") ctx.fill()
    else            ctx.stroke()
    break
   case "line":
    n = num(s); ck("line", n.length==4)
    ctx.beginPath()
    ctx.moveTo(n[0], n[1])
    ctx.lineTo(n[2], n[3])
    ctx.stroke()
    break
   case "poly":
   case "fillpoly":
    ck(s, a.length==2)
    let x = vec(a[0]), y = vec(a[1])
    if((x.length>1)&&(x.length==y.length)){
     ctx.moveTo(x[0], y[0])
     for(let i=1;i<x.length;i++) ctx.lineTo(x[i], y[i])
     if(s=="poly") ctx.stroke()
     else         {ctx.closePath(); ctx.fill()}
    }
    break
   case "text":
   case "rtext":
    ck(s, (a.length==3)&&(K.TK(a[2])=='C'))
    n = nums(a.slice(0,2))
    let nx=n[1], ny=n[2]
    if(s=="rtext"){ ctx.save(); ctx.translate(nx,ny); ctx.rotate(-Math.PI/2); nx=0; ny=0 }
    ctx.fillText(K.CK(a[2]), nx, ny)
    if(s=="rtext"){ ctx.restore() }
    break
   default:
    K.KE("draw: command: "+s)
    break
   }
 }
 
 let s="", c=[]
 for(let i=0;i<n;i++){
  let e=K.TK(l[i])
  if(K.TK(e) == 's'){
   if(s != ""){
    cmd(s, c)
    c = []
    s = K.sK(e)
   }
  } else c.push(e)
 }
 if(s != "") cmd(s, c)
}


// display an image on a canvas with id="_cnv"
// show(20;(20*50)#255)
D.show = function(x){
 unhandle()
 if(K.TK(x)=='i') return x  // hide image: show 0
 return setimg(img(x))
}

function setimg(x){
 let cnv = ge("_cnv"); cnv.width=x.w; cnv.height=x.h
 let ctx = cnv.getContext("2d")
 let d = ctx.createImageData(x.w, x.h)
 let u = new Uint8Array(d.data.buffer)
 u.set(new Uint8Array(x.I.buffer))
 ctx.putImageData(d,0,0)
 cnv.style.display = "block"
 return K.Ki(x.w*x.h)
}

let click, zoom // k-values: callback functions

function unhandle(){
 if(click!=undefined){K.unref(click); click=undefined}
 if(zoom !=undefined){K.unref(zoom);  zoom =undefined}
 let cnv = ge("_cnv")
 cnv.ondblclick  = undefined
 cnv.onmousedown = undefined
 cnv.onmouseup   = undefined
 cnv.style.display = "none"
}


//im:{(20;(20*50)#*1?255*256*256)};showev[im[];{`click \(x;y);im[]};{[x;y;w;h]`zoom \(x;y;w;h);im[]}]
D.showev = function(x,y,z){
 let r=D.show(x); click=y; zoom=z
 
 let cnv = ge("_cnv");
 let ctx = cnv.getContext("2d")
 console.log("set ondblclick")
 cnv.ondblclick = function(ev){
  // console.log("dblclick", ev.offsetX, ev.offsetY)
  let x=K.Kx(".", K.ref(click), K.KI([ev.offsetX, ev.offsetY]))
  K.unref( (K.TK(x)==="L") ? setimg(img(x)) : x )
 }
 
 let zd = false, zm, zs, ze, bg
 zs = function(ev){ //zoom-start
  bg = ctx.getImageData(0,0,cnv.width,cnv.height)
  cnv.onmousemove = zm
  cnv.style.cursor="crosshair"
  zd = [ev.offsetX,ev.offsetY,0,0]
 }
 zm = function(ev){ //zoom-move
  if(zd!==false){
   zd=[zd[0],zd[1],ev.offsetX-zd[0],ev.offsetY-zd[1]]
   ctx.putImageData(bg,0,0)
   ctx.beginPath()
   ctx.rect(...zd)
   ctx.strokeStyle='red'
   ctx.stroke()
  }
 }
 ze = function(ev){ //zoom-end
  ctx.putImageData(bg,0,0)
  cnv.style.cursor=""
  cnv.onmousemove = undefined
  if(Math.abs(zd[2])< 5||Math.abs(zd[3])<5){zd=false;return}
  if(zd[2]<0){ zd[0]+=zd[2]; zd[2]=-zd[2] }
  if(zd[3]<0){ zd[1]+=zd[3]; zd[3]=-zd[3] }
  
  console.log("zoom", zd)
  
  let x=K.Kx(".", K.ref(zoom), K.KI(zd))
  K.unref( (K.TK(x)==="L") ? setimg(img(x)) : x )
  
  zd = false
  return
 }
 cnv.onmousedown = zs
 cnv.onmousemove = zm
 cnv.onmouseup   = ze
 
 return r
}

function img(x){
 if(K.TK(x) != 'L') return K.KE("img: type")
 if(K.NK(x) != 2)   return K.KE("img: L2")
 let r = K.LK(x)
 if(K.TK(r[0]) != 'i') return K.KE("img: h-type")
 if(K.TK(r[1]) != 'I') return K.KE("img: I-type")
 let h = K.iK(r[0])
 let I = K.IK(r[1])
 let w = Math.floor(I.length / h)
 if(I.length != w*h) return K.KE("img:rect")
 let u = new Uint8Array(I.buffer)
 for(let i=3;i<u.length;i+=4)u[i]=255 //alpha
 return {w:w, h:h, I:I}
}

function ge(x){return document.getElementById(x)}
function ce(x){return document.createElement(x)}

// cnv=document.getElementById("_cnv");ctx=cnv.getContext("2d")



export { D }
