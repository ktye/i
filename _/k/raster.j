raster:{[wh;a]
 s:`o`w`h`ix`iy`fg`pa`r!((*/wh)#0;w;h;w:wh 0;h:wh 1;0;();0)
 s:.s{@[y;1;x]}/{(x;;,y)}'[!a;.a]
 s`o}
 

strokeStyle:{[s;c]s[`fg]:icolor c;s}
beginPath:  {[s;u]s[`pa]:();s} /pa:(x0 y0;x1 y1;..)
moveTo:     {[s;u]s[`pa]:,_u;s}
lineTo:     {[s;u]s[`pa],:,_u;s}
clip:       {[s;u]s[`cl]:,/0 0-':(&/;|/)@\:s`p;s}
stroke:     {[s;u]f:iline[;;s`w;s`h];s[`o;$[0~s`r;,/1_f':s`pa;&(s`r)>+iarc[s`r;s`pa;s`w;s`h]]]:s`fg;s[`r]:0;s}
fill:       {[s;u]s[`o;&0>iarc[s`r;s`pa;s`w;s`h]]:s`fg;s[`r]:0;s}
rect:       {[s;u]dx:,/(u 2;0);dy:,/(0;u 3);o:u 0 1;s[`pa]:(o;o+dx;o+dx+dy;o+dy;o);s}
arc:        {[s;u]s[`pa`r]:_u(0 1;2);s}
fillText:   {[s;u]s[`o;istr[u 1;u 2;s`w;u 0]]:s`fg;s}
fillStyle:strokeStyle
clip:lineWidth:textBaseline:font:{y;x}


ihex:  {[c]+/i*16\*\(n:#i:|"0123456789abcdef"?-c)#16}
irgb:  {[rgb]+/1 256 65536*rgb}
icolor:{c:`white`black`gray`red`green`blue!irgb'(3#255;3#0;3#128;255 0 0;0 255 0;0 0 255);$[(#c)>i:(!c)?k:`$x;c k;ihex 1_x]}
iline: {[a;b;w;h]dx:*d:-a-b;dy:d 1;m:dy%0.+dx;$[(+dx)>+dy;x+w*y:_(a 1)+m*x:(a 0)+(1-2*dx<0)*!dx;(w*y)+x:_(a 0)+(dx%0.+dy)*y:(a 1)+(1-2*dy<0)*!dy]}
iarc:  {[r;xy;w;h]`r`xy`w`h!(r;xy;w;h);x:(w/i:!w*h)-xy 0;y:(w\i)-xy 1;((x*x)+y*y)-r*r}
istr:  {[x;y;w;s],/(10*!#s)+(x+w*y-20)+(10/i)+w*10\i:glyphs@0+s} /font.k (10x20)
\font.k

rasterT:,/(`beginPath`strokeStyle`lineWidth`rect`stroke!(();"red";2;(10 10 80 50);())
 `fillText!,("text";10;30;40)
 `beginPath`strokeStyle`arc`stroke!(();"red";80 50 10 0 2p;())
 `beginPath`strokeStyle`arc`fill!(();"red";20 50 10 0 2p;()))
rasterTest:{[](*wh) 'draster[wh:200 100;rasterT]}

rasterTest[]
