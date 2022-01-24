import { K } from './k.js'

let D = {} // png rgb draw show

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
 return {w:w, h:h, I:I}
}
function urlpng(x){
 var u = K.BK(x)
 return URL.createObjectURL(new Blob([u.buffer], {type: "image/png"}))
 //for(let i=0;i<u.length;i++) s+=String.fromCharCode(u[i])
 //s = "data:data:image/png;base64," + window.btoa(s)
}



// png(20;(20*30)#255)
D.png = function(x){
 let im = img(x)
 let cnv = ce("canvas")
 let ctx = cnv.getContext("2d")
 let d = ctx.createImageData(im.w, im.h)
 let m = new Uint8ClampedArray(im.I.buffer)
 for(let i=3;i<m.length;i+=4) m[i]=255
 d.data.set(m)
 
 let s = cnv.toDataURL("image/png")
 let u = b46(s.slice(7+s.indexOf("base64,")))
 //let u = new Uint8Array(s.length)
 //for(let i=0;i<u.length;i++) u[i]=s.charCodeAt(i);
 cnv.remove()
 return K.KC(u)
}

D.show = function(x){
 var cnv = ge("_cnv")
 if(K.TK(x)==="C"){
  var im = new Image()
  im.src = "data:;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAgMAAABinRfyAAAACVBMVEX/AAAAAAD////KksOZAAAAMElEQVR4nGJYtWrVKoYFq1ZxMSyYhkZMgxNRXAwLpmbBCDAXSRZEgAwAGQUIAAD//+QzHr+8V1EyAAAAAElFTkSuQmCC"
  im.onload = function(){
   console.log("drawImage..")
   cnv.width = im.width
   cnv.height = im.height
   let ctx = cnv.getContext("2d")
   ctx.drawImage(im, 0, 0)
  }
  return K.Ki(0)
 }else{return KE("show type")}
}

function ge(x){return document.getElementById(x)}
function ce(x){return document.createElement(x)}



// https://gist.github.com/enepomnyaschih/72c423f727d395eeaa09697058238727 (mit)
const _b64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_".split("")
const _b46 = [255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,62,255,255,52,53,54,55,56,57,58,59,60,61,255,255,255,0,255,255,255,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,255,255,255,255,63,255,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255]
function b64(u) {
	let r = '', i, l = u.length;
	for (i=2; i<l; i+= 3){
		r += _b64[u[i - 2] >> 2];
		r += _b64[((u[i - 2] & 0x03) << 4) | (u[i - 1] >> 4)];
		r += _b64[((u[i - 1] & 0x0F) << 2) | (u[i] >> 6)];
		r += _b64[u[i] & 0x3F];
	}
	if (i === l + 1){
		r += _b64[u[i - 2] >> 2];
		r += _b64[(u[i - 2] & 0x03) << 4];
		r += "==";
	}
	if (i === l){
		r += _b64[u[i - 2] >> 2];
		r += _b64[((u[i - 2] & 0x03) << 4) | (u[i - 1] >> 4)];
		r += _b64[(u[i - 1] & 0x0F) << 2];
		r += "=";
	}
	return r
}
function b46(s) {
	const index = s.indexOf("=");
	let m = s.endsWith("==") ? 2 : s.endsWith("=") ? 1 : 0, n=s.length, r=new Uint8Array(3 * (n / 4)), b
	for (let i = 0, j = 0; i < n; i += 4, j += 3) {
		b = _b46[s.charCodeAt(i)] << 18 | _b46[s.charCodeAt(i + 1)] << 12 | _b46[s.charCodeAt(i + 2)] << 6 | _b46[s.charCodeAt(i + 3)];
		r[j] = b >> 16;
		r[j + 1] = (b >> 8) & 0xFF;
		r[j + 2] = b & 0xFF;
	}
	return r.subarray(0, r.length-m);
}




export { D }
