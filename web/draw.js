import { K } from './k.js'

let D = {} // draw show Show


// draw[L;wh]               // draw[(`color;255*256;`rect;10 10 300 200);400 300]
// draw[L;image]            // draw over bg image
// `color;rgb               // i
// `font;"20px monospace"   // C
// `linewidth;w             // i or f
// `rect;(x;y;w;h)          // I or F        stroke
// `Rect;(x;y;w;h)          // I or F        fill
// `circle;(x;y;r)          // I or F        stroke
// `Circle;(x;y;r)          // I or F        fill
// `clip;(x;y;r)or(x;y;w;h) // I or F
// `line;(x0;y0;x1;y1)      // I or F
// `poly;(X;Y)              // X, Y: I or F  stroke
// `Poly;(X;Y)              // I or F        close fill
// `text;(x;y;"text")       // i,i,C
// `Text;(x;y;"text")       // i,i,C         rotated
// draw returns an image object (h;I)
D.draw = function(x,y){
 if(K.TK(x) != 'L'){ K.unref(x); K.unref(y); return K.KE("draw: x type") }
 
 let wh, bg
 if(K.TK(y)=='L'){
  let m = img(y);
  wh=[m.w, m.h]; bg=m.I;
 } else if((K.TK(y)=='I')&&(K.NK(y) == 2)) wh=K.IK(y);
 else {
  K.unref(x); return K.KE("draw: y type"); }
 
 let l=K.LK(x)
 let n=K.NK(x)
 if(n%2 != 0) return K.KE("draw: #x")
 
 let cnv = ce("canvas"); 
 cnv.width=wh[0]; cnv.height=wh[1];
 let ctx = cnv.getContext("2d")
 
 ctx.fillStyle = "white"
 ctx.fillRect(0,0,wh[0],wh[1])
 ctx.fillStyle = "black"
 ctx.strokeStyle = "black"
 ctx.font = "20px monospace"
 
 let state = {}
 let saveState      = function(){
  state.fillStyle   = ctx.fillStyle
  state.strokeStyle = ctx.strokeStyle
  state.font        = ctx.font
  state.lineWidth   = ctx.lineWidth
 }
 let restoreState = function(){
  ctx.fillStyle   = state.fillStyle
  ctx.strokeStyle = state.strokeStyle
  ctx.font        = state.font
  ctx.lineWidth   = state.lineWidth
 }
 saveState()
 
 // reset clip needs save/restore which also resets all other changes
 let resetClip = function(){
   ctx.restore()
   ctx.save()
   restoreState()
 }
 
 if(bg !== undefined){
  let d = ctx.createImageData(wh[0], wh[1])
  let u = new Uint8Array(d.data.buffer)
  u.set(new Uint8Array(bg.buffer))
  ctx.putImageData(d,0,0)
 }
 
 let ck = function(s,c){ if(c==false) K.KE(s+" arg") }
 
 let setcolor = function(i){ 
  let s="rgb(" + String(i&255) + ", " + String((i>>>8)&255) + ", " + String((i>>>16)&255) + ")"
  ctx.fillStyle=s; ctx.strokeStyle=s; saveState()
 }
 
 let num = function(x){ let t=K.TK(x); return (t=='f') ? K.fK(x) : (t=='i') ? K.iK(x) : K.KE("num type") }
 let vec = function(x){ let t=K.TK(x); return (t=='F') ? K.FK(x) : (t=='I') ? K.IK(x) : K.KE("vec type") }
 
 ctx.save() //for clipping
 let cmd=function(s,a){
  let n
  switch(s){
   case "color":
    setcolor(num(a))
    break
   case "font":
    ck(s, ((K.TK(a)=="L") && (K.NK(a) == 2)))
    a = K.LK(a)
    ck(s, K.TK(a[0]) == "C")
    let px = num(a[1])
    ctx.font = String(px) + "px " + K.CK(a[0]); saveState()
    break
   case "linewidth":
    ctx.lineWidth=num(a); saveState()
    break
   case "rect":
   case "Rect":
    if(s=="rect") ctx.strokeRect(...vec(a))
    else          ctx.fillRect(...vec(a))
    break
   case "circle":
   case "Circle":
    a = vec(a)
    ck(s, a.length==3)
    ctx.beginPath()
    ctx.arc(a[0], a[1], a[2], 0, 2 * Math.PI)
    if(s=="circle") ctx.stroke()
    else            ctx.fill()
    break
   case "clip":
    a = vec(a)
    resetClip()
    ctx.beginPath()
    if(a.length==3) ctx.arc(a[0], a[1], a[2], 0, 2 * Math.PI)
    else{
     ck(a.length==4)
                    ctx.rect(a[0], a[1], a[2], a[3])
    }
    ctx.clip()
    break
   case "line":
    a = vec(a)
    ck(s, a.length==4)
    ctx.beginPath()
    ctx.moveTo(a[0], a[1])
    ctx.lineTo(a[2], a[3])
    ctx.stroke()
    break
   case "poly":
   case "Poly":
    ck(s, (K.TK(a)=="L" && K.NK(a) == 2))
    a = K.LK(a)
    let x = vec(a[0]), y = vec(a[1])
    if((x.length>1)&&(x.length==y.length)){
     ctx.beginPath()
     ctx.moveTo(x[0], y[0])
     for(let i=1;i<x.length;i++) ctx.lineTo(x[i], y[i])
     if(s=="poly") ctx.stroke()
     else         {ctx.closePath(); ctx.fill()}
    }
    break
   case "text":
   case "Text":
    ck(s, ((K.TK(a)=="L") && (K.NK(a) == 3)))
    a = K.LK(a)
    let xy = vec(a[0])
    ck(s, xy.length==2)
    ck(s, K.TK(a[1]) == 'i')
    ck(s, K.TK(a[2]) == 'C')
    align(ctx, num(a[1]))
    if(s=="Text"){ ctx.save(); ctx.translate(xy[0], xy[1]); ctx.rotate(-Math.PI/2); xy[0]=0; xy[1]=0 }
    ctx.fillText(K.CK(a[2]), xy[0], xy[1])
    if(s=="Text"){ ctx.restore() }
    break
   default:
    K.KE("draw: command: "+s)
    break
   }
 }
 
 for(let i=0;i<n;i+=2){
  if(K.TK(l[i]) != 's') return K.KE("draw cmd type")
  cmd(K.CK(K.Kx("$",l[i])), l[1+i])
 }
 
 let d = ctx.getImageData(0,0,wh[0],wh[1])
 let r = K.KL([K.Ki(wh[1]), K.KI(new Int32Array(d.data.buffer))])
 cnv.remove()
 return r
}

function align(ctx, a){
 ctx.textBaseline  = ["bottom","bottom","bottom","middle","top","top","top","middle","middle"][a]
 ctx.textAlign     = ["left","center","right","right","right","center","left","left","center"][a]
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


//im:{(20;(20*50)#*1?255*256*256)};Show[im[];{`click \(x;y);im[]};{[x;y;w;h]`zoom \(x;y;w;h);im[]}]
D.Show = function(x,y,z){
 let r=D.show(x); click=y; zoom=z
 
 let cnv = ge("_cnv");
 let ctx = cnv.getContext("2d")
 cnv.ondblclick = function(ev){
  // console.log("dblclick", ev.offsetX, ev.offsetY)
  let X=ev.offsetX, Y=ev.offsetY
  ctx.beginPath();ctx.arc(X,Y,3,0,2*Math.PI);ctx.fillStyle="red";ctx.fill()
  let x=K.Kx(".", K.ref(click), K.KI([X, Y]))
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

// loadfont is ignored in js. it's only needed for the c version.
D.loadfont = function(name, bytes) { K.unref(bytes); return name }

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

export { D }
