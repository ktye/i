/k\ui webserver

/js client program: sends requests to k for each event(mouse,key,size) and draws result to canvas
j:"\n"/:("w=window;d=document;b=d.body;b.style.margin=0;b.style.padding=0;b.style.overflow='hidden';N=Number;function pd(e){if(e){e.preventDefault();e.stopPropagation()}};"
 "c=d.createElement('canvas');b.appendChild(c);ctx=c.getContext('2d');"
 "function draw(s){ctx.putImageData(new ImageData(new Uint8ClampedArray(s), w.innerWidth),0,0)};" /TODO: send w in header?
 "function debounce(f,w,i){var t;return function e(){var c=this;var a=arguments;var l=function(){t=null;if(!i)f.apply(c,a);};var n=i&&!t;clearTimeout(t);t=setTimeout(l,w);if(n)f.apply(c,a);}}"
 "function get(p,f){var r = new XMLHttpRequest();r.responseType='arraybuffer';r.onreadystatechange=function(){if(this.readyState==4&&this.status == 200){if(f)f(this.response,this);}};r.open('GET',p);r.send()};"
 "function mod(e){return ','+[N(e.shiftKey),N(e.altKey),N(e.ctrlKey)]};"
 "xd=0;yd=0;down=function(e){xd=e.clientX;yd=e.clientY;pd(e)};nomenu=function(e){pd(e)};"
 "up=function(e){pd(e);get('m,'+[e.button,xd,e.clientX,yd,e.clientY]+mod(e),draw)};" /bs tab ret esc   delete     page;end;home;arrows->14-21
 "keycode=function(e){var k=e.keyCode;return (e.key.length==1)?e.key.charCodeAt():(k==8)?8:(k==9)?9:(k==13)?13:(k==27)?27:(k==46)?127:(k>32&&k<41)?k-19:null};"
 "key=function(e){var k=keycode(e);if(!k)return;get('k,'+k+mod(e),draw);pd(e);};"
 "wheel=function(e){var x=e.clientX;var y=e.clientY;var m=(e.deltaY>0)?4:(e.deltaY<0)?5:null;if(m)get('m,'+m+','+[x,y,x,y]+mod(e),draw)};"
 "size=function(e){c.width=w.innerWidth;c.height=w.innerHeight;get('s,'+[c.width,c.height],draw);pd(e)};"
 "function ae(x,y,z){x.addEventListener(y,z)};ae(w,'contextmenu',nomenu);ae(w,'mousedown',down);ae(w,'mouseup',up);ae(w,'wheel',wheel);ae(w,'keydown',key);ae(w,'resize',debounce(size,100));size()")
p:"<!DOCTYPE html>\n<html><head><link rel='icon' type='image/png' href='k.png'></head><body><script>",j,"</script></body></html>"

.Z.G:{ \x
 x:","\:*x;u:*x;a:`i$'1_x
 $[(,"/")  ~u;p
   "/k.png"~u;((,"Content-type")!,"image/png";icon)
   "/k"~u;flush uk[a 0;a 1 2 3]
   "/m"~u;flush um[a 1;a 2]
   "/s"~u;flush us[a 0;a 1]
  404]}
flush:{8_`@x+255*256*256*256}
icon:0x89504e470d0a1a0a0000000d49484452000000100000001008060000001ff3ff61000000017352474200aece1ce90000000467414d410000b18f0bfc6105000000097048597300000ec400000ec401952b0e1b0000006349444154384fad93e10a8020108377bdff3b9793a428dc16f681723fdcc71d6aed0d902ae02c2db7b35bdf177809aad9b952feefe02b91408d650523382eeb8914b830990a9230911db8308946504c05d70bd7926804259102e22456409464f13b03070b7f28230cf1c9ad0000000049454e44ae426082
